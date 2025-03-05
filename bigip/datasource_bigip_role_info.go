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

func dataSourceBigipRoleInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipRoleInfoRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the role info",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the role info",
			},
			"attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The attribute of the role info",
			},
			"console": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The console of the role info",
			},
			"deny": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The deny of the role info",
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The role of the role info",
			},
			"user_partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user partition of the role info",
			},
			"line_order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The line order of the role info",
			},
		},
	}
}

func dataSourceBigipRoleInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	roleInfo, err := client.GetRoleInfo(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving Role Info %s: %v ", name, err))
	}
	if roleInfo == nil {
		log.Printf("[DEBUG] Role Info (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", roleInfo.Name)
	_ = d.Set("attribute", roleInfo.Attribute)
	_ = d.Set("console", roleInfo.Console)
	_ = d.Set("deny", roleInfo.Deny)
	_ = d.Set("description", roleInfo.Description)
	_ = d.Set("line_order", roleInfo.LineOrder)
	_ = d.Set("role", roleInfo.Role)
	_ = d.Set("user_partition", roleInfo.UserPartition)

	d.SetId(roleInfo.Name)

	return nil

}
