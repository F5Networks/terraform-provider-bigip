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

func dataSourceBigipAs3Adc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipAs3AdcRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of Adc",
			},
			"identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier for this declaration (max 255 printable chars with no spaces, quotation marks, angle brackets, nor backslashes)",
			},
			"common": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"remark": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Arbitrary (brief) text pertaining to this object (optional)",
						},
					},
				},
			},
			"constants": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Version number of declaration; update when you change contents but not ID (optional but recommended)",
						},
						"timestamp": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Date+time (this version of) declaration was created (optional but recommended)",
						},
					},
				},
			},
			"controls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"archive_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"archive_timestamp": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"log_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "error",
							Description: "",
						},
						"trace": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "",
						},
					},
				},
			},

			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional friendly name for this object",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Arbitrary (brief) text pertaining to this object (optional)",
			},
			"schema_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "3.15.0",
				Description: "Version of ADC Declaration schema this declaration uses",
			},
			"scratch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Holds some system data during declaration processing",
			},
			"target": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP address of managed device to be configured",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Host name of managed device to be configured",
						},
						"ssg_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of shared service group to be configured",
						},
					},
				},
			},
			"update_mode": {
				Type:     schema.TypeString,
				Optional: true,
				//       Default:     "selective",
				Description: "When set to ‘selective’ (default) AS3 does not modify Tenants not referenced in the declaration. Otherwise (‘complete’) AS3 removes unreferenced Tenants.",
			},
			"tenant_class_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenant_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of tenants to be deleted",
			},
		},
	}
}

func dataSourceBigipAs3AdcRead(d *schema.ResourceData, meta interface{}) error {
	map1 := make(map[string]interface{})
	var result map[string]interface{}
	map1["class"] = "ADC"
	map1["schemaVersion"] = d.Get("schema_version").(string)
	if m, ok := d.GetOk("tenant_class_list"); ok {
		for _, kv := range m.([]interface{}) {
			err := json.Unmarshal([]byte(kv.(string)), &result)
			if err != nil {
				log.Printf("Json Unmarshall failed with: %v", err)
				return err
			}
			var tenantResult map[string]interface{}
			for k, v := range result {
				err := json.Unmarshal([]byte(v.(string)), &tenantResult)
				if err != nil {
					log.Printf("Json Unmarshall failed with: %v", err)
					return err
				}
				map1[k] = tenantResult
			}
		}
	}
	out, err := json.Marshal(map1)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	resultMap := make(map[string]interface{})
	resultMap[name] = string(out)
	out, err = json.Marshal(resultMap)
	if err != nil {
		return err
	}
	log.Printf("ADC class string:%+v\n", string(out))
	d.SetId(string(out))
	d.Set("tenant_list", bigip.As3Tenant.TenantList)
	return nil
}

/*
func getTenantList() string {
	log.Printf("Tenant name from getTenantList :%+v", teList)
	return teList
}
*/
