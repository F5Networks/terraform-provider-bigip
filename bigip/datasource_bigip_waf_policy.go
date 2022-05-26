/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBigipWafPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipWafPolicyRead,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the policy",
			},
			"policy_json": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Exported JSON policy",
			},
		},
	}
}

func dataSourceBigipWafPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	policyID := d.Get("policy_id").(string)

	log.Printf("[DEBUG] Reading AWAF Policy with ID: %+v", policyID)

	wafpolicy, err := client.GetWafPolicy(policyID)
	if err != nil {
		return fmt.Errorf("error retrieving waf policy %+v: %v", wafpolicy, err)
	}

	policyJson, err := client.ExportPolicy(policyID)
	if err != nil {
		return fmt.Errorf("error Exporting waf policy ID %v with : %+v", policyID, err)
	}

	log.Printf("[DEBUG] Policy Json : %+v", policyJson.Policy)

	plJson, err := json.Marshal(policyJson.Policy)
	if err != nil {
		return err
	}

	_ = d.Set("policy_id", wafpolicy.ID)
	_ = d.Set("policy_json", string(plJson))

	d.SetId(wafpolicy.ID)
	return nil
}
