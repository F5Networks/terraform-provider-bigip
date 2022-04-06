/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipCmDevicegroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipCmDevicegroupCreate,
		Update: resourceBigipCmDevicegroupUpdate,
		Read:   resourceBigipCmDevicegroupRead,
		Delete: resourceBigipCmDevicegroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the Device group",
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Device administrative partition",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of Device group",
			},

			"auto_sync": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Specifies if the device-group will automatically sync configuration data to its members",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "sync-only",
				Description: "Specifies if the device-group will be used for failover or resource syncing",
			},
			"full_load_on_sync": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "false",
				Description: "Specifies if the device-group will perform a full-load upon sync",
			},
			"save_on_auto_sync": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "false",
				Description: "Specifies whether the configuration should be saved upon auto-sync.",
			},

			"network_failover": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "Specifies if the device-group will use a network connection for failover",
			},

			"incremental_config": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1024,
				Description: "Specifies the maximum size (in KB) to devote to incremental config sync cached transactions. The default is 1024 KB.",
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
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
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
		log.Printf("[ERROR] Unable to Create Devicegroup (%s) (%v) ", name, err)
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
	err := client.UpdateDevicegroup(name, &p)
	if err != nil {
		log.Printf("[ERROR] Unable to Update Devicegroup (%s) (%v) ", name, err)
		return err
	}
	return resourceBigipCmDevicegroupRead(d, meta)
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
		if _, err := client.DevicegroupsDevices(name, Rname); err != nil {
			log.Printf("[ERROR] Unable to retrieve DevicegroupsDevices (%s,%s) (%v) ", name, Rname, err)
		}
	}

	p, err := client.Devicegroups(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Devicegroup (%s) (%v) ", name, err)
		return err
	}

	if p == nil {
		log.Printf("[WARN] Devicegroup (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", p.Name)
	_ = d.Set("description", p.Description)
	if err := d.Set("auto_sync", p.AutoSync); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AutoSync  to state for Devicegroup (%s): %s", d.Id(), err)
	}

	if err := d.Set("type", p.Type); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Type  to state for Devicegroup (%s): %s", d.Id(), err)
	}
	_ = d.Set("fullLoadOnSync", p.FullLoadOnSync)
	_ = d.Set("saveOnAutoSync", p.SaveOnAutoSync)
	_ = d.Set("incrementalConfigSyncSizeMax", p.IncrementalConfigSyncSizeMax)
	_ = d.Set("networkFailover", p.NetworkFailover)
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
		err := client.DeleteDevicegroupDevices(name, Rname)
		if err != nil {
			log.Printf("[ERROR] Unable to Delete Deviceg (%s)  (%v) ", Rname, err)
			return err
		}

	}

	err := client.DeleteDevicegroup(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Devicegroup (%s)  (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
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
	_ = d.Set("name", p.Name)
	_ = d.Set("partition", p.Partition)
	_ = d.Set("auto_sync", p.AutoSync)
	_ = d.Set("description", p.Description)
	_ = d.Set("type", p.Type)
	_ = d.Set("full_load_on_sync", p.FullLoadOnSync)
	_ = d.Set("save_on_auto_sync", p.SaveOnAutoSync)
	_ = d.Set("network_failover", p.NetworkFailover)
	_ = d.Set("incremental_config", p.IncrementalConfigSyncSizeMax)

	for i, r := range p.Deviceb {
		device := fmt.Sprintf("device.%d", i)
		_ = d.Set(fmt.Sprintf("%s.name", device), r.Name)
		_ = d.Set(fmt.Sprintf("%s.set_sync_leader", device), r.SetSyncLeader)

	}
	return nil
}
