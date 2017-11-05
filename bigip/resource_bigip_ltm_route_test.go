package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_ROUTE_NAME = fmt.Sprintf("/%s/test-route", TEST_PARTITION)

var TEST_ROUTE_RESOURCE = `

resource "bigip_ltm_vlan" "test-vlan" {
	name = "` + TEST_VLAN_NAME + `"
	tag = 101
	interfaces = {
		vlanport = 1.2,
		tagged = false
	}
}

resource "bigip_ltm_selfip" "test-selfip" {
	name = "` + TEST_SELFIP_NAME + `"
	ip = "11.1.1.1/24"
	vlan = "/Common/test-vlan"
	depends_on = ["bigip_ltm_vlan.test-vlan"]
		}

resource "bigip_route" "test-route" {
  name = "/Common/test-route"
  network = "10.10.10.0/24"
  gw      = "11.1.1.2"
	depends_on = ["bigip_ltm_selfip.test-selfip"]
}
`

func TestBigipLtmroute_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckroutesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ROUTE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckrouteExists(TEST_ROUTE_NAME, true),
					resource.TestCheckResourceAttr("bigip_route.test-route", "name", TEST_ROUTE_NAME),
					resource.TestCheckResourceAttr("bigip_route.test-route", "network", "10.10.10.0/24"),
					resource.TestCheckResourceAttr("bigip_route.test-route", "gw", "11.1.1.2"),
				),
			},
		},
	})
}

func TestBigipLtmroute_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckroutesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ROUTE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckrouteExists(TEST_ROUTE_NAME, true),
				),
				ResourceName:      TEST_ROUTE_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckrouteExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.Routes()
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("route ", name, " was not created.")
		}
		if !exists && p != nil {
			return fmt.Errorf("route ", name, " still exists.")
		}
		return nil
	}
}

func testCheckroutesDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_route" {
			continue
		}

		name := rs.Primary.ID
		route, err := client.Routes()
		if err != nil {
			return err
		}
		if route == nil {
			return fmt.Errorf("route ", name, " not destroyed.")
		}
	}
	return nil
}
