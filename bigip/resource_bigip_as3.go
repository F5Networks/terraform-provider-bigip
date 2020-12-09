/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"log"
	"strings"
	"sync"
)

var x = 0
var m sync.Mutex

func resourceBigipAs3() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipAs3Create,
		Read:   resourceBigipAs3Read,
		Update: resourceBigipAs3Update,
		Delete: resourceBigipAs3Delete,
		Exists: resourceBigipAs3Exists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// d.Id() here is the last argument passed to the `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
				// Here we use a function to parse the import ID (like the example above) to simplify our logic
				//if err != nil {
				//    return nil, err
				//}
				_ = d.Set("tenant_list", d.Id())
				_ = d.Set("tenant_filter", d.Id())

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"as3_json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "AS3 json",
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)

					return json
				},
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if _, err := structure.NormalizeJsonString(v); err != nil {
						errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
					}
					as3json := v.(string)
					resp := []byte(as3json)
					jsonRef := make(map[string]interface{})
					json.Unmarshal(resp, &jsonRef)
					for key, value := range jsonRef {
						if key == "class" && value != "AS3" {
							errors = append(errors, fmt.Errorf("Json must have AS3 class"))
						}
						if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
							for k, v := range rec {
								if k == "class" && v != "ADC" {
									errors = append(errors, fmt.Errorf("Json must have ADC class"))
								}
							}
						}
					}
					return
				},
			},
			"tenant_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "this attribute is no longer in use",
				Description: "Name of Tenant",
			},
			"tenant_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of Tenant",
			},
			"tenant_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of Tenant",
			},
			"application_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of Application",
			},
		},
	}
}

func resourceBigipAs3Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	as3Json := d.Get("as3_json").(string)
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Creating As3 config")
	tenantFilter := d.Get("tenant_filter").(string)
	tenantList, _, applicationList := client.GetTenantList(as3Json)
	tenantCount := strings.Split(tenantList, ",")
	if tenantFilter != "" {
		tenantList = tenantFilter
	}
	_ = d.Set("tenant_list", tenantList)
	d.Set("application_list", applicationList)
	strTrimSpace, err := client.AddTeemAgent(as3Json)
	if err != nil {
		return err
	}
	//log.Printf("[INFO] Tenants in Json:%+v", tenantList)
	log.Printf("[INFO] Creating as3 config in bigip:%s", strTrimSpace)
	err, successfulTenants := client.PostAs3Bigip(strTrimSpace, tenantList)
	if err != nil {
		if successfulTenants == "" {
			return fmt.Errorf("posting as3 config failed for tenants:(%s) with error: %v", tenantList, err)
		}
		_ = d.Set("tenant_list", successfulTenants)
		if len(successfulTenants) != len(tenantList) {
			log.Printf("%v", err)
		}
	}
	if !client.Teem {
		id := uuid.New()
		uniqueID := id.String()
		assetInfo := f5teem.AssetInfo{
			"Terraform-provider-bigip",
			client.UserAgent,
			uniqueID,
		}
		teemDevice := f5teem.AnonymousClient(assetInfo, "")
		f := map[string]interface{}{
			"Number_of_tenants": len(tenantCount),
			"Terraform Version": client.UserAgent,
		}
		err = teemDevice.Report(f, "bigip_as3", "1")
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	d.SetId(tenantList)
	x = x + 1
	//m.Unlock()
	return resourceBigipAs3Read(d, meta)
}
func resourceBigipAs3Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading As3 config")
	name := d.Get("tenant_list").(string)
	applicationList := d.Get("application_list").(string)
	log.Printf("[DEBUG] Tenants in AS3 get call : %s", name)
	as3Resp, err := client.GetAs3(name, applicationList)
	log.Printf("[DEBUG] AS3 json retreived from the GET call in Read function : %s", as3Resp)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return nil
		}
		return err
	}
	if as3Resp == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("as3_json", as3Resp)
	_ = d.Set("tenant_list", name)
	return nil
}

func resourceBigipAs3Exists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Checking if As3 config exists in BIGIP")
	name := d.Get("tenant_list").(string)
	applicationList := d.Get("application_list").(string)
	tenantFilter := d.Get("tenant_filter").(string)
	if tenantFilter != "" {
		name = tenantFilter
	}
	as3Resp, err := client.GetAs3(name, applicationList)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return false, nil
		}
		return false, err
	}
	log.Printf("[INFO] AS3 response Body:%+v", as3Resp)
	if as3Resp == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		return false, nil
	}
	return true, nil
}

func resourceBigipAs3Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	as3Json := d.Get("as3_json").(string)
	//m.Lock()
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Updating As3 Config :%s", as3Json)
	name := d.Get("tenant_list").(string)
	tenantList, _, _ := client.GetTenantList(as3Json)
	tenantFilter := d.Get("tenant_filter").(string)
	if tenantFilter == "" {
		if tenantList != name {
			_ = d.Set("tenant_list", tenantList)
			newList := strings.Split(tenantList, ",")
			oldList := strings.Split(name, ",")
			deletedTenants := client.TenantDifference(oldList, newList)
			if deletedTenants != "" {
				err, _ := client.DeleteAs3Bigip(deletedTenants)
				if err != nil {
					log.Printf("[ERROR] Unable to Delete removed tenants: %v :", err)
					return err
				}
			}
		}
	} else {
		tenantList = tenantFilter
	}
	strTrimSpace, err := client.AddTeemAgent(as3Json)
	if err != nil {
		return err
	}
	err, successfulTenants := client.PostAs3Bigip(strTrimSpace, tenantList)
	if err != nil {
		if successfulTenants == "" {
			return fmt.Errorf("Error updating json  %s: %v", tenantList, err)
		}
		_ = d.Set("tenant_list", successfulTenants)
		if len(successfulTenants) != len(tenantList) {
			log.Printf("%v", err)
		}
	}
	x = x + 1
	//m.Unlock()
	return resourceBigipAs3Read(d, meta)
}

func resourceBigipAs3Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	//as3Json := d.Get("as3_json").(string)
	m.Lock()
	defer m.Unlock()
	//m.Lock()
	log.Printf("[INFO] Deleting As3 config")
	name := d.Get("tenant_list").(string)
	err, failedTenants := client.DeleteAs3Bigip(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete: %v :", err)
		return err
	}
	if failedTenants != "" {
		_ = d.Set("tenant_list", name)
		return resourceBigipAs3Read(d, meta)
	}
	x = x + 1
	//m.Unlock()
	d.SetId("")
	return nil
}
