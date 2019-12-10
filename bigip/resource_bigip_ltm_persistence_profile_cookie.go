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

func resourceBigipLtmPersistenceProfileCookie() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPersistenceProfileCookieCreate,
		Read:   resourceBigipLtmPersistenceProfileCookieRead,
		Update: resourceBigipLtmPersistenceProfileCookieUpdate,
		Delete: resourceBigipLtmPersistenceProfileCookieDelete,
		Exists: resourceBigipLtmPersistenceProfileCookieExists,
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
				Description:  "To enable _ disable match across virtual servers with given persistence record",
				ValidateFunc: validateEnabledDisabled,
			},

			"mirror": {
				Type:         schema.TypeString,
				Default:      "disabled",
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

			// Specific to CookiePersistenceProfile
			"always_send": {
				Type:         schema.TypeString,
				Default:      "default",
				Optional:     true,
				Description:  "To enable _ disable always sending cookies",
				ValidateFunc: validateEnabledDisabled,
			},

			"cookie_encryption": {
				Type:         schema.TypeString,
				Default:      "disabled",
				Optional:     true,
				Description:  "To required, preferred, or disabled policy for cookie encryption",
				ValidateFunc: validateReqPrefDisabled,
			},

			"cookie_encryption_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Passphrase for encrypted cookies",
			},

			"cookie_name": {
				Type:        schema.TypeString,
				Default:     "disabled",
				Optional:    true,
				Description: "Name of the cookie to track persistence",
			},

			"expiration": {
				Type:        schema.TypeString,
				Default:     "0",
				Optional:    true,
				Description: "Expiration TTL for cookie specified in D:H:M:S or in seconds",
			},

			"hash_length": {
				Type:        schema.TypeInt,
				Default:     0,
				Optional:    true,
				Description: "Length of hash to apply to cookie",
			},

			"hash_offset": {
				Type:        schema.TypeInt,
				Default:     0,
				Optional:    true,
				Description: "Number of characters to skip in the cookie for the hash",
			},

			"httponly": {
				Type:         schema.TypeString,
				Default:      "disabled",
				Optional:     true,
				Description:  "To enable _ disable sending only over http",
				ValidateFunc: validateEnabledDisabled,
			},
		},
	}
}

func resourceBigipLtmPersistenceProfileCookieCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	err := client.CreateCookiePersistenceProfile(
		name,
		parent,
	)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Cookie Persistence Profile %s %v :", name, err)
		return err
	}

	d.SetId(name)

	err = resourceBigipLtmPersistenceProfileCookieUpdate(d, meta)
	if err != nil {
		client.DeleteCookiePersistenceProfile(name)
		return err
	}

	return resourceBigipLtmPersistenceProfileCookieRead(d, meta)

}

func resourceBigipLtmPersistenceProfileCookieRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Cookie Persistence Profile " + name)

	pp, err := client.GetCookiePersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Cookie Persistence Profile %s  %v : ", name, err)
		return err
	}
	if pp == nil {
		log.Printf("[WARN] Cookie Persistence Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	if err := d.Set("app_service", pp.AppService); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AppService to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
	}

	if err := d.Set("defaults_from", pp.DefaultsFrom); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
	}
	d.Set("match_across_pools", pp.MatchAcrossPools)
	d.Set("match_across_services", pp.MatchAcrossServices)
	d.Set("match_across_virtuals", pp.MatchAcrossVirtuals)
	if err := d.Set("mirror", pp.Mirror); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mirror to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
	}
	d.Set("timeout", pp.Timeout)
	d.Set("override_conn_limit", pp.OverrideConnectionLimit)

	// Specific to CookiePersistenceProfile
	d.Set("always_send", pp.AlwaysSend)
	if err := d.Set("cookie_encryption", pp.CookieEncryption); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CookieEncryption to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
	}
	if err := d.Set("cookie_encryption_passphrase", pp.CookieEncryptionPassphrase); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CookieEncryptionPassphrase to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
	}
	d.Set("cookie_name", pp.CookieName)
	d.Set("expiration", pp.Expiration)
	d.Set("hash_length", pp.HashLength)
	d.Set("hash_offset", pp.HashOffset)
	d.Set("httponly", pp.HTTPOnly)

	return nil
}

func resourceBigipLtmPersistenceProfileCookieUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	pp := &bigip.CookiePersistenceProfile{
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
		// Specific to CookiePersistenceProfile
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
		return err
	}

	return resourceBigipLtmPersistenceProfileCookieRead(d, meta)
}

func resourceBigipLtmPersistenceProfileCookieDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Cookie Persistence Profile " + name)
	err := client.DeleteCookiePersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Cookie Persistence Profile %s  %v : ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

func resourceBigipLtmPersistenceProfileCookieExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Cookie Persistence Profile " + name)

	pp, err := client.GetCookiePersistenceProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Cookie Persistence Profile %s  %v : ", name, err)
		return false, err
	}
	if pp == nil {
		log.Printf("[WARN] persistance profile cookie (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, nil
	}

	return pp != nil, nil
}
