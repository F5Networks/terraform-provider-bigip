/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_SERVERSSL_NAME = fmt.Sprintf("/%s/test-ServerSsl", TEST_PARTITION)

var TEST_SERVERSSL_RESOURCE = `
resource "bigip_ltm_profile_server_ssl" "test-ServerSsl" {
  name = "/Common/test-ServerSsl"
  partition = "Common"
  defaults_from = "/Common/serverssl"
  authenticate = "always"
  ciphers = "DEFAULT"
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
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "alert_timeout", "indefinite"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "authenticate", "always"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "authenticate_depth", "9"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "cache_size", "262144"),
					//resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "cache_timeout", "3600"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ca_file", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "cert", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "chain", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ciphers", "DEFAULT"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "expire_cert_response_control", "drop"),
					//resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "generic_alert", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "handshake_timeout", "10"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "key", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "mod_ssl_methods", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "mode", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "peer_cert_mode", "ignore"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "proxy_ssl", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "renegotiate_period", "indefinite"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "renegotiate_size", "indefinite"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "renegotiation", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "retain_certificate", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "secure_renegotiation", "require-strict"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "server_name", "none"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "session_mirroring", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "session_ticket", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "sni_default", "false"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "sni_require", "false"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ssl_forward_proxy", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ssl_forward_proxy_bypass", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "ssl_sign_hash", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "strict_resume", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "unclean_shutdown", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_server_ssl.test-ServerSsl", "untrusted_cert_response_control", "drop"),
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
