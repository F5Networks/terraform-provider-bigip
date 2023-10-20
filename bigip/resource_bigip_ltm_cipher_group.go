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
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the configuration of the allowed groups of ciphers. You can select a cipher rule from the Available Cipher Rules list",
			},
			"require": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
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
	cipherGroup := getCipherGroupConfig(d, cipherGrouptmp)

	log.Printf("[INFO] cipherGroup config :%+v", cipherGroup)
	err := client.AddLtmCipherGroup(cipherGroup)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating cipher rule (%s): %s", name, err))
	}
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
		err = teemDevice.Report(f, "bigip_ltm_cipher_group", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
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
	var allowList []interface{}
	for _, val := range cipherGroup.Allow {
		tmpCipher := fmt.Sprintf("/%s/%s", val.(map[string]interface{})["partition"].(string), val.(map[string]interface{})["name"].(string))
		allowList = append(allowList, tmpCipher)
	}
	_ = d.Set("allow", allowList)
	var requireList []interface{}
	for _, val := range cipherGroup.Require {
		tmpCipher := fmt.Sprintf("/%s/%s", val.(map[string]interface{})["partition"].(string), val.(map[string]interface{})["name"].(string))
		requireList = append(requireList, tmpCipher)
	}
	_ = d.Set("require", requireList)
	return nil
}

func resourceBigipLtmCipherGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	cipherGrouptmp := &bigip.CipherGroupReq{}
	cipherGrouptmp.Name = name
	cipherGroupconfig := getCipherGroupConfig(d, cipherGrouptmp)
	if p, ok := d.GetOk("require"); ok {
		for _, r := range p.(*schema.Set).List() {
			cipherGroupconfig.Require = append(cipherGroupconfig.Require, r.(string))
		}
	}
	type CipherGroupReqnew struct {
		bigip.CipherGroupReq
		Require []interface{} `json:"require"`
		Allow   []interface{} `json:"allow"`
	}
	new := &CipherGroupReqnew{}
	new.Require = cipherGroupconfig.Require
	new.Name = cipherGroupconfig.Name
	new.Ordering = cipherGroupconfig.Ordering
	new.Allow = cipherGroupconfig.Allow

	if err := client.ModifyLtmCipherGroupNew(name, new); err != nil {
		return diag.FromErr(fmt.Errorf("error modifying cipher group %s: %v", name, err))
	}

	//
	// if err := client.ModifyLtmCipherGroup(name, cipherGroupconfig); err != nil {
	//	return diag.FromErr(fmt.Errorf("error modifying cipher group %s: %v", name, err))
	// }

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

func getCipherGroupConfig(d *schema.ResourceData, cipherGroup *bigip.CipherGroupReq) *bigip.CipherGroupReq {
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
	return cipherGroup
}
