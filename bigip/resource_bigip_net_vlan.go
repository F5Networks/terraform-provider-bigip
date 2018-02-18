package bigip

import (
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipNetVlan() *schema.Resource {

	return &schema.Resource{
		Create: resourceBigipNetVlanCreate,
		Read:   resourceBigipNetVlanRead,
		Update: resourceBigipNetVlanUpdate,
		Delete: resourceBigipNetVlanDelete,
		//Exists: resourceBigipNetVlanExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipNetVlanImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the vlan",
				//			ValidateFunc: validateF5Name,
			},

			"tag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Tagged number",
			},

			"interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlanport": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Vlan name",
						},

						"tagged": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Interface tagged",
						},
					},
				},
			},
		},
	}

}

func resourceBigipNetVlanCreate(d *schema.ResourceData, meta interface{}) error {
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
		iface := d.Get(prefix + ".vlanport").(string)
		tagged := d.Get(prefix + ".tagged").(bool)

		err = client.AddInterfaceToVlan(name, iface, tagged)
		if err != nil {
			return err
		}
	}

	d.SetId(name)

	return resourceBigipNetVlanRead(d, meta)

	//	return resourceBigipNetVlanRead(d, meta)
}

func resourceBigipNetVlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading vlan " + name)

	vlans, err := client.Vlans()
	if err != nil {
		return err
	}
	if vlans == nil {
		log.Printf("[WARN] Node (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	for _, vlan := range vlans.Vlans {
		log.Println(vlan.Name)
		if vlan.Name == name {
			d.Set("name", vlan.Name)
		}
	}

	return nil
}

func resourceBigipNetVlanExists(d *schema.ResourceData, meta interface{}) (bool, error) {
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

func resourceBigipNetVlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Vlan " + name)

	r := &bigip.Vlan{
		Name: name,
		Tag:  d.Get("tag").(int),
	}

	err := client.ModifyVlan(name, r)
	if err != nil {
		return err
	}
	return resourceBigipNetVlanRead(d, meta)

}

func resourceBigipNetVlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Deleting vlan " + name)

	err := client.DeleteVlan(name)
	if err == nil {
		log.Printf("[WARN] Vlan (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

func resourceBigipNetVlanImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
