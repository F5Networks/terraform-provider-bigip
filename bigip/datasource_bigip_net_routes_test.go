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

func TestAccDataSourceBigipNetRoutes_basic(t *testing.T) {
	vlanFullPath := fmt.Sprintf("/%s/test-routes-ds-vlan", TestPartition)
	selfipFullPath := fmt.Sprintf("/%s/test-routes-ds-selfip", TestPartition)
	routeFullPath := fmt.Sprintf("/%s/test-routes-ds", TestPartition)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckroutesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipNetRoutesConfig(vlanFullPath, selfipFullPath, routeFullPath),
				Check: resource.ComposeTestCheckFunc(
					testCheckrouteExists(routeFullPath, true),
					resource.TestCheckResourceAttrSet("data.bigip_net_routes.all", "routes.#"),
					resource.TestCheckTypeSetElemNestedAttrs("data.bigip_net_routes.all", "routes.*", map[string]string{
						"full_path": routeFullPath,
						"network":   "10.20.30.0/24",
						"gw":        "13.1.1.2",
					}),
				),
			},
		},
	})
}

func testAccDataSourceBigipNetRoutesConfig(vlanFullPath, selfipFullPath, routeFullPath string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-routes-ds-vlan" {
  name = "%s"
  tag  = 1097
  interfaces {
    vlanport = 1.1
    tagged   = true
  }
}

resource "bigip_net_selfip" "test-routes-ds-selfip" {
  name          = "%s"
  ip            = "13.1.1.1/24"
  vlan          = bigip_net_vlan.test-routes-ds-vlan.name
  port_lockdown = ["none"]
  depends_on    = [bigip_net_vlan.test-routes-ds-vlan]
}

resource "bigip_net_route" "test-routes-ds" {
  name       = "%s"
  network    = "10.20.30.0/24"
  gw         = "13.1.1.2"
  depends_on = [bigip_net_selfip.test-routes-ds-selfip]
}

data "bigip_net_routes" "all" {
  depends_on = [bigip_net_route.test-routes-ds]
}
`, vlanFullPath, selfipFullPath, routeFullPath)
}
