// Copyright 2023 F5 Networks Inc.
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.

package bigip

import (
	"context"
	"fmt"
	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceBigipLtmCipherGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmCipherGroupCreate,
		ReadContext:   resourceBigipLtmCipherGroupRead,
		UpdateContext: resourceBigipLtmCipherGroupUpdate,
		DeleteContext: resourceBigipLtmCipherGroupDelete,
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
			"ordering": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies one or more Cipher Suites used.Note: For SM2, type the following cipher suite string: ECC-SM4-SM3.",
			},
			"allow": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the DH Groups Elliptic Curve Diffie-Hellman key exchange algorithms, separated by colons (:).Note: You can also type a special keyword, DEFAULT, which represents the recommended set of named groups",
			},
			"require": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the DH Groups Elliptic Curve Diffie-Hellman key exchange algorithms, separated by colons (:).Note: You can also type a special keyword, DEFAULT, which represents the recommended set of named groups",
			},
		},
	}
}

func resourceBigipLtmCipherGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	log.Printf("[INFO] Creating Cipher rule:%+v", name)

	cipherGrouptmp := &bigip.CipherGroupReq{}
	cipherGrouptmp.Name = name
	cipherGroup, err := getCipherGroupConfig(d, cipherGrouptmp)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading input config failed(%s): %s", name, err))
	}
	log.Printf("[INFO] cipherGroup config :%+v", cipherGroup)
	err = client.AddLtmCipherGroup(cipherGroup)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating cipher rule (%s): %s", name, err))
	}
	d.SetId(name)
	return resourceBigipLtmCipherGroupRead(ctx, d, meta)
}

func resourceBigipLtmCipherGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Fetching Cipher group :%+v", name)

	cipherRule, err := client.GetLtmCipherGroup(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve cipher rule %s  %v :", name, err)
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Cipher rule response :%+v", cipherRule)
	return nil
}

func resourceBigipLtmCipherGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	cipherGrouptmp := &bigip.CipherGroupReq{}
	cipherGrouptmp.Name = name
	cipherGroupconfig, err := getCipherGroupConfig(d, cipherGrouptmp)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading input config failed(%s): %s", name, err))
	}
	if err := client.ModifyLtmCipherGroup(name, cipherGroupconfig); err != nil {
		return diag.FromErr(fmt.Errorf("error modifying cipher group %s: %v", name, err))
	}

	return resourceBigipLtmCipherGroupRead(ctx, d, meta)
}

func resourceBigipLtmCipherGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Printf("[INFO] Deleting cipher group :%+v", name)
	err := client.DeleteLtmCipherGroup(name)

	if err != nil {
		log.Printf("[ERROR] Unable to Delete cipher rule %s  %v : ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getCipherGroupConfig(d *schema.ResourceData, cipherGroup *bigip.CipherGroupReq) (*bigip.CipherGroupReq, error) {
	cipherGroup.Ordering = d.Get("ordering").(string)
	if p, ok := d.GetOk("allow"); ok {
		for _, r := range p.(*schema.Set).List() {
			cipherGroup.Allow = append(cipherGroup.Allow, r.(string))
		}
	}
	if p, ok := d.GetOk("require"); ok {
		for _, r := range p.(*schema.Set).List() {
			cipherGroup.Require = append(cipherGroup.Require, r.(string))
		}
	}
	return cipherGroup, nil
}
