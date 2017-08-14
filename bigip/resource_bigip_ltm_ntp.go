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

	err := client.CreateNTP(
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
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Updating NTP " + description)

	r := &bigip.NTP{
		Description: description,
		Servers:     setToStringSlice(d.Get("servers").(*schema.Set)),
		Timezone:    d.Get("timezone").(string),
	}

	return client.ModifyNTP(r)
}

func resourceBigipLtmNtpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading NTP " + description)

	ntp, err := client.NTPs()
	if err != nil {
		return err
	}

	d.Set("description", ntp.Description)
	d.Set("servers", ntp.Servers)
	d.Set("timezone", ntp.Timezone)

	return nil
}

func resourceBigipLtmNtpDelete(d *schema.ResourceData, meta interface{}) error {
	/* This function is not supported on BIG-IP, you cannot DELETE NTP API is not supported */
	return nil
}

func resourceBigipLtmNtpImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
