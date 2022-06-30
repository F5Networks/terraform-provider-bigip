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

var httpsAppName = "fast_https_app"
var httpsTenantName = "fast_https_tenant"

func TestAccFastHTTPSAppCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastHTTPSAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastHTTPSAppConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(httpsAppName, httpsTenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "application", "fast_https_app"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "tenant", "fast_https_tenant"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.ip", "10.30.40.44"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.port", "443"),
				),
			},
		},
	})
}

func getFastHTTPSAppConfig() string {
	return fmt.Sprintf(`
resource "bigip_fast_https_app" "fast_https_app" {
  tenant      = "%v"
  application = "%v"
  virtual_server {
    ip   = "10.30.40.44"
    port = 443
  }
}
`, httpsTenantName, httpsAppName)
}

func testCheckFastHTTPSAppDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_fast_https_app" {
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
