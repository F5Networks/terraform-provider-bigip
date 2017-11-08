package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmDevicename() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmDevicenameCreate,
		Update: resourceBigipLtmDevicenameUpdate,
		Read:   resourceBigipLtmDevicenameRead,
		Delete: resourceBigipLtmDevicenameDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmDevicenameImporter,
		},

		Schema: map[string]*schema.Schema{

			"command": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of BIG-IP",
				//	ValidateFunc: validateF5Name,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Device which needs to be Devicenameensed",
			},

			"target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
		},
	}

}

func resourceBigipLtmDevicenameCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	command := d.Get("command").(string)
	name := d.Get("name").(string)
	target := d.Get("target").(string)
	log.Println("[INFO] Creating Devicename ")

	err := client.CreateDevicename(
		command,
		name,
		target,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipLtmDevicenameUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Devicename " + name)

	r := &bigip.Devicename{
		Name:    name,
		Command: d.Get("command").(string),
		Target:  d.Get("target").(string),
	}

	return client.ModifyDevicename(r)
}

func resourceBigipLtmDevicenameRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Devicename " + name)

	members, err := client.Devicenames()
	if err != nil {
		return err
	}

	d.Set("name", members.Name)
	d.Set("command", members.Command)
	d.Set("target", members.Target)

	return nil
}

func resourceBigipLtmDevicenameDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmDevicenameImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
