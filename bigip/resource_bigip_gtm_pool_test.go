package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TEST_POOL_NAME = "test_pool"
var TEST_POOL_TYPE = "a"

func TestAccBigipGtmPool_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmPoolConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmPoolExists(TEST_POOL_NAME, TEST_POOL_TYPE, true),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "name", TEST_POOL_NAME),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "type", TEST_POOL_TYPE),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "load_balancing_mode", "round-robin"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "enabled", "true"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "ttl", "30"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "max_answers_returned", "1"),
				),
			},
		},
	})
}

func TestAccBigipGtmPool_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmPoolConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmPoolExists(TEST_POOL_NAME, TEST_POOL_TYPE, true),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "load_balancing_mode", "round-robin"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "ttl", "30"),
				),
			},
			{
				Config: testAccBigipGtmPoolConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmPoolExists(TEST_POOL_NAME, TEST_POOL_TYPE, true),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "load_balancing_mode", "ratio"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "ttl", "60"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "max_answers_returned", "2"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool", "fallback_mode", "return-to-dns"),
				),
			},
		},
	})
}

func TestAccBigipGtmPool_withMembers(t *testing.T) {
	poolName := "test_pool_with_members"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmPoolConfigWithMembers(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmPoolExists(poolName, TEST_POOL_TYPE, true),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-members", "name", poolName),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-members", "load_balancing_mode", "ratio"),
				),
			},
		},
	})
}

func TestAccBigipGtmPool_withQoS(t *testing.T) {
	poolName := "test_pool_qos"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmPoolConfigWithQoS(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmPoolExists(poolName, TEST_POOL_TYPE, true),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-qos", "qos_hit_ratio", "10"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-qos", "qos_hops", "5"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-qos", "qos_kilobytes_second", "5"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-qos", "qos_lcs", "50"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-qos", "qos_packet_rate", "5"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-qos", "qos_rtt", "100"),
				),
			},
		},
	})
}

func TestAccBigipGtmPool_withLimits(t *testing.T) {
	poolName := "test_pool_limits"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmPoolConfigWithLimits(poolName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmPoolExists(poolName, TEST_POOL_TYPE, true),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-limits", "limit_max_connections", "1000"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-limits", "limit_max_connections_status", "enabled"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-limits", "limit_max_bps", "10000000"),
					resource.TestCheckResourceAttr("bigip_gtm_pool.test-pool-limits", "limit_max_bps_status", "enabled"),
				),
			},
		},
	})
}

func TestAccBigipGtmPool_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmPoolConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmPoolExists(TEST_POOL_NAME, TEST_POOL_TYPE, true),
				),
			},
			{
				ResourceName:      "bigip_gtm_pool.test-pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckGtmPoolExists(name, poolType string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		fullPath := fmt.Sprintf("/Common/%s", name)

		pool, err := client.GetGTMPool(fullPath, poolType)
		if err != nil {
			return err
		}
		if exists && pool == nil {
			return fmt.Errorf("GTM Pool %s does not exist", fullPath)
		}
		if !exists && pool != nil {
			return fmt.Errorf("GTM Pool %s still exists", fullPath)
		}
		return nil
	}
}

func testCheckGtmPoolDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_pool" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		poolType := rs.Primary.Attributes["type"]
		partition := rs.Primary.Attributes["partition"]
		fullPath := fmt.Sprintf("/%s/%s", partition, name)

		pool, err := client.GetGTMPool(fullPath, poolType)
		if err != nil {
			return err
		}
		if pool != nil {
			return fmt.Errorf("GTM Pool %s still exists", fullPath)
		}
	}
	return nil
}

func testAccBigipGtmPoolConfig() string {
	return fmt.Sprintf(`
resource "bigip_gtm_pool" "test-pool" {
  name      = "%s"
  type      = "%s"
  partition = "Common"
  
  load_balancing_mode   = "round-robin"
  alternate_mode        = "round-robin"
  fallback_mode         = "return-to-dns"
  fallback_ip           = "any"
  
  enabled               = true
  dynamic_ratio         = "disabled"
  manual_resume         = "disabled"
  
  max_answers_returned  = 1
  ttl                   = 30
  
  verify_member_availability = "enabled"
}
`, TEST_POOL_NAME, TEST_POOL_TYPE)
}

func testAccBigipGtmPoolConfigUpdated() string {
	return fmt.Sprintf(`
resource "bigip_gtm_pool" "test-pool" {
  name      = "%s"
  type      = "%s"
  partition = "Common"
  
  load_balancing_mode   = "ratio"
  alternate_mode        = "topology"
  fallback_mode         = "return-to-dns"
  fallback_ip           = "any"
  
  enabled               = true
  dynamic_ratio         = "disabled"
  manual_resume         = "disabled"
  
  max_answers_returned  = 2
  ttl                   = 60
  
  verify_member_availability = "enabled"
}
`, TEST_POOL_NAME, TEST_POOL_TYPE)
}

func testAccBigipGtmPoolConfigWithMembers(poolName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_pool" "test-pool-members" {
  name      = "%s"
  type      = "%s"
  partition = "Common"
  
  load_balancing_mode = "ratio"
  
  members {
    name         = "test_server1:/Common/vs1"
    enabled      = true
    ratio        = 2
    member_order = 0
    monitor      = "default"
  }
  
  members {
    name         = "test_server2:/Common/vs2"
    enabled      = true
    ratio        = 1
    member_order = 1
    monitor      = "default"
  }
}
`, poolName, TEST_POOL_TYPE)
}

func testAccBigipGtmPoolConfigWithQoS(poolName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_pool" "test-pool-qos" {
  name      = "%s"
  type      = "%s"
  partition = "Common"
  
  load_balancing_mode = "round-robin"
  
  qos_hit_ratio        = 10
  qos_hops             = 5
  qos_kilobytes_second = 5
  qos_lcs              = 50
  qos_packet_rate      = 5
  qos_rtt              = 100
  qos_topology         = 0
  qos_vs_capacity      = 0
  qos_vs_score         = 0
}
`, poolName, TEST_POOL_TYPE)
}

func testAccBigipGtmPoolConfigWithLimits(poolName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_pool" "test-pool-limits" {
  name      = "%s"
  type      = "%s"
  partition = "Common"
  
  load_balancing_mode          = "round-robin"
  
  limit_max_connections        = 1000
  limit_max_connections_status = "enabled"
  limit_max_bps                = 10000000
  limit_max_bps_status         = "enabled"
  limit_max_pps                = 1000
  limit_max_pps_status         = "enabled"
}
`, poolName, TEST_POOL_TYPE)
}
