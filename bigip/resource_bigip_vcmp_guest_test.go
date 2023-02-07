/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2022 F5 Networks Inc.
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

var TEST_VCMP_GUEST_NAME = "tf-guest"

var TEST_VCMP_GUEST_RESOURCE = `
resource "bigip_vcmp_guest" "test-guest" {
  name = "` + TEST_VCMP_GUEST_NAME + `"
  initial_image = "12.1.2.iso"
  mgmt_network = "bridged"
  mgmt_address = "10.1.1.1/24"
  mgmt_route = "none"
  state = "provisioned"
  cores_per_slot = 2
  number_of_slots = 1
  min_number_of_slots = 1
  vlans = ["/Common/testvlan"]
}

`

func TestAccBigipVcmpguest_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckvcmpguestDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_VCMP_GUEST_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testBigipVcmpguestExists(TEST_VCMP_GUEST_NAME, true),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "name", TEST_VCMP_GUEST_NAME),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "initial_image", "12.1.2.iso"),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "mgmt_network", "bridged"),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "mgmt_address", "10.1.1.1/24"),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "mgmt_route", "none"),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "state", "provisioned"),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "cores_per_slot", "2"),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "number_of_slots", "1"),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "min_number_of_slots", "1"),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "min_number_of_slots", "1"),
					resource.TestCheckResourceAttr("bigip_vcmp_guest.test-guest", "vlans.0", "/Common/testvlan"),
				),
			},
		},
	})
}

func TestAccBigipVcmpguest_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckvcmpguestDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_VCMP_GUEST_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testBigipVcmpguestExists(TEST_VCMP_GUEST_NAME, true),
				),
				ResourceName:      TEST_VCMP_GUEST_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testBigipVcmpguestExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetVcmpGuest(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("vcmp guest %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("vcmp guest %s still exists.", name)
		}
		return nil
	}
}

func testCheckvcmpguestDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_vcmp_guest" {
			continue
		}

		name := rs.Primary.ID
		snatpool, err := client.GetVcmpGuest(name)
		if err != nil {
			return err
		}
		if snatpool == nil {
			return fmt.Errorf("vcmp guest %s not destroyed.", name)
		}
	}
	return nil
}
