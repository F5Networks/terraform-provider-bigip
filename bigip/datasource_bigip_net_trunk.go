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

// dataSourceBigipNetTrunk looks up a single trunk (LAG) by name, returning
// its member interfaces and LACP settings. This is the single-trunk
// counterpart to bigip_net_trunks (which lists every trunk); use
// bigip_net_trunk when the name is already known, for example when mapping
// a specific TMOS trunk to its F5OS LAG equivalent during migration (see
// docs/guides/interface-trunk-mapping.md).
func dataSourceBigipNetTrunk() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipNetTrunkRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the trunk to look up.",
			},
			"interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Physical interfaces that are members of the trunk.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lacp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether LACP is enabled on the trunk.",
			},
			"lacp_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LACP mode (active/passive).",
			},
			"lacp_timeout": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LACP timeout (long/short).",
			},
			"distribution_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancing distribution hash used across trunk members.",
			},
			"link_select_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Link selection policy for the trunk.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Configured bandwidth of the trunk.",
			},
			"trunk_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Trunk ID.",
			},
			"stp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether Spanning Tree Protocol is enabled on the trunk.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Trunk type.",
			},
			"working_member_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of currently-working trunk members.",
			},
		},
	}
}

func dataSourceBigipNetTrunkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	name := d.Get("name").(string)

	log.Printf("[DEBUG] Reading trunk: %s", name)

	var trunk bigip.Trunk
	found, err := restGet(client, "net/trunk/"+mangleFullPath(name), &trunk)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving trunk %s: %v", name, err))
	}
	if !found {
		return diag.FromErr(fmt.Errorf("trunk not found: %s", name))
	}

	_ = d.Set("name", trunk.Name)
	_ = d.Set("interfaces", trunk.Interfaces)
	_ = d.Set("lacp", trunk.LACP)
	_ = d.Set("lacp_mode", trunk.LACPMode)
	_ = d.Set("lacp_timeout", trunk.LACPTimeout)
	_ = d.Set("distribution_hash", trunk.DistributionHash)
	_ = d.Set("link_select_policy", trunk.LinkSelectPolicy)
	_ = d.Set("bandwidth", trunk.Bandwidth)
	_ = d.Set("trunk_id", trunk.ID)
	_ = d.Set("stp", trunk.STP)
	_ = d.Set("type", trunk.Type)
	_ = d.Set("working_member_count", trunk.WorkingMemberCount)

	d.SetId(trunk.FullPath)

	return nil
}
