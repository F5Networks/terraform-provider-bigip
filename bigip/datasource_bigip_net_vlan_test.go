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

func TestAccDataSourceBigipNetVlan_byName(t *testing.T) {
	vlanFullPath := fmt.Sprintf("/%s/test-vlan-ds-name", TestPartition)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckvlansDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipNetVlanByNameConfig(vlanFullPath),
				Check: resource.ComposeTestCheckFunc(
					testCheckvlanExists(vlanFullPath, true),
					resource.TestCheckResourceAttr("data.bigip_net_vlan.by_name", "full_path", vlanFullPath),
					resource.TestCheckResourceAttr("data.bigip_net_vlan.by_name", "vlan_id", "1098"),
					resource.TestCheckResourceAttr("data.bigip_net_vlan.by_name", "tagged_interfaces.#", "1"),
					resource.TestCheckResourceAttr("data.bigip_net_vlan.by_name", "tagged_interfaces.0", "1.1"),
					resource.TestCheckResourceAttr("data.bigip_net_vlan.by_name", "untagged_interfaces.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceBigipNetVlan_byVlanID(t *testing.T) {
	vlanFullPath := fmt.Sprintf("/%s/test-vlan-ds-id", TestPartition)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckvlansDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipNetVlanByVlanIDConfig(vlanFullPath),
				Check: resource.ComposeTestCheckFunc(
					testCheckvlanExists(vlanFullPath, true),
					resource.TestCheckResourceAttr("data.bigip_net_vlan.by_id", "name", vlanFullPath),
					resource.TestCheckResourceAttr("data.bigip_net_vlan.by_id", "full_path", vlanFullPath),
					resource.TestCheckResourceAttr("data.bigip_net_vlan.by_id", "tagged_interfaces.#", "1"),
					resource.TestCheckResourceAttr("data.bigip_net_vlan.by_id", "tagged_interfaces.0", "1.1"),
				),
			},
		},
	})
}

func testAccDataSourceBigipNetVlanByNameConfig(vlanFullPath string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-vlan-ds-name" {
  name = "%s"
  tag  = 1098
  interfaces {
    vlanport = 1.1
    tagged   = true
  }
}

data "bigip_net_vlan" "by_name" {
  name = bigip_net_vlan.test-vlan-ds-name.name
}
`, vlanFullPath)
}

func testAccDataSourceBigipNetVlanByVlanIDConfig(vlanFullPath string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-vlan-ds-id" {
  name = "%s"
  tag  = 1097
  interfaces {
    vlanport = 1.1
    tagged   = true
  }
}

data "bigip_net_vlan" "by_id" {
  vlan_id    = bigip_net_vlan.test-vlan-ds-id.tag
  depends_on = [bigip_net_vlan.test-vlan-ds-id]
}
`, vlanFullPath)
}
