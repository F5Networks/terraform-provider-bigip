/*
Copyright 2021 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBigipIpsecPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipIpsecPolicyCreate,
		Read:   resourceBigipIpsecPolicyRead,
		Update: resourceBigipIpsecPolicyUpdate,
		Delete: resourceBigipIpsecPolicyDelete,
		Exists: resourceBigipIpsecPolicyExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specifies the name of the IPsec policy.",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Description of the IPsec policy.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies the IPsec protocol.",
				ValidateFunc: validation.StringInSlice([]string{"ah", "esp"}, false),
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies the processing mode.",
				ValidateFunc: validation.StringInSlice([]string{"transport", "interface", "isession", "tunnel"}, false),
			},
			"tunnel_local_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the local endpoint IP address of the IPsec tunnel. This parameter is only valid when mode is tunnel.",
			},
			"tunnel_remote_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the remote endpoint IP address of the IPsec tunnel. This parameter is only valid when mode is tunnel.",
			},
			"encrypt_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the algorithm to use for IKE encryption.",
				ValidateFunc: validation.StringInSlice([]string{"null", "3des", "aes128", "aes192", "aes256", "aes-gmac256",
					"aes-gmac192", "aes-gmac128", "aes-gcm256", "aes-gcm192", "aes-gcm256", "aes-gcm128"}, false),
			},
			"auth_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the algorithm to use for IKE authentication.",
				ValidateFunc: validation.StringInSlice([]string{"sha1", "sha256", "sha384", "sha512", "aes-gcm128",
					"aes-gcm192", "aes-gcm256", "aes-gmac128", "aes-gmac192",
					"aes-gmac256"}, false),
			},
			"lifetime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the length of time before the IKE security association expires, in minutes.",
			},
			"kb_lifetime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the length of time before the IKE security association expires, in kilobytes.",
			},
			"perfect_forward_secrecy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the Diffie-Hellman group to use for IKE Phase 2 negotiation.",
				ValidateFunc: validation.StringInSlice([]string{"none", "modp768", "modp1024", "modp1536", "modp2048", "modp3072",
					"modp4096", "modp6144", "modp8192"}, false),
			},
			"ipcomp": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies whether to use IPComp encapsulation.",
				ValidateFunc: validation.StringInSlice([]string{"none", "null", "deflate"}, false),
			},
		},
	}
}

func resourceBigipIpsecPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating IPSec Policy " + name)

	selectorConfig := &bigip.IPSecPolicy{
		Name:                           d.Get("name").(string),
		Description:                    d.Get("description").(string),
		Protocol:                       d.Get("protocol").(string),
		Mode:                           d.Get("mode").(string),
		TunnelLocalAddress:             d.Get("tunnel_local_address").(string),
		TunnelRemoteAddress:            d.Get("tunnel_remote_address").(string),
		IkePhase2EncryptAlgorithm:      d.Get("encrypt_algorithm").(string),
		IkePhase2AuthAlgorithm:         d.Get("auth_algorithm").(string),
		IkePhase2Lifetime:              d.Get("lifetime").(int),
		IkePhase2LifetimeKilobytes:     d.Get("kb_lifetime").(int),
		IkePhase2PerfectForwardSecrecy: d.Get("perfect_forward_secrecy").(string),
		Ipcomp:                         d.Get("ipcomp").(string),
	}

	err := client.CreateIPSecPolicy(selectorConfig)
	if err != nil {
		log.Printf("[ERROR] Unable to Create IPSec policy (%s) (%v)", name, err)
		return err
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
		err = teemDevice.Report(f, "bigip_ipsec_policy", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipIpsecPolicyRead(d, meta)
}

func resourceBigipIpsecPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Reading IPSec policy :%+v", name)
	ipsec, err := client.GetIPSecPolicy(name)
	if err != nil {
		return err
	}
	if ipsec == nil {
		d.SetId("")
		return fmt.Errorf("[ERROR] IPSec policy (%s) not found, removing from state", d.Id())
	}
	log.Printf("[DEBUG] IPSec Policy:%+v", ipsec)
	if err := d.Set("protocol", ipsec.Protocol); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Protocol to state for IPSec policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("mode", ipsec.Mode); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mode to state for IPSec policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("tunnel_local_address", ipsec.TunnelLocalAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving TunnelLocalAddress to state for IPSec policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("tunnel_remote_address", ipsec.TunnelRemoteAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving TunnelRemoteAddress to state for IPSec policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("encrypt_algorithm", ipsec.IkePhase2EncryptAlgorithm); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IkePhase2EncryptAlgorithm to state for IPSec Policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("auth_algorithm", ipsec.IkePhase2AuthAlgorithm); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IkePhase2AuthAlgorithm to state for IPSec Policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("lifetime", ipsec.IkePhase2Lifetime); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IkePhase2Lifetime to state for IPSec Policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("kb_lifetime", ipsec.IkePhase2LifetimeKilobytes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IkePhase2LifetimeKilobytes to state for IPSec Policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("perfect_forward_secrecy", ipsec.IkePhase2PerfectForwardSecrecy); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IkePhase2PerfectForwardSecrecy to state for IPSec Policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("ipcomp", ipsec.Ipcomp); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Ipcomp to state for IPSec Policy (%s): %s", d.Id(), err)
	}
	_ = d.Set("description", ipsec.Description)
	_ = d.Set("name", name)
	return nil
}

func resourceBigipIpsecPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Check existence of IPSec Policy: %+v ", name)
	ipsec, err := client.GetTrafficselctor(name)
	if err != nil {
		return false, err
	}
	if ipsec == nil {
		log.Printf("[WARN] IPSec Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, fmt.Errorf("[ERROR] IPSec Policy (%s) not found, removing from state", d.Id())
	}
	return true, nil
}

func resourceBigipIpsecPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating IPSec Policy:%+v ", name)
	ipsec := &bigip.IPSecPolicy{
		Name:                           name,
		Description:                    d.Get("description").(string),
		Protocol:                       d.Get("protocol").(string),
		Mode:                           d.Get("mode").(string),
		TunnelLocalAddress:             d.Get("tunnel_local_address").(string),
		TunnelRemoteAddress:            d.Get("tunnel_remote_address").(string),
		IkePhase2EncryptAlgorithm:      d.Get("encrypt_algorithm").(string),
		IkePhase2AuthAlgorithm:         d.Get("auth_algorithm").(string),
		IkePhase2Lifetime:              d.Get("lifetime").(int),
		IkePhase2LifetimeKilobytes:     d.Get("kb_lifetime").(int),
		IkePhase2PerfectForwardSecrecy: d.Get("perfect_forward_secrecy").(string),
		Ipcomp:                         d.Get("ipcomp").(string),
	}
	err := client.ModifyIPSecPolicy(name, ipsec)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify IPSec Policy   (%s) (%v) ", name, err)
		return err
	}
	return resourceBigipIpsecPolicyRead(d, meta)
}
func resourceBigipIpsecPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Deleting IPSec Policy:%+v ", name)
	err := client.DeleteIPSecPolicy(name)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to Delete IPSec Policy (%s) (%v) ", name, err)
	}
	d.SetId("")
	return nil
}
