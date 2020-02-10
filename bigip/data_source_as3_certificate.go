/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceBigipAs3Certificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipAs3CertificateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of certificate",
			},
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "X.509 public-key certificate",
			},
			"chain_ca": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bundle of one or more CA certificates in trust-chain from root CA to certificate",
			},
			"pkcs12": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The pkcs12 value which may be a url to fetch the binary file from or base64 encoded string",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Private key matching certificate’s public key",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional friendly name for this object",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Arbitrary (brief) text pertaining to this object",
			},
			"passphrase": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_reuse": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "If true, other declaration objects may reuse this value",
						},
						"ciphertext": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Put base64url(data_value) here",
						},
						"ignore_changes": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "If false (default), the system updates the ciphertext in every AS3 declaration deployment. If true, AS3 creates the ciphertext on first deployment, and leaves it untouched afterwards",
						},
						"mini_jwe": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "If true (default), object is an f5 mini-JWE",
						},
						"protected": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"reuse_from": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AS3 pointer to another JWE cryptogram in this declaration to copy",
						},
					},
				},
				Description: "If supplied, used to decrypt privateKey at runtime",
			},
		},
	}

}
func dataSourceBigipAs3CertificateRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	var as3cert = &bigip.As3Certificate{}
	as3cert.Class = "Certificate"
	as3cert.Certificate = d.Get("certificate").(string)
	as3cert.ChainCA = d.Get("chain_ca").(string)
	as3cert.Pkcs12 = d.Get("pkcs12").(string)
	as3cert.Label = d.Get("label").(string)
	as3cert.Remark = d.Get("remark").(string)
	as3cert.PrivateKey = d.Get("private_key").(string)

	var passphraseList []bigip.As3Passphrase
	if m, ok := d.GetOk("passphrase"); ok {
		var as3Passphrase = bigip.As3Passphrase{}
		for _, v := range m.(*schema.Set).List() {
			as3Passphrase.AllowReuse = v.(map[string]interface{})["allow_reuse"].(bool)
			as3Passphrase.Ciphertext = v.(map[string]interface{})["ciphertext"].(string)
			as3Passphrase.IgnoreChanges = v.(map[string]interface{})["ignore_changes"].(bool)
			as3Passphrase.MiniJWE = v.(map[string]interface{})["mini_jwe"].(bool)
			as3Passphrase.Protected = v.(map[string]interface{})["protected"].(string)
			as3Passphrase.ReuseFrom = v.(map[string]interface{})["reuse_from"].(string)
		}
		passphraseList = append(passphraseList, as3Passphrase)
	}

	out, err := json.Marshal(as3cert)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	resultMap := make(map[string]interface{})
	resultMap[name] = string(out)
	out, err = json.Marshal(resultMap)
	if err != nil {
		return err
	}
	log.Printf("Certificate class string:%+v\n", string(out))
	d.SetId(string(out))
	return nil
}
