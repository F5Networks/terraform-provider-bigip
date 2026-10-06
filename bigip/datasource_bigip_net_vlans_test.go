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

func TestAccDataSourceBigipNetVlans_basic(t *testing.T) {
	vlanFullPath := fmt.Sprintf("/%s/test-vlans-ds", TestPartition)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckvlansDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipNetVlansConfig(vlanFullPath),
				Check: resource.ComposeTestCheckFunc(
					testCheckvlanExists(vlanFullPath, true),
					resource.TestCheckResourceAttrSet("data.bigip_net_vlans.all", "vlans.#"),
					resource.TestCheckTypeSetElemNestedAttrs("data.bigip_net_vlans.all", "vlans.*", map[string]string{
						"full_path":           vlanFullPath,
						"tag":                 "1099",
						"mtu":                 "1400",
						"interfaces.#":        "1",
						"interfaces.0.name":   "1.1",
						"interfaces.0.tagged": "true",
					}),
				),
			},
		},
	})
}

func testAccDataSourceBigipNetVlansConfig(vlanFullPath string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-vlans-ds" {
  name = "%s"
  tag  = 1099
  mtu  = 1400
  interfaces {
    vlanport = 1.1
    tagged   = true
  }
}

data "bigip_net_vlans" "all" {
  depends_on = [bigip_net_vlan.test-vlans-ds]
}
`, vlanFullPath)
}
