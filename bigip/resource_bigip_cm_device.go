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

func resourceBigipCmDevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipCmDeviceCreate,
		UpdateContext: resourceBigipCmDeviceUpdate,
		ReadContext:   resourceBigipCmDeviceRead,
		DeleteContext: resourceBigipCmDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{

			"configsync_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address used for config sync",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Device which needs to be Deviceensed",
			},

			"mirror_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address used for state mirroring",
			},
			"mirror_secondary_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secondary IP address used for state mirroring",
			},
		},
	}

}

func resourceBigipCmDeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	configsyncIp := d.Get("configsync_ip").(string)
	name := d.Get("name").(string)
	mirrorIp := d.Get("mirror_ip").(string)
	mirrorSecondaryIp := d.Get("mirror_secondary_ip").(string)

	log.Println("[INFO] Creating Device ")

	err := client.CreateDevice(
		name,
		configsyncIp,
		mirrorIp,
		mirrorSecondaryIp,
	)

	if err != nil {
		log.Printf("[ERROR] Unable to Create Device %s %v ", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipCmDeviceRead(ctx, d, meta)

}

func resourceBigipCmDeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Device " + name)

	r := &bigip.Device{
		Name:              name,
		ConfigsyncIp:      d.Get("configsync_ip").(string),
		MirrorIp:          d.Get("mirror_ip").(string),
		MirrorSecondaryIp: d.Get("mirror_secondary_ip").(string),
	}

	err := client.ModifyDevice(r)
	if err != nil {
		log.Printf("[ERROR] Unable to Modidy Device (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	return resourceBigipCmDeviceRead(ctx, d, meta)
}

func resourceBigipCmDeviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Device " + name)

	members, err := client.Devices(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Device (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	if members == nil {
		log.Printf("[WARN] Device (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("name", members.Name)

	if err := d.Set("mirror_ip", members.MirrorIp); err != nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Error saving mirror_ip  to state for Device (%s): %s", d.Id(), err))
	}

	if err := d.Set("configsync_ip", members.ConfigsyncIp); err != nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Error saving configsync_ip  to state for Device (%s): %s", d.Id(), err))
	}

	if err := d.Set("mirror_secondary_ip", members.MirrorSecondaryIp); err != nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Error saving mirror_secondary_ip  to state for Device (%s): %s", d.Id(), err))
	}

	return nil
}

func resourceBigipCmDeviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeleteDevice(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Device (%s)  (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
