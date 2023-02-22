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

func resourceBigipIpsecProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipIpsecProfileCreate,
		ReadContext:   resourceBigipIpsecProfileRead,
		UpdateContext: resourceBigipIpsecProfileUpdate,
		DeleteContext: resourceBigipIpsecProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Displays the name of the IPsec interface tunnel profile",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies descriptive text that identifies the IPsec interface tunnel profile",
			},
			"parent_profile": {
				Type:        schema.TypeString,
				Default:     "/Common/ipsec",
				Optional:    true,
				Description: "Specifies the profile from which this profile inherits settings. The default is the system-supplied `/Common/ipsec` profile",
			},
			"traffic_selector": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateF5Name,
				Description:  "Specifies the traffic selector for the IPsec interface tunnel to which the profile is applied",
			},
		},
	}
}

func resourceBigipIpsecProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating IPSec profile " + name)

	pss := &bigip.IPSecProfile{
		Name: name,
	}
	selectorConfig := getIPSecProfileConfig(d, pss)

	err := client.CreateIPSecProfile(selectorConfig)
	if err != nil {
		log.Printf("[ERROR] Unable to Create IPsec profile (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipIpsecProfileRead(ctx, d, meta)
}

func resourceBigipIpsecProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Reading IPsec profile :%+v", name)
	ts, err := client.GetIPSecProfile(name)
	log.Printf("IPsec Profile:%+v", ts)
	if err != nil {
		return diag.FromErr(err)
	}
	if ts == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("[ERROR] IPsec profile (%s) not found, removing from state", d.Id()))
	}
	if err := d.Set("parent_profile", ts.DefaultsFrom); err != nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Error saving IPsec parent profile (%s): %s", d.Id(), err))
	}
	if err := d.Set("traffic_selector", ts.TrafficSelector); err != nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Error saving IPsec profile (%s): %s", d.Id(), err))
	}
	_ = d.Set("description", ts.Description)
	_ = d.Set("name", name)
	return nil
}

func resourceBigipIpsecProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating IPsec Profile:%+v ", name)
	pss := &bigip.IPSecProfile{
		Name: name,
	}
	config := getIPSecProfileConfig(d, pss)

	err := client.ModifyIPSecProfile(name, config)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify IPsec Profile   (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	return resourceBigipIpsecProfileRead(ctx, d, meta)
}
func resourceBigipIpsecProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Deleting IPsec Profile :%+v ", name)
	err := client.DeleteIPSecProfile(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Unable to Delete IPsec Profile (%s) (%v) ", name, err))
	}
	d.SetId("")
	return nil
}

func getIPSecProfileConfig(d *schema.ResourceData, config *bigip.IPSecProfile) *bigip.IPSecProfile {
	config.DefaultsFrom = d.Get("parent_profile").(string)
	config.Description = d.Get("description").(string)
	config.TrafficSelector = d.Get("traffic_selector").(string)
	return config
}
