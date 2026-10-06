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

// interfaceStatsResponse mirrors the shape of a
// GET /mgmt/tm/net/interface/stats (bulk, all interfaces at once) response.
// Each key under "entries" is that particular interface's stats selfLink
// URL (not a fixed name), so it's decoded as a map; each entry's own
// "tmName" nested field identifies which interface it belongs to.
type interfaceStatsResponse struct {
	Entries map[string]struct {
		NestedStats struct {
			Entries map[string]struct {
				Description string `json:"description,omitempty"`
			} `json:"entries"`
		} `json:"nestedStats"`
	} `json:"entries"`
}

// interfaceStatuses retrieves the operational link status (e.g. "up",
// "down", "uninit") for every interface in a single request via the bulk
// net/interface/stats endpoint, keyed by interface name -- the
// interfacestate object itself (net/interface/<name>) has no such field;
// "mediaActive" there reflects negotiated media, not link status, and
// go-bigip has no client wrapper for interface stats at all. A per-interface
// request (net/interface/<name>/stats) would also work but turns an O(1)
// call into an O(n) one for no benefit, since the bulk endpoint already
// returns every interface's stats in one response. Returns an empty map and
// logs a warning (rather than failing the whole data source read) if the
// stats endpoint doesn't return anything usable, since this is supplemental
// information layered on top of the primary net/interface listing.
func interfaceStatuses(client *bigip.BigIP) map[string]string {
	var stats interfaceStatsResponse
	found, err := restGet(client, "net/interface/stats", &stats)
	if err != nil || !found {
		log.Printf("[WARN] Unable to retrieve interface stats/status: %v", err)
		return map[string]string{}
	}
	result := make(map[string]string, len(stats.Entries))
	for _, entry := range stats.Entries {
		name := entry.NestedStats.Entries["tmName"].Description
		if name == "" {
			continue
		}
		result[name] = entry.NestedStats.Entries["status"].Description
	}
	return result
}

// dataSourceBigipNetInterfaces lists every physical interface on the BIG-IP
// system, for use in mapping i-Series interface layout (e.g. "1.1") to the
// equivalent r-Series/F5OS interface naming -- rSeries (F5OS-A) uses
// "<port>.0" (e.g. "1.0"), while VELOS (F5OS-C, chassis-based) uses
// "<blade>/<port>.<subport>" (e.g. "1/1.0") -- see
// docs/guides/interface-trunk-mapping.md for the full naming-convention
// differences between the two platforms.
func dataSourceBigipNetInterfaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipNetInterfacesRead,
		Schema: map[string]*schema.Schema{
			"interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of every physical interface currently present on the BIG-IP system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TMOS interface name, e.g. \"1.1\" (blade.port) or \"mgmt\".",
						},
						"full_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Full path of the interface.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the interface is administratively enabled.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operational link status from the interface's stats (net/interface/stats), e.g. \"up\", \"down\", \"uninit\". Not present on the net/interface object itself; retrieved via a separate, single bulk stats request covering every interface.",
						},
						"media_active": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Currently negotiated media/speed/duplex, e.g. \"10000T-FD\", or \"none\" if the link is down/not negotiated. This is the closest equivalent to a link status/speed indicator TMOS exposes for this object.",
						},
						"media_fixed": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Media type configured for the interface when not using auto-negotiation.",
						},
						"media_max": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Maximum media/speed capability advertised for auto-negotiation.",
						},
						"media_sfp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Media setting for SFP/SFP+ transceiver slots, where applicable.",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC address of the interface.",
						},
						"mtu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Configured MTU of the interface.",
						},
						"bundle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the interface supports bundling (trunk membership), e.g. \"not-supported\", \"enabled\", \"disabled\".",
						},
						"if_index": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "SNMP ifIndex of the interface.",
						},
						"flow_control": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flow control setting, e.g. \"tx-rx\", \"none\".",
						},
						"lldp_admin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "LLDP administrative mode for the interface, e.g. \"txonly\", \"disable\".",
						},
						"stp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether Spanning Tree Protocol is enabled on the interface.",
						},
						"stp_link_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Spanning Tree link type for the interface.",
						},
						"prefer_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "For combination ports, which physical media (sfp/copper) is preferred.",
						},
					},
				},
			},
		},
	}
}

func dataSourceBigipNetInterfacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	log.Println("[INFO] Reading all interfaces")

	interfaces, err := client.Interfaces()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving interfaces: %v", err))
	}

	var result []map[string]interface{}
	if interfaces != nil {
		statuses := interfaceStatuses(client)
		for _, iface := range interfaces.Interfaces {
			result = append(result, map[string]interface{}{
				"name":          iface.Name,
				"full_path":     iface.FullPath,
				"enabled":       iface.Enabled,
				"status":        statuses[iface.FullPath],
				"media_active":  iface.MediaActive,
				"media_fixed":   iface.MediaFixed,
				"media_max":     iface.MediaMax,
				"media_sfp":     iface.MediaSFP,
				"mac_address":   iface.MACAddress,
				"mtu":           iface.MTU,
				"bundle":        iface.Bundle,
				"if_index":      iface.IfIndex,
				"flow_control":  iface.FlowControl,
				"lldp_admin":    iface.LLDPAdmin,
				"stp":           iface.STP,
				"stp_link_type": iface.STPLinkType,
				"prefer_port":   iface.PreferPort,
			})
		}
	}

	if err := d.Set("interfaces", result); err != nil {
		return diag.FromErr(fmt.Errorf("error setting interfaces: %v", err))
	}

	d.SetId(hashForState(fmt.Sprintf("%v", result)))

	return nil
}
