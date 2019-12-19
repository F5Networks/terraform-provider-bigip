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

func resourceBigipLtmProfileClientSsl() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmProfileClientSSLCreate,
		Update: resourceBigipLtmProfileClientSSLUpdate,
		Read:   resourceBigipLtmProfileClientSSLRead,
		Delete: resourceBigipLtmProfileClientSSLDelete,
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

			"allow_non_ssl": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Allow non ssl",
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

			"cert_key_chain": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name",
						},

						"cert": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cert file name",
						},

						"chain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Chain file name",
						},

						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key filename",
						},

						"passphrase": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key passphrase",
						},
					},
				},
			},

			"cert_extension_includes": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Optional:    true,
				Computed:    true,
				Description: "Cert extension includes for ssl forward proxy",
			},

			"cert_life_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Life span of the certificate in days for ssl forward proxy",
			},

			"cert_lookup_by_ipaddr_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Cert lookup by ip address and port enabled / disabled",
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

			"client_cert_ca": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "client certificate name",
			},

			"crl_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Certificate revocation file name",
			},

			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/clientssl",
				Description: "Profile name that this profile defaults from.",
			},

			"forward_proxy_bypass_default_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Forward proxy bypass default action. (enabled / disabled)",
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

			"inherit_cert_keychain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Inherit cert key chain",
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the Server SSL profile key",
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
				Optional: true,
				Computed: true,
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

			"proxy_ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Proxy CA Cert",
			},

			"proxy_ca_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Proxy CA Key",
			},

			"proxy_ca_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Proxy CA Passphrase",
			},

			"proxy_ssl": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Proxy SSL enabled / disabled.  Default is disabled.",
			},

			"proxy_ssl_passthrough": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Proxy SSL passthrough enabled / disabled.  Default is disabled.",
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
		},
	}
}

func resourceBigipLtmProfileClientSSLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	parent := d.Get("defaults_from").(string)
	log.Println("[INFO] Creating Client Ssl Profile " + name)

	err := client.CreateClientSSLProfile(
		name,
		parent,
	)

	if err != nil {
		log.Printf("[ERROR] Unable to Create Client Ssl Profile (%s) (%v)", name, err)
		return err
	}

	d.SetId(name)

	err = resourceBigipLtmProfileClientSSLUpdate(d, meta)
	if err != nil {
		client.DeleteClientSSLProfile(name)
		return err
	}

	return resourceBigipLtmProfileClientSSLRead(d, meta)
}

func resourceBigipLtmProfileClientSSLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//log.Println("[INFO] Updating Route " + description)

	//var tmOptions []string
	//if t, ok := d.GetOk("tm_options"); ok {
	//	tmOptions = setToStringSlice(t.(*schema.Set))
	//}

	var CertExtensionIncludes []string
	if cei, ok := d.GetOk("cert_extension_includes"); ok {
		CertExtensionIncludes = setToStringSlice(cei.(*schema.Set))
	}

	type certKeyChain struct {
		Name       string "json:\"name,omitempty\""
		Cert       string "json:\"cert,omitempty\""
		Chain      string "json:\"chain,omitempty\""
		Key        string "json:\"key,omitempty\""
		Passphrase string "json:\"passphrase,omitempty\""
	}

	var certKeyChains []struct {
		Name       string "json:\"name,omitempty\""
		Cert       string "json:\"cert,omitempty\""
		Chain      string "json:\"chain,omitempty\""
		Key        string "json:\"key,omitempty\""
		Passphrase string "json:\"passphrase,omitempty\""
	}

	certKeyChainCount := d.Get("cert_key_chain.#").(int)
	for i := 0; i < certKeyChainCount; i++ {
		prefix := fmt.Sprintf("cert_key_chain.%d", i)
		certKeyChains = append(certKeyChains, certKeyChain{
			Name:       d.Get(prefix + ".name").(string),
			Cert:       d.Get(prefix + ".cert").(string),
			Chain:      d.Get(prefix + ".chain").(string),
			Key:        d.Get(prefix + ".key").(string),
			Passphrase: d.Get(prefix + ".passphrase").(string),
		})
	}

	pss := &bigip.ClientSSLProfile{
		Name:                            d.Get("name").(string),
		Partition:                       d.Get("partition").(string),
		FullPath:                        d.Get("full_path").(string),
		Generation:                      d.Get("generation").(int),
		AlertTimeout:                    d.Get("alert_timeout").(string),
		AllowNonSsl:                     d.Get("allow_non_ssl").(string),
		Authenticate:                    d.Get("authenticate").(string),
		AuthenticateDepth:               d.Get("authenticate_depth").(int),
		CaFile:                          d.Get("ca_file").(string),
		CacheSize:                       d.Get("cache_size").(int),
		CacheTimeout:                    d.Get("cache_timeout").(int),
		Cert:                            d.Get("cert").(string),
		CertExtensionIncludes:           CertExtensionIncludes,
		CertKeyChain:                    certKeyChains,
		CertLifespan:                    d.Get("cert_life_span").(int),
		CertLookupByIpaddrPort:          d.Get("cert_lookup_by_ipaddr_port").(string),
		Chain:                           d.Get("chain").(string),
		Ciphers:                         d.Get("ciphers").(string),
		ClientCertCa:                    d.Get("client_cert_ca").(string),
		CrlFile:                         d.Get("crl_file").(string),
		DefaultsFrom:                    d.Get("defaults_from").(string),
		ForwardProxyBypassDefaultAction: d.Get("forward_proxy_bypass_default_action").(string),
		GenericAlert:                    d.Get("generic_alert").(string),
		HandshakeTimeout:                d.Get("handshake_timeout").(string),
		InheritCertkeychain:             d.Get("inherit_cert_keychain").(string),
		Key:                             d.Get("key").(string),
		ModSslMethods:                   d.Get("mod_ssl_methods").(string),
		Mode:                            d.Get("mode").(string),
		//TmOptions:                       tmOptions,
		Passphrase:            d.Get("passphrase").(string),
		PeerCertMode:          d.Get("peer_cert_mode").(string),
		ProxyCaPassphrase:     d.Get("proxy_ca_passphrase").(string),
		ProxySsl:              d.Get("proxy_ssl").(string),
		ProxySslPassthrough:   d.Get("proxy_ssl_passthrough").(string),
		RenegotiatePeriod:     d.Get("renegotiate_period").(string),
		RenegotiateSize:       d.Get("renegotiate_size").(string),
		Renegotiation:         d.Get("renegotiation").(string),
		RetainCertificate:     d.Get("retain_certificate").(string),
		SecureRenegotiation:   d.Get("secure_renegotiation").(string),
		ServerName:            d.Get("server_name").(string),
		SessionMirroring:      d.Get("session_mirroring").(string),
		SessionTicket:         d.Get("session_ticket").(string),
		SniDefault:            d.Get("sni_default").(string),
		SniRequire:            d.Get("sni_require").(string),
		SslForwardProxy:       d.Get("ssl_forward_proxy").(string),
		SslForwardProxyBypass: d.Get("ssl_forward_proxy_bypass").(string),
		SslSignHash:           d.Get("ssl_sign_hash").(string),
		StrictResume:          d.Get("strict_resume").(string),
		UncleanShutdown:       d.Get("unclean_shutdown").(string),
	}

	err := client.ModifyClientSSLProfile(name, pss)
	if err != nil {
		return fmt.Errorf("Error create profile Ssl (%s): %s", name, err)
	}
	return resourceBigipLtmProfileClientSSLRead(d, meta)
}

func resourceBigipLtmProfileClientSSLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Fetching Client SSL Profile " + name)
	obj, err := client.GetClientSSLProfile(name)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Client SSL Profile   (%s) (%v) ", name, err)
		return err
	}

	if obj == nil {
		log.Printf("[WARN] Client SSL Profile (%s) not found, removing from state", d.Id())
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

	if err := d.Set("allow_non_ssl", obj.AllowNonSsl); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AllowNonSsl to state for Ssl profile  (%s): %s", d.Id(), err)
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

	if err := d.Set("cache_size", obj.CacheSize); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CacheSize to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("cache_timeout", obj.CacheTimeout); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CacheTimeout to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("cert", obj.Cert); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Cert to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	for i, c := range obj.CertKeyChain {
		ckc := fmt.Sprintf("cert_key_chain.%d", i)
		d.Set(fmt.Sprintf("%s.name", ckc), c.Name)
		d.Set(fmt.Sprintf("%s.cert", ckc), c.Cert)
		d.Set(fmt.Sprintf("%s.chain", ckc), c.Chain)
		d.Set(fmt.Sprintf("%s.key", ckc), c.Key)
		d.Set(fmt.Sprintf("%s.passphrase", ckc), c.Passphrase)

	}

	if err := d.Set("cert_extension_includes", obj.CertExtensionIncludes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CertExtensionIncludes to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("cert_life_span", obj.CertLifespan); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CertLifespan to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("cert_lookup_by_ipaddr_port", obj.CertLookupByIpaddrPort); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CertLookupByIpaddrPort to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("chain", obj.Chain); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Chain to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("ciphers", obj.Ciphers); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Ciphers to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("client_cert_ca", obj.ClientCertCa); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ClientCertCa to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("crl_file", obj.CrlFile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving CrlFile to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("forward_proxy_bypass_default_action", obj.ForwardProxyBypassDefaultAction); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ForwardProxyBypassDefaultAction to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("generic_alert", obj.GenericAlert); err != nil {
		return fmt.Errorf("[DEBUG] Error saving GenericAlert to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("handshake_timeout", obj.HandshakeTimeout); err != nil {
		return fmt.Errorf("[DEBUG] Error saving HandshakeTimeout to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("inherit_cert_keychain", obj.InheritCertkeychain); err != nil {
		return fmt.Errorf("[DEBUG] Error saving InheritCertkeychain to state for Ssl profile  (%s): %s", d.Id(), err)
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

	//if err := d.Set("tm_options", obj.TmOptions); err != nil {
	//	return fmt.Errorf("[DEBUG] Error saving TmOptions to state for Ssl profile  (%s): %s", d.Id(), err)
	//}

	if err := d.Set("passphrase", obj.Passphrase); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Passphrase to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("peer_cert_mode", obj.PeerCertMode); err != nil {
		return fmt.Errorf("[DEBUG] Error saving PeerCertMode to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("proxy_ca_passphrase", obj.ProxyCaPassphrase); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ProxyCaPassphrase to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("proxy_ssl", obj.ProxySsl); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ProxySsl to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("proxy_ssl_passthrough", obj.ProxySslPassthrough); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ProxySslPassthrough to state for Ssl profile  (%s): %s", d.Id(), err)
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

	return nil
}

func resourceBigipLtmProfileClientSSLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Ssl Client Profile " + name)

	err := client.DeleteClientSSLProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Ssl Profile (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}
