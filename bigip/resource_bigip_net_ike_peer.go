/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipNetIkePeer() *schema.Resource {

	return &schema.Resource{
		Create: resourceBigipNetIkePeerCreate,
		Read:   resourceBigipNetIkePeerRead,
		Update: resourceBigipNetIkePeerUpdate,
		Delete: resourceBigipNetIkePeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
func resourceBigipNetIkePeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	r := &bigip.IkePeer{
		Name: name,
	}
	config := getIkeConfig(d, r)

	err := client.CreateIkePeer(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create IkePeer %s %v :", name, err)
		return err
	}
	d.SetId(name)

	return resourceBigipNetIkePeerRead(d, meta)
}
func resourceBigipNetIkePeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Reading IkePeer %s", name)
	ikepeer, err := client.GetIkePeer(name)
	if err != nil {
		return err
	}
	if ikepeer == nil {
		d.SetId("")
		return fmt.Errorf("[ERROR] IkePeer (%s) not found, removing from state", d.Id())
	}
	log.Printf("[DEBUG] IkePeer:%+v", ikepeer)
	if err := d.Set("app_service", ikepeer.AppService); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AppService to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("my_cert_file", ikepeer.MyCertFile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving MyCertFile to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("my_cert_key_file", ikepeer.MyCertKeyFile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving MyCertKeyFile to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("my_cert_key_passphrase", ikepeer.MyCertKeyPassphrase); err != nil {
		return fmt.Errorf("[DEBUG] Error saving MyCertKeyPassphrase to state for IkePeer (%s): %s", d.Id(), err)
	}
	if ikepeer.PresharedKey != "" && d.Get("preshared_key").(string) != "" {
		_ = d.Set("preshared_key", ikepeer.PresharedKey)
	}
	if err := d.Set("preshared_key_encrypted", ikepeer.PresharedKeyEncrypted); err != nil {
		return fmt.Errorf("[DEBUG] Error saving PresharedKeyEncrypted to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("ca_cert_file", ikepeer.CaCertFile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CaCertFile to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("crl_file", ikepeer.CrlFile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving  to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("description", ikepeer.Description); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CrlFile to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("generate_policy", ikepeer.GeneratePolicy); err != nil {
		return fmt.Errorf("[DEBUG] Error saving GeneratePolicy to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("mode", ikepeer.Mode); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mode to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("my_id_type", ikepeer.MyIdType); err != nil {
		return fmt.Errorf("[DEBUG] Error saving MyIdType to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("my_id_value", ikepeer.MyIdValue); err != nil {
		return fmt.Errorf("[DEBUG] Error saving MyIdValue to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("nat_traversal", ikepeer.NatTraversal); err != nil {
		return fmt.Errorf("[DEBUG] Error saving NatTraversal to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("passive", ikepeer.Passive); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Passive to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("peers_cert_file", ikepeer.PeersCertFile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving PeersCertFile to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("peers_cert_type", ikepeer.PeersCertType); err != nil {
		return fmt.Errorf("[DEBUG] Error saving PeersCertType to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("peers_id_type", ikepeer.PeersIdType); err != nil {
		return fmt.Errorf("[DEBUG] Error saving PeersIdType to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("peers_id_value", ikepeer.PeersIdValue); err != nil {
		return fmt.Errorf("[DEBUG] Error saving PeersIdValue to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("phase1_auth_method", ikepeer.Phase1AuthMethod); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Phase1AuthMethod to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("phase1_encrypt_algorithm", ikepeer.Phase1EncryptAlgorithm); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Phase1EncryptAlgorithm to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("phase1_hash_algorithm", ikepeer.Phase1HashAlgorithm); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Phase1HashAlgorithm to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("phase1_perfect_forward_secrecy", ikepeer.Phase1PerfectForwardSecrecy); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Phase1PerfectForwardSecrecy to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("prf", ikepeer.Prf); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Prf to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("proxy_support", ikepeer.ProxySupport); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ProxySupport to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("remote_address", ikepeer.RemoteAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving RemoteAddress to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("state", ikepeer.State); err != nil {
		return fmt.Errorf("[DEBUG] Error saving State to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("traffic_selector", ikepeer.TrafficSelector); err != nil {
		return fmt.Errorf("[DEBUG] Error saving TrafficSelector to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("verify_cert", ikepeer.VerifyCert); err != nil {
		return fmt.Errorf("[DEBUG] Error saving VerifyCert to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("version", ikepeer.Version); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Version to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("dpd_delay", ikepeer.DpdDelay); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DpdDelay to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("lifetime", ikepeer.Lifetime); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Lifetime to state for IkePeer (%s): %s", d.Id(), err)
	}
	if err := d.Set("replay_window_size", ikepeer.ReplayWindowSize); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ReplayWindowSize to state for IkePeer (%s): %s", d.Id(), err)
	}
	_ = d.Set("name", name)
	return nil
}
func resourceBigipNetIkePeerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Updating IkePeer %s", name)

	r := &bigip.IkePeer{
		Name: name,
	}
	config := getIkeConfig(d, r)

	err := client.ModifyIkePeer(name, config)
	if err != nil {
		return fmt.Errorf(" Error modifying IkePeer %s: %v", name, err)
	}

	return resourceBigipNetIkePeerRead(d, meta)
}
func resourceBigipNetIkePeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Deleting IkePeer %s", name)

	err := client.DeleteIkePeer(name)
	if err != nil {
		return fmt.Errorf(" Error Deleting IkePeer : %s", err)
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
