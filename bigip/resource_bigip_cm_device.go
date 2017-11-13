package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipCmDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipCmDeviceCreate,
		Update: resourceBigipCmDeviceUpdate,
		Read:   resourceBigipCmDeviceRead,
		Delete: resourceBigipCmDeviceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipCmDeviceImporter,
		},

		Schema: map[string]*schema.Schema{

			"configsync_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of BIG-IP",
				//	ValidateFunc: validateF5Name,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Device which needs to be Deviceensed",
			},

			"mirror_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"mirror_secondary_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
		},
	}

}

func resourceBigipCmDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	configsyncIp := d.Get("configsync_ip").(string)
	name := d.Get("name").(string)
	mirrorIp := d.Get("mirror_ip").(string)
	mirrorSecondaryIp := d.Get("mirror_secondary_ip").(string)

	log.Println("[INFO] Creating Device ")

	err := client.CreateDevice(
		name,
		configsyncIp,
		mirrorIp,
		mirrorSecondaryIp,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipCmDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Device " + name)

	r := &bigip.Device{
		Name:              name,
		ConfigsyncIp:      d.Get("configsync_ip").(string),
		MirrorIp:          d.Get("mirror_ip").(string),
		MirrorSecondaryIp: d.Get("mirror_secondary_ip").(string),
	}

	return client.ModifyDevice(r)
}

func resourceBigipCmDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Device " + name)

	members, err := client.Devices()
	if err != nil {
		return err
	}

	d.Set("name", members.Name)
	d.Set("mirror_ip", members.MirrorIp)
	d.Set("configsync_ip", members.ConfigsyncIp)
	d.Set("mirror_secondary_ip", members.MirrorSecondaryIp)
	return nil
}

func resourceBigipCmDeviceDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipCmDeviceImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
