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

func dataSourceBigipLtmMonitor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipLtmMonitorRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the LTM Monitor",
			},
			"partition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition of LTM Monitor",
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Existing monitor to inherit from. Must be one of /Common/http, /Common/https, /Common/icmp, /Common/gateway-icmp or /Common/tcp-half-open.",
			},
			"interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"receive_disable": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expected response string.",
			},
			"reverse": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transparent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manual_resume": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_dscp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"time_until_up": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Time in seconds",
			},
			"destination": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alias for the destination",
			},
			"filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the full path and file name of the file that the system attempts to download. The health check is successful if the system can download the file.",
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the data transfer process (DTP) mode. The default value is passive.",
			},
			"adaptive": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ftp adaptive",
			},
			"adaptive_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Integer value",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the user name if the monitored target requires authentication",
			},
			"database": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the database in which your user is created",
			},

			"base": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the location in the LDAP tree from which the monitor starts the health check",
			},
			"filter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies an LDAP key for which the monitor searches",
			},
			"mandatory_attributes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies whether the target must include attributes in its response to be considered up. The options are no (Specifies that the system performs only a one-level search (based on the Filter setting), and does not require that the target returns any attributes.) and yes (Specifies that the system performs a sub-tree search, and if the target returns no attributes, the target is considered down.)",
			},
			"chase_referrals": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies whether the system will query the LDAP servers pointed to by any referrals in the query results.",
			},
			"security": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the secure communications protocol that the monitor uses to communicate with the target. The options are none (Specifies that the system does not use a security protocol for communications with the target.), ssl (Specifies that the system uses the SSL protocol for communications with the target.), and tls (Specifies that the system uses the TLS protocol for communications with the target.)",
			},
		},
	}
}

func dataSourceBigipLtmMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")
	name := fmt.Sprintf("/%s/%s", d.Get("partition").(string), d.Get("name").(string))
	log.Printf("[DEBUG] Retrieving Monitor: %s", name)
	monitors, err := client.Monitors()
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Monitor (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	if monitors == nil {
		log.Printf("[WARN] Monitor (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	for _, m := range monitors {
		if m.FullPath == name {
			_ = d.Set("defaults_from", m.ParentMonitor)
			_ = d.Set("interval", m.Interval)
			_ = d.Set("timeout", m.Timeout)
			_ = d.Set("receive_disable", m.ReceiveDisable)
			_ = d.Set("reverse", m.Reverse)
			_ = d.Set("transparent", m.Transparent)
			_ = d.Set("ip_dscp", m.IPDSCP)
			_ = d.Set("time_until_up", m.TimeUntilUp)
			_ = d.Set("manual_resume", m.ManualResume)
			_ = d.Set("destination", m.Destination)
			_ = d.Set("filename", m.Filename)
			_ = d.Set("mode", m.Mode)
			_ = d.Set("adaptive", m.Adaptive)
			_ = d.Set("adaptive_limit", m.AdaptiveLimit)
			_ = d.Set("username", m.Username)
			_ = d.Set("name", name)
			_ = d.Set("database", m.Database)
			_ = d.Set("base", m.Base)
			_ = d.Set("filter", m.Filter)
			_ = d.Set("mandatory_attributes", m.MandatoryAttributes)
			_ = d.Set("chase_referrals", m.ChaseReferrals)
			_ = d.Set("security", m.Security)
			d.SetId(m.FullPath)
			return nil
		}
	}
	return nil
}
