package bigip

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var p = 0
var q sync.Mutex

func resourceBigiqAs3() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigiqAs3Create,
		Read:   resourceBigiqAs3Read,
		Update: resourceBigiqAs3Update,
		Delete: resourceBigiqAs3Delete,
		Exists: resourceBigiqAs3Exists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bigiq_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key pool to use",
			},
			"bigiq_user": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_token_auth": {
				Type:      schema.TypeBool,
				Optional:  true,
				Sensitive: true,
				//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//	//log.Printf("Value of k=%v,old=%v,new%v", k, old, new)
				//	if old != new {
				//		return true
				//	}
				//	return false
				//},
				Description: "Enable to use an external authentication source (LDAP, TACACS, etc)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_TOKEN_AUTH", true),
			},
			"bigiq_login_ref": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//	//log.Printf("Value of k=%v,old=%v,new%v", k, old, new)
				//	if old != new {
				//		return true
				//	}
				//	return false
				//},
				Description: "Login reference for token authentication (see BIG-IQ REST docs for details)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_LOGIN_REF", "local"),
			},
			"as3_json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "AS3 json",
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				ValidateFunc: validation.ValidateJsonString,
			},
			"tenant_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of Tenant",
			},
		},
	}
}

func resourceBigiqAs3Create(d *schema.ResourceData, meta interface{}) error {
	//bigipRef := meta.(*bigip.BigIP)
	//log.Println(bigipRef)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	q.Lock()
	defer q.Unlock()
	as3Json := d.Get("as3_json").(string)
	tenantList, _, _ := bigiqRef.GetTenantList(as3Json)
	targetInfo := bigiqRef.GetTarget(as3Json)
	_ = d.Set("tenant_list", tenantList)
	err, successfulTenants := bigiqRef.PostAs3Bigiq(as3Json)
	if err != nil {
		if successfulTenants == "" {
			return fmt.Errorf("Error creating json  %s: %v", tenantList, err)
		}
		_ = d.Set("tenant_list", successfulTenants)
	}
	as3ID := fmt.Sprintf("%s_%s", targetInfo, successfulTenants)
	d.SetId(as3ID)
	p = p + 1
	return resourceBigiqAs3Read(d, meta)
}

func resourceBigiqAs3Read(d *schema.ResourceData, meta interface{}) error {
	time.Sleep(20 * time.Second)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	tenantRef := d.Id()
	log.Println("[INFO] Reading As3 config")
	targetRef := strings.Split(tenantRef, "_")[0]
	name := strings.Split(tenantRef, "_")[1]
	if name != d.Get("tenant_list").(string) {
		as3Resp, err := bigiqRef.GetAs3Bigiq(targetRef, d.Get("tenant_list").(string))
		if err != nil {
			log.Printf("[ERROR] Unable to retrieve json ")
			return err
		}
		if as3Resp == "" {
			log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		_ = d.Set("as3_json", as3Resp)
	} else {
		as3Resp, err := bigiqRef.GetAs3Bigiq(targetRef, name)
		if err != nil {
			log.Printf("[ERROR] Unable to retrieve json ")
			return err
		}
		if as3Resp == "" {
			log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		_ = d.Set("as3_json", as3Resp)

	}
	return nil
}

func resourceBigiqAs3Exists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func resourceBigiqAs3Update(d *schema.ResourceData, meta interface{}) error {
	time.Sleep(20 * time.Second)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	as3Json := d.Get("as3_json").(string)
	q.Lock()
	defer q.Unlock()
	log.Printf("[INFO] Updating As3 Config :%s", as3Json)
	name := d.Get("tenant_list").(string)
	tenantList, _, _ := bigiqRef.GetTenantList(as3Json)
	if tenantList != name {
		_ = d.Set("tenant_list", tenantList)
		newList := strings.Split(tenantList, ",")
		oldList := strings.Split(name, ",")
		deletedTenants := bigiqRef.TenantDifference(oldList, newList)
		if deletedTenants != "" {
			//err, _ := bigiqRef.DeleteAs3Bigip(deleted_tenants)
			//if err != nil {
			//	log.Printf("[ERROR] Unable to Delete removed tenants: %v :", err)
			//	return err
			//}
		}
	}
	err, successfulTenants := bigiqRef.PostAs3Bigiq(as3Json)
	if err != nil {
		if successfulTenants == "" {
			return fmt.Errorf("Error creating json  %s: %v", tenantList, err)
		}
		_ = d.Set("tenant_list", successfulTenants)
	}
	p = p + 1
	return resourceBigiqAs3Read(d, meta)
}

func resourceBigiqAs3Delete(d *schema.ResourceData, meta interface{}) error {
	time.Sleep(20 * time.Second)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return err
	}
	q.Lock()
	defer q.Unlock()
	log.Printf("[INFO] Deleting As3 config")
	name := d.Get("tenant_list").(string)
	as3Json := d.Get("as3_json").(string)
	err, failedTenants := bigiqRef.DeleteAs3Bigiq(as3Json, name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete: %v :", err)
		return err
	}
	if failedTenants != "" {
		_ = d.Set("tenant_list", name)
		return resourceBigipAs3Read(d, meta)
	}
	p = p + 1
	//m.Unlock()
	d.SetId("")
	return nil
}
