/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipNetRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipNetRouteCreate,
		UpdateContext: resourceBigipNetRouteUpdate,
		ReadContext:   resourceBigipNetRouteRead,
		DeleteContext: resourceBigipNetRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5Name,
				Description:  "Name of the route",
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination network",
			},
			"gw": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Gateway address",
			},
			"tunnel_ref": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateF5Name,
				Description:  "tunnel_ref to route traffic",
			},
			"reject": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "reject route",
			},
		},
	}
}

func resourceBigipNetRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	network := d.Get("network").(string)
	gw := d.Get("gw").(string)
	tunnelRef := d.Get("tunnel_ref").(string)
	reject := d.Get("reject").(bool)

	log.Println("[INFO] Creating Route")
	config := &bigip.Route{
		Name:    name,
		Network: network,
	}
	if gw != "" {
		config.Gateway = gw
	}
	if tunnelRef != "" {
		config.TmInterface = tunnelRef
	}
	if reject {
		config.Blackhole = reject
	}

	err := client.CreateRoute(config)

	if err != nil {
		log.Printf("[ERROR] Unable to Create Route  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipNetRouteRead(ctx, d, meta)
}

func resourceBigipNetRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Route " + name)
	network := d.Get("network").(string)
	gw := d.Get("gw").(string)
	tunnelRef := d.Get("tunnel_ref").(string)
	reject := d.Get("reject").(bool)

	config := &bigip.Route{
		Name:    name,
		Network: network,
	}
	if gw != "" {
		config.Gateway = gw
	}
	if tunnelRef != "" {
		config.TmInterface = tunnelRef
	}
	if reject {
		config.Blackhole = reject
	}

	err := client.ModifyRoute(name, config)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Route  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	return resourceBigipNetRouteRead(ctx, d, meta)
}

func resourceBigipNetRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Reading Net Route config :%+v", name)
	obj, err := client.GetRoute(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Route  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	if obj == nil {
		log.Printf("[WARN] Route (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("name", obj.FullPath)

	_ = d.Set("network", obj.Network)

	if obj.Gateway != "" || d.Get("gw").(string) != "" {
		_ = d.Set("gw", obj.Gateway)
	}
	if obj.TmInterface != "" || d.Get("tunnel_ref").(string) != "" {
		_ = d.Set("tunnel_ref", obj.TmInterface)
	}
	if obj.Blackhole {
		_ = d.Set("reject", obj.Blackhole)
	}
	return nil
}

func resourceBigipNetRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Route " + name)

	err := client.DeleteRoute(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Route  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
