/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipLtmProfileClientSsl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmProfileClientSSLCreate,
		UpdateContext: resourceBigipLtmProfileClientSSLUpdate,
		ReadContext:   resourceBigipLtmProfileClientSSLRead,
		DeleteContext: resourceBigipLtmProfileClientSSLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//Default:      "/Common/default.crt",
				ValidateFunc: validateF5NameWithDirectory,
				Description:  "Name of the server certificate.",
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//Default:      "/Common/default.key",
				ValidateFunc: validateF5NameWithDirectory,
				Description:  "Name of the Server SSL profile key",
			},
			"chain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//Default:  "none",
				// ValidateFunc: validateF5NameWithDirectory,
				Description: "Client certificate chain name.",
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "Client Certificate Constrained Delegation CA passphrase",
			},
			"cert_key_chain": {
				Type:       schema.TypeList,
				Optional:   true,
				MaxItems:   1,
				Deprecated: "This Field 'cert_key_chain' going to deprecate in future version, please specify with cert,key,chain,passphrase as separate attribute.",
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
							Computed:    true,
							Sensitive:   true,
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

			"ciphers": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "BigIP Cipher string.",
			},

			"cipher_group": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "none",
				Description:   "Cipher group for the ssl client profile",
				ConflictsWith: []string{"ciphers"},
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

			"ocsp_stapling": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Specifies whether the system uses OCSP stapling.",
			},

			"tm_options": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
				Computed: true,
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

func resourceBigipLtmProfileClientSSLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
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

	return resourceBigipLtmProfileClientSSLRead(ctx, d, meta)
}

func resourceBigipLtmProfileClientSSLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Updating Clientssl Profile : %v", name)

	pss := &bigip.ClientSSLProfile{
		Name: name,
	}
	config := getClientSslConfig(d, pss)
	err := client.ModifyClientSSLProfile(name, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error create profile Ssl (%s): %s", name, err))
	}
	return resourceBigipLtmProfileClientSSLRead(ctx, d, meta)
}

func resourceBigipLtmProfileClientSSLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Fetching Client SSL Profile " + name)
	obj, err := client.GetClientSSLProfile(name)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Client SSL Profile   (%s) (%v) ", name, err)
		return diag.FromErr(err)
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
		_ = d.Set("alert_timeout", obj.AlertTimeout)
	}
	if _, ok := d.GetOk("allow_non_ssl"); ok {
		_ = d.Set("allow_non_ssl", obj.AllowNonSsl)
	}
	if _, ok := d.GetOk("authenticate"); ok {
		_ = d.Set("authenticate", obj.Authenticate)
	}
	if _, ok := d.GetOk("authenticate_depth"); ok {
		_ = d.Set("authenticate_depth", obj.AuthenticateDepth)
	}
	if _, ok := d.GetOk("c3d_client_fallback_cert"); ok {
		_ = d.Set("c3d_client_fallback_cert", obj.C3dClientFallbackCert)
	}
	if _, ok := d.GetOk("c3d_drop_unknown_ocsp_status"); ok {
		_ = d.Set("c3d_drop_unknown_ocsp_status", obj.C3dDropUnknownOcspStatus)
	}
	if _, ok := d.GetOk("c3d_ocsp"); ok {
		_ = d.Set("c3d_ocsp", obj.C3dOcsp)
	}
	if _, ok := d.GetOk("ca_file"); ok {
		_ = d.Set("ca_file", obj.CaFile)
	}
	if _, ok := d.GetOk("cache_size"); ok {
		_ = d.Set("cache_size", obj.CacheSize)
	}
	if _, ok := d.GetOk("cache_timeout"); ok {
		_ = d.Set("cache_timeout", obj.CacheTimeout)
	}
	if _, ok := d.GetOk("cert"); ok {
		_ = d.Set("cert", obj.Cert)
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
	log.Printf("certMapList:%+v", certMapList)

	if _, ok := d.GetOk("cert_extension_includes"); ok {
		_ = d.Set("cert_extension_includes", obj.CertExtensionIncludes)
	}
	if _, ok := d.GetOk("cert_life_span"); ok {
		_ = d.Set("cert_life_span", obj.CertLifespan)
	}
	if _, ok := d.GetOk("cert_lookup_by_ipaddr_port"); ok {
		_ = d.Set("cert_lookup_by_ipaddr_port", obj.CertLookupByIpaddrPort)
	}
	if _, ok := d.GetOk("chain"); ok {
		_ = d.Set("chain", obj.Chain)
	}
	if _, ok := d.GetOk("key"); ok {
		_ = d.Set("key", obj.Key)
	}
	if _, ok := d.GetOk("ciphers"); ok {
		_ = d.Set("ciphers", obj.Ciphers)
	}
	if _, ok := d.GetOk("cipher_group"); ok {
		_ = d.Set("cipher_group", obj.CipherGroup)
	}
	if _, ok := d.GetOk("client_cert_ca"); ok {
		_ = d.Set("client_cert_ca", obj.ClientCertCa)
	}

	if _, ok := d.GetOk("crl_file"); ok {
		_ = d.Set("crl_file", obj.CrlFile)
	}
	if _, ok := d.GetOk("forward_proxy_bypass_default_action"); ok {
		_ = d.Set("forward_proxy_bypass_default_action", obj.ForwardProxyBypassDefaultAction)
	}
	if _, ok := d.GetOk("generic_alert"); ok {
		_ = d.Set("generic_alert", obj.GenericAlert)
	}
	if _, ok := d.GetOk("handshake_timeout"); ok {
		_ = d.Set("handshake_timeout", obj.HandshakeTimeout)
	}
	if _, ok := d.GetOk("inherit_cert_keychain"); ok {
		_ = d.Set("inherit_cert_keychain", obj.InheritCertkeychain)
	}
	if _, ok := d.GetOk("mod_ssl_methods"); ok {
		_ = d.Set("mod_ssl_methods", obj.ModSslMethods)
	}
	if _, ok := d.GetOk("mode"); ok {
		_ = d.Set("mode", obj.Mode)
	}
	xt := reflect.TypeOf(obj.TmOptions).Kind()
	if obj.TmOptions != "none" {

		if xt == reflect.String {
			tmOptions := strings.Split(obj.TmOptions.(string), " ")
			if len(tmOptions) > 0 {
				tmOptions = tmOptions[1:]
				tmOptions = tmOptions[:len(tmOptions)-1]
			}
			_ = d.Set("tm_options", tmOptions)
		} else {
			var newObj []string
			for _, v := range obj.TmOptions.([]interface{}) {
				newObj = append(newObj, v.(string))
			}
			_ = d.Set("tm_options", newObj)
		}
	} else {
		var tmOptions []string
		_ = d.Set("tm_options", tmOptions)
	}

	if _, ok := d.GetOk("ocsp_stapling"); ok {
		_ = d.Set("ocsp_stapling", obj.OcspStapling)
	}

	if _, ok := d.GetOk("proxy_ca_cert"); ok {
		_ = d.Set("proxy_ca_cert", obj.ProxyCaCert)
	}

	if _, ok := d.GetOk("proxy_ca_key"); ok {
		_ = d.Set("proxy_ca_key", obj.ProxyCaKey)
	}

	if _, ok := d.GetOk("peer_cert_mode"); ok {
		_ = d.Set("peer_cert_mode", obj.PeerCertMode)
	}

	if _, ok := d.GetOk("proxy_ca_passphrase"); ok {
		_ = d.Set("proxy_ca_passphrase", obj.ProxyCaPassphrase)
	}

	if _, ok := d.GetOk("proxy_ssl"); ok {
		_ = d.Set("proxy_ssl", obj.ProxySsl)
	}

	if _, ok := d.GetOk("proxy_ssl_passthrough"); ok {
		_ = d.Set("proxy_ssl_passthrough", obj.ProxySslPassthrough)
	}

	if _, ok := d.GetOk("renegotiate_period"); ok {
		_ = d.Set("renegotiate_period", obj.RenegotiatePeriod)
	}

	if _, ok := d.GetOk("renegotiate_size"); ok {
		_ = d.Set("renegotiate_size", obj.RenegotiateSize)
	}
	if _, ok := d.GetOk("renegotiation"); ok {
		_ = d.Set("renegotiation", obj.Renegotiation)
	}

	if _, ok := d.GetOk("retain_certificate"); ok {
		_ = d.Set("retain_certificate", obj.RetainCertificate)
	}

	if _, ok := d.GetOk("secure_renegotiation"); ok {
		_ = d.Set("secure_renegotiation", obj.SecureRenegotiation)
	}

	if _, ok := d.GetOk("server_name"); ok {
		_ = d.Set("server_name", obj.ServerName)
	}

	if _, ok := d.GetOk("session_mirroring"); ok {
		_ = d.Set("session_mirroring", obj.SessionMirroring)
	}

	if _, ok := d.GetOk("session_ticket"); ok {
		_ = d.Set("session_ticket", obj.SessionTicket)
	}

	if _, ok := d.GetOk("sni_default"); ok {
		_ = d.Set("sni_default", obj.SniDefault)
	}

	if _, ok := d.GetOk("sni_require"); ok {
		_ = d.Set("sni_require", obj.SniRequire)
	}

	if _, ok := d.GetOk("ssl_c3d"); ok {
		_ = d.Set("ssl_c3d", obj.SslC3d)
	}

	if _, ok := d.GetOk("ssl_forward_proxy"); ok {
		_ = d.Set("ssl_forward_proxy", obj.SslForwardProxy)
	}

	if _, ok := d.GetOk("ssl_forward_proxy_bypass"); ok {
		_ = d.Set("ssl_forward_proxy_bypass", obj.SslForwardProxyBypass)
	}

	if _, ok := d.GetOk("ssl_sign_hash"); ok {
		_ = d.Set("ssl_sign_hash", obj.SslSignHash)
	}

	if _, ok := d.GetOk("strict_resume"); ok {
		_ = d.Set("strict_resume", obj.StrictResume)
	}

	if _, ok := d.GetOk("unclean_shutdown"); ok {
		_ = d.Set("unclean_shutdown", obj.UncleanShutdown)
	}

	return nil
}

func resourceBigipLtmProfileClientSSLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Ssl Client Profile " + name)

	err := client.DeleteClientSSLProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Ssl Profile (%s) (%v)", name, err)
		return diag.FromErr(err)
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
	config.OcspStapling = d.Get("ocsp_stapling").(string)
	log.Printf("[DEBUG] Length of certKeyChains :%+v", len(certKeyChains))
	log.Printf("[DEBUG] certKeyChains :%+v", certKeyChains)
	if len(certKeyChains) == 0 {
		config.Cert = d.Get("cert").(string)
		config.Key = d.Get("key").(string)
		config.Chain = d.Get("chain").(string)
		config.Passphrase = d.Get("passphrase").(string)
	} else {
		config.CertKeyChain = certKeyChains
	}
	config.CertExtensionIncludes = CertExtensionIncludes
	config.CertLifespan = d.Get("cert_life_span").(int)
	config.CertLookupByIpaddrPort = d.Get("cert_lookup_by_ipaddr_port").(string)
	if ciphers, ok := d.GetOk("ciphers"); ok {
		config.Ciphers = ciphers.(string)
		config.CipherGroup = "none"
	}
	if cipherGrp, ok := d.GetOk("cipher_group"); ok && cipherGrp != "none" {
		config.CipherGroup = cipherGrp.(string)
		config.Ciphers = "none"
	}
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
