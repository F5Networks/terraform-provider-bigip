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

func resourceBigipLtmProfileFastl4() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipProfileLtmFastl4Create,
		Update: resourceBigipLtmProfileFastl4Update,
		Read:   resourceBigipLtmProfileFastl4Read,
		Delete: resourceBigipLtmProfileFastl4Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5Name,
				Description:  "Name of the Fastl4 Profile",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateF5Name,
				Description:  "Use the parent Fastl4 profile",
			},
			"client_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds allowed for a client to transmit enough data to select a server when you have late binding enabled. Value -1 means indefinite (not recommended)",
			},
			"explicitflow_migration": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"hardware_syncookie": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"idle_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 300; may not be 0) connection may remain idle before it becomes eligible for deletion. Value -1 (not recommended) means infinite",
			},
			"iptos_toclient": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"iptos_toserver": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"keepalive_interval": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent Fastl4 profile",
			},
		},
	}
}

func resourceBigipProfileLtmFastl4Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	clientTimeout := d.Get("client_timeout").(int)
	explicitFlowMigration := d.Get("explicitflow_migration").(string)
	hardwareSynCookie := d.Get("hardware_syncookie").(string)
	idleTimeout := d.Get("idle_timeout").(string)
	ipTosToClient := d.Get("iptos_toclient").(string)
	ipTosToServer := d.Get("iptos_toserver").(string)
	keepAliveInterval := d.Get("keepalive_interval").(string)

	log.Println("[INFO] Creating Fastl4 profile")

	err := client.CreateFastl4(
		name,
		partition,
		defaultsFrom,
		clientTimeout,
		explicitFlowMigration,
		hardwareSynCookie,
		idleTimeout,
		ipTosToClient,
		ipTosToServer,
		keepAliveInterval,
	)
	if err != nil {
		log.Printf("[ERROR] Unable to Create FastL4  (%s) (%v) ", name, err)
		return fmt.Errorf("Error retrieving profile fastl4 (%s): %s", name, err)
	}

	d.SetId(name)
	return resourceBigipLtmProfileFastl4Read(d, meta)
}

func resourceBigipLtmProfileFastl4Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	r := &bigip.Fastl4{
		Name:                  name,
		Partition:             d.Get("partition").(string),
		DefaultsFrom:          d.Get("defaults_from").(string),
		ClientTimeout:         d.Get("client_timeout").(int),
		ExplicitFlowMigration: d.Get("explicitflow_migration").(string),
		HardwareSynCookie:     d.Get("hardware_syncookie").(string),
		IdleTimeout:           d.Get("idle_timeout").(string),
		IpTosToClient:         d.Get("iptos_toclient").(string),
		IpTosToServer:         d.Get("iptos_toserver").(string),
		KeepAliveInterval:     d.Get("keepalive_interval").(string),
	}

	err := client.ModifyFastl4(name, r)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify FastL4  (%s) (%v) ", name, err)
		return err
	}
	return resourceBigipLtmProfileFastl4Read(d, meta)
}

func resourceBigipLtmProfileFastl4Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetFastl4(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve FastL4  (%s) (%v) ", name, err)
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Fastl4 profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("partition", obj.Partition)
	_ = d.Set("defaults_from", obj.DefaultsFrom)

	if _, ok := d.GetOk("client_timeout"); ok {
		if err := d.Set("client_timeout", obj.ClientTimeout); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ClientTimeout to state for FastL4 profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("explicitflow_migration"); ok {
		if err := d.Set("explicitflow_migration", obj.ExplicitFlowMigration); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ExplicitFlowMigration to state for FastL4 profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("iptos_toclient"); ok {
		_ = d.Set("iptos_toclient", obj.IpTosToClient)
	}
	if _, ok := d.GetOk("iptos_toserver"); ok {
		_ = d.Set("iptos_toserver", obj.IpTosToServer)
	}
	if _, ok := d.GetOk("hardware_syncookie"); ok {
		_ = d.Set("hardware_syncookie", obj.HardwareSynCookie)
	}
	if _, ok := d.GetOk("idle_timeout"); ok {
		_ = d.Set("idle_timeout", obj.IdleTimeout)
	}
	if _, ok := d.GetOk("keepalive_interval"); ok {
		_ = d.Set("keepalive_interval", obj.KeepAliveInterval)
	}

	return nil
}

func resourceBigipLtmProfileFastl4Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Fastl4 Profile " + name)

	err := client.DeleteFastl4(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve node (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}
