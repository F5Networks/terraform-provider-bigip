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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func resourceBigipSysNtpCreate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	d.SetId(description)
	return resourceBigipSysNtpRead(d, meta)
}

func resourceBigipSysNtpUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	return resourceBigipSysNtpRead(d, meta)
}

func resourceBigipSysNtpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading NTP Config" + description)

	ntp, err := client.NTPs()
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve NTP Config (%s) ", err)
		return err
	}
	if ntp == nil {
		log.Printf("[WARN] NTP Config (%s) not found, removing from state", d.Id())
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
		return err
	}
	d.SetId("")
	return nil
}

func getSysNTPConfig(d *schema.ResourceData, config *bigip.NTP) *bigip.NTP {
	config.Servers = listToStringSlice(d.Get("servers").([]interface{}))
	config.Timezone = d.Get("timezone").(string)
	return config
}
