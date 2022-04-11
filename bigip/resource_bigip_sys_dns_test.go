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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_DNS_NAME = fmt.Sprintf("/%s/test-dns", TEST_PARTITION)

var TEST_DNS_RESOURCE = `
resource "bigip_sys_dns" "test-dns" {
   description = "` + TEST_DNS_NAME + `"
   name_servers = ["1.1.1.1"]
   number_of_dots = 2
   search = ["f5.com"]
}

`

func TestAccBigipSysdns_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TEST_DNS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TEST_DNS_NAME, true),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns", "description", TEST_DNS_NAME),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns", "number_of_dots", "2"),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns",
						fmt.Sprintf("name_servers.%d", schema.HashString("1.1.1.1")),
						"1.1.1.1"),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns",
						fmt.Sprintf("search.%d", schema.HashString("f5.com")),
						"f5.com"),
				),
			},
		},
	})
}

func TestAccBigipSysdns_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TEST_DNS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TEST_DNS_NAME, true),
				),
				ResourceName:      TEST_DNS_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckdnsExists(description string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		dns, err := client.DNSs()
		if err != nil {
			return err
		}
		if exists && dns == nil {
			return fmt.Errorf("dns %s was not created.", description)

		}
		if !exists && dns != nil {
			return fmt.Errorf("dns %s still exists.", description)

		}
		return nil
	}
}
