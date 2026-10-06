package bigip

import (
	"fmt"
	"log"
	"os"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testResourceSSLKeyCert = `
resource "bigip_ssl_key_cert" "testkeycert" {
  partition   = "Common"
  key_name    = "ssl-test-key"
	key_content = "${file("` + sslKeyCertFolder + `/../examples/serverkey.key")}"
  cert_name    = "ssl-test-cert"
	cert_content = "${file("` + sslKeyCertFolder + `/../examples/servercert.crt")}"
}
`

var sslProfileCertKey = `
resource "bigip_ssl_key_cert" "testkeycert" {
  partition   = "Common"
  key_name    = "ssl-test-key"
	key_content = "${file("` + sslKeyCertFolder + `/../examples/%s")}"
  cert_name    = "ssl-test-cert"
	cert_content = "${file("` + sslKeyCertFolder + `/../examples/%s")}"
}

resource "bigip_ltm_profile_server_ssl" "test-ServerSsl" {
  name          = "/Common/test-ServerSsl"
  defaults_from = "/Common/serverssl"
  authenticate  = "always"
  cipher_group  = "/Common/f5-default"
  cert          = "/Common/ssl-test-cert"
  key           = "/Common/ssl-test-key"

  depends_on = [
	bigip_ssl_key_cert.testkeycert
  ]
}
`

// sslProfileCertKeyOCSP embeds leafCertPEM/leafKeyPEM directly (rather
// than file()-ing static fixtures) since issuer_cert = "/Common/MyCA"
// only validates against a certificate actually signed by whatever CA key
// was installed as /Common/MyCA -- see installMyCA's doc comment in
// provider_test.go. leafCertPEM/leafKeyPEM must come from the same
// generateCAAndLeafCert call whose CA half gets passed to
// testAcctPreCheckOCSP for that same test.
func sslProfileCertKeyOCSP(leafCertPEM, leafKeyPEM []byte) string {
	return fmt.Sprintf(`
resource "bigip_ssl_key_cert" "testkeycert" {
  partition            = "Common"
  key_name             = "ssl-test-key"
  key_content          = <<-EOT
%s
EOT
  cert_name            = "ssl-test-cert"
  cert_content         = <<-EOT
%s
EOT
  cert_monitoring_type = "ocsp"
  issuer_cert          = "/Common/MyCA"
  cert_ocsp            = "/Common/testocsp1"
}
`, leafKeyPEM, leafCertPEM)
}

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
	crt1Content, _ := os.ReadFile(sslKeyCertFolder + `/../examples/` + "servercert.crt")
	key1Content, _ := os.ReadFile(sslKeyCertFolder + `/../examples/` + "serverkey.key")
	crt2Content, _ := os.ReadFile(sslKeyCertFolder + `/../examples/` + "servercert2.crt")
	key2Content, _ := os.ReadFile(sslKeyCertFolder + `/../examples/` + "serverkey2.key")

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
	// Generated once, up front: the CA half is installed as /Common/MyCA
	// by testAcctPreCheckOCSP, and the leaf half (cryptographically signed
	// by that same CA key) is embedded directly in sslProfileCertKeyOCSP
	// below. See installMyCA's doc comment in provider_test.go for why
	// these must be generated together as a matched pair rather than
	// using static fixture files.
	caCertPEM, leafCertPEM, leafKeyPEM := generateCAAndLeafCert(t, "MyCA", "ocsp-leaf.f5.com")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheckOCSP(caCertPEM)(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: sslProfileCertKeyOCSP(leafCertPEM, leafKeyPEM),
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

// Write-only attribute tests for SSL key-cert resource

var sslKeyCertFolder, _ = os.Getwd()

var TestSSLKeyCertResourceWriteOnly = `
resource "bigip_ssl_key_cert" "testkeycert-wo" {
	key_name = "ssl-test-key-wo"
	key_content_wo = "${file("` + sslKeyCertFolder + `/../examples/serverkey.key")}" 
	cert_name = "ssl-test-cert-wo"
	cert_content_wo = "${file("` + sslKeyCertFolder + `/../examples/servercert.crt") }"
	partition = "Common"
}
`

var TestSSLKeyCertResourceWriteOnlyWithVersion = `
resource "bigip_ssl_key_cert" "testkeycert-wo" {
	key_name = "ssl-test-key-wo"
	key_content_wo = "${file("` + sslKeyCertFolder + `/../examples/serverkey.key") }"
	key_content_wo_version = 1
	cert_name = "ssl-test-cert-wo"
	cert_content_wo = "${file("` + sslKeyCertFolder + `/../examples/servercert.crt") }"
	cert_content_wo_version = 1
	partition = "Common"
}
`

// Test write-only content_wo attributes for key-cert
func TestAccBigipSSLKeyCertWriteOnlyContent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslKeyCertDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestSSLKeyCertResourceWriteOnly,
				Check: resource.ComposeTestCheckFunc(
					testChecksslKeyCertExists("ssl-test-key-wo", "ssl-test-cert-wo", true),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert-wo", "key_name", "ssl-test-key-wo"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert-wo", "cert_name", "ssl-test-cert-wo"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert-wo", "partition", "Common"),
					// Verify write-only attributes are not stored in state
					testCheckResourceAttrNotSetKeyCert("bigip_ssl_key_cert.testkeycert-wo", "key_content_wo"),
					testCheckResourceAttrNotSetKeyCert("bigip_ssl_key_cert.testkeycert-wo", "cert_content_wo"),
				),
			},
		},
	})
}

