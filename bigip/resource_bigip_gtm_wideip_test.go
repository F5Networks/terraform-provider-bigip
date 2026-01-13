package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TEST_WIDEIP_NAME = "testwideip.local"
var TEST_WIDEIP_TYPE = "a"

func TestAccBigipGtmWideip_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmWideipDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmWideipConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmWideipExists(TEST_WIDEIP_NAME, TEST_WIDEIP_TYPE, true),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "name", TEST_WIDEIP_NAME),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "type", TEST_WIDEIP_TYPE),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "description", "test_wideip_a"),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "enabled", "true"),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "minimal_response", "enabled"),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "pool_lb_mode", "round-robin"),
				),
			},
		},
	})
}

func TestAccBigipGtmWideip_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmWideipDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmWideipConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmWideipExists(TEST_WIDEIP_NAME, TEST_WIDEIP_TYPE, true),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "description", "test_wideip_a"),
				),
			},
			{
				Config: testAccBigipGtmWideipConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmWideipExists(TEST_WIDEIP_NAME, TEST_WIDEIP_TYPE, true),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "description", "updated_description"),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "pool_lb_mode", "topology"),
					resource.TestCheckResourceAttr("bigip_gtm_wideip.test-wideip", "ttl_persistence", "7200"),
				),
			},
		},
	})
}

func TestAccBigipGtmWideip_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmWideipDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmWideipConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmWideipExists(TEST_WIDEIP_NAME, TEST_WIDEIP_TYPE, true),
				),
			},
			{
				ResourceName:      "bigip_gtm_wideip.test-wideip",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s:/Common/%s", TEST_WIDEIP_TYPE, TEST_WIDEIP_NAME),
			},
		},
	})
}

func testCheckGtmWideipExists(name, wideipType string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		fullPath := fmt.Sprintf("/Common/%s", name)

		wideip, err := client.GetGTMWideIP(fullPath, wideipType)
		if err != nil {
			return err
		}
		if exists && wideip == nil {
			return fmt.Errorf("WideIP %s does not exist", fullPath)
		}
		if !exists && wideip != nil {
			return fmt.Errorf("WideIP %s still exists", fullPath)
		}
		return nil
	}
}

func testCheckGtmWideipDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_wideip" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		wideipType := rs.Primary.Attributes["type"]
		partition := rs.Primary.Attributes["partition"]
		fullPath := fmt.Sprintf("/%s/%s", partition, name)

		wideip, err := client.GetGTMWideIP(fullPath, wideipType)
		if err != nil {
			return err
		}
		if wideip != nil {
			return fmt.Errorf("WideIP %s still exists", fullPath)
		}
	}
	return nil
}

func testAccBigipGtmWideipConfig() string {
	return fmt.Sprintf(`
resource "bigip_gtm_wideip" "test-wideip" {
  name      = "%s"
  type      = "%s"
  partition = "Common"
  
  description              = "test_wideip_a"
  enabled                  = true
  failure_rcode            = "noerror"
  failure_rcode_response   = "disabled"
  failure_rcode_ttl        = 0
  minimal_response         = "enabled"
  persist_cidr_ipv4        = 32
  persist_cidr_ipv6        = 128
  persistence              = "disabled"
  pool_lb_mode             = "round-robin"
  ttl_persistence          = 3600
}
`, TEST_WIDEIP_NAME, TEST_WIDEIP_TYPE)
}

func testAccBigipGtmWideipConfigUpdated() string {
	return fmt.Sprintf(`
resource "bigip_gtm_wideip" "test-wideip" {
  name      = "%s"
  type      = "%s"
  partition = "Common"
  
  description              = "updated_description"
  enabled                  = true
  failure_rcode            = "servfail"
  failure_rcode_response   = "enabled"
  failure_rcode_ttl        = 300
  minimal_response         = "disabled"
  persist_cidr_ipv4        = 24
  persist_cidr_ipv6        = 64
  persistence              = "enabled"
  pool_lb_mode             = "topology"
  ttl_persistence          = 7200
}
`, TEST_WIDEIP_NAME, TEST_WIDEIP_TYPE)
}
