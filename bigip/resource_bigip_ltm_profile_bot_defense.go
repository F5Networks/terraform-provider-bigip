/*
Copyright 2024 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipLtmProfileBotDefense() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmProfileBotDefenseCreate,
		ReadContext:   resourceBigipLtmProfileBotDefenseRead,
		UpdateContext: resourceBigipLtmProfileBotDefenseUpdate,
		DeleteContext: resourceBigipLtmProfileBotDefenseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the Bot Defense profile",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "/Common/bot-defense",
				Description:  "Specifies the profile from which this profile inherits settings. The default is the system-supplied `request-log` profile",
				ValidateFunc: validateF5Name,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User defined description for Bot Defense profile",
			},
			"template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"relaxed", "balanced", "strict"}, false),
				Description: "Profile templates specify Mitigation and Verification Settings default values. possible ptions `balanced`,`relaxed` and `strict`",
			},
			"enforcement_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"transparent",
					"blocking"}, false),
				Description: "Select the enforcement mode, possible values are `transparent` and `blocking`.",
			},
		},
	}
}

func resourceBigipLtmProfileBotDefenseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Printf("[INFO] Creating Bot Defense Profile:%+v ", name)
	pss := &bigip.BotDefenseProfile{
		Name: name,
	}
	config := getProfileBotDefenseConfig(d, pss)
	log.Printf("[DEBUG] Bot Defense Profile config :%+v ", config)
	err := client.AddBotDefenseProfile(config)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipLtmProfileBotDefenseRead(ctx, d, meta)
}

func resourceBigipLtmProfileBotDefenseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading Bot Defense Profile:%+v ", client)
	name := d.Id()
	log.Printf("[INFO] Reading Bot Defense Profile:%+v ", name)
	botProfile, err := client.GetBotDefenseProfile(name)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Bot Defense Profile config :%+v ", botProfile)
	d.Set("name", botProfile.FullPath)
	d.Set("defaults_from", botProfile.DefaultsFrom)
	d.Set("description", botProfile.Description)
	d.Set("template", botProfile.Template)
	d.Set("enforcement_mode", botProfile.EnforcementMode)
	return nil
}

func resourceBigipLtmProfileBotDefenseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating Bot Defense Profile:%+v ", name)
	pss := &bigip.BotDefenseProfile{
		Name: name,
	}
	config := getProfileBotDefenseConfig(d, pss)

	err := client.ModifyBotDefenseProfile(name, config)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceBigipLtmProfileBotDefenseRead(ctx, d, meta)
}

func resourceBigipLtmProfileBotDefenseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Bot Defense Profile " + name)
	err := client.DeleteBotDefenseProfile(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func getProfileBotDefenseConfig(d *schema.ResourceData, config *bigip.BotDefenseProfile) *bigip.BotDefenseProfile {
	config.Name = d.Get("name").(string)
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.Description = d.Get("description").(string)
	config.Template = d.Get("template").(string)
	config.EnforcementMode = d.Get("enforcement_mode").(string)
	log.Printf("[INFO][getProfileBotDefenseConfig] config:%+v ", config)
	return config
}
