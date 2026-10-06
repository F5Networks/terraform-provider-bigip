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

// netRouteConfig is the response body for net/route/<name>. go-bigip's Route
// struct does not include description (see net.go), so this data source
// talks to the endpoint directly via restGet (see bigip/rest_helpers.go)
// instead of client.GetRoute, matching the precedent in
// resource_bigip_auth_ldap.go.
type netRouteConfig struct {
	Name        string `json:"name,omitempty"`
	Partition   string `json:"partition,omitempty"`
	FullPath    string `json:"fullPath,omitempty"`
	Description string `json:"description,omitempty"`
	Network     string `json:"network,omitempty"`
	Gateway     string `json:"gw,omitempty"`
	TmInterface string `json:"tmInterface,omitempty"`
	MTU         int    `json:"mtu,omitempty"`
	Blackhole   bool   `json:"blackhole,omitempty"`
}

// dataSourceBigipNetRoute looks up a single static route by name, returning
// its destination network, gateway, and description. This is the
// single-route counterpart to bigip_net_routes (which lists every route);
// use bigip_net_route when the name is already known, for example when
// migrating an existing tmsh-based workflow.
func dataSourceBigipNetRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipNetRouteRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the route to look up. May be a bare name (combined with `partition`) or a full path (e.g. `/Common/external-route`).",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				Description: "Partition of the route when `name` is a bare name. Ignored when `name` is already a full path. Defaults to `Common`.",
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
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the route.",
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
	}
}

func dataSourceBigipNetRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	name := d.Get("name").(string)
	fullPath := name
	if !strings.HasPrefix(fullPath, "/") {
		fullPath = fmt.Sprintf("/%s/%s", d.Get("partition").(string), name)
	}

	log.Printf("[DEBUG] Reading route: %s", fullPath)

	var route netRouteConfig
	found, err := restGet(client, "net/route/"+mangleFullPath(fullPath), &route)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving route %s: %v", fullPath, err))
	}
	if !found {
		return diag.FromErr(fmt.Errorf("route not found: %s", fullPath))
	}

	_ = d.Set("name", route.FullPath)
	_ = d.Set("partition", route.Partition)
	_ = d.Set("full_path", route.FullPath)
	_ = d.Set("network", route.Network)
	_ = d.Set("gw", route.Gateway)
	_ = d.Set("description", route.Description)
	_ = d.Set("tm_interface", route.TmInterface)
	_ = d.Set("mtu", route.MTU)
	_ = d.Set("blackhole", route.Blackhole)

	d.SetId(route.FullPath)

	return nil
}
