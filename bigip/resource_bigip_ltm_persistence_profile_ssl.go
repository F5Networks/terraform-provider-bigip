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

func resourceBigipLtmPersistenceProfileSSL() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPersistenceProfileSSLCreate,
		Read:   resourceBigipLtmPersistenceProfileSSLRead,
		Update: resourceBigipLtmPersistenceProfileSSLUpdate,
		Delete: resourceBigipLtmPersistenceProfileSSLDelete,
		Exists: resourceBigipLtmPersistenceProfileSSLExists,
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
		},
	}
}

func resourceBigipLtmPersistenceProfileSSLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	err := client.CreateSSLPersistenceProfile(
		name,
		parent,
	)
	if err != nil {
		return err
	}
	d.SetId(name)

	err = resourceBigipLtmPersistenceProfileSSLUpdate(d, meta)
	if err != nil {
		client.DeleteSSLPersistenceProfile(name)
		return err
	}

	return resourceBigipLtmPersistenceProfileSSLRead(d, meta)

}

func resourceBigipLtmPersistenceProfileSSLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching SSL Persistence Profile " + name)

	pp, err := client.GetSSLPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrive SSL Persistence Profile  (%s) ", err)
		return err
	}
	if pp == nil {
		log.Printf("[WARN] SSL  Persistence Profile (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	d.Set("defaults_from", pp.DefaultsFrom)
	if err := d.Set("match_across_pools", pp.MatchAcrossPools); err != nil {
		return fmt.Errorf("[DEBUG] Error saving MatchAcrossPools to state for PersistenceProfile SSL  (%s): %s", d.Id(), err)
	}
	if err := d.Set("match_across_services", pp.MatchAcrossServices); err != nil {
		return fmt.Errorf("[DEBUG] Error saving MatchAcrossServices to state for PersistenceProfile SSL  (%s): %s", d.Id(), err)
	}
	if err := d.Set("match_across_virtuals", pp.MatchAcrossVirtuals); err != nil {
		return fmt.Errorf("[DEBUG] Error saving MatchAcrossVirtuals to state for PersistenceProfile SSL  (%s): %s", d.Id(), err)
	}
	d.Set("mirror", pp.Mirror)
	d.Set("timeout", pp.Timeout)
	d.Set("override_conn_limit", pp.OverrideConnectionLimit)

	return nil
}

func resourceBigipLtmPersistenceProfileSSLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

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
		return err
	}

	return resourceBigipLtmPersistenceProfileSSLRead(d, meta)
}

func resourceBigipLtmPersistenceProfileSSLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting SSL Persistence Profile " + name)
	err := client.DeleteSSLPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete SSL Persistence Profile  (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

func resourceBigipLtmPersistenceProfileSSLExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching SSL Persistence Profile " + name)

	pp, err := client.GetSSLPersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrive SSL Persistence Profile (%s) (%v) ", name, err)
		return false, err
	}

	if pp == nil {
		log.Printf("[WARN] persistance profile SSL  (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return pp != nil, nil
}
