/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"
)

func resourceBigipAs3() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipAs3Create,
		Read:   resourceBigipAs3Read,
		Update: resourceBigipAs3Update,
		Delete: resourceBigipAs3Delete,
		Exists: resourceBigipAs3Exists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				ValidateFunc: validation.ValidateJsonString,
			},
                        "tenant_name": {
                               Type:        schema.TypeString,
                                Optional:    true, 
                                Deprecated:    "this attribute is no longer in use",
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
		},
	}
}

func resourceBigipAs3Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	as3Json := d.Get("as3_json").(string)
        tenantFilter := d.Get("tenant_filter").(string)
	if ok := bigip.ValidateAS3Template(as3Json); !ok {
		return fmt.Errorf("[AS3] Error validating template \n")
	}
	//strTrimSpace := strings.TrimSpace(as3Json)
	tenantList, _ := client.GetTenantList(as3Json)
        if tenantFilter != "" {
           tenantList = tenantFilter
        }
	d.Set("tenant_list", tenantList)
	strTrimSpace := client.AddTeemAgent(as3Json)
	log.Printf("[INFO] Tenants in Json:%+v", tenantList)
	log.Printf("[INFO] Creating as3 config in bigip:%s", strTrimSpace)
	err := client.PostAs3Bigip(strTrimSpace, tenantList)
	if err != nil {
		client.DeleteAs3Bigip(tenantList)
		return fmt.Errorf("Error creating json  %s: %v", tenantList, err)
	}
	d.SetId(tenantList)
	return resourceBigipAs3Read(d, meta)
}
func resourceBigipAs3Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading As3 config")
	name := d.Get("tenant_list").(string)
	as3Resp, err := client.GetAs3(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		return err
	}
	if as3Resp == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("as3_json", as3Resp)
	d.Set("tenant_list", name)
	return nil
}

func resourceBigipAs3Exists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Checking if As3 config exists in BIGIP")
	name := d.Get("tenant_list").(string)
        tenantFilter := d.Get("tenant_filter").(string)
        if tenantFilter!= "" {
           name = tenantFilter
        }
	as3Resp, err := client.GetAs3(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
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
	log.Printf("[INFO] Updating As3 Config :%s", as3Json)
	name := d.Get("tenant_list").(string)
	tenantList, _ := client.GetTenantList(as3Json)
        tenantFilter := d.Get("tenant_filter").(string)
        if tenantFilter == "" {
	if tenantList != name {
		d.Set("tenant_list", tenantList)
		new_list := strings.Split(tenantList, ",")
		old_list := strings.Split(name, ",")
		deleted_tenants := client.TenantDifference(old_list, new_list)
		if deleted_tenants != "" {
			err := client.DeleteAs3Bigip(deleted_tenants)
			if err != nil {
				log.Printf("[ERROR] Unable to Delete removed tenants: %v :", err)
				return err
			}
		}
	}
        } else {
           tenantList = tenantFilter
        }
	err := client.ModifyAs3(tenantList, as3Json)

	if err != nil {
		return fmt.Errorf("Error modifying json %s: %v", name, err)
	}
	return resourceBigipAs3Read(d, meta)
}

func resourceBigipAs3Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Deleting As3 config")
	name := d.Get("tenant_list").(string)
	err := client.DeleteAs3Bigip(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete: %v :", err)
		return err
	}
	d.SetId("")
	return nil
}
