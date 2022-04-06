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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
				Description:  "Name of the TCP Profile",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateF5Name,
				Description:  "Use the parent tcp profile",
			},
			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 300; may not be 0) connection may remain idle before it becomes eligible for deletion. Value -1 (not recommended) means infinite",
			},
			"close_wait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 5) connection will remain in LAST-ACK state before exiting. Value -1 means indefinite, limited by maximum retransmission timeout",
			},
			"finwait_2timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 300) connection will remain in LAST-ACK state before closing. Value -1 means indefinite, limited by maximum retransmission timeout",
			},
			"finwait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 5) connection will remain in FIN-WAIT-1 or closing state before exiting. Value -1 means indefinite, limited by maximum retransmission timeout",
			},
			"keepalive_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 1800) between keep-alive probes",
			},
			"deferred_accept": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "If enabled, ADC will defer allocating resources to a connection until some payload data has arrived from the client (default false). This may help minimize the impact of certain DoS attacks but adds undesirable latency under normal conditions. Note: ‘deferredAccept’ is incompatible with server-speaks-first application protocols,Default : disabled",
			},
			"fast_open": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "If enabled (default), the system can use the TCP Fast Open protocol extension to reduce latency by sending payload data with initial SYN",
			},
		},
	}
}

func resourceBigipLtmProfileTcpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	tcpProfileConfig := &bigip.Tcp{
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
	log.Println("[INFO] Creating TCP profile")
	err := client.CreateTcp(tcpProfileConfig)
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
	log.Println("[INFO] Updating TCP Profile Route " + name)
	tcpProfileConfig := &bigip.Tcp{
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
	err := client.ModifyTcp(name, tcpProfileConfig)
	if err != nil {
		return fmt.Errorf("Error create profile tcp (%s): %s ", name, err)
	}
	return resourceBigipLtmProfileTcpRead(d, meta)
}

func resourceBigipLtmProfileTcpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetTcp(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve tcp Profile  (%s) (%v)", name, err)
		return err
	}
	if obj == nil {
		log.Printf("[WARN] tcp  Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	if _, ok := d.GetOk("defaults_from"); ok {
		if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
			return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for tcp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("idle_timeout"); ok {
		if err := d.Set("idle_timeout", obj.IdleTimeout); err != nil {
			return fmt.Errorf("[DEBUG] Error saving IdleTimeout to state for tcp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("close_wait_timeout"); ok {
		if err := d.Set("close_wait_timeout", obj.CloseWaitTimeout); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CloseWaitTimeout to state for tcp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("finwait_2timeout"); ok {
		if err := d.Set("finwait_2timeout", obj.FinWait_2Timeout); err != nil {
			return fmt.Errorf("[DEBUG] Error saving FinWait_2Timeout to state for tcp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("finwait_timeout"); ok {
		if err := d.Set("finwait_timeout", obj.FinWaitTimeout); err != nil {
			return fmt.Errorf("[DEBUG] Error saving FinWaitTimeout to state for tcp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("keepalive_interval"); ok {
		_ = d.Set("keepalive_interval", obj.KeepAliveInterval)
	}
	if _, ok := d.GetOk("deferred_accept"); ok {
		if err := d.Set("deferred_accept", obj.DeferredAccept); err != nil {
			return fmt.Errorf("[DEBUG] Error saving DeferredAccept to state for tcp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("fast_open"); ok {
		_ = d.Set("fast_open", obj.FastOpen)
	}
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
