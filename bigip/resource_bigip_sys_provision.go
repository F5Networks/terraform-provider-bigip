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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of module to provision in BIG-IP.",
				ValidateFunc: validation.StringInSlice([]string{"afm", "am", "apm", "asm", "avr", "cgnat", "dos", "fps", "gtm", "ilx", "lc", "ltm", "pem", "sslo", "swg", "urldb"}, false),
			},
			"full_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//ValidateFunc: validation.StringInSlice([]string{"afm", "am", "apm","asm","avr","dos","fps","gtm","ilx","lc","ltm","pem", "sslo" ,"swg","urldb"}, false),
			},
			"cpu_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use this option only when the level option is set to custom.F5 Networks recommends that you do not modify this option. The default value is none",
			},
			"disk_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use this option only when the level option is set to custom.F5 Networks recommends that you do not modify this option. The default value is none",
			},
			"level": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Sets the provisioning level for the requested modules. Changing the level for one module may require modifying the level of another module. For example, changing one module to dedicated requires setting all others to none. Setting the level of a module to none means the module is not activated.",
				Default:      "nominal",
				ValidateFunc: validation.StringInSlice([]string{"nominal", "none", "minimum", "dedicated"}, false),
			},
			"memory_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use this option only when the level option is set to custom.F5 Networks recommends that you do not modify this option. The default value is none",
			},
		},
	}
}

func resourceBigipSysProvisionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	log.Printf("[INFO] Provisioning for %v module", name)

	pss := &bigip.Provision{
		Name: name,
	}
	config := getsysProvisionConfig(d, pss)

	err := client.ProvisionModule(config)
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

	pss := &bigip.Provision{
		Name: name,
	}
	config := getsysProvisionConfig(d, pss)

	err := client.ProvisionModule(config)
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
	_ = d.Set("cpu_ratio", p.CpuRatio)
	_ = d.Set("disk_ratio", p.DiskRatio)
	_ = d.Set("level", p.Level)
	_ = d.Set("memory_ratio", p.MemoryRatio)

	return nil
}

func resourceBigipSysProvisionDelete(d *schema.ResourceData, meta interface{}) error {
	// API is not supported for Deleting
	return nil
}

func getsysProvisionConfig(d *schema.ResourceData, config *bigip.Provision) *bigip.Provision {
	config.FullPath = d.Get("full_path").(string)
	config.CpuRatio = d.Get("cpu_ratio").(int)
	config.DiskRatio = d.Get("disk_ratio").(int)
	config.Level = d.Get("level").(string)
	config.MemoryRatio = d.Get("memory_ratio").(int)
	return config
}
