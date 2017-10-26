package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_SNATPOOL_NAME = fmt.Sprintf("/%s/test-snatpool", TEST_PARTITION)

var TEST_SNATPOOL_RESOURCE = `
resource "bigip_snatpool" "test-snatpool" {
  name = "/Common/snatpool_sanjose"
  members = ["191.1.1.1","194.2.2.2"]
}

`

func TestBigipLtmsnatpool_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksnatpoolsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_SNATPOOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnatpoolExists(TEST_SNATPOOL_NAME, true),
					resource.TestCheckResourceAttr("bigip_snatpool.test-snatpool", "name", "/Common/snatpool_sanjose"),
					resource.TestCheckResourceAttr("bigip_snatpool.test-snatpool",
						fmt.Sprintf("members.%d", schema.HashString("191.1.1.1")),
						"191.1.1.1"),
					resource.TestCheckResourceAttr("bigip_snatpool.test-snatpool",
						fmt.Sprintf("members.%d", schema.HashString("194.2.2.2")),
						"194.2.2.2"),
				),
			},
		},
	})
}

func TestBigipLtmsnatpool_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksnatpoolsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_SNATPOOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnatpoolExists(TEST_SNATPOOL_NAME, true),
				),
				ResourceName:      TEST_SNATPOOL_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testChecksnatpoolExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.Snatpools(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("snatpool ", name, " was not created.")
		}
		if !exists && p != nil {
			return fmt.Errorf("snatpool ", name, " still exists.")
		}
		return nil
	}
}

func testChecksnatpoolsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_snatpool" {
			continue
		}

		name := rs.Primary.ID
		snatpool, err := client.Snatpools(name)
		if err != nil {
			return err
		}
		if snatpool == nil {
			return fmt.Errorf("snatpool ", name, " not destroyed.")
		}
	}
	return nil
}
