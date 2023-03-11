/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipLtmDataGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipLtmDataGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Data Group List",
			},
			"partition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition of resource group",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Data Group type (string, ip, integer)",
			},
			"record": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
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

func dataSourceBigipLtmDataGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")
	var records []map[string]interface{}
	name := fmt.Sprintf("/%s/%s", d.Get("partition").(string), d.Get("name").(string))
	log.Printf("[INFO] Retrieving Data Group List %s", name)
	dataGroup, err := client.GetInternalDataGroup(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving Data Group List %s: %v ", name, err))
	}
	if dataGroup == nil {
		log.Printf("[DEBUG] Data Group List %s not found, removing from state", name)
		d.SetId("")
		return nil
	}
	_ = d.Set("name", dataGroup.Name)
	_ = d.Set("partition", dataGroup.Partition)
	_ = d.Set("type", dataGroup.Type)
	for _, record := range dataGroup.Records {
		dgRecord := map[string]interface{}{
			"name": record.Name,
			"data": record.Data,
		}
		records = append(records, dgRecord)
	}
	if err := d.Set("record", records); err != nil {
		return diag.FromErr(fmt.Errorf("Error updating records in state for Data Group List %s: %v ", name, err))
	}
	d.SetId(dataGroup.FullPath)
	return nil
}
