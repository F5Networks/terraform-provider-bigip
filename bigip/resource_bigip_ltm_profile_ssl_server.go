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

func resourceBigipLtmProfileServerSsl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmProfileServerSslCreate,
		UpdateContext: resourceBigipLtmProfileServerSslUpdate,
		ReadContext:   resourceBigipLtmProfileServerSslRead,
		DeleteContext: resourceBigipLtmProfileServerSslDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the Ssl Profile",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/serverssl",
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

			"c3d_ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "none",
				Description: "CA Certificate. Default none.",
			},

			"c3d_ca_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "none",
				Description: "CA Key.  Default none.",
			},

			"c3d_ca_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CA Passphrase. Default",
			},

			"c3d_certificate_extensions": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CA Passphrase. Default enabled",
			},

			"c3d_cert_extension_custom_oids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Certificate Extensions List.  Default",
			},

			"c3d_cert_extension_includes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Certificate Extensions Includes. Default Extensions List",
			},

			"c3d_cert_lifespan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Certificate Lifespan.  Default",
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
				Default:  "/Common/default.crt",
				// ValidateFunc: validateF5NameWithDirectory,
				Description: "Name of the server certificate.",
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/Common/default.key",
				// ValidateFunc: validateF5NameWithDirectory,
				Description: "Name of the Server SSL profile key",
			},
			"chain": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "none",
				// ValidateFunc: validateF5NameWithDirectory,
				Description: "Server certificate chain name.",
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "Client Certificate Constrained Delegation CA passphrase",
			},

			"ciphers": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "BigIP Cipher string.",
			},

			"cipher_group": {
				Type:     schema.TypeString,
				Optional: true,
				// Computed rather than Default: "none" -- BIG-IP's own
				// default for an unconfigured cipher_group varies by
				// TMOS version (confirmed directly: 21.0.x defaults to
				// "none" with ciphers = "DEFAULT", while 21.1.x defaults
				// to "/Common/f5-default" with ciphers = "none", since
				// TLS 1.3 support now requires a cipher group to be
				// configured) -- a hardcoded client-side default can't
				// track that, and Default+Computed cannot coexist on the
				// same schema field (the SDK's own InternalValidate
				// rejects it) so this must be Computed-only, matching
				// resourceBigipLtmProfileServerSslRead's own
				// unconditional d.Set("cipher_group", ...) below.
				Computed: true,
				// No ConflictsWith against "ciphers" -- BIG-IP 21.1.x
				// legitimately requires setting both together
				// (cipher_group = "none" alongside an explicit ciphers
				// value) when overriding an inherited parent profile's
				// cipher_group, confirmed directly against a live
				// 21.1.0.1 device (same finding as
				// resource_bigip_ltm_profile_ssl_client.go's identical
				// cipher_group schema): setting ciphers alone (inheriting
				// cipher_group from the parent) is rejected with "cannot
				// contain both ciphers and a cipher-group", so the one
				// combination that actually works on 21.1.x is exactly
				// what a client-side ConflictsWith here would block.
				// Whether a given combination is actually valid now
				// depends on TMOS version and the profile's own
				// defaults_from chain, so BIG-IP itself (via an
				// apply-time API error) is the right place to enforce
				// this, not a static client-side rule.
				Description: "Cipher group for the ssl server profile",
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

			"ssl_c3d": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Client Certificate Constrained Delegation. Default disabled",
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

func resourceBigipLtmProfileServerSslCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Server Ssl Profile " + name)

	pss := &bigip.ServerSSLProfile{
		Name: name,
	}
	config := getServerSslConfig(d, pss)

	err := client.CreateServerSSLProfile(config)

	if err != nil {
		log.Printf("[ERROR] Unable to Create Server Ssl Profile (%s) (%v)", name, err)
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
		err = teemDevice.Report(f, "bigip_ltm_profile_server_ssl", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}

	return resourceBigipLtmProfileServerSslRead(ctx, d, meta)
}

func resourceBigipLtmProfileServerSslUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Updating ServerSSL Profile:%+v ", name)

	pss := &bigip.ServerSSLProfile{
		Name: name,
	}
	config := getServerSslConfig(d, pss)

	err := client.ModifyServerSSLProfile(name, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error create profile Ssl (%s): %s", name, err))
	}
	return resourceBigipLtmProfileServerSslRead(ctx, d, meta)
}

func resourceBigipLtmProfileServerSslRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Fetching Server SSL Profile " + name)
	obj, err := client.GetServerSSLProfile(name)

	if err != nil && strings.Contains(err.Error(), "not found") {
		log.Printf("[WARN] Server SSL Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Server SSL Profile   (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}

	if obj == nil {
		log.Printf("[WARN] Server SSL Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("name", name)
	_ = d.Set("partition", obj.Partition)

	_ = d.Set("defaults_from", obj.DefaultsFrom)
	_ = d.Set("alert_timeout", obj.AlertTimeout)
	_ = d.Set("authenticate", obj.Authenticate)
	_ = d.Set("authenticate_depth", obj.AuthenticateDepth)
	_ = d.Set("c3d_ca_cert", obj.C3dCaCert)
	_ = d.Set("c3d_ca_key", obj.C3dCaKey)
	_ = d.Set("c3d_ca_passphrase", obj.C3dCaPassphrase)
	_ = d.Set("c3d_cert_extension_custom_oids", obj.C3dCertExtensionCustomOids)
	_ = d.Set("c3d_cert_extension_includes", obj.C3dCertExtensionIncludes)
	_ = d.Set("c3d_cert_lifespan", obj.C3dCertLifespan)
	_ = d.Set("ca_file", obj.CaFile)
	_ = d.Set("cert", obj.Cert)
	_ = d.Set("chain", obj.Chain)
	// ciphers and cipher_group are both Optional+Computed (see their
	// schema comments for why neither can use a static Default anymore);
	// both must be set unconditionally here, the same way their
	// neighboring fields are, so a fresh Read after Create -- where the
	// user never set either explicitly -- still reflects BIG-IP's actual
	// returned value instead of leaving the attribute entirely absent.
	_ = d.Set("ciphers", obj.Ciphers)
	_ = d.Set("cipher_group", obj.CipherGroup)
	_ = d.Set("expire_cert_response_control", obj.ExpireCertResponseControl)
	_ = d.Set("cache_size", obj.CacheSize)
	_ = d.Set("handshake_timeout", obj.HandshakeTimeout)
	_ = d.Set("key", obj.Key)
	_ = d.Set("mod_ssl_methods", obj.ModSslMethods)
	_ = d.Set("mode", obj.Mode)
	_ = d.Set("proxy_ca_cert", obj.ProxyCaCert)
	_ = d.Set("proxy_ca_key", obj.ProxyCaKey)
	if obj.TmOptions == nil || obj.TmOptions == "none" {
		_ = d.Set("tm_options", []string{})
	} else {
		switch reflect.TypeOf(obj.TmOptions).Kind() {
		case reflect.String:
			tmOptions := strings.Split(obj.TmOptions.(string), " ")
			if len(tmOptions) > 0 {
				tmOptions = tmOptions[1:]
				tmOptions = tmOptions[:len(tmOptions)-1]
			}
			if err := d.Set("tm_options", tmOptions); err != nil {
				return diag.FromErr(fmt.Errorf("[DEBUG] Error saving TmOptions to state for Ssl profile  (%s): %s", d.Id(), err))
			}
		default:
			var newObj []string
			for _, v := range obj.TmOptions.([]interface{}) {
				newObj = append(newObj, v.(string))
			}
			_ = d.Set("tm_options", newObj)
		}
	}

	_ = d.Set("passphrase", obj.Passphrase)
	_ = d.Set("proxy_ssl", obj.ProxySsl)
	_ = d.Set("peer_cert_mode", obj.PeerCertMode)
	_ = d.Set("renegotiate_period", obj.RenegotiatePeriod)
	_ = d.Set("renegotiate_size", obj.RenegotiateSize)
	_ = d.Set("renegotiation", obj.Renegotiation)
	_ = d.Set("retain_certificate", obj.RetainCertificate)
	_ = d.Set("secure_renegotiation", obj.SecureRenegotiation)
	_ = d.Set("server_name", obj.ServerName)
	_ = d.Set("session_mirroring", obj.SessionMirroring)
	_ = d.Set("session_ticket", obj.SessionTicket)
	_ = d.Set("sni_default", obj.SniDefault)
	_ = d.Set("sni_require", obj.SniRequire)
	_ = d.Set("ssl_c3d", obj.SslC3d)
	_ = d.Set("ssl_forward_proxy", obj.SslForwardProxy)
	_ = d.Set("ssl_forward_proxy_bypass", obj.SslForwardProxyBypass)
	_ = d.Set("ssl_sign_hash", obj.SslSignHash)
	_ = d.Set("strict_resume", obj.StrictResume)
	_ = d.Set("unclean_shutdown", obj.UncleanShutdown)
	_ = d.Set("untrusted_cert_response_control", obj.UntrustedCertResponseControl)

	return nil
}

func resourceBigipLtmProfileServerSslDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Ssl Server Profile " + name)

	err := client.DeleteServerSSLProfile(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Ssl Profile (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getServerSslConfig(d *schema.ResourceData, config *bigip.ServerSSLProfile) *bigip.ServerSSLProfile {

	sslForwardProxyEnabled := d.Get("ssl_forward_proxy").(string)
	proxyCaCert := d.Get("proxy_ca_cert").(string)
	proxyCaKey := d.Get("proxy_ca_key").(string)
	sslForwardProxyBypass := d.Get("ssl_forward_proxy_bypass").(string)
	if sslForwardProxyEnabled == "enabled" {
		proxyCaCert = "/Common/default.crt"
		proxyCaKey = "/Common/default.key"
		if sslForwardProxyBypass == "" {
			sslForwardProxyBypass = "disabled"
		}
	}

	// tm_options has the same Optional+Computed defaults_from
	// inheritance-drift bug fixed in getClientSslConfig (SFDC #01262589;
	// see that function's comment in resource_bigip_ltm_profile_ssl_client.go
	// for the full root-cause explanation) -- server-ssl shares the
	// identical Optional+Computed tm_options schema/Read/Update pattern.
	// Read only from the raw config -- what the user actually wrote in
	// their current .tf -- to decide whether to include tm_options in the
	// outgoing payload.
	//
	// Unlike getClientSslConfig, tmOptionsIsSet is captured separately from
	// tmOptions itself (rather than just checking len(tmOptions) > 0 below)
	// so that a practitioner who explicitly writes tm_options = [] in
	// their .tf -- expressing intent to clear any inherited options rather
	// than merely omitting the attribute -- still has that empty list
	// assigned to config.TmOptions and sent to BIG-IP as an explicit
	// "options {}" clear, rather than being silently dropped and left at
	// its previously inherited value. getClientSslConfig does not yet make
	// this same distinction (it still checks len(tmOptions) > 0), so an
	// explicit tm_options = [] on bigip_ltm_profile_client_ssl is currently
	// silently not sent; this asymmetry should be resolved by applying the
	// same tmOptionsIsSet pattern there if/when that gap is addressed.
	tmOptionsIsSet := rawConfigAttrIsSet(d, "tm_options")
	var tmOptions []string
	if tmOptionsIsSet {
		tmOptions = setToStringSlice(d.Get("tm_options").(*schema.Set))
	}

	// Funky stuff toi get c3d to work
	var c3dCertExtensionCustomOids []string
	if t, ok := d.GetOk("c3d_cert_extension_custom_oids"); ok {
		c3dCertExtensionCustomOids = listToStringSlice(t.([]interface{}))
	}
	var c3dCertExtensionIncludes []string
	if t, ok := d.GetOk("c3d_cert_extension_includes"); ok {
		c3dCertExtensionIncludes = listToStringSlice(t.([]interface{}))
	}
	sslC3d := d.Get("ssl_c3d").(string)
	c3dCaCert := d.Get("c3d_ca_cert").(string)
	c3dCaKey := d.Get("c3d_ca_key").(string)

	if sslC3d == "enabled" {
		if c3dCaCert == "none" {
			c3dCaCert = "/Common/default.crt"
		}
		if c3dCaKey == "none" {
			c3dCaKey = "/Common/default.key"
		}
		if len(c3dCertExtensionIncludes) == 0 {
			c3dCertExtensionIncludes = []string{"basic-constraints", "extended-key-usage", "key-usage", "subject-alternative-name"}
		}
	} else {
		c3dCaCert = "none"
		c3dCaKey = "none"
	}

	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.Partition = d.Get("partition").(string)
	config.FullPath = d.Get("full_path").(string)
	config.Generation = d.Get("generation").(int)
	config.AlertTimeout = d.Get("alert_timeout").(string)
	config.Authenticate = d.Get("authenticate").(string)
	config.AuthenticateDepth = d.Get("authenticate_depth").(int)
	config.C3dCaCert = c3dCaCert
	config.C3dCaKey = c3dCaKey
	config.C3dCaPassphrase = d.Get("c3d_ca_passphrase").(string)
	config.C3dCertExtensionCustomOids = c3dCertExtensionCustomOids
	config.C3dCertExtensionIncludes = c3dCertExtensionIncludes
	config.C3dCertLifespan = d.Get("c3d_cert_lifespan").(int)
	config.CaFile = d.Get("ca_file").(string)
	config.CacheSize = d.Get("cache_size").(int)
	config.CacheTimeout = d.Get("cache_timeout").(int)
	config.Cert = d.Get("cert").(string)
	config.Key = d.Get("key").(string)
	config.Chain = d.Get("chain").(string)
	config.Passphrase = d.Get("passphrase").(string)
	// ciphers and cipher_group are both Optional+Computed: once Read
	// backfills state with an inherited value from a parent profile's
	// defaults_from, d.GetOk would return ok=true for that backfilled
	// value even though the practitioner's .tf never set it themselves,
	// incorrectly re-sending it as an explicit override on the very next
	// Update triggered by any other attribute changing -- see the
	// identical fix (and its rawConfigAttrIsSet precedent) in
	// resource_bigip_ltm_profile_ssl_client.go's getClientSslConfig,
	// confirmed there via
	// TestAccBigipLtmProfileClientSsl_ChildInheritanceAuditOtherFields.
	if rawConfigAttrIsSet(d, "ciphers") {
		config.Ciphers = d.Get("ciphers").(string)
		config.CipherGroup = "none"
	}
	if rawConfigAttrIsSet(d, "cipher_group") && d.Get("cipher_group").(string) != "none" {
		config.CipherGroup = d.Get("cipher_group").(string)
		config.Ciphers = "none"
	}
	config.ExpireCertResponseControl = d.Get("expire_cert_response_control").(string)
	config.GenericAlert = d.Get("generic_alert").(string)
	config.HandshakeTimeout = d.Get("handshake_timeout").(string)
	config.ModSslMethods = d.Get("mod_ssl_methods").(string)
	config.Mode = d.Get("mode").(string)
	config.ProxyCaCert = proxyCaCert
	config.ProxyCaKey = proxyCaKey
	config.PeerCertMode = d.Get("peer_cert_mode").(string)
	config.ProxySsl = d.Get("proxy_ssl").(string)
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
	config.SslC3d = sslC3d
	config.SslForwardProxy = sslForwardProxyEnabled
	config.SslForwardProxyBypass = sslForwardProxyBypass
	config.SslSignHash = d.Get("ssl_sign_hash").(string)
	config.StrictResume = d.Get("strict_resume").(string)
	config.UncleanShutdown = d.Get("unclean_shutdown").(string)
	config.UntrustedCertResponseControl = d.Get("untrusted_cert_response_control").(string)

	if tmOptionsIsSet {
		config.TmOptions = tmOptions
	}
	return config
}
