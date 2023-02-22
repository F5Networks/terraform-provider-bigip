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
	"strconv"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipLtmPersistenceProfileCookie() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmPersistenceProfileCookieCreate,
		ReadContext:   resourceBigipLtmPersistenceProfileCookieRead,
		UpdateContext: resourceBigipLtmPersistenceProfileCookieUpdate,
		DeleteContext: resourceBigipLtmPersistenceProfileCookieDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the persistence profile",
				ValidateFunc: validateF5Name,
			},

			"app_service": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"defaults_from": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Inherit defaults from parent profile",
				ValidateFunc: validateF5Name,
			},

			"match_across_pools": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To enable _ disable match across pools with given persistence record",
			},

			"match_across_services": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To enable _ disable match across services with given persistence record",
			},

			"match_across_virtuals": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To enable _ disable match across virtual servers with given persistence record",
			},
			"method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the type of cookie processing that the system uses",
			},

			"mirror": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To enable _ disable",
			},

			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Timeout for persistence of the session",
			},

			"override_conn_limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To enable _ disable that pool member connection limits are overridden for persisted clients. Per-virtual connection limits remain hard limits and are not overridden.",
			},

			// Specific to CookiePersistenceProfile
			"always_send": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "To enable _ disable always sending cookies",
			},

			"cookie_encryption": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To required, preferred, or disabled policy for cookie encryption",
			},

			"cookie_encryption_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Passphrase for encrypted cookies",
			},

			"cookie_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the cookie to track persistence",
			},

			"expiration": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Expiration TTL for cookie specified in D:H:M:S or in seconds",
			},

			"hash_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Length of hash to apply to cookie",
			},

			"hash_offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of characters to skip in the cookie for the hash",
			},

			"httponly": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable sending only over http",
				Computed:    true,
			},
		},
	}
}

func resourceBigipLtmPersistenceProfileCookieCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	/*err := client.CreateCookiePersistenceProfile(
		name,
		parent,
	)*/
	config := &bigip.PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}
	err := client.CreateCookiePersistenceProfile(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Cookie Persistence Profile %s %v :", name, err)
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceBigipLtmPersistenceProfileCookieUpdate(ctx, d, meta)
}

func resourceBigipLtmPersistenceProfileCookieRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Cookie Persistence Profile " + name)

	pp, err := client.GetCookiePersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Cookie Persistence Profile %s  %v : ", name, err)
		return diag.FromErr(err)
	}
	if pp == nil {
		log.Printf("[WARN] Cookie Persistence Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("defaults_from", pp.DefaultsFrom)
	_ = d.Set("match_across_pools", pp.MatchAcrossPools)
	_ = d.Set("match_across_services", pp.MatchAcrossServices)
	_ = d.Set("match_across_virtuals", pp.MatchAcrossVirtuals)
	_ = d.Set("mirror", pp.Mirror)
	_ = d.Set("override_conn_limit", pp.OverrideConnectionLimit)
	_ = d.Set("httponly", pp.HTTPOnly)
	_ = d.Set("expiration", pp.Expiration)
	_ = d.Set("always_send", pp.AlwaysSend)
	_ = d.Set("hash_length", pp.HashLength)
	_ = d.Set("hash_offset", pp.HashOffset)
	_ = d.Set("method", pp.Method)
	if timeout, err := strconv.Atoi(pp.Timeout); err == nil {
		_ = d.Set("timeout", timeout)
	}

	if _, ok := d.GetOk("app_service"); ok {
		_ = d.Set("app_service", pp.AppService)
	}
	// Specific to CookiePersistenceProfile
	if _, ok := d.GetOk("cookie_encryption"); ok {
		_ = d.Set("cookie_encryption", pp.CookieEncryption)
	}
	if _, ok := d.GetOk("cookie_encryption_passphrase"); ok {
		_ = d.Set("cookie_encryption_passphrase", pp.CookieEncryptionPassphrase)
	}
	if _, ok := d.GetOk("cookie_name"); ok {
		_ = d.Set("cookie_name", pp.CookieName)
	}

	return nil
}

func resourceBigipLtmPersistenceProfileCookieUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	pp := &bigip.CookiePersistenceProfile{
		PersistenceProfile: bigip.PersistenceProfile{
			AppService:          d.Get("app_service").(string),
			DefaultsFrom:        d.Get("defaults_from").(string),
			MatchAcrossPools:    d.Get("match_across_pools").(string),
			MatchAcrossServices: d.Get("match_across_services").(string),
			MatchAcrossVirtuals: d.Get("match_across_virtuals").(string),
			Mirror:              d.Get("mirror").(string),
			//  Method:                  d.Get("method").(string),
			OverrideConnectionLimit: d.Get("override_conn_limit").(string),
			Timeout:                 strconv.Itoa(d.Get("timeout").(int)),
		},
		// Specific to CookiePersistenceProfile
		Method:                     d.Get("method").(string),
		AlwaysSend:                 d.Get("always_send").(string),
		CookieEncryption:           d.Get("cookie_encryption").(string),
		CookieEncryptionPassphrase: d.Get("cookie_encryption_passphrase").(string),
		CookieName:                 d.Get("cookie_name").(string),
		Expiration:                 d.Get("expiration").(string),
		HashLength:                 d.Get("hash_length").(int),
		HashOffset:                 d.Get("hash_offset").(int),
		HTTPOnly:                   d.Get("httponly").(string),
	}

	err := client.ModifyCookiePersistenceProfile(name, pp)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify Cookie Persistence Profile %s %v ", name, err)
		if errdel := client.DeleteCookiePersistenceProfile(name); errdel != nil {
			return diag.FromErr(errdel)
		}
		return diag.FromErr(err)
	}
	return resourceBigipLtmPersistenceProfileCookieRead(ctx, d, meta)
}

func resourceBigipLtmPersistenceProfileCookieDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Cookie Persistence Profile " + name)
	err := client.DeleteCookiePersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Cookie Persistence Profile %s  %v : ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
