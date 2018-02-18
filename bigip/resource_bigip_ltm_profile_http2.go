package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmProfileHttp2() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileHttp2Create,
		Update: resourceBigipLtmProfileHttp2Update,
		Read:   resourceBigipLtmProfileHttp2Read,
		Delete: resourceBigipLtmProfileHttp2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmProfileHttp2Importer,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Http2 Profile",
				//ValidateFunc: validateF5Name,
			},

			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"concurrent_streams_per_connection": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"connection_idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},
			"header_table_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"activation_modes": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},
		},
	}
}

func resourceBigipLtmProfileHttp2Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	concurrentStreamsPerConnection := d.Get("concurrent_streams_per_connection").(int)
	connectionIdleTimeout := d.Get("connection_idle_timeout").(int)
	headerTableSize := d.Get("header_table_size").(int)
	activationModes := setToStringSlice(d.Get("activation_modes").(*schema.Set))

	log.Println("[INFO] Creating Http2 profile")

	err := client.CreateHttp2(
		name,
		defaultsFrom,
		concurrentStreamsPerConnection,
		connectionIdleTimeout,
		headerTableSize,
		activationModes,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return resourceBigipLtmProfileHttp2Read(d, meta)
}

func resourceBigipLtmProfileHttp2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Http2{
		Name:                           name,
		DefaultsFrom:                   d.Get("defaults_from").(string),
		ConcurrentStreamsPerConnection: d.Get("concurrentr_streams_perr_connection").(int),
		ConnectionIdleTimeout:          d.Get("connection_idle_timeout").(int),
		HeaderTableSize:                d.Get("header_table_size").(int),
		ActivationModes:                setToStringSlice(d.Get("activation_modes").(*schema.Set)),
	}

	err := client.ModifyHttp2(name, r)
	if err != nil {
		return err
	}
	 return resourceBigipLtmProfileHttp2Read(d, meta)
}

func resourceBigipLtmProfileHttp2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetHttp2(name)
	if err != nil {
		d.SetId("")
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Http2 Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	d.Set("concurrent_streams_per_connection", obj.ConcurrentStreamsPerConnection)
	d.Set("connection_idle_timeout", obj.ConnectionIdleTimeout)
	d.Set("activation_modes", obj.ActivationModes)
	return nil
}

func resourceBigipLtmProfileHttp2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Http2 Profile " + name)

	err := client.DeleteHttp2(name)
	if err != nil {
		return err
	}
	if err == nil {
		log.Printf("[WARN] Http2 profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

func resourceBigipLtmProfileHttp2Importer(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
