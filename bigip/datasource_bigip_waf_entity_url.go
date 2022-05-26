/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceBigipWafEntityUrl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipWafEntityUrlRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the URL.",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the URL.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies whether the parameter is an 'explicit' or a 'wildcard' attribute.",
				Default:      "wildcard",
				ValidateFunc: validation.StringInSlice([]string{"explicit", "wildcard"}, false),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies whether the protocol for the URL is 'http' or 'https'.",
				Default:      "http",
				ValidateFunc: validation.StringInSlice([]string{"http", "https"}, true),
			},
			"method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Select a method for the URL to create an API endpoint.",
				Default:     "*",
			},
			"perform_staging": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "If true then any violation associated to the respective URL will not be enforced, " +
					"and the request will not be considered illegal.",
			},
			"method_overrides": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of methods that are allowed or disallowed for a specific URL.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Specifies that the system allows or disallows a method for this URL.",
						},
						"method": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies an HTTP method.",
						},
					},
				},
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
	_ = meta // we never call the device, to avoid compiler errors we zero this out here

	name := d.Get("name").(string)
	d.SetId(name)
	log.Println("[INFO] Creating URL " + name)

	urlJson := &bigip.WafUrlJson{
		Name:                      name,
		Description:               d.Get("description").(string),
		Type:                      d.Get("type").(string),
		Protocol:                  d.Get("protocol").(string),
		Method:                    d.Get("method").(string),
		PerformStaging:            d.Get("perform_staging").(bool),
		AttackSignaturesCheck:     true,
		IsAllowed:                 true,
		MethodsOverrideOnUrlCheck: false,
	}
	sigCount := d.Get("signature_overrides_disable.#").(int)
	urlJson.SignatureOverrides = make([]bigip.WafUrlSig, 0, sigCount)
	for i := 0; i < sigCount; i++ {
		var s bigip.WafUrlSig
		prefix := fmt.Sprintf("signature_overrides_disable.%d", i)
		s.Enabled = false
		s.Id = d.Get(prefix).(int)
		urlJson.SignatureOverrides = append(urlJson.SignatureOverrides, s)
	}

	methodCount := d.Get("method_overrides.#").(int)
	urlJson.MethodOverrides = make([]bigip.MethodOverrides, 0, methodCount)
	for i := 0; i < methodCount; i++ {
		var m bigip.MethodOverrides
		prefix := fmt.Sprintf("method_overrides.%d", i)
		m.Allowed = d.Get(prefix + ".allow").(bool)
		m.Method = d.Get(prefix + ".method").(string)
		urlJson.MethodOverrides = append(urlJson.MethodOverrides, m)
	}

	jsonString, err := json.Marshal(urlJson)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] URL Json:%+v", string(jsonString))
	_ = d.Set("json", string(jsonString))
	return nil
}
