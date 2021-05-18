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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var resName = "bigip_ltm_profile_client_ssl"

func TestAccBigipLtmProfileClientSsl_Default_create(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resName, instName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileclientsslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName, true),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr(resFullName, "alert_timeout", "indefinite"),
					resource.TestCheckResourceAttr(resFullName, "allow_non_ssl", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckResourceAttr(resFullName, "authenticate_depth", "9"),
					resource.TestCheckResourceAttr(resFullName, "ca_file", "none"),
					resource.TestCheckResourceAttr(resFullName, "cache_size", "262144"),
					resource.TestCheckResourceAttr(resFullName, "cache_timeout", "3600"),
					resource.TestCheckResourceAttr(resFullName, "cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr(resFullName,
						fmt.Sprintf("cert_extension_includes.%d", schema.HashString("basic-constraints")),
						"basic-constraints"),
					resource.TestCheckResourceAttr(resFullName,
						fmt.Sprintf("cert_extension_includes.%d", schema.HashString("subject-alternative-name")),
						"subject-alternative-name"),
					resource.TestCheckResourceAttr(resFullName, "cert_life_span", "30"),
					resource.TestCheckResourceAttr(resFullName, "cert_lookup_by_ipaddr_port", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "chain", "none"),
					resource.TestCheckResourceAttr(resFullName, "ciphers", "DEFAULT"),
					resource.TestCheckResourceAttr(resFullName, "client_cert_ca", "none"),
					resource.TestCheckResourceAttr(resFullName, "crl_file", "none"),
					resource.TestCheckResourceAttr(resFullName, "forward_proxy_bypass_default_action", "intercept"),
					resource.TestCheckResourceAttr(resFullName, "generic_alert", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "handshake_timeout", "10"),
					resource.TestCheckResourceAttr(resFullName, "inherit_cert_keychain", "true"),
					resource.TestCheckResourceAttr(resFullName, "key", "/Common/default.key"),
					resource.TestCheckResourceAttr(resFullName, "mod_ssl_methods", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "mode", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "peer_cert_mode", "ignore"),
					resource.TestCheckResourceAttr(resFullName, "proxy_ssl", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "proxy_ssl_passthrough", "disabled"),
					resource.TestCheckResourceAttr(resFullName, "renegotiate_period", "indefinite"),
					resource.TestCheckResourceAttr(resFullName, "renegotiate_size", "indefinite"),
					resource.TestCheckResourceAttr(resFullName, "renegotiation", "enabled"),
					resource.TestCheckResourceAttr(resFullName, "retain_certificate", "true"),
					resource.TestCheckResourceAttr(resFullName, "secure_renegotiation", "require"),
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
				),
			},
		},
	})
}

//
//This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/213
//
func TestAccBigipLtmProfileClientSsl_UpdateAuthenticate(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateAuthenticate"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileclientsslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName, true),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testAccBigipLtmProfileClientSsl_UpdateParam(instName, "authenticate"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName, true),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "always"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

//
//This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/213
//
func TestAccBigipLtmProfileClientSsl_UpdateTmoptions(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateTmoptions"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileclientsslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName, true),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("tm_options.%d", schema.HashString("dont-insert-empty-fragments")), "dont-insert-empty-fragments"),
					resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("tm_options.%d", schema.HashString("no-tlsv1.3")), "no-tlsv1.3"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testAccBigipLtmProfileClientSsl_UpdateParam(instName, "tm_options"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName, true),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("tm_options.%d", schema.HashString("no-tlsv1.3")), "no-tlsv1.3"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

//
//This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/318
//
func TestAccBigipLtmProfileClientSsl_NonDefaultCert_Create(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl"
	//var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileclientsslNondefaultcertconfigbasic("Common", instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists("/Common/lbeform_INT", true),
					resource.TestCheckResourceAttr(resFullName, "name", "/Common/lbeform_INT"),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr(resFullName, "cert", "/Common/lbeform_2020_INT.crt"),
					resource.TestCheckResourceAttr(resFullName, "key", "/Common/lbeform_2020_INT.key"),
				),
			},
		},
	})
}

