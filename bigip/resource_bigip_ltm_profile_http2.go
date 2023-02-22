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
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipLtmProfileHttp2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmProfileHttp2Create,
		UpdateContext: resourceBigipLtmProfileHttp2Update,
		ReadContext:   resourceBigipLtmProfileHttp2Read,
		DeleteContext: resourceBigipLtmProfileHttp2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description:  "Name of Parent Http2 profile",
			},
			"concurrent_streams_per_connection": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the number of outstanding concurrent requests that are allowed on a single HTTP/2 connection. The default is 10",
			},
			"connection_idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the number of seconds that an HTTP/2 connection is idly left open before being shut down. The default is 300 seconds",
			},
			"insert_header": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
				Description:  "Specifies whether an HTTP header indicating the use of HTTP/2 should be inserted into the request that goes to the server. The default value is Disabled",
			},
			"insert_header_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the name of the HTTP header controlled by Insert Header. The default X-HTTP2.",
			},
			"enforce_tls_requirements": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
				Description:  "Enable or disable enforcement of TLS requirements,Default:enabled",
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
				Description:  "Enable to include content-length in HTTP/2 headers,Default : disabled",
			},
			"activation_modes": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether to enable all HTTP/2 modes, or only enable the Selected Modes listed in the Enabled column",
			},
		},
	}
}

func resourceBigipLtmProfileHttp2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	log.Printf("[INFO] Creating HTTP2 Profile:%+v ", name)

	pss := &bigip.Http2{
		Name: name,
	}
	config := getHttp2ProfileConfig(d, pss)

	err := client.CreateHttp2(config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating profile Http2 (%s): %s", name, err))
	}
	d.SetId(name)
	if !client.Teem {
		id := uuid.New()
		uniqueID := id.String()
		assetInfo := f5teem.AssetInfo{
			Name:    "Terraform-provider-bigip",
			Version: client.UserAgent,
			Id:      uniqueID,
		}
		apiKey := os.Getenv("TEEM_API_KEY")
		teemDevice := f5teem.AnonymousClient(assetInfo, apiKey)
		f := map[string]interface{}{
			"Terraform Version": client.UserAgent,
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err = teemDevice.Report(f, "bigip_ltm_profile_http2", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipLtmProfileHttp2Read(ctx, d, meta)
}

func resourceBigipLtmProfileHttp2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Updating HTTP2 Profile Profile:%+v ", name)
	pss := &bigip.Http2{
		Name: name,
	}
	config := getHttp2ProfileConfig(d, pss)

	err := client.ModifyHttp2(name, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error modifying profile Http2 (%s): %s ", name, err))
	}
	return resourceBigipLtmProfileHttp2Read(ctx, d, meta)
}

func resourceBigipLtmProfileHttp2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading http2 profile " + name)
	obj, err := client.GetHttp2(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve http2  (%s) (%v) ", name, err)

		return diag.FromErr(err)
	}
	if obj == nil {
		log.Printf("[WARN] Http2 Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("defaults_from", obj.DefaultsFrom)
	if _, ok := d.GetOk("concurrent_streams_per_connection"); ok {
		_ = d.Set("concurrent_streams_per_connection", obj.ConcurrentStreamsPerConnection)
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
		_ = d.Set("activation_modes", obj.ActivationModes)
	}
	return nil
}

func resourceBigipLtmProfileHttp2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Http2 Profile " + name)

	err := client.DeleteHttp2(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting  profile Http2 (%s): %s", name, err))
	}
	d.SetId("")
	return nil
}

func getHttp2ProfileConfig(d *schema.ResourceData, config *bigip.Http2) *bigip.Http2 {
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.ConcurrentStreamsPerConnection = d.Get("concurrent_streams_per_connection").(int)
	config.ConnectionIdleTimeout = d.Get("connection_idle_timeout").(int)
	config.HeaderTableSize = d.Get("header_table_size").(int)
	config.ActivationModes = setToStringSlice(d.Get("activation_modes").(*schema.Set))
	config.EnforceTLSRequirements = d.Get("enforce_tls_requirements").(string)
	config.FrameSize = d.Get("frame_size").(int)
	config.IncludeContentLength = d.Get("include_content_length").(string)
	config.InsertHeader = d.Get("insert_header").(string)
	config.InsertHeaderName = d.Get("insert_header_name").(string)
	config.ReceiveWindow = d.Get("receive_window").(int)
	config.WriteSize = d.Get("write_size").(int)

	return config
}
