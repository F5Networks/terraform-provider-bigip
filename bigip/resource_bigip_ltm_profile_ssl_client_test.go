/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resName = "bigip_ltm_profile_client_ssl"

func TestAccBigipLtmProfileClientSsl_Default_create(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

// This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/505
func TestAccBigipLtmProfileClientSsl_UpdateName(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateName"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateName(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(fmt.Sprintf("%s-%s", instFullName, "new")),
					resource.TestCheckResourceAttr(resFullName, "name", fmt.Sprintf("%s-%s", instFullName, "new")),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

// This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/213
func TestAccBigipLtmProfileClientSsl_UpdateAuthenticate(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateAuthenticate"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "authenticate"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate", "always"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileClientSsl_UpdateAuthenticateDepth(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateAuthenticateDepth"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "authenticate_depth"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate_depth", "8"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "cache_size"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "authenticate_depth", "8"),
					resource.TestCheckResourceAttr(resFullName, "cache_size", "262100"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

// This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/213
func TestAccBigipLtmProfileClientSsl_UpdateTmoptions(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateTmoptions"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					// Newer BIG-IP defaults no longer consistently include inherited no-tlsv1.
					resource.TestCheckTypeSetElemAttr(resFullName, "tm_options.*", "dont-insert-empty-fragments"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "tm_options"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckTypeSetElemAttr(resFullName, "tm_options.*", "no-tlsv1.3"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

// This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/318
func TestAccBigipLtmProfileClientSsl_NonDefaultCert_Create(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl"
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
					testCheckClientSslExists("/Common/lbeform_INT"),
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

// This TC is added baseddded based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/449
// cert_key_chain field is going to be deprecated in near future.
func TestAccBigipLtmProfileClientSsl_CertkeyChain(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-CertkeyChain"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr(resFullName, "cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr(resFullName, "key", "/Common/default.key"),
					resource.TestCheckResourceAttr(resFullName, "chain", "/Common/ca-bundle.crt"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslCerkeychainissue449(instName),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr(resFullName, "cert", "/Common/default.crt"),
					resource.TestCheckResourceAttr(resFullName, "key", "/Common/default.key"),
					resource.TestCheckResourceAttr(resFullName, "chain", "/Common/ca-bundle.crt"),
				),
			},
		},
	})
}

// This TC is added based on ref: https://github.com/F5Networks/terraform-provider-bigip/issues/213
func TestAccBigipLtmProfileClientSsl_UpdateCachetimeout(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateCachetimeout"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "cache_timeout"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "cache_timeout", "2400"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileClientSsl_UpdateCertlifespan(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateCertlifespan"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "cert_life_span"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "cert_life_span", "40"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "handshake_timeout"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "cert_life_span", "40"),
					resource.TestCheckResourceAttr(resFullName, "handshake_timeout", "40"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileClientSsl_UpdateCipher(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateCipher"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "ciphers"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "ciphers", "AES"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "ciphers", "AES"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileClientSsl_UpdateCipherGroup(t *testing.T) {
	t.Parallel()
	var instName = "test-ClientSsl-UpdateCipherGroup"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
					testCheckClientSSLDefaultsMatchStock(resFullName),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "cipher_group"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr(resFullName, "cipher_group", "/Common/f5-aes"),
				),
			},
			{
				Config: testaccbigipltmprofileclientsslUpdateparam(instName, "cipher_group"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "partition", "Common"),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/clientssl"),
					resource.TestCheckResourceAttr(resFullName, "cipher_group", "/Common/f5-aes"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileClientSsl_import(t *testing.T) {
	var instName = "test-ClientSsl"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
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
					testCheckClientSslExists(instFullName),
				),
				ResourceName:      instFullName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckClientSslExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetClientSSLProfile(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("ClientSsl Profile %s was not created ", name)
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
			return fmt.Errorf("ClientSsl Profile %s not destroyed. ", name)
		}
	}
	return nil
}
func testaccbigipltmprofileclientsslDefaultcreate(instName string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name          = "/Common/%[2]s"
  defaults_from = "/Common/clientssl"
}
		`, resName, instName)
}

func testaccbigipltmprofileclientsslUpdateName(instName string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name          = "/Common/%[2]s-new"
  defaults_from = "/Common/clientssl"
}
		`, resName, instName)
}

func testaccbigipltmprofileclientsslCerkeychain(instName string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name         = "/Common/%[2]s"
  authenticate = "always"
  cert         = "/Common/default.crt"
  key          = "/Common/default.key"
  chain        = "/Common/ca-bundle.crt"
  passphrase   = "test123"
}
`, resName, instName)
}

func testaccbigipltmprofileclientsslCerkeychainissue449(instName string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name         = "/Common/%[2]s"
  authenticate = "once"
  cert         = "/Common/default.crt"
  key          = "/Common/default.key"
  chain        = "/Common/ca-bundle.crt"
  passphrase   = "test123"
}
`, resName, instName)
}

