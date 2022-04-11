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
	"os"
	"regexp"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
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
			"internal": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set flase if you want to create External Datagroup",
				Default:     true,
			},
			"records_src": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"record"},
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
				ConflictsWith: []string{"records_src"},
			},
		},
	}
}

func resourceBigipLtmDataGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var name string

	dgtype := d.Get("type").(string)
	rs := d.Get("record").(*schema.Set)

	tmplPath := d.Get("records_src").(string)

	name = d.Get("name").(string)
	log.Printf("[DEBUG] Creating Data Group List %s", name)
	if d.Get("internal").(bool) {
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
			return fmt.Errorf("Error creating Data Group List %s: %v ", name, err)
		}
	} else {
		res := strings.Split(name, "/")
		file, fail := os.OpenFile(tmplPath, os.O_RDWR, 0644)
		if fail != nil {
			return fmt.Errorf("error in reading file: %s", fail)
		}
		err := client.UploadDatagroup(file, res[2], res[1], dgtype, true)
		defer file.Close()
		if err != nil {
			return fmt.Errorf("error in creating External Datagroup (%s): %s", name, err)
		}
	}
	d.SetId(name)
	return resourceBigipLtmDataGroupRead(d, meta)
}

func resourceBigipLtmDataGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var records []map[string]interface{}

	name := d.Id()
	log.Printf("[DEBUG] Retrieving Data Group List %s", name)
	if d.Get("internal").(bool) {
		datagroup, err := client.GetInternalDataGroup(name)
		if err != nil {
			return fmt.Errorf("Error retrieving Data Group List %s: %v ", name, err)
		}

		if datagroup == nil {
			log.Printf("[DEBUG] Data Group List %s not found, removing from state", name)
			d.SetId("")
			return nil
		}
		_ = d.Set("name", datagroup.FullPath)
		_ = d.Set("type", datagroup.Type)
		for _, record := range datagroup.Records {
			dgRecord := map[string]interface{}{
				"name": record.Name,
				"data": record.Data,
			}
			records = append(records, dgRecord)
		}
		if err := d.Set("record", records); err != nil {
			return fmt.Errorf("Error updating records in state for Data Group List %s: %v ", name, err)
		}
	} else {
		datagroup, err := client.GetExternalDataGroup(name)
		if err != nil {
			return fmt.Errorf("Error retrieving Data Group List %s: %v ", name, err)
		}

		if datagroup == nil {
			log.Printf("[DEBUG] Data Group List %s not found, removing from state", name)
			d.SetId("")
			return nil
		}
		_ = d.Set("name", datagroup.FullPath)
		_ = d.Set("type", datagroup.Type)
	}

	return nil
}

func resourceBigipLtmDataGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[DEBUG] Checking if Data Group List (%s) exists", name)

	if d.Get("internal").(bool) {
		datagroup, err := client.GetInternalDataGroup(name)
		if err != nil {
			return false, fmt.Errorf("Error retrieving Data Group List %s: %v ", name, err)
		}

		if datagroup == nil {
			log.Printf("[DEBUG] Data Group List (%s) not found, removing from state", name)
			d.SetId("")
			return false, nil
		}
	} else {
		datagroup, err := client.GetExternalDataGroup(name)
		if err != nil {
			return false, fmt.Errorf("Error retrieving Data Group List %s: %v ", name, err)
		}
		if datagroup == nil {
			log.Printf("[DEBUG] Data Group List %s not found, removing from state", name)
			d.SetId("")
			return false, nil
		}
	}
	return true, nil
}

func resourceBigipLtmDataGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[DEBUG] Modifying Data Group List %s", name)

	rs := d.Get("record").(*schema.Set)
	dgtype := d.Get("type").(string)

	if d.Get("internal").(bool) {
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
		if err != nil {
			return fmt.Errorf("Could not get BigipVersion: %v ", err)
		}

		bigipversion := ver.Entries.HTTPSLocalhostMgmtTmCliVersion0.NestedStats.Entries.Active.Description
		re := regexp.MustCompile(`^(12)|(13).*`)
		matchresult := re.MatchString(bigipversion)
		regversion := re.FindAllString(bigipversion, -1)
		if matchresult {
			log.Printf("[DEBUG] Bigip version is : %s", regversion)
			if err := client.ModifyInternalDataGroupRecords(dgver1213); err != nil {
				return fmt.Errorf("Error modifying Data Group List %s: %v ", name, err)
			}
		} else {
			log.Printf("[DEBUG] Bigip version is : %s", regversion)
			if err := client.ModifyInternalDataGroupRecords(dgver); err != nil {
				return fmt.Errorf("Error modifying Data Group List %s: %v ", name, err)
			}
		}
	} else {
		tmplPath := d.Get("records_src").(string)
		res := strings.Split(name, "/")
		file, fail := os.OpenFile(tmplPath, os.O_RDWR, 0644)
		if fail != nil {
			return fmt.Errorf("error in reading file: %s", fail)
		}
		err := client.UploadDatagroup(file, res[2], res[1], dgtype, true)
		defer file.Close()
		if err != nil {
			return fmt.Errorf("error in creating External Datagroup (%s): %s", name, err)
		}
	}
	return resourceBigipLtmDataGroupRead(d, meta)
}

func resourceBigipLtmDataGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[DEBUG] Deleting Data Group List %s", name)
	if d.Get("internal").(bool) {
		err := client.DeleteInternalDataGroup(name)
		if err != nil {
			return fmt.Errorf("Error deleting Data Group List %s: %v ", name, err)
		}
	} else {
		err := client.DeleteExternalDataGroup(name)
		if err != nil {
			return fmt.Errorf("Error deleting Data Group List %s: %v ", name, err)
		}
		err = client.DeleteExternalDatagroupfile(name)
		if err != nil {
			log.Printf("[ERROR] Unable to Delete External Datagroup file   (%s) (%v) ", name, err)
			return err
		}
	}
	d.SetId("")
	return nil
}
