/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resLmName = "bigip_ltm_monitor"

var TestMonitorName = fmt.Sprintf("/%s/test-monitor", TestPartition)
var TestHttpsMonitorName = fmt.Sprintf("/%s/test-https-monitor", TestPartition)
var TestFtpMonitorName = fmt.Sprintf("/%s/test-ftp-monitor", TestPartition)
var TestUdpMonitorName = fmt.Sprintf("/%s/test-udp-monitor", TestPartition)
var TestPostgresqlMonitorName = fmt.Sprintf("/%s/test-postgresql-monitor", TestPartition)
var TestGatewayIcmpMonitorName = fmt.Sprintf("/%s/test-gateway", TestPartition)
var TestTcpHalfOpenMonitorName = fmt.Sprintf("/%s/test-tcp-half-open", TestPartition)

var TestMonitorResource = `
resource "bigip_ltm_monitor" "test-monitor" {
	name = "` + TestMonitorName + `"
	parent = "/Common/http"
	send = "GET /some/path\r\n"
	timeout = 999
	interval = 998
	receive = "HTTP 1.1 302 Found"
	receive_disable = "HTTP/1.1 429"
	reverse = "disabled"
	transparent = "disabled"
	manual_resume = "disabled"
	ip_dscp = 0
	time_until_up = 0
	destination = "1.2.3.4:1234"
}
`

var TestHttpsMonitorResource = `
resource "bigip_ltm_monitor" "test-https-monitor" {
	name = "` + TestHttpsMonitorName + `"
	parent = "/Common/https"
	interval          = 5
	time_until_up     = 0
	timeout           = 16
	send = "GET /some/path\r\n"
	reverse = "disabled"
	destination       = "*:8008"
	compatibility    = "enabled"
	ssl_profile      = "/Common/serverssl"
}
`

var TestFtpMonitorResource = `
resource "bigip_ltm_monitor" "test-ftp-monitor" {
	name = "` + TestFtpMonitorName + `"
	parent = "/Common/ftp"
	interval          = 5
	time_until_up     = 0
	timeout           = 16
	destination       = "*:8008"
	filename = "somefile"
	mode = "passive"
	adaptive = ""
	adaptive_limit = "0"
	transparent = ""
}
`

var TestUdpMonitorResource = `
resource "bigip_ltm_monitor" "test-udp-monitor" {
        name = "` + TestUdpMonitorName + `"
        parent = "/Common/udp"
        interval          = 5
        time_until_up     = 0
        timeout           = 16
        reverse = "disabled"
        send = "default send string"
}
`

var TestPostgresqlMonitorResource = `
resource "bigip_ltm_monitor" "test-postgresql-monitor" {
        name = "` + TestPostgresqlMonitorName + `"
        parent = "/Common/postgresql"
        interval          = 5
        time_until_up     = 0
        timeout           = 16
        database          = "postgres"
}
`

var TestGatewayIcmpMonitorResource = `
resource "bigip_ltm_monitor" "test-gateway-icmp-monitor" {
  name        = "` + TestGatewayIcmpMonitorName + `"
  parent      = "/Common/gateway_icmp"
  timeout     = "16"
  interval    = "5"
  destination = "10.10.10.10:*"
}
`

var TestTcpHalfOpenMonitorResource = `
resource "bigip_ltm_monitor" "test-tcp-half-open-monitor" {
  name        = "` + TestTcpHalfOpenMonitorName + `"
  parent      = "/Common/tcp_half_open"
  timeout     = "16"
  interval    = "5"
  destination = "10.10.10.10:1234"
}
`

func TestAccBigipLtmMonitor_GatewayIcmpCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestGatewayIcmpMonitorResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TestGatewayIcmpMonitorName),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-gateway-icmp-monitor", "parent", "/Common/gateway_icmp"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-gateway-icmp-monitor", "timeout", "16"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-gateway-icmp-monitor", "interval", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-gateway-icmp-monitor", "destination", "10.10.10.10:*"),
				),
			},
		},
	})
}

func TestAccBigipLtmMonitor_TcpHalfOpenCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestTcpHalfOpenMonitorResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TestTcpHalfOpenMonitorName),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-tcp-half-open-monitor", "parent", "/Common/tcp_half_open"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-tcp-half-open-monitor", "timeout", "16"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-tcp-half-open-monitor", "interval", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-tcp-half-open-monitor", "destination", "10.10.10.10:1234"),
				),
			},
		},
	})
}

