package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmHttp2() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmHttp2Create,
		Update: resourceBigipLtmHttp2Update,
		Read:   resourceBigipLtmHttp2Read,
		Delete: resourceBigipLtmHttp2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmHttp2Importer,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Http2 Profile",
				//ValidateFunc: validateF5Name,
			},

			"defaults_from": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"concurrent_streams_per_connection": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"connection_idle_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},
			"header_table_size": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"activation_modes": &schema.Schema{
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},
		},
	}
}

func resourceBigipLtmHttp2Create(d *schema.ResourceData, meta interface{}) error {
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
	return nil
}

func resourceBigipLtmHttp2Update(d *schema.ResourceData, meta interface{}) error {
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

	return client.ModifyHttp2(name, r)
}

func resourceBigipLtmHttp2Read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBigipLtmHttp2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Http2 Profile " + name)

	return client.DeleteHttp2(name)
}

func resourceBigipLtmHttp2Importer(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
