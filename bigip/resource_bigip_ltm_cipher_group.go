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
				Description:  "Name of the cipher group,name should be in pattern ``partition` + `cipher group name``",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies descriptive text that identifies the cipher rule",
			},
			"ordering": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//Default:     "default",
				Description: "Controls the order of the Cipher String list in the Cipher Audit section. Options are Default, Speed, Strength, FIPS, and Hardware. The rules are processed in the order listed",
			},
			"allow": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the configuration of the allowed groups of ciphers. You can select a cipher rule from the Available Cipher Rules list",
			},
			"require": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the configuration of the restrict groups of ciphers. You can select a cipher rule from the Available Cipher Rules list",
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

	cipherGroup, err := client.GetLtmCipherGroup(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve cipher group %s  %v :", name, err)
		return diag.FromErr(err)
	}
	_ = d.Set("name", cipherGroup.FullPath)
	_ = d.Set("ordering", cipherGroup.Ordering)
	log.Printf("[INFO] Cipher group response :%+v", cipherGroup)
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
