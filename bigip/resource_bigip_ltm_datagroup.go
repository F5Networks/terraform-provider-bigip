package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmDataGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmDataGroupCreate,
		Read:   resourceBigipLtmDataGroupRead,
		Update: resourceBigipLtmDataGroupUpdate,
		Delete: resourceBigipLtmDataGroupDelete,
		Exists: resourceBigipLtmDataGroupExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the Data Group List",
				ValidateFunc: validateF5Name,
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The Data Group type (string, ip, integer)",
				ValidateFunc: validateDataGroupType,
			},

			"record": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"data": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceBigipLtmDataGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Println("[INFO] Creating Data Group List " + name)

	dgtype := d.Get("type").(string)
	rs := d.Get("record").(*schema.Set)

	var records []bigip.DataGroupRecord
	if rs.Len() > 0 {
		for _, r := range rs.List() {
			record := r.(map[string]interface{})
			records = append(records, bigip.DataGroupRecord{Name: record["name"].(string), Data: record["data"].(string)})
		}
	} else {
		records = nil
	}

	dg := &bigip.DataGroup{
		Name:    name,
		Type:    dgtype,
		Records: records,
	}

	err := client.AddInternalDataGroup(dg)
	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceBigipLtmDataGroupRead(d, meta)
}

func resourceBigipLtmDataGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	datagroup, err := client.GetInternalDataGroup(name)
	if err != nil {
		log.Printf("Error reading Internal Data group : %s", err)
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if datagroup == nil {
		log.Printf("[WARN] Data Group List (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	return nil
}

func resourceBigipLtmDataGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching Data Group " + name)

	datagroup, err := client.GetInternalDataGroup(name)
	if err != nil {
		return false, err
	}
	if datagroup == nil {
		log.Printf("[WARN] Data Group List (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, nil
	}
	return datagroup != nil, nil
}

func resourceBigipLtmDataGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Modifying Data Group " + name)

	rs := d.Get("record").(*schema.Set)

	var records []bigip.DataGroupRecord
	if rs.Len() > 0 {
		for _, r := range rs.List() {
			record := r.(map[string]interface{})
			records = append(records, bigip.DataGroupRecord{Name: record["name"].(string), Data: record["data"].(string)})
		}
	} else {
		records = nil
	}

	err := client.ModifyInternalDataGroupRecords(name, records)
	if err != nil {
		return err
	}
	if err == nil {
		log.Printf("[WARN] Data Group List (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return resourceBigipLtmDataGroupRead(d, meta)
}

func resourceBigipLtmDataGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting Data Group " + name)

	err := client.DeleteInternalDataGroup(name)
	if err != nil {
		log.Printf("Error Destroying  Internal Data Group : %s", err)
		d.SetId("")
		return nil
	}

	if err == nil {
		log.Printf("[WARN] Datag Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil
}