//
//This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/449
//
func TestAccBigipLtmProfileClientSsl_CertkeyChain(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-CertkeyChain"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileclientsslCerkeychain(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName, true),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr(resFullName, "cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr(resFullName, "key", "/Common/default.key"),
					//resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("cert_key_chain.0.%d", schema.HashString("cert")), "default.crt"),
					//resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("cert_key_chain.0.%d", schema.HashString("key")), "default.key"),
					//resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("cert_key_chain.0.%d", schema.HashString("name")), "default"),
					//resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("cert_key_chain.0.%d", schema.HashString("chain")), "ca-bundle.crt"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslCerkeychainissue449(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName, true),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "once"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr(resFullName, "cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr(resFullName, "key", "/Common/default.key"),
					//resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("cert_key_chain.0.%d", schema.HashString("cert")), "default.crt"),
					//resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("cert_key_chain.0.%d", schema.HashString("key")), "default.key"),
					//resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("cert_key_chain.0.%d", schema.HashString("name")), "default"),
					//resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("cert_key_chain.0.%d", schema.HashString("chain")), "ca-bundle.crt"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileClientSsl_import(t *testing.T) {
	var instName = "test-ClientSsl"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	//resFullName := fmt.Sprintf("%s.%s", resName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				//Config: TestClientsslResource,
				Config: testaccbigipltmprofileclientsslDefaultcreate(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName, true),
				),
				ResourceName:      instFullName,
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
func testaccbigipltmprofileclientsslDefaultcreate(instName string) string {
	return fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			  name = "/Common/%[2]s"
			  defaults_from = "/Common/clientssl"
		}`, resName, instName)
}

func testaccbigipltmprofileclientsslCerkeychain(instName string) string {
	return fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			name = "/Common/%[2]s"
			authenticate = "always"
  			cert_key_chain {
    			name = "default"
    			cert = "default.crt"
				key  = "default.key"
    			chain = "ca-bundle.crt"
  			}
		}`, resName, instName)
}

func testaccbigipltmprofileclientsslCerkeychainissue449(instName string) string {
	return fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			name = "/Common/%[2]s"
			authenticate = "once"
  			cert_key_chain {
    			name = "default"
    			cert = "default.crt"
				key  = "default.key"
    			chain = "ca-bundle.crt"
  			}
		}`, resName, instName)
}

//func testAccBigipLtmProfileClientSsl_UpdateAuthenticate(instName string) string {
//	return fmt.Sprintf(`
//		resource "%[1]s" "%[2]s" {
//			  name = "/Common/%[2]s"
//			  defaults_from = "/Common/clientssl"
//			  authenticate = "always"
//			  //ciphers = "DEFAULT"
//		}`, resName, instName)
//}

func testAccBigipLtmProfileClientSsl_UpdateParam(instName, updateParam string) string {
	resPrefix := fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			  name = "/Common/%[2]s"
			  defaults_from = "/Common/clientssl"`, resName, instName)
	switch updateParam {
	case "authenticate":
		resPrefix = fmt.Sprintf(`%s
			  authenticate = "always"`, resPrefix)
	case "tm_options":
		resPrefix = fmt.Sprintf(`%s
			  tm_options = ["no-tlsv1.3"]`, resPrefix)
	}
	return fmt.Sprintf(`%s
		}`, resPrefix)
}

func testaccbigipltmprofileclientsslNondefaultcertconfigbasic(partition, instName string) string {
	return fmt.Sprintf(`
	variable vs_lb {
		type = object({
		client_profile = string
	})
	default = { "client_profile" = "lbeform" }
	}
	variable env {
		type    = string
		default = "INT"
	}
	resource "bigip_ssl_certificate" "test-cert" {
		name      = "${lookup(var.vs_lb, "client_profile")}_2020_${var.env}.crt"
		content   = file("`+dir+`/../examples/servercert.crt")
		partition = "%[1]s"
	}
	resource "bigip_ssl_key" "test-key" {
		name      = "${lookup(var.vs_lb, "client_profile")}_2020_${var.env}.key"
		content   = file("`+dir+`/../examples/serverkey.key")
		partition = "%[1]s"
	}
	resource "%[2]s" "%[3]s" {
		name = "/%[1]s/${lookup(var.vs_lb, "client_profile")}_${var.env}"
		cert = "/%[1]s/${lookup(var.vs_lb, "client_profile")}_2020_${var.env}.crt"
		key  = "/%[1]s/${lookup(var.vs_lb, "client_profile")}_2020_${var.env}.key"
		depends_on = [bigip_ssl_certificate.test-cert, bigip_ssl_key.test-key]
	}`, partition, resName, instName)
}
