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

func TestAccDataSourceBigipNetSelfips_basic(t *testing.T) {
	vlanFullPath := fmt.Sprintf("/%s/test-selfips-ds-vlan", TestPartition)
	selfipFullPath := fmt.Sprintf("/%s/test-selfips-ds", TestPartition)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckselfipsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipNetSelfipsConfig(vlanFullPath, selfipFullPath),
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(selfipFullPath),
					resource.TestCheckResourceAttrSet("data.bigip_net_selfips.all", "self_ips.#"),
					resource.TestCheckTypeSetElemNestedAttrs("data.bigip_net_selfips.all", "self_ips.*", map[string]string{
						"full_path": selfipFullPath,
						"address":   "12.1.1.1/24",
						"vlan":      vlanFullPath,
					}),
				),
			},
		},
	})
}

func testAccDataSourceBigipNetSelfipsConfig(vlanFullPath, selfipFullPath string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-selfips-ds-vlan" {
  name = "%s"
  tag  = 1098
  interfaces {
    vlanport = 1.1
    tagged   = true
  }
}

resource "bigip_net_selfip" "test-selfips-ds" {
  name          = "%s"
  ip            = "12.1.1.1/24"
  vlan          = bigip_net_vlan.test-selfips-ds-vlan.name
  port_lockdown = ["none"]
  depends_on    = [bigip_net_vlan.test-selfips-ds-vlan]
}

data "bigip_net_selfips" "all" {
  depends_on = [bigip_net_selfip.test-selfips-ds]
}
`, vlanFullPath, selfipFullPath)
}
