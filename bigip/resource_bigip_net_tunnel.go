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

func resourceBigipNetTunnel() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceBigipNetTunnelCreate,
		ReadContext:   resourceBigipNetTunnelRead,
		UpdateContext: resourceBigipNetTunnelUpdate,
		DeleteContext: resourceBigipNetTunnelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the TUNNEL",
			},
			"app_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application service that the object belongs to",
			},
			"auto_last_hop": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether auto lasthop is enabled or not",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description",
			},
			"local_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies a local IP address. This option is required",
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies how the tunnel carries traffic",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Displays the admin-partition within which this component resides",
			},
			"profile": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the profile that you want to associate with the tunnel",
			},
			"remote_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a remote IP address",
			},
			"secondary_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a secondary non-floating IP address when the local-address is set to a floating address",
			},
			"tos": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a value for insertion into the Type of Service (ToS) octet within the IP header of the encapsulating header of transmitted packets",
			},
			"traffic_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a traffic-group for use with the tunnel",
			},
			"transparent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enables or disables the tunnel to be transparent",
			},
			"use_pmtu": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enables or disables the tunnel to use the PMTU (Path MTU) information provided by ICMP NeedFrag error messages",
			},
			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies an idle timeout for wildcard tunnels in seconds",
			},

			"key": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The key field may represent different values depending on the type of the tunnel",
			},
			"mtu": {
				Type:     schema.TypeInt,
				Optional: true,
				//  Default:     0,
				Description: "Specifies the maximum transmission unit (MTU) of the tunnel",
			},
		},
	}

}

func resourceBigipNetTunnelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	log.Printf("[INFO] Creating TUNNEL %s", name)

	r := &bigip.Tunnel{
		Name: name,
	}
	config := getConfig(d, r)

	err := client.CreateTunnel(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Tunnel %s %v :", name, err)
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceBigipNetTunnelRead(ctx, d, meta)
}

func resourceBigipNetTunnelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Reading TUNNEL %s", name)
	tunnel, err := client.GetTunnel(name)
	log.Printf("[DEBUG] TUNNEL Output :%+v", tunnel)
	if err != nil {
		return diag.FromErr(err)
	}
	if tunnel == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("[ERROR] Tunnel (%s) not found, removing from state", d.Id()))
	}
	_ = d.Set("app_service", tunnel.AppService)

	_ = d.Set("auto_last_hop", tunnel.AutoLasthop)

	_ = d.Set("description", tunnel.Description)

	_ = d.Set("idle_timeout", tunnel.IdleTimeout)

	_ = d.Set("key", tunnel.Key)

	_ = d.Set("local_address", tunnel.LocalAddress)

	_ = d.Set("mode", tunnel.Mode)

	_ = d.Set("mtu", tunnel.Mtu)

	_ = d.Set("partition", tunnel.Partition)

	_ = d.Set("profile", tunnel.Profile)

	_ = d.Set("remote_address", tunnel.RemoteAddress)

	_ = d.Set("secondary_address", tunnel.SecondaryAddress)

	_ = d.Set("tos", tunnel.Tos)

	_ = d.Set("traffic_group", tunnel.TrafficGroup)

	_ = d.Set("transparent", tunnel.Transparent)

	_ = d.Set("use_pmtu", tunnel.UsePmtu)

	_ = d.Set("name", name)
	return nil
}

func resourceBigipNetTunnelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Updating Tunnel %s", name)

	r := &bigip.Tunnel{
		Name: name,
	}
	config := getConfig(d, r)

	err := client.ModifyTunnel(name, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error modifying TUNNEL %s: %v ", name, err))
	}

	return resourceBigipNetTunnelRead(ctx, d, meta)
}

func resourceBigipNetTunnelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Deleting TUNNEL %s", name)

	err := client.DeleteTunnel(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error Deleting Tunnel : %s ", err))
	}

	d.SetId("")
	return nil
}

func getConfig(d *schema.ResourceData, config *bigip.Tunnel) *bigip.Tunnel {
	config.AppService = d.Get("app_service").(string)
	config.AutoLasthop = d.Get("auto_last_hop").(string)
	config.Description = d.Get("description").(string)
	config.LocalAddress = d.Get("local_address").(string)
	config.Profile = d.Get("profile").(string)
	config.IdleTimeout = d.Get("idle_timeout").(int)
	//IfIndex:d.Get("if_index").(int),
	config.Key = d.Get("key").(int)
	config.Mode = d.Get("mode").(string)
	config.Mtu = d.Get("mtu").(int)
	config.Partition = d.Get("partition").(string)
	config.RemoteAddress = d.Get("remote_address").(string)
	config.SecondaryAddress = d.Get("secondary_address").(string)
	config.Tos = d.Get("tos").(string)
	config.TrafficGroup = d.Get("traffic_group").(string)
	config.Transparent = d.Get("transparent").(string)
	config.UsePmtu = d.Get("use_pmtu").(string)
	return config
}
