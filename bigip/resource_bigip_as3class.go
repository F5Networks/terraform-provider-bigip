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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
)

func resourceBigipAs3Class() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipAs3ClassCreate,
		Read:   resourceBigipAs3ClassRead,
		Update: resourceBigipAs3ClassUpdate,
		Delete: resourceBigipAs3ClassDelete,
		//Exists: resourceAs3Exists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of AS3",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier for this declaration (max 255 printable chars with no spaces, quotation marks, angle brackets, nor backslashes)",
			},
			"persist": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Unique identifier for this declaration (max 255 printable chars with no spaces, quotation marks, angle brackets, nor backslashes)",
			},
			"declaration": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Declaration string",
			},
			"tenants": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of tenants to be deleted",
			},
		},
	}
}

func resourceBigipAs3ClassCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var as3Main = &bigip.As3Main{}
	as3Main.Class = "AS3"
	as3Main.Action = "deploy"
	as3Main.Persist = true

	var result map[string]interface{}
	var adcResult map[string]interface{}
	if kv, ok := d.GetOk("declaration"); ok {
		err := json.Unmarshal([]byte(kv.(string)), &adcResult)
		if err != nil {
			log.Printf("Json Unmarshall failed with: %v", err)
			return err
		}
		for _, v := range adcResult {
			err := json.Unmarshal([]byte(v.(string)), &result)
			if err != nil {
				log.Printf("Json Unmarshall failed with: %v", err)
				return err
			}
		}
	}
	as3Main.Declaration = result
	name := d.Get("name").(string)
	log.Printf("AS3 JSON to be Posted:%+v", as3Main)
	err := client.PostAs3Bigip(as3Main)
	if err != nil {
		log.Printf("[ERROR] Unable to Post: %v :", err)
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipAs3ClassRead(d *schema.ResourceData, meta interface{}) error {
	var result map[string]interface{}
	var adcResult map[string]interface{}
	if kv, ok := d.GetOk("declaration"); ok {
		err := json.Unmarshal([]byte(kv.(string)), &adcResult)
		if err != nil {
			log.Printf("Json Unmarshall failed with: %v", err)
			return err
		}
		for _, v := range adcResult {
			err := json.Unmarshal([]byte(v.(string)), &result)
			if err != nil {
				log.Printf("Json Unmarshall failed with: %v", err)
				return err
			}
		}
	}
	log.Printf("Result in Read:%+v", result)
	return nil
}

func resourceBigipAs3ClassUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var as3Main = &bigip.As3Main{}
	as3Main.Class = "AS3"
	as3Main.Action = "deploy"
	as3Main.Persist = true

	var result map[string]interface{}
	var adcResult map[string]interface{}
	if kv, ok := d.GetOk("declaration"); ok {
		err := json.Unmarshal([]byte(kv.(string)), &adcResult)
		if err != nil {
			log.Printf("Json Unmarshall failed with: %v", err)
			return err
		}
		for _, v := range adcResult {
			err := json.Unmarshal([]byte(v.(string)), &result)
			if err != nil {
				log.Printf("Json Unmarshall failed with: %v", err)
				return err
			}
		}
	}
	as3Main.Declaration = result
	log.Printf("[DEBUG] AS3 JSON to be Posted:%+v", as3Main)
	err := client.PostAs3Bigip(as3Main)
	if err != nil {
		log.Printf("[ERROR] Unable to Post: %v :", err)
		return err
	}
	return nil
}

func resourceBigipAs3ClassDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var tenantList []string
	if m, ok := d.GetOk("tenants"); ok {
		for _, vs := range m.([]interface{}) {
			tenantList = append(tenantList, vs.(string))
		}
	} else {
		return fmt.Errorf("[ERROR] Please specify list of tenants to be deleted using tenant_list")
	}
	result := strings.Join(tenantList, ",")
	log.Printf("[INFO] Deleting AS3 Partition:%v", result)
	err := client.DeleteAs3Bigip(result)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete: %v :", err)
		return err
	}
	return nil
}
