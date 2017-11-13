package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipSysProvision() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysProvisionCreate,
		Update: resourceBigipSysProvisionUpdate,
		Read:   resourceBigipSysProvisionRead,
		Delete: resourceBigipSysProvisionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipSysProvisionImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the module to be provisioned",
				ValidateFunc: validateF5Name,
			},

			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
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

	log.Println("[INFO] Provisioning  ")

	err := client.CreateProvision(
		name,
		fullPath,
		cpuRatio,
		diskRatio,
		level,
		memoryRatio,
	)

	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func resourceBigipSysProvisionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Provsioning " + name)

	r := &bigip.Provision{
		Name:        name,
		FullPath:    d.Get("full_path").(string),
		CpuRatio:    d.Get("cpu_ratio").(int),
		DiskRatio:   d.Get("disk_ratio").(int),
		Level:       d.Get("level").(string),
		MemoryRatio: d.Get("memory_ratio").(int),
	}

	return client.ModifyProvision(r)
}

func resourceBigipSysProvisionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Provision " + name)

	provision, err := client.Provisions()
	if err != nil {
		return err
	}

	d.Set("name", provision.Name)
	d.Set("full_path", provision.FullPath)
	d.Set("cpu_ratio", provision.CpuRatio)
	d.Set("disk_ratio", provision.DiskRatio)
	d.Set("level", provision.Level)
	d.Set("memory_ratio", provision.MemoryRatio)

	return nil
}

func resourceBigipSysProvisionDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipSysProvisionImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
