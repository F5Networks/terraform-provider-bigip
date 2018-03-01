package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
	load_balancing_mode = "round-robin"
	slow_ramp_time = "5"
	service_down_action = "reset"
	nodes = ["` + TEST_POOLNODE_NAMEPORT + `"]
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
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "load_balancing_mode", "round-robin"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "slow_ramp_time", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "service_down_action", "reset"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "nodes.#", "1"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool",
						fmt.Sprintf("nodes.%d", schema.HashString(TEST_POOLNODE_NAMEPORT)),
						TEST_POOLNODE_NAMEPORT),
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

/* func testCheckPoolMember(poolName, memberName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		members, err := client.PoolMembers(poolName)
		if err != nil {
			return err
		}

		for _, member := range members {
			if member.Name == memberName {
				return nil
			}
		}

		return fmt.Errorf("Member %s not found in %s", memberName, poolName)
	}
}

func testCheckEmptyPool(poolName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		members, err := client.PoolMembers(poolName)
		if err != nil {
			return err
		}
		if len(members) != 0 {
			return fmt.Errorf("Pool %s not empty (%d members))", poolName, len(members))
		}
		return nil
	}
}
*/

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
