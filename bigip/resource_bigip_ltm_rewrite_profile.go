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
	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
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
				Required:     true,
				Description:  "Specifies the type of Client caching.",
				ValidateFunc: validation.StringInSlice([]string{"cache-css-js", "cache-all", "no-cache", "cache-img-css-js"}, false),
			},
			"ca_file": {
				Type:         schema.TypeString,
				Optional:     true,
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
				Description:  "Specifies a certificate to use for re-signing of signed Java applets after patching.",
				ValidateFunc: validateF5Name,
			},
			"signing_key": {
				Type:         schema.TypeString,
				Optional:     true,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"insert_xfwd_for": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
						},
						"insert_xfwd_host": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
						},
						"insert_xfwd_protocol": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
						},
						"rewrite_headers": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
						},
					},
				},
			},
			"response": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rewrite_content": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
						},
						"rewrite_headers": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
							Optional:     true,
						},
					},
				},
			},
			//"cookie_rules": {
			//	Type:     schema.TypeList,
			//	Optional: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"name": {
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//			"client": {
			//				Type:     schema.TypeSet,
			//				Required: true,
			//				Elem: &schema.Resource{
			//					Schema: map[string]*schema.Schema{
			//						"domain": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//						"path": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//					},
			//				},
			//			},
			//			"server": {
			//				Type:     schema.TypeSet,
			//				Required: true,
			//				Elem: &schema.Resource{
			//					Schema: map[string]*schema.Schema{
			//						"domain": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//						"path": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//					},
			//				},
			//			},
			//		},
			//	},
			//},
			//"uri_rules": {
			//	Type:     schema.TypeList,
			//	Optional: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"name": {
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//			"client": {
			//				Type:     schema.TypeSet,
			//				Required: true,
			//				Elem: &schema.Resource{
			//					Schema: map[string]*schema.Schema{
			//						"host": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//						"path": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//						"scheme": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//						"port": {
			//							Type:     schema.TypeString,
			//							Default:  "none",
			//							Optional: true,
			//						},
			//					},
			//				},
			//			},
			//			"server": {
			//				Type:     schema.TypeSet,
			//				Required: true,
			//				Elem: &schema.Resource{
			//					Schema: map[string]*schema.Schema{
			//						"host": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//						"path": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//						"scheme": {
			//							Type:     schema.TypeString,
			//							Required: true,
			//						},
			//						"port": {
			//							Type:     schema.TypeString,
			//							Default:  "none",
			//							Optional: true,
			//						},
			//					},
			//				},
			//			},
			//		},
			//	},
			//},
		},
	}
}

func resourceBigipLtmRewriteProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	partition := strings.Split(name, "/")[1]

	pss := &bigip.RewriteProfile{
		Name:      name,
		Partition: partition,
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
	//if _, ok := d.GetOk("uri_rules"); ok {
	//	ruleCount := d.Get("uri_rules.#").(int)
	//	for i := 0; i < ruleCount; i++ {
	//		var ruleConfig bigip.RewriteProfileUriRule
	//		prefix := fmt.Sprintf("uri_rules.%d", i)
	//		ruleConfig.Name = d.Get(prefix + ".name").(string)
	//		ruleConfig.Type = d.Get(prefix + ".type").(string)
	//		if val, ok := d.GetOk(prefix + ".client"); ok {
	//			var rc bigip.RewriteProfileUrlClSrv
	//			for _, item := range val.(*schema.Set).List() {
	//				log.Printf("Value:%+v", item.(map[string]interface{})["host"].(string))
	//				rc.Host = item.(map[string]interface{})["host"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["path"].(string))
	//				rc.Path = item.(map[string]interface{})["path"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["port"].(string))
	//				rc.Port = item.(map[string]interface{})["port"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["scheme"].(string))
	//				rc.Scheme = item.(map[string]interface{})["scheme"].(string)
	//			}
	//			ruleConfig.Client = rc
	//		}
	//		if val, ok := d.GetOk(prefix + ".server"); ok {
	//			var sc bigip.RewriteProfileUrlClSrv
	//			for _, item := range val.(*schema.Set).List() {
	//				log.Printf("Value:%+v", item.(map[string]interface{})["host"].(string))
	//				sc.Host = item.(map[string]interface{})["host"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["path"].(string))
	//				sc.Path = item.(map[string]interface{})["path"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["port"].(string))
	//				sc.Port = item.(map[string]interface{})["port"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["scheme"].(string))
	//				sc.Scheme = item.(map[string]interface{})["scheme"].(string)
	//			}
	//			ruleConfig.Server = sc
	//		}
	//		err := client.AddRewriteProfileUriRule(name, &ruleConfig)
	//		if err != nil {
	//			log.Printf("[ERROR] Unable to Create Url Rewrite Profile Rule %s %v :", name, err)
	//			return diag.FromErr(err)
	//		}
	//	}
	//}

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

	//uriRules, err := client.GetRewriteProfileUriRules(name)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//if uriRules == nil {
	//	log.Printf("No URI rules found associated with this LTM Rewrite Profile")
	//}
	//if uriRules != nil {
	//	rules := make([]interface{}, len(uriRules.Uri))
	//	for i, v := range uriRules.Uri {
	//		obj := make(map[string]interface{})
	//		if v.Name != "" {
	//			obj["name"] = v.Name
	//		}
	//		cl := map[string]interface{}{
	//			"host":   v.Client.Host,
	//			"path":   v.Client.Path,
	//			"scheme": v.Client.Scheme,
	//			"port":   v.Client.Port,
	//		}
	//		obj["client"] = cl
	//
	//		sr := map[string]interface{}{
	//			"host":   v.Client.Host,
	//			"path":   v.Client.Path,
	//			"scheme": v.Client.Scheme,
	//			"port":   v.Client.Port,
	//		}
	//		obj["server"] = sr
	//
	//		rules[i] = obj
	//	}
	//	err := d.Set("uri_rules", rules)
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//}
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

	//if _, ok := d.GetOk("uri_rules"); ok {
	//	ruleCount := d.Get("uri_rules.#").(int)
	//	for i := 0; i < ruleCount; i++ {
	//		var ruleConfig bigip.RewriteProfileUriRule
	//		prefix := fmt.Sprintf("uri_rules.%d", i)
	//		ruleConfig.Name = d.Get(prefix + ".name").(string)
	//		ruleConfig.Type = d.Get(prefix + ".type").(string)
	//		if val, ok := d.GetOk(prefix + ".client"); ok {
	//			var rc bigip.RewriteProfileUrlClSrv
	//			for _, item := range val.(*schema.Set).List() {
	//				log.Printf("Value:%+v", item.(map[string]interface{})["host"].(string))
	//				rc.Host = item.(map[string]interface{})["host"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["path"].(string))
	//				rc.Path = item.(map[string]interface{})["path"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["port"].(string))
	//				rc.Port = item.(map[string]interface{})["port"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["scheme"].(string))
	//				rc.Scheme = item.(map[string]interface{})["scheme"].(string)
	//			}
	//			ruleConfig.Client = rc
	//		}
	//		if val, ok := d.GetOk(prefix + ".server"); ok {
	//			var sc bigip.RewriteProfileUrlClSrv
	//			for _, item := range val.(*schema.Set).List() {
	//				log.Printf("Value:%+v", item.(map[string]interface{})["host"].(string))
	//				sc.Host = item.(map[string]interface{})["host"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["path"].(string))
	//				sc.Path = item.(map[string]interface{})["path"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["port"].(string))
	//				sc.Port = item.(map[string]interface{})["port"].(string)
	//				log.Printf("Value:%+v", item.(map[string]interface{})["scheme"].(string))
	//				sc.Scheme = item.(map[string]interface{})["scheme"].(string)
	//			}
	//			ruleConfig.Server = sc
	//		}
	//		err := client.ModifyRewriteProfileUriRule(name, ruleConfig.Name, &ruleConfig)
	//		if err != nil {
	//			log.Printf("[ERROR] Unable to Modify LTM Rewrite Profile Uri Rule %s %v :", name, err)
	//			return diag.FromErr(err)
	//		}
	//	}
	//}
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
	_ = d.Set("defaults_from", data.DefaultsFrom)
	_ = d.Set("rewrite_mode", data.Mode)
	_ = d.Set("ca_file", data.CaFile)
	_ = d.Set("crl_file", data.CrlFile)
	_ = d.Set("signing_cert", data.SigningCert)
	_ = d.Set("signing_key", data.SigningKey)
	_ = d.Set("signing_key_password", data.SigningKeyPass)
	_ = d.Set("cache_type", data.CachingType)
	_ = d.Set("split_tunneling", data.SplitTunnel)
	_ = d.Set("rewrite_list", data.RewriteList)
	_ = d.Set("bypass_list", data.BypassList)

	var reqList []interface{}
	req := make(map[string]interface{})
	req["insert_xfwd_for"] = data.Request.XfwdFor
	req["insert_xfwd_host"] = data.Request.XfwdHost
	req["insert_xfwd_protocol"] = data.Request.XfwdProtocol
	req["rewrite_headers"] = data.Request.RewriteHeaders
	reqList = append(reqList, req)
	_ = d.Set("request", reqList)
	var resList []interface{}
	res := make(map[string]interface{})
	res["rewrite_content"] = data.Response.RewriteContent
	res["rewrite_headers"] = data.Response.RewriteHeaders
	resList = append(resList, res)
	_ = d.Set("response", resList)
	//cookies := make([]interface{}, len(data.Cookies))
	//for i, v := range data.Cookies {
	//	obj := make(map[string]interface{})
	//	if v.Name != "" {
	//		obj["name"] = v.Name
	//	}
	//	cl := map[string]interface{}{
	//		"domain": v.Client.Domain,
	//		"path":   v.Client.Path,
	//	}
	//	obj["client"] = cl
	//	sr := map[string]interface{}{
	//		"domain": v.Server.Domain,
	//		"path":   v.Server.Path,
	//	}
	//	obj["server"] = sr
	//	cookies[i] = obj
	//}
	//err := d.Set("cookie_rules", cookies)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	return nil
}

