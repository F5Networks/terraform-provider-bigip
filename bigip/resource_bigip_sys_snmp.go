package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

// this module does not have DELETE function as there is no API for Delete
func resourceBigipNetSnmp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipNetSnmpCreate,
		Update: resourceBigipNetSnmpUpdate,
		Read:   resourceBigipNetSnmpRead,
		Delete: resourceBigipNetSnmpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipNetSnmpImporter,
		},

		Schema: map[string]*schema.Schema{
			"sys_contact": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact Person email",
				//ValidateFunc: validateF5Name,
			},
			"sys_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of the F5 ",
			},
			"allowedaddresses": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of SNMP addresses",
			},
		},
	}

}

func resourceBigipNetSnmpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	sysContact := d.Get("sys_contact").(string)
	sysLocation := d.Get("sys_location").(string)
	allowedAddresses := setToStringSlice(d.Get("allowedaddresses").(*schema.Set))

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

func resourceBigipNetSnmpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	sysContact := d.Id()

	log.Println("[INFO] Updating SNMP " + sysContact)

	r := &bigip.SNMP{
		SysContact:       sysContact,
		SysLocation:      d.Get("sys_location").(string),
		AllowedAddresses: setToStringSlice(d.Get("allowedaddresses").(*schema.Set)),
	}

	return client.ModifySNMP(r)
}

func resourceBigipNetSnmpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	sysContact := d.Id()

	log.Println("[INFO] Reading SNMP " + sysContact)

	snmp, err := client.SNMPs()
	if err != nil {
		return err
	}

	d.Set("sys_contact", snmp.SysContact)
	d.Set("sys_location", snmp.SysLocation)
	d.Set("allowedaddresses", snmp.AllowedAddresses)

	return nil
}

func resourceBigipNetSnmpDelete(d *schema.ResourceData, meta interface{}) error {
	// No API support for Delete
	return nil
}

func resourceBigipNetSnmpImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
