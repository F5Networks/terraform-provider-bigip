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

// dataSourceBigipNetVlan looks up a single VLAN by either its name or its
// VLAN ID (tag), returning the tagged/untagged interface membership. This is
// the single-VLAN counterpart to bigip_net_vlans (which lists every VLAN);
// use bigip_net_vlan when the identifier is already known (for example, when
// migrating an existing tmsh-based workflow that referenced a VLAN by tag).
func dataSourceBigipNetVlan() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipNetVlanRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "vlan_id"},
				Description:  "Name of the VLAN to look up. May be a bare name (combined with `partition`) or a full path (e.g. `/Common/external`). Exactly one of `name` or `vlan_id` must be specified.",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				Description: "Partition of the VLAN when looking up by `name`. Ignored when `name` is already a full path or when looking up by `vlan_id`. Defaults to `Common`.",
			},
			"vlan_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "vlan_id"},
				Description:  "VLAN ID (tag) to look up. Exactly one of `name` or `vlan_id` must be specified.",
			},
			"full_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Full path of the VLAN (partition and name).",
			},
			"tagged_interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Names of the interfaces or trunks attached to the VLAN as tagged members.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"untagged_interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Names of the interfaces or trunks attached to the VLAN as untagged members.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceBigipNetVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	var vlan *bigip.Vlan

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		fullPath := name
		if !strings.HasPrefix(fullPath, "/") {
			fullPath = fmt.Sprintf("/%s/%s", d.Get("partition").(string), name)
		}

		log.Printf("[DEBUG] Reading VLAN by name: %s", fullPath)

		found, err := client.Vlan(fullPath)
		if err != nil {
			if restIsNotFound(err) {
				return diag.FromErr(fmt.Errorf("VLAN not found: %s", fullPath))
			}
			return diag.FromErr(fmt.Errorf("error retrieving VLAN %s: %v", fullPath, err))
		}
		vlan = found
	} else if v, ok := d.GetOk("vlan_id"); ok {
		tag := v.(int)

		log.Printf("[DEBUG] Reading VLAN by vlan_id: %d", tag)

		vlans, err := client.Vlans()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error retrieving VLANs: %v", err))
		}
		if vlans != nil {
			for i := range vlans.Vlans {
				if vlans.Vlans[i].Tag == tag {
					vlan = &vlans.Vlans[i]
					break
				}
			}
		}
	} else {
		return diag.FromErr(fmt.Errorf("one of 'name' or 'vlan_id' must be specified"))
	}

	if vlan == nil {
		if v, ok := d.GetOk("vlan_id"); ok {
			return diag.FromErr(fmt.Errorf("VLAN not found: vlan_id %d", v.(int)))
		}
		return diag.FromErr(fmt.Errorf("VLAN not found"))
	}

	vlanInterfaces, err := client.GetVlanInterfaces(vlan.FullPath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving interfaces for VLAN %s: %v", vlan.FullPath, err))
	}

	var tagged, untagged []string
	if vlanInterfaces != nil {
		for _, iface := range vlanInterfaces.VlanInterfaces {
			if iface.Tagged {
				tagged = append(tagged, iface.Name)
			} else {
				untagged = append(untagged, iface.Name)
			}
		}
	}

	_ = d.Set("name", vlan.FullPath)
	_ = d.Set("partition", vlan.Partition)
	_ = d.Set("full_path", vlan.FullPath)
	_ = d.Set("vlan_id", vlan.Tag)
	_ = d.Set("tagged_interfaces", tagged)
	_ = d.Set("untagged_interfaces", untagged)

	d.SetId(vlan.FullPath)

	return nil
}
