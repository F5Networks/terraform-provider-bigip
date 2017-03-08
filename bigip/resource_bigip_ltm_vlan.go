package bigip

import (
	"log"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmVlan() *schema.Resource {
	log.Println("Resource schema");

	return &schema.Resource{
		Create: resourceBigipLtmVlanCreate,
		Read:   resourceBigipLtmVlanRead,
		Update: resourceBigipLtmVlanUpdate,
		Delete: resourceBigipLtmVlanDelete,
		Exists: resourceBigipLtmVlanExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmVlanImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the vlan",
				ValidateFunc: validateF5Name,
			},

			"tag": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Tagged number",
			},

			"interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							Description: "Interface name",
						},

						"tagged": &schema.Schema{
							Type:     schema.TypeBool,
							Required: true,
							Description: "Interface tagged",
						},
					},
				},
			},
		},
	}

}

func resourceBigipLtmVlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	tag := d.Get("tag").(int)

	log.Println("[INFO] Creating vlan ")

	err := client.CreateVlan(
		name,
		tag,
	)

	if err != nil {
        return err
    }
	
	ifaceCount := d.Get("interfaces.#").(int)
	for i := 0; i < ifaceCount; i++ {
		prefix := fmt.Sprintf("interfaces.%d", i)
		iface := d.Get(prefix + ".name").(string)
		tagged := d.Get(prefix + ".tagged").(bool)

		err = client.AddInterfaceToVlan(name, iface, tagged)
		if err != nil {
			return err
		}
	}
	
	d.SetId(name)

	return resourceBigipLtmVlanRead(d, meta)
}

func resourceBigipLtmVlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

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
			d.Set("tag", vlan.Tag)
		}
	}

	return nil
}

func resourceBigipLtmVlanExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

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

	return false, nil
}

func resourceBigipLtmVlanUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBigipLtmVlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Deleting vlan " + name)

	err := client.DeleteVlan(
		name,
	)

	if err != nil {
        return err
    }

	return nil
}

func resourceBigipLtmVlanImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
