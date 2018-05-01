package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceBigipLtmProfileFasthttp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileFasthttpCreate,
		Update: resourceBigipLtmProfileFasthttpUpdate,
		Read:   resourceBigipLtmProfileFasthttpRead,
		Delete: resourceBigipLtmProfileFasthttpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmProfileFasthttpImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Fasthttp Profile",
			},

			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Fasthttp profile",
			},

			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer value",
				Default:     300,
			},

			"connpoolidle_timeoutoverride": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "idle_timeout can be given value",
			},

			"connpool_maxreuse": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "connpool_maxreuse timer",
				Default:     0,
			},

			"connpool_maxsize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "timer integer",
			},

			"connpool_minsize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Pool min size",
			},

			"connpool_replenish": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "enabled or disabled",
				Default:     "enabled",
			},

			"connpool_step": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer value",
				Default:     4,
			},
			"forcehttp_10response": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "disabled or enabled ",
			},

			"maxheader_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer value",
				Default:     32768,
			},
		},
	}

}

func resourceBigipLtmProfileFasthttpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	idleTimeout := d.Get("idle_timeout").(int)
	connpoolIdleTimeoutOverride := d.Get("connpoolidle_timeoutoverride").(int)
	connpoolMaxReuse := d.Get("connpool_maxreuse").(int)
	connpoolMaxSize := d.Get("connpool_maxsize").(int)
	connpoolMinSize := d.Get("connpool_minsize").(int)
	connpoolReplenish := d.Get("connpool_replenish").(string)
	connpoolStep := d.Get("connpool_step").(int)
	forceHttp_10Response := d.Get("forcehttp_10response").(string)
	maxHeaderSize := d.Get("maxheader_size").(int)
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
	return resourceBigipLtmProfileFasthttpRead(d, meta)
}

func resourceBigipLtmProfileFasthttpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Fasthttp{
		Name:                        name,
		DefaultsFrom:                d.Get("defaults_from").(string),
		IdleTimeout:                 d.Get("idle_timeout").(int),
		ConnpoolIdleTimeoutOverride: d.Get("connpoolidle_timeoutoverride").(int),
		ConnpoolMaxReuse:            d.Get("connpool_maxreuse").(int),
		ConnpoolMaxSize:             d.Get("connpool_maxsize").(int),
		ConnpoolMinSize:             d.Get("connpool_minsize").(int),
		ConnpoolReplenish:           d.Get("connpool_replenish").(string),
		ConnpoolStep:                d.Get("connpool_step").(int),
		ForceHttp_10Response:        d.Get("forcehttp_10response").(string),
		MaxHeaderSize:               d.Get("maxheader_size").(int),
	}

	err := client.ModifyFasthttp(name, r)
	if err != nil {
		return err
	}
	return resourceBigipLtmProfileFasthttpRead(d, meta)

}

func resourceBigipLtmProfileFasthttpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetFasthttp(name)
	if err != nil {
		d.SetId("")
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Fasthttp profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for Fasthttp profile  (%s): %s", d.Id(), err)
	}
	d.Set("idle_timeout", obj.IdleTimeout)
	d.Set("connpoolidle_timeoutoverride", obj.ConnpoolIdleTimeoutOverride)
	d.Set("connpool_maxreuse", obj.ConnpoolMaxReuse)
	d.Set("connpool_maxsize", obj.ConnpoolMaxSize)
	d.Set("connpool_minsize", obj.ConnpoolMinSize)
	d.Set("connpool_replenish", obj.ConnpoolReplenish)
	d.Set("connpool_step", obj.ConnpoolStep)
	d.Set("forceHttp_10Response", obj.ForceHttp_10Response)
	d.Set("maxheader_size", obj.MaxHeaderSize)
	return nil
}

func resourceBigipLtmProfileFasthttpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Fasthttp Profile " + name)

	err := client.DeleteFasthttp(name)
	if err != nil {
		return err
	}
	if err == nil {
		log.Printf("[WARN] Fasthttp Profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

func resourceBigipLtmProfileFasthttpImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
