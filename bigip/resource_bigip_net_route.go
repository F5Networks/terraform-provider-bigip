/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipNetRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipNetRouteCreate,
		Update: resourceBigipNetRouteUpdate,
		Read:   resourceBigipNetRouteRead,
		Delete: resourceBigipNetRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceBigipNetRouteCreate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	d.SetId(name)
	return resourceBigipNetRouteRead(d, meta)
}

func resourceBigipNetRouteUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	return resourceBigipNetRouteRead(d, meta)
}

func resourceBigipNetRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[DEBUG] Reading Net Route config :%+v", name)
	obj, err := client.GetRoute(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Route  (%s) (%v)", name, err)
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Route (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", obj.FullPath)

	if err := d.Set("network", obj.Network); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Network to state for Route (%s): %s", d.Id(), err)
	}
	if obj.Gateway != "" || d.Get("gw").(string) != "" {
		d.Set("gw", obj.Gateway)
	}
	if obj.TmInterface != "" || d.Get("tunnel_ref").(string) != "" {
		d.Set("tunnel_ref", obj.TmInterface)
	}
	if obj.Blackhole {
		d.Set("reject", obj.Blackhole)
	}
	return nil
}

func resourceBigipNetRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Route " + name)

	err := client.DeleteRoute(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Route  (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}
