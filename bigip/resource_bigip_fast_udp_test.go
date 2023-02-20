/*
Copyright 2022 F5 Networks Inc.
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

func TestAccFastUDPAppCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastUDPAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastUDPAppConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists("fast_udp_app", "fast_udp_tenant", true),
					resource.TestCheckResourceAttr("bigip_fast_udp_app.fast_udp_app", "application", "fast_udp_app"),
					resource.TestCheckResourceAttr("bigip_fast_udp_app.fast_udp_app", "tenant", "fast_udp_tenant"),
					resource.TestCheckResourceAttr("bigip_fast_udp_app.fast_udp_app", "virtual_server.0.ip", "10.99.11.88"),
					resource.TestCheckResourceAttr("bigip_fast_udp_app.fast_udp_app", "virtual_server.0.port", "80"),
				),
			},
		},
	})
}

func getFastUDPAppConfig() string {
	return fmt.Sprintf(`
resource "bigip_fast_udp_app" "fast_udp_app" {
  application = "%v"
  tenant      = "%v"
  virtual_server {
    ip   = "10.99.11.88"
    port = 80
  }
}
`, "fast_udp_app", "fast_udp_tenant")
}

func testCheckFastUDPAppDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_fast_udp_app" {
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
