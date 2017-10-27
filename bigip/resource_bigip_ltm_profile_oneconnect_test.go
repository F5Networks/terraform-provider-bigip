package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_ONECONNECT_NAME = fmt.Sprintf("/%s/test-oneconnect", TEST_PARTITION)

var TEST_ONECONNECT_RESOURCE = `
resource "bigip_ltm_oneconnect" "test-oneconnect"
        {
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

func TestBigipLtmoneconnect_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckoneconnectsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_ONECONNECT_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckoneconnectExists(TEST_ONECONNECT_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_oneconnect.test-oneconnect", "name", "/Common/test-oneconnect"),
					resource.TestCheckResourceAttr("bigip_ltm_oneconnect.test-oneconnect", "defaults_from", "/Common/oneconnect"),
					resource.TestCheckResourceAttr("bigip_ltm_oneconnect.test-oneconnect", "idle_timeout_override", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_oneconnect.test-oneconnect", "max_age", "3600"),
					resource.TestCheckResourceAttr("bigip_ltm_oneconnect.test-oneconnect", "max_reuse", "1000"),
					resource.TestCheckResourceAttr("bigip_ltm_oneconnect.test-oneconnect", "max_size", "1000"),
					resource.TestCheckResourceAttr("bigip_ltm_oneconnect.test-oneconnect", "share_pools", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_oneconnect.test-oneconnect", "source_mask", "255.255.255.255"),
				),
			},
		},
	})
}

func TestBigipLtmoneconnect_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckoneconnectsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
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
		p, err := client.Oneconnect(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("oneconnects", name, " was not created.")
		}
		if !exists && p != nil {
			return fmt.Errorf("oneconnects ", name, " still exists.")
		}
		return nil
	}
}

func testCheckoneconnectsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_oneconnect" {
			continue
		}

		name := rs.Primary.ID
		oneconnect, err := client.Oneconnect(name)
		if err != nil {
			return err
		}
		if oneconnect == nil {
			return fmt.Errorf("oneconnects ", name, " not destroyed.")
		}
	}
	return nil
}
