package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_SERVERSSL_NAME = fmt.Sprintf("/%s/test-ServerSsl", TEST_PARTITION)

var TEST_SERVERSSL_RESOURCE = `
resource "bigip_ltm_profile_server_ssl" "test-ServerSsl" {
	name                            = "/Common/test-ServerSsl"
	partition                       = "Common"
	defaults_from                   = "/Common/serverssl"
	alert_timeout                   = "60"
	authenticate                    = "always"
	authenticate_depth              = 7
	cache_size                      = 131072
	cache_timeout					= 2400
	ca_file							= "none"
	cert 							= "none"
	chain                           = "none"
	ciphers                         = "ALL"
	expire_cert_response_control    = "drop"
	full_path						= "/Common/test-ServerSsl"
	generic_alert					= "disabled"
	handshake_timeout               = "30"
	key                             = "none"
	mod_ssl_methods                 = "enabled"
	mode                            = "disabled"
	passphrase						= ""
	peer_cert_mode					= "require"
	proxy_ssl                       = "enabled"
	renegotiate_period              = "30"
	renegotiate_size                = "2"
	renegotiation                   = "enabled"
	retain_certificate              = "false"
	secure_renegotiation            = "request"
	server_name                     = "testservername"
	session_mirroring               = "disabled"
	session_ticket                  = "disabled"
	sni_default                     = "true"
	sni_require                     = "true"
	ssl_forward_proxy               = "disabled"
	ssl_forward_proxy_bypass        = "disabled"
	ssl_sign_hash                   = "sha1"
	strict_resume                   = "enabled"
    tm_options                      = [
	  "dont-insert-empty-fragments",
	]
	unclean_shutdown                = "disabled"
	untrusted_cert_response_control = "ignore"
}
`

func TestAccBigipLtmProfileServerSsl_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SERVERSSL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(TEST_SERVERSSL_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "name", "/Common/test-ServerSsl"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "defaults_from", "/Common/serverssl"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "alert_timeout", "60"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "authenticate", "always"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "authenticate_depth", "7"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "cache_size", "131072"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "cache_timeout", "2400"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ca_file", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "cert", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "chain", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ciphers", "ALL"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "expire_cert_response_control", "drop"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "generic_alert", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "handshake_timeout", "30"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "key", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "mod_ssl_methods", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "mode", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "passphrase", ""),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "peer_cert_mode", "require"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "proxy_ssl", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "renegotiate_period", "30"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "renegotiate_size", "2"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "renegotiation", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "retain_certificate", "false"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "secure_renegotiation", "request"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "server_name", "testservername"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "session_mirroring", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "session_ticket", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "sni_default", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "sni_require", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ssl_forward_proxy", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ssl_forward_proxy_bypass", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ssl_sign_hash", "sha1"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "strict_resume", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl",
						fmt.Sprintf("tm_options.%d", schema.HashString("dont-insert-empty-fragments")),
						"dont-insert-empty-fragments"),

					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "unclean_shutdown", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "untrusted_cert_response_control", "ignore"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileServerSsl_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SERVERSSL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(TEST_SERVERSSL_NAME, true),
				),
				ResourceName:      TEST_SERVERSSL_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckServerSslExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetServerSSLProfile(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("ServerSsl Profile %s was not created.", name)
		}
		if !exists && p == nil {
			return fmt.Errorf("ServerSsl Profile %s still exists.", name)
		}
		return nil
	}
}

func testCheckServerSslDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_serverssl" {
			continue
		}

		name := rs.Primary.ID
		ServerSsl, err := client.GetServerSSLProfile(name)
		if err != nil {
			return err
		}
		if ServerSsl != nil {
			return fmt.Errorf("ServerSsl Profile %s not destroyed.", name)
		}
	}
	return nil
}
