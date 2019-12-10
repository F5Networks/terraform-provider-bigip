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

	"github.com/f5devcentral/go-bigip"
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "To enable _ disable match across pools with given persistence record",
				ValidateFunc: validateEnabledDisabled,
			},

			"match_across_services": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "To enable _ disable match across services with given persistence record",
				ValidateFunc: validateEnabledDisabled,
			},

			"match_across_virtuals": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "To enable _ disable match across services with given persistence record",
				ValidateFunc: validateEnabledDisabled,
			},

			"mirror": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "To enable _ disable",
				ValidateFunc: validateEnabledDisabled,
			},

			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout for persistence of the session",
			},

			"override_conn_limit": {
				Type:         schema.TypeString,
				Default:      false,
				Optional:     true,
				Description:  "To enable _ disable that pool member connection limits are overridden for persisted clients. Per-virtual connection limits remain hard limits and are not overridden.",
				ValidateFunc: validateEnabledDisabled,
			},

			// Specific to SourceAddrPersistenceProfile
			"hash_algorithm": {
				Type:        schema.TypeString,
				Default:     "default",
				Optional:    true,
				Description: "Specify the hash algorithm",
			},

			"map_proxies": {
				Type:         schema.TypeString,
				Default:      true,
				Optional:     true,
				Description:  "To enable _ disable directs all to the same single pool member",
				ValidateFunc: validateEnabledDisabled,
			},

			"mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identify a range of source IP addresses to manage together as a single source address affinity persistent connection when connecting to the pool. Must be a valid IPv4 or IPv6 mask.",
			},
		},
	}
}

func resourceBigipLtmPersistenceProfileSrcAddrCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	err := client.CreateSourceAddrPersistenceProfile(
		name,
		parent,
	)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Source Address Persistence Profile  (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)

	err = resourceBigipLtmPersistenceProfileSrcAddrUpdate(d, meta)
	if err != nil {
		client.DeleteSourceAddrPersistenceProfile(name)
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
	d.Set("name", name)
	if err := d.Set("app_service", pp.AppService); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AppService to state for PersistenceProfileSrcAddr (%s): %s", d.Id(), err)
	}
	if err := d.Set("defaults_from", pp.DefaultsFrom); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for PersistenceProfileSrcAddr (%s): %s", d.Id(), err)
	}
	d.Set("match_across_pools", pp.MatchAcrossPools)
	d.Set("match_across_services", pp.MatchAcrossServices)
	d.Set("match_across_virtuals", pp.MatchAcrossVirtuals)
	if err := d.Set("mirror", pp.Mirror); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mirror to state for PersistenceProfileSrcAddr (%s): %s", d.Id(), err)
	}
	d.Set("timeout", pp.Timeout)
	d.Set("override_conn_limit", pp.OverrideConnectionLimit)

	// Specific to SourceAddrPersistenceProfile
	if err := d.Set("hash_algorithm", pp.HashAlgorithm); err != nil {
		return fmt.Errorf("[DEBUG] Error saving HashAlgorithm to state for PersistenceProfileSrcAddr (%s): %s", d.Id(), err)
	}
	d.Set("map_proxies", pp.MapProxies)
	d.Set("mask", pp.Mask)

	return nil
}

func resourceBigipLtmPersistenceProfileSrcAddrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

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
		log.Printf("[WARN] persistance profile src_addr  (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return pp != nil, nil
}
