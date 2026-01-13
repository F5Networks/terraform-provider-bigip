package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TEST_DATACENTER_NAME = "test_datacenter"

func TestAccBigipGtmDatacenter_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmDatacenterDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmDatacenterConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmDatacenterExists(TEST_DATACENTER_NAME, true),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "name", TEST_DATACENTER_NAME),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "location", "Seattle"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "contact", "admin@example.com"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "enabled", "true"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "prober_preference", "inside-datacenter"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "prober_fallback", "any-available"),
				),
			},
		},
	})
}

func TestAccBigipGtmDatacenter_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmDatacenterDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmDatacenterConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmDatacenterExists(TEST_DATACENTER_NAME, true),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "location", "Seattle"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "contact", "admin@example.com"),
				),
			},
			{
				Config: testAccBigipGtmDatacenterConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmDatacenterExists(TEST_DATACENTER_NAME, true),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "location", "San Francisco"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "contact", "ops@example.com"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "description", "Updated datacenter"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "prober_preference", "outside-datacenter"),
				),
			},
		},
	})
}

func TestAccBigipGtmDatacenter_withProberSettings(t *testing.T) {
	dcName := "test_datacenter_prober"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmDatacenterDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmDatacenterConfigWithProber(dcName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmDatacenterExists(dcName, true),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter-prober", "prober_preference", "pool"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter-prober", "prober_fallback", "outside-datacenter"),
				),
			},
		},
	})
}

func TestAccBigipGtmDatacenter_disabled(t *testing.T) {
	dcName := "test_datacenter_disabled"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmDatacenterDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmDatacenterConfigDisabled(dcName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmDatacenterExists(dcName, true),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter-disabled", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccBigipGtmDatacenter_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmDatacenterDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmDatacenterConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmDatacenterExists(TEST_DATACENTER_NAME, true),
				),
			},
			{
				ResourceName:      "bigip_gtm_datacenter.test-datacenter",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckGtmDatacenterExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		fullPath := fmt.Sprintf("/Common/%s", name)

		datacenter, err := client.GetGTMDatacenter(fullPath)
		if err != nil {
			return err
		}
		if exists && datacenter == nil {
			return fmt.Errorf("GTM Datacenter %s does not exist", fullPath)
		}
		if !exists && datacenter != nil {
			return fmt.Errorf("GTM Datacenter %s still exists", fullPath)
		}
		return nil
	}
}

func testCheckGtmDatacenterDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_datacenter" {
			continue
		}

		fullPath := rs.Primary.ID

		datacenter, err := client.GetGTMDatacenter(fullPath)
		if err != nil {
			return err
		}
		if datacenter != nil {
			return fmt.Errorf("GTM Datacenter %s still exists", fullPath)
		}
	}
	return nil
}

func testAccBigipGtmDatacenterConfig() string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "%s"
  partition = "Common"
  
  location           = "Seattle"
  contact            = "admin@example.com"
  description        = "Test datacenter"
  
  enabled            = true
  prober_preference  = "inside-datacenter"
  prober_fallback    = "any-available"
}
`, TEST_DATACENTER_NAME)
}

func testAccBigipGtmDatacenterConfigUpdated() string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "%s"
  partition = "Common"
  
  location           = "San Francisco"
  contact            = "ops@example.com"
  description        = "Updated datacenter"
  
  enabled            = true
  prober_preference  = "outside-datacenter"
  prober_fallback    = "any-available"
}
`, TEST_DATACENTER_NAME)
}

func testAccBigipGtmDatacenterConfigWithProber(dcName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter-prober" {
  name      = "%s"
  partition = "Common"
  
  location           = "New York"
  contact            = "network@example.com"
  description        = "Datacenter with custom prober settings"
  
  enabled            = true
  prober_preference  = "pool"
  prober_fallback    = "outside-datacenter"
}
`, dcName)
}

func testAccBigipGtmDatacenterConfigDisabled(dcName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter-disabled" {
  name      = "%s"
  partition = "Common"
  
  location           = "London"
  contact            = "support@example.com"
  description        = "Disabled datacenter"
  
  enabled            = false
  prober_preference  = "inside-datacenter"
  prober_fallback    = "any-available"
}
`, dcName)
}
