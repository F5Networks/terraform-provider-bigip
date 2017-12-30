package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
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
		},
	}

}

func resourceBigipCmDevicegroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	autoSync := d.Get("auto_sync").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	typo := d.Get("type").(string)
	fullLoadOnSync := d.Get("full_load_on_sync").(string)
	saveOnAutoSync := d.Get("save_on_auto_sync").(string)
	networkFailover := d.Get("network_failover").(string)
	incrementalConfigSyncSizeMax := d.Get("incremental_config").(int)


	log.Println("[INFO] Creating Devicegroup ")

	err := client.CreateDevicegroup(
		name,
		description,
		autoSync,
		typo,
		fullLoadOnSync,
		saveOnAutoSync,
		networkFailover,
		incrementalConfigSyncSizeMax,
	)

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

	r := &bigip.Devicegroup{
		Name:           name,
		Description:    d.Get("description").(string),
		AutoSync:       d.Get("auto_sync").(string),
		Type:           d.Get("type").(string),
		FullLoadOnSync: d.Get("full_load_on_sync").(string),
		SaveOnAutoSync: d.Get("save_on_auto_sync").(string),
		NetworkFailover: d.Get("network_failover").(string),
		IncrementalConfigSyncSizeMax: d.Get("incremental_config").(int),
	}

	return client.ModifyDevicegroup(r)
}

func resourceBigipCmDevicegroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Devicegroup " + name)

	members, err := client.Devicegroups(name)
	if err != nil {
		return err
	}
log.Println("i am in read @@@@@ @@@@@@@@@@@@@@@@ @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@   ", members)
	d.Set("name", members.Name)
	d.Set("description", members.Description)
	d.Set("auto_sync", members.AutoSync)
	d.Set("type", members.Type)
	d.Set("full_load_on_sync", members.FullLoadOnSync)
	d.Set("save_on_auto_sync", members.SaveOnAutoSync)
	d.Set("network_failover", members.NetworkFailover)
	d.Set("incremental_config", members.IncrementalConfigSyncSizeMax)
	 return nil
}

func resourceBigipCmDevicegroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	return client.DeleteDevicegroup(name)
}

func resourceBigipCmDevicegroupImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
