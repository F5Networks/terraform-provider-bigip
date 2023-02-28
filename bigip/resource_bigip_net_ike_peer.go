/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipNetIkePeer() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceBigipNetIkePeerCreate,
		ReadContext:   resourceBigipNetIkePeerRead,
		UpdateContext: resourceBigipNetIkePeerUpdate,
		DeleteContext: resourceBigipNetIkePeerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5Name,
				Description:  "Name of the IKE PEER",
			},
			"app_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application service that the object belongs to",
			},
			"ca_cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The trusted root and intermediate certificate authorities",
			},
			"crl_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the file name of the Certificate Revocation List. Only supported in IKEv1",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User defined description",
			},
			"generate_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable the generation of Security Policy Database entries(SPD) when the device is the responder of the IKE remote node",
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the exchange mode for phase 1 when racoon is the initiator, or the acceptable exchange mode when racoon is the responder",
			},
			"my_cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the name of the certificate file object",
			},
			"my_cert_key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the name of the certificate key file object",
			},
			"my_cert_key_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the passphrase of the key used for my-cert-key-file",
			},
			"my_id_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the identifier type sent to the remote host to use in the phase 1 negotiation",
			},
			"my_id_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the identifier value sent to the remote host in the phase 1 negotiation",
			},
			"nat_traversal": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enables use of the NAT-Traversal IPsec extension",
			},
			"passive": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether the local IKE agent can be the initiator of the IKE negotiation with this ike-peer",
			},
			"peers_cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the peer’s certificate for authentication",
			},
			"peers_cert_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies that the only peers-cert-type supported is certfile",
			},
			"peers_id_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies which of address, fqdn, asn1dn, user-fqdn or keyid-tag types to use as peers-id-type",
			},
			"peers_id_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the peer’s identifier to be received",
			},
			"phase1_auth_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the authentication method used for phase 1 negotiation",
			},
			"phase1_encrypt_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the encryption algorithm used for the isakmp phase 1 negotiation",
			},
			"phase1_hash_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the hash algorithm used for the isakmp phase 1 negotiation",
			},
			"phase1_perfect_forward_secrecy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the Diffie-Hellman group for key exchange to provide perfect forward secrecy",
			},
			"preshared_key": {
				Type:     schema.TypeString,
				Optional: true,
				//Computed:    true,
				Description: "Specifies the preshared key for ISAKMP SAs",
			},
			"preshared_key_encrypted": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Display the encrypted preshared-key for the IKE remote node",
			},
			"prf": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the pseudo-random function used to derive keying material for all cryptographic operations",
			},
			"proxy_support": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "If this value is enabled, both values of ID payloads in the phase 2 exchange are used as the addresses of end-point of IPsec-SAs",
			},
			"remote_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the IP address of the IKE remote node",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables this IKE remote node",
			},
			"traffic_selector": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				Description: "Specifies the names of the traffic-selector objects associated with this ike-peer",
			},
			"verify_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether to verify the certificate chain of the remote peer based on the trusted certificates in ca-cert-file",
			},
			"version": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				Description: "Specifies which version of IKE to be used",
			},
			"dpd_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the number of seconds between Dead Peer Detection messages",
			},
			"lifetime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Defines the lifetime in minutes of an IKE SA which will be proposed in the phase 1 negotiations",
			},
			"replay_window_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the replay window size of the IPsec SAs negotiated with the IKE remote node",
			},
		},
	}

}
func resourceBigipNetIkePeerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	r := &bigip.IkePeer{
		Name: name,
	}
	config := getIkeConfig(d, r)

	err := client.CreateIkePeer(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create IkePeer %s %v :", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)

	return resourceBigipNetIkePeerRead(ctx, d, meta)
}
func resourceBigipNetIkePeerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Reading IkePeer %s", name)
	ikepeer, err := client.GetIkePeer(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if ikepeer == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("[ERROR] IkePeer (%s) not found, removing from state", d.Id()))
	}
	_ = d.Set("app_service", ikepeer.AppService)

	_ = d.Set("my_cert_file", ikepeer.MyCertFile)

	_ = d.Set("my_cert_key_file", ikepeer.MyCertKeyFile)

	_ = d.Set("my_cert_key_passphrase", ikepeer.MyCertKeyPassphrase)

	if ikepeer.PresharedKey != "" && d.Get("preshared_key").(string) != "" {
		_ = d.Set("preshared_key", ikepeer.PresharedKey)
	}
	_ = d.Set("preshared_key_encrypted", ikepeer.PresharedKeyEncrypted)

	_ = d.Set("ca_cert_file", ikepeer.CaCertFile)

	_ = d.Set("crl_file", ikepeer.CrlFile)

	_ = d.Set("description", ikepeer.Description)

	_ = d.Set("generate_policy", ikepeer.GeneratePolicy)

	_ = d.Set("mode", ikepeer.Mode)

	_ = d.Set("my_id_type", ikepeer.MyIdType)

	_ = d.Set("my_id_value", ikepeer.MyIdValue)

	_ = d.Set("nat_traversal", ikepeer.NatTraversal)

	_ = d.Set("passive", ikepeer.Passive)

	_ = d.Set("peers_cert_file", ikepeer.PeersCertFile)

	_ = d.Set("peers_cert_type", ikepeer.PeersCertType)

	_ = d.Set("peers_id_type", ikepeer.PeersIdType)

	_ = d.Set("peers_id_value", ikepeer.PeersIdValue)

	_ = d.Set("phase1_auth_method", ikepeer.Phase1AuthMethod)

	_ = d.Set("phase1_encrypt_algorithm", ikepeer.Phase1EncryptAlgorithm)

	_ = d.Set("phase1_hash_algorithm", ikepeer.Phase1HashAlgorithm)

	_ = d.Set("phase1_perfect_forward_secrecy", ikepeer.Phase1PerfectForwardSecrecy)

	_ = d.Set("prf", ikepeer.Prf)

	_ = d.Set("proxy_support", ikepeer.ProxySupport)

	_ = d.Set("remote_address", ikepeer.RemoteAddress)

	_ = d.Set("state", ikepeer.State)

	_ = d.Set("traffic_selector", ikepeer.TrafficSelector)

	_ = d.Set("verify_cert", ikepeer.VerifyCert)

	_ = d.Set("version", ikepeer.Version)

	_ = d.Set("dpd_delay", ikepeer.DpdDelay)

	_ = d.Set("lifetime", ikepeer.Lifetime)
	_ = d.Set("replay_window_size", ikepeer.ReplayWindowSize)
	_ = d.Set("name", name)
	return nil
}
func resourceBigipNetIkePeerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Updating IkePeer %s", name)

	r := &bigip.IkePeer{
		Name: name,
	}
	config := getIkeConfig(d, r)

	err := client.ModifyIkePeer(name, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf(" Error modifying IkePeer %s: %v", name, err))
	}

	return resourceBigipNetIkePeerRead(ctx, d, meta)
}
func resourceBigipNetIkePeerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Deleting IkePeer %s", name)

	err := client.DeleteIkePeer(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf(" Error Deleting IkePeer : %s", err))
	}

	d.SetId("")
	return nil
}

