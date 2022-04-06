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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
				Description:  "Name of the Oneconnect Profile",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent oneconnect profile",
			},
			"idle_timeout_override": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "idleTimeoutOverride can be enabled or disabled",
			},
			"share_pools": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "sharePools can be enabled or disabled",
			},
			"source_mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "source_mask can be 255.255.255.255",
			},
			"limit_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Controls how connection limits are enforced in conjunction with OneConnect. The default is None. Supported Values: [None,idle,strict]",
			},
			"max_age": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "max_age has integer value typical 3600 sec",
			},
			"max_reuse": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "max_reuse has integer value typical 1000 sec",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
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

	oneConnectconfig := &bigip.Oneconnect{
		Name:                name,
		IdleTimeoutOverride: idleTimeoutOverride,
		Partition:           partition,
		LimitType:           d.Get("limit_type").(string),
		DefaultsFrom:        defaultsFrom,
		SharePools:          sharePools,
		SourceMask:          sourceMask,
		MaxAge:              maxAge,
		MaxReuse:            maxReuse,
		MaxSize:             maxSize,
	}

	err := client.CreateOneconnect(oneConnectconfig)

	if err != nil {
		return fmt.Errorf("Error create profile oneConnect (%s): %s ", name, err)
	}
	d.SetId(name)
	return resourceBigipLtmProfileOneconnectRead(d, meta)
}

func resourceBigipLtmProfileOneconnectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating OneConnect Profile :%+v", name)
	r := &bigip.Oneconnect{
		Name:                name,
		IdleTimeoutOverride: d.Get("idle_timeout_override").(string),
		Partition:           d.Get("partition").(string),
		LimitType:           d.Get("limit_type").(string),
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
	log.Printf("[INFO] Reading OneConnect Profile :%+v", name)
	obj, err := client.GetOneconnect(name)
	if err != nil {
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Onceconnect Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("partition", obj.Partition)
	}
	if _, ok := d.GetOk("defaults_from"); ok {
		if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
			return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for Onceconnect profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("share_pools"); ok {
		if err := d.Set("share_pools", obj.SharePools); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SharePools to state for Onceconnect profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("source_mask"); ok {
		if err := d.Set("source_mask", obj.SourceMask); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SourceMask to state for Onceconnect profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("max_age"); ok {
		if err := d.Set("max_age", obj.MaxAge); err != nil {
			return fmt.Errorf("[DEBUG] Error saving MaxAge to state for Onceconnect profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("max_size"); ok {
		_ = d.Set("max_size", obj.MaxSize)
	}
	if _, ok := d.GetOk("limit_type"); ok {
		_ = d.Set("limit_type", obj.LimitType)
	}
	if _, ok := d.GetOk("max_reuse"); ok {
		_ = d.Set("max_reuse", obj.MaxReuse)
	}
	if _, ok := d.GetOk("idle_timeout_override"); ok {
		_ = d.Set("idle_timeout_override", obj.IdleTimeoutOverride)
	}
	return nil
}
func resourceBigipLtmProfileOneconnectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting OneConnect Profile " + name)
	err := client.DeleteOneconnect(name)
	if err != nil {
		return fmt.Errorf("Error Deleting profile oneConnect (%s): %s ", name, err)
	}
	d.SetId("")
	return nil
}
