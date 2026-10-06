/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceBigipNetInterfaces_basic exercises the data source
// against a real device. There's no bigip_net_interface resource (physical
// interfaces can't be created/destroyed), so this only verifies the read
// succeeds and returns the interfaces that exist on the test device.
func TestAccDataSourceBigipNetInterfaces_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "bigip_net_interfaces" "all" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bigip_net_interfaces.all", "interfaces.#"),
				),
			},
		},
	})
}
