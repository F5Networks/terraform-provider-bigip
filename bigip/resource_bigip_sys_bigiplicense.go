/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipSysBigiplicense() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysBigiplicenseCreate,
		UpdateContext: resourceBigipSysBigiplicenseUpdate,
		ReadContext:   resourceBigipSysBigiplicenseRead,
		DeleteContext: resourceBigipSysBigiplicenseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"command": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tmsh command to execute tmsh commands like install",
			},
			"registration_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique Key F5 provides for Licensing BIG-IP",
			},
		},
	}
}

func resourceBigipSysBigiplicenseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	command := d.Get("command").(string)
	registrationKey := d.Get("registration_key").(string)
	log.Println("[INFO] Creating BigipLicense ")

	err := client.CreateBigiplicense(
		command,
		registrationKey,
	)
	time.Sleep(300 * time.Second)
	if err != nil {
		log.Printf("[ERROR] Unable to Apply License to Bigip  (%v) ", err)
		return diag.FromErr(err)
	}
	d.SetId(registrationKey)
	return resourceBigipSysBigiplicenseRead(ctx, d, meta)
}

func resourceBigipSysBigiplicenseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	registrationKey := d.Id()

	log.Println("[INFO] Updating Bigiplicense " + registrationKey)

	r := &bigip.Bigiplicense{
		Registration_key: registrationKey,
		Command:          d.Get("command").(string),
	}
	err := client.ModifyBigiplicense(r)
	if err != nil {
		log.Printf("[ERROR] Unable to Apply License to Bigip  (%v) ", err)
		return diag.FromErr(err)
	}
	return resourceBigipSysBigiplicenseRead(ctx, d, meta)
}

func resourceBigipSysBigiplicenseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Bigiplicense " + name)

	licenses, err := client.Bigiplicenses()
	if err != nil {
		log.Printf("[ERROR] Unable to Read License from Bigip  (%v) ", err)
		return diag.FromErr(err)
	}
	if licenses == nil {
		log.Printf("[WARN] License (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceBigipSysBigiplicenseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// API does not Exists
	return nil
}
