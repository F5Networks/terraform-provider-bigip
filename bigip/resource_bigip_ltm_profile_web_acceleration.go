/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipLtmProfileWebAcceleration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmProfileWebAccelerationCreate,
		ReadContext:   resourceBigipLtmProfileWebAccelerationRead,
		UpdateContext: resourceBigipLtmProfileWebAccelerationUpdate,
		DeleteContext: resourceBigipLtmProfileWebAccelerationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the Web Acceleration profile",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.",
				ValidateFunc: validateF5Name,
			},
			"cache_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the maximum size for the cache. When the cache reaches the maximum size, the system starts removing the oldest entries. The default value is 100 megabytes.",
			},
			"cache_max_entries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the maximum number of entries that can be in the cache. The default value is 0 (zero), which means that the system does not limit the maximum entries.",
			},
			"cache_max_age": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how long the system considers the cached content to be valid. The default value is 3600 seconds.",
			},
			"cache_object_min_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the smallest object that the system considers eligible for caching. The default value is 500 bytes.",
			},
			"cache_object_max_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the largest object that the system considers eligible for caching. The default value is 50000 bytes.",
			},
			"cache_uri_exclude": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Configures a list of URIs to exclude from the cache. The default value of none specifies no URIs are excluded.",
			},
			"cache_uri_include": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Configures a list of URIs to include in the cache. The default value of .* specifies that all URIs are cacheable.",
			},
			"cache_uri_include_override": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Configures a list of URIs to include in the cache even if they would normally be excluded due to factors like object size or HTTP request type. The default value of none specifies no URIs are to be forced into the cache.",
			},
			"cache_uri_pinned": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Configures a list of URIs to keep in the cache. The pinning process keeps URIs in cache when they would normally be evicted to make room for more active URIs.",
			},
			"cache_client_cache_control_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies which cache disabling headers sent by clients the system ignores. The default value is all.",
			},
			"cache_insert_age_header": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Inserts Age and Date headers in the response. The default value is enabled.",
			},
			"cache_aging_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how quickly the system ages a cache entry. The aging rate ranges from 0 (slowest aging) to 10 (fastest aging). The default value is 9.",
			},
		},
	}
}

func resourceBigipLtmProfileWebAccelerationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating Profile Web Acceleration Service:%+v ", name)

	pss := &bigip.WebAccelerationProfileService{
		Name: name,
	}
	config := getHttpProfileWebAccelerationConfig(d, pss)

	err := client.AddWebAcceleration(config)
	if err != nil {
		return diag.FromErr(err)
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
		err = teemDevice.Report(f, "bigip_ltm_profile_web_acceleration", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipLtmProfileWebAccelerationRead(ctx, d, meta)
}

func resourceBigipLtmProfileWebAccelerationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching HTTP Profile Web Acceleration" + name)

	wap, err := client.GetWebAccelerationProfile(name)

	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Profile Web Acceleration  (%s) ", err)
		return diag.FromErr(err)
	}
	if wap == nil {
		log.Printf("[WARN] Web Acceleration Profile (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("defaults_from", wap.DefaultsFrom)

	if _, ok := d.GetOk("cache_size"); ok {
		_ = d.Set("cache_size", wap.CacheSize)
	}
	if _, ok := d.GetOk("cache_max_entries"); ok {
		_ = d.Set("cache_max_entries", wap.CacheMaxEntries)
	}
	if _, ok := d.GetOk("cache_max_age"); ok {
		_ = d.Set("cache_max_age", wap.CacheMaxAge)
	}
	if _, ok := d.GetOk("cache_object_min_size"); ok {
		_ = d.Set("cache_object_min_size", wap.CacheObjectMinSize)
	}
	if _, ok := d.GetOk("cache_object_max_size"); ok {
		_ = d.Set("cache_object_max_size", wap.CacheObjectMaxSize)
	}
	if _, ok := d.GetOk("cache_uri_exclude"); ok {
		_ = d.Set("cache_uri_exclude", wap.CacheUriExclude)
	}
	if _, ok := d.GetOk("cache_uri_include"); ok {
		_ = d.Set("cache_uri_include", wap.CacheUriInclude)
	}
	if _, ok := d.GetOk("cache_uri_include_override"); ok {
		_ = d.Set("cache_uri_include_override", wap.CacheUriIncludeOverride)
	}
	if _, ok := d.GetOk("cache_uri_pinned"); ok {
		_ = d.Set("cache_uri_pinned", wap.CacheUriPinned)
	}
	if _, ok := d.GetOk("cache_client_cache_control_mode"); ok {
		_ = d.Set("cache_client_cache_control_mode", wap.CacheClientCacheControlMode)
	}
	if _, ok := d.GetOk("cache_insert_age_header"); ok {
		_ = d.Set("cache_insert_age_header", wap.CacheInsertAgeHeader)
	}
	if _, ok := d.GetOk("cache_aging_rate"); ok {
		_ = d.Set("cache_aging_rate", wap.CacheAgingRate)
	}

	return nil
}

func resourceBigipLtmProfileWebAccelerationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating Profile Web Acceleration:%+v ", name)

	pss := &bigip.WebAccelerationProfileService{
		Name: name,
	}
	config := getHttpProfileWebAccelerationConfig(d, pss)

	err := client.ModifyWebAccelerationProfile(name, config)

	if err != nil {
		log.Printf("[ERROR] Unable to Modify HTTP Profile  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}

	return resourceBigipLtmProfileWebAccelerationRead(ctx, d, meta)

}

func resourceBigipLtmProfileWebAccelerationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Profile Web Acceleration " + name)
	err := client.DeleteWebAccelerationProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Profile Web Acceleration  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getHttpProfileWebAccelerationConfig(d *schema.ResourceData, config *bigip.WebAccelerationProfileService) *bigip.WebAccelerationProfileService {
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.CacheSize = d.Get("cache_size").(int)
	config.CacheMaxEntries = d.Get("cache_max_entries").(int)
	config.CacheMaxAge = d.Get("cache_max_age").(int)
	config.CacheObjectMinSize = d.Get("cache_object_min_size").(int)
	config.CacheObjectMaxSize = d.Get("cache_object_max_size").(int)
	config.CacheUriExclude = setToStringSlice(d.Get("cache_uri_exclude").(*schema.Set))
	config.CacheUriInclude = setToStringSlice(d.Get("cache_uri_include").(*schema.Set))
	config.CacheUriIncludeOverride = setToStringSlice(d.Get("cache_uri_include_override").(*schema.Set))
	config.CacheUriPinned = setToStringSlice(d.Get("cache_uri_pinned").(*schema.Set))
	config.CacheClientCacheControlMode = d.Get("cache_client_cache_control_mode").(string)
	config.CacheInsertAgeHeader = d.Get("cache_insert_age_header").(string)
	config.CacheAgingRate = d.Get("cache_aging_rate").(int)

	return config
}
