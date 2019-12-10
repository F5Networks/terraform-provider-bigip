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
	"regexp"

	"github.com/f5devcentral/go-bigip"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the route",
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
		},
	}

}

func resourceBigipNetRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	network := d.Get("network").(string)
	gw := d.Get("gw").(string)

	log.Println("[INFO] Creating Route")

	err := client.CreateRoute(
		name,
		network,
		gw,
	)

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

	r := &bigip.Route{
		Name:    name,
		Network: d.Get("network").(string),
	}

	err := client.ModifyRoute(name, r)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Route  (%s) (%v)", name, err)
		return err
	}
	return resourceBigipNetRouteRead(d, meta)
}

func resourceBigipNetRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
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
	d.Set("name", name)

	regex := regexp.MustCompile(`(default|(?:[0-9]{1,3}\.){3}[0-9]{1,3}\/[0-9]{1,2})(?:\%\d+)?`)
	network := regex.FindStringSubmatch(obj.Network)

	regex = regexp.MustCompile(`((?:[0-9]{1,3}\.){3}[0-9]{1,3})(?:\%\d+)?`)
	gw := regex.FindStringSubmatch(obj.Gateway)

	if err := d.Set("network", network[1]); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Network to state for Route (%s): %s", d.Id(), err)
	}

	if err := d.Set("gw", gw[1]); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Gateway to state for Route (%s): %s", d.Id(), err)
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
