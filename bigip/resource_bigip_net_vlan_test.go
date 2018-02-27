package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_VLAN_NAME = fmt.Sprintf("/%s/test-vlan", TEST_PARTITION)

var TEST_VLAN_RESOURCE = `
resource "bigip_net_vlan" "test-vlan" {
	name = "/Common/test-vlan"
	tag = 101
	interfaces = {
		vlanport = 1.1,
		tagged = false
	}
}
`

func TestAccBigipNetvlan_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckvlansDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_VLAN_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckvlanExists(TEST_VLAN_NAME, true),
					resource.TestCheckResourceAttr("bigip_net_vlan.test-vlan", "name", "/Common/test-vlan"),
					resource.TestCheckResourceAttr("bigip_net_vlan.test-vlan", "tag", "101"),
					resource.TestCheckResourceAttr("bigip_net_vlan.test-vlan", "interfaces.0.vlanport", "1.1"),
					resource.TestCheckResourceAttr("bigip_net_vlan.test-vlan", "interfaces.0.tagged", "false"),
				),
			},
		},
	})
}

func TestAccBigipNetvlan_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckvlansDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_VLAN_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckvlanExists(TEST_VLAN_NAME, true),
				),
				ResourceName:      TEST_VLAN_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckvlanExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.Vlans()
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("vlan %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("vlan %s still exists.", name)
		}
		return nil
	}
}

func testCheckvlansDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_net_vlan" {
			continue
		}

		name := rs.Primary.ID
		vlan, err := client.Vlans()
		if err != nil {
			return err
		}
		if vlan == nil {
			return fmt.Errorf("vlan %s not destroyed.", name)
		}
	}
	return nil
}
