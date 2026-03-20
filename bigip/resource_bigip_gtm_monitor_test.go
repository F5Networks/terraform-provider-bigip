/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"strings"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestGtmMonitorHttpResource = `
resource "bigip_gtm_monitor_http" "test-gtm-http-monitor" {
  name                  = "/Common/test-gtm-http-monitor"
  defaults_from         = "/Common/http"
  destination           = "*:*"
  interval              = 30
  timeout               = 120
  probe_timeout         = 5
  ignore_down_response  = "disabled"
  transparent           = "disabled"
  reverse               = "disabled"
  send                  = "GET /\\r\\n"
  receive              = "200 OK"
}
`

var TestGtmMonitorHttpsResource = `
resource "bigip_gtm_monitor_https" "test-gtm-https-monitor" {
  name                  = "/Common/test-gtm-https-monitor"
  defaults_from         = "/Common/https"
  destination           = "*:*"
  interval              = 30
  timeout               = 120
  probe_timeout         = 5
  ignore_down_response  = "disabled"
  transparent           = "disabled"
  reverse               = "disabled"
  send                  = "GET /\\r\\n"
  receive              = "200 OK"
  cipherlist           = "DEFAULT:+SHA:+3DES:+kEDH"
  compatibility        = "enabled"
}
`

var TestGtmMonitorTcpResource = `
resource "bigip_gtm_monitor_tcp" "test-gtm-tcp-monitor" {
  name                  = "/Common/test-gtm-tcp-monitor"
  defaults_from         = "/Common/tcp"
  destination           = "*:*"
  interval              = 30
  timeout               = 120
  probe_timeout         = 5
  ignore_down_response  = "disabled"
  transparent           = "disabled"
  reverse               = "disabled"
}
`

var TestGtmMonitorPostgresqlResource = `
resource "bigip_gtm_monitor_postgresql" "test-gtm-postgresql-monitor" {
  name                  = "/Common/test-gtm-postgresql-monitor"
  defaults_from         = "/Common/postgresql"
  destination           = "*:5432"
  interval              = 30
  timeout               = 120
  probe_timeout         = 5
  ignore_down_response  = "disabled"
  database              = "testdb"
  username              = "testuser"
  debug                 = "no"
}
`

var TestGtmMonitorBigipResource = `
resource "bigip_gtm_monitor_bigip" "test-gtm-bigip-monitor" {
  name                  = "/Common/test-gtm-bigip-monitor"
  defaults_from         = "/Common/bigip"
  destination           = "*:*"
  interval              = 30
  timeout               = 90
  ignore_down_response  = "disabled"
  aggregation_type      = "none"
}
`

var TestGtmMonitorHttpResourceUpdate = `
resource "bigip_gtm_monitor_http" "test-gtm-http-monitor" {
  name                  = "/Common/test-gtm-http-monitor"
  defaults_from         = "/Common/http"
  destination           = "10.1.1.100:8080"
  interval              = 60
  timeout               = 180
  probe_timeout         = 10
  ignore_down_response  = "enabled"
  transparent           = "enabled"
  reverse               = "enabled"
  send                  = "GET /health\\r\\n"
  receive              = "healthy"
}
`

var TestGtmMonitorHttpsResourceUpdate = `
resource "bigip_gtm_monitor_https" "test-gtm-https-monitor" {
  name                  = "/Common/test-gtm-https-monitor"
  defaults_from         = "/Common/https"
  destination           = "10.1.1.100:8443"
  interval              = 45
  timeout               = 180
  probe_timeout         = 10
  ignore_down_response  = "enabled"
  transparent           = "enabled"
  reverse               = "disabled"
  send                  = "GET /api/health\\r\\n"
  receive              = "status: ok"
  cipherlist           = "HIGH:!ADH:!MD5"
  compatibility        = "disabled"
}
`

var TestGtmMonitorTcpResourceMinimal = `
resource "bigip_gtm_monitor_tcp" "test-gtm-tcp-monitor-minimal" {
  name = "/Common/test-gtm-tcp-monitor-minimal"
}
`

var TestGtmMonitorPostgresqlResourceUpdate = `
resource "bigip_gtm_monitor_postgresql" "test-gtm-postgresql-monitor" {
  name                  = "/Common/test-gtm-postgresql-monitor"
  defaults_from         = "/Common/postgresql"
  destination           = "192.168.1.50:5432"
  interval              = 45
  timeout               = 91
  probe_timeout         = 10
  ignore_down_response  = "enabled"
  database              = "production"
  username              = "produser"
  debug                 = "yes"
}
`

func TestAccBigipGtmMonitorHttp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorHttpDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorHttpResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorHttpExists("/Common/test-gtm-http-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "name", "/Common/test-gtm-http-monitor"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "interval", "30"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "timeout", "120"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "send", "GET /\\r\\n"),
				),
			},
		},
	})
}

