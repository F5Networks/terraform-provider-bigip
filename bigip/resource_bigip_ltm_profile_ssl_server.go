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
				ForceNew:     true,
				Description:  "Name of the Ssl Profile",
				ValidateFunc: validateF5NameWithDirectory,
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

func resourceBigipLtmProfileServerSslCreate(d *schema.ResourceData, meta interface{}) error {
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
		err = teemDevice.Report(f, "bigip_ltm_profile_server_ssl", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}

	return resourceBigipLtmProfileServerSslRead(d, meta)
}

func resourceBigipLtmProfileServerSslUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Updating ServerSSL Profile:%+v ", name)

	pss := &bigip.ServerSSLProfile{
		Name: name,
	}
	config := getServerSslConfig(d, pss)

	err := client.ModifyServerSSLProfile(name, config)
	if err != nil {
		return fmt.Errorf(" Error create profile Ssl (%s): %s", name, err)
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

	_ = d.Set("name", name)
	_ = d.Set("partition", obj.Partition)

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

	if err := d.Set("c3d_ca_cert", obj.C3dCaCert); err != nil {
		return fmt.Errorf("[DEBUG] Error saving C3dCaCert to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("c3d_ca_key", obj.C3dCaKey); err != nil {
		return fmt.Errorf("[DEBUG] Error saving C3dCaKey to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("c3d_ca_passphrase", obj.C3dCaPassphrase); err != nil {
		return fmt.Errorf("[DEBUG] Error saving C3dCaPassphrase to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("c3d_cert_extension_custom_oids", obj.C3dCertExtensionCustomOids); err != nil {
		return fmt.Errorf("[DEBUG] Error saving C3dCertExtensionCustomOids to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("c3d_cert_extension_includes", obj.C3dCertExtensionIncludes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving C3dCertExtensionIncludes to state for Ssl profile  (%s): %s", d.Id(), err)
	}

	if err := d.Set("c3d_cert_lifespan", obj.C3dCertLifespan); err != nil {
		return fmt.Errorf("[DEBUG] Error saving C3dCertLifespan to state for Ssl profile  (%s): %s", d.Id(), err)
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

	if err := d.Set("proxy_ca_cert", obj.ProxyCaCert); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mode to state for Ssl profile  (%s): %s", d.Id(), err)
	}
	if err := d.Set("proxy_ca_key", obj.ProxyCaKey); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mode to state for Ssl profile  (%s): %s", d.Id(), err)
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
		tmOptions := []string{}
		if err := d.Set("tm_options", tmOptions); err != nil {
			return fmt.Errorf("[DEBUG] Error saving TmOptions to state for Ssl profile  (%s): %s", d.Id(), err)
		}
	}

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

	if err := d.Set("ssl_c3d", obj.SslC3d); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SslC3d to state for Ssl profile  (%s): %s", d.Id(), err)
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

	var tmOptions []string
	if t, ok := d.GetOk("tm_options"); ok {
		tmOptions = setToStringSlice(t.(*schema.Set))
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
	config.Chain = d.Get("chain").(string)
	config.Ciphers = d.Get("ciphers").(string)
	config.ExpireCertResponseControl = d.Get("expire_cert_response_control").(string)
	config.GenericAlert = d.Get("generic_alert").(string)
	config.HandshakeTimeout = d.Get("handshake_timeout").(string)
	config.Key = d.Get("key").(string)
	config.ModSslMethods = d.Get("mod_ssl_methods").(string)
	config.Mode = d.Get("mode").(string)
	config.ProxyCaCert = proxyCaCert
	config.ProxyCaKey = proxyCaKey
	config.Passphrase = d.Get("passphrase").(string)
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

	if len(tmOptions) > 0 {
		config.TmOptions = tmOptions
	}
	return config
}
