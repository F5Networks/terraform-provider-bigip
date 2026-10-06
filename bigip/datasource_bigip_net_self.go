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
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceBigipNetSelf looks up a single self IP by name, returning its
// address, VLAN, traffic group, and allowed services. This is the
// single-self-IP counterpart to bigip_net_selfips (which lists every self
// IP); use bigip_net_self when the name is already known, for example when
// migrating an existing tmsh-based workflow.
func dataSourceBigipNetSelf() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipNetSelfRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the self IP to look up. May be a bare name (combined with `partition`) or a full path (e.g. `/Common/external-self`).",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				Description: "Partition of the self IP when `name` is a bare name. Ignored when `name` is already a full path. Defaults to `Common`.",
			},
			"full_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Full path of the self IP (partition and name).",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address (with subnet/route-domain suffix) of the self IP.",
			},
			"vlan": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VLAN the self IP is associated with.",
			},
			"traffic_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Traffic group the self IP belongs to.",
			},
			"floating": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the self IP is a floating address.",
			},
			"unit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unit ID for non-floating self IPs in a redundant pair.",
			},
			"allow_service": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Port lockdown / allowed services for the self IP (e.g. \"all\", \"none\", or a list of service:port entries).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceBigipNetSelfRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	name := d.Get("name").(string)
	fullPath := name
	if !strings.HasPrefix(fullPath, "/") {
		fullPath = fmt.Sprintf("/%s/%s", d.Get("partition").(string), name)
	}

	log.Printf("[DEBUG] Reading self IP: %s", fullPath)

	selfIP, err := client.SelfIP(fullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving self IP %s: %v", fullPath, err))
	}
	if selfIP == nil {
		return diag.FromErr(fmt.Errorf("self IP %s not found", fullPath))
	}

	var allowService []string
	switch v := selfIP.AllowService.(type) {
	case string:
		if v != "" {
			allowService = []string{v}
		} else {
			allowService = []string{"none"}
		}
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok {
				allowService = append(allowService, s)
			}
		}
	case nil:
		allowService = []string{"none"}
	}

	_ = d.Set("name", selfIP.FullPath)
	_ = d.Set("partition", selfIP.Partition)
	_ = d.Set("full_path", selfIP.FullPath)
	_ = d.Set("address", selfIP.Address)
	_ = d.Set("vlan", selfIP.Vlan)
	_ = d.Set("traffic_group", selfIP.TrafficGroup)
	_ = d.Set("floating", selfIP.Floating)
	_ = d.Set("unit", selfIP.Unit)
	_ = d.Set("allow_service", allowService)

	d.SetId(selfIP.FullPath)

	return nil
}
