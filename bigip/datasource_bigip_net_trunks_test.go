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

// TestAccDataSourceBigipNetTrunks_basic exercises the data source against a
// real device. There is no bigip_net_trunk resource to create a fixture
// trunk with (trunk/LAG membership requires real physical interfaces), so
// this only verifies the read succeeds and returns a well-formed (possibly
// empty) list.
func TestAccDataSourceBigipNetTrunks_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipNetTrunksConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bigip_net_trunks.all", "trunks.#"),
				),
			},
		},
	})
}

func testAccDataSourceBigipNetTrunksConfig() string {
	return `data "bigip_net_trunks" "all" {}`
}
