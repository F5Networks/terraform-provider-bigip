package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

// this module does not have DELETE function as there is no API for Delete
func resourceBigipLtmSnmp() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmSnmpCreate,
		Update: resourceBigipLtmSnmpUpdate,
		Read:   resourceBigipLtmSnmpRead,
		Delete: resourceBigipLtmSnmpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmSnmpImporter,
		},

		Schema: map[string]*schema.Schema{
			"sysContact": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact Person email",
				//ValidateFunc: validateF5Name,
			},
			"sysLocation": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of the F5 ",
			},
			"allowedAddresses": &schema.Schema{
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of SNMP addresses",
			},
		},
	}

}

func resourceBigipLtmSnmpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	sysContact := d.Get("sysContact").(string)
	sysLocation := d.Get("sysLocation").(string)
	allowedAddresses := setToStringSlice(d.Get("allowedAddresses").(*schema.Set))

	log.Println("[INFO] Creating Snmp ")

	err := client.CreateSNMP(
		sysContact,
		sysLocation,
		allowedAddresses,
	)

	if err != nil {
		return err
	}
	d.SetId(sysContact)
	return nil
}

func resourceBigipLtmSnmpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	sysContact := d.Id()

	log.Println("[INFO] Updating SNMP " + sysContact)

	r := &bigip.SNMP{
		SysContact:       sysContact,
		SysLocation:      d.Get("sysLocation").(string),
		AllowedAddresses: setToStringSlice(d.Get("allowedAddresses").(*schema.Set)),
	}

	return client.ModifySNMP(r)
	return nil
}

func resourceBigipLtmSnmpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	sysContact := d.Id()

	log.Println("[INFO] Reading SNMP " + sysContact)

	snmp, err := client.SNMPs()
	if err != nil {
		return err
	}

	d.Set("sysContact", snmp.SysContact)
	d.Set("sysContact", snmp.SysLocation)
	d.Set("allowedAddresses", snmp.AllowedAddresses)

	return nil
}

func resourceBigipLtmSnmpDelete(d *schema.ResourceData, meta interface{}) error {
	/* This function is not supported on BIG-IP, you cannot DELETE NTP API is not supported */
	return nil
}

func resourceBigipLtmSnmpImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
