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
	"regexp"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	log.Printf("[DEBUG] Creating Data Group List %s", name)

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
		return fmt.Errorf("Error creating Data Group List %s: %v", name, err)
	}

	d.SetId(name)

	return resourceBigipLtmDataGroupRead(d, meta)
}

func resourceBigipLtmDataGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var records []map[string]interface{}

	name := d.Id()
	log.Printf("[DEBUG] Retrieving Data Group List %s", name)

	datagroup, err := client.GetInternalDataGroup(name)
	if err != nil {
		return fmt.Errorf("Error retrieving Data Group List %s: %v", name, err)
	}

	if datagroup == nil {
		log.Printf("[DEBUG] Data Group List %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	d.Set("name", datagroup.FullPath)
	d.Set("type", datagroup.Type)

	for _, record := range datagroup.Records {
		dgRecord := map[string]interface{}{
			"name": record.Name,
			"data": record.Data,
		}
		records = append(records, dgRecord)
	}

	if err := d.Set("record", records); err != nil {
		return fmt.Errorf("Error updating records in state for Data Group List %s: %v", name, err)
	}

	return nil
}

func resourceBigipLtmDataGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[DEBUG] Checking if Data Group List (%s) exists", name)

	datagroup, err := client.GetInternalDataGroup(name)
	if err != nil {
		return false, fmt.Errorf("Error retrieving Data Group List %s: %v", name, err)
	}

	if datagroup == nil {
		log.Printf("[DEBUG] Data Group List (%s) not found, removing from state", name)
		d.SetId("")
		return false, nil
	}

	return datagroup != nil, nil
}

func resourceBigipLtmDataGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[DEBUG] Modifying Data Group List %s", name)

	rs := d.Get("record").(*schema.Set)
	dgtype := d.Get("type").(string)

	var records []bigip.DataGroupRecord
	if rs.Len() > 0 {
		for _, r := range rs.List() {
			record := r.(map[string]interface{})
			records = append(records, bigip.DataGroupRecord{Name: record["name"].(string), Data: record["data"].(string)})
		}
	} else {
		records = nil
	}

	dgver := &bigip.DataGroup{
		Name:    name,
		Type:    dgtype,
		Records: records,
	}

	dgver1213 := &bigip.DataGroup{
		Name:    name,
		Records: records,
	}

	ver, err := client.BigipVersion()

	bigipversion := ver.Entries.HTTPSLocalhostMgmtTmCliVersion0.NestedStats.Entries.Active.Description
	re := regexp.MustCompile(`^(12)|(13).*`)
	matchresult := re.MatchString(bigipversion)
	regversion := re.FindAllString(bigipversion, -1)
	if matchresult == true {
		log.Printf("[DEBUG] Bigip version is : %s", regversion)
		err = client.ModifyInternalDataGroupRecords(dgver1213)
		if err != nil {
			return fmt.Errorf("Error modifying Data Group List %s: %v", name, err)
		}
	} else {
		err = client.ModifyInternalDataGroupRecords(dgver)
		if err != nil {
			return fmt.Errorf("Error modifying Data Group List %s: %v", name, err)
		}
	}
	return resourceBigipLtmDataGroupRead(d, meta)
}

func resourceBigipLtmDataGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[DEBUG] Deleting Data Group List %s", name)

	err := client.DeleteInternalDataGroup(name)
	if err != nil {
		return fmt.Errorf("Error deleting Data Group List %s: %v", name, err)
	}

	d.SetId("")
	return nil
}
