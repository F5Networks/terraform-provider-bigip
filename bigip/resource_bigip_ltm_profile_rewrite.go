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

func resourceBigipLtmRewriteProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmRewriteProfileCreate,
		ReadContext:   resourceBigipLtmProfileRewriteRead,
		UpdateContext: resourceBigipLtmProfileRewriteUpdate,
		DeleteContext: resourceBigipLtmProfileRewriteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the rewrite profile.",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inherit defaults from parent profile.",
			},
			"bypass_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies a list of URIs to bypass inside a web page when the page is accessed using Portal Access.",
			},
			"cache_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "cache-img-css-js",
				Description:  "Specifies the type of client caching.",
				ValidateFunc: validation.StringInSlice([]string{"cache-css-js", "cache-all", "no-cache", "cache-img-css-js"}, false),
			},
			"ca_file": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies a CA against which to verify signed Java applets signatures.",
				ValidateFunc: validateF5Name,
			},
			"crl_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "none",
				Description: "Specifies a CRL against which to verify signed Java applets signature certificates.",
			},
			"signing_cert": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies a certificate to use for re-signing of signed Java applets after patching.",
				ValidateFunc: validateF5Name,
			},
			"signing_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies a private key for re-signing of signed Java applets after patching.",
				ValidateFunc: validateF5Name,
			},
			"split_tunneling": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies a private key for re-signing of signed Java applets after patching.",
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"signing_key_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				Description:  "Specifies a pass phrase to use for encrypting the private signing key.",
				ValidateFunc: validateF5Name,
			},
			"rewrite_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies a list of URIs to rewrite inside a web page when the page is accessed using Portal Access.",
			},
			"rewrite_mode": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specifies the type of rewrite operations.",
				ValidateFunc: validation.StringInSlice([]string{"portal", "uri-translation"}, false),
			},
			"request": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"insert_xfwd_for": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
							Description:  "Enable to add the X-Forwarded For (XFF) header, to specify the originating IP address of the client.",
						},
						"insert_xfwd_host": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
							Description:  "Enable to add the X-Forwarded Host header, to specify the originating host of the client.",
						},
						"insert_xfwd_protocol": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
							Description:  "Enable to add the X-Forwarded Proto header, to specify the originating protocol of the client.",
						},
						"rewrite_headers": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
							Description:  "Enable to rewrite headers in Request settings.",
						},
					},
				},
			},
			"response": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rewrite_content": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
							Description:  "Enable to rewrite links in content in the response.",
						},
						"rewrite_headers": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
							Description:  "Enable to rewrite headers in the response.",
						},
					},
				},
			},
			"cookie_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the cookie rewrite rule.",
						},
						"client_domain": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server_domain": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server_path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceBigipLtmRewriteProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	// partition := strings.Split(name, "/")[1]

	pss := &bigip.RewriteProfile{
		Name: name,
		// Partition: partition,
	}
	log.Printf("[INFO] Creating LTM rewrite profile config")
	config := getRewriteProfileConfig(d, pss)
	log.Printf("Config value:%+v", config)

	log.Printf("[INFO] Creating LTM rewrite profile")
	err := client.AddRewriteProfile(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Rewrite Profile %s %v :", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)

	return resourceBigipLtmProfileRewriteRead(ctx, d, meta)
}
func resourceBigipLtmProfileRewriteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Reading LTM rewrite profile config")
	profile, err := client.GetRewriteProfile(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if profile == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("[ERROR] LTM Rewrite Profile (%s) not found, removing from state", d.Id()))
	}
	log.Printf("[DEBUG] LTM rewrite profile:%+v", profile)
	setRewriteProfileData(d, profile)
	return nil
}

func resourceBigipLtmProfileRewriteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	profileConfig := &bigip.RewriteProfile{
		Name: name,
	}
	log.Println("[INFO] Updating LTM rewrite profile")
	rewriteProfileConfig := getRewriteProfileConfig(d, profileConfig)
	log.Printf("Config value:%+v", rewriteProfileConfig)
	err := client.ModifyRewriteProfile(name, rewriteProfileConfig)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error modifying LTM Rewrite Profile (%s): %s", name, err))
	}
	return resourceBigipLtmProfileRewriteRead(ctx, d, meta)
}

func resourceBigipLtmProfileRewriteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting LTM Rewrite Profile " + name)
	err := client.DeleteRewriteProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete LTM Rewrite Profile (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func setRewriteProfileData(d *schema.ResourceData, data *bigip.RewriteProfile) diag.Diagnostics {

	_ = d.Set("name", data.FullPath)
	_ = d.Set("rewrite_mode", data.Mode)
	if _, ok := d.GetOk("defaults_from"); ok {
		_ = d.Set("defaults_from", data.DefaultsFrom)
	}
	if _, ok := d.GetOk("ca_file"); ok {
		_ = d.Set("ca_file", data.CaFile)
	}
	if _, ok := d.GetOk("crl_file"); ok {
		_ = d.Set("crl_file", data.CrlFile)
	}
	if _, ok := d.GetOk("signing_cert"); ok {
		_ = d.Set("signing_cert", data.SigningCert)
	}
	if _, ok := d.GetOk("signing_key"); ok {
		_ = d.Set("signing_key", data.SigningKey)
	}
	if _, ok := d.GetOk("signing_key_password"); ok {
		_ = d.Set("signing_key_password", data.SigningKeyPass)
	}
	if _, ok := d.GetOk("cache_type"); ok {
		_ = d.Set("cache_type", data.CachingType)
	}
	if _, ok := d.GetOk("split_tunneling"); ok {
		_ = d.Set("split_tunneling", data.SplitTunnel)
	}
	if _, ok := d.GetOk("rewrite_list"); ok {
		_ = d.Set("rewrite_list", data.RewriteList)
	}
	if _, ok := d.GetOk("bypass_list"); ok {
		_ = d.Set("bypass_list", data.BypassList)
	}
	var reqList []interface{}
	req := make(map[string]interface{})
	if val, ok := d.GetOk("request"); ok {
		for _, item := range val.(*schema.Set).List() {
			if val, ok := item.(map[string]interface{})["insert_xfwd_for"].(string); ok && val != "" {
				req["insert_xfwd_for"] = val
			}
			if val, ok := item.(map[string]interface{})["insert_xfwd_host"].(string); ok && val != "" {
				req["insert_xfwd_host"] = val
			}
			if val, ok := item.(map[string]interface{})["insert_xfwd_protocol"].(string); ok && val != "" {
				req["insert_xfwd_protocol"] = val
			}
			if val, ok := item.(map[string]interface{})["rewrite_headers"].(string); ok && val != "" {
				req["rewrite_headers"] = val
			}
		}
	}
	reqList = append(reqList, req)
	_ = d.Set("request", reqList)
	var resList []interface{}
	res := make(map[string]interface{})
	if val, ok := d.GetOk("response"); ok {
		for _, item := range val.(*schema.Set).List() {
			if val, ok := item.(map[string]interface{})["rewrite_content"].(string); ok && val != "" {
				res["rewrite_content"] = val
			}
			if val, ok := item.(map[string]interface{})["rewrite_headers"].(string); ok && val != "" {
				res["rewrite_headers"] = val
			}
		}
		resList = append(resList, res)
		_ = d.Set("response", resList)
	}
	cookies := make([]interface{}, len(data.Cookies))
	for i, v := range data.Cookies {
		obj := make(map[string]interface{})
		if v.Name != "" {
			obj["rule_name"] = v.Name
		}
		if v.Client.Domain != "" {
			obj["client_domain"] = v.Client.Domain
		}
		if v.Client.Path != "" {
			obj["client_path"] = v.Client.Path
		}
		if v.Server.Domain != "" {
			obj["server_domain"] = v.Server.Domain
		}
		if v.Server.Path != "" {
			obj["server_path"] = v.Server.Path
		}
		cookies[i] = obj
	}
	err := d.Set("cookie_rules", cookies)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getRewriteProfileConfig(d *schema.ResourceData, config *bigip.RewriteProfile) *bigip.RewriteProfile {
	config.DefaultsFrom = d.Get("defaults_from").(string)
	// config.Partition = d.Get("partition").(string)
	config.Mode = d.Get("rewrite_mode").(string)
	config.CaFile = d.Get("ca_file").(string)
	config.CrlFile = d.Get("crl_file").(string)
	config.SigningCert = d.Get("signing_cert").(string)
	config.SigningKey = d.Get("signing_key").(string)
	config.SigningKeyPass = d.Get("signing_key_password").(string)
	config.CachingType = d.Get("cache_type").(string)
	config.SplitTunnel = d.Get("split_tunneling").(string)
	config.RewriteList = listToStringSlice(d.Get("rewrite_list").([]interface{}))
	config.BypassList = listToStringSlice(d.Get("bypass_list").([]interface{}))

	if val, ok := d.GetOk("request"); ok {
		var reqAttrs bigip.RewriteProfileRequestd
		for _, item := range val.(*schema.Set).List() {
			reqAttrs.XfwdFor = item.(map[string]interface{})["insert_xfwd_for"].(string)
			reqAttrs.XfwdHost = item.(map[string]interface{})["insert_xfwd_host"].(string)
			reqAttrs.XfwdProtocol = item.(map[string]interface{})["insert_xfwd_protocol"].(string)
			reqAttrs.RewriteHeaders = item.(map[string]interface{})["rewrite_headers"].(string)
		}
		config.Request = reqAttrs
	}

	if val, ok := d.GetOk("response"); ok {
		var resAttrs bigip.RewriteProfileResponsed
		for _, item := range val.(*schema.Set).List() {
			resAttrs.RewriteContent = item.(map[string]interface{})["rewrite_content"].(string)
			resAttrs.RewriteHeaders = item.(map[string]interface{})["rewrite_headers"].(string)
		}
		config.Response = resAttrs
	}

	if val, ok := d.GetOk("cookie_rules"); ok {
		var cookieRules []bigip.RewriteProfileCookieRules
		for _, item := range val.(*schema.Set).List() {
			cookieRule := bigip.RewriteProfileCookieRules{}
			log.Printf("[DEBUG] Value:%+v", item.(map[string]interface{})["rule_name"].(string))
			cookieRule.Name = item.(map[string]interface{})["rule_name"].(string)
			log.Printf("[DEBUG] Value:%+v", item.(map[string]interface{})["client_domain"].(string))
			cookieRule.Client.Domain = item.(map[string]interface{})["client_domain"].(string)
			log.Printf("[DEBUG] Value:%+v", item.(map[string]interface{})["client_path"].(string))
			cookieRule.Client.Path = item.(map[string]interface{})["client_path"].(string)
			log.Printf("[DEBUG] Value:%+v", item.(map[string]interface{})["server_domain"].(string))
			cookieRule.Server.Domain = item.(map[string]interface{})["server_domain"].(string)
			log.Printf("[DEBUG] Value:%+v", item.(map[string]interface{})["server_path"].(string))
			cookieRule.Server.Path = item.(map[string]interface{})["server_path"].(string)
			cookieRules = append(cookieRules, cookieRule)
		}
		config.Cookies = cookieRules
	}
	return config
}
