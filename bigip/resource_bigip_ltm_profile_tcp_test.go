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
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_TCP_NAME = fmt.Sprintf("/%s/test-tcp", TEST_PARTITION)

var TEST_TCP_RESOURCE = `
resource "bigip_ltm_profile_tcp" "test-tcp" {
            name = "/Common/sanjose-tcp-wan-profile"
            defaults_from = "/Common/tcp-wan-optimized"
						partition = "Common"
            idle_timeout = 300
            close_wait_timeout = 5
            finwait_2timeout = 5
            finwait_timeout = 300
            keepalive_interval = 1700
            deferred_accept = "enabled"
            fast_open = "enabled"
        }
`

func TestAccBigipLtmProfileTcp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_TCP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(TEST_TCP_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "name", "/Common/sanjose-tcp-wan-profile"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "defaults_from", "/Common/tcp-wan-optimized"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "idle_timeout", "300"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "close_wait_timeout", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "finwait_2timeout", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "finwait_timeout", "300"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "keepalive_interval", "1700"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "deferred_accept", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "fast_open", "enabled"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileTcp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_TCP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(TEST_TCP_NAME, true),
				),
				ResourceName:      TEST_TCP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckTcpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetTcp(name)
		if err != nil {
			return err
		}
		if exists && p != nil {
			return fmt.Errorf("tcp %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("tcp %s still exists.", name)
		}
		return nil
	}
}

func testCheckTcpsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_tcp" {
			continue
		}

		name := rs.Primary.ID
		tcp, err := client.GetTcp(name)
		if err != nil {
			return err
		}
		if tcp != nil {
			return fmt.Errorf("tcp %s not destroyed.", name)
		}
	}
	return nil
}
