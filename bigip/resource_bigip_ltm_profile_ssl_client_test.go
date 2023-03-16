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
					resource.TestCheckTypeSetElemAttr(resFullName, "tm_options.*", "dont-insert-empty-fragments"),
					resource.TestCheckTypeSetElemAttr(resFullName, "tm_options.*", "no-tlsv1.3"),
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
					resource.TestCheckResourceAttr(resFullName, "cipher_group", "none"),
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
		resPrefix = fmt.Sprintf(`%s
			  ciphers = "AES"`, resPrefix)
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
