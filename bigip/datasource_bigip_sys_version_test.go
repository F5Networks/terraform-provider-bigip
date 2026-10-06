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

// TestAccDataSourceBigipSysVersion_basic exercises the data source against
// a real device. sys/version (and sys/hardware, for the platform field) are
// read-only system information with nothing to create/destroy, so this
// only asserts that the expected fields come back non-empty rather than
// pinning exact version/build/platform strings, which will vary by
// whatever device runs the acceptance suite.
func TestAccDataSourceBigipSysVersion_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipSysVersionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bigip_sys_version.test", "version"),
					resource.TestCheckResourceAttrSet("data.bigip_sys_version.test", "build"),
					resource.TestCheckResourceAttr("data.bigip_sys_version.test", "product", "BIG-IP"),
					resource.TestCheckResourceAttrSet("data.bigip_sys_version.test", "edition"),
					resource.TestCheckResourceAttrSet("data.bigip_sys_version.test", "platform"),
				),
			},
		},
	})
}

func testAccDataSourceBigipSysVersionConfig() string {
	return `
data "bigip_sys_version" "test" {}
`
}
