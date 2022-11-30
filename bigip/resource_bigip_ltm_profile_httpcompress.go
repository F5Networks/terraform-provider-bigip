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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipLtmProfileHttpcompress() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileHttpcompressCreate,
		Update: resourceBigipLtmProfileHttpcompressUpdate,
		Read:   resourceBigipLtmProfileHttpcompressRead,
		Delete: resourceBigipLtmProfileHttpcompressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
				Description:  "Name of the Httpcompress Profile",
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent Httpcompress profile",
			},
			"uri_exclude": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Servers Address",
			},
			"uri_include": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Servers Address",
			},
			"content_type_include": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Specifies a list of content types for compression of HTTP Content-Type responses. Use a string list to specify a list of content types you want to compress.",
			},
			"content_type_exclude": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Specifies a list of content types for compression of HTTP Content-Type responses. Use a string list to specify a list of content types you want to exclude.",
			},
			"compression_buffersize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the maximum number of compressed bytes that the system buffers before inserting a Content-Length header (which specifies the compressed size) into the response. The default is 4096 bytes.",
			},
			"gzip_compression_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the degree to which the system compresses the content. Higher compression levels cause the compression process to be slower. The default is 1 - Least Compression (Fastest)",
			},
			"gzip_memory_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144}),
				Description:  "Specifies the number of bytes of memory that the system uses for internal compression buffers when compressing a server response. The default is 8 kilobytes/8192 bytes.",
			},
			"gzip_window_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072}),
				Description:  "Specifies the number of kilobytes in the window size that the system uses when compressing a server response. The default is 16 kilobytes",
			},
			"keep_accept_encoding": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies, when checked (enabled), that the system does not remove the Accept-Encoding: header from an HTTP request. The default is disabled.",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"vary_header": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies, when checked (enabled), that the system inserts a Vary header into cacheable server responses. The default is enabled.",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"cpu_saver": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies, when checked (enabled), that the system monitors the percent CPU usage and adjusts compression rates automatically when the CPU usage reaches either the CPU Saver High Threshold or the CPU Saver Low Threshold. The default is enabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
		},
	}
}

func resourceBigipLtmProfileHttpcompressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	httpcompressConfig := &bigip.Httpcompress{
		Name: name,
	}
	htpcompProfileConfig := getHTTPCompressProfileConfig(d, httpcompressConfig)

	log.Println("[INFO] Creating Httpcompress profile")
	obj, _ := client.GetHttpcompress(name)
	if obj != nil && obj.FullPath == name {
		d.SetId(name)
		return resourceBigipLtmProfileHttpcompressRead(d, meta)
	}
	err := client.CreateHttpcompress(htpcompProfileConfig)

	if err != nil {
		return fmt.Errorf("Error retrieving profile Http compress (%s): %s ", name, err)
	}
	d.SetId(name)
	return resourceBigipLtmProfileHttpcompressRead(d, meta)
}

func resourceBigipLtmProfileHttpcompressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	httpcompressConfig := &bigip.Httpcompress{
		Name: name,
	}
	log.Println("[INFO] Updating Httpcompress profile")
	htpcompProfileConfig := getHTTPCompressProfileConfig(d, httpcompressConfig)

	err := client.ModifyHttpcompress(name, htpcompProfileConfig)
	if err != nil {
		return fmt.Errorf("Error modifying  profile Http compress (%s): %s", name, err)
	}
	return resourceBigipLtmProfileHttpcompressRead(d, meta)
}

func resourceBigipLtmProfileHttpcompressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetHttpcompress(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Http Compress Profile (%s) (%v)", name, err)
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Httpcompress Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", obj.FullPath)
	if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for Http Compress profile  (%s): %s", d.Id(), err)
	}
	if _, ok := d.GetOk("uri_include"); ok {
		if err := d.Set("uri_include", obj.UriInclude); err != nil {
			return fmt.Errorf("[DEBUG] Error saving UriInclude to state for  Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("uri_exclude"); ok {
		if err := d.Set("uri_exclude", obj.UriExclude); err != nil {
			return fmt.Errorf("[DEBUG] Error saving UriExclude to state for  Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("content_type_include"); ok {
		if err := d.Set("content_type_include", obj.ContentTypeInclude); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ContentTypeInclude to state for  Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("content_type_exclude"); ok {
		if err := d.Set("content_type_exclude", obj.ContentTypeExclude); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ContentTypeExclude to state for  Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("compression_buffersize"); ok {
		if err := d.Set("compression_buffersize", obj.BufferSize); err != nil {
			return fmt.Errorf("[DEBUG] Error saving compression_buffersize to state for Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("gzip_compression_level"); ok {
		if err := d.Set("gzip_compression_level", obj.GzipLevel); err != nil {
			return fmt.Errorf("[DEBUG] Error saving gzip_compression_level to state for Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("gzip_memory_level"); ok {
		if err := d.Set("gzip_memory_level", obj.GzipMemoryLevel); err != nil {
			return fmt.Errorf("[DEBUG] Error saving gzip_memory_level to state for Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("gzip_window_size"); ok {
		if err := d.Set("gzip_window_size", obj.GzipWindowSize); err != nil {
			return fmt.Errorf("[DEBUG] Error saving gzip_window_size to state for Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("keep_accept_encoding"); ok {
		if err := d.Set("keep_accept_encoding", obj.KeepAcceptEncoding); err != nil {
			return fmt.Errorf("[DEBUG] Error saving keep_accept_encoding to state for Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("vary_header"); ok {
		if err := d.Set("vary_header", obj.VaryHeader); err != nil {
			return fmt.Errorf("[DEBUG] Error saving vary_header to state for Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("cpu_saver"); ok {
		if err := d.Set("cpu_saver", obj.CPUSaver); err != nil {
			return fmt.Errorf("[DEBUG] Error saving cpu_saver to state for Http Compress profile  (%s): %s", d.Id(), err)
		}
	}
	return nil
}

func resourceBigipLtmProfileHttpcompressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Httpcompress Profile " + name)

	err := client.DeleteHttpcompress(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Httpcompress  (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

func getHTTPCompressProfileConfig(d *schema.ResourceData, config *bigip.Httpcompress) *bigip.Httpcompress {
	uriExclude := setToStringSlice(d.Get("uri_exclude").(*schema.Set))
	uriInclude := setToStringSlice(d.Get("uri_include").(*schema.Set))
	contentTypeInclude := setToStringSlice(d.Get("content_type_include").(*schema.Set))
	contentTypeExclude := setToStringSlice(d.Get("content_type_exclude").(*schema.Set))
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.UriExclude = uriExclude
	config.UriInclude = uriInclude
	config.ContentTypeInclude = contentTypeInclude
	config.ContentTypeExclude = contentTypeExclude
	config.BufferSize = d.Get("compression_buffersize").(int)
	config.GzipLevel = d.Get("gzip_compression_level").(int)
	config.GzipMemoryLevel = d.Get("gzip_memory_level").(int)
	config.GzipWindowSize = d.Get("gzip_window_size").(int)
	config.KeepAcceptEncoding = d.Get("keep_accept_encoding").(string)
	config.VaryHeader = d.Get("vary_header").(string)
	config.CPUSaver = d.Get("cpu_saver").(string)
	return config
}
