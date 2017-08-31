package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmFasthttp() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmFasthttpCreate,
		Update: resourceBigipLtmFasthttpUpdate,
		Read:   resourceBigipLtmFasthttpRead,
		Delete: resourceBigipLtmFasthttpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmFasthttpImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Fasthttp Profile",
				//ValidateFunc: validateF5Name,
			},

			"defaultsFrom": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fasthttp profile",
			},

			"idleTimeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer value",
			},

			"connpoolIdleTimeoutOverride": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "idleTimeout can be given value",
			},

			"connpoolMaxReuse": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "connpoolMaxReuse timer",
			},

			"connpoolMaxSize": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "timer integer",
			},

			"connpoolMinSize": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Pool min size",
			},

			"connpoolReplenish": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "enabled or disabled",
			},

			"connpoolStep": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer value",
			},
			"forceHttp_10Response": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "disabled or enabled ",
			},

			"maxHeaderSize": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer value",
			},
		},
	}

}

func resourceBigipLtmFasthttpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	defaultsFrom := d.Get("defaultsFrom").(string)
	idleTimeout := d.Get("idleTimeout").(int)
	connpoolIdleTimeoutOverride := d.Get("connpoolIdleTimeoutOverride").(int)
	connpoolMaxReuse := d.Get("connpoolMaxReuse").(int)
	connpoolMaxSize := d.Get("connpoolMaxSize").(int)
	connpoolMinSize := d.Get("connpoolMinSize").(int)
	connpoolReplenish := d.Get("connpoolReplenish").(string)
	connpoolStep := d.Get("connpoolStep").(int)
	forceHttp_10Response := d.Get("forceHttp_10Response").(string)
	maxHeaderSize := d.Get("maxHeaderSize").(int)
	log.Println("[INFO] Creating Fasthttp profile")

	err := client.CreateFasthttp(
		name,
		defaultsFrom,
		idleTimeout,
		connpoolIdleTimeoutOverride,
		connpoolMaxReuse,
		connpoolMaxSize,
		connpoolMinSize,
		connpoolReplenish,
		connpoolStep,
		forceHttp_10Response,
		maxHeaderSize,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return resourceBigipLtmFasthttpRead(d, meta)
}

func resourceBigipLtmFasthttpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Fasthttp{
		Name:                        name,
		DefaultsFrom:                d.Get("defaultsFrom").(string),
		IdleTimeout:                 d.Get("idleTimeout").(int),
		ConnpoolIdleTimeoutOverride: d.Get("connpoolIdleTimeoutOverride").(int),
		ConnpoolMaxReuse:            d.Get("connpoolMaxReuse").(int),
		ConnpoolMaxSize:             d.Get("connpoolMaxSize").(int),
		ConnpoolMinSize:             d.Get("connpoolMinSize").(int),
		ConnpoolReplenish:           d.Get("connpoolReplenish").(string),
		ConnpoolStep:                d.Get("connpoolStep").(int),
		ForceHttp_10Response:        d.Get("forceHttp_10Response").(string),
		MaxHeaderSize:               d.Get("maxHeaderSize").(int),
	}

	return client.ModifyFasthttp(name, r)
}

func resourceBigipLtmFasthttpRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmFasthttpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Fasthttp Profile " + name)

	return client.DeleteFasthttp(name)
}

func resourceBigipLtmFasthttpImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
