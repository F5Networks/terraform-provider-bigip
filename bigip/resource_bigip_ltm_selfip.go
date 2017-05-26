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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the SelfIP",
				ValidateFunc: validateF5Name,
			},

			"ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "SelfIP IP address",
			},

			"vlan": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the vlan",
				ValidateFunc: validateF5Name,
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

	err := client.CreateSelfIP(name+"-self", ip, vlan)

	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceBigipLtmSelfIPRead(d, meta)
}

func resourceBigipLtmSelfIPRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmSelfIPExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	return false, nil
}

func resourceBigipLtmSelfIPUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBigipLtmSelfIPDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmSelfIPImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
