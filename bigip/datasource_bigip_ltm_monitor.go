/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
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
			_ = d.Set("defaults_from", m.DefaultsFrom)
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
			_ = d.Set("password", m.Password)
			_ = d.Set("name", name)
			_ = d.Set("database", m.Database)
			d.SetId(m.FullPath)
			return nil
		}
	}
	return nil
}
