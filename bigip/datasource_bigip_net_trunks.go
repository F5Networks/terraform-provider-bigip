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

// dataSourceBigipNetTrunks lists every trunk (LAG) configured on the BIG-IP
// system, for use in mapping i-Series trunk membership to the equivalent
// r-Series/F5OS LAG configuration (see docs/guides/interface-trunk-mapping.md
// for the naming-convention differences between the two platforms).
func dataSourceBigipNetTrunks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipNetTrunksRead,
		Schema: map[string]*schema.Schema{
			"trunks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of every trunk (LAG) currently configured on the BIG-IP system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the trunk.",
						},
						"full_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Full path of the trunk.",
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
						"id": {
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
				},
			},
		},
	}
}

func dataSourceBigipNetTrunksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	log.Println("[INFO] Reading all trunks")

	trunks, err := client.Trunks()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving trunks: %v", err))
	}

	var result []map[string]interface{}
	if trunks != nil {
		for _, trunk := range trunks.Trunks {
			result = append(result, map[string]interface{}{
				"name":                 trunk.Name,
				"full_path":            trunk.FullPath,
				"interfaces":           trunk.Interfaces,
				"lacp":                 trunk.LACP,
				"lacp_mode":            trunk.LACPMode,
				"lacp_timeout":         trunk.LACPTimeout,
				"distribution_hash":    trunk.DistributionHash,
				"link_select_policy":   trunk.LinkSelectPolicy,
				"bandwidth":            trunk.Bandwidth,
				"id":                   trunk.ID,
				"stp":                  trunk.STP,
				"type":                 trunk.Type,
				"working_member_count": trunk.WorkingMemberCount,
			})
		}
	}

	if err := d.Set("trunks", result); err != nil {
		return diag.FromErr(fmt.Errorf("error setting trunks: %v", err))
	}

	d.SetId(hashForState(fmt.Sprintf("%v", result)))

	return nil
}