func getRewriteProfileConfig(d *schema.ResourceData, config *bigip.RewriteProfile) *bigip.RewriteProfile {
	config.DefaultsFrom = d.Get("defaults_from").(string)
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
			log.Printf("Value:%+v", item.(map[string]interface{})["insert_xfwd_for"].(string))
			reqAttrs.XfwdFor = item.(map[string]interface{})["insert_xfwd_for"].(string)
			log.Printf("Value:%+v", item.(map[string]interface{})["insert_xfwd_host"].(string))
			reqAttrs.XfwdHost = item.(map[string]interface{})["insert_xfwd_host"].(string)
			log.Printf("Value:%+v", item.(map[string]interface{})["insert_xfwd_protocol"].(string))
			reqAttrs.XfwdProtocol = item.(map[string]interface{})["insert_xfwd_protocol"].(string)
			log.Printf("Value:%+v", item.(map[string]interface{})["rewrite_headers"].(string))
			reqAttrs.RewriteHeaders = item.(map[string]interface{})["rewrite_headers"].(string)
		}
		config.Request = reqAttrs
	}

	if val, ok := d.GetOk("response"); ok {
		var resAttrs bigip.RewriteProfileResponsed
		for _, item := range val.(*schema.Set).List() {
			log.Printf("Value:%+v", item.(map[string]interface{})["rewrite_content"].(string))
			resAttrs.RewriteContent = item.(map[string]interface{})["rewrite_content"].(string)
			log.Printf("Value:%+v", item.(map[string]interface{})["rewrite_headers"].(string))
			resAttrs.RewriteHeaders = item.(map[string]interface{})["rewrite_headers"].(string)
		}
		config.Response = resAttrs
	}

	//cookieCount := d.Get("cookie_rules.#").(int)
	//log.Printf("[INFO] Cookie Count:%+v", cookieCount)
	//config.Cookies = make([]bigip.RewriteProfileCookieRules, 0, cookieCount)
	//for i := 0; i < cookieCount; i++ {
	//	var cookieRules bigip.RewriteProfileCookieRules
	//	prefix := fmt.Sprintf("cookie_rules.%d", i)
	//	cookieRules.Name = d.Get(prefix + ".name").(string)
	//	if val, ok := d.GetOk(prefix + ".client"); ok {
	//		var cc bigip.RewriteProfileCookieClSrv
	//		for _, item := range val.(*schema.Set).List() {
	//			log.Printf("Value:%+v", item.(map[string]interface{})["domain"].(string))
	//			cc.Domain = item.(map[string]interface{})["domain"].(string)
	//			log.Printf("Value:%+v", item.(map[string]interface{})["path"].(string))
	//			cc.Path = item.(map[string]interface{})["path"].(string)
	//		}
	//		cookieRules.Client = cc
	//	}
	//	if val, ok := d.GetOk(prefix + ".server"); ok {
	//		var sc bigip.RewriteProfileCookieClSrv
	//		for _, item := range val.(*schema.Set).List() {
	//			log.Printf("Value:%+v", item.(map[string]interface{})["domain"].(string))
	//			sc.Domain = item.(map[string]interface{})["domain"].(string)
	//			log.Printf("Value:%+v", item.(map[string]interface{})["path"].(string))
	//			sc.Path = item.(map[string]interface{})["path"].(string)
	//		}
	//		cookieRules.Server = sc
	//	}
	//	config.Cookies = append(config.Cookies, cookieRules)
	//}
	return config
}
