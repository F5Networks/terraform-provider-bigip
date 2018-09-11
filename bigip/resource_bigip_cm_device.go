package bigip

import (
	"fmt"
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
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"configsync_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address used for config sync",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Device which needs to be Deviceensed",
			},

			"mirror_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address used for state mirroring",
			},
			"mirror_secondary_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secondary IP address used for state mirroring",
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
		log.Printf("[WARN] Not able to Create Device (%s) ", err)
		return err
	}
	d.SetId(name)
	return resourceBigipCmDeviceRead(d, meta)

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

	err := client.ModifyDevice(r)
	if err != nil {
		log.Printf("[WARN] Not able to Modidy Device (%s) ", err)
		return err
	}
	return resourceBigipCmDeviceRead(d, meta)
}

func resourceBigipCmDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Device " + name)

	members, err := client.Devices(name)
	if err != nil {
		log.Printf("[WARN] Not able to retrive Device (%s) ", err)
		return err
	}
	if members == nil {
		log.Printf("[WARN] Device (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", members.Name)

	if err := d.Set("mirror_ip", members.MirrorIp); err != nil {
		return fmt.Errorf("[DEBUG] Error saving mirror_ip  to state for Device (%s): %s", d.Id(), err)
	}

	if err := d.Set("configsync_ip", members.ConfigsyncIp); err != nil {
		return fmt.Errorf("[DEBUG] Error saving configsync_ip  to state for Device (%s): %s", d.Id(), err)
	}

	if err := d.Set("mirror_secondary_ip", members.MirrorSecondaryIp); err != nil {
		return fmt.Errorf("[DEBUG] Error saving mirror_secondary_ip  to state for Device (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceBigipCmDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeleteDevice(name)
	if err != nil {
		log.Printf("[WARN] Not able to Delete Device (%s)  (%s) ", name, err)
		return err
	}
	if err == nil {
		log.Printf("[WARN] Device (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}
