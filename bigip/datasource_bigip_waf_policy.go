/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipWafPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipWafPolicyRead,
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

func dataSourceBigipWafPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	d.SetId("")

	policyID := d.Get("policy_id").(string)

	log.Printf("[DEBUG] Reading AWAF Policy with ID: %+v", policyID)

	wafpolicy, err := client.GetWafPolicy(policyID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving waf policy %+v: %v", wafpolicy, err))
	}

	policyJson, err := client.ExportPolicy(policyID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error Exporting waf policy ID %v with : %+v", policyID, err))
	}

	log.Printf("[DEBUG] Policy Json : %+v", policyJson.Policy)

	plJson, err := json.Marshal(policyJson.Policy)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("policy_id", wafpolicy.ID)
	_ = d.Set("policy_json", string(plJson))

	d.SetId(wafpolicy.ID)
	return nil
}
