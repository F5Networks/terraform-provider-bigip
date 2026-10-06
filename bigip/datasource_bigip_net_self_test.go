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

func TestAccDataSourceBigipNetSelf_basic(t *testing.T) {
	vlanFullPath := fmt.Sprintf("/%s/test-vlan-self-ds", TestPartition)
	selfFullPath := fmt.Sprintf("/%s/test-self-ds", TestPartition)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckselfipsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipNetSelfConfig(vlanFullPath, selfFullPath),
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(selfFullPath),
					resource.TestCheckResourceAttr("data.bigip_net_self.test", "full_path", selfFullPath),
					resource.TestCheckResourceAttr("data.bigip_net_self.test", "address", "11.2.1.1/24"),
					resource.TestCheckResourceAttr("data.bigip_net_self.test", "vlan", vlanFullPath),
					resource.TestCheckResourceAttr("data.bigip_net_self.test", "traffic_group", "/Common/traffic-group-local-only"),
					resource.TestCheckResourceAttr("data.bigip_net_self.test", "allow_service.#", "1"),
					resource.TestCheckResourceAttr("data.bigip_net_self.test", "allow_service.0", "all"),
				),
			},
		},
	})
}

func testAccDataSourceBigipNetSelfConfig(vlanFullPath, selfFullPath string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-vlan-self-ds" {
  name = "%s"
  tag  = 1096
  interfaces {
    vlanport = 1.1
    tagged   = true
  }
}

resource "bigip_net_selfip" "test-self-ds" {
  name          = "%s"
  ip            = "11.2.1.1/24"
  vlan          = bigip_net_vlan.test-vlan-self-ds.name
  port_lockdown = ["all"]
  depends_on    = [bigip_net_vlan.test-vlan-self-ds]
}

data "bigip_net_self" "test" {
  name       = bigip_net_selfip.test-self-ds.name
  depends_on = [bigip_net_selfip.test-self-ds]
}
`, vlanFullPath, selfFullPath)
}
