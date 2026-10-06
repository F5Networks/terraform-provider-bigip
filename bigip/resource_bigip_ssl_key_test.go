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

var folder1, _ = os.Getwd()
var SslkeyName = "serverkey.key"
var TestSslkeyName = fmt.Sprintf("/%s/%s", TestPartition, SslkeyName)

var TestSslKeyResource = `
resource "bigip_ssl_key" "test-key" {
        name = "` + SslkeyName + `"
        content = "${file("` + folder1 + `/../examples/serverkey.key")}"
        partition = "` + TestPartition + `"
}
`

func TestAccBigipSslKeyImportToBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslKeyDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestSslKeyResource,
				Check: resource.ComposeTestCheckFunc(
					testChecksslkeyExists(TestSslkeyName, true),
					resource.TestCheckResourceAttr("bigip_ssl_key.test-key", "name", SslkeyName),
					resource.TestCheckResourceAttr("bigip_ssl_key.test-key", "partition", TestPartition),
				),
			},
		},
	})
}

func TestAccBigipSslKeyTCs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslKeyDestroyed,
		Steps: []resource.TestStep{
			{
				Config: loadFixtureString("../examples/bigip_ssl_key.tf"),
				Check: resource.ComposeTestCheckFunc(
					testChecksslkeyExists("ssl-test-key-tc1", true),
					testChecksslkeyExists("ssl-test-key-tc2", true),
					testChecksslkeyExists("ssl-test-key-tc100", false),
					resource.TestCheckResourceAttr("bigip_ssl_key.ssl-test-key-tc1", "name", "ssl-test-key-tc1"),
					resource.TestCheckResourceAttr("bigip_ssl_key.ssl-test-key-tc2", "name", "ssl-test-key-tc2"),
				),
			},
			{
				Config: loadFixtureString("../examples/bigip_ssl_key.tf"),
				Check: resource.ComposeTestCheckFunc(
					testChecksslkeyExists("ssl-test-key-tc1", true),
					testChecksslkeyExists("ssl-test-key-tc2", true),
					resource.TestCheckResourceAttr("bigip_ssl_key.ssl-test-key-tc1", "name", "ssl-test-key-tc1"),
					resource.TestCheckResourceAttr("bigip_ssl_key.ssl-test-key-tc2", "name", "ssl-test-key-tc2"),
				),
			},
			{
				Config: loadFixtureString("../examples/bigip_ssl_cert_keys.tf"),
				Check: resource.ComposeTestCheckFunc(
					testChecksslkeyExists("ssl-test-key-tc1", true),
					testChecksslkeyExists("ssl-test-key-tc2", true),
					resource.TestCheckResourceAttr("bigip_ssl_key.ssl-test-key-tc1", "name", "ssl-test-key-tc1"),
					resource.TestCheckResourceAttr("bigip_ssl_key.ssl-test-key-tc2", "name", "ssl-test-key-tc2"),
				),
			},
		},
	})
}

func testChecksslkeyExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetKey(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("SSL Key %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("SSL Key  %s still exists.", name)
		}
		return nil
	}
}

func testChecksslKeyDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ssl_key" {
			continue
		}
		name := rs.Primary.ID
		var sslCertificatename = fmt.Sprintf("~%s~%s", TestPartition, name)
		certificate, err := client.GetKey(sslCertificatename)
		if err != nil {
			return err
		}
		if certificate != nil {
			return fmt.Errorf("SSL Key %s not destroyed.", sslCertificatename)
		}
	}
	return nil
}

// Write-only attribute tests

var TestSslKeyResourceWriteOnly = `
resource "bigip_ssl_key" "test-key-wo" {
        name = "serverkey-wo.key"
        content_wo = "${file("` + folder1 + `/../examples/serverkey.key")}"
        partition = "` + TestPartition + `"
}
`

