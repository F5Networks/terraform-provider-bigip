package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipSysBigiplicense() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysBigiplicenseCreate,
		Update: resourceBigipSysBigiplicenseUpdate,
		Read:   resourceBigipSysBigiplicenseRead,
		Delete: resourceBigipSysBigiplicenseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"command": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tmsh command to execute tmsh commands like install",
			},
			"registration_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique Key F5 provides for Licensing BIG-IP",
			},
		},
	}

}

func resourceBigipSysBigiplicenseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	command := d.Get("command").(string)
	registration_key := d.Get("registration_key").(string)
	log.Println("[INFO] Creating BigipLicense ")

	err := client.CreateBigiplicense(
		command,
		registration_key,
	)

	if err != nil {
		return err
	}
	d.SetId(registration_key)
	return nil
}

func resourceBigipSysBigiplicenseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	registration_key := d.Id()

	log.Println("[INFO] Updating Bigiplicense " + registration_key)

	r := &bigip.Bigiplicense{
		Registration_key: registration_key,
		Command:          d.Get("command").(string),
	}

	return client.ModifyBigiplicense(r)
}

func resourceBigipSysBigiplicenseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Bigiplicense " + name)

	members, err := client.Bigiplicenses()
	if err != nil {
		return err
	}

	d.Set("registration_key", members.Registration_key)
	d.Set("command", members.Command)

	return nil
}

func resourceBigipSysBigiplicenseDelete(d *schema.ResourceData, meta interface{}) error {
	//API does not Exists
	return nil
}

func resourceBigipSysBigiplicenseImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
