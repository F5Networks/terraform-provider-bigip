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

func resourceBigipLtmProfileTcp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileTcpCreate,
		Update: resourceBigipLtmProfileTcpUpdate,
		Read:   resourceBigipLtmProfileTcpRead,
		Delete: resourceBigipLtmProfileTcpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the TCP Profile",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/tcp",
				Description: "Use the parent tcp profile",
			},

			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "idle_timeout can be given value",
			},

			"close_wait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "close wait timer integer",
			},

			"finwait_2timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "timer integer",
			},

			"finwait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "fin wait timer integer",
			},

			"keepalive_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "keepalive_interval timer integer",
			},

			"deferred_accept": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Defferred accept",
			},
			"fast_open": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "fast_open value ",
			},
		},
	}

}

func resourceBigipLtmProfileTcpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	idleTimeout := d.Get("idle_timeout").(int)
	closeWaitTimeout := d.Get("close_wait_timeout").(int)
	finWait_2Timeout := d.Get("finwait_2timeout").(int)
	finWaitTimeout := d.Get("finwait_timeout").(int)
	keepAliveInterval := d.Get("keepalive_interval").(int)
	deferredAccept := d.Get("deferred_accept").(string)
	fastOpen := d.Get("fast_open").(string)
	log.Println("[INFO] Creating TCP profile")

	err := client.CreateTcp(
		name,
		partition,
		defaultsFrom,
		idleTimeout,
		closeWaitTimeout,
		finWait_2Timeout,
		finWaitTimeout,
		keepAliveInterval,
		deferredAccept,
		fastOpen,
	)

	if err != nil {
		log.Printf("[ERROR] Unable to Create tcp Profile  (%s) (%v)", name, err)
		return err
	}
	d.SetId(name)
	return resourceBigipLtmProfileTcpRead(d, meta)
}

func resourceBigipLtmProfileTcpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Tcp{
		Name:              name,
		Partition:         d.Get("partition").(string),
		DefaultsFrom:      d.Get("defaults_from").(string),
		IdleTimeout:       d.Get("idle_timeout").(int),
		CloseWaitTimeout:  d.Get("close_wait_timeout").(int),
		FinWait_2Timeout:  d.Get("finwait_2timeout").(int),
		FinWaitTimeout:    d.Get("finwait_timeout").(int),
		KeepAliveInterval: d.Get("keepalive_interval").(int),
		DeferredAccept:    d.Get("deferred_accept").(string),
		FastOpen:          d.Get("fast_open").(string),
	}

	err := client.ModifyTcp(name, r)
	if err != nil {
		return fmt.Errorf("Error create profile tcp (%s): %s", name, err)
	}
	return resourceBigipLtmProfileTcpRead(d, meta)
}

func resourceBigipLtmProfileTcpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetTcp(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrive tcp Profile  (%s) (%v)", name, err)
		return err
	}
	if obj == nil {
		log.Printf("[WARN] tcp  Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	d.Set("partition", obj.Partition)
	if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for tcp profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("idle_timeout", obj.IdleTimeout); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IdleTimeout to state for tcp profile  (%s): %s", d.Id(), err)
	}
	if err := d.Set("close_wait_timeout", obj.CloseWaitTimeout); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CloseWaitTimeout to state for tcp profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("finwait_2timeout", obj.FinWait_2Timeout); err != nil {
		return fmt.Errorf("[DEBUG] Error saving FinWait_2Timeout to state for tcp profile  (%s): %s", d.Id(), err)
	}
	if err := d.Set("finwait_timeout", obj.FinWaitTimeout); err != nil {
		return fmt.Errorf("[DEBUG] Error saving FinWaitTimeout to state for tcp profile  (%s): %s", d.Id(), err)
	}

	d.Set("keepalive_interval", obj.KeepAliveInterval)
	if err := d.Set("deferred_accept", obj.DeferredAccept); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DeferredAccept to state for tcp profile  (%s): %s", d.Id(), err)
	}
	d.Set("fast_open", obj.FastOpen)

	return nil
}

func resourceBigipLtmProfileTcpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Tcp Profile " + name)

	err := client.DeleteTcp(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete tcp Profile (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}
