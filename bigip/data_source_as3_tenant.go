/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceBigipAs3Tenant() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipAs3TenantRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of Tenant",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Tenant handles traffic only when enabled (default)",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional friendly name for this object",
			},
			"optimisticlockkey": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Arbitrary (brief) text pertaining to this object (optional)",
			},
			"app_class_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Reference to  list of Application class objects",
			},
			"defaultroutedomain": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Name of application",
			},
			"tenant_map": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceBigipAs3TenantRead(d *schema.ResourceData, meta interface{}) error {
	map1 := make(map[string]interface{})
	var result map[string]interface{}
	map1["class"] = "Tenant"
	if m, ok := d.GetOk("app_class_list"); ok {
		for _, vs := range m.([]interface{}) {
			err := json.Unmarshal([]byte(vs.(string)), &result)
			if err != nil {
				log.Printf("Json Unmarshall failed with: %v", err)
				return err
			}
			for k, v := range result {
				var result1 map[string]interface{}
				err := json.Unmarshal([]byte(v.(string)), &result1)
				if err != nil {
					log.Printf("Json Unmarshall failed with: %v", err)
					return err
				}
				map1[k] = result1
				//log.Printf("result application map in tenant :%+v ", map1)
			}
		}
	}
	out, err := json.Marshal(map1)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	bigip.As3Tenant.TenantList = append(bigip.As3Tenant.TenantList, name)
	log.Printf("Tenant name:%+v", bigip.As3Tenant.TenantList)
	resultMap := make(map[string]interface{})
	resultMap[name] = string(out)
	log.Printf("resultMap in Tenant Class:%v", resultMap)
	d.Set("tenant_map", resultMap)
	d.SetId(name)
	return nil
}
