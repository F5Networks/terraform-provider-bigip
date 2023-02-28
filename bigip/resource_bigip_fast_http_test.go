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

var httpAppName = "fast_http_app"
var httpTenantName = "fast_http_tenant"

func TestAccFastHTTPAppCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastHTTPAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastHTTPAppConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(httpAppName, httpTenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app", "application", "fast_http_app"),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app", "tenant", "fast_http_tenant"),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app", "virtual_server.0.ip", "10.30.30.44"),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app", "virtual_server.0.port", "443"),
				),
			},
		},
	})
}

func TestAccFastHTTPAppCreateTC02(t *testing.T) {
	var httpApp1Name = "fast_http_apptc2"
	var httpTenant1Name = "fast_http_tenanttc2"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastHTTPAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastHTTPAppConfigTC02(httpTenant1Name, httpApp1Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(httpApp1Name, httpTenant1Name, true),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app_tc2", "application", httpApp1Name),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app_tc2", "tenant", httpTenant1Name),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app_tc2", "virtual_server.0.ip", "10.200.21.2"),
					resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app_tc2", "virtual_server.0.port", "443"),
					// resource.TestCheckResourceAttr("bigip_fast_http_app.fast_http_app", "endpoint_ltm_policy.0", "/Common/testpolicy1"),
				),
			},
		},
	})
}

func getFastHTTPAppConfig() string {
	return fmt.Sprintf(`
resource "bigip_fast_http_app" "fast_http_app" {
  tenant      = "%v"
  application = "%v"
  virtual_server {
    ip   = "10.30.30.44"
    port = 443
  }
}
`, httpTenantName, httpAppName)
}

func getFastHTTPAppConfigTC02(httpTenantName, httpAppName string) string {
	return fmt.Sprintf(`
resource "bigip_fast_http_app" "fast_http_app_tc2" {
  tenant      = "%v"
  application = "%v"
  virtual_server {
    ip   = "10.200.21.2"
    port = 443
  }
  pool_members {
    addresses = ["10.1.20.120", "10.1.10.121", "10.1.10.122"]
    port      = 80
  }
  load_balancing_mode = "least-connections-member"
  endpoint_ltm_policy = ["/Common/testpolicy1"]
}
`, httpTenantName, httpAppName)
}

func testCheckFastHTTPAppDestroyed(s *terraform.State) error {
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
