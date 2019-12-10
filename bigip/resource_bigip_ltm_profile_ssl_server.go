/*
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

func resourceBigipLtmProfileServerSsl() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileServerSslCreate,
		Update: resourceBigipLtmProfileServerSslUpdate,
		Read:   resourceBigipLtmProfileServerSslRead,
		Delete: resourceBigipLtmProfileServerSslDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the Ssl Profile",
				ValidateFunc: validateF5Name,
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "name of partition",
			},

			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "full path of the profile",
			},

			"generation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "generation",
			},

			"alert_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Alert time out",
			},

			"authenticate": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Server authentication once / always (default is once).",
			},

			"authenticate_depth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Client certificate chain traversal depth.  Default 9.",
			},

			"ca_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Client certificate file path.  Default None.",
			},

			"cache_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache size (sessions).",
			},

			"cache_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time out",
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the server certificate.",
			},

			"chain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Client certificate chain name.",
			},

			"ciphers": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "BigIP Cipher string.",
			},

			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/serverssl",
				Description: "Profile name that this profile defaults from.",
			},

			"expire_cert_response_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Response if the cert is expired (drop / ignore). ",
			},

			"generic_alert": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Generic alerts enabled / disabled.",
			},

			"handshake_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Handshake time out (seconds)",
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the Serer SSL profile key",
			},

			"mod_ssl_methods": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ModSSL Methods enabled / disabled.  Default is disabled.",
			},

			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ModSSL Methods enabled / disabled.  Default is disabled.",
			},

			"tm_options": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
				Optional: true,
			},

			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Client Certificate Constrained Delegation CA passphrase",
			},

			"peer_cert_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Peer Cert Mode",
			},

			"proxy_ssl": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Proxy SSL enabled / disabled.  Default is disabled.",
			},

			"renegotiate_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Renogotiate Period (seconds)",
			},

			"renegotiate_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Renogotiate Size",
			},

			"renegotiation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Renegotiation (enabled / disabled)",
			},

			"retain_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Retain certificate.",
			},

			"secure_renegotiation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Secure reneogotiaton (request / require / require-strict).",
			},

			"server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Server name",
			},

			"session_mirroring": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Session Mirroring (enabled / disabled)",
			},

			"session_ticket": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Session Ticket (enabled / disabled)",
			},

			"sni_default": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SNI Default (true / false)",
			},

			"sni_require": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SNI Require (true / false)",
			},

			"ssl_forward_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SSL forward Proxy (enabled / disabled)",
			},

			"ssl_forward_proxy_bypass": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SSL forward Proxy Bypass (enabled / disabled)",
			},

			"ssl_sign_hash": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SSL sign hash (any, sha1, sha256, sha384)",
			},

			"strict_resume": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Strict Resume (enabled / disabled)",
			},

			"unclean_shutdown": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Unclean Shutdown (enabled / disabled)",
			},

			"untrusted_cert_response_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Unclean Shutdown (drop / ignore)",
			},
		},
	}
}

func resourceBigipLtmProfileServerSslCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)
	log.Println("[INFO] Creating Server Ssl Profile " + name)

	err := client.CreateServerSSLProfile(
		name,
		parent,
	)

	if err != nil {
		log.Printf("[ERROR] Unable to Create Server Ssl Profile (%s) (%v)", name, err)
		return err
	}

	d.SetId(name)

	err = resourceBigipLtmProfileServerSslUpdate(d, meta)
	if err != nil {
		client.DeleteServerSSLProfile(name)
		return err
	}

	return resourceBigipLtmProfileServerSslRead(d, meta)
}

func resourceBigipLtmProfileServerSslUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	/*var tmOptions []string
	if t, ok := d.GetOk("tm_options"); ok {
		tmOptions = setToStringSlice(t.(*schema.Set))
	}*/

	pss := &bigip.ServerSSLProfile{
		Name:                      d.Get("name").(string),
		Partition:                 d.Get("partition").(string),
		FullPath:                  d.Get("full_path").(string),
		Generation:                d.Get("generation").(int),
		AlertTimeout:              d.Get("alert_timeout").(string),
		Authenticate:              d.Get("authenticate").(string),
		AuthenticateDepth:         d.Get("authenticate_depth").(int),
		CaFile:                    d.Get("ca_file").(string),
		CacheSize:                 d.Get("cache_size").(int),
		CacheTimeout:              d.Get("cache_timeout").(int),
		Cert:                      d.Get("cert").(string),
		Chain:                     d.Get("chain").(string),
		Ciphers:                   d.Get("ciphers").(string),
		DefaultsFrom:              d.Get("defaults_from").(string),
		ExpireCertResponseControl: d.Get("expire_cert_response_control").(string),
		GenericAlert:              d.Get("generic_alert").(string),
		HandshakeTimeout:          d.Get("handshake_timeout").(string),
		Key:                       d.Get("key").(string),
		ModSslMethods:             d.Get("mod_ssl_methods").(string),
		Mode:                      d.Get("mode").(string),
		//TmOptions:                    tmOptions,
		Passphrase:                   d.Get("passphrase").(string),
		PeerCertMode:                 d.Get("peer_cert_mode").(string),
		ProxySsl:                     d.Get("proxy_ssl").(string),
		RenegotiatePeriod:            d.Get("renegotiate_period").(string),
		RenegotiateSize:              d.Get("renegotiate_size").(string),
		Renegotiation:                d.Get("renegotiation").(string),
		RetainCertificate:            d.Get("retain_certificate").(string),
		SecureRenegotiation:          d.Get("secure_renegotiation").(string),
		ServerName:                   d.Get("server_name").(string),
		SessionMirroring:             d.Get("session_mirroring").(string),
		SessionTicket:                d.Get("session_ticket").(string),
		SniDefault:                   d.Get("sni_default").(string),
		SniRequire:                   d.Get("sni_require").(string),
		SslForwardProxy:              d.Get("ssl_forward_proxy").(string),
		SslForwardProxyBypass:        d.Get("ssl_forward_proxy_bypass").(string),
		SslSignHash:                  d.Get("ssl_sign_hash").(string),
		StrictResume:                 d.Get("strict_resume").(string),
		UncleanShutdown:              d.Get("unclean_shutdown").(string),
		UntrustedCertResponseControl: d.Get("untrusted_cert_response_control").(string),
	}

	err := client.ModifyServerSSLProfile(name, pss)
	if err != nil {
		return fmt.Errorf("Error create profile Ssl (%s): %s", name, err)
	}
	return resourceBigipLtmProfileServerSslRead(d, meta)
}

func resourceBigipLtmProfileServerSslRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Fetching Server SSL Profile " + name)
	obj, err := client.GetServerSSLProfile(name)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Server SSL Profile   (%s) (%v) ", name, err)
		return err
	}

	if obj == nil {
		log.Printf("[WARN] Server SSL Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	d.Set("partition", obj.Partition)

	if err := d.Set("defaults_from", obj.DefaultsFrom); err != nil {
		return fmt.Errorf("[DEBUG] Error saving DefaultsFrom to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("alert_timeout", obj.AlertTimeout); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AlertTimeout to state for Ssl profile  (%s): %s", d.Id(), err)
	}
	if err := d.Set("authenticate", obj.Authenticate); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Authenticate to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("authenticate_depth", obj.AuthenticateDepth); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AuthenticateDepth to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("ca_file", obj.CaFile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CaFile to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("cert", obj.Cert); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Cert to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("chain", obj.Chain); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Chain to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("ciphers", obj.Ciphers); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Ciphers to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("expire_cert_response_control", obj.ExpireCertResponseControl); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ExpireCertResponseControl to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("cache_size", obj.CacheSize); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CacheSize to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("handshake_timeout", obj.HandshakeTimeout); err != nil {
		return fmt.Errorf("[DEBUG] Error saving HandshakeTimeout to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("key", obj.Key); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Key to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("mod_ssl_methods", obj.ModSslMethods); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ModSslMethods to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("mode", obj.Mode); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mode to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	/*if err := d.Set("tm_options", obj.TmOptions); err != nil {
		return fmt.Errorf("[DEBUG] Error saving TmOptions to state for Ssl profile  (%s): %s", d.Id(), err)
	}*/

	if err := d.Set("passphrase", obj.Passphrase); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Passphrase to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("proxy_ssl", obj.ProxySsl); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ProxySsl to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("peer_cert_mode", obj.PeerCertMode); err != nil {
		return fmt.Errorf("[DEBUG] Error saving PeerCertMode to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("renegotiate_period", obj.RenegotiatePeriod); err != nil {
		return fmt.Errorf("[DEBUG] Error saving RenegotiatePeriod to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("renegotiate_size", obj.RenegotiateSize); err != nil {
		return fmt.Errorf("[DEBUG] Error saving RenegotiateSize to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("renegotiation", obj.Renegotiation); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Renegotiation to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("retain_certificate", obj.RetainCertificate); err != nil {
		return fmt.Errorf("[DEBUG] Error saving RetainCertificate to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("secure_renegotiation", obj.SecureRenegotiation); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SecureRenegotiation to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("server_name", obj.ServerName); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ServerName to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("session_mirroring", obj.SessionMirroring); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SessionMirroring to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("session_ticket", obj.SessionTicket); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SessionTicket to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("sni_default", obj.SniDefault); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SniDefault to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("sni_require", obj.SniRequire); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SniRequire to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("ssl_forward_proxy", obj.SslForwardProxy); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SslForwardProxy to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("ssl_forward_proxy_bypass", obj.SslForwardProxyBypass); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SslForwardProxyBypass to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("ssl_sign_hash", obj.SslSignHash); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SslSignHash to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("strict_resume", obj.StrictResume); err != nil {
		return fmt.Errorf("[DEBUG] Error saving StrictResume to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("unclean_shutdown", obj.UncleanShutdown); err != nil {
		return fmt.Errorf("[DEBUG] Error saving UncleanShutdown to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("untrusted_cert_response_control", obj.UntrustedCertResponseControl); err != nil {
		return fmt.Errorf("[DEBUG] Error saving UntrustedCertResponseControl to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	return nil
}

/*
func resourceBigipLtmProfileServerSslRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetServerSSLProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrive Ssl Profile  (%s) (%v)", name, err)
		return err
	}
	if obj == nil {
		log.Printf("[WARN] Ssl  Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	d.Set("partition", obj.Partition)



	return nil
}
*/

func resourceBigipLtmProfileServerSslDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Ssl Server Profile " + name)

	err := client.DeleteServerSSLProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Ssl Profile (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}
