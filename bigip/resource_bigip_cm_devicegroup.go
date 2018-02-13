package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceBigipCmDevicegroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipCmDevicegroupCreate,
		Update: resourceBigipCmDevicegroupUpdate,
		Read:   resourceBigipCmDevicegroupRead,
		Delete: resourceBigipCmDevicegroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipCmDevicegroupImporter,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the Devicegroup which needs to be Devicegroupensed",
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the Devicegroup which needs to be Devicegroupensed",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the Devicegroup which needs to be Devicegroupensed",
			},

			"auto_sync": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "BIG-IP password",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "sync-only",
				Description: "BIG-IP password",
			},
			"full_load_on_sync": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "false",
				Description: "BIG-IP password",
			},
			"save_on_auto_sync": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "false",
				Description: "BIG-IP password",
			},

			"network_failover": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "BIG-IP password",
			},

			"incremental_config": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1024,
				Description: "BIG-IP password",
			},
			"device": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"set_sync_leader": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
							//ValidateFunc: validateF5Name,
						},
					},
				},
			},
		},
	}
}

func resourceBigipCmDevicegroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating Device Group" + name)

	p := dataToDevicegroup(name, d)
	d.SetId(name)
	err := client.CreateDevicegroup(&p)

	log.Println("[INFO] Creating Devicegroup ")

	if err != nil {
		return err
	}
	d.SetId(name)
	return resourceBigipCmDevicegroupRead(d, meta)
}

func resourceBigipCmDevicegroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating Devicegroup " + name)
	p := dataToDevicegroup(name, d)
	return client.UpdateDevicegroup(name, &p)
}

func resourceBigipCmDevicegroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Devicegroup " + name)
	deviceCount := d.Get("device.#").(int)
	for i := 0; i < deviceCount; i++ {
		var r bigip.Devicerecord
		prefix := fmt.Sprintf("device.%d", i)
		r.Name = d.Get(prefix + ".name").(string)
		Rname := r.Name
		log.Println(" my rname is  ", Rname)
		client.DevicegroupsDevices(name, Rname)
	}

	p, err := client.Devicegroups(name)
	if err != nil {
		return err
	}

	if p == nil {
			log.Printf("[WARN] Devicegroup (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
	d.Set("name", p.Name)
	d.Set("description", p.Description)
	d.Set("type", p.Type)
	d.Set("fullLoadOnSync", p.FullLoadOnSync)
	d.Set("saveOnAutoSync", p.SaveOnAutoSync)
	d.Set("incrementalConfigSyncSizeMax", p.IncrementalConfigSyncSizeMax)
	d.Set("networkFailover", p.NetworkFailover)

	return nil

}

func resourceBigipCmDevicegroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	deviceCount := d.Get("device.#").(int)
	for i := 0; i < deviceCount; i++ {
		var r bigip.Devicerecord
		prefix := fmt.Sprintf("device.%d", i)
		r.Name = d.Get(prefix + ".name").(string)
		Rname := r.Name
		log.Println(" my rname is  ", Rname)
		client.DeleteDevicegroupDevices(name, Rname)
	}

	err := client.DeleteDevicegroup(name)
	if err != nil {
		return err
	}
	if err == nil {
		log.Printf("[WARN] Devicegroup  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}

func resourceBigipCmDevicegroupImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func dataToDevicegroup(name string, d *schema.ResourceData) bigip.Devicegroup {
	var p bigip.Devicegroup

	p.Name = name
	p.Partition = d.Get("partition").(string)
	p.AutoSync = d.Get("auto_sync").(string)
	p.Description = d.Get("description").(string)
	p.Type = d.Get("type").(string)
	p.FullLoadOnSync = d.Get("full_load_on_sync").(string)
	p.SaveOnAutoSync = d.Get("save_on_auto_sync").(string)
	p.NetworkFailover = d.Get("network_failover").(string)
	p.IncrementalConfigSyncSizeMax = d.Get("incremental_config").(int)
	deviceCount := d.Get("device.#").(int)
	p.Deviceb = make([]bigip.Devicerecord, 0, deviceCount)
	for i := 0; i < deviceCount; i++ {
		var r bigip.Devicerecord
		log.Println("I am in dattodevicegroup policy ", p, deviceCount, i)
		prefix := fmt.Sprintf("device.%d", i)
		r.Name = d.Get(prefix + ".name").(string)
		p.Deviceb = append(p.Deviceb, r)
	}

	return p
}

func DevicegroupToData(p *bigip.Devicegroup, d *schema.ResourceData) error {
	d.Set("name", p.Name)
	d.Set("partition", p.Partition)
	d.Set("auto_sync", p.AutoSync)
	d.Set("description", p.Description)
	d.Set("type", p.Type)
	d.Set("full_load_on_sync", p.FullLoadOnSync)
	d.Set("save_on_auto_sync", p.SaveOnAutoSync)
	d.Set("network_failover", p.NetworkFailover)
	d.Set("incremental_config", p.IncrementalConfigSyncSizeMax)

	for i, r := range p.Deviceb {
		device := fmt.Sprintf("device.%d", i)

		d.Set(fmt.Sprintf("%s.name", device), r.Name)
		d.Set(fmt.Sprintf("%s.set_sync_leader", device), r.SetSyncLeader)

	}
	return nil
}
