/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
)

func dataSourceBigipWafEntityUrl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipWafEntityUrlRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the URL",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the URL.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Specifies whether the parameter is an 'explicit' or a 'wildcard' attribute. " +
					"Default is: wildcard",
				Default: "wildcard",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the protocol for the URL is 'http' or 'https'. Default is: http",
				Default:     "http",
			},
			"method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Select a Method for the URL to create an API endpoint. Default is : *",
				Default:     "*",
			},
			"perform_staging": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "If true then any violation associated to the respective URL will not be enforced, " +
					"and the request will not be considered illegal.",
			},
			"signature_overrides_disable": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of Attack Signature Ids which are disabled for this particular URL.",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The created JSON for WAF URL entity.",
			},
		},
	}
}

func dataSourceBigipWafEntityUrlRead(d *schema.ResourceData, meta interface{}) error {
	_ = meta // we never call the device, to avoid complier errors we zero this out here

	name := d.Get("name").(string)
	d.SetId(name)
	log.Println("[INFO] Creating URL " + name)
	data := make(map[string]interface{})
	data["name"] = name

	if _, ok := d.GetOk("description"); ok {
		data["description"] = d.Get("description")
	}
	if _, ok := d.GetOk("type"); ok {
		data["type"] = d.Get("type")
	}
	if _, ok := d.GetOk("protocol"); ok {
		data["protocol"] = strings.ToUpper(d.Get("protocol").(string))
	}
	if _, ok := d.GetOk("method"); ok {
		data["method"] = d.Get("method")
	}
	if _, ok := d.GetOk("perform_staging"); ok {
		data["performStaging"] = d.Get("perform_staging")
	}
	if _, ok := d.GetOk("signature_overrides_disable"); ok {
		sigids := d.Get("signature_overrides_disable")
		var sigs []map[string]interface{}
		for _, s := range sigids.([]interface{}) {
			s1 := map[string]interface{}{"enabled": false, "signatureId": s}
			sigs = append(sigs, s1)
		}
		data["signatureOverrides"] = sigs
	}
	data["attackSignaturesCheck"] = true
	data["isAllowed"] = true
	data["methodsOverrideOnUrlCheck"] = false

	jsonString, err := json.Marshal(data)
	if err != nil {
		return err
	}
	log.Printf(string(jsonString))
	_ = d.Set("json", string(jsonString))

	return nil
}
