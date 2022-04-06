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

func resourceBigipLtmPersistenceProfileCookieCreate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}

	d.SetId(name)

	err = resourceBigipLtmPersistenceProfileCookieUpdate(d, meta)
	if err != nil {
		if errdel := client.DeleteCookiePersistenceProfile(name); errdel != nil {
			return errdel
		}
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
	_ = d.Set("name", name)
	if _, ok := d.GetOk("partition"); ok {
		if err := d.Set("app_service", pp.AppService); err != nil {
			return fmt.Errorf("[DEBUG] Error saving AppService to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("partition"); ok {
		if err := d.Set("defaults_from", pp.DefaultsFrom); err != nil {
			return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("match_across_pools", pp.MatchAcrossPools)
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("match_across_services", pp.MatchAcrossServices)
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("match_across_virtuals", pp.MatchAcrossVirtuals)
	}
	if _, ok := d.GetOk("partition"); ok {
		if err := d.Set("mirror", pp.Mirror); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Mirror to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("method", pp.Method)
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("timeout", pp.Timeout)
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("override_conn_limit", pp.OverrideConnectionLimit)
	}
	// Specific to CookiePersistenceProfile
	if _, ok := d.GetOk("partition"); ok {

		_ = d.Set("always_send", pp.AlwaysSend)
	}
	if _, ok := d.GetOk("partition"); ok {
		if err := d.Set("cookie_encryption", pp.CookieEncryption); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CookieEncryption to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("partition"); ok {
		if err := d.Set("cookie_encryption_passphrase", pp.CookieEncryptionPassphrase); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CookieEncryptionPassphrase to state for PersistenceProfileCookie (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("cookie_name", pp.CookieName)
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("expiration", pp.Expiration)
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("hash_length", pp.HashLength)
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("hash_offset", pp.HashOffset)
	}
	if _, ok := d.GetOk("partition"); ok {
		_ = d.Set("httponly", pp.HTTPOnly)
	}

	return nil
}

func resourceBigipLtmPersistenceProfileCookieUpdate(d *schema.ResourceData, meta interface{}) error {
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
		log.Printf("[WARN] persistence profile cookie (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, nil
	}

	return pp != nil, nil
}
