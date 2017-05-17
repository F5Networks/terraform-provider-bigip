package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmProvision() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmProvisionCreate,
		Update: resourceBigipLtmProvisionUpdate,
		Read:   resourceBigipLtmProvisionRead,
		Delete: resourceBigipLtmProvisionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmProvisionImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the module to be provisioned",
				ValidateFunc: validateF5Name,
			},

			"fullPath": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "path",
			},

			"cpuRatio": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "cpu Ratio",
			},

			"diskRatio": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "disk Ratio",
			},

			"level": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "what level nominal or dedicated",
			},

			"memoryRatio": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "memory Ratio",
			},
		},
	}

}

func resourceBigipLtmProvisionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	fullPath := d.Get("fullPath").(string)
	cpuRatio := d.Get("cpuRatio").(int)
	diskRatio := d.Get("diskRatio").(int)
	level := d.Get("level").(string)
	memoryRatio := d.Get("memoryRatio").(int)

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

func resourceBigipLtmProvisionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Provsioning " + name)

	r := &bigip.Provision{
		Name:        name,
		FullPath:    d.Get("fullPath").(string),
		CpuRatio:    d.Get("cpuRatio").(int),
		DiskRatio:   d.Get("diskRatio").(int),
		Level:       d.Get("level").(string),
		MemoryRatio: d.Get("memoryRatio").(int),
	}

	return client.ModifyProvision(r)
}

func resourceBigipLtmProvisionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Reading Provision " + name)

	provision, err := client.Provisions()
	if err != nil {
		return err
	}

	d.Set("name", provision.Name)
	d.Set("fullPath", provision.FullPath)
	d.Set("cpuRatio", provision.CpuRatio)
	d.Set("diskRatio", provision.DiskRatio)
	d.Set("level", provision.Level)
	d.Set("memoryRatio", provision.MemoryRatio)

	return nil
}

func resourceBigipLtmProvisionDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmProvisionImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
