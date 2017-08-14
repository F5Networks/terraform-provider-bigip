package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmSelfIP() *schema.Resource {
	log.Println("Resource schema")

	return &schema.Resource{
		Create: resourceBigipLtmSelfIPCreate,
		Read:   resourceBigipLtmSelfIPRead,
		Update: resourceBigipLtmSelfIPUpdate,
		Delete: resourceBigipLtmSelfIPDelete,
		Exists: resourceBigipLtmSelfIPExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmSelfIPImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the SelfIP",
				//ValidateFunc: validateF5Name,
			},

			"ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "SelfIP IP address",
			},

			"vlan": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the vlan",
				//ValidateFunc: validateF5Name,
			},
		},
	}
}

func resourceBigipLtmSelfIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	ip := d.Get("ip").(string)
	vlan := d.Get("vlan").(string)

	log.Println("[INFO] Creating SelfIP ")

	err := client.CreateSelfIP(name, ip, vlan)
	// err := client.CreateSelfIP(name+"-self", ip, vlan)

	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceBigipLtmSelfIPRead(d, meta)
	// return nil
}

func resourceBigipLtmSelfIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching SelfIP " + name)

	selfIPs, err := client.SelfIPs()
	if err != nil {
		return err
	}
	for _, selfip := range selfIPs.SelfIPs {
		log.Println(selfip.Name)
		if selfip.Name == name {
			return nil
		}
	}

	return nil
}

func resourceBigipLtmSelfIPExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching SelfIP " + name)

	selfIPs, err := client.SelfIPs()
	if err != nil {
		return false, err
	}
	for _, selfip := range selfIPs.SelfIPs {
		log.Println(selfip.Name)
		if selfip.Name == name {
			return true, nil
		}
	}

	return false, nil

}

func resourceBigipLtmSelfIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating SelfIP " + name)

	r := &bigip.SelfIP{
		Name:    name,
		Address: d.Get("ip").(string),
		Vlan:    d.Get("vlan").(string),
	}

	return client.ModifySelfIP(name, r)

}

func resourceBigipLtmSelfIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Deleting selfIP " + name)

	return client.DeleteSelfIP(name)
}

func resourceBigipLtmSelfIPImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
