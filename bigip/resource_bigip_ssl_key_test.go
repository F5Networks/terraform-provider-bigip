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

var folder1, _ = os.Getwd()
var SslkeyName = "serverkey.key"
var TestSslkeyName = fmt.Sprintf("/%s/%s", TEST_PARTITION, SslkeyName)

var TestSslKeyResource = `
resource "bigip_ssl_key" "test-key" {
        name = "` + SslkeyName + `"
        content = "${file("` + folder1 + `/../examples/serverkey.key")}"
        partition = "` + TEST_PARTITION + `"
}
`

func TestAccSslKeyImportToBigip(t *testing.T) {
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
					resource.TestCheckResourceAttr("bigip_ssl_key.test-key", "partition", TEST_PARTITION),
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
		var sslCertificatename = fmt.Sprintf("~%s~%s", TEST_PARTITION, name)
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
