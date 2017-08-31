package bigip

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmDatagroup() *schema.Resource {
	log.Println("Resource schema")

	return &schema.Resource{
		Create: resourceBigipLtmDatagroupCreate,
		Read:   resourceBigipLtmDatagroupRead,
		Update: resourceBigipLtmDatagroupUpdate,
		Delete: resourceBigipLtmDatagroupDelete,
		//Exists: resourceBigipLtmVlanExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmDatagroupImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the vlan",
				//			ValidateFunc: validateF5Name,
			},

			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of datagroup",
			},

			"records": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: " name field in datagroup",
						},

						"data": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "value",
						},
					},
				},
			},
		},
	}

}

func resourceBigipLtmDatagroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	typo := d.Get("type").(string)

	log.Println("[INFO] Creating datagroup ")
	fmt.Println("[INFO] Creating datagroup ")
	records := listToSlice2(d)
	err := client.CreateDatagroup(
		typo,
		name,
		records,
	)

	if err != nil {
		return err
	}

	/*recordCount := d.Get("records.#").(int)
	for i := 0; i < recordCount; i++ {
		prefix := fmt.Sprintf("records.%d", i)
		rname := d.Get(prefix + "name").(string)
		data := d.Get(prefix + "data").(string)

		err = client.AddRecords(name, rname, data)
		if err != nil {
			return err
		}
	}
	*/
	// New cod

	d.SetId(name)

	return nil

	//	return resourceBigipLtmVlanRead(d, meta)
}

func resourceBigipLtmDatagroupRead(d *schema.ResourceData, meta interface{}) error {
	/*	client := meta.(*bigip.BigIP)

		name := d.Id()

		log.Println("[INFO] Reading vlan " + name)

		vlans, err := client.Vlans()
		if err != nil {
			return err
		}

		for _, vlan := range vlans.Vlans {
			log.Println(vlan.Name)
			if vlan.Name == name {
				d.Set("name", vlan.Name)
			}
		}
	*/
	return nil
}

func resourceBigipLtmDatagroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	/* client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Vlan " + name)

	vlans, err := client.Vlans()
	if err != nil {
		return false, err
	}
	for _, vlan := range vlans.Vlans {
		log.Println(vlan.Name)
		if vlan.Name == name {
			return true, nil
		}
	}
	*/
	return false, nil
}

func resourceBigipLtmDatagroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Vlan " + name)

	r := &bigip.Vlan{
		Name: name,
		Tag:  d.Get("tag").(int),
	}

	return client.ModifyVlan(name, r)

}

func resourceBigipLtmDatagroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Deleting vlan " + name)

	return client.DeleteVlan(name)
}

func resourceBigipLtmDatagroupImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func listToSlice2(d *schema.ResourceData) []bigip.Records {
	addrecordCount := d.Get("records.#").(int)
	var r = make([]bigip.Records, addrecordCount, addrecordCount)

	for i := 0; i < addrecordCount; i++ {
		prefix := fmt.Sprintf("records.%d", i)
		r[i].Name = d.Get(prefix + ".name").(string)
		r[i].Data = d.Get(prefix + ".data").(string)
	}

	return r
}
