/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_ROUTE_NAME = fmt.Sprintf("/%s/test-route", TEST_PARTITION)

var TEST_ROUTE_RESOURCE = `

resource "bigip_net_vlan" "test-vlan" {
	name = "` + TEST_VLAN_NAME + `"
	tag = 101
	interfaces {
		vlanport = 1.1
		tagged = true
	}
}
resource "bigip_net_selfip" "test-selfip" {
	name = "` + TEST_SELFIP_NAME + `"
	ip = "11.1.1.1/24"
	vlan = "/Common/test-vlan"
	depends_on = ["bigip_net_vlan.test-vlan"]
}
resource "bigip_net_route" "test-route" {
	  name = "` + TEST_ROUTE_NAME + `"
	  network = "10.10.10.0/24"
	  gw      = "11.1.1.2"
	  depends_on = ["bigip_net_selfip.test-selfip"]
}
`
var TEST_ROUTE_RESOURCE_UPDATE = `

resource "bigip_net_vlan" "test-vlan" {
        name = "` + TEST_VLAN_NAME + `"
        tag = 101
        interfaces {
                vlanport = 1.1
                tagged = true
        }
}
resource "bigip_net_selfip" "test-selfip" {
        name = "` + TEST_SELFIP_NAME + `"
        ip = "11.1.1.1/24"
        vlan = "/Common/test-vlan"
        depends_on = ["bigip_net_vlan.test-vlan"]
}
resource "bigip_net_route" "test-route" {
          name = "` + TEST_ROUTE_NAME + `"
          network = "10.10.10.0/24"
          gw      = "11.1.1.3"
          depends_on = ["bigip_net_selfip.test-selfip"]
}
`

func TestAccBigipNetroute_create(t *testing.T) {
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
					resource.TestCheckResourceAttr("bigip_net_route.test-route", "name", "/Common/test-route"),
					resource.TestCheckResourceAttr("bigip_net_route.test-route", "network", "10.10.10.0/24"),
					resource.TestCheckResourceAttr("bigip_net_route.test-route", "gw", "11.1.1.2"),
				),
			},
		},
	})
}
func TestAccBigipNetroute_update(t *testing.T) {
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
					resource.TestCheckResourceAttr("bigip_net_route.test-route", "gw", "11.1.1.2"),
				),
			},
			{
				Config: TEST_ROUTE_RESOURCE_UPDATE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_net_route.test-route", "gw", "11.1.1.3"),
				),
			},
		},
	})
}
func TestAccBigipNetroute_import(t *testing.T) {
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

		p, err := client.GetRoute(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("route %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("route %s still exists.", name)
		}
		return nil
	}
}

func testCheckroutesDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_net_route" {
			continue
		}

		name := rs.Primary.ID
		route, err := client.GetRoute(name)
		if err != nil {
			return err
		}
		if route != nil {
			return fmt.Errorf("route %s not destroyed.", name)
		}
	}
	return nil
}
