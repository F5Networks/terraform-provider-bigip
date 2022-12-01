/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"os"
	"regexp"
	"strings"
)

// Warning or errors can be collected in a slice type
var diags diag.Diagnostics

func resourceBigipLtmDataGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages internal (in-line)/external datagroup configuration",
		CreateContext: resourceBigipLtmDataGroupCreate,
		ReadContext:   resourceBigipLtmDataGroupRead,
		UpdateContext: resourceBigipLtmDataGroupUpdate,
		DeleteContext: resourceBigipLtmDataGroupDelete,
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
				Description: "Set false if you want to create External Datagroup",
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

func resourceBigipLtmDataGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	var name string
	dgtype := d.Get("type").(string)
	rs := d.Get("record").(*schema.Set)
	tmplPath := d.Get("records_src").(string)
	name = d.Get("name").(string)

	tflog.Info(ctx, fmt.Sprintf("Creating Data Group List:%+v", name))
	var diags diag.Diagnostics
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
			return diag.Errorf("Error creating Data Group List %s: %v ", name, err)
		}
	} else {
		res := strings.Split(name, "/")
		file, fail := os.OpenFile(tmplPath, os.O_RDWR, 0644)
		if fail != nil {
			return diag.Errorf("error in reading file: %s", fail)
		}
		err := client.UploadDatagroup(file, res[2], res[1], dgtype, true)
		defer file.Close()
		if err != nil {
			return diag.Errorf("error in creating External Datagroup (%s): %s", name, err)
		}
	}
	d.SetId(name)
	resourceBigipLtmDataGroupRead(ctx, d, meta)
	return diags
}

func resourceBigipLtmDataGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	var records []map[string]interface{}

	name := strings.Split(d.Id(), ":")

	log.Printf("[INFO] Retrieving Data Group List %s", name)
	if d.Get("internal").(bool) {
		datagroup, err := client.GetInternalDataGroup(name[0])
		if err != nil {
			return diag.Errorf("Error retrieving Data Group List %s: %v ", name, err)
		}
		if datagroup == nil {
			log.Printf("[DEBUG] Data Group List %s not found, removing from state", name)
			d.SetId("")
			return diags
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
			return diag.Errorf("Error updating records in state for Data Group List %s: %v ", name, err)
		}
	} else {
		if len(name) > 1 && name[1] == "external" {
			datagroup, err := client.GetExternalDataGroup(name[0])
			if err != nil {
				return diag.Errorf("Error retrieving Data Group List %s: %v ", name, err)
			}
			if datagroup == nil {
				log.Printf("[DEBUG] Data Group List %s not found, removing from state", name)
				d.SetId("")
				return diags
			}
			_ = d.Set("name", datagroup.FullPath)
			_ = d.Set("type", datagroup.Type)
		} else if len(name) > 1 && name[1] == "internal" {
			datagroup, err := client.GetInternalDataGroup(name[0])
			if err != nil {
				return diag.Errorf("Error retrieving Data Group List %s: %v ", name, err)
			}
			if datagroup == nil {
				log.Printf("[DEBUG] Data Group List %s not found, removing from state", name)
				d.SetId("")
				return diags
			}
			_ = d.Set("name", datagroup.FullPath)
			_ = d.Set("type", datagroup.Type)
		} else {
			datagroup, err := client.GetExternalDataGroup(name[0])
			if err != nil {
				return diag.Errorf("Error retrieving Data Group List %s: %v ", name, err)
			}
			if datagroup == nil {
				log.Printf("[DEBUG] Data Group List %s not found, removing from state", name)
				d.SetId("")
				return diags
			}
			_ = d.Set("name", datagroup.FullPath)
			_ = d.Set("type", datagroup.Type)
		}
	}
	return diags
}

func resourceBigipLtmDataGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[INFO] Modifying Data Group List %s", name)

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
			return diag.Errorf("Could not get BigipVersion: %v ", err)
		}

		bigipversion := ver.Entries.HTTPSLocalhostMgmtTmCliVersion0.NestedStats.Entries.Active.Description
		re := regexp.MustCompile(`^(12)|(13).*`)
		matchresult := re.MatchString(bigipversion)
		regversion := re.FindAllString(bigipversion, -1)
		if matchresult {
			log.Printf("[DEBUG] Bigip version is : %s", regversion)
			if err := client.ModifyInternalDataGroupRecords(dgver1213); err != nil {
				return diag.Errorf("Error modifying Data Group List %s: %v ", name, err)
			}
		} else {
			log.Printf("[DEBUG] Bigip version is : %s", regversion)
			if err := client.ModifyInternalDataGroupRecords(dgver); err != nil {
				return diag.Errorf("Error modifying Data Group List %s: %v ", name, err)
			}
		}
	} else {
		tmplPath := d.Get("records_src").(string)
		res := strings.Split(name, "/")
		file, fail := os.OpenFile(tmplPath, os.O_RDWR, 0644)
		if fail != nil {
			return diag.Errorf("error in reading file: %s", fail)
		}
		err := client.UploadDatagroup(file, res[2], res[1], dgtype, true)
		defer file.Close()
		if err != nil {
			return diag.Errorf("error in creating External Datagroup (%s): %s", name, err)
		}
	}
	return resourceBigipLtmDataGroupRead(ctx, d, meta)
}

func resourceBigipLtmDataGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[INFO] Deleting Data Group List %s", name)
	if d.Get("internal").(bool) {
		err := client.DeleteInternalDataGroup(name)
		if err != nil {
			return diag.Errorf("Error deleting Data Group List %s: %v ", name, err)
		}
	} else {
		err := client.DeleteExternalDataGroup(name)
		if err != nil {
			return diag.Errorf("Error deleting Data Group List %s: %v ", name, err)
		}
		err = client.DeleteExternalDatagroupfile(name)
		if err != nil {
			log.Printf("[ERROR] Unable to Delete External Datagroup file   (%s) (%v) ", name, err)
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}
