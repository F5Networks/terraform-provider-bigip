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

			"configsyncIp": &schema.Schema{
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

			"mirrorIp": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"mirrorSecondaryIp": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
		},
	}

}

func resourceBigipLtmDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	configsyncIp := d.Get("configsyncIp").(string)
	name := d.Get("name").(string)
	mirrorIp := d.Get("mirrorIp").(string)
	mirrorSecondaryIp := d.Get("mirrorSecondaryIp").(string)

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
		ConfigsyncIp:      d.Get("configsyncIp").(string),
		MirrorIp:          d.Get("mirrorIp").(string),
		MirrorSecondaryIp: d.Get("mirrorSecondaryIp").(string),
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
	d.Set("mirrorIp", members.MirrorIp)
	d.Set("configsyncIp", members.ConfigsyncIp)
	d.Set("mirrorSecondaryIp", members.MirrorSecondaryIp)
	return nil
}

func resourceBigipLtmDeviceDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmDeviceImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
