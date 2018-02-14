package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipSourceAddrPersistenceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSourceAddrPersistenceProfileCreate,
		Read:   resourceBigipSourceAddrPersistenceProfileRead,
		Update: resourceBigipSourceAddrPersistenceProfileUpdate,
		Delete: resourceBigipSourceAddrPersistenceProfileDelete,
		Exists: resourceBigipSourceAddrPersistenceProfileExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipSourceAddrPersistenceProfileImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the persistence profile",
				ValidateFunc: validateF5Name,
			},

			"defaults_from": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Inherit defaults from parent profile",
			},

			"match_across_pools": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable match across pools with given persistence record",
			},

			"match_across_services": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable match across services with given persistence record",
			},

			"match_across_virtuals": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable match across services with given persistence record",
			},

			"mirror": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable _ disable ??",
			},

			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout for persistence of the session",
			},

			"override_conn_limit": {
				Type:        schema.TypeString,
				Default:     false,
				Optional:    true,
				Description: "To enable _ disable that pool member connection limits are overridden for persisted clients. Per-virtual connection limits remain hard limits and are not overridden.",
			},

			// Specific to SourceAddrPersistenceProfile
			"hash_algorithm": {
				Type:        schema.TypeString,
				Default:     "default",
				Optional:    true,
				Description: "Specify the hash algorithm",
			},

			"map_proxies": {
				Type:        schema.TypeString,
				Default:     true,
				Optional:    true,
				Description: "To enable _ disable directs all to the same single pool member",
			},

			"mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identify a range of source IP addresses to manage together as a single source address affinity persistent connection when connecting to the pool. Must be a valid IPv4 or IPv6 mask.",
			},
		},
	}
}

func resourceBigipSourceAddrPersistenceProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)

	err := client.CreateSourceAddrPersistenceProfile(
		name,
		parent,
	)

	d.SetId(name)

	err = resourceBigipSourceAddrPersistenceProfileUpdate(d, meta)
	if err != nil {
		client.DeleteSourceAddrPersistenceProfile(name)
		return err
	}

	return resourceBigipSourceAddrPersistenceProfileRead(d, meta)

}

func resourceBigipSourceAddrPersistenceProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Source Address Persistence Profile " + name)

	pp, err := client.GetSourceAddrPersistenceProfile(name)
	if err != nil {
		return err
	}

	d.Set("name", name)
	d.Set("defaults_from", pp.DefaultsFrom)
	d.Set("match_across_pools", pp.MatchAcrossPools)
	d.Set("match_across_services", pp.MatchAcrossServices)
	d.Set("match_across_virtuals", pp.MatchAcrossVirtuals)
	d.Set("mirror", pp.Mirror)
	d.Set("timeout", pp.Timeout)
	d.Set("override_conn_limit", pp.OverrideConnectionLimit)

	// Specific to SourceAddrPersistenceProfile
	d.Set("hash_algorithm", pp.HashAlgorithm)
	d.Set("map_proxies", pp.MapProxies)
	d.Set("mask", pp.Mask)

	return nil
}

func resourceBigipSourceAddrPersistenceProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	pp := &bigip.SourceAddrPersistenceProfile{
		PersistenceProfile: bigip.PersistenceProfile{
			DefaultsFrom:            d.Get("defaults_from").(string),
			MatchAcrossPools:        d.Get("match_across_pools").(string),
			MatchAcrossServices:     d.Get("match_across_services").(string),
			MatchAcrossVirtuals:     d.Get("match_across_virtuals").(string),
			Mirror:                  d.Get("mirror").(string),
			OverrideConnectionLimit: d.Get("override_conn_limit").(string),
			Timeout:                 d.Get("timeout").(string),
		},

		// Specific to SourceAddrPersistenceProfile
		HashAlgorithm: d.Get("hash_algorithm").(string),
		MapProxies:    d.Get("map_proxies").(string),
		Mask:          d.Get("mask").(string),
	}

	err := client.ModifySourceAddrPersistenceProfile(name, pp)
	if err != nil {
		return err
	}

	return nil
}

func resourceBigipSourceAddrPersistenceProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Source Address Persistence Profile " + name)

	return client.DeleteSourceAddrPersistenceProfile(name)
}

func resourceBigipSourceAddrPersistenceProfileExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Source Address Persistence Profile " + name)

	pp, err := client.GetSourceAddrPersistenceProfile(name)
	if err != nil {
		return false, err
	}

	if pp == nil {
		d.SetId("")
	}

	return pp != nil, nil
}

func resourceBigipSourceAddrPersistenceProfileImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
