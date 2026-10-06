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

func dataSourceBigipNetRoutes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipNetRoutesRead,
		Schema: map[string]*schema.Schema{
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of every static route currently configured on the BIG-IP system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the route.",
						},
						"partition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Partition the route belongs to.",
						},
						"full_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Full path of the route (partition and name).",
						},
						"network": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination network of the route.",
						},
						"gw": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway address for the route.",
						},
						"tm_interface": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tunnel/interface used to route traffic, if configured instead of a gateway.",
						},
						"mtu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MTU configured for the route.",
						},
						"blackhole": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the route is configured to reject/blackhole traffic.",
						},
					},
				},
			},
		},
	}
}

func dataSourceBigipNetRoutesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	log.Println("[INFO] Reading all routes")

	routes, err := client.Routes()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving routes: %v", err))
	}

	var result []map[string]interface{}
	if routes != nil {
		for _, route := range routes.Routes {
			result = append(result, map[string]interface{}{
				"name":         route.Name,
				"partition":    route.Partition,
				"full_path":    route.FullPath,
				"network":      route.Network,
				"gw":           route.Gateway,
				"tm_interface": route.TmInterface,
				"mtu":          route.MTU,
				"blackhole":    route.Blackhole,
			})
		}
	}

	if err := d.Set("routes", result); err != nil {
		return diag.FromErr(fmt.Errorf("error setting routes: %v", err))
	}

	d.SetId(fmt.Sprintf("%d", schema.HashString(fmt.Sprintf("%v", result))))

	return nil
}
