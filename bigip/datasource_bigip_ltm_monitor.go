/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBigipLtmMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipLtmMonitorRead,
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
				Description: "Existing monitor to inherit from. Must be one of /Common/http, /Common/https, /Common/icmp or /Common/gateway-icmp.",
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
		},
	}
}

func dataSourceBigipLtmMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	d.SetId("")
	name := fmt.Sprintf("/%s/%s", d.Get("partition").(string), d.Get("name").(string))
	log.Printf("[DEBUG] Retrieving Monitor: %s", name)
	monitors, err := client.Monitors()
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Monitor (%s) (%v) ", name, err)
		return err
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
			d.SetId(m.FullPath)
			return nil
		}
	}
	return nil
}
