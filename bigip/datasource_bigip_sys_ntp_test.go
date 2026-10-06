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

// TestAccDataSourceBigipSysNtp_basic exercises the data source against a
// real device. sys/ntp is a system singleton (not created/destroyed), so
// this configures it via the bigip_sys_ntp resource and reads it back
// through the data source; there is no CheckDestroy since the underlying
// object always exists on the device and has no Delete API.
func TestAccDataSourceBigipSysNtp_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipSysNtpConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckntpExists("test-ntp-ds", true),
					resource.TestCheckResourceAttr("data.bigip_sys_ntp.test", "description", "test-ntp-ds"),
					resource.TestCheckResourceAttr("data.bigip_sys_ntp.test", "ntp_servers.#", "2"),
					resource.TestCheckResourceAttr("data.bigip_sys_ntp.test", "ntp_servers.0", "0.pool.ntp.org"),
					resource.TestCheckResourceAttr("data.bigip_sys_ntp.test", "ntp_servers.1", "1.pool.ntp.org"),
					resource.TestCheckResourceAttr("data.bigip_sys_ntp.test", "timezone", "America/Los_Angeles"),
				),
			},
		},
	})
}

func testAccDataSourceBigipSysNtpConfig() string {
	return `
resource "bigip_sys_ntp" "test-ntp-ds" {
  description = "test-ntp-ds"
  servers     = ["0.pool.ntp.org", "1.pool.ntp.org"]
  timezone    = "America/Los_Angeles"
}

data "bigip_sys_ntp" "test" {
  depends_on = [bigip_sys_ntp.test-ntp-ds]
}
`
}
