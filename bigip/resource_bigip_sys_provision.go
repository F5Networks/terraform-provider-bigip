/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBigipSysProvision() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysProvisionCreate,
		Update: resourceBigipSysProvisionUpdate,
		Read:   resourceBigipSysProvisionRead,
		Delete: resourceBigipSysProvisionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the module to be provisioned",
			},
			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "path",
			},
			"cpu_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "cpu Ratio",
			},
			"disk_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "disk Ratio",
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "what level nominal or dedicated",
				Default:     "nominal",
			},
			"memory_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "memory Ratio",
			},
		},
	}
}

func resourceBigipSysProvisionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	fullPath := d.Get("full_path").(string)
	cpuRatio := d.Get("cpu_ratio").(int)
	diskRatio := d.Get("disk_ratio").(int)
	level := d.Get("level").(string)
	memoryRatio := d.Get("memory_ratio").(int)

	log.Printf("[INFO] Provisioning for %v module", name)

	r := &bigip.Provision{
		Name:        name,
		FullPath:    fullPath,
		CpuRatio:    cpuRatio,
		DiskRatio:   diskRatio,
		Level:       level,
		MemoryRatio: memoryRatio,
	}
	err := client.ProvisionModule(r)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Provision  (%s) ", err)
		return err
	}
	d.SetId(name)
	return resourceBigipSysProvisionRead(d, meta)
}

func resourceBigipSysProvisionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating Provisioning for :%v module", name)
	r := &bigip.Provision{
		Name:        name,
		FullPath:    d.Get("full_path").(string),
		CpuRatio:    d.Get("cpu_ratio").(int),
		DiskRatio:   d.Get("disk_ratio").(int),
		Level:       d.Get("level").(string),
		MemoryRatio: d.Get("memory_ratio").(int),
	}
	err := client.ProvisionModule(r)
	if err != nil {
		log.Printf("[ERROR] Unable to Update Provision (%v) ", err)
		return err
	}
	return resourceBigipSysProvisionRead(d, meta)
}

func resourceBigipSysProvisionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading Provisions " + name)
	p, err := client.Provisions(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Provision (%s) (%v) ", name, err)
		return err
	}
	if p == nil {
		log.Printf("[WARN] Provision (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err := d.Set("full_path", p.FullPath); err != nil {
		return fmt.Errorf("[DEBUG] Error saving FullPath to state for Provision  (%s): %s", d.Id(), err)
	}
	d.Set("cpu_ratio", p.CpuRatio)
	d.Set("disk_ratio", p.DiskRatio)
	d.Set("level", p.Level)
	d.Set("memory_ratio", p.MemoryRatio)

	return nil
}

func resourceBigipSysProvisionDelete(d *schema.ResourceData, meta interface{}) error {
	//API is not supported for Deleting
	return nil
}
