/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"log"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipSysBigiplicense() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysBigiplicenseCreate,
		Update: resourceBigipSysBigiplicenseUpdate,
		Read:   resourceBigipSysBigiplicenseRead,
		Delete: resourceBigipSysBigiplicenseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"timeout": {
				Type:     schema.TypeInt,
				Default:  300,
				Optional: true,
			},
		},
	}
}
func resourceBigipSysBigiplicenseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	command := d.Get("command").(string)
	registrationKey := d.Get("registration_key").(string)
	log.Println("[INFO] Creating Bigip License ")

	err := client.CreateBigiplicense(
		command,
		registrationKey,
	)
	timeOut := time.Duration(d.Get("timeout").(int)) * time.Second
	time.Sleep(timeOut)
	if err != nil {
		log.Printf("[ERROR] Unable to Apply License to Bigip  (%v) ", err)
		return err
	}
	d.SetId(registrationKey)
	return resourceBigipSysBigiplicenseRead(d, meta)
}

func resourceBigipSysBigiplicenseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	registrationKey := d.Id()

	log.Println("[INFO] Updating Bigiplicense " + registrationKey)

	r := &bigip.Bigiplicense{
		Registration_key: registrationKey,
		Command:          d.Get("command").(string),
	}

	return client.ModifyBigiplicense(r)
}

func resourceBigipSysBigiplicenseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Bigiplicense " + name)

	licenses, err := client.Bigiplicenses()
	if err != nil {
		log.Printf("[ERROR] Unable to Read License from Bigip  (%v) ", err)
		d.SetId("")
		return err
	}
	if licenses == nil {
		log.Printf("[WARN] License (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceBigipSysBigiplicenseDelete(d *schema.ResourceData, meta interface{}) error {
	// API does not Exists
	return nil
}
