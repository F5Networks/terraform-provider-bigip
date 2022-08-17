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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var appName = "fast_tcp_app"
var tenantName = "fast_tcp_tenant"

func TestAccFastTCPAppCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastTCPAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastTCPAppConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(appName, tenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "application", "fast_tcp_app"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "tenant", "fast_tcp_tenant"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "virtual_server.0.ip", "10.99.11.88"),
					resource.TestCheckResourceAttr("bigip_fast_tcp_app.fast_tcp_app", "virtual_server.0.port", "80"),
				),
			},
		},
	})
}

func getFastTCPAppConfig() string {
	return fmt.Sprintf(`
resource "bigip_fast_tcp_app" "fast_tcp_app" {
  application = "%v"
  tenant      = "%v"
  virtual_server {
    ip   = "10.99.11.88"
    port = 80
  }
}
`, appName, tenantName)
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