func TestAccBigipGtmMonitorHttps_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorHttpsResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorHttpsExists("/Common/test-gtm-https-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "name", "/Common/test-gtm-https-monitor"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "defaults_from", "/Common/https"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "interval", "30"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "timeout", "120"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "cipherlist", "DEFAULT:+SHA:+3DES:+kEDH"),
				),
			},
		},
	})
}

func TestAccBigipGtmMonitorTcp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorTcpDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorTcpResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorTcpExists("/Common/test-gtm-tcp-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_tcp.test-gtm-tcp-monitor", "name", "/Common/test-gtm-tcp-monitor"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_tcp.test-gtm-tcp-monitor", "defaults_from", "/Common/tcp"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_tcp.test-gtm-tcp-monitor", "interval", "30"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_tcp.test-gtm-tcp-monitor", "timeout", "120"),
				),
			},
		},
	})
}

func TestAccBigipGtmMonitorPostgresql_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorPostgresqlDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorPostgresqlResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorPostgresqlExists("/Common/test-gtm-postgresql-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "name", "/Common/test-gtm-postgresql-monitor"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "defaults_from", "/Common/postgresql"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "database", "testdb"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "username", "testuser"),
				),
			},
		},
	})
}

func TestAccBigipGtmMonitorBigip_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorBigipDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorBigipResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorBigipExists("/Common/test-gtm-bigip-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_bigip.test-gtm-bigip-monitor", "name", "/Common/test-gtm-bigip-monitor"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_bigip.test-gtm-bigip-monitor", "defaults_from", "/Common/bigip"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_bigip.test-gtm-bigip-monitor", "aggregation_type", "none"),
				),
			},
		},
	})
}

func testCheckGtmMonitorHttpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		monitor, err := client.GetGtmMonitor(name, "http")
		if err != nil {
			return err
		}
		if exists && monitor == nil {
			return fmt.Errorf("GTM HTTP Monitor %s was not created", name)
		}
		if !exists && monitor != nil {
			return fmt.Errorf("GTM HTTP Monitor %s still exists", name)
		}
		return nil
	}
}

func testCheckGtmMonitorHttpDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_monitor_http" {
			continue
		}
		name := rs.Primary.ID
		monitor, err := client.GetGtmMonitor(name, "http")
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return err
		}
		if monitor != nil {
			return fmt.Errorf("GTM HTTP Monitor %s not destroyed", name)
		}
	}
	return nil
}

func testCheckGtmMonitorHttpsExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		monitor, err := client.GetGtmMonitor(name, "https")
		if err != nil {
			return err
		}
		if exists && monitor == nil {
			return fmt.Errorf("GTM HTTPS Monitor %s was not created", name)
		}
		if !exists && monitor != nil {
			return fmt.Errorf("GTM HTTPS Monitor %s still exists", name)
		}
		return nil
	}
}

func testCheckGtmMonitorHttpsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_monitor_https" {
			continue
		}
		name := rs.Primary.ID
		monitor, err := client.GetGtmMonitor(name, "https")
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return err
		}
		if monitor != nil {
			return fmt.Errorf("GTM HTTPS Monitor %s not destroyed", name)
		}
	}
	return nil
}

func testCheckGtmMonitorTcpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		monitor, err := client.GetGtmMonitor(name, "tcp")
		if err != nil {
			return err
		}
		if exists && monitor == nil {
			return fmt.Errorf("GTM TCP Monitor %s was not created", name)
		}
		if !exists && monitor != nil {
			return fmt.Errorf("GTM TCP Monitor %s still exists", name)
		}
		return nil
	}
}

func testCheckGtmMonitorTcpDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_monitor_tcp" {
			continue
		}
		name := rs.Primary.ID
		monitor, err := client.GetGtmMonitor(name, "tcp")
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return err
		}
		if monitor != nil {
			return fmt.Errorf("GTM TCP Monitor %s not destroyed", name)
		}
	}
	return nil
}

func testCheckGtmMonitorPostgresqlExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		monitor, err := client.GetGtmMonitor(name, "postgresql")
		if err != nil {
			return err
		}
		if exists && monitor == nil {
			return fmt.Errorf("GTM PostgreSQL Monitor %s was not created", name)
		}
		if !exists && monitor != nil {
			return fmt.Errorf("GTM PostgreSQL Monitor %s still exists", name)
		}
		return nil
	}
}

func testCheckGtmMonitorPostgresqlDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_monitor_postgresql" {
			continue
		}
		name := rs.Primary.ID
		monitor, err := client.GetGtmMonitor(name, "postgresql")
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return err
		}
		if monitor != nil {
			return fmt.Errorf("GTM PostgreSQL Monitor %s not destroyed", name)
		}
	}
	return nil
}

func testCheckGtmMonitorBigipExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		monitor, err := client.GetGtmMonitor(name, "bigip")
		if err != nil {
			return err
		}
		if exists && monitor == nil {
			return fmt.Errorf("GTM BIG-IP Monitor %s was not created", name)
		}
		if !exists && monitor != nil {
			return fmt.Errorf("GTM BIG-IP Monitor %s still exists", name)
		}
		return nil
	}
}

func testCheckGtmMonitorBigipDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_monitor_bigip" {
			continue
		}
		name := rs.Primary.ID
		monitor, err := client.GetGtmMonitor(name, "bigip")
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return err
		}
		if monitor != nil {
			return fmt.Errorf("GTM BIG-IP Monitor %s not destroyed", name)
		}
	}
	return nil
}

// Update tests - verify resource modifications work correctly

func TestAccBigipGtmMonitorHttp_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorHttpDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorHttpResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorHttpExists("/Common/test-gtm-http-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "interval", "30"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "timeout", "120"),
				),
			},
			{
				Config: TestGtmMonitorHttpResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorHttpExists("/Common/test-gtm-http-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "destination", "10.1.1.100:8080"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "interval", "60"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "timeout", "180"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "probe_timeout", "10"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "ignore_down_response", "enabled"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "transparent", "enabled"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "reverse", "enabled"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "send", "GET /health\\r\\n"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_http.test-gtm-http-monitor", "receive", "healthy"),
				),
			},
		},
	})
}

func TestAccBigipGtmMonitorHttps_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorHttpsResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorHttpsExists("/Common/test-gtm-https-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "interval", "30"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "cipherlist", "DEFAULT:+SHA:+3DES:+kEDH"),
				),
			},
			{
				Config: TestGtmMonitorHttpsResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorHttpsExists("/Common/test-gtm-https-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "destination", "10.1.1.100:8443"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "interval", "45"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "cipherlist", "HIGH:!ADH:!MD5"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_https.test-gtm-https-monitor", "compatibility", "disabled"),
				),
			},
		},
	})
}

func TestAccBigipGtmMonitorPostgresql_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorPostgresqlDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorPostgresqlResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorPostgresqlExists("/Common/test-gtm-postgresql-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "database", "testdb"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "debug", "no"),
				),
			},
			{
				Config: TestGtmMonitorPostgresqlResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorPostgresqlExists("/Common/test-gtm-postgresql-monitor", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "destination", "192.168.1.50:5432"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "database", "production"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "username", "produser"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor", "debug", "yes"),
				),
			},
		},
	})
}

// Import tests - verify resources can be imported correctly

func TestAccBigipGtmMonitorHttp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorHttpDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorHttpResource,
			},
			{
				ResourceName:      "bigip_gtm_monitor_http.test-gtm-http-monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigipGtmMonitorHttps_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorHttpsResource,
			},
			{
				ResourceName:      "bigip_gtm_monitor_https.test-gtm-https-monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigipGtmMonitorTcp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorTcpDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorTcpResource,
			},
			{
				ResourceName:      "bigip_gtm_monitor_tcp.test-gtm-tcp-monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigipGtmMonitorPostgresql_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorPostgresqlDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorPostgresqlResource,
			},
			{
				ResourceName:      "bigip_gtm_monitor_postgresql.test-gtm-postgresql-monitor",
				ImportState:       true,
				ImportStateVerify: true,
				// Password is sensitive and won't be returned in read
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccBigipGtmMonitorBigip_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorBigipDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorBigipResource,
			},
			{
				ResourceName:      "bigip_gtm_monitor_bigip.test-gtm-bigip-monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Minimal configuration tests - verify defaults work correctly

func TestAccBigipGtmMonitorTcp_minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmMonitorTcpDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGtmMonitorTcpResourceMinimal,
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmMonitorTcpExists("/Common/test-gtm-tcp-monitor-minimal", true),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_tcp.test-gtm-tcp-monitor-minimal", "name", "/Common/test-gtm-tcp-monitor-minimal"),
					// Check defaults are applied
					resource.TestCheckResourceAttr("bigip_gtm_monitor_tcp.test-gtm-tcp-monitor-minimal", "defaults_from", "/Common/tcp"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_tcp.test-gtm-tcp-monitor-minimal", "destination", "*:*"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_tcp.test-gtm-tcp-monitor-minimal", "interval", "30"),
					resource.TestCheckResourceAttr("bigip_gtm_monitor_tcp.test-gtm-tcp-monitor-minimal", "timeout", "120"),
				),
			},
		},
	})
}
