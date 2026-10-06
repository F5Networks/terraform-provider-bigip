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

// dataSourceBigipSysNtp reads the BIG-IP system's NTP configuration
// (sys/ntp) -- a singleton object, so this data source takes no arguments.
// Use bigip_sys_ntp (the resource) to manage this configuration.
func dataSourceBigipSysNtp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipSysNtpRead,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-defined description of the system NTP configuration.",
			},
			// Named ntp_servers rather than matching resource_bigip_sys_ntp.go's
			// "servers" field -- this is a deliberate, explicitly-requested
			// naming choice for this data source, not an inconsistency with
			// bigip_sys_dns (which does mirror its resource's field names).
			"ntp_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Time servers the system uses to update the system time.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"timezone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time zone used for the system time.",
			},
		},
	}
}

func dataSourceBigipSysNtpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	log.Println("[INFO] Reading system NTP configuration")

	ntp, err := client.NTPs()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving system NTP configuration: %v", err))
	}
	if ntp == nil {
		return diag.FromErr(fmt.Errorf("system NTP configuration not found"))
	}

	_ = d.Set("description", ntp.Description)
	_ = d.Set("ntp_servers", ntp.Servers)
	_ = d.Set("timezone", ntp.Timezone)

	// sys/ntp has no path/name; resource_bigip_sys_ntp.go uses description
	// as the ID for this same underlying object, so mirror that here
	// rather than hashing the struct. description is a required field on
	// the resource, but it's optional at the API level (BIG-IP omits it
	// entirely from the response if unset), so fall back to a hash of the
	// struct in that case to avoid setting an empty ID.
	id := ntp.Description
	if id == "" {
		id = hashForState(fmt.Sprintf("%v", ntp))
	}
	d.SetId(id)

	return nil
}
