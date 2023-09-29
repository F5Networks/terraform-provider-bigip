/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testCipherGroupConfigTC1 = `
resource "bigip_ltm_cipher_group" "test-cipher-group" {
  name     = "/Common/test-cipher-group-01"
  allow    = ["/Common/f5-aes"]
  require  = ["/Common/f5-quic"]
  ordering = "speed"
}
`

func TestAccBigipLtmCipherGroupCreateTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckCipherGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testCipherGroupConfigTC1,
				Check: resource.ComposeTestCheckFunc(
					testCheckCipherGroupExists("/Common/test-cipher-group-01"),
					resource.TestCheckResourceAttr("bigip_ltm_cipher_group.test-cipher-group", "name", "/Common/test-cipher-group-01"),
				),
			},
		},
	})
}

func testCheckCipherGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.GetLtmCipherGroup(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("Pool %s does not exist ", name)
		}

		return nil
	}
}

func testCheckCipherGroupDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_cipher_group" {
			continue
		}
		name := rs.Primary.ID
		pool, err := client.GetLtmCipherGroup(name)
		if err != nil {
			return err
		}
		if pool != nil {
			return fmt.Errorf("Cipher rule %s not destroyed ", name)
		}
	}
	return nil
}
