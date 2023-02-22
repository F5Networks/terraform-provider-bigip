/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipSysNtp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysNtpCreate,
		UpdateContext: resourceBigipSysNtpUpdate,
		ReadContext:   resourceBigipSysNtpRead,
		DeleteContext: resourceBigipSysNtpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User defined description.",
				//ValidateFunc: validateF5Name,
			},
			"servers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Specifies the time servers that the system uses to update the system time",
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the time zone that you want to use for the system time",
			},
		},
	}

}

func resourceBigipSysNtpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Get("description").(string)

	log.Println("[INFO] Configuring NTP Servers ")

	configSysNTP := &bigip.NTP{
		Description: description,
	}
	sysNTPConfig := getSysNTPConfig(d, configSysNTP)
	err := client.ModifyNTP(sysNTPConfig)

	if err != nil {
		log.Printf("[ERROR] Unable to Configure  NTP Servers  (%s) ", err)
		return diag.FromErr(err)
	}
	d.SetId(description)
	return resourceBigipSysNtpRead(ctx, d, meta)
}

func resourceBigipSysNtpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Updating NTP Servers" + description)

	configSysNTP := &bigip.NTP{
		Description: description,
	}
	sysNTPConfig := getSysNTPConfig(d, configSysNTP)
	err := client.ModifyNTP(sysNTPConfig)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify  NTP Servers (%v) ", err)
		return diag.FromErr(err)
	}
	return resourceBigipSysNtpRead(ctx, d, meta)
}

func resourceBigipSysNtpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading NTP Config" + description)

	ntp, err := client.NTPs()
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve NTP Config (%s) ", err)
		return diag.FromErr(err)
	}
	if ntp == nil {
		log.Printf("[WARN] NTP Config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("description", ntp.Description)
	_ = d.Set("servers", ntp.Servers)
	_ = d.Set("timezone", ntp.Timezone)

	return nil
}

func resourceBigipSysNtpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	description := d.Id()
	log.Println("[INFO] Deleting System NTP Config:" + description)
	configSysNTP := &bigip.NTP{
		Description: description,
		Servers:     []string{},
		Timezone:    "America/Los_Angeles",
	}
	err := client.ModifyNTP(configSysNTP)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete NTP Config (%s) (%v) ", description, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getSysNTPConfig(d *schema.ResourceData, config *bigip.NTP) *bigip.NTP {
	config.Servers = listToStringSlice(d.Get("servers").([]interface{}))
	config.Timezone = d.Get("timezone").(string)
	return config
}
