package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmLic() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmLicCreate,
		Update: resourceBigipLtmLicUpdate,
		Read:   resourceBigipLtmLicRead,
		Delete: resourceBigipLtmLicDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmLicImporter,
		},

		Schema: map[string]*schema.Schema{

			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of BIG-IP",
				//	ValidateFunc: validateF5Name,
			},
			"device_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Device which needs to be licensed",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
		},
	}

}

func resourceBigipLtmLicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	deviceAddress := d.Get("device_address").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	log.Println("[INFO] Creating Lic ")

	err := client.CreateLIC(
		deviceAddress,
		username,
		password,
	)

	if err != nil {
		return err
	}
	d.SetId(deviceAddress)
	return nil
}

func resourceBigipLtmLicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	deviceAddress := d.Id()

	log.Println("[INFO] Updating LIC " + deviceAddress)

	r := &bigip.LIC{
		DeviceAddress: deviceAddress,
		Username:      d.Get("username").(string),
		Password:      d.Get("password").(string),
	}

	return client.ModifyLIC(r)
}

func resourceBigipLtmLicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	deviceAddress := d.Id()

	log.Println("[INFO] Reading LIC " + deviceAddress)

	members, err := client.LICs()
	if err != nil {
		return err
	}

	d.Set("device_address", members.DeviceAddress)
	d.Set("username", members.Username)
	d.Set("password", members.Password)

	return nil
}

func resourceBigipLtmLicDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmLicImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
