package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmHttpcompress() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmHttpcompressCreate,
		Update: resourceBigipLtmHttpcompressUpdate,
		Read:   resourceBigipLtmHttpcompressRead,
		Delete: resourceBigipLtmHttpcompressDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmHttpcompressImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Httpcompress Profile",
				//ValidateFunc: validateF5Name,
			},

			"defaults_from": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Httpcompress profile",
			},

			"uri_exclude": &schema.Schema{
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},
			"uri_include": &schema.Schema{
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},
		},
	}
}

func resourceBigipLtmHttpcompressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	defaultsFrom := d.Get("defaults_from").(string)
	uriExclude := setToStringSlice(d.Get("uri_exclude").(*schema.Set))
	uriInclude := setToStringSlice(d.Get("uri_include").(*schema.Set))

	log.Println("[INFO] Creating Httpcompress profile")

	err := client.CreateHttpcompress(
		name,
		defaultsFrom,
		uriExclude,
		uriInclude,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipLtmHttpcompressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	r := &bigip.Httpcompress{
		Name:         name,
		DefaultsFrom: d.Get("defaults_from").(string),
		UriExclude:   setToStringSlice(d.Get("uri_exclude").(*schema.Set)),
		UriInclude:   setToStringSlice(d.Get("uri_include").(*schema.Set)),
	}

	return client.ModifyHttpcompress(name, r)
}

func resourceBigipLtmHttpcompressRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBigipLtmHttpcompressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Httpcompress Profile " + name)

	return client.DeleteHttpcompress(name)
}

func resourceBigipLtmHttpcompressImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
