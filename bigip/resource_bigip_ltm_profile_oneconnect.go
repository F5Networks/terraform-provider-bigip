package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmOneconnect() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmOneconnectCreate,
		Update: resourceBigipLtmOneconnectUpdate,
		Read:   resourceBigipLtmOneconnectRead,
		Delete: resourceBigipLtmOneconnectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmOneconnectImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Oneconnect Profile",
				//ValidateFunc: validateF5Name,
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent oneconnect profile",
			},

			"idle_timeout_override": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "idleTimeoutOverride can be enabled or disabled",
			},

			"share_pools": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "sharePools can be enabled or disabled",
			},
			"source_mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source_mask can be 255.255.255.255",
			},

			"max_age": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "max_age has integer value typical 3600 sec",
			},
			"max_reuse": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "max_reuse has integer value typical 1000 sec",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "max_size has integer value typical 1000 sec",
			},
		},
	}

}

func resourceBigipLtmOneconnectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	sharePools := d.Get("share_pools").(string)
	maxAge := d.Get("max_age").(int)
	maxReuse := d.Get("max_reuse").(int)
	maxSize := d.Get("max_size").(int)
	sourceMask := d.Get("source_mask").(string)
	idleTimeoutOverride := d.Get("idle_timeout_override").(string)

	log.Println("[INFO] Creating OneConnect profile")

	err := client.CreateOneconnect(
		name,
		idleTimeoutOverride,
		partition,
		defaultsFrom,
		sharePools,
		sourceMask,
		maxAge,
		maxReuse,
		maxSize,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipLtmOneconnectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Oneconnect{
		Name:                name,
		IdleTimeoutOverride: d.Get("idle_timeout_override").(string),
		Partition:           d.Get("partition").(string),
		DefaultsFrom:        d.Get("defaults_from").(string),
		SharePools:          d.Get("share_pools").(string),
		SourceMask:          d.Get("source_mask").(string),
		MaxAge:              d.Get("max_age").(int),
		MaxSize:             d.Get("max_size").(int),
		MaxReuse:            d.Get("max_reuse").(int),
	}

	return client.ModifyOneconnect(name, r)
}

func resourceBigipLtmOneconnectRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmOneconnectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting OneConnect Profile " + name)

	return client.DeleteOneconnect(name)
}

func resourceBigipLtmOneconnectImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
