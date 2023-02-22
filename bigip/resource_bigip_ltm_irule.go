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
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipLtmIRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmIRuleCreate,
		ReadContext:   resourceBigipLtmIRuleRead,
		UpdateContext: resourceBigipLtmIRuleUpdate,
		DeleteContext: resourceBigipLtmIRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the iRule",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},

			"irule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The iRule body",
				StateFunc: func(s interface{}) string {
					return strings.TrimSpace(s.(string))
				},
			},
		},
	}
}

func resourceBigipLtmIRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating iRule %s", name)

	err := client.CreateIRule(name, d.Get("irule").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating iRule %s: %v", name, err))
	}

	d.SetId(name)

	return resourceBigipLtmIRuleRead(ctx, d, meta)
}

func resourceBigipLtmIRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[INFO] Retrieving iRule %s", name)

	irule, err := client.IRule(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving iRule %s: %v", name, err))
	}
	if irule == nil {
		log.Printf("[DEBUG] iRule (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", irule.FullPath)
	_ = d.Set("irule", irule.Rule)

	return nil
}

func resourceBigipLtmIRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	r := &bigip.IRule{
		FullPath: name,
		Rule:     d.Get("irule").(string),
	}

	err := client.ModifyIRule(name, r)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error modifying iRule %s: %v", name, err))
	}
	return resourceBigipLtmIRuleRead(ctx, d, meta)
}

func resourceBigipLtmIRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeleteIRule(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting iRule %s: %v", name, err))
	}
	d.SetId("")
	return nil
}
