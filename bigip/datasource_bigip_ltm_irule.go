/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBigipLtmIrule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipLtmIruleRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the irule",
			},
			"partition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition of resource group",
			},

			"irule": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The iRule body",
				StateFunc: func(s interface{}) string {
					return strings.TrimSpace(s.(string))
				},
			},
		},
	}
}

func dataSourceBigipLtmIruleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	d.SetId("")
	name := fmt.Sprintf("/%s/%s", d.Get("partition").(string), d.Get("name").(string))

	irule, err := client.IRule(name)
	if err != nil {
		return fmt.Errorf("Error retrieving iRule %s: %v", name, err)
	}
	if irule == nil {
		log.Printf("[DEBUG] iRule (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", irule.FullPath)
	_ = d.Set("partition", irule.Partition)
	_ = d.Set("irule", irule.Rule)

	d.SetId(irule.FullPath)

	return nil

}
