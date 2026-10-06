/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBigipNetRoute_basic(t *testing.T) {
	vlanFullPath := fmt.Sprintf("/%s/test-vlan-route-ds", TestPartition)
	selfFullPath := fmt.Sprintf("/%s/test-self-route-ds", TestPartition)
	routeFullPath := fmt.Sprintf("/%s/test-route-ds", TestPartition)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckrouteExists(routeFullPath, false),
			testCheckselfipsDestroyed,
			testCheckvlansDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipNetRouteConfig(vlanFullPath, selfFullPath, routeFullPath),
				Check: resource.ComposeTestCheckFunc(
					testCheckrouteExists(routeFullPath, true),
					resource.TestCheckResourceAttr("data.bigip_net_route.test", "full_path", routeFullPath),
					resource.TestCheckResourceAttr("data.bigip_net_route.test", "network", "10.50.50.0/24"),
					resource.TestCheckResourceAttr("data.bigip_net_route.test", "gw", "11.4.1.254"),
					resource.TestCheckResourceAttr("data.bigip_net_route.test", "description", "route created for data source acceptance test"),
				),
			},
		},
	})
}

func testAccDataSourceBigipNetRouteConfig(vlanFullPath, selfFullPath, routeFullPath string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-vlan-route-ds" {
  name = "%[1]s"
  tag  = 1095
  interfaces {
    vlanport = 1.1
    tagged   = true
  }
}

resource "bigip_net_selfip" "test-self-route-ds" {
  name          = "%[2]s"
  ip            = "11.4.1.1/24"
  vlan          = bigip_net_vlan.test-vlan-route-ds.name
  port_lockdown = ["none"]
  depends_on    = [bigip_net_vlan.test-vlan-route-ds]
}

resource "bigip_net_route" "test-route-ds" {
  name       = "%[3]s"
  network    = "10.50.50.0/24"
  gw         = "11.4.1.254"
  depends_on = [bigip_net_selfip.test-self-route-ds]
}

# bigip_net_route has no description argument; set it directly so the data
# source's description field has something non-empty to read back.
resource "bigip_command" "set_route_description" {
  commands   = ["modify net route %[3]s description \"route created for data source acceptance test\""]
  when       = "apply"
  depends_on = [bigip_net_route.test-route-ds]
}

data "bigip_net_route" "test" {
  name       = bigip_net_route.test-route-ds.name
  depends_on = [bigip_command.set_route_description]
}
`, vlanFullPath, selfFullPath, routeFullPath)
}
