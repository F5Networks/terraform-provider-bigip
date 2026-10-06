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

// TestAccDataSourceBigipSysDns_basic exercises the data source against a
// real device. sys/dns is a system singleton (not created/destroyed), so
// this configures it via the bigip_sys_dns resource and reads it back
// through the data source; there is no CheckDestroy since the underlying
// object always exists on the device.
func TestAccDataSourceBigipSysDns_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipSysDnsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists("test-dns-ds", true),
					resource.TestCheckResourceAttr("data.bigip_sys_dns.test", "description", "test-dns-ds"),
					resource.TestCheckResourceAttr("data.bigip_sys_dns.test", "name_servers.#", "2"),
					resource.TestCheckResourceAttr("data.bigip_sys_dns.test", "name_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.bigip_sys_dns.test", "name_servers.1", "2.2.2.2"),
					resource.TestCheckResourceAttr("data.bigip_sys_dns.test", "search.#", "2"),
					resource.TestCheckResourceAttr("data.bigip_sys_dns.test", "search.0", "f5.com"),
					resource.TestCheckResourceAttr("data.bigip_sys_dns.test", "search.1", "f5.net"),
				),
			},
		},
	})
}

func testAccDataSourceBigipSysDnsConfig() string {
	return `
resource "bigip_sys_dns" "test-dns-ds" {
  description  = "test-dns-ds"
  name_servers = ["1.1.1.1", "2.2.2.2"]
  search       = ["f5.com", "f5.net"]
}

data "bigip_sys_dns" "test" {
  depends_on = [bigip_sys_dns.test-dns-ds]
}
`
}
