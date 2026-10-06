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

func dataSourceBigipNetSelfips() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipNetSelfipsRead,
		Schema: map[string]*schema.Schema{
			"self_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of every self IP currently configured on the BIG-IP system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the self IP.",
						},
						"partition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Partition the self IP belongs to.",
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
				},
			},
		},
	}
}

func dataSourceBigipNetSelfipsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	log.Println("[INFO] Reading all self IPs")

	selfIPs, err := client.SelfIPs()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving self IPs: %v", err))
	}

	var result []map[string]interface{}
	if selfIPs != nil {
		for _, selfIP := range selfIPs.SelfIPs {
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

			result = append(result, map[string]interface{}{
				"name":          selfIP.Name,
				"partition":     selfIP.Partition,
				"full_path":     selfIP.FullPath,
				"address":       selfIP.Address,
				"vlan":          selfIP.Vlan,
				"traffic_group": selfIP.TrafficGroup,
				"floating":      selfIP.Floating,
				"unit":          selfIP.Unit,
				"allow_service": allowService,
			})
		}
	}

	if err := d.Set("self_ips", result); err != nil {
		return diag.FromErr(fmt.Errorf("error setting self_ips: %v", err))
	}

	d.SetId(fmt.Sprintf("%d", schema.HashString(fmt.Sprintf("%v", result))))

	return nil
}
