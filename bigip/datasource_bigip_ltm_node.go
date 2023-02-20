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
	"regexp"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipLtmNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipLtmNodeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the node.",
			},
			"partition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Partition of resource group",
			},
			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Full path of the node (partition and name)",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description of the node.",
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address of the node of the node.",
			},
			"connection_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Node connection limit.",
			},
			"dynamic_ratio": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The dynamic ratio number for the node.",
			},
			"monitor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the health monitors the system currently uses to monitor this node.",
			},
			"rate_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node rate limit.",
			},
			"ratio": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Node ratio weight.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the node.",
			},
			"session": {
				Type:        schema.TypeString,
				Description: "Enables or disables the node for new sessions.",
				Computed:    true,
			},
			"fqdn": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_family": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The FQDN node's address family.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The fully qualified domain name of the node.",
						},
						"interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The amount of time before sending the next DNS query.",
						},
						"downinterval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The number of attempts to resolve a domain name.",
						},
						"autopopulate": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies if the node should scale to the IP address set returned by DNS.",
						},
					},
				},
			},
		},
	}
}
func dataSourceBigipLtmNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")
	name := fmt.Sprintf("/%s/%s", d.Get("partition").(string), d.Get("name").(string))
	log.Println("[DEBUG] Reading Node : " + name)
	node, err := client.GetNode(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving node %s: %v", name, err))
	}
	if node == nil {
		log.Printf("[DEBUG] Node %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	if node.FQDN.Name != "" {
		_ = d.Set("address", node.FQDN.Name)
	} else {
		// xxx.xxx.xxx.xxx(%x)
		// x:x(%x)
		regex := regexp.MustCompile(`((?:(?:[0-9]{1,3}\.){3}[0-9]{1,3})|(?:.*:[^%]*))(?:\%\d+)?`)
		address := regex.FindStringSubmatch(node.Address)
		log.Println("[INFO] Address: " + address[1])
		_ = d.Set("address", node.Address)
	}

	_ = d.Set("name", node.Name)
	_ = d.Set("partition", node.Partition)
	_ = d.Set("full_path", node.FullPath)
	_ = d.Set("connection_limit", node.ConnectionLimit)
	_ = d.Set("dynamic_ratio", node.DynamicRatio)
	_ = d.Set("monitor", node.Monitor)
	_ = d.Set("rate_limit", node.RateLimit)
	_ = d.Set("ratio", node.Ratio)
	_ = d.Set("state", node.State)
	_ = d.Set("session", node.Session)
	var fqdn []map[string]interface{}
	fqdnelements := map[string]interface{}{
		"interval":       node.FQDN.Interval,
		"downinterval":   node.FQDN.DownInterval,
		"autopopulate":   node.FQDN.AutoPopulate,
		"address_family": node.FQDN.AddressFamily,
	}
	fqdn = append(fqdn, fqdnelements)
	_ = d.Set("fqdn", fqdn)

	//	_ = d.Set("fqdn.0.interval", node.FQDN.Interval)
	//	_ = d.Set("fqdn.0.downinterval", node.FQDN.DownInterval)
	//	_ = d.Set("fqdn.0.autopopulate", node.FQDN.AutoPopulate)
	//	_ = d.Set("fqdn.0.address_family", node.FQDN.AddressFamily)
	d.SetId(node.Name)
	return nil
}
