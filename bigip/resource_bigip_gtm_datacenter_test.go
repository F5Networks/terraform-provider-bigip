package bigip

import (
	"fmt"
	"testing"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_DATACENTER_NAME = fmt.Sprintf("/%s/test-datacenter", TEST_PARTITION)

var TEST_DATACENTER_RESOURCE = `
resource "bigip_gtm_datacenter" "test-datacenter"

{
name = "/Common/test-datacenter"
description = "This is DC located in San jose"
contact = "shitole"
enabled = true
}
`

func TestBigipGtmdatacenter_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: testCheckgtmdatacentersDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DATACENTER_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckgtmdatacenterExists(TEST_DATACENTER_NAME, true),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "name","/Common/test-datacenter"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "description", "This is DC located in San jose"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "contact", "shitole"),
					resource.TestCheckResourceAttr("bigip_gtm_datacenter.test-datacenter", "enabled", "true"),

				),
			},
		},
	})
}

func TestBigipGtmdatacenter_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: testCheckgtmdatacentersDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DATACENTER_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckgtmdatacenterExists(TEST_DATACENTER_NAME, true),
				),
				ResourceName:      TEST_DATACENTER_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckgtmdatacenterExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.Datacenters()
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("gtmdatacenter %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("gtmdatacenter %s still exists.", name)
		}
		return nil
	}
}

func testCheckgtmdatacentersDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_datacenter" {
			continue
		}

		name := rs.Primary.ID
		gtmdatacenter, err := client.Datacenters()
		if err != nil {
			return err
		}
		if gtmdatacenter == nil {
			return fmt.Errorf("gtmdatacenter ", name, " not destroyed.")
		}
	}
	return nil
}
