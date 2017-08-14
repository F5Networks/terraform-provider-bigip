package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmDatagroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmDatagroupCreate,
		Read:   resourceBigipLtmDatagroupRead,
		Update: resourceBigipLtmDatagroupUpdate,
		Delete: resourceBigipLtmDatagroupDelete,
		//Exists: resourceBigipLtmDatagroupExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmDatagroupImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the datagroup",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},

			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type string or address",
			},
			"records": &schema.Schema{
				Type:        schema.TypeSet,
				Set:       schema.HashString,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Give the record key and values",
			},
		},
	}
}

func resourceBigipLtmDatagroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	typos := d.Get("type").(string)
	//partition := d.Get("partition").(string)

	p := dataToRecords(name, d)

	log.Println("[INFO] Creating Data Group" + name)
  d.SetId(name)
	err := client.CreateDatagroup(&p)

	if err != nil {
		return err
	}
	d.SetId(name)
	return resourceBigipLtmDatagroupRead(d, meta)
}

func dataToRecords(name string, d *schema.ResourceData) bigip.Datagroup {
	var p bigip.Datagroup
	p.Name = name
	p.Typos = d.Get("type").(string)
	p.Records = setToStringSlice(d.Get("records").(*schema.Set))

func resourceBigipLtmDatagroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Datagroup " + description)

	r := &bigip.Datagroup{
		Name:      name,
		Type:      d.Get("type").(string),
		Records:   d.Get("records").(string),
	}

	return client.ModifyDatagroup(name, r)
}

func resourceBigipLtmDatagroupRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmDatagroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Datagroup " + name)

	return client.DeleteDatagroup(name)
}

func resourceBigipLtmDatagroupImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
