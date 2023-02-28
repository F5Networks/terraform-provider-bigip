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

func TestAccFastHTTPSAppSSLProfileTC1(t *testing.T) {
	httpsAppName = "fast_https_apptc1"
	httpsTenantName = "fast_https_tenanttc1"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastHTTPSAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastHTTPSAppSSLConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(httpsAppName, httpsTenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "application", "fast_https_apptc1"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "tenant", "fast_https_tenanttc1"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.ip", "10.30.40.44"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.port", "443"),
				),
			},
		},
	})
}

func TestAccFastHTTPSAppSSLProfileTC2(t *testing.T) {
	httpsAppName = "fast_https_apptc2"
	httpsTenantName = "fast_https_tenanttc2"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastHTTPSAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastHTTPSAppSSLConfigTC2(),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(httpsAppName, httpsTenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "application", "fast_https_apptc2"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "tenant", "fast_https_tenanttc2"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.ip", "10.30.40.44"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.port", "443"),
				),
			},
		},
	})
}

func TestAccFastHTTPSAppProfileTC3(t *testing.T) {
	httpsAppName = "fast_https_apptc3"
	httpsTenantName = "fast_https_tenanttc3"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastHTTPSAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastHTTPSAppConfigTC3(),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(httpsAppName, httpsTenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "application", "fast_https_apptc3"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "tenant", "fast_https_tenanttc3"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.ip", "10.30.40.45"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.port", "443"),
				),
			},
		},
	})
}

func TestAccFastHTTPSAppProfileTC4(t *testing.T) {
	httpsAppName = "fast_https_apptc4"
	httpsTenantName = "fast_https_tenanttc4"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastHTTPSAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastHTTPSAppConfigTC4(),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(httpsAppName, httpsTenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "application", "fast_https_apptc4"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "tenant", "fast_https_tenanttc4"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.ip", "10.30.41.45"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.port", "443"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "endpoint_ltm_policy.0", "/Common/testpolicy1"),
				),
			},
		},
	})
}

func TestAccFastHTTPSAppProfileTC5(t *testing.T) {
	httpsAppName = "fast_https_apptc5"
	httpsTenantName = "fast_https_tenanttc5"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastHTTPSAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getFastHTTPSAppConfigTC5(),
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists(httpsAppName, httpsTenantName, true),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "application", "fast_https_apptc5"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "tenant", "fast_https_tenanttc5"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.ip", "10.30.41.45"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "virtual_server.0.port", "443"),
					resource.TestCheckResourceAttr("bigip_fast_https_app.fast_https_app", "endpoint_ltm_policy.0", "/Common/testpolicy1"),
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

func getFastHTTPSAppSSLConfig() string {
	return fmt.Sprintf(`
resource "bigip_fast_https_app" "fast_https_app" {
  tenant      = "%v"
  application = "%v"
  virtual_server {
    ip   = "10.30.40.44"
    port = 443
  }
  tls_server_profile {
    tls_cert_name = "/Common/default.crt"
    tls_key_name  = "/Common/default.key"
  }
}
`, httpsTenantName, httpsAppName)
}

func getFastHTTPSAppSSLConfigTC2() string {
	return fmt.Sprintf(`
resource "bigip_fast_https_app" "fast_https_app" {
  tenant      = "%v"
  application = "%v"
  virtual_server {
    ip   = "10.30.40.44"
    port = 443
  }
  tls_client_profile {
    tls_cert_name = "/Common/default.crt"
    tls_key_name  = "/Common/default.key"
  }
}
`, httpsTenantName, httpsAppName)
}

func getFastHTTPSAppConfigTC3() string {
	return fmt.Sprintf(`
resource "bigip_fast_https_app" "fast_https_app" {
  tenant      = "%v"
  application = "%v"
  virtual_server {
    ip   = "10.30.40.45"
    port = 443
  }
  tls_server_profile {
    tls_cert_name = "/Common/default.crt"
    tls_key_name  = "/Common/default.key"
  }
  tls_client_profile {
    tls_cert_name = "/Common/default.crt"
    tls_key_name  = "/Common/default.key"
  }
  waf_security_policy {
    enable = true
  }
}
`, httpsTenantName, httpsAppName)
}

func getFastHTTPSAppConfigTC4() string {
	return fmt.Sprintf(`
resource "bigip_fast_https_app" "fast_https_app" {
  tenant      = "%v"
  application = "%v"
  virtual_server {
    ip   = "10.30.41.45"
    port = 443
  }
  tls_server_profile {
    tls_cert_name = "/Common/default.crt"
    tls_key_name  = "/Common/default.key"
  }
  tls_client_profile {
    tls_cert_name = "/Common/default.crt"
    tls_key_name  = "/Common/default.key"
  }
  waf_security_policy {
    enable = true
  }
  endpoint_ltm_policy = ["/Common/testpolicy1"]
}
`, httpsTenantName, httpsAppName)
}

func getFastHTTPSAppConfigTC5() string {
	return fmt.Sprintf(`
resource "bigip_fast_https_app" "fast_https_app" {
  tenant      = "%v"
  application = "%v"
  virtual_server {
    ip   = "10.30.41.45"
    port = 443
  }
  tls_server_profile {
    tls_cert_name = "/Common/default.crt"
    tls_key_name  = "/Common/default.key"
  }
  tls_client_profile {
    tls_cert_name = "/Common/default.crt"
    tls_key_name  = "/Common/default.key"
  }
  waf_security_policy {
    enable = true
  }
  endpoint_ltm_policy = ["/Common/testpolicy1"]
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
