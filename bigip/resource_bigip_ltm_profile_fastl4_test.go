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

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_FASTL4_NAME = fmt.Sprintf("/%s/test-fastl4", TEST_PARTITION)

var TEST_FASTL4_RESOURCE = `
resource "bigip_ltm_profile_fastl4" "test-fastl4" {
            name = "` + TEST_FASTL4_NAME + `"
            partition = "Common"
            defaults_from = "/Common/fastL4"
			client_timeout = 40
			idle_timeout = "200"
            explicitflow_migration = "enabled"
            hardware_syncookie = "enabled"
            iptos_toclient = "pass-through"
            iptos_toserver = "pass-through"
            keepalive_interval = "disabled"
 }
`

func TestAccBigipLtmProfileFastl4_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfastl4sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FASTL4_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(TEST_FASTL4_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "name", TEST_FASTL4_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "defaults_from", "/Common/fastL4"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "client_timeout", "40"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "explicitflow_migration", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "hardware_syncookie", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "idle_timeout", "200"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "hardware_syncookie", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "iptos_toclient", "pass-through"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "iptos_toserver", "pass-through"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "keepalive_interval", "disabled"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileFastl4_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfastl4sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FASTL4_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(TEST_FASTL4_NAME, true),
				),
				ResourceName:      TEST_FASTL4_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckfastl4Exists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetFastl4(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("fastl4 %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("fastl4 %s still exists.", name)
		}
		return nil
	}
}

func testCheckfastl4sDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_fastl4" {
			continue
		}

		name := rs.Primary.ID
		fastl4, err := client.GetFastl4(name)
		if err != nil {
			return err
		}
		if fastl4 != nil {
			return fmt.Errorf("fastl4 %s not destroyed.", name)
		}
	}
	return nil
}
