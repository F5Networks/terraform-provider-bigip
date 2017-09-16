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

			"defaultsFrom": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"concurrentStreamsPerConnection": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"connectionIdleTimeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},
			"headerTableSize": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use the parent Http2 profile",
			},

			"activationModes": &schema.Schema{
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
	defaultsFrom := d.Get("defaultsFrom").(string)
	concurrentStreamsPerConnection := d.Get("concurrentStreamsPerConnection").(int)
	connectionIdleTimeout := d.Get("connectionIdleTimeout").(int)
	headerTableSize := d.Get("headerTableSize").(int)
	activationModes := setToStringSlice(d.Get("activationModes").(*schema.Set))

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
		DefaultsFrom:                   d.Get("defaultsFrom").(string),
		ConcurrentStreamsPerConnection: d.Get("concurrentStreamsPerConnection").(int),
		ConnectionIdleTimeout:          d.Get("connectionIdleTimeout").(int),
		HeaderTableSize:                d.Get("headerTableSize").(int),
		ActivationModes:                setToStringSlice(d.Get("activationModes").(*schema.Set)),
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
