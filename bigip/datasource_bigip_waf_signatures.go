/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

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
				Computed:    true,
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
		},
	}
}

func dataSourceBigipWafSignatureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	d.SetId("")
	sid := d.Get("signature_id").(int)
	signatures, err := client.GetWafSignature(sid)
	if err != nil {
		return fmt.Errorf("error retrieving signature %d: %v", sid, err)
	}
	// filter query always returns a list if the list is empty it means the signature is not found
	if len(signatures.Signatures) == 0 {
		log.Printf("[DEBUG] Signature %d not found, removing from state", sid)
		d.SetId("")
		return nil
	}
	// if successful filter query will return a list with a single item
	sign := signatures.Signatures[0]

	_ = d.Set("name", sign.Name)
	_ = d.Set("description", sign.Description)
	_ = d.Set("system_signature_id", sign.ResourceId)
	_ = d.Set("signature_id", sign.SignatureId)
	_ = d.Set("type", sign.Type)
	_ = d.Set("accuracy", sign.Accuracy)
	_ = d.Set("risk", sign.Risk)

	d.SetId(sign.Name)
	return nil
}
