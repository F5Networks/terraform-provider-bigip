package bigip

import (
	"fmt"
	"testing"

	"github.com/pirotrav/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_ClientSsl_NAME = fmt.Sprintf("/%s/test-ClientSsl", TEST_PARTITION)

var TEST_ClientSsl_RESOURCE = `
resource "bigip_ltm_profile_client_ssl" "profile_mutualssl2" {
	alert_timeout                       = "indefinite"
	allow_non_ssl                       = "disabled"
	authenticate                        = "once"
	authenticate_depth                  = 9
	ca_file                             = "none"
	cache_size                          = 262144
	cache_timeout                       = 3600
	cert                                = "/Common/default.crt"
	cert_extension_includes             = [
		"basic-constraints",
		"subject-alternative-name",
	  ]
	cert_life_span                      = 30
	cert_lookup_by_ipaddr_port          = "disabled"
	chain                               = "none"
	ciphers                             = "DEFAULT"
	client_cert_ca                      = "none"
	crl_file                            = "none"
	defaults_from                       = "/Common/clientssl"
	forward_proxy_bypass_default_action = "intercept"
	generic_alert                       = "enabled"
	handshake_timeout                   = "10"
	id                                  = "terraform_test_client"
	inherit_cert_keychain               = "false"
	key                                 = "/Common/default.key"
	mod_ssl_methods                     = "disabled"
	mode                                = "enabled"
	name                                = "terraform_test_client"
	partition                           = "Common"
	peer_cert_mode                      = "ignore"
	proxy_ca_cert                       = "none"
	proxy_ca_key                        = "none"
	proxy_ssl                           = "disabled"
	proxy_ssl_passthrough               = "disabled"
	renegotiate_period                  = "indefinite"
	renegotiate_size                    = "indefinite"
	renegotiation                       = "enabled"
	retain_certificate                  = "true"
	secure_renegotiation                = "require"
	server_name                         = "none"
	session_mirroring                   = "disabled"
	session_ticket                      = "disabled"
	sni_default                         = "false"
	sni_require                         = "false"
	ssl_forward_proxy                   = "disabled"
	ssl_forward_proxy_bypass            = "disabled"
	ssl_sign_hash                       = "any"
	strict_resume                       = "disabled"
	tm_options                          = [
		"dont-insert-empty-fragments",
	  ]
	unclean_shutdown                    = "enabled"

	  cert_key_chain {
		  cert = "/Common/default.crt"
		  key  = "/Common/default.key"
		  name = "default"
	  }
  }
`

func TestAccBigipLtmProfileClientSsl_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ClientSsl_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(TEST_ClientSsl_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "name", "/Common/test-ClientSsl"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-ClientSsl", "defaults_from", "/Common/ClientSsl"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "alert_timeout", "indefinite"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "allow_non_ssl", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "authenticate", "once"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "authenticate_depth", "9"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "ca_file", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "cache_size", "262144"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "cache_timeout", "3600"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "cert_extension_includes", ["basic-constraints","subject-alternative-name",]),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "cert_life_span", "30"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "cert_lookup_by_ipaddr_port", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "chain", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "ciphers", "DEFAULT"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "client_cert_ca", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "crl_file", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "forward_proxy_bypass_default_action", "intercept"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "generic_alert", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "handshake_timeout", "10"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "id", "terraform_test_client"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "inherit_cert_keychain", "FALSE"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "key", "/Common/default.key"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "mod_ssl_methods", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "mode", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "name", "terraform_test_client"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "peer_cert_mode", "ignore"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "proxy_ca_cert", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "proxy_ca_key", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "proxy_ssl", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "proxy_ssl_passthrough", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "renegotiate_period", "indefinite"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "renegotiate_size", "indefinite"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "renegotiation", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "retain_certificate", "TRUE"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "secure_renegotiation", "require"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "server_name", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "session_mirroring", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "session_ticket", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "sni_default", "FALSE"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "sni_require", "FALSE"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "ssl_forward_proxy", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "ssl_forward_proxy_bypass", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "ssl_sign_hash", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "strict_resume", "disabled"
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "m_options", "["dont-insert-empty-fragments",]"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "unclean_shutdown", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_client_ssl.test-Client", "cert_key_chain", {cert="/Common/default.crt"
						key  = "/Common/default.key"
						name = "default"
					}),
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
		CheckDestroy: testCheckClientSslsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ClientSsl_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(TEST_ClientSsl_NAME, true),
				),
				ResourceName:      TEST_ClientSsl_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckClientSslExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetClientSsl(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("ClientSsl %s was not created.", name)
		}
		if !exists && p == nil {
			return fmt.Errorf("ClientSsl %s still exists.", name)
		}
		return nil
	}
}

func testCheckClientSslsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_ClientSsl" {
			continue
		}

		name := rs.Primary.ID
		Client, err := client.GetClientSsl(name)
		if err != nil {
			return err
		}
		if ClientSsl != nil {
			return fmt.Errorf("ClientSsl %s not destroyed.", name)
		}
	}
	return nil
}