func testaccbigipltmprofileclientsslUpdateparam(instName, updateParam string) string {
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
	case "authenticate_depth":
		resPrefix = fmt.Sprintf(`%s
			  authenticate_depth = 8`, resPrefix)
	case "cache_size":
		resPrefix = fmt.Sprintf(`%s
			  cache_size = 262100`, resPrefix)
	case "cache_timeout":
		resPrefix = fmt.Sprintf(`%s
			  cache_timeout = 2400`, resPrefix)
	case "cert_life_span":
		resPrefix = fmt.Sprintf(`%s
			  cert_life_span = 40`, resPrefix)
	case "handshake_timeout":
		resPrefix = fmt.Sprintf(`%s
			  handshake_timeout = 40`, resPrefix)
	case "ciphers":
		// BIG-IP 21.1.x rejects setting the legacy ciphers string
		// while TLS 1.3 remains implicitly enabled (inherited from
		// /Common/clientssl, whose default tm_options no longer
		// disables it -- see this test's own Default_create case):
		// "A cipher group must be configured when TLS 1.3 is enabled".
		// Confirmed directly against a live 21.1.0.1 device that using
		// legacy ciphers now requires both cipher_group = "none" (to
		// opt out of the new cipher-group-based model) and
		// tm_options = ["no-tlsv1.3"] (to explicitly disable TLS 1.3,
		// since ciphers has no TLS 1.3 cipher suite syntax of its own).
		// Neither was required pre-21.1.x.
		resPrefix = fmt.Sprintf(`%s
			  ciphers      = "AES"
			  cipher_group = "none"
			  tm_options   = ["no-tlsv1.3"]`, resPrefix)
	case "cipher_group":
		resPrefix = fmt.Sprintf(`%s
				cipher_group = "/Common/f5-aes"`, resPrefix)
	default:
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
  name       = "/%[1]s/${lookup(var.vs_lb, "client_profile")}_${var.env}"
  cert       = "/%[1]s/${lookup(var.vs_lb, "client_profile")}_2020_${var.env}.crt"
  key        = "/%[1]s/${lookup(var.vs_lb, "client_profile")}_2020_${var.env}.key"
  depends_on = [bigip_ssl_certificate.test-cert, bigip_ssl_key.test-key]
}
	`, partition, resName, instName)
}

// TestAccBigipLtmProfileClientSsl_ChildInheritanceSurvivesUnrelatedUpdate is
// the regression test for SFDC #01262589 (ENI SpA): a child profile that
// inherits tm_options from its parent via defaults_from -- and never sets
// tm_options itself -- must keep inheriting even after an unrelated
// attribute on the child (here, authenticate, matching the customer's
// cert/key-only change) is updated. Before the fix in getClientSslConfig,
// Read would backfill the child's state with the parent's inherited
// tm_options, and the very next Update (triggered by the unrelated change)
// would re-send that backfilled value to BIG-IP as an explicit "options"
// line, permanently breaking inheritance.
func TestAccBigipLtmProfileClientSsl_ChildInheritanceSurvivesUnrelatedUpdate(t *testing.T) {
	t.Parallel()
	parentName := "test-ClientSsl-InheritParent"
	childName := "test-ClientSsl-InheritChild"
	parentFullName := fmt.Sprintf("/%s/%s", TestPartition, parentName)
	childFullName := fmt.Sprintf("/%s/%s", TestPartition, childName)
	childResFullName := fmt.Sprintf("%s.%s", resName, childName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				// Step 1: parent has tm_options set; child has defaults_from
				// pointing at the parent and does NOT set tm_options itself
				// in its own .tf. tm_options is Optional+Computed, so it is
				// expected (and correct) for the child's Terraform state to
				// reflect the inherited effective value here -- that is not
				// the bug. testCheckClientSslNoExplicitOptions instead
				// checks BIG-IP's own view of the child object directly:
				// the device-side tmOptions must still be "none", meaning
				// the child has not had an explicit options line written to
				// it and is purely inheriting from the parent.
				Config: testaccbigipltmprofileclientsslInheritanceConfig(parentName, childName, "once"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(parentFullName),
					testCheckClientSslExists(childFullName),
					resource.TestCheckResourceAttr(childResFullName, "defaults_from", parentFullName),
					testCheckClientSslNoExplicitOptions(childFullName),
				),
			},
			{
				// Step 2: update ONLY an unrelated attribute on the child
				// (authenticate), matching the customer's cert/key-only
				// update. Before the fix, this Update would have read the
				// tm_options value Step 1's Read backfilled into state
				// (the parent's inherited value) and re-sent it to BIG-IP
				// as an explicit override, permanently breaking
				// inheritance. The device-side check below is the actual
				// regression assertion.
				Config: testaccbigipltmprofileclientsslInheritanceConfig(parentName, childName, "always"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(childFullName),
					resource.TestCheckResourceAttr(childResFullName, "authenticate", "always"),
					testCheckClientSslNoExplicitOptions(childFullName),
				),
			},
		},
	})
}

func testaccbigipltmprofileclientsslInheritanceConfig(parentName, childName, childAuthenticate string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name          = "/Common/%[2]s"
  defaults_from = "/Common/clientssl"
  tm_options    = ["no-tlsv1.3"]
}

resource "%[1]s" "%[3]s" {
  name          = "/Common/%[3]s"
  defaults_from = %[1]s.%[2]s.name
  authenticate  = "%[4]s"
  depends_on    = [%[1]s.%[2]s]
}
`, resName, parentName, childName, childAuthenticate)
}

