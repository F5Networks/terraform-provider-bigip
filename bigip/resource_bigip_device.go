package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmDevice() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmDeviceCreate,
		Update: resourceBigipLtmDeviceUpdate,
		Read:   resourceBigipLtmDeviceRead,
		Delete: resourceBigipLtmDeviceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmDeviceImporter,
		},

		Schema: map[string]*schema.Schema{

			"configsync_ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of BIG-IP",
				//	ValidateFunc: validateF5Name,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Device which needs to be Deviceensed",
			},

			"mirror_ip": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"mirror_secondary_ip": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
		},
	}

}

func resourceBigipLtmDeviceCreate(d *schema.ResourceData, meta interface{}) error {
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

func resourceBigipLtmDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
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

func resourceBigipLtmDeviceRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceBigipLtmDeviceDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmDeviceImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
