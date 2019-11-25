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

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmProfileOneconnect() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileOneconnectCreate,
		Update: resourceBigipLtmProfileOneconnectUpdate,
		Read:   resourceBigipLtmProfileOneconnectRead,
		Delete: resourceBigipLtmProfileOneconnectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Oneconnect Profile",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent oneconnect profile",
			},

			"idle_timeout_override": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				Description: "idleTimeoutOverride can be enabled or disabled",
			},

			"share_pools": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "sharePools can be enabled or disabled",
			},
			"source_mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source_mask can be 255.255.255.255",
			},

			"max_age": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "max_age has integer value typical 3600 sec",
			},
			"max_reuse": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "max_reuse has integer value typical 1000 sec",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "max_size has integer value typical 1000 sec",
			},
		},
	}

}

func resourceBigipLtmProfileOneconnectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	sharePools := d.Get("share_pools").(string)
	maxAge := d.Get("max_age").(int)
	maxReuse := d.Get("max_reuse").(int)
	maxSize := d.Get("max_size").(int)
	sourceMask := d.Get("source_mask").(string)
	idleTimeoutOverride := d.Get("idle_timeout_override").(string)

	log.Println("[INFO] Creating OneConnect profile")

	err := client.CreateOneconnect(
		name,
		idleTimeoutOverride,
		partition,
		defaultsFrom,
		sharePools,
		sourceMask,
		maxAge,
		maxReuse,
		maxSize,
	)

	if err != nil {
		return fmt.Errorf("Error create profile oneConnect (%s): %s", name, err)
	}
	d.SetId(name)
	return resourceBigipLtmProfileOneconnectRead(d, meta)
}

func resourceBigipLtmProfileOneconnectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Oneconnect{
		Name:                name,
		IdleTimeoutOverride: d.Get("idle_timeout_override").(string),
		Partition:           d.Get("partition").(string),
		DefaultsFrom:        d.Get("defaults_from").(string),
		SharePools:          d.Get("share_pools").(string),
		SourceMask:          d.Get("source_mask").(string),
		MaxAge:              d.Get("max_age").(int),
		MaxSize:             d.Get("max_size").(int),
		MaxReuse:            d.Get("max_reuse").(int),
	}

	err := client.ModifyOneconnect(name, r)
	if err != nil {

		log.Printf("[ERROR] Unable to Modify OneConnect profile   (%s) (%v) ", name, err)

		return err
	}
	return resourceBigipLtmProfileOneconnectRead(d, meta)
}

func resourceBigipLtmProfileOneconnectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetOneconnect(name)
	if err != nil {
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Onceconnect Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	d.Set("partition", obj.Partition)
	if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for Onceconnect profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("share_pools", obj.SharePools); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SharePools to state for Onceconnect profile  (%s): %s", d.Id(), err)
	}
	if err := d.Set("source_mask", obj.SourceMask); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SourceMask to state for Onceconnect profile  (%s): %s", d.Id(), err)
	}
	if err := d.Set("max_age", obj.MaxAge); err != nil {
		return fmt.Errorf("[DEBUG] Error saving MaxAge to state for Onceconnect profile  (%s): %s", d.Id(), err)
	}
	d.Set("max_size", obj.MaxSize)
	d.Set("max_reuse", obj.MaxReuse)
	d.Set("idle_timeout_override", obj.IdleTimeoutOverride)
	return nil
}

func resourceBigipLtmProfileOneconnectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting OneConnect Profile " + name)

	err := client.DeleteOneconnect(name)
	if err != nil {
		return fmt.Errorf("Error Deleting profile oneConnect (%s): %s", name, err)
	}
	if err == nil {
		log.Printf("[WARN] OneConnect profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}
