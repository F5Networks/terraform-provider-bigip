/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2024 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipLtmRewriteProfileUriRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmRewriteProfileUriRuleCreate,
		ReadContext:   resourceBigipLtmProfileRewriteUriRuleRead,
		UpdateContext: resourceBigipLtmProfileRewriteUriRuleUpdate,
		DeleteContext: resourceBigipLtmProfileRewriteUriRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"profile_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the rewrite profile.",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the name of the uri rule.",
			},
			"rule_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "both",
				Description:  "Specifies the type of the uri rule.",
				ValidateFunc: validation.StringInSlice([]string{"request", "response", "both"}, false),
			},
			"client": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Host part of the uri, e.g. www.foo.com.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "/",
							Description: "Path part of the uri, when unspecified a trailing `/` is assumed.",
						},
						"scheme": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Scheme part of the uri, e.g. https, ftp etc.",
						},
						"port": {
							Type:        schema.TypeString,
							Default:     "none",
							Optional:    true,
							Description: "Port part of the uri, when not defined 'none' value is assumed.",
						},
					},
				},
			},
			"server": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Host part of the uri, e.g. www.foo.com",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "/",
							Description: "Path part of the uri, when unspecified a trailing `/` is assumed.",
						},
						"scheme": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Scheme part of the uri, e.g. https, ftp etc.",
						},
						"port": {
							Type:        schema.TypeString,
							Default:     "none",
							Optional:    true,
							Description: "Port part of the uri, when not defined 'none' value is assumed.",
						},
					},
				},
			},
		},
	}
}

func resourceBigipLtmRewriteProfileUriRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	profileName := d.Get("profile_name").(string)
	ruleName := d.Get("rule_name").(string)
	pss := &bigip.RewriteProfileUriRule{
		Name: ruleName,
	}

	log.Printf("[INFO] Creating LTM rewrite URI rule")
	config := getUriRulesConfig(d, pss)
	log.Printf("Config value:%+v", config)

	log.Printf("[INFO] Creating LTM rewrite URI rule")
	err := client.AddRewriteProfileUriRule(profileName, config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Rewrite URI Rule for profile %s %v :", profileName, err)
		return diag.FromErr(err)
	}
	d.SetId(ruleName)

	return resourceBigipLtmProfileRewriteUriRuleRead(ctx, d, meta)
}
func resourceBigipLtmProfileRewriteUriRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	ruleName := d.Id()
	profileName := d.Get("profile_name").(string)

	log.Printf("[INFO] Reading LTM rewrite URI rule: %s", ruleName)
	rules, err := client.GetRewriteProfileUriRule(profileName, ruleName)
	if err != nil {
		return diag.FromErr(err)
	}
	if rules == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("[ERROR] LTM Rewrite Profile URI rule (%s) not found, removing from state", d.Id()))
	}
	log.Printf("[DEBUG] LTM rewrite profile URI rule:%+v", rules)

	setUriRulesData(d, rules)
	return nil
}

func resourceBigipLtmProfileRewriteUriRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	ruleName := d.Id()
	profileName := d.Get("profile_name").(string)
	uriConfig := &bigip.RewriteProfileUriRule{
		Name: ruleName,
	}
	log.Println("[INFO] Updating LTM rewrite URI rule")
	uriRules := getUriRulesConfig(d, uriConfig)
	log.Printf("Config value:%+v", uriRules)
	err := client.ModifyRewriteProfileUriRule(profileName, ruleName, uriRules)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error modifying LTM Rewrite URI rule (%s): %s", ruleName, err))
	}
	return resourceBigipLtmProfileRewriteRead(ctx, d, meta)
}

func resourceBigipLtmProfileRewriteUriRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	ruleName := d.Id()
	profileName := d.Get("profile_name").(string)
	log.Println("[INFO] Deleting LTM Rewrite Profile URI rule " + ruleName)
	err := client.DeleteRewriteProfileUriRule(profileName, ruleName)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete LTM Rewrite Profile URI rule (%s) (%v) ", ruleName, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
func setUriRulesData(d *schema.ResourceData, data *bigip.RewriteProfileUriRule) {
	_ = d.Set("rule_name", data.Name)
	_ = d.Set("rule_type", data.Type)
	var clList []interface{}
	cl := make(map[string]interface{})
	cl["host"] = data.Client.Host
	cl["path"] = data.Client.Path
	cl["scheme"] = data.Client.Scheme
	cl["port"] = data.Client.Port
	clList = append(clList, cl)
	_ = d.Set("client", clList)

	var srvList []interface{}
	srv := make(map[string]interface{})
	srv["host"] = data.Server.Host
	srv["path"] = data.Server.Path
	srv["scheme"] = data.Server.Scheme
	srv["port"] = data.Server.Port
	srvList = append(srvList, srv)
	_ = d.Set("server", srvList)
}

func getUriRulesConfig(d *schema.ResourceData, config *bigip.RewriteProfileUriRule) *bigip.RewriteProfileUriRule {
	config.Type = d.Get("rule_type").(string)
	client := d.Get("client")
	server := d.Get("server")

	for _, item := range client.(*schema.Set).List() {
		config.Client.Host = item.(map[string]interface{})["host"].(string)
		config.Client.Path = item.(map[string]interface{})["path"].(string)
		config.Client.Scheme = item.(map[string]interface{})["scheme"].(string)
		config.Client.Port = item.(map[string]interface{})["port"].(string)
	}

	for _, item := range server.(*schema.Set).List() {
		config.Server.Host = item.(map[string]interface{})["host"].(string)
		config.Server.Path = item.(map[string]interface{})["path"].(string)
		config.Server.Scheme = item.(map[string]interface{})["scheme"].(string)
		config.Server.Port = item.(map[string]interface{})["port"].(string)
	}
	return config
}
