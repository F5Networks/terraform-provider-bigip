package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipGtmMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipGtmMonitorCreate,
		Update: resourceBigipGtmMonitorUpdate,
		Read:   resourceBigipGtmMonitorRead,
		Delete: resourceBigipGtmMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipGtmMonitorImporter,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Gtmmonitor  Name",
				//	ValidateFunc: validateF5Name,
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Which partition on BIG-IP",
			},

			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "BIG-IP autolasthop",
			},
			"probe_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"recv": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"send": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},

		},
	}
}

func resourceBigipGtmMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	defaults_from := d.Get("defaults_from").(string)
	interval := d.Get("interval").(int)
	probeTimeout := d.Get("probe_timeout").(int)
	recv := d.Get("recv").(string)
	send := d.Get("send").(string)
	log.Println("[INFO] Creating Gtmmonitor ")

	err := client.CreateGtmmonitor(
		name,
		defaults_from,
		interval,
		probeTimeout,
		recv,
		send,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipGtmMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Gtmmonitor " + name)

	r := &bigip.Gtmmonitor{
		Name:          d.Get("name").(string),
		Defaults_from:   d.Get("defaults_from").(string),
		Interval:   d.Get("interval").(int),
		Probe_timeout:        d.Get("probe_timeout").(int),
		Recv:    d.Get("recv").(string),
		Send:   d.Get("send").(string),
	}

	return client.ModifyGtmmonitor(r)
}

func resourceBigipGtmMonitorRead(d *schema.ResourceData, meta interface{}) error {
	/*client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Monitorlist " + name)

	Monitor, err := client.GetMonitor(name)
	if err != nil {
		return err
	}
	d.Set("origins", Monitor.Origins)
	d.Set("name", name)
	*/
	return nil
}

func resourceBigipGtmMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	return client.DeleteGtmmonitor(name)
	//return nil
}

func resourceBigipGtmMonitorImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
