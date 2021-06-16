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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5Name,
				Description:  "Name of the Http2 Profile",
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateF5Name,
				Description:  "Use the parent Http2 profile",
			},
			"concurrent_streams_per_connection": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of concurrent connections to allow on a single HTTP/2 connection.Default is 10",
			},
			"connection_idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of seconds that a HTTP/2 connection is left open idly before it is closed",
			},
			"insert_header": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This setting specifies whether the BIG-IP system should add an HTTP header to the HTTP request to show that the request was received over HTTP/2,Default:disabled",
			},
			"insert_header_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This setting specifies the name of the header that the BIG-IP system will add to the HTTP request when the Insert Header is enabled.",
			},
			"enforce_tls_requirements": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable enforcement of TLS requirements,Default:enabled",
			},
			"header_table_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The size of the header table, in KB, for the HTTP headers that the HTTP/2 protocol compresses to save bandwidth.Default: 4096",
			},
			"frame_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The size of the data frames, in bytes, that the HTTP/2 protocol sends to the client. Default: 2048",
			},
			"receive_window": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The flow-control size for upload streams, in KB. Default: 32",
			},
			"write_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The total size of combined data frames, in bytes, that the HTTP/2 protocol sends in a single write function. Default: 16384",
			},
			"include_content_length": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enable to include content-length in HTTP/2 headers,Default : disabled",
			},
			"activation_modes": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "This setting specifies the condition that will cause the BIG-IP system to handle an incoming connection as an HTTP/2 connection.",
			},
		},
	}
}

func resourceBigipLtmProfileHttp2Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	//concurrentStreamsPerConnection := d.Get("concurrent_streams_per_connection").(int)
	//connectionIdleTimeout := d.Get("connection_idle_timeout").(int)
	//headerTableSize := d.Get("header_table_size").(int)
	//activationModes := setToStringSlice(d.Get("activation_modes").(*schema.Set))
	log.Println("[INFO] Creating Http2 profile")
	r := &bigip.Http2{
		Name:                           name,
		DefaultsFrom:                   defaultsFrom,
		ConcurrentStreamsPerConnection: d.Get("concurrent_streams_per_connection").(int),
		ConnectionIdleTimeout:          d.Get("connection_idle_timeout").(int),
		HeaderTableSize:                d.Get("header_table_size").(int),
		ActivationModes:                setToStringSlice(d.Get("activation_modes").(*schema.Set)),
		EnforceTLSRequirements:         d.Get("enforce_tls_requirements").(string),
		FrameSize:                      d.Get("frame_size").(int),
		IncludeContentLength:           d.Get("include_content_length").(string),
		InsertHeader:                   d.Get("insert_header").(string),
		InsertHeaderName:               d.Get("insert_header_name").(string),
		ReceiveWindow:                  d.Get("receive_window").(int),
		WriteSize:                      d.Get("write_size").(int),
	}
	err := client.CreateHttp2(r)
	if err != nil {
		return fmt.Errorf("Error creating profile Http2 (%s): %s ", name, err)
	}
	d.SetId(name)
	return resourceBigipLtmProfileHttp2Read(d, meta)
}

func resourceBigipLtmProfileHttp2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating http2 profile " + name)
	r := &bigip.Http2{
		Name:                           name,
		DefaultsFrom:                   d.Get("defaults_from").(string),
		ConcurrentStreamsPerConnection: d.Get("concurrent_streams_per_connection").(int),
		ConnectionIdleTimeout:          d.Get("connection_idle_timeout").(int),
		HeaderTableSize:                d.Get("header_table_size").(int),
		ActivationModes:                setToStringSlice(d.Get("activation_modes").(*schema.Set)),
		EnforceTLSRequirements:         d.Get("enforce_tls_requirements").(string),
		FrameSize:                      d.Get("frame_size").(int),
		IncludeContentLength:           d.Get("include_content_length").(string),
		InsertHeader:                   d.Get("insert_header").(string),
		InsertHeaderName:               d.Get("insert_header_name").(string),
		ReceiveWindow:                  d.Get("receive_window").(int),
		WriteSize:                      d.Get("write_size").(int),
	}
	err := client.ModifyHttp2(name, r)
	if err != nil {
		return fmt.Errorf("Error modifying profile Http2 (%s): %s ", name, err)
	}
	return resourceBigipLtmProfileHttp2Read(d, meta)
}

func resourceBigipLtmProfileHttp2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading http2 profile " + name)
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
	_ = d.Set("name", name)
	_ = d.Set("defaults_from", obj.DefaultsFrom)
	if _, ok := d.GetOk("concurrent_streams_per_connection"); ok {
		if err := d.Set("concurrent_streams_per_connection", obj.ConcurrentStreamsPerConnection); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ConcurrentStreamsPerConnection to state for Http2 profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("connection_idle_timeout"); ok {
		_ = d.Set("connection_idle_timeout", obj.ConnectionIdleTimeout)
	}
	if _, ok := d.GetOk("header_table_size"); ok {
		_ = d.Set("header_table_size", obj.HeaderTableSize)
	}
	if _, ok := d.GetOk("enforce_tls_requirements"); ok {
		_ = d.Set("enforce_tls_requirements", obj.EnforceTLSRequirements)
	}
	if _, ok := d.GetOk("frame_size"); ok {
		_ = d.Set("frame_size", obj.FrameSize)
	}
	if _, ok := d.GetOk("receive_window"); ok {
		_ = d.Set("receive_window", obj.ReceiveWindow)
	}
	if _, ok := d.GetOk("write_size"); ok {
		_ = d.Set("write_size", obj.WriteSize)
	}
	if _, ok := d.GetOk("insert_header"); ok {
		_ = d.Set("insert_header", obj.InsertHeader)
	}
	if _, ok := d.GetOk("insert_header_name"); ok {
		_ = d.Set("insert_header_name", obj.InsertHeaderName)
	}
	if _, ok := d.GetOk("connection_idle_timeout"); ok {
		if err := d.Set("activation_modes", obj.ActivationModes); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ActivationModes to state for Http2 profile  (%s): %s", d.Id(), err)
		}
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
