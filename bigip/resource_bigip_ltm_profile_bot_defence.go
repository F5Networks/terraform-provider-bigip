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

func resourceBigipLtmProfileBotDefence() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmProfileBotDefenceCreate,
		ReadContext:   resourceBigipLtmProfileBotDefenceRead,
		UpdateContext: resourceBigipLtmProfileBotDefenceUpdate,
		DeleteContext: resourceBigipLtmProfileBotDefenceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the Bot Defence profile",
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
				Description: "User defined description for Bot Defence profile",
			},
			"template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"relaxed",
					"enabled"}, false),
				Description: "Enables or disables Bot Defence. The default is `disabled`",
			},
			"enforcement_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"transparent",
					"blocking"}, false),
				Description: "Specifies the protocol to be used for high-speed logging of requests. The default is `mds-udp`",
			},
		},
	}
}

func resourceBigipLtmProfileBotDefenceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Printf("[INFO] Creating Bot Defence Profile:%+v ", name)
	pss := &bigip.BotDefenseProfile{
		Name: name,
	}
	config := getProfileBotDefenceConfig(d, pss)
	log.Printf("[DEBUG] Bot Defence Profile config :%+v ", config)
	err := client.AddBotDefenseProfile(config)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipLtmProfileBotDefenceRead(ctx, d, meta)
}

func resourceBigipLtmProfileBotDefenceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading Bot Defence Profile:%+v ", client)
	name := d.Id()
	log.Printf("[INFO] Reading Bot Defence Profile:%+v ", name)
	botProfile, err := client.GetBotDefenseProfile(name)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Bot Defence Profile config :%+v ", botProfile)
	d.Set("name", botProfile.FullPath)
	d.Set("defaults_from", botProfile.DefaultsFrom)
	d.Set("description", botProfile.Description)
	d.Set("template", botProfile.Template)
	d.Set("enforcement_mode", botProfile.EnforcementMode)
	return nil
}

func resourceBigipLtmProfileBotDefenceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating Bot Defence Profile:%+v ", name)
	pss := &bigip.BotDefenseProfile{
		Name: name,
	}
	config := getProfileBotDefenceConfig(d, pss)

	err := client.ModifyBotDefenseProfile(name, config)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceBigipLtmProfileBotDefenceRead(ctx, d, meta)
}

func resourceBigipLtmProfileBotDefenceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Bot Defence Profile " + name)
	err := client.DeleteBotDefenseProfile(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func getProfileBotDefenceConfig(d *schema.ResourceData, config *bigip.BotDefenseProfile) *bigip.BotDefenseProfile {
	config.Name = d.Get("name").(string)
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.Description = d.Get("description").(string)
	config.Template = d.Get("template").(string)
	config.EnforcementMode = d.Get("enforcement_mode").(string)
	log.Printf("[INFO][getProfileBotDefenceConfig] config:%+v ", config)
	return config
}
