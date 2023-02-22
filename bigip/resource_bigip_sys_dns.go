/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipSysDns() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysDnsCreate,
		UpdateContext: resourceBigipSysDnsUpdate,
		ReadContext:   resourceBigipSysDnsRead,
		DeleteContext: resourceBigipSysDnsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User defined description",
				//ValidateFunc: validateF5Name,
			},
			"name_servers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Specifies the name servers that the system uses to validate DNS lookups, and resolve host names",
			},
			"number_of_dots": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "how many DNS Servers",
			},
			"search": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies the domains that the system searches for local domain lookups, to resolve local host names",
			},
		},
	}
}

func resourceBigipSysDnsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Get("description").(string)
	log.Println("[INFO] Configuring System DNS Server: " + description)
	configSysDns := &bigip.DNS{
		Description: description,
	}
	sysDNSConfig := getSysDNSConfig(d, configSysDns)

	err := client.ModifyDNS(sysDNSConfig)

	if err != nil {
		log.Printf("[ERROR] Unable to Create DNS (%s) (%v) ", description, err)
		return diag.FromErr(err)
	}
	d.SetId(description)

	return resourceBigipSysDnsRead(ctx, d, meta)
}

func resourceBigipSysDnsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Updating System DNS Server:" + description)

	configSysDns := &bigip.DNS{
		Description: description,
	}
	sysDNSConfig := getSysDNSConfig(d, configSysDns)

	err := client.ModifyDNS(sysDNSConfig)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify DNS (%s) (%v) ", description, err)
		return diag.FromErr(err)
	}
	return resourceBigipSysDnsRead(ctx, d, meta)
}

func resourceBigipSysDnsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading DNS " + description)

	dns, err := client.DNSs()
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve DNS (%s) (%v) ", description, err)
		return diag.FromErr(err)
	}
	if dns == nil {
		log.Printf("[WARN] DNS (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("description", dns.Description)

	_ = d.Set("name_servers", dns.NameServers)
	_ = d.Set("number_of_dots", dns.NumberOfDots)
	_ = d.Set("search", dns.Search)

	return nil
}

func resourceBigipSysDnsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no Delete API for this operation
	client := meta.(*bigip.BigIP)
	description := d.Id()
	log.Println("[INFO] Deleting System DNS Server:" + description)
	configSysDns := &bigip.DNS{
		Description:  description,
		NameServers:  []string{},
		Search:       []string{},
		NumberOfDots: 0,
	}
	err := client.ModifyDNS(configSysDns)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete DNS (%s) (%v) ", description, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getSysDNSConfig(d *schema.ResourceData, config *bigip.DNS) *bigip.DNS {
	config.NameServers = listToStringSlice(d.Get("name_servers").([]interface{}))
	config.NumberOfDots = d.Get("number_of_dots").(int)
	config.Search = listToStringSlice(d.Get("search").([]interface{}))
	return config
}
