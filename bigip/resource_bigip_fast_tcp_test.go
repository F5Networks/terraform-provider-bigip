/*
Copyright 2021 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var appName = "fast_tcp_app"
var tenantName = "fast_tcp_tenant"

var cfg1 = fmt.Sprintf(`
resource "bigip_fast_tcp_app" "fast_tcp_app" {
  application = "%v"
  tenant      = "%v"
  virtual_server {
    ip   = "10.99.11.88"
    port = 80
  }
  existing_monitor = "/Common/tcp"
}
`, appName, tenantName)

var cfg2 = fmt.Sprintf(`
resource "bigip_fast_tcp_app" "fast_tcp_app" {
  application = "%v"
  tenant      = "%v"
  virtual_server {
    ip   = "10.99.11.88"
    port = 80
  }
  monitor {
	interval = 40
  }
}
`, appName, tenantName)

func TestAccFastTCPAppCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastTCPAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: cfg1,
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(appName, tenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "application", "fast_tcp_app"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "tenant", "fast_tcp_tenant"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "virtual_server.0.ip", "10.99.11.88"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "virtual_server.0.port", "80"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "existing_monitor", "/Common/tcp"),
				),
			},
			{
				Config: cfg2,
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(appName, tenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "application", "fast_tcp_app"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "tenant", "fast_tcp_tenant"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "virtual_server.0.ip", "10.99.11.88"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "virtual_server.0.port", "80"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "monitor.0.interval", "40"),
				),
			},
		},
	})
}

func testCheckFastTCPAppDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_fast_tcp_app" {
			continue
		}
		name := rs.Primary.ID
		template, err := client.GetFastApp(tenant, name)
		if err != nil {
			return err
		}
		if template != "" {
			return fmt.Errorf("Fast Application  %s not destroyed.", name)
		}
	}
	return nil
}
