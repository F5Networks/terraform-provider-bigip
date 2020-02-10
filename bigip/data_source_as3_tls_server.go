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
	//"google.golang.org/api/option"
	"log"
	//"strconv"
	//"reflect"
)

func dataSourceBigipAs3TlsServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipAs3TlsServerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of TLS_Server",
			},
			"certificates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AS3 pointer to Certificate declaration",
						},
						"match_to_sni": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If value is FQDN (wildcard okay), ignore all names in certificate and select this cert when SNI matches value (or by default)",
						},
					},
				},
				Description: "Primary and (optional) additional certificates (order is significant, element 0 is primary cert)",
			},
		},
	}

}
func dataSourceBigipAs3TlsServerRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	var tlsserver = &bigip.TlsServer{}
	tlsserver.Class = "TLS_Server"
	//tlsserver.AuthenticationFrequency = d.Get("authentication_frequency").(string)
	//tlsserver.C3dCertificateLifespan = d.Get("c3d_certificate_lifespan").(int)
	//tlsserver.C3dCertificateAuthority = d.Get("c3d_certificate_authority").(string)
	//tlsserver.Ciphers = d.Get("ciphers").(string)
	//tlsserver.ClientCertificate = d.Get("client_certificate").(string)
	//tlsserver.AllowExpiredCRL = d.Get("allow_expired_crl").(bool)
	//tlsserver.SendSNI = d.Get("send_sni").(string)
	//tlsserver.LdapStartTLS = d.Get("ldap_start_tls").(string)
	//tlsserver.Label = d.Get("label").(string)
	//tlsserver.Remark = d.Get("remark").(string)
	//tlsserver.C3dEnabled = d.Get("c3d_enabled").(bool)
	//tlsserver.IgnoreExpired = d.Get("ignore_expired").(bool)
	//tlsserver.IgnoreUntrusted = d.Get("ignore_untrusted").(bool)
	//tlsserver.SessionTickets = d.Get("session_tickets").(bool)
	//tlsserver.ValidateCertificate = d.Get("validate_certificate").(bool)
	//tlsserver.ServerName = d.Get("server_name").(string)
	//tlsserver.TrustCA = d.Get("trust_ca").(string)

	var certificateList []bigip.As3Certificates
	if m, ok := d.GetOk("certificates"); ok {
		var as3Cert = bigip.As3Certificates{}
		for _, v := range m.(*schema.Set).List() {
			as3Cert.Certificate = v.(map[string]interface{})["certificate"].(string)
			as3Cert.MatchToSNI = v.(map[string]interface{})["match_to_sni"].(string)
		}
		certificateList = append(certificateList, as3Cert)
	}
	tlsserver.Certificates = certificateList

	out, err := json.Marshal(tlsserver)
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
	d.SetId(string(out))
	log.Printf("[DEBUG] TLS Server Class string :%+v", string(out))
	return nil
}