// TestAccBigipLtmProfileClientSsl_ChildInheritanceAuditOtherFields is the
// device-level regression coverage for the audit of every Optional+Computed
// field on this resource beyond tm_options (see
// resource_bigip_ltm_profile_ssl_client_inheritance_audit_test.go for the
// unit-test-level coverage of all 53 fields and the full writeup of why
// this audit concluded no code fix was needed for them). This test proves
// the same "no leak" property end-to-end against a real device for three
// representative fields explicitly named in the two historical GitHub
// issues this audit cross-references (authenticate: #450; cipher_group and
// secure_renegotiation-style enum fields: #902): a parent with distinctive
// non-default values set, a child that never configures any of these three
// fields itself, and an unrelated attribute (session_mirroring) updated on
// the child across two applies.
func TestAccBigipLtmProfileClientSsl_ChildInheritanceAuditOtherFields(t *testing.T) {
	t.Parallel()
	parentName := "test-ClientSsl-AuditParent"
	childName := "test-ClientSsl-AuditChild"
	parentFullName := fmt.Sprintf("/%s/%s", TestPartition, parentName)
	childFullName := fmt.Sprintf("/%s/%s", TestPartition, childName)
	childResFullName := fmt.Sprintf("%s.%s", resName, childName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckClientSslDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofileclientsslInheritanceAuditConfig(parentName, childName, "disabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(parentFullName),
					testCheckClientSslExists(childFullName),
					resource.TestCheckResourceAttr(childResFullName, "defaults_from", parentFullName),
					testCheckClientSslNoExplicitFields(childFullName, "authenticate", "cipher-group", "secure-renegotiation"),
				),
			},
			{
				// Update ONLY session_mirroring on the child; authenticate,
				// cipher_group, and secure_renegotiation must still be
				// purely inherited from the parent afterward.
				Config: testaccbigipltmprofileclientsslInheritanceAuditConfig(parentName, childName, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckClientSslExists(childFullName),
					resource.TestCheckResourceAttr(childResFullName, "session_mirroring", "enabled"),
					testCheckClientSslNoExplicitFields(childFullName, "authenticate", "cipher-group", "secure-renegotiation"),
				),
			},
		},
	})
}

func testaccbigipltmprofileclientsslInheritanceAuditConfig(parentName, childName, childSessionMirroring string) string {
	return fmt.Sprintf(`
resource "%[1]s" "%[2]s" {
  name                 = "/Common/%[2]s"
  defaults_from        = "/Common/clientssl"
  authenticate         = "always"
  cipher_group         = "/Common/f5-secure"
  secure_renegotiation = "require-strict"
}

resource "%[1]s" "%[3]s" {
  name              = "/Common/%[3]s"
  defaults_from     = %[1]s.%[2]s.name
  session_mirroring = "%[4]s"
  depends_on        = [%[1]s.%[2]s]
}
`, resName, parentName, childName, childSessionMirroring)
}

