package bigip

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testResourceSSLKeyCert = `
resource "bigip_ssl_key_cert" "testkeycert" {
  partition   = "Common"
  key_name    = "ssl-test-key"
  key_content = "${file("` + folder + `/../examples/serverkey.key")}"
  cert_name    = "ssl-test-cert"
  cert_content = "${file("` + folder + `/../examples/servercert.crt")}"
}
`

var sslProfileCertKey = `
resource "bigip_ssl_key_cert" "testkeycert" {
  partition   = "Common"
  key_name    = "ssl-test-key"
  key_content = "${file("` + folder + `/../examples/%s")}"
  cert_name    = "ssl-test-cert"
  cert_content = "${file("` + folder + `/../examples/%s")}"
}

resource "bigip_ltm_profile_server_ssl" "test-ServerSsl" {
  name          = "/Common/test-ServerSsl"
  defaults_from = "/Common/serverssl"
  authenticate  = "always"
  ciphers       = "DEFAULT"
  cert          = "/Common/ssl-test-cert"
  key           = "/Common/ssl-test-key"

  depends_on = [
	bigip_ssl_key_cert.testkeycert
  ]
}
`

var sslProfileCertKeyOCSP = `
resource "bigip_ssl_key_cert" "testkeycert" {
  partition            = "Common"
  key_name             = "ssl-test-key"
  key_content          = "${file("` + folder + `/../examples/mycertocspv2.pem")}"
  cert_name            = "ssl-test-cert"
  cert_content         = "${file("` + folder + `/../examples/mycertocspv2.crt")}"
  cert_monitoring_type = "ocsp"
  issuer_cert          = "/Common/MyCA"
  cert_ocsp            = "/Common/testocsp1"
}
`

func TestAccBigipSSLCertKeyCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		// CheckDestroy:
		Steps: []resource.TestStep{
			{
				Config: testResourceSSLKeyCert,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "key_name", "ssl-test-key"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_name", "ssl-test-cert"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "partition", "Common"),
				),
				Destroy: false,
			},
			{
				Config: testResourceSSLKeyCert,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "key_name", "ssl-test-key"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_name", "ssl-test-cert"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "partition", "Common"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccBigipSSLCertKeyCreateCertKeyProfile(t *testing.T) {
	create := fmt.Sprintf(sslProfileCertKey, "serverkey.key", "servercert.crt")
	modify := fmt.Sprintf(sslProfileCertKey, "serverkey2.key", "servercert2.crt")
	crt1Content, _ := os.ReadFile(folder + `/../examples/` + "servercert.crt")
	key1Content, _ := os.ReadFile(folder + `/../examples/` + "serverkey.key")
	crt2Content, _ := os.ReadFile(folder + `/../examples/` + "servercert2.crt")
	key2Content, _ := os.ReadFile(folder + `/../examples/` + "serverkey2.key")

	log.Println(create)
	log.Println(modify)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: create,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "key_name", "ssl-test-key"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_name", "ssl-test-cert"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "key_content", string(key1Content)),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_content", string(crt1Content)),
				),
				Destroy: false,
			},
			{
				Config: modify,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "key_name", "ssl-test-key"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_name", "ssl-test-cert"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "key_content", string(key2Content)),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_content", string(crt2Content)),
				),
			},
		},
	})
}

func TestAccBigipSSLCertKeyCreateCertKeyProfileOCSP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: sslProfileCertKeyOCSP,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "key_name", "ssl-test-key"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_name", "ssl-test-cert"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_monitoring_type", "ocsp"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "issuer_cert", "/Common/MyCA"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_ocsp", "/Common/testocsp1"),
				),
			},
		},
	})
}
