/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBigipLtmPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipLtmPoolRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the certificate",
			},
			"partition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "partition of resource group",
			},
			"full_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full path, if found.",
			},
		},
	}
}
func dataSourceBigipLtmPoolRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*bigip.BigIP)
	d.SetId("")
	name := fmt.Sprintf("/%s/%s", d.Get("partition").(string), d.Get("name").(string))

	log.Println("[INFO] Reading Pool : " + name)
	pool, err := client.GetPool(name)
	if err != nil {
		return fmt.Errorf("Error retrieving pool %s: %v", name, err)
	}
	if pool == nil {
		return fmt.Errorf("Pool (%s) not found", name)
	}

	d.Set("name", pool.Name)
	d.Set("partition", pool.Partition)
	d.Set("full_path", pool.FullPath)
	d.SetId(pool.FullPath)

	return nil

}
