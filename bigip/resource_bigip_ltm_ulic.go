package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmULic() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmULicCreate,
		Update: resourceBigipLtmULicUpdate,
		Read:   resourceBigipLtmULicRead,
		Delete: resourceBigipLtmULicDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmULicImporter,
		},

		Schema: map[string]*schema.Schema{

			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of BIG-IP",
				//	ValidateFunc: validateF5Name,
			},
			"deviceAddress": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Device which needs to be licensed",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"unitOfMeasure": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP subscription yearly/hourly",
			},
		},
	}

}

func resourceBigipLtmULicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	deviceAddress := d.Get("deviceAddress").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	unitOfMeasure := d.Get("unitOfMeasure").(string)
	log.Println("[INFO] Creating Lic ")

	err := client.CreateULIC(
		deviceAddress,
		username,
		password,
		unitOfMeasure,
	)

	if err != nil {
		return err
	}
	d.SetId(deviceAddress)
	return nil
}

func resourceBigipLtmULicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	deviceAddress := d.Id()

	log.Println("[INFO] Updating ULIC " + deviceAddress)

	r := &bigip.ULIC{
		DeviceAddress: deviceAddress,
		Username:      d.Get("username").(string),
		Password:      d.Get("password").(string),
		UnitOfMeasure: d.Get("unitOfMeasure").(string),
	}

	return client.ModifyULIC(r)
}

func resourceBigipLtmULicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	deviceAddress := d.Id()

	log.Println("[INFO] Reading ULIC " + deviceAddress)

	members, err := client.ULICs()
	if err != nil {
		return err
	}

	d.Set("deviceAddress", members.DeviceAddress)
	d.Set("username", members.Username)
	d.Set("password", members.Password)
	//d.set("unitOfMeasure", members.UnitOfMeasure)

	return nil
}

func resourceBigipLtmULicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	deviceAddress := d.Id()

	log.Println("[INFO] Updating ULIC " + deviceAddress)

	r := &bigip.ULIC{
		DeviceAddress: deviceAddress,
		Username:      d.Get("username").(string),
		Password:      d.Get("password").(string),
		UnitOfMeasure: d.Get("unitOfMeasure").(string),
	}

	return client.DeleteULIC(r)
}

func resourceBigipLtmULicImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
