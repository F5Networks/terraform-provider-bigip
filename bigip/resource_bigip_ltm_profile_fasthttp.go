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

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmProfileFasthttp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileFasthttpCreate,
		Update: resourceBigipLtmProfileFasthttpUpdate,
		Read:   resourceBigipLtmProfileFasthttpRead,
		Delete: resourceBigipLtmProfileFasthttpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
				//	Default:     300,
				Computed: true,
			},

			"connpoolidle_timeoutoverride": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "idle_timeout can be given value",
				Computed:    true,
			},

			"connpool_maxreuse": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "connpool_maxreuse timer",
				//	Default:     0,
				Computed: true,
			},

			"connpool_maxsize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "timer integer",
				Computed:    true,
			},

			"connpool_minsize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Pool min size",
				Computed:    true,
			},

			"connpool_replenish": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "enabled or disabled",
				//	Default:     "enabled",
				Computed: true,
			},

			"connpool_step": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer value",
				//	Default:     4,
				Computed: true,
			},
			"forcehttp_10response": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "disabled or enabled ",
				//	Default:     "",
				Computed: true,
			},

			"maxheader_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer value",
				//	Default:     32768,
				Computed: true,
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

	r := &bigip.Fasthttp{
		Name:                        name,
		DefaultsFrom:                defaultsFrom,
		IdleTimeout:                 idleTimeout,
		ConnpoolIdleTimeoutOverride: connpoolIdleTimeoutOverride,
		ConnpoolMaxReuse:            connpoolMaxReuse,
		ConnpoolMaxSize:             connpoolMaxSize,
		ConnpoolMinSize:             connpoolMinSize,
		ConnpoolReplenish:           connpoolReplenish,
		ConnpoolStep:                connpoolStep,
		ForceHttp_10Response:        forceHttp_10Response,
		MaxHeaderSize:               maxHeaderSize,
	}
	err := client.CreateFasthttp(r)

	if err != nil {
		log.Printf("[ERROR] Unable to Create Fasthttp   (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceBigipLtmProfileFasthttpRead(d, meta)
}

func resourceBigipLtmProfileFasthttpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

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
		log.Printf("[ERROR] Unable to Modify Fasthttp   (%s) (%v) ", name, err)
		return err
	}
	return resourceBigipLtmProfileFasthttpRead(d, meta)

}

func resourceBigipLtmProfileFasthttpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetFasthttp(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Fasthttp   (%s) (%v) ", name, err)
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Fasthttp profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	if _, ok := d.GetOk("defaults_from"); ok {
		if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
			return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("idle_timeout"); ok {
		if err := d.Set("idle_timeout", obj.IdleTimeout); err != nil {
			return fmt.Errorf("[DEBUG] Error saving IdleTimeout to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("connpoolidle_timeoutoverride"); ok {
		if err := d.Set("connpoolidle_timeoutoverride", obj.ConnpoolIdleTimeoutOverride); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ConnpoolIdleTimeoutOverride to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("connpool_maxreuse"); ok {

		if err := d.Set("connpool_maxreuse", obj.ConnpoolMaxReuse); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ConnpoolMaxReuse to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("connpool_maxsize"); ok {
		if err := d.Set("connpool_maxsize", obj.ConnpoolMaxSize); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ConnpoolMaxSize to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("connpool_minsize"); ok {
		if err := d.Set("connpool_minsize", obj.ConnpoolMinSize); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ConnpoolMinSize to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("connpool_replenish"); ok {
		if err := d.Set("connpool_replenish", obj.ConnpoolReplenish); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ConnpoolReplenish to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("connpool_step"); ok {
		if err := d.Set("connpool_step", obj.ConnpoolStep); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ConnpoolStep to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("forcehttp_10response"); ok {
		if err := d.Set("forcehttp_10response", obj.ForceHttp_10Response); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ForceHttp_10Response to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk("maxheader_size"); ok {
		if err := d.Set("maxheader_size", obj.MaxHeaderSize); err != nil {
			return fmt.Errorf("[DEBUG] Error saving MaxHeaderSize to state for Fasthttp profile  (%s): %s", d.Id(), err)
		}
	}
	return nil
}

func resourceBigipLtmProfileFasthttpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Fasthttp Profile " + name)

	err := client.DeleteFasthttp(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Fasthttp   (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}
