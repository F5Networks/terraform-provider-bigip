/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"os"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var folder, _ = os.Getwd()
var SslcertificateName = "servercert.crt"
var TestSslcertificateName = fmt.Sprintf("/%s/%s", TestPartition, SslcertificateName)

var TestSslCertificateResource = `
resource "bigip_ssl_certificate" "test-cert" {
        name = "` + SslcertificateName + `"
        content = "${file("` + folder + `/../examples/servercert.crt")}"
        partition = "Common"
}
`

// testSslCertOCSPResource embeds leafCertPEM directly (rather than
// file()-ing a static fixture) since issuer_cert = "/Common/MyCA" only
// validates against a certificate actually signed by whatever CA key was
// installed as /Common/MyCA -- see installMyCA's doc comment. leafCertPEM
// must come from the same generateCAAndLeafCert call whose CA half gets
// passed to testAcctPreCheckOCSP for that same test.
func testSslCertOCSPResource(leafCertPEM []byte) string {
	return fmt.Sprintf(`
resource "bigip_ssl_certificate" "ssl-test-certificate-tc1" {
  name            = "test-certificate"
  content         = <<-EOT
%s
EOT
  partition       = "Common"
  monitoring_type = "ocsp"
  issuer_cert     = "/Common/MyCA"
  ocsp            = "/Common/testocsp1"
}
`, leafCertPEM)
}

func TestAccBigipSslCertificateImportToBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslcertificateDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestSslCertificateResource,
				Check: resource.ComposeTestCheckFunc(
					testChecksslcertificateExists(TestSslcertificateName, true),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.test-cert", "name", SslcertificateName),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.test-cert", "partition", TestPartition),
				),
			},
		},
	})
}

func TestAccBigipSslCertificateTCs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslcertificateDestroyed,
		Steps: []resource.TestStep{
			{
				Config: loadFixtureString("../examples/bigip_ssl_cert_keys.tf"),
				Check: resource.ComposeTestCheckFunc(
					testChecksslcertificateExists("ssl-test-certificate-tc1", true),
					testChecksslcertificateExists("ssl-test-certificate-tc2", true),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc1", "name", "ssl-test-certificate-tc1"),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc2", "name", "ssl-test-certificate-tc2"),
				),
			},
			{
				Config: loadFixtureString("../examples/bigip_ssl_cert_keys.tf"),
				Check: resource.ComposeTestCheckFunc(
					testChecksslcertificateExists("ssl-test-certificate-tc1", true),
					testChecksslcertificateExists("ssl-test-certificate-tc2", true),
					testChecksslcertificateExists("ssl-test-certificate-tc10", false),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc1", "name", "ssl-test-certificate-tc1"),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc2", "name", "ssl-test-certificate-tc2"),
				),
			},
			{
				Config: loadFixtureString("../examples/bigip_ssl_certificate.tf"),
				Check: resource.ComposeTestCheckFunc(
					testChecksslcertificateExists("ssl-test-certificate-tc1", true),
					testChecksslcertificateExists("ssl-test-certificate-tc2", true),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc1", "name", "ssl-test-certificate-tc1"),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc2", "name", "ssl-test-certificate-tc2"),
				),
			},
		},
	})
}

func TestAccBigipSslCertificateOCSP(t *testing.T) {
	// Generated once, up front: the CA half is installed as /Common/MyCA
	// by testAcctPreCheckOCSP, and the leaf half (cryptographically signed
	// by that same CA key) is embedded directly in testSslCertOCSPResource
	// below. See installMyCA's doc comment for why these must be
	// generated together as a matched pair rather than using a static
	// fixture file.
	caCertPEM, leafCertPEM, _ := generateCAAndLeafCert(t, "MyCA", "ocsp-leaf.f5.com")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheckOCSP(caCertPEM)(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslcertificateDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testSslCertOCSPResource(leafCertPEM),
				Check: resource.ComposeTestCheckFunc(
					testChecksslcertificateExists("test-certificate", true),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc1", "name", "test-certificate"),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc1", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc1", "monitoring_type", "ocsp"),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc1", "issuer_cert", "/Common/MyCA"),
					resource.TestCheckResourceAttr("bigip_ssl_certificate.ssl-test-certificate-tc1", "ocsp", "/Common/testocsp1"),
				),
			},
		},
	})
}

func testChecksslcertificateExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetCertificate(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf(" SSL Certificate %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("SSL Certificate %s still exists.", name)
		}
		return nil
	}
}

func testChecksslcertificateDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ssl_certificate" {
			continue
		}
		name := rs.Primary.ID
		var sslCertificatename = fmt.Sprintf("~%s~%s", TestPartition, name)
		certificate, err := client.GetCertificate(sslCertificatename)
		if err != nil {
			return err
		}
		if certificate != nil {
			return fmt.Errorf("SSL Certificate %s not destroyed.", sslCertificatename)
		}
	}
	return nil
}