func getIkeConfig(d *schema.ResourceData, config *bigip.IkePeer) *bigip.IkePeer {
	var version []string
	if t, ok := d.GetOk("version"); ok {
		version = listToStringSlice(t.([]interface{}))
	}
	var trafficSelectors []string
	if t, ok := d.GetOk("traffic_selector"); ok {
		trafficSelectors = listToStringSlice(t.([]interface{}))
	}
	config.AppService = d.Get("app_service").(string)
	config.CaCertFile = d.Get("ca_cert_file").(string)
	config.CrlFile = d.Get("crl_file").(string)
	config.DpdDelay = d.Get("dpd_delay").(int)
	config.Lifetime = d.Get("lifetime").(int)
	config.Description = d.Get("description").(string)
	config.GeneratePolicy = d.Get("generate_policy").(string)
	config.Mode = d.Get("mode").(string)
	config.MyCertFile = d.Get("my_cert_file").(string)
	config.MyCertKeyFile = d.Get("my_cert_key_file").(string)
	config.MyCertKeyPassphrase = d.Get("my_cert_key_passphrase").(string)
	config.MyIdType = d.Get("my_id_type").(string)
	config.MyIdValue = d.Get("my_id_value").(string)
	config.NatTraversal = d.Get("nat_traversal").(string)
	config.Passive = d.Get("passive").(string)
	config.PeersCertFile = d.Get("peers_cert_file").(string)
	config.PeersCertType = d.Get("peers_cert_type").(string)
	config.PeersIdType = d.Get("peers_id_type").(string)
	config.PeersIdValue = d.Get("peers_id_value").(string)
	config.Phase1AuthMethod = d.Get("phase1_auth_method").(string)
	config.Phase1EncryptAlgorithm = d.Get("phase1_encrypt_algorithm").(string)
	config.Phase1HashAlgorithm = d.Get("phase1_hash_algorithm").(string)
	config.Phase1PerfectForwardSecrecy = d.Get("phase1_perfect_forward_secrecy").(string)
	config.PresharedKey = d.Get("preshared_key").(string)
	config.Prf = d.Get("prf").(string)
	config.ProxySupport = d.Get("proxy_support").(string)
	config.RemoteAddress = d.Get("remote_address").(string)
	config.ReplayWindowSize = d.Get("replay_window_size").(int)
	config.State = d.Get("state").(string)
	config.TrafficSelector = trafficSelectors
	config.VerifyCert = d.Get("verify_cert").(string)
	config.Version = version

	return config
}
