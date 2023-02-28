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

func resourceBigipLtmPersistenceProfileSSL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmPersistenceProfileSSLCreate,
		ReadContext:   resourceBigipLtmPersistenceProfileSSLRead,
		UpdateContext: resourceBigipLtmPersistenceProfileSSLUpdate,
		DeleteContext: resourceBigipLtmPersistenceProfileSSLDelete,
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
				Default:  "",
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
			},

			"override_conn_limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable that pool member connection limits are overridden for persisted clients. Per-virtual connection limits remain hard limits and are not overridden.",
				Computed:    true,
			},
		},
	}
}

func resourceBigipLtmPersistenceProfileSSLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	config := &bigip.PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	err := client.CreateSSLPersistenceProfile(config)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)

	return resourceBigipLtmPersistenceProfileSSLUpdate(ctx, d, meta)

}

func resourceBigipLtmPersistenceProfileSSLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching SSL Persistence Profile " + name)

	pp, err := client.GetSSLPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve SSL Persistence Profile  (%s) ", err)
		return diag.FromErr(err)
	}
	if pp == nil {
		log.Printf("[WARN] SSL  Persistence Profile (%s) not found, removing from state", name)
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
	return nil
}

func resourceBigipLtmPersistenceProfileSSLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	timeout := d.Get("timeout").(int)
	if timeout != 0 {
		pp := &bigip.SSLPersistenceProfile{
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
		}

		err := client.ModifySSLPersistenceProfile(name, pp)
		if err != nil {
			log.Printf("[ERROR] Unable to Modify SSL Persistence Profile  (%s) (%v)", name, err)
			if errdel := client.DeleteSSLPersistenceProfile(name); errdel != nil {
				return diag.FromErr(errdel)
			}
			return diag.FromErr(err)
		}
	} else {
		pp := &bigip.SSLPersistenceProfile{
			PersistenceProfile: bigip.PersistenceProfile{
				AppService:              d.Get("app_service").(string),
				DefaultsFrom:            d.Get("defaults_from").(string),
				MatchAcrossPools:        d.Get("match_across_pools").(string),
				MatchAcrossServices:     d.Get("match_across_services").(string),
				MatchAcrossVirtuals:     d.Get("match_across_virtuals").(string),
				Mirror:                  d.Get("mirror").(string),
				OverrideConnectionLimit: d.Get("override_conn_limit").(string),
			},
		}

		err := client.ModifySSLPersistenceProfile(name, pp)
		if err != nil {
			log.Printf("[ERROR] Unable to Modify SSL Persistence Profile  (%s) (%v)", name, err)
			if errdel := client.DeleteSSLPersistenceProfile(name); errdel != nil {
				return diag.FromErr(errdel)
			}
			return diag.FromErr(err)
		}

	}
	return resourceBigipLtmPersistenceProfileSSLRead(ctx, d, meta)
}

func resourceBigipLtmPersistenceProfileSSLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting SSL Persistence Profile " + name)
	err := client.DeleteSSLPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete SSL Persistence Profile  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
