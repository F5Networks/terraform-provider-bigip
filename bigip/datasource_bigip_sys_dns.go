/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceBigipSysDns reads the BIG-IP system's DNS resolver configuration
// (sys/dns) -- a singleton object, so this data source takes no arguments.
// Use bigip_sys_dns (the resource) to manage this configuration.
func dataSourceBigipSysDns() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipSysDnsRead,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-defined description of the system DNS configuration.",
			},
			"name_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name servers the system uses to validate DNS lookups and resolve host names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"search": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Domains the system searches for local domain lookups, to resolve local host names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"number_of_dots": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of dots that must appear in a host name before an initial absolute lookup is performed.",
			},
		},
	}
}

func dataSourceBigipSysDnsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	log.Println("[INFO] Reading system DNS configuration")

	dns, err := client.DNSs()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving system DNS configuration: %v", err))
	}
	if dns == nil {
		return diag.FromErr(fmt.Errorf("system DNS configuration not found"))
	}

	_ = d.Set("description", dns.Description)
	_ = d.Set("name_servers", dns.NameServers)
	_ = d.Set("search", dns.Search)
	_ = d.Set("number_of_dots", dns.NumberOfDots)

	// sys/dns has no path/name; resource_bigip_sys_dns.go uses description
	// as the ID for this same underlying object, so mirror that here
	// rather than hashing the struct. description is a required field on
	// the resource, but it's optional at the API level (BIG-IP omits it
	// entirely from the response if unset), so fall back to a hash of the
	// struct in that case to avoid setting an empty ID.
	id := dns.Description
	if id == "" {
		id = hashForState(fmt.Sprintf("%v", dns))
	}
	d.SetId(id)

	return nil
}
