/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

//TODO: delete not implemented in virtual address

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TEST_VA_NAME = fmt.Sprintf("/%s/test-va", TestPartition)
var TEST_VA_NAME_CHANGED = fmt.Sprintf("/%s/test-va-changed", TestPartition)
var TEST_VA_CONFIG = `
resource "bigip_ltm_virtual_address" "test-va" {
	name          = "%s"
	traffic_group = "/Common/none"
}
`
var TEST_VA_RESOURCE = fmt.Sprintf(TEST_VA_CONFIG, TEST_VA_NAME)
var TEST_VA_RESOURCE_NAME_CHANGED = fmt.Sprintf(TEST_VA_CONFIG, TEST_VA_NAME_CHANGED)

func TestAccBigipLtmVA_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVAsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_VA_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVAExists(TEST_VA_NAME, true),
				),
			},
			{
				Config:    TEST_VA_RESOURCE_NAME_CHANGED,
				PreConfig: func() { testCheckVAExists(TEST_VA_NAME, true) },
				Check: resource.ComposeTestCheckFunc(
					testCheckVAExists(TEST_VA_NAME, false),
					testCheckVAExists(TEST_VA_NAME_CHANGED, true),
				),
			},
		},
	})
}

func TestAccBigipLtmVA_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVAsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_VA_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVAExists(TEST_VA_NAME, true),
				),
				ResourceName:      TEST_VA_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckVAExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		vas, err := client.VirtualAddresses()
		if err != nil {
			return err
		}

		for _, va := range vas.VirtualAddresses {
			if va.FullPath == name {
				if !exists {
					return fmt.Errorf("Virtual address %s exists.", name)
				}
				return nil
			}
		}
		if exists {
			return fmt.Errorf("Virtual address %s does not exist.", name)
		}

		return nil
	}
}

func testCheckVAsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_virtual_address" {
			continue
		}

		name := rs.Primary.ID
		vas, err := client.VirtualAddresses()
		if err != nil {
			return err
		}
		for _, va := range vas.VirtualAddresses {
			if va.FullPath == name {
				return fmt.Errorf("Virtual address %s not destroyed.", name)
			}
		}
	}
	return nil
}
