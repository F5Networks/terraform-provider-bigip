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

var httpAppName = "fast_http_app"
var httpTenantName = "fast_http_tenant"

func TestAccFastHTTPAppCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastTCPAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastHTTPAppConfig(""),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(httpAppName, httpTenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app", "application", "fast_http_app"),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app", "tenant", "fast_http_tenant"),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app", "virtual_server.ip", "10.30.30.44"),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app", "virtual_server.port", "443"),
				),
			},
		},
	})
}

func getFastHTTPAppConfig(attrs string) string {
	return fmt.Sprintf(`
resource "bigip_fast_http_app" "fast-http" {
  tenant = "%v"
  application= "%v"
  virtual_server = {
    ip = "10.30.30.44"
    port = 443
  }
}
`, httpTenantName, httpAppName)
}

func testCheckFastTCPAppDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_fast_http_app" {
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