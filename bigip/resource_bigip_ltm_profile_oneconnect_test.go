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

var TEST_ONECONNECT_NAME = fmt.Sprintf("/%s/test-oneconnect", TestPartition)

var TEST_ONECONNECT_RESOURCE = `
resource "bigip_ltm_profile_oneconnect" "test-oneconnect" {
            name = "/Common/test-oneconnect"
            partition = "Common"
            defaults_from = "/Common/oneconnect"
            idle_timeout_override = "disabled"
            max_age = 3600
            max_reuse = 1000
            max_size = 1000
            share_pools = "disabled"
            source_mask = "255.255.255.255"
        }
`

func TestAccBigipLtmProfileoneconnect_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckoneconnectsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ONECONNECT_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckoneconnectExists(TEST_ONECONNECT_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_oneconnect.test-oneconnect", "name", "/Common/test-oneconnect"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_oneconnect.test-oneconnect", "defaults_from", "/Common/oneconnect"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_oneconnect.test-oneconnect", "idle_timeout_override", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_oneconnect.test-oneconnect", "max_age", "3600"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_oneconnect.test-oneconnect", "max_reuse", "1000"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_oneconnect.test-oneconnect", "max_size", "1000"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_oneconnect.test-oneconnect", "share_pools", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_oneconnect.test-oneconnect", "source_mask", "255.255.255.255"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileoneconnect_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckoneconnectsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ONECONNECT_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckoneconnectExists(TEST_ONECONNECT_NAME, true),
				),
				ResourceName:      TEST_ONECONNECT_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckoneconnectExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetOneconnect(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("oneconnect %s was not created.", name)
		}
		if !exists && p != nil {

			return fmt.Errorf("oneconnect %s still exists.", name)
		}
		return nil
	}
}

func testCheckoneconnectsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_oneconnect" {
			continue
		}

		name := rs.Primary.ID
		oneconnect, err := client.GetOneconnect(name)
		if err != nil {
			return err
		}
		if oneconnect != nil {
			return fmt.Errorf("oneconnect %s not destroyed.", name)
		}
	}
	return nil
}
