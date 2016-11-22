package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
	"testing"
)

var TEST_POOL_NAME = fmt.Sprintf("/%s/test-pool", TEST_PARTITION)

var TEST_POOL_RESOURCE = `
resource "bigip_ltm_pool" "test-pool" {
	name = "` + TEST_POOL_NAME + `"
	monitors = ["/Common/http"]
	allow_nat = true
	allow_snat = true
	load_balancing_mode = "round-robin"
}
`

func TestBigipLtmPool_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_POOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "name", TEST_POOL_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "allow_nat", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "allow_snat", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_pool.test-pool", "load_balancing_mode", "round-robin"),
				),
			},
		},
	})
}

func TestBigipLtmPool_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_POOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_POOL_NAME, true),
				),
				ResourceName:      TEST_POOL_RESOURCE,
				ImportState:       true,
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
			return fmt.Errorf("Pool ", name, " does not exist.")
		}
		if !exists && p != nil {
			return fmt.Errorf("Pool ", name, " exists.")
		}
		return nil
	}
}

func testCheckPoolMember(pool_name, member_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		members, err := client.PoolMembers(pool_name)
		if err != nil {
			return err
		}

		for _, member := range members {
			if member == member_name {
				return nil
			}
		}

		return fmt.Errorf("Member %s not found in %s", member_name, pool_name)
	}
}

func testCheckEmptyPool(pool_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		members, err := client.PoolMembers(pool_name)
		if err != nil {
			return err
		}
		if len(members) != 0 {
			return fmt.Errorf("Pool %s not empty (%d members))", pool_name, len(members))
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
			return fmt.Errorf("Pool ", name, " not destroyed.")
		}
	}
	return nil
}
