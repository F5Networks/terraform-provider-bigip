/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
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

func resourceBigipLtmVirtualAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmVirtualAddressCreate,
		ReadContext:   resourceBigipLtmVirtualAddressRead,
		UpdateContext: resourceBigipLtmVirtualAddressUpdate,
		DeleteContext: resourceBigipLtmVirtualAddressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the virtual address",
				ValidateFunc: validateVirtualAddressName,
			},

			"arp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable or disable ARP for the virtual address",
				Default:     true,
			},

			"auto_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Automatically delete the virtual address with the virtual server",
			},

			"conn_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Max number of connections for virtual address",
			},

			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable or disable the virtual address",
			},

			"icmp_echo": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "Enable/Disable ICMP response to the virtual address",
			},

			"advertize_route": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Enabled dynamic routing of the address",
			},

			"traffic_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "/Common/traffic-group-1",
				Description:  "Specify the partition and traffic group",
				ValidateFunc: validateF5Name,
			},
		},
	}
}

func resourceBigipLtmVirtualAddressCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Println("[INFO] Creating virtual address " + name)

	if err := client.CreateVirtualAddress(name, hydrateVirtualAddress(d)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	return resourceBigipLtmVirtualAddressRead(ctx, d, meta)
}

func resourceBigipLtmVirtualAddressRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching virtual address " + name)

	var va bigip.VirtualAddress
	vas, err := client.VirtualAddresses()
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Virtual Address (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	if vas == nil {
		log.Printf("[WARN] VirtualAddress (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	for _, va = range vas.VirtualAddresses {
		if va.FullPath == name {
			break
		}
	}
	if va.FullPath != name {
		return diag.FromErr(fmt.Errorf("virtual address %s not found", name))
	}
	log.Printf("[DEBUG] virtual address configured on bigip is :%v", vas)
	_ = d.Set("name", name)
	_ = d.Set("arp", va.ARP)
	_ = d.Set("auto_delete", va.AutoDelete)
	_ = d.Set("conn_limit", va.ConnectionLimit)
	_ = d.Set("enabled", va.Enabled)
	_ = d.Set("icmp_echo", va.ICMPEcho)
	_ = d.Set("advertize_route", va.RouteAdvertisement)
	trafficGroup := va.TrafficGroup
	if trafficGroup == "none" {
		trafficGroup = fmt.Sprintf("/Common/%s", trafficGroup)
	}
	_ = d.Set("traffic_group", trafficGroup)

	return nil
}

func resourceBigipLtmVirtualAddressExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching virtual address " + name)

	var va *bigip.VirtualAddress
	vas, err := client.VirtualAddresses()
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Virtual Address  (%s) (%v) ", name, err)
		return false, err
	}
	for _, cand := range vas.VirtualAddresses {
		if cand.FullPath == name {
			va = &cand
			break
		}
	}

	if va == nil {
		d.SetId("")
	}

	return va != nil, nil
}

func resourceBigipLtmVirtualAddressUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	va := hydrateVirtualAddress(d)

	err := client.ModifyVirtualAddress(name, va)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Virtual Address  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}

	return resourceBigipLtmVirtualAddressRead(ctx, d, meta)
}

func hydrateVirtualAddress(d *schema.ResourceData) *bigip.VirtualAddress {
	return &bigip.VirtualAddress{
		Name:               d.Id(),
		ARP:                d.Get("arp").(bool),
		ConnectionLimit:    d.Get("conn_limit").(int),
		Enabled:            d.Get("enabled").(bool),
		ICMPEcho:           d.Get("icmp_echo").(string),
		RouteAdvertisement: d.Get("advertize_route").(string),
		TrafficGroup:       d.Get("traffic_group").(string),
		AutoDelete:         d.Get("auto_delete").(bool),
	}
}

func resourceBigipLtmVirtualAddressDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting virtual address " + name)
	client := meta.(*bigip.BigIP)
	vs, errCheck := resourceBigipLtmVirtualAddressExists(d, meta)

	if !vs {
		log.Printf("[ERROR] Unable to get Virtual Address  (%v)  (%v) ", vs, errCheck)
		d.SetId("")
		return nil
	}
	err := client.DeleteVirtualAddress(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Virtual Address  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
