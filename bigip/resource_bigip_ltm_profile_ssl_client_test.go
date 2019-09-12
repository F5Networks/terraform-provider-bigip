/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_CLIENTSSL_NAME = fmt.Sprintf("/%s/test-ClientSsl", TEST_PARTITION)

//Many of the values below are intentionally non default settings so that the unit tests will avoid defaults that are set by f5 but possibly not updated by the REST call (potential hidden bug)
var TEST_CLIENTSSL_RESOURCE = `
resource "bigip_ltm_profile_client_ssl" "test-ClientSsl" {
	name                                = "/Common/test-ClientSsl"
	partition                           = "Common"
	defaults_from                       = "/Common/clientssl"
	alert_timeout                       = "60"
	allow_non_ssl                       = "enabled"
	authenticate                        = "always"
	authenticate_depth                  = 7
	ca_file                             = "none"
	cache_size                          = 131072
	cache_timeout                       = 2400
	cert                                = "/Common/default.crt"
	cert_key_chain {
		cert = "/Common/default.crt"
		key  = "/Common/default.key"
		name = "default"
	}
	cert_extension_includes             = [
		"basic-constraints",
		"subject-alternative-name",
	  ]
	cert_life_span                      = 60
	cert_lookup_by_ipaddr_port          = "enabled"
	chain                               = "none"
	ciphers                             = "DEFAULT"
	client_cert_ca                      = "none"
	crl_file                            = "none"
	forward_proxy_bypass_default_action = "intercept"
	generic_alert                       = "disabled"
	handshake_timeout                   = "30"
	inherit_cert_keychain               = "false"
	key                                 = "/Common/default.key"
	mod_ssl_methods                     = "enabled"
	mode                                = "disabled"
	peer_cert_mode                      = "require"
	proxy_ca_cert                       = "none"
	proxy_ca_key                        = "none"
	proxy_ca_passphrase                 = ""
	proxy_ssl                           = "enabled"
	proxy_ssl_passthrough               = "enabled"
	renegotiate_period                  = "60"
	renegotiate_size                    = "2"
	renegotiation                       = "enabled"
	retain_certificate                  = "false"
	secure_renegotiation                = "request"
	server_name                         = "testservername"
	session_mirroring                   = "disabled"
	session_ticket                      = "enabled"
	sni_default                         = "true"
	sni_require                         = "true"
	ssl_forward_proxy                   = "disabled"
	ssl_forward_proxy_bypass            = "disabled"
	ssl_sign_hash                       = "sha1"
	strict_resume                       = "enabled"
	tm_options                          = [
		"dont-insert-empty-fragments",
	  ]
	unclean_shutdown                    = "disabled"
  }
`

func TestAccBigipLtmProfileClientSsl_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_CLIENTSSL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(TEST_CLIENTSSL_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "name", "/Common/test-ClientSsl"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "alert_timeout", "60"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "allow_non_ssl", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "authenticate", "always"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "authenticate_depth", "7"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "ca_file", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "cache_size", "131072"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "cache_timeout", "2400"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl",
						fmt.Sprintf("cert_extension_includes.%d", schema.HashString("basic-constraints")),
						"basic-constraints"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl",
						fmt.Sprintf("cert_extension_includes.%d", schema.HashString("subject-alternative-name")),
						"subject-alternative-name"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "cert_life_span", "60"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "cert_lookup_by_ipaddr_port", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "chain", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "ciphers", "DEFAULT"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "client_cert_ca", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "crl_file", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "forward_proxy_bypass_default_action", "intercept"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "generic_alert", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "handshake_timeout", "30"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "inherit_cert_keychain", "false"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "key", "/Common/default.key"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "mod_ssl_methods", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "mode", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "peer_cert_mode", "require"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "proxy_ca_cert", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "proxy_ca_key", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "proxy_ca_passphrase", ""),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "proxy_ssl", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "proxy_ssl_passthrough", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "renegotiate_period", "60"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "renegotiate_size", "2"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "renegotiation", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "retain_certificate", "false"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "secure_renegotiation", "request"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "server_name", "testservername"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "session_mirroring", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "session_ticket", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "sni_default", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "sni_require", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "ssl_forward_proxy", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "ssl_forward_proxy_bypass", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "ssl_sign_hash", "sha1"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "strict_resume", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl",
						fmt.Sprintf("tm_options.%d", schema.HashString("dont-insert-empty-fragments")),
						"dont-insert-empty-fragments"),

					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "unclean_shutdown", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "cert_key_chain.0.name", "default"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "cert_key_chain.0.cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "cert_key_chain.0.key", "/Common/default.key"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileClientSsl_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_CLIENTSSL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(TEST_CLIENTSSL_NAME, true),
				),
				ResourceName:      TEST_CLIENTSSL_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckClientSslExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetClientSSLProfile(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("ClientSsl Profile %s was not created.", name)
		}
		if !exists && p == nil {
			return fmt.Errorf("ClientSsl Profile %s still exists.", name)
		}
		return nil
	}
}

func testCheckClientSslDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_clientssl" {
			continue
		}

		name := rs.Primary.ID
		ClientSsl, err := client.GetClientSSLProfile(name)
		if err != nil {
			return err
		}
		if ClientSsl != nil {
			return fmt.Errorf("ClientSsl Profile %s not destroyed.", name)
		}
	}
	return nil
}
