/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBigipWafSignatures() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipWafSignatureRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the signature",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the signature",
			},
			"system_signature_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "System generated ID of the signature",
			},
			"signature_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the signature in the database",
			},
			"perform_staging": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The relative detection accuracy of the signature",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies, if true, that the signature is enabled on the security policy. When false, the signature is disable on the security policy.",
			},
			"tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The signature tag which, along with the signature name, identifies the signature.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Type of the signature",
			},
			"accuracy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The relative detection accuracy of the signature",
			},
			"risk": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The relative risk level of the attack that matches this signature",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The created JSON for WAF Signature block",
			},
		},
	}
}

func dataSourceBigipWafSignatureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	d.SetId("")
	sid := d.Get("signature_id").(int)
	provision := "asm"
	p, err := client.Provisions(provision)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Provision (%s) (%v) ", provision, err)
		return err
	}
	if p.Level == "none" {
		return fmt.Errorf("[ERROR] ASM Module is not provisioned, it is set to : (%s) ", p.Level)
	}

	signatures, err := client.GetWafSignature(sid)
	if err != nil {
		return fmt.Errorf("error retrieving signature %d: %v", sid, err)
	}

	// filter query always returns a list if the list is empty it means the signature is not found
	if len(signatures.Signatures) != 0 {
		// log.Printf("[DEBUG] Signature %d not found, removing from state", sid)

		// if successful filter query will return a list with a single item
		sign := signatures.Signatures[0]

		_ = d.Set("name", sign.Name)
		_ = d.Set("description", sign.Description)
		_ = d.Set("system_signature_id", sign.ResourceId)
		_ = d.Set("signature_id", sign.SignatureId)
		_ = d.Set("type", sign.Type)
		_ = d.Set("accuracy", sign.Accuracy)
		_ = d.Set("risk", sign.Risk)
	}

	sigJson := &bigip.WafSignature{
		SignatureID:    d.Get("signature_id").(int),
		Enabled:        d.Get("enabled").(bool),
		PerformStaging: d.Get("perform_staging").(bool),
	}
	jsonString, err := json.Marshal(sigJson)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Signature Json:%+v", string(jsonString))
	_ = d.Set("json", string(jsonString))

	d.SetId(strconv.Itoa(d.Get("signature_id").(int)))
	return nil
}