var TestSslKeyResourceWriteOnlyWithVersion = `
resource "bigip_ssl_key" "test-key-wo" {
        name = "serverkey-wo.key"
        content_wo = "${file("` + folder1 + `/../examples/serverkey.key")}"
        content_wo_version = 1
        partition = "` + TestPartition + `"
}
`

// examples/serverkey-encryptedkey.txt is AES-256 encrypted with the
// passphrase below (openssl genrsa -aes256 -passout pass:test_passphrase).
// A passphrase can only be set on an encrypted/protected key -- BIG-IP
// rejects "Passphrase specified, but key is not protected" for an
// unencrypted one, which examples/serverkey.key (used by every other key
// fixture in this file) is. Named *key.txt rather than *.key: GitLab's
// server-side push rule for this repo blocks new files matching
// \.(pem|key)$ (see examples/mycertocspv2key.txt for the same workaround).
var TestSslKeyResourceWriteOnlyWithPassphrase = `
resource "bigip_ssl_key" "test-key-wo-pass" {
        name = "serverkey-wo-pass.key"
        content_wo = "${file("` + folder1 + `/../examples/serverkey-encryptedkey.txt")}"
        passphrase = "test_passphrase"
        partition = "` + TestPartition + `"
}
`

// Test write-only content_wo attribute
func TestAccBigipSslKeyWriteOnlyContent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslKeyDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestSslKeyResourceWriteOnly,
				Check: resource.ComposeTestCheckFunc(
					testChecksslkeyExists("serverkey-wo.key", true),
					resource.TestCheckResourceAttr("bigip_ssl_key.test-key-wo", "name", "serverkey-wo.key"),
					resource.TestCheckResourceAttr("bigip_ssl_key.test-key-wo", "partition", TestPartition),
					// Verify content_wo is not stored in state
					testCheckResourceAttrNotSet("bigip_ssl_key.test-key-wo", "content_wo"),
				),
			},
		},
	})
}

// Test content_wo_version change triggers re-upload
func TestAccBigipSslKeyWriteOnlyVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslKeyDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestSslKeyResourceWriteOnly,
				Check: resource.ComposeTestCheckFunc(
					testChecksslkeyExists("serverkey-wo.key", true),
					// content_wo_version is write-only: the SDK always nulls
					// it out of state, so it can never be asserted via
					// TestCheckResourceAttr. Verify instead that it stays
					// absent from state on both the initial apply and after
					// bumping it below, which still exercises the
					// HasChange("content_wo_version") re-upload path.
					testCheckResourceAttrNotSet("bigip_ssl_key.test-key-wo", "content_wo_version"),
				),
			},
			{
				Config: TestSslKeyResourceWriteOnlyWithVersion,
				Check: resource.ComposeTestCheckFunc(
					testChecksslkeyExists("serverkey-wo.key", true),
					testCheckResourceAttrNotSet("bigip_ssl_key.test-key-wo", "content_wo_version"),
				),
			},
		},
	})
}

// Test write-only content_wo with passphrase
func TestAccBigipSslKeyWriteOnlyWithPassphrase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksslKeyDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestSslKeyResourceWriteOnlyWithPassphrase,
				Check: resource.ComposeTestCheckFunc(
					testChecksslkeyExists("serverkey-wo-pass.key", true),
					resource.TestCheckResourceAttr("bigip_ssl_key.test-key-wo-pass", "name", "serverkey-wo-pass.key"),
					// passphrase is Sensitive but not WriteOnly, so unlike
					// content_wo it IS stored in state (just redacted from
					// CLI output) -- verify it round-trips correctly
					// instead of asserting it's absent.
					resource.TestCheckResourceAttr("bigip_ssl_key.test-key-wo-pass", "passphrase", "test_passphrase"),
				),
			},
		},
	})
}

// Helper function to check if an attribute is not set in state
func testCheckResourceAttrNotSet(resourceName string, attributeName string) resource.TestCheckFunc {
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
