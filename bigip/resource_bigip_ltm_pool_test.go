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
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_POOL_NAME = fmt.Sprintf("/%s/test-pool", TEST_PARTITION)
var TEST_POOLNODE_NAME = fmt.Sprintf("/%s/test-node", TEST_PARTITION)
var TEST_POOLNODE_NAMEPORT = fmt.Sprintf("%s:443", TEST_POOLNODE_NAME)

var TEST_POOL_RESOURCE = `
resource "bigip_ltm_node" "test-node" {
	name = "` + TEST_NODE_NAME + `"
	address = "10.10.10.10"
	connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "default"
	rate_limit = "disabled"
}

resource "bigip_ltm_pool" "test-pool" {
	name = "` + TEST_POOL_NAME + `"
	monitors = ["/Common/http"]
	allow_nat = "yes"
	allow_snat = "yes"
	description = "Test-Pool-Sample"
	load_balancing_mode = "round-robin"
	slow_ramp_time = "5"
	service_down_action = "reset"
	reselect_tries = "2"
}

resource "bigip_ltm_pool_attachment" "test-pool_test-node" {
	pool = "` + TEST_POOL_NAME + `"
	node = "` + TEST_POOLNODE_NAMEPORT + `"
	depends_on = ["bigip_ltm_node.test-node", "bigip_ltm_pool.test-pool"]
}
`

func TestAccBigipLtmPool_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_POOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "name", TEST_POOL_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "allow_nat", "yes"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "allow_snat", "yes"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "description", "Test-Pool-Sample"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "load_balancing_mode", "round-robin"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "slow_ramp_time", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "service_down_action", "reset"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "reselect_tries", "2"),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TEST_POOL_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", TEST_POOLNODE_NAMEPORT),
				),
			},
		},
	})
}

func TestAccBigipLtmPool_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_POOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
				),
				ResourceName:      TEST_POOL_RESOURCE,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

//TODO: test adding/removing nodes

func testCheckPoolExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.GetPool(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("Pool %s does not exist.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("Pool %s exists.", name)
		}
		return nil
	}
}

func testCheckPoolsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_pool" {
			continue
		}

		name := rs.Primary.ID
		pool, err := client.GetPool(name)
		if err != nil {
			return err
		}
		if pool != nil {
			return fmt.Errorf("Pool %s not destroyed.", name)
		}
	}
	return nil
}
