/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
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
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
			},

			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/clientssl",
				Description: "Profile name that this profile defaults from.",
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

			"c3d_client_fallback_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Client Fallback Certificate. Default None.",
			},

			"c3d_drop_unknown_ocsp_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Unknown OCSP Response Control. Default Drop.",
			},

			"c3d_ocsp": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "OCSP. Default None.",
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateF5Name,
				Description:  "Name of the server certificate.",
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
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								oldList := strings.Split(old, "/")
								newList := strings.Split(new, "/")
								if old == new {
									return true
								}
								if oldList[len(oldList)-1] == newList[len(newList)-1] {
									return true
								}
								return false
							},
						},

						"chain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Chain file name",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								oldList := strings.Split(old, "/")
								newList := strings.Split(new, "/")
								if old == new {
									return true
								}
								if oldList[len(oldList)-1] == newList[len(newList)-1] {
									return true
								}
								return false
							},
						},

						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key filename",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								oldList := strings.Split(old, "/")
								newList := strings.Split(new, "/")
								if old == new {
									return true
								}
								if oldList[len(oldList)-1] == newList[len(newList)-1] {
									return true
								}
								return false
							},
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateF5Name,
				Description:  "Name of the Server SSL profile key",
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

			"ssl_c3d": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Client Certificate Constrained Delegation enabled / disabled.  Default is disabled.",
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
	log.Printf("[INFO] Creating Client Ssl Profile:%+v ", name)

	pss := &bigip.ClientSSLProfile{
		Name: name,
	}
	config := getClientSslConfig(d, pss)
	err := client.CreateClientSSLProfile(config)

	if err != nil {
		log.Printf("[ERROR] Unable to Create Client Ssl Profile (%s) (%v)", name, err)
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
		err = teemDevice.Report(f, "bigip_ltm_profile_client_ssl", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}

	return resourceBigipLtmProfileClientSSLRead(d, meta)
}

func resourceBigipLtmProfileClientSSLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Updating Clientssl Profile : %v", name)

	pss := &bigip.ClientSSLProfile{
		Name: name,
	}
	config := getClientSslConfig(d, pss)
	err := client.ModifyClientSSLProfile(name, config)
	if err != nil {
		return fmt.Errorf(" Error create profile Ssl (%s): %s", name, err)
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

	_ = d.Set("name", name)
	_ = d.Set("partition", obj.Partition)
	_ = d.Set("defaults_from", obj.DefaultsFrom)

	if _, ok := d.GetOk("alert_timeout"); ok {
		if err := d.Set("alert_timeout", obj.AlertTimeout); err != nil {
			return fmt.Errorf("[DEBUG] Error saving AlertTimeout to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("allow_non_ssl"); ok {
		if err := d.Set("allow_non_ssl", obj.AllowNonSsl); err != nil {
			return fmt.Errorf("[DEBUG] Error saving AllowNonSsl to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("authenticate"); ok {
		if err := d.Set("authenticate", obj.Authenticate); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Authenticate to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("authenticate_depth"); ok {
		if err := d.Set("authenticate_depth", obj.AuthenticateDepth); err != nil {
			return fmt.Errorf("[DEBUG] Error saving AuthenticateDepth to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("c3d_client_fallback_cert"); ok {
		if err := d.Set("c3d_client_fallback_cert", obj.C3dClientFallbackCert); err != nil {
			return fmt.Errorf("[DEBUG] Error saving C3dClientFallbackCert to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("c3d_drop_unknown_ocsp_status"); ok {
		if err := d.Set("c3d_drop_unknown_ocsp_status", obj.C3dDropUnknownOcspStatus); err != nil {
			return fmt.Errorf("[DEBUG] Error saving C3dDropUnknownOcspStatus to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("c3d_ocsp"); ok {
		if err := d.Set("c3d_ocsp", obj.C3dOcsp); err != nil {
			return fmt.Errorf("[DEBUG] Error saving C3dOcsp to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("ca_file"); ok {
		if err := d.Set("ca_file", obj.CaFile); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CaFile to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("cache_size"); ok {
		if err := d.Set("cache_size", obj.CacheSize); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CacheSize to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("cache_timeout"); ok {
		if err := d.Set("cache_timeout", obj.CacheTimeout); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CacheTimeout to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("cert"); ok {
		if err := d.Set("cert", obj.Cert); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Cert to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	certMap := make(map[string]interface{})
	var certMapList []interface{}
	for _, c := range obj.CertKeyChain {
		certMap["name"] = c.Name
		certMap["cert"] = c.Cert
		certMap["key"] = c.Key
		certMap["chain"] = c.Chain
		certMap["passphrase"] = c.Passphrase
		certMapList = append(certMapList, certMap)
	}
	_ = d.Set("cert_key_chain", certMapList)

	if _, ok := d.GetOk("cert_extension_includes"); ok {
		if err := d.Set("cert_extension_includes", obj.CertExtensionIncludes); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CertExtensionIncludes to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("cert_life_span"); ok {
		if err := d.Set("cert_life_span", obj.CertLifespan); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CertLifespan to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("cert_lookup_by_ipaddr_port"); ok {
		if err := d.Set("cert_lookup_by_ipaddr_port", obj.CertLookupByIpaddrPort); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CertLookupByIpaddrPort to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("chain"); ok {
		if err := d.Set("chain", obj.Chain); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Chain to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("ciphers"); ok {
		if err := d.Set("ciphers", obj.Ciphers); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Ciphers to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("client_cert_ca"); ok {
		if err := d.Set("client_cert_ca", obj.ClientCertCa); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ClientCertCa to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("crl_file"); ok {
		if err := d.Set("crl_file", obj.CrlFile); err != nil {
			return fmt.Errorf("[DEBUG] Error saving CrlFile to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("forward_proxy_bypass_default_action"); ok {
		if err := d.Set("forward_proxy_bypass_default_action", obj.ForwardProxyBypassDefaultAction); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ForwardProxyBypassDefaultAction to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("generic_alert"); ok {
		if err := d.Set("generic_alert", obj.GenericAlert); err != nil {
			return fmt.Errorf("[DEBUG] Error saving GenericAlert to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("handshake_timeout"); ok {
		if err := d.Set("handshake_timeout", obj.HandshakeTimeout); err != nil {
			return fmt.Errorf("[DEBUG] Error saving HandshakeTimeout to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("inherit_cert_keychain"); ok {
		if err := d.Set("inherit_cert_keychain", obj.InheritCertkeychain); err != nil {
			return fmt.Errorf("[DEBUG] Error saving InheritCertkeychain to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("key"); ok {
		if err := d.Set("key", obj.Key); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Key to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("mod_ssl_methods"); ok {
		if err := d.Set("mod_ssl_methods", obj.ModSslMethods); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ModSslMethods to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("mode"); ok {
		if err := d.Set("mode", obj.Mode); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Mode to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}
	xt := reflect.TypeOf(obj.TmOptions).Kind()
	if obj.TmOptions != "none" {

		if xt == reflect.String {
			tmOptions := strings.Split(obj.TmOptions.(string), " ")
			if len(tmOptions) > 0 {
				tmOptions = tmOptions[1:]
				tmOptions = tmOptions[:len(tmOptions)-1]
			}
			if err := d.Set("tm_options", tmOptions); err != nil {
				return fmt.Errorf("[DEBUG] Error saving TmOptions to state for Ssl profile  (%s): %s", d.Id(), err)
			}
		} else {
			var newObj []string
			for _, v := range obj.TmOptions.([]interface{}) {
				newObj = append(newObj, v.(string))
			}
			if err := d.Set("tm_options", newObj); err != nil {
				return fmt.Errorf("[DEBUG] Error saving TmOptions to state for Ssl profile  (%s): %s", d.Id(), err)
			}
		}
	} else {
		var tmOptions []string
		if err := d.Set("tm_options", tmOptions); err != nil {
			return fmt.Errorf("[DEBUG] Error saving TmOptions to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("proxy_ca_cert"); ok {
		if err := d.Set("proxy_ca_cert", obj.ProxyCaCert); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Mode to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("proxy_ca_key"); ok {
		if err := d.Set("proxy_ca_key", obj.ProxyCaKey); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Mode to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("passphrase"); ok {
		if err := d.Set("passphrase", obj.Passphrase); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Passphrase to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("peer_cert_mode"); ok {
		if err := d.Set("peer_cert_mode", obj.PeerCertMode); err != nil {
			return fmt.Errorf("[DEBUG] Error saving PeerCertMode to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("proxy_ca_passphrase"); ok {
		if err := d.Set("proxy_ca_passphrase", obj.ProxyCaPassphrase); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ProxyCaPassphrase to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("proxy_ssl"); ok {
		if err := d.Set("proxy_ssl", obj.ProxySsl); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ProxySsl to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("proxy_ssl_passthrough"); ok {
		if err := d.Set("proxy_ssl_passthrough", obj.ProxySslPassthrough); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ProxySslPassthrough to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("renegotiate_period"); ok {
		if err := d.Set("renegotiate_period", obj.RenegotiatePeriod); err != nil {
			return fmt.Errorf("[DEBUG] Error saving RenegotiatePeriod to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("renegotiate_size"); ok {
		if err := d.Set("renegotiate_size", obj.RenegotiateSize); err != nil {
			return fmt.Errorf("[DEBUG] Error saving RenegotiateSize to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("renegotiation"); ok {
		if err := d.Set("renegotiation", obj.Renegotiation); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Renegotiation to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("retain_certificate"); ok {
		if err := d.Set("retain_certificate", obj.RetainCertificate); err != nil {
			return fmt.Errorf("[DEBUG] Error saving RetainCertificate to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("secure_renegotiation"); ok {
		if err := d.Set("secure_renegotiation", obj.SecureRenegotiation); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SecureRenegotiation to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("server_name"); ok {
		if err := d.Set("server_name", obj.ServerName); err != nil {
			return fmt.Errorf("[DEBUG] Error saving ServerName to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("session_mirroring"); ok {
		if err := d.Set("session_mirroring", obj.SessionMirroring); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SessionMirroring to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("session_ticket"); ok {
		if err := d.Set("session_ticket", obj.SessionTicket); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SessionTicket to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("sni_default"); ok {
		if err := d.Set("sni_default", obj.SniDefault); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SniDefault to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("sni_require"); ok {
		if err := d.Set("sni_require", obj.SniRequire); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SniRequire to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("ssl_c3d"); ok {
		if err := d.Set("ssl_c3d", obj.SslC3d); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SslC3d to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("ssl_forward_proxy"); ok {
		if err := d.Set("ssl_forward_proxy", obj.SslForwardProxy); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SslForwardProxy to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("ssl_forward_proxy_bypass"); ok {
		if err := d.Set("ssl_forward_proxy_bypass", obj.SslForwardProxyBypass); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SslForwardProxyBypass to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("ssl_sign_hash"); ok {
		if err := d.Set("ssl_sign_hash", obj.SslSignHash); err != nil {
			return fmt.Errorf("[DEBUG] Error saving SslSignHash to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("strict_resume"); ok {
		if err := d.Set("strict_resume", obj.StrictResume); err != nil {
			return fmt.Errorf("[DEBUG] Error saving StrictResume to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("unclean_shutdown"); ok {
		if err := d.Set("unclean_shutdown", obj.UncleanShutdown); err != nil {
			return fmt.Errorf("[DEBUG] Error saving UncleanShutdown to state for Ssl profile  (%s): %s", d.Id(), err)
		}
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

func getClientSslConfig(d *schema.ResourceData, config *bigip.ClientSSLProfile) *bigip.ClientSSLProfile {

	var tmOptions []string
	if t, ok := d.GetOk("tm_options"); ok {
		tmOptions = setToStringSlice(t.(*schema.Set))
	}
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

	sslForwardProxyEnabled := d.Get("ssl_forward_proxy").(string)
	sslForwardProxyBypass := d.Get("ssl_forward_proxy_bypass").(string)
	inheritCertkeychain := d.Get("inherit_cert_keychain").(string)
	proxyCaCert := d.Get("proxy_ca_cert").(string)
	proxyCaKey := d.Get("proxy_ca_key").(string)
	if sslForwardProxyEnabled == "enabled" {
		proxyCaCert = "/Common/default.crt"
		proxyCaKey = "/Common/default.key"
		inheritCertkeychain = "true"
		if sslForwardProxyBypass == "" {
			sslForwardProxyBypass = "disabled"
		}
	}
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.Partition = d.Get("partition").(string)
	config.FullPath = d.Get("full_path").(string)
	config.Generation = d.Get("generation").(int)
	config.AlertTimeout = d.Get("alert_timeout").(string)
	config.AllowNonSsl = d.Get("allow_non_ssl").(string)
	config.Authenticate = d.Get("authenticate").(string)
	config.AuthenticateDepth = d.Get("authenticate_depth").(int)
	config.C3dClientFallbackCert = d.Get("c3d_client_fallback_cert").(string)
	config.C3dDropUnknownOcspStatus = d.Get("c3d_drop_unknown_ocsp_status").(string)
	config.C3dOcsp = d.Get("c3d_ocsp").(string)
	config.CaFile = d.Get("ca_file").(string)
	config.CacheSize = d.Get("cache_size").(int)
	config.CacheTimeout = d.Get("cache_timeout").(int)
	log.Printf("[DEBUG] Length of certKeyChains :%+v", len(certKeyChains))
	log.Printf("[DEBUG] certKeyChains :%+v", certKeyChains)
	if len(certKeyChains) == 0 {
		config.Cert = d.Get("cert").(string)
		config.Key = d.Get("key").(string)
		config.Chain = d.Get("chain").(string)
		config.Passphrase = d.Get("passphrase").(string)
	}
	config.CertExtensionIncludes = CertExtensionIncludes
	config.CertKeyChain = certKeyChains
	config.CertLifespan = d.Get("cert_life_span").(int)
	config.CertLookupByIpaddrPort = d.Get("cert_lookup_by_ipaddr_port").(string)
	config.Ciphers = d.Get("ciphers").(string)
	config.ClientCertCa = d.Get("client_cert_ca").(string)
	config.CrlFile = d.Get("crl_file").(string)
	config.ForwardProxyBypassDefaultAction = d.Get("forward_proxy_bypass_default_action").(string)
	config.GenericAlert = d.Get("generic_alert").(string)
	config.HandshakeTimeout = d.Get("handshake_timeout").(string)
	config.InheritCertkeychain = inheritCertkeychain
	config.ModSslMethods = d.Get("mod_ssl_methods").(string)
	config.Mode = d.Get("mode").(string)
	config.PeerCertMode = d.Get("peer_cert_mode").(string)
	config.ProxyCaPassphrase = d.Get("proxy_ca_passphrase").(string)
	config.ProxySsl = d.Get("proxy_ssl").(string)
	config.ProxySslPassthrough = d.Get("proxy_ssl_passthrough").(string)
	config.RenegotiatePeriod = d.Get("renegotiate_period").(string)
	config.RenegotiateSize = d.Get("renegotiate_size").(string)
	config.Renegotiation = d.Get("renegotiation").(string)
	config.RetainCertificate = d.Get("retain_certificate").(string)
	config.SecureRenegotiation = d.Get("secure_renegotiation").(string)
	config.ServerName = d.Get("server_name").(string)
	config.SessionMirroring = d.Get("session_mirroring").(string)
	config.SessionTicket = d.Get("session_ticket").(string)
	config.SniDefault = d.Get("sni_default").(string)
	config.SniRequire = d.Get("sni_require").(string)
	config.SslC3d = d.Get("ssl_c3d").(string)
	config.SslForwardProxy = sslForwardProxyEnabled
	config.SslForwardProxyBypass = sslForwardProxyBypass
	config.SslSignHash = d.Get("ssl_sign_hash").(string)
	config.StrictResume = d.Get("strict_resume").(string)
	config.UncleanShutdown = d.Get("unclean_shutdown").(string)

	if len(tmOptions) > 0 {
		config.TmOptions = tmOptions
	}
	if proxyCaCert != "none" {
		config.ProxyCaCert = proxyCaCert
	}
	if proxyCaKey != "none" {
		config.ProxyCaKey = proxyCaKey
	}
	return config
}