// Test content_wo_version change triggers re-upload for key-cert
func TestAccBigipSSLKeyCertWriteOnlyVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslKeyCertDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestSSLKeyCertResourceWriteOnly,
				Check: resource.ComposeTestCheckFunc(
					testChecksslKeyCertExists("ssl-test-key-wo", "ssl-test-cert-wo", true),
					// key_content_wo_version/cert_content_wo_version are
					// write-only: the SDK always nulls them out of state, so
					// they can never be asserted via TestCheckResourceAttr.
					// Verify instead that they stay absent from state on both
					// the initial apply and after bumping them below, which
					// still exercises the HasChange(..._wo_version) re-upload
					// path.
					testCheckResourceAttrNotSetKeyCert("bigip_ssl_key_cert.testkeycert-wo", "key_content_wo_version"),
					testCheckResourceAttrNotSetKeyCert("bigip_ssl_key_cert.testkeycert-wo", "cert_content_wo_version"),
				),
			},
			{
				Config: TestSSLKeyCertResourceWriteOnlyWithVersion,
				Check: resource.ComposeTestCheckFunc(
					testChecksslKeyCertExists("ssl-test-key-wo", "ssl-test-cert-wo", true),
					testCheckResourceAttrNotSetKeyCert("bigip_ssl_key_cert.testkeycert-wo", "key_content_wo_version"),
					testCheckResourceAttrNotSetKeyCert("bigip_ssl_key_cert.testkeycert-wo", "cert_content_wo_version"),
				),
			},
		},
	})
}

// Helper function to check if an attribute is not set in state (for key-cert resource)
func testCheckResourceAttrNotSetKeyCert(resourceName string, attributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		val, exists := rs.Primary.Attributes[attributeName]
		if exists && val != "" {
			return fmt.Errorf("expected %s to not be set, but got value: %s", attributeName, val)
		}
		return nil
	}
}

// Check that the ssl key-cert exists (by resource state entries)
func testChecksslKeyCertExists(keyName, certName string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "bigip_ssl_key_cert" {
				continue
			}
			if rs.Primary.Attributes["key_name"] == keyName || rs.Primary.Attributes["cert_name"] == certName {
				if exists {
					return nil
				}
				return fmt.Errorf("ssl key-cert %s/%s unexpectedly exists", keyName, certName)
			}
		}
		if exists {
			return fmt.Errorf("ssl key-cert %s/%s not found in state", keyName, certName)
		}
		return nil
	}
}

func testChecksslKeyCertDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ssl_key_cert" {
			continue
		}
		partition := rs.Primary.Attributes["partition"]
		keyName := fqdn(partition, rs.Primary.Attributes["key_name"])
		certName := fqdn(partition, rs.Primary.Attributes["cert_name"])

		key, err := client.GetKey(keyName)
		if err != nil {
			return err
		}
		if key != nil {
			return fmt.Errorf("ssl key %s not destroyed", keyName)
		}

		cert, err := client.GetCertificate(certName)
		if err != nil {
			return err
		}
		if cert != nil {
			return fmt.Errorf("ssl certificate %s not destroyed", certName)
		}
	}
	return nil
}
