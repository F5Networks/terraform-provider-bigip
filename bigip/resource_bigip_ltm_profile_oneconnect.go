package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmOneconnect() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmOneconnectCreate,
		Update: resourceBigipLtmOneconnectUpdate,
		Read:   resourceBigipLtmOneconnectRead,
		Delete: resourceBigipLtmOneconnectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmOneconnectImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Oneconnect Profile",
				//ValidateFunc: validateF5Name,
			},
			"partition": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaultsFrom": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent oneconnect profile",
			},

			"idleTimeoutOverride": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "idleTimeoutOverride can be enabled or disabled",
			},

			"sharePools": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "sharePools can be enabled or disabled",
			},
			"sourceMask": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "sourceMask can be 255.255.255.255",
			},

			"maxAge": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "maxAge has integer value typical 3600 sec",
			},
			"maxReuse": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "maxReuse has integer value typical 1000 sec",
			},
			"maxSize": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "maxSize has integer value typical 1000 sec",
			},
		},
	}

}

func resourceBigipLtmOneconnectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	defaultsFrom := d.Get("defaultsFrom").(string)
	sharePools := d.Get("sharePools").(string)
	maxAge := d.Get("maxAge").(int)
	maxReuse := d.Get("maxReuse").(int)
	maxSize := d.Get("maxSize").(int)
	sourceMask := d.Get("sourceMask").(string)
	idleTimeoutOverride := d.Get("idleTimeoutOverride").(string)

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
		IdleTimeoutOverride: d.Get("idleTimeoutOverride").(string),
		Partition:           d.Get("partition").(string),
		DefaultsFrom:        d.Get("defaultsFrom").(string),
		SharePools:          d.Get("sharePools").(string),
		SourceMask:          d.Get("sourceMask").(string),
		MaxAge:              d.Get("maxAge").(int),
		MaxSize:             d.Get("maxSize").(int),
		MaxReuse:            d.Get("maxReuse").(int),
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
