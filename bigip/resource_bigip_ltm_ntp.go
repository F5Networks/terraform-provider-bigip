package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmNtp() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmNtpCreate,
		Update: resourceBigipLtmNtpUpdate,
		Read:   resourceBigipLtmNtpRead,
		Delete: resourceBigipLtmNtpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmNtpImporter,
		},

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the ntp Servers",
				ValidateFunc: validateF5Name,
			},

			"servers": &schema.Schema{
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},

			"timezone": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Servers timezone",
			},
		},
	}

}

func resourceBigipLtmNtpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Get("description").(string)
	servers := setToStringSlice(d.Get("servers").(*schema.Set))
	timezone := d.Get("timezone").(string)

	log.Println("[INFO] Creating Ntp ")

	err := client.CreateNtp(
		description,
		servers,
		timezone,
	)

	if err != nil {
		return err
	}
	d.SetId(description)
	return nil
}

func resourceBigipLtmNtpUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmNtpRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmNtpDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmNtpImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
