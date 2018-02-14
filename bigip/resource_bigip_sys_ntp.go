package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipSysNtp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysNtpCreate,
		Update: resourceBigipSysNtpUpdate,
		Read:   resourceBigipSysNtpRead,
		Delete: resourceBigipSysNtpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipSysNtpImporter,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the ntp Servers",
				ValidateFunc: validateF5Name,
			},

			"servers": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},

			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Servers timezone",
			},
		},
	}

}

func resourceBigipSysNtpCreate(d *schema.ResourceData, meta interface{}) error {
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

func resourceBigipSysNtpUpdate(d *schema.ResourceData, meta interface{}) error {
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

func resourceBigipSysNtpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading NTP " + description)

	ntp, err := client.NTPs()
	if err != nil {
		return err
	}
	if ntp == nil {
			log.Printf("[WARN] NTP (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

	d.Set("description", ntp.Description)
	d.Set("servers", ntp.Servers)
	d.Set("timezone", ntp.Timezone)

	return nil
}

func resourceBigipSysNtpDelete(d *schema.ResourceData, meta interface{}) error {
	/* This function is not supported on BIG-IP, you cannot DELETE NTP API is not supported */
	return nil
}

func resourceBigipSysNtpImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
