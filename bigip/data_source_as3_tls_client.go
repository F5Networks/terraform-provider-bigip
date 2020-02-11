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

func dataSourceBigipAs3TlsClient() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipAs3TlsClientRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Nme of TLS_Client",
			},
			"allow_expired_crl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies if the CRL can be used even if it has expired",
			},
			"authentication_frequency": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "one-time",
				Description: "Client certificate authentication frequency",
			},
			"c3d_certificate_authority": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pointer to a Certificate class which specifies the Certificate Authority values for C3D",
			},
			"c3d_certificate_lifespan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     24,
				Description: "Specifies the lifespan of the certificate generated using the SSL client certificate constrained delegation",
			},
			"ciphers": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "DEFAULT",
				Description: "Ciphersuite selection string",
			},
			"client_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AS3 pointer to client Certificate declaration",
			},
			"c3d_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables or disables SSL Client certificate constrained delegation (C3D). Using C3D eliminates the need for requiring users to provide credentials twice for certain authentication actions",
			},
			"ignore_expired": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If false (default) drop connections with expired server certificates",
			},
			"ignore_untrusted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If false (default) drop connections with untrusted server certificates",
			},
			"c3d_certificate_extensions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Specifies the custom extension OID of the client certificates to be included in the generated certificates using SSL C3D",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional friendly name for this object",
			},
			"ldap_start_tls": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Creates a client LDAP profile with the specified activation mode STARTTLS",
			},
			"session_tickets": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If false (default) do not use rfc5077 session tickets",
			},
			"validate_certificate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If false (default) accept any cert from server, else validate server cert against trusted CA bundle",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Arbitrary (brief) text pertaining to this object",
			},
			"send_sni": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "none",
				Description: "FQDN to send in SNI",
			},
			"server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "none",
				Description: "FQDN which server certificate must match",
			},
			"trust_ca": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CA’s trusted to validate server certificate; ‘generic’ (default) or else AS3 pointer to declaration of CA Bundle",
			},
			"result_map": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				//Description:  "Name of service",
			},
		},
	}

}
func dataSourceBigipAs3TlsClientRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	var tlsclient = &bigip.TlsClient{}
	tlsclient.Class = "Service_HTTP"
	tlsclient.AuthenticationFrequency = d.Get("authentication_frequency").(string)
	tlsclient.C3dCertificateLifespan = d.Get("c3d_certificate_lifespan").(int)
	tlsclient.C3dCertificateAuthority = d.Get("c3d_certificate_authority").(string)
	tlsclient.Ciphers = d.Get("ciphers").(string)
	tlsclient.ClientCertificate = d.Get("client_certificate").(string)
	tlsclient.AllowExpiredCRL = d.Get("allow_expired_crl").(bool)
	tlsclient.SendSNI = d.Get("send_sni").(string)
	tlsclient.LdapStartTLS = d.Get("ldap_start_tls").(string)
	tlsclient.Label = d.Get("label").(string)
	tlsclient.Remark = d.Get("remark").(string)
	tlsclient.C3dEnabled = d.Get("c3d_enabled").(bool)
	tlsclient.IgnoreExpired = d.Get("ignore_expired").(bool)
	tlsclient.IgnoreUntrusted = d.Get("ignore_untrusted").(bool)
	tlsclient.SessionTickets = d.Get("session_tickets").(bool)
	tlsclient.ValidateCertificate = d.Get("validate_certificate").(bool)
	tlsclient.ServerName = d.Get("server_name").(string)
	tlsclient.TrustCA = d.Get("trust_ca").(string)

	out, err := json.Marshal(tlsclient)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	resultMap := make(map[string]interface{})
	resultMap[name] = string(out)
	d.Set("result_map", resultMap)
	log.Printf("resultMap in tls client class :%+v\n", resultMap)
	d.SetId(name)
	return nil
}
