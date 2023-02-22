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
	"strconv"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipLtmPersistenceProfileDstAddr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmPersistenceProfileDstAddrCreate,
		ReadContext:   resourceBigipLtmPersistenceProfileDstAddrRead,
		UpdateContext: resourceBigipLtmPersistenceProfileDstAddrUpdate,
		DeleteContext: resourceBigipLtmPersistenceProfileDstAddrDelete,
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
				Optional: true,
				Computed: true,
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
				Description: "To enable _ disable match across pools with given persistence record",
				Computed:    true,
			},

			"match_across_services": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable match across services with given persistence record",
				Computed:    true,
			},

			"match_across_virtuals": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable match across services with given persistence record",
				Computed:    true,
			},

			"mirror": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable",
				Computed:    true,
			},

			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout for persistence of the session",
				Computed:    true,
			},

			"override_conn_limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable that pool member connection limits are overridden for persisted clients. Per-virtual connection limits remain hard limits and are not overridden.",
				Computed:    true,
			},

			// Specific to DestAddrPersistenceProfile
			"hash_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the hash algorithm",
				Computed:    true,
			},

			"mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identify a range of source IP addresses to manage together as a single source address affinity persistent connection when connecting to the pool. Must be a valid IPv4 or IPv6 mask.",
				Computed:    true,
			},
		},
	}
}

func resourceBigipLtmPersistenceProfileDstAddrCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	config := &bigip.PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}
	err := client.CreateDestAddrPersistenceProfile(config)
	if err != nil {
		log.Printf("[ERROR] Unable to create Dst Address Persistence profile %s  %v : ", name, err)
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceBigipLtmPersistenceProfileDstAddrUpdate(ctx, d, meta)
}

func resourceBigipLtmPersistenceProfileDstAddrRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Destination Address Persistence Profile " + name)

	pp, err := client.GetDestAddrPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve DestAdd Persistence Profile %s %v :", name, err)
		return diag.FromErr(err)
	}
	if pp == nil {
		log.Printf("[WARN] Destination Address Persistence Profile (%s) not found, removing from state", d.Id())
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
	if timeout, err := strconv.Atoi(pp.Timeout); err == nil {
		_ = d.Set("timeout", timeout)
	}

	if _, ok := d.GetOk("app_service"); ok {
		if err := d.Set("app_service", pp.AppService); err != nil {
			return diag.FromErr(fmt.Errorf("[DEBUG] Error saving AppService to state for resourceBigipLtmPersistenceProfileDstAddr (%s): %s", d.Id(), err))
		}
	}
	// Specific to DestAddrPersistenceProfile
	if _, ok := d.GetOk("hash_algorithm"); ok {
		if err := d.Set("hash_algorithm", pp.HashAlgorithm); err != nil {
			return diag.FromErr(fmt.Errorf("[DEBUG] Error saving HashAlgorithm to state for resourceBigipLtmPersistenceProfileDstAddr (%s): %s", d.Id(), err))
		}
	}
	if _, ok := d.GetOk("mask"); ok {
		_ = d.Set("mask", pp.Mask)
	}
	return nil
}

func resourceBigipLtmPersistenceProfileDstAddrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	timeout := d.Get("timeout").(int)
	if timeout != 0 {
		pp := &bigip.DestAddrPersistenceProfile{
			PersistenceProfile: bigip.PersistenceProfile{
				AppService:              d.Get("app_service").(string),
				DefaultsFrom:            d.Get("defaults_from").(string),
				MatchAcrossPools:        d.Get("match_across_pools").(string),
				MatchAcrossServices:     d.Get("match_across_services").(string),
				MatchAcrossVirtuals:     d.Get("match_across_virtuals").(string),
				Mirror:                  d.Get("mirror").(string),
				OverrideConnectionLimit: d.Get("override_conn_limit").(string),
				Timeout:                 strconv.Itoa(d.Get("timeout").(int)),
			},

			// Specific to DestAddrPersistenceProfile
			HashAlgorithm: d.Get("hash_algorithm").(string),
			Mask:          d.Get("mask").(string),
		}

		err := client.ModifyDestAddrPersistenceProfile(name, pp)
		if err != nil {
			log.Printf("[ERROR] Unable to Modify DestAdd Persistence Profile %s %v :", name, err)
			if errdel := client.DeleteDestAddrPersistenceProfile(name); errdel != nil {
				return diag.FromErr(errdel)
			}
			return diag.FromErr(err)
		}
	} else {
		pp := &bigip.DestAddrPersistenceProfile{
			PersistenceProfile: bigip.PersistenceProfile{
				AppService:              d.Get("app_service").(string),
				DefaultsFrom:            d.Get("defaults_from").(string),
				MatchAcrossPools:        d.Get("match_across_pools").(string),
				MatchAcrossServices:     d.Get("match_across_services").(string),
				MatchAcrossVirtuals:     d.Get("match_across_virtuals").(string),
				Mirror:                  d.Get("mirror").(string),
				OverrideConnectionLimit: d.Get("override_conn_limit").(string),
			},

			// Specific to DestAddrPersistenceProfile
			HashAlgorithm: d.Get("hash_algorithm").(string),
			Mask:          d.Get("mask").(string),
		}

		err := client.ModifyDestAddrPersistenceProfile(name, pp)
		if err != nil {
			log.Printf("[ERROR] Unable to Modify DestAdd Persistence Profile %s %v :", name, err)
			if errdel := client.DeleteDestAddrPersistenceProfile(name); errdel != nil {
				return diag.FromErr(errdel)
			}
			return diag.FromErr(err)
		}
	}

	return resourceBigipLtmPersistenceProfileDstAddrRead(ctx, d, meta)
}

func resourceBigipLtmPersistenceProfileDstAddrDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Destination Address Persistence Profile " + name)

	err := client.DeleteDestAddrPersistenceProfile(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error deleting DestAddPersistence profile  %s: %s", name, err))
	}
	d.SetId("")
	return nil
}
