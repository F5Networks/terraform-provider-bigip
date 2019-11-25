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

func resourceBigipLtmProfileHttp2() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileHttp2Create,
		Update: resourceBigipLtmProfileHttp2Update,
		Read:   resourceBigipLtmProfileHttp2Read,
		Delete: resourceBigipLtmProfileHttp2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Http2 Profile",
			},

			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"concurrent_streams_per_connection": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Use the parent Http2 profile",
			},

			"connection_idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Use the parent Http2 profile",
			},
			"header_table_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     4096,
				Description: "Use the parent Http2 profile",
			},

			"activation_modes": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},
		},
	}
}

func resourceBigipLtmProfileHttp2Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	concurrentStreamsPerConnection := d.Get("concurrent_streams_per_connection").(int)
	connectionIdleTimeout := d.Get("connection_idle_timeout").(int)
	headerTableSize := d.Get("header_table_size").(int)
	activationModes := setToStringSlice(d.Get("activation_modes").(*schema.Set))

	log.Println("[INFO] Creating Http2 profile")

	err := client.CreateHttp2(
		name,
		defaultsFrom,
		concurrentStreamsPerConnection,
		connectionIdleTimeout,
		headerTableSize,
		activationModes,
	)
	if err != nil {
		return fmt.Errorf("Error creating profile Http2 (%s): %s", name, err)
	}
	d.SetId(name)
	return resourceBigipLtmProfileHttp2Read(d, meta)
}

func resourceBigipLtmProfileHttp2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Http2{
		Name:                           name,
		DefaultsFrom:                   d.Get("defaults_from").(string),
		ConcurrentStreamsPerConnection: d.Get("concurrent_streams_per_connection").(int),
		ConnectionIdleTimeout:          d.Get("connection_idle_timeout").(int),
		HeaderTableSize:                d.Get("header_table_size").(int),
		ActivationModes:                setToStringSlice(d.Get("activation_modes").(*schema.Set)),
	}

	err := client.ModifyHttp2(name, r)
	if err != nil {
		return fmt.Errorf("Error modifying profile Http2 (%s): %s", name, err)
	}
	return resourceBigipLtmProfileHttp2Read(d, meta)
}

func resourceBigipLtmProfileHttp2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetHttp2(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve http2  (%s) (%v) ", name, err)

		return err
	}
	if obj == nil {
		log.Printf("[WARN] Http2 Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for Http2 profile  (%s): %s", d.Id(), err)
	}
	if err := d.Set("concurrent_streams_per_connection", obj.ConcurrentStreamsPerConnection); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ConcurrentStreamsPerConnection to state for Http2 profile  (%s): %s", d.Id(), err)
	}
	d.Set("connection_idle_timeout", obj.ConnectionIdleTimeout)
	if err := d.Set("activation_modes", obj.ActivationModes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ActivationModes to state for Http2 profile  (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceBigipLtmProfileHttp2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Http2 Profile " + name)

	err := client.DeleteHttp2(name)
	if err != nil {
		return fmt.Errorf("Error deleting  profile Http2 (%s): %s", name, err)
	}
	d.SetId("")
	return nil
}
