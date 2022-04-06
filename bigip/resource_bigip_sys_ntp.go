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

func resourceBigipSysNtp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysNtpCreate,
		Update: resourceBigipSysNtpUpdate,
		Read:   resourceBigipSysNtpRead,
		Delete: resourceBigipSysNtpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the ntp Servers",
				ValidateFunc: validateF5Name,
			},

			"servers": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},

			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Servers timezone",
			},
		},
	}

}

func resourceBigipSysNtpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Get("description").(string)
	servers := setToStringSlice(d.Get("servers").(*schema.Set))
	timezone := d.Get("timezone").(string)

	log.Println("[INFO] Configuring Ntp ")

	err := client.CreateNTP(
		description,
		servers,
		timezone,
	)

	if err != nil {
		log.Printf("[ERROR] Unable to Configure  NTP   (%s) ", err)
		return err
	}
	d.SetId(description)
	return resourceBigipSysNtpRead(d, meta)
}

func resourceBigipSysNtpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Updating NTP " + description)

	r := &bigip.NTP{
		Description: description,
		Servers:     setToStringSlice(d.Get("servers").(*schema.Set)),
		Timezone:    d.Get("timezone").(string),
	}

	err := client.ModifyNTP(r)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify  NTP  (%v) ", err)
		return err
	}
	return resourceBigipSysNtpRead(d, meta)
}

func resourceBigipSysNtpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading NTP " + description)

	ntp, err := client.NTPs()
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve NTP   (%s) ", err)
		return err
	}
	if ntp == nil {
		log.Printf("[WARN] NTP (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("description", ntp.Description); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Description to state for NTP  (%s): %s", d.Id(), err)
	}
	if err := d.Set("servers", ntp.Servers); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Servers to state for NTP  (%s): %s", d.Id(), err)
	}
	if err := d.Set("timezone", ntp.Timezone); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Timezone  state for NTP  (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceBigipSysNtpDelete(d *schema.ResourceData, meta interface{}) error {
	/* This function is not supported on BIG-IP, you cannot DELETE NTP API is not supported */
	return nil
}
