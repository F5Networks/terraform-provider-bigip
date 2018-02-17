package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_SELFIP_NAME = fmt.Sprintf("/%s/test-selfip", TEST_PARTITION)

var TEST_SELFIP_RESOURCE = `

resource "bigip_net_vlan" "test-vlan" {
	name = "` + TEST_VLAN_NAME + `"
	tag = 101
	interfaces = {
		vlanport = 1.1,
		tagged = false
	}
}

resource "bigip_net_selfip" "test-selfip" {
	name = "/Common/test-selfip"
	ip = "11.1.1.1/24"
	vlan = "/Common/test-vlan"
	depends_on = ["bigip_net_vlan.test-vlan"]
		}
`

func TestAccBigipNetselfip_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckselfipsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SELFIP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME, true),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "name", "/Common/test-selfip"),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "ip", "11.1.1.1/24"),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "vlan", "/Common/test-vlan"),
				),
			},
		},
	})
}

func TestAccBigipNetselfip_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckselfipsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SELFIP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME, true),
				),
				ResourceName:      TEST_SELFIP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckselfipExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.SelfIPs()
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("selfip %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("selfip %s still exists.", name)
		}
		return nil
	}
}

func testCheckselfipsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_net_selfip" {
			continue
		}

		name := rs.Primary.ID
		selfip, err := client.SelfIPs()
		if err != nil {
			return err
		}
		if selfip == nil {
			return fmt.Errorf("selfip %s not destroyed.", name)
		}
	}
	return nil
}
