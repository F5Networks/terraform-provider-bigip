/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceBigipAs3App() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipAs3AppRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of app",
			},
			"template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of template",
				//ForceNew:    true,
			},
			"servicemain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of application",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Tenant handles traffic only when enabled (default)",
			},
			"schema_overlay": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IQ name for a supplemental validation schema is applied to the Application class definition before the main AS3 schema",
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
			"pool_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional friendly name for this object",
			},
			"service_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "this will be reference to service class object",
			},
			"cert_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "this will be reference to Certificate class object",
			},
			"tls_server_class": {
				Type:     schema.TypeString,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
func dataSourceBigipAs3AppRead(d *schema.ResourceData, meta interface{}) error {
	map1 := make(map[string]interface{})
	var serviceTypeRead string
	map1["class"] = "Application"
	map1["template"] = d.Get("template").(string)
	var serviceResult map[string]interface{}
	var certresult map[string]interface{}
	var tlsresult map[string]interface{}
	var poolResult map[string]interface{}
	if kv, ok := d.GetOk("service_class"); ok {
		err := json.Unmarshal([]byte(kv.(string)), &serviceResult)
		if err != nil {
			log.Printf("Json Unmarshall failed with: %v", err)
			return err
		}
		var serviceResult1 map[string]interface{}
		for k, v := range serviceResult {
			if k != "service_type" {
				err := json.Unmarshal([]byte(v.(string)), &serviceResult1)
				if err != nil {
					log.Printf("Json Unmarshall failed with: %v", err)
					return err
				}
				map1[k] = serviceResult1
			} else {
				serviceTypeRead = v.(string)
				if map1["template"] != serviceTypeRead && map1["template"] != "shared" && map1["template"] != "generic" {
					return errors.New(fmt.Sprintf("Incorrect Template Type"))
				}
			}
		}
	}
	if kv, ok := d.GetOk("pool_class"); ok {
		err := json.Unmarshal([]byte(kv.(string)), &poolResult)
		if err != nil {
			log.Printf("Json Unmarshall failed with: %v", err)
			return err
		}
		var poolResult1 map[string]interface{}
		for k, v := range poolResult {

			err := json.Unmarshal([]byte(v.(string)), &poolResult1)
			if err != nil {
				log.Printf("Json Unmarshall failed with: %v", err)
				return err
			}
			map1[k] = poolResult1
		}
	}
	if kv, ok := d.GetOk("tls_server_class"); ok {
		err := json.Unmarshal([]byte(kv.(string)), &tlsresult)
		if err != nil {
			log.Printf("Json Unmarshall failed with: %v", err)
			return err
		}
		var tlsresult1 map[string]interface{}
		for k, v := range tlsresult {
			err := json.Unmarshal([]byte(v.(string)), &tlsresult1)
			if err != nil {
				log.Printf("Json Unmarshall failed with: %v", err)
				return err
			}
			map1[k] = tlsresult1
		}
	}
	if kv, ok := d.GetOk("cert_class"); ok {
		err := json.Unmarshal([]byte(kv.(string)), &certresult)
		if err != nil {
			log.Printf("Json Unmarshall failed with: %v", err)
			return err
		}
		var certresult1 map[string]interface{}
		for k, v := range certresult {
			err := json.Unmarshal([]byte(v.(string)), &certresult1)
			if err != nil {
				log.Printf("Json Unmarshall failed with: %v", err)
				return err
			}
			map1[k] = certresult1
		}
	}
	out, err := json.Marshal(map1)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	resultMap := make(map[string]interface{})
	resultMap[name] = string(out)
	out1, err := json.Marshal(resultMap)
	if err != nil {
		return err
	}
	log.Printf("Application Class string :%+v", string(out1))
	d.SetId(string(out1))
	return nil
}
