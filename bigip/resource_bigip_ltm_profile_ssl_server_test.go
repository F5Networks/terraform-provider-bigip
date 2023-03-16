/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resNameserver = "bigip_ltm_profile_server_ssl"

func TestAccBigipLtmProfileServerSsl_Default_create(t *testing.T) {
	t.Parallel()
	var instName = "test-ServerSsl"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resNameserver, instName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
					resource.TestCheckResourceAttr(resFullName, "alert_timeout", "indefinite"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckResourceAttr(resFullName, "authenticate_depth", "9"),
					resource.TestCheckResourceAttr(resFullName, "cache_size", "262144"),
					resource.TestCheckResourceAttr(resFullName, "ca_file", "none"),
					resource.TestCheckResourceAttr(resFullName, "cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr(resFullName, "key", "/Common/default.key"),
					resource.TestCheckResourceAttr(resFullName, "chain", "none"),
					resource.TestCheckResourceAttr(resFullName, "ciphers", "DEFAULT"),
					resource.TestCheckResourceAttr(resFullName, "expire_cert_response_control", "drop"),
					resource.TestCheckResourceAttr(resFullName, "handshake_timeout", "10"),
					resource.TestCheckResourceAttr(resFullName, "mod_ssl_methods", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "mode", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "peer_cert_mode", "ignore"),
					resource.TestCheckResourceAttr(resFullName, "proxy_ssl", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "renegotiate_period", "indefinite"),
					resource.TestCheckResourceAttr(resFullName, "renegotiate_size", "indefinite"),
					resource.TestCheckResourceAttr(resFullName, "renegotiation", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "retain_certificate", "true"),
					resource.TestCheckResourceAttr(resFullName, "secure_renegotiation", "require-strict"),
					resource.TestCheckResourceAttr(resFullName, "server_name", "none"),
					resource.TestCheckResourceAttr(resFullName, "session_mirroring", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "session_ticket", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "sni_default", "false"),
					resource.TestCheckResourceAttr(resFullName, "sni_require", "false"),
					resource.TestCheckResourceAttr(resFullName, "ssl_forward_proxy", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "ssl_forward_proxy_bypass", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "ssl_sign_hash", "any"),
					resource.TestCheckResourceAttr(resFullName, "strict_resume", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "unclean_shutdown", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "untrusted_cert_response_control", "drop"),
				),
			},
		},
	})
}

// This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/213
func TestAccBigipLtmProfileServerSsl_UpdateAuthenticate(t *testing.T) {
	t.Parallel()
	var instName = "test-ServerSsl-UpdateAuthenticate"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resNameserver, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
				),
			},
			{
				Config: testAccBigipLtmProfileServerSsl_UpdateParam(instName, "authenticate"),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "always"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
				),
			},
		},
	})
}

// This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/213
func TestAccBigipLtmProfileServerSsl_UpdateTmoptions(t *testing.T) {
	t.Parallel()
	var instName = "test-ServerSsl-UpdateTmoptions"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resNameserver, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckTypeSetElemAttr(resFullName, "tm_options.*", "dont-insert-empty-fragments"),
					resource.TestCheckTypeSetElemAttr(resFullName, "tm_options.*", "no-tlsv1.3"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
				),
			},
			{
				Config: testAccBigipLtmProfileServerSsl_UpdateParam(instName, "tm_options"),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckTypeSetElemAttr(resFullName, "tm_options.*", "no-tlsv1.3"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileServerSsl_UpdateCipherGroup(t *testing.T) {
	t.Parallel()
	var instName = "test-ServerSsl-UpdateCipherGroup"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resNameserver, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
					resource.TestCheckResourceAttr(resFullName, "cipher_group", "none"),
				),
			},
			{
				Config: testAccBigipLtmProfileServerSsl_UpdateParam(instName, "cipher_group"),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/serverssl"),
					resource.TestCheckResourceAttr(resFullName, "cipher_group", "/Common/f5-aes"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileServerSsl_import(t *testing.T) {
	var instName = "test-ServerSsl"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckServerSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileserversslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckServerSslExists(instFullName),
				),
				ResourceName:      instFullName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckServerSslExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetServerSSLProfile(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("ServerSsl Profile %s was not created.", name)
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
func testaccbigipltmprofileserversslDefaultcreate(instName string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name = "/Common/%[2]s"
  //defaults_from = "/Common/serverssl"
}
		`, resNameserver, instName)
}

func testAccBigipLtmProfileServerSsl_UpdateParam(instName, updateParam string) string {
	resPrefix := fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			  name = "/Common/%[2]s"
			  defaults_from = "/Common/serverssl"`, resNameserver, instName)
	switch updateParam {
	case "authenticate":
		resPrefix = fmt.Sprintf(`%s
			  authenticate = "always"`, resPrefix)
	case "tm_options":
		resPrefix = fmt.Sprintf(`%s
			  tm_options = ["no-tlsv1.3"]`, resPrefix)
	case "cipher_group":
		resPrefix = fmt.Sprintf(`%s
			  cipher_group = "/Common/f5-aes"`, resPrefix)
	}
	return fmt.Sprintf(`%s
		}`, resPrefix)
}
