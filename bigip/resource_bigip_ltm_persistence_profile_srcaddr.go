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
	"strconv"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmPersistenceProfileSrcAddr() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPersistenceProfileSrcAddrCreate,
		Read:   resourceBigipLtmPersistenceProfileSrcAddrRead,
		Update: resourceBigipLtmPersistenceProfileSrcAddrUpdate,
		Delete: resourceBigipLtmPersistenceProfileSrcAddrDelete,
		Exists: resourceBigipLtmPersistenceProfileSrcAddrExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

			// Specific to SourceAddrPersistenceProfile
			"hash_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the hash algorithm",
				Computed:    true,
			},

			"map_proxies": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable directs all to the same single pool member",
				Computed:    true,
			},

			"mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Identify a range of source IP addresses to manage together as a single source address affinity persistent connection when connecting to the pool. Must be a valid IPv4 or IPv6 mask.",
			},
		},
	}
}

func resourceBigipLtmPersistenceProfileSrcAddrCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	config := &bigip.PersistenceProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	err := client.CreateSourceAddrPersistenceProfile(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Source Address Persistence Profile  (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)

	err = resourceBigipLtmPersistenceProfileSrcAddrUpdate(d, meta)
	if err != nil {
		if errdel := client.DeleteSourceAddrPersistenceProfile(name); errdel != nil {
			return errdel
		}
		return err
	}

	return resourceBigipLtmPersistenceProfileSrcAddrRead(d, meta)

}

func resourceBigipLtmPersistenceProfileSrcAddrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Source Address Persistence Profile " + name)

	pp, err := client.GetSourceAddrPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Source Address Persistence Profile  (%s)(%v) ", name, err)
		return err
	}
	if pp == nil {
		log.Printf("[WARN] Source Address Persistence Profile (%s) not found, removing from state", d.Id())
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
		d.Set("timeout", timeout)
	}

	if _, ok := d.GetOk("app_service"); ok {
		if err := d.Set("app_service", pp.AppService); err != nil {
			return fmt.Errorf("[DEBUG] Error saving AppService to state for PersistenceProfileSrcAddr (%s): %s", d.Id(), err)
		}
	}

	// Specific to SourceAddrPersistenceProfile
	if _, ok := d.GetOk("hash_algorithm"); ok {
		if err := d.Set("hash_algorithm", pp.HashAlgorithm); err != nil {
			return fmt.Errorf("[DEBUG] Error saving HashAlgorithm to state for PersistenceProfileSrcAddr (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("map_proxies"); ok {
		_ = d.Set("map_proxies", pp.MapProxies)
	}
	if _, ok := d.GetOk("mask"); ok {
		_ = d.Set("mask", pp.Mask)
	}

	return nil
}

func resourceBigipLtmPersistenceProfileSrcAddrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	timeout := d.Get("timeout").(int)
	if timeout != 0 {
		pp := &bigip.SourceAddrPersistenceProfile{
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

			// Specific to SourceAddrPersistenceProfile
			HashAlgorithm: d.Get("hash_algorithm").(string),
			MapProxies:    d.Get("map_proxies").(string),
			Mask:          d.Get("mask").(string),
		}
		err := client.ModifySourceAddrPersistenceProfile(name, pp)
		if err != nil {
			log.Printf("[ERROR] Unable to Modify Source Address Persistence Profile  (%s) ", err)
			return err
		}
	} else {
		pp := &bigip.SourceAddrPersistenceProfile{
			PersistenceProfile: bigip.PersistenceProfile{
				AppService:              d.Get("app_service").(string),
				DefaultsFrom:            d.Get("defaults_from").(string),
				MatchAcrossPools:        d.Get("match_across_pools").(string),
				MatchAcrossServices:     d.Get("match_across_services").(string),
				MatchAcrossVirtuals:     d.Get("match_across_virtuals").(string),
				Mirror:                  d.Get("mirror").(string),
				OverrideConnectionLimit: d.Get("override_conn_limit").(string),
			},

			// Specific to SourceAddrPersistenceProfile
			HashAlgorithm: d.Get("hash_algorithm").(string),
			MapProxies:    d.Get("map_proxies").(string),
			Mask:          d.Get("mask").(string),
		}
		err := client.ModifySourceAddrPersistenceProfile(name, pp)
		if err != nil {
			log.Printf("[ERROR] Unable to Modify Source Address Persistence Profile  (%s) ", err)
			return err
		}
	}

	return resourceBigipLtmPersistenceProfileSrcAddrRead(d, meta)
}

func resourceBigipLtmPersistenceProfileSrcAddrDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Source Address Persistence Profile " + name)
	err := client.DeleteSourceAddrPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Source Address Persistence Profile (%s)  (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

func resourceBigipLtmPersistenceProfileSrcAddrExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Source Address Persistence Profile " + name)

	pp, err := client.GetSourceAddrPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Source Address Persistence Profile  (%s) (%v)", name, err)
		return false, err
	}

	if pp == nil {
		log.Printf("[WARN] persistence profile src_addr  (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return pp != nil, nil
}
