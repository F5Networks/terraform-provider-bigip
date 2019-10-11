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
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmPersistenceProfileDstAddr() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPersistenceProfileDstAddrCreate,
		Read:   resourceBigipLtmPersistenceProfileDstAddrRead,
		Update: resourceBigipLtmPersistenceProfileDstAddrUpdate,
		Delete: resourceBigipLtmPersistenceProfileDstAddrDelete,
		Exists: resourceBigipLtmPersistenceProfileDstAddrExists,
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

			// Specific to DestAddrPersistenceProfile
			"hash_algorithm": {
				Type:        schema.TypeString,
				Default:     "default",
				Optional:    true,
				Description: "Specify the hash algorithm",
			},

			"mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identify a range of source IP addresses to manage together as a single source address affinity persistent connection when connecting to the pool. Must be a valid IPv4 or IPv6 mask.",
			},
		},
	}
}

func resourceBigipLtmPersistenceProfileDstAddrCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	err := client.CreateDestAddrPersistenceProfile(
		name,
		parent,
	)
	if err != nil {
		log.Printf("[ERROR] Unable to create Dst Address Persistence profile %s  %v : ", name, err)
		return err
	}

	d.SetId(name)

	err = resourceBigipLtmPersistenceProfileDstAddrUpdate(d, meta)
	if err != nil {
		client.DeleteDestAddrPersistenceProfile(name)
		return err
	}

	return resourceBigipLtmPersistenceProfileDstAddrRead(d, meta)

}

func resourceBigipLtmPersistenceProfileDstAddrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Destination Address Persistence Profile " + name)

	pp, err := client.GetDestAddrPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrive DestAdd Persistence Profile %s %v :", name, err)
		return err
	}
	if pp == nil {
		log.Printf("[WARN] Destination Address Persistence Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	if err := d.Set("app_service", pp.AppService); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AppService to state for resourceBigipLtmPersistenceProfileDstAddr (%s): %s", d.Id(), err)
	}
	if err := d.Set("defaults_from", pp.DefaultsFrom); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for resourceBigipLtmPersistenceProfileDstAddr (%s): %s", d.Id(), err)
	}
	d.Set("match_across_pools", pp.MatchAcrossPools)
	d.Set("match_across_services", pp.MatchAcrossServices)
	d.Set("match_across_virtuals", pp.MatchAcrossVirtuals)
	if err := d.Set("mirror", pp.Mirror); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mirror to state for resourceBigipLtmPersistenceProfileDstAddr (%s): %s", d.Id(), err)
	}
	d.Set("timeout", pp.Timeout)
	d.Set("override_conn_limit", pp.OverrideConnectionLimit)

	// Specific to DestAddrPersistenceProfile
	if err := d.Set("hash_algorithm", pp.HashAlgorithm); err != nil {
		return fmt.Errorf("[DEBUG] Error saving HashAlgorithm to state for resourceBigipLtmPersistenceProfileDstAddr (%s): %s", d.Id(), err)
	}
	d.Set("mask", pp.Mask)

	return nil
}

func resourceBigipLtmPersistenceProfileDstAddrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

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
		return err
	}
	return resourceBigipLtmPersistenceProfileDstAddrRead(d, meta)
}

func resourceBigipLtmPersistenceProfileDstAddrDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Destination Address Persistence Profile " + name)

	err := client.DeleteDestAddrPersistenceProfile(name)
	if err != nil {
		return fmt.Errorf("Error deleting DestAddPersistence profile  %s: %s", name, err)
	}
	d.SetId("")
	return nil
}

func resourceBigipLtmPersistenceProfileDstAddrExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Destination Address Persistence Profile " + name)

	pp, err := client.GetDestAddrPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrive Destination Address Persistence Profile  (%s) ", err)
		return false, err
	}

	if pp == nil {
		log.Printf("[WARN] DestAddpersistance profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, nil
	}

	return pp != nil, nil
}
