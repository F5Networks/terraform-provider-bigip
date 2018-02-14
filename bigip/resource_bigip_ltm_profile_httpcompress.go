package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmProfileHttpcompress() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileHttpcompressCreate,
		Update: resourceBigipLtmProfileHttpcompressUpdate,
		Read:   resourceBigipLtmProfileHttpcompressRead,
		Delete: resourceBigipLtmProfileHttpcompressDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmProfileHttpcompressImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Httpcompress Profile",
				//ValidateFunc: validateF5Name,
			},

			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Use the parent Httpcompress profile",
			},

			"uri_exclude": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},
			"uri_include": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},
		},
	}
}

func resourceBigipLtmProfileHttpcompressCreate(d *schema.ResourceData, meta interface{}) error {
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

func resourceBigipLtmProfileHttpcompressUpdate(d *schema.ResourceData, meta interface{}) error {
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

func resourceBigipLtmProfileHttpcompressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetHttpcompress(name)
	if err != nil {
	 d.SetId("")
	return err
	}
	if obj == nil {
 			log.Printf("[WARN] Httpcompress Profile (%s) not found, removing from state", d.Id())
 			d.SetId("")
 			return nil
 		}
	d.Set("name", name)
	d.Set("uri_include", obj.UriInclude)
	//d.Set("uri_exclude", obj.UriInclude)

	return nil
}

func resourceBigipLtmProfileHttpcompressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Httpcompress Profile " + name)

	err := client.DeleteHttpcompress(name)
	if err != nil {
		return err
	}
	if err == nil {
		log.Printf("[WARN] Httpcompress profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

func resourceBigipLtmProfileHttpcompressImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
