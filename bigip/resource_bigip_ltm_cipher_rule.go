// Copyright 2023 F5 Networks Inc.
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.

package bigip

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipLtmCipherRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmCipherRuleCreate,
		ReadContext:   resourceBigipLtmCipherRuleRead,
		UpdateContext: resourceBigipLtmCipherRuleUpdate,
		DeleteContext: resourceBigipLtmCipherRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the cipher rule,name should be in pattern ``partition` + `cipher rule name``",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies descriptive text that identifies the cipher rule",
			},
			"cipher": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies one or more Cipher Suites used.Note: For SM2, type the following cipher suite string: ECC-SM4-SM3.",
			},
			"dh_groups": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the DH Groups Elliptic Curve Diffie-Hellman key exchange algorithms, separated by colons (:).Note: You can also type a special keyword, DEFAULT, which represents the recommended set of named groups",
			},
			"signature_algorithms": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the Signature Algorithms, separated by colons (:), that you want to include in the cipher rule. You can also type a special keyword, DEFAULT, which represents the recommended set of signature algorithms",
			},
		},
	}
}

func resourceBigipLtmCipherRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	log.Printf("[INFO] Creating Cipher rule:%+v", name)

	cipherRuletmp := &bigip.CipherRuleReq{}
	cipherRuletmp.Name = name

	cipherRule := getCipherRuleConfig(d, cipherRuletmp)

	log.Printf("[INFO] cipherRule config :%+v", cipherRule)
	err := client.AddLtmCipherRule(cipherRule)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating cipher rule (%s): %s", name, err))
	}
	d.SetId(name)
	if !client.Teem {
		id := uuid.New()
		uniqueID := id.String()
		assetInfo := f5teem.AssetInfo{
			Name:    "Terraform-provider-bigip",
			Version: client.UserAgent,
			Id:      uniqueID,
		}
		apiKey := os.Getenv("TEEM_API_KEY")
		teemDevice := f5teem.AnonymousClient(assetInfo, apiKey)
		f := map[string]interface{}{
			"Terraform Version": client.UserAgent,
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err = teemDevice.Report(f, "bigip_ltm_cipher_rule", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipLtmCipherRuleRead(ctx, d, meta)
}

func resourceBigipLtmCipherRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Fetching Cipher rule :%+v", name)
	cipherRule, err := client.GetLtmCipherRule(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve cipher rule %s  %v :", name, err)
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Cipher rule response :%+v", cipherRule)
	_ = d.Set("name", cipherRule.FullPath)
	_ = d.Set("cipher", cipherRule.Cipher)
	_ = d.Set("dh_groups", cipherRule.DhGroups)
	_ = d.Set("signature_algorithms", cipherRule.SignatureAlgorithms)
	return nil
}

func resourceBigipLtmCipherRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	cipherRuletmp := &bigip.CipherRuleReq{}
	cipherRuletmp.Name = name
	cipheRuleconfig := getCipherRuleConfig(d, cipherRuletmp)
	if err := client.ModifyLtmCipherRule(name, cipheRuleconfig); err != nil {
		return diag.FromErr(fmt.Errorf("error modifying cipher rule %s: %v", name, err))
	}
	return resourceBigipLtmCipherRuleRead(ctx, d, meta)
}

func resourceBigipLtmCipherRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Deleting cipher rule :%+v", name)
	err := client.DeleteLtmCipherRule(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete cipher rule %s  %v : ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getCipherRuleConfig(d *schema.ResourceData, cipherRule *bigip.CipherRuleReq) *bigip.CipherRuleReq {
	cipherRule.Cipher = d.Get("cipher").(string)
	cipherRule.DhGroups = d.Get("dh_groups").(string)
	cipherRule.SignatureAlgorithms = d.Get("signature_algorithms").(string)
	cipherRule.Description = d.Get("description").(string)
	return cipherRule
}