// testCheckClientSslNoExplicitFields asserts, via "tmsh list" (the same
// config-authoritative source used by testCheckClientSslNoExplicitOptions),
// that none of the given tmsh field names (e.g. "authenticate",
// "cipher-group") appear as a top-level explicit line on the named
// client-ssl profile -- i.e. every listed field is purely inherited from
// defaults_from rather than having been written explicitly onto this
// specific object.
//
// Caveat: this does exact first-token matching per line, but does NOT
// track brace/nesting depth, so it cannot distinguish a genuinely
// top-level field from a same-named sub-field nested inside a block (e.g.
// "cert"/"key"/"chain"/"name"/"passphrase" are all legitimate sub-keys
// inside the "cert-key-chain { default { ... } }" block that every
// client-ssl profile always has, in addition to being legitimate top-level
// field names in their own right). Do not pass any of those five names to
// this helper without first adding depth tracking -- doing so today would
// produce a false positive on the always-present cert-key-chain block
// regardless of whether the top-level field was ever set explicitly.
func testCheckClientSslNoExplicitFields(name string, tmshFields ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		resp, err := client.RunCommand(&bigip.BigipCommand{
			Command:     "run",
			UtilCmdArgs: fmt.Sprintf("-c \"tmsh list ltm profile client-ssl %s\"", name),
		})
		if err != nil {
			return err
		}
		if strings.Contains(resp.CommandResult, "not found") {
			return fmt.Errorf("ClientSsl Profile %s was not found: %s", name, resp.CommandResult)
		}
		wantFields := make(map[string]bool, len(tmshFields))
		for _, field := range tmshFields {
			wantFields[field] = true
		}
		for _, line := range strings.Split(resp.CommandResult, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || trimmed == "{" || trimmed == "}" {
				continue
			}
			// tmsh list output lines are "<field-name>[ <value>][ {]" --
			// take only the first whitespace-delimited token as the field
			// name, so e.g. searching for "authenticate" can never match
			// a line that actually starts with "authenticate-depth" (an
			// exact-token comparison rather than a substring/prefix match
			// avoids false positives when one field name happens to be a
			// prefix of another top-level tmsh field name -- see the
			// function's doc comment for the separate nested-block caveat
			// this does NOT solve).
			fieldName := trimmed
			if idx := strings.IndexAny(trimmed, " \t"); idx != -1 {
				fieldName = trimmed[:idx]
			}
			if wantFields[fieldName] {
				return fmt.Errorf("expected profile %s to have no explicit %q line (pure inheritance from parent), got tmsh output:\n%s", name, fieldName, resp.CommandResult)
			}
		}
		return nil
	}
}

// testCheckClientSslNoExplicitOptions asserts, via "tmsh list" (exactly the
// diagnostic the customer used in SFDC #01262589), that the named client-ssl
// profile has no explicit "options { ... }" line.
//
// The plain iControl REST GET response's tmOptions field is NOT a reliable
// signal for this: BIG-IP's GET always echoes the *effective* (possibly
// inherited) tmOptions value for a profile with defaults_from set, even
// when nothing was ever explicitly configured on that specific object --
// confirmed by direct device testing (creating a child with no tmOptions in
// its own create payload still returns the parent's tmOptions verbatim
// from a subsequent GET). "tmsh list" is the config-authoritative view:
// it only prints an "options { ... }" line when tmOptions has actually
// been written explicitly onto that object, which is exactly what must
// not happen here.

// explicitOptionsLineRegexp matches a tmsh-list "options { ... }" stanza
// line, tolerant of the leading indentation tmsh uses (varies with nesting
// depth) and of brace-spacing variance across TMOS versions/patches (e.g.
// "options {" vs "options{"). Anchored so it only matches "options" as a
// standalone tmsh keyword at the start of a (trimmed) line -- not merely
// wherever the substring "options" appears -- so it can't be tripped up by
// an unrelated line that happens to contain "options" elsewhere in the
// object dump (e.g. as part of a value or a differently-named key).
var explicitOptionsLineRegexp = regexp.MustCompile(`(?m)^\s*options\s*\{`)

func testCheckClientSslNoExplicitOptions(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		resp, err := client.RunCommand(&bigip.BigipCommand{
			Command:     "run",
			UtilCmdArgs: fmt.Sprintf("-c \"tmsh list ltm profile client-ssl %s\"", name),
		})
		if err != nil {
			return err
		}
		if strings.Contains(resp.CommandResult, "not found") {
			return fmt.Errorf("ClientSsl Profile %s was not found: %s", name, resp.CommandResult)
		}
		if explicitOptionsLineRegexp.MatchString(resp.CommandResult) {
			return fmt.Errorf("expected profile %s to have no explicit \"options\" line (pure inheritance from parent), got tmsh output:\n%s", name, resp.CommandResult)
		}
		return nil
	}
}
