package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmDevicegroup() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmDevicegroupCreate,
		Update: resourceBigipLtmDevicegroupUpdate,
		Read:   resourceBigipLtmDevicegroupRead,
		Delete: resourceBigipLtmDevicegroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmDevicegroupImporter,
		},

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Devicegroup which needs to be Devicegroupensed",
			},

			"autoSync": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"fullLoadOnSync": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
		},
	}

}

func resourceBigipLtmDevicegroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	autoSync := d.Get("autoSync").(string)
	name := d.Get("name").(string)
	typo := d.Get("type").(string)
	fullLoadOnSync := d.Get("fullLoadOnSync").(string)

	log.Println("[INFO] Creating Devicegroup ")

	err := client.CreateDevicegroup(
		name,
		autoSync,
		typo,
		fullLoadOnSync,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipLtmDevicegroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Devicegroup " + name)

	r := &bigip.Devicegroup{
		Name:           name,
		AutoSync:       d.Get("autoSync").(string),
		Type:           d.Get("type").(string),
		FullLoadOnSync: d.Get("fullLoadOnSync").(string),
	}

	return client.ModifyDevicegroup(r)
}

func resourceBigipLtmDevicegroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Devicegroup " + name)

	members, err := client.Devicegroups()
	if err != nil {
		return err
	}

	d.Set("name", members.Name)
	d.Set("autoSync", members.AutoSync)
	d.Set("type", members.Type)
	d.Set("fullLoadOnSync", members.FullLoadOnSync)
	return nil
}

func resourceBigipLtmDevicegroupDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmDevicegroupImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
