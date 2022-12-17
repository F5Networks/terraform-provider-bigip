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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var folder, _ = os.Getwd()
var SslcertificateName = "servercert.crt"
var TestSslcertificateName = fmt.Sprintf("/%s/%s", TEST_PARTITION, SslcertificateName)

var TestSslCertificateResource = `
resource "bigip_ssl_certificate" "test-cert" {
        name = "` + SslcertificateName + `"
        content = "${file("` + folder + `/../examples/servercert.crt")}"
        partition = "Common"
}
`

func TestAccSslCertificateImportToBigip(t *testing.T) {
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
					resource.TestCheckResourceAttr("bigip_ssl_certificate.test-cert", "partition", TEST_PARTITION),
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
			return fmt.Errorf("SSL Certificate %s was not created.", name)
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
		var sslCertificatename = fmt.Sprintf("~%s~%s", TEST_PARTITION, name)
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
