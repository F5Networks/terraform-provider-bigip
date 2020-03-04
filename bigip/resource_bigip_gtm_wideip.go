package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBigipGtmWideip() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmWideipCreate,
		Read:   resourceBigipGtmWideipRead,
		Update: resourceBigipGtmWideipUpdate,
		Delete: resourceBigipGtmWideipDelete,
		//		Exists: resourceBigipGtmWideipExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Wideip",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"partition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"full_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_service": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"generation": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"failure_rcode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"failure_rcode_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"failure_rcode_response": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_resort_pool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"minimal_response": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"persistence": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pool_lb_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"persist_cidr_ipv4": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"persist_cidr_ipv6": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ttl_persistence": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}
func resourceBigipGtmWideipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	d.SetId(name)
	gtmtype := d.Get("type").(string)
	err := client.AddGTMWideIP(name, gtmtype)
	if err != nil {
		return fmt.Errorf("Error creating wideip (%s): %s", name, err)
	}
	return resourceBigipGtmWideipRead(d, meta)
}

func resourceBigipGtmWideipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	recordtype := d.Get("type").(string)

	log.Printf("[INFO] Fetching Wideip " + name)
	wideip, err := client.GetGTMWideIP(name, recordtype)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve wideip %s  %v :", name, err)
		return err
	}
	if wideip == nil {
		log.Printf("[WARN] wideip (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	d.Set("fullPath", wideip.FullPath)
	d.Set("generation", wideip.Generation)
	d.Set("enabled", wideip.Enabled)
	d.Set("failureRcode", wideip.FailureRcode)
	d.Set("failureRcodeResponse", wideip.FailureRcodeResponse)
	d.Set("failureRcodeTtl", wideip.FailureRcodeTTL)
	d.Set("lastResortPool", wideip.LastResortPool)
	d.Set("minimalResponse", wideip.MinimalResponse)
	d.Set("persistCidrIpv4", wideip.PersistCidrIpv4)
	d.Set("persistCidrIpv6", wideip.PersistCidrIpv6)
	d.Set("persistence", wideip.Persistence)
	d.Set("poolLbMode", wideip.PoolLbMode)
	d.Set("ttlPersistence", wideip.TTLPersistence)
	return nil
}

func resourceBigipGtmWideipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	fullpath := d.Get("full_path").(string)
	gtmtype := d.Get("type").(string)
	gtmwideip := &bigip.GTMWideIP{
		FullPath:             d.Get("full_path").(string),
		Generation:           d.Get("generation").(int),
		AppService:           d.Get("app_service").(string),
		Description:          d.Get("description").(string),
		Disabled:             d.Get("disabled").(bool),
		Enabled:              d.Get("enabled").(bool),
		FailureRcode:         d.Get("failure_rcode").(string),
		FailureRcodeResponse: d.Get("failure_rcode_response").(string),
		FailureRcodeTTL:      d.Get("failure_rcode_ttl").(int),
		LastResortPool:       d.Get("last_resort_pool").(string),
		//LoadBalancingDecisionLogVerbosity: d.Get("loadBalancingDecisionLogVerbosity,omitempty").
		MinimalResponse: d.Get("minimal_response").(string),
		PersistCidrIpv4: d.Get("persist_cidr_ipv4").(int),
		PersistCidrIpv6: d.Get("persist_cidr_ipv6").(int),
		Persistence:     d.Get("persistence").(string),
		PoolLbMode:      d.Get("pool_lb_mode").(string),
		TTLPersistence:  d.Get("ttl_persistence").(int),
	}
	err := client.ModifyGTMWideIP(fullpath, gtmwideip, gtmtype)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify WideIp   (%s) (%v) ", name, err)
		return err
	}

	return resourceBigipGtmWideipRead(d, meta)
}
func resourceBigipGtmWideipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	gtmtype := d.Get("type").(string)
	name := d.Id()
	log.Println("[INFO] Deleting wideip " + name)
	err := client.DeleteGTMWideIP(name, gtmtype)

	if err != nil {
		log.Printf("[ERROR] Unable to Delete wideip %s  %v : ", name, err)
		return err
	}
	d.SetId("")
	return nil
}
