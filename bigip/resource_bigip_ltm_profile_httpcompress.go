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

			"defaultsFrom": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Httpcompress profile",
			},

			"uriExclude": &schema.Schema{
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},
			"uriInclude": &schema.Schema{
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
	defaultsFrom := d.Get("defaultsFrom").(string)
	uriExclude := setToStringSlice(d.Get("uriExclude").(*schema.Set))
	uriInclude := setToStringSlice(d.Get("uriInclude").(*schema.Set))

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
		DefaultsFrom: d.Get("defaultsFrom").(string),
		UriExclude:   setToStringSlice(d.Get("uriExclude").(*schema.Set)),
		UriInclude:   setToStringSlice(d.Get("uriInclude").(*schema.Set)),
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
