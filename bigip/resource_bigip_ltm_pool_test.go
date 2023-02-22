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

var poolMember1 = fmt.Sprintf("%s:443", "10.10.10.10")
var TestPoolName = fmt.Sprintf("/%s/test-pool", TestPartition)

var TestPoolResource = `
/*resource "bigip_ltm_node" "test-node" {
	name = "` + TestNodeName + `"
	address = "10.10.10.10"
	connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "default"
	rate_limit = "disabled"
}*/

resource "bigip_ltm_pool" "test-pool" {
	name = "` + TestPoolName + `"
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
         pool = bigip_ltm_pool.test-pool.name 
         node = "` + poolMember1 + `"
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
				Config: TestPoolResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "name", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "allow_nat", "yes"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "allow_snat", "yes"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "description", "Test-Pool-Sample"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "load_balancing_mode", "round-robin"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "slow_ramp_time", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "service_down_action", "reset"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "reselect_tries", "2"),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", poolMember1),
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
				Config: TestPoolResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
				),
				ResourceName:      TestPoolResource,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

//TODO: test adding/removing nodes

func testCheckPoolExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.GetPool(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("Pool %s does not exist ", name)
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
			return fmt.Errorf("Pool %s not destroyed ", name)
		}
	}
	return nil
}