func TestAccBigipLtmMonitor_HttpCreate(t *testing.T) {
	t.Parallel()
	var instName = "test-monitor-http"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resLmName, instName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmmonitorUpdateparam(instName, "http", ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "send", "GET /\\r\\n"),
					resource.TestCheckResourceAttr(resFullName, "timeout", "16"),
					resource.TestCheckResourceAttr(resFullName, "interval", "5"),
				),
			},
		},
	})
}
func TestAccBigipLtmMonitor_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestMonitorResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TestMonitorName),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "parent", "/Common/http"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "send", "GET /some/path\\r\\n"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "timeout", "999"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "interval", "998"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "receive", "HTTP 1.1 302 Found"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "receive_disable", "HTTP/1.1 429"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "reverse", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "transparent", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "manual_resume", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "ip_dscp", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "time_until_up", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-monitor", "destination", "1.2.3.4:1234"),
				),
			},
		},
	})
}
func TestAccBigipLtmMonitor_HttpsCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestHttpsMonitorResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TestHttpsMonitorName),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "parent", "/Common/https"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "timeout", "16"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "interval", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "time_until_up", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "send", "GET /some/path\\r\\n"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "destination", "*:8008"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "compatibility", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "reverse", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "ssl_profile", "/Common/serverssl"),
				),
			},
		},
	})
}
func TestAccBigipLtmMonitor_FtpCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestFtpMonitorResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TestFtpMonitorName),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "parent", "/Common/ftp"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "timeout", "16"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "interval", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "time_until_up", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "destination", "*:8008"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "filename", "somefile"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "mode", "passive"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "adaptive", ""),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "adaptive_limit", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-ftp-monitor", "transparent", ""),
				),
			},
		},
	})
}
func TestAccBigipLtmMonitor_UdpCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestUdpMonitorResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TestUdpMonitorName),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-udp-monitor", "parent", "/Common/udp"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-udp-monitor", "timeout", "16"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-udp-monitor", "interval", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-udp-monitor", "time_until_up", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-udp-monitor", "reverse", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-udp-monitor", "send", "default send string"),
				),
			},
		},
	})
}
func TestAccBigipLtmMonitor_PostgresqlCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPostgresqlMonitorResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TestPostgresqlMonitorName),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-postgresql-monitor", "parent", "/Common/postgresql"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-postgresql-monitor", "timeout", "16"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-postgresql-monitor", "interval", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-postgresql-monitor", "time_until_up", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-postgresql-monitor", "database", "postgres"),
				),
			},
		},
	})
}

func TestAccBigipLtmMonitorTestCases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: loadFixtureString("../examples/bigip_ltm_monitor.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists("/Common/test_monitor_tc1"),
					testCheckMonitorExists("/Common/test_monitor_tc2"),
					testCheckMonitorExists("/Common/test_monitor_tc3"),
					testCheckMonitorExists("/Common/test_monitor_tc4"),
					testCheckMonitorExists("/Common/test_monitor_tc5"),
				),
			},
		},
	})
}

func TestAccBigipLtmMonitor_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestMonitorResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TestMonitorName),
				),
				ResourceName:      TestMonitorName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckMonitorExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		monitors, err := client.Monitors()
		if err != nil {
			return err
		}

		for _, m := range monitors {
			if m.FullPath == name {
				return nil
			}
		}
		return fmt.Errorf("Monitor %s was not created ", name)
	}
}

func testMonitorsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	monitors, err := client.Monitors()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_monitor" {
			continue
		}

		name := rs.Primary.ID
		for _, m := range monitors {
			if m.FullPath == name {
				return fmt.Errorf("Monitor %s not destroyed ", name)
			}
		}
	}
	return nil
}

func testaccbigipltmmonitorUpdateparam(instName, parentM, updateParam string) string {
	resPrefix := fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			  name = "/Common/%[2]s"`, resLmName, instName)
	switch parentM {
	case "http":
		resPrefix = fmt.Sprintf(`%s
			  parent = "/Common/http"`, resPrefix)
	case "https":
		resPrefix = fmt.Sprintf(`%s
			  parent = "/Common/https"`, resPrefix)
	case "udp":
		resPrefix = fmt.Sprintf(`%s
			  parent = "/Common/udp"`, resPrefix)
	case "ftp":
		resPrefix = fmt.Sprintf(`%s
			  parent = "/Common/ftp"`, resPrefix)
	case "postgresql":
		resPrefix = fmt.Sprintf(`%s
			  parent = "/Common/postgresql"`, resPrefix)
	default:
		resPrefix = fmt.Sprintf(`%s
			  parent = "/Common/http"`, resPrefix)
	}
	switch updateParam {
	case "timeout":
		resPrefix = fmt.Sprintf(`%s
			  timeout = "always"`, resPrefix)
	case "interval":
		resPrefix = fmt.Sprintf(`%s
			  interval = ["no-tlsv1.3"]`, resPrefix)
	case "send":
		resPrefix = fmt.Sprintf(`%s
			  send = 8`, resPrefix)
	case "receive":
		resPrefix = fmt.Sprintf(`%s
			  receive = 262100`, resPrefix)
	default:
	}
	return fmt.Sprintf(`%s
		}`, resPrefix)
}
