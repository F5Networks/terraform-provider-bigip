package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_SNAT_NAME = fmt.Sprintf("/%s/test-snat", TEST_PARTITION)

var TEST_SNAT_RESOURCE = `
resource "bigip_snat" "test-snat" {
 name = "/Common/NewSnatList"
 translation = "136.1.1.1"
 origins = ["2.2.2.2", "3.3.3.3"]
}

`

func TestBigipLtmsnat_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksnatsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SNAT_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnatExists(TEST_SNAT_NAME, true),
					resource.TestCheckResourceAttr("bigip_snat.test-snat", "name", "/Common/NewSnatList"),
					resource.TestCheckResourceAttr("bigip_snat.test-snat", "translation", "136.1.1.1"),
					resource.TestCheckResourceAttr("bigip_snat.test-snat",
						fmt.Sprintf("origins.%d", schema.HashString("2.2.2.2")),
						"2.2.2.2"),
					resource.TestCheckResourceAttr("bigip_snat.test-snat",
						fmt.Sprintf("origins.%d", schema.HashString("3.3.3.3")),
						"3.3.3.3"),
				),
			},
		},
	})
}

func TestBigipLtmsnat_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksnatsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SNAT_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnatExists(TEST_SNAT_NAME, true),
				),
				ResourceName:      TEST_SNAT_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testChecksnatExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.Snats(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("Snat ", name, " was not created.")
		}
		if !exists && p != nil {
			return fmt.Errorf("Snat ", name, " still exists.")
		}
		return nil
	}
}

func testChecksnatsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_snat" {
			continue
		}

		name := rs.Primary.ID
		snat, err := client.Snats(name)
		if err != nil {
			return err
		}
		if snat == nil {
			return fmt.Errorf("Snat ", name, " not destroyed.")
		}
	}
	return nil
}
