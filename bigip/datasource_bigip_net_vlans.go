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

func dataSourceBigipNetVlans() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipNetVlansRead,
		Schema: map[string]*schema.Schema{
			"vlans": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of every VLAN currently configured on the BIG-IP system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the VLAN.",
						},
						"partition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Partition the VLAN belongs to.",
						},
						"full_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Full path of the VLAN (partition and name).",
						},
						"tag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VLAN ID (tag).",
						},
						"mtu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum Transmission Unit (MTU) for the VLAN.",
						},
						"cmp_hash": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies how traffic on the VLAN is disaggregated.",
						},
						"auto_lasthop": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies whether auto lasthop is enabled or disabled on the VLAN.",
						},
						"failsafe": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether failsafe is enabled on the VLAN.",
						},
						"failsafe_action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action to take when failsafe is triggered.",
						},
						"failsafe_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Failsafe timeout, in seconds.",
						},
						"learning": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the VLAN's learning mode.",
						},
						"source_checking": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies whether source checking is enabled on the VLAN.",
						},
						"interfaces": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Interfaces (or trunks) attached to the VLAN.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the interface or trunk attached to the VLAN.",
									},
									"tagged": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the interface is tagged on this VLAN.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceBigipNetVlansRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	log.Println("[INFO] Reading all VLANs")

	vlans, err := client.Vlans()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving VLANs: %v", err))
	}

	var result []map[string]interface{}
	if vlans != nil {
		for _, vlan := range vlans.Vlans {
			var interfaces []map[string]interface{}
			vlanInterfaces, ifErr := client.GetVlanInterfaces(vlan.FullPath)
			if ifErr != nil {
				return diag.FromErr(fmt.Errorf("error retrieving interfaces for VLAN %s: %v", vlan.FullPath, ifErr))
			}
			if vlanInterfaces != nil {
				for _, iface := range vlanInterfaces.VlanInterfaces {
					interfaces = append(interfaces, map[string]interface{}{
						"name":   iface.Name,
						"tagged": iface.Tagged,
					})
				}
			}

			result = append(result, map[string]interface{}{
				"name":             vlan.Name,
				"partition":        vlan.Partition,
				"full_path":        vlan.FullPath,
				"tag":              vlan.Tag,
				"mtu":              vlan.MTU,
				"cmp_hash":         vlan.CMPHash,
				"auto_lasthop":     vlan.AutoLastHop,
				"failsafe":         vlan.Failsafe,
				"failsafe_action":  vlan.FailsafeAction,
				"failsafe_timeout": vlan.FailsafeTimeout,
				"learning":         vlan.Learning,
				"source_checking":  vlan.SourceChecking,
				"interfaces":       interfaces,
			})
		}
	}

	if err := d.Set("vlans", result); err != nil {
		return diag.FromErr(fmt.Errorf("error setting vlans: %v", err))
	}

	d.SetId(fmt.Sprintf("%d", schema.HashString(fmt.Sprintf("%v", result))))

	return nil
}
