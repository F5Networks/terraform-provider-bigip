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

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_MONITOR_NAME = fmt.Sprintf("/%s/test-monitor", TEST_PARTITION)
var TEST_HTTPS_MONITOR_NAME = fmt.Sprintf("/%s/test-https-monitor", TEST_PARTITION)
var TEST_FTP_MONITOR_NAME = fmt.Sprintf("/%s/test-ftp-monitor", TEST_PARTITION)
var TEST_UDP_MONITOR_NAME = fmt.Sprintf("/%s/test-udp-monitor", TEST_PARTITION)
var TEST_POSTGRESQL_MONITOR_NAME = fmt.Sprintf("/%s/test-postgresql-monitor", TEST_PARTITION)

var TEST_MONITOR_RESOURCE = `
resource "bigip_ltm_monitor" "test-monitor" {
	name = "` + TEST_MONITOR_NAME + `"
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

var TEST_HTTPS_MONITOR_RESOURCE = `
resource "bigip_ltm_monitor" "test-https-monitor" {
	name = "` + TEST_HTTPS_MONITOR_NAME + `"
	parent = "/Common/https"
	interval          = 5
	time_until_up     = 0
	timeout           = 16
	send = "GET /some/path\r\n"
	reverse = "disabled"
	destination       = "*:8008"
	compatibility    = "enabled"
}
`

var TEST_FTP_MONITOR_RESOURCE = `
resource "bigip_ltm_monitor" "test-ftp-monitor" {
	name = "` + TEST_FTP_MONITOR_NAME + `"
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

var TEST_UDP_MONITOR_RESOURCE = `
resource "bigip_ltm_monitor" "test-udp-monitor" {
        name = "` + TEST_UDP_MONITOR_NAME + `"
        parent = "/Common/udp"
        interval          = 5
        time_until_up     = 0
        timeout           = 16
        reverse = "disabled"
        send = "default send string"
}
`

var TEST_POSTGRESQL_MONITOR_RESOURCE = `
resource "bigip_ltm_monitor" "test-postgresql-monitor" {
        name = "` + TEST_POSTGRESQL_MONITOR_NAME + `"
        parent = "/Common/postgresql"
        interval          = 5
        time_until_up     = 0
        timeout           = 16
        database          = "postgres"
}
`

func TestAccBigipLtmMonitor_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_MONITOR_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TEST_MONITOR_NAME),
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

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_HTTPS_MONITOR_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TEST_HTTPS_MONITOR_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "parent", "/Common/https"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "timeout", "16"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "interval", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "time_until_up", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "send", "GET /some/path\\r\\n"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "destination", "*:8008"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "compatibility", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_monitor.test-https-monitor", "reverse", "disabled"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FTP_MONITOR_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TEST_FTP_MONITOR_NAME),
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_UDP_MONITOR_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TEST_UDP_MONITOR_NAME),
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_POSTGRESQL_MONITOR_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TEST_POSTGRESQL_MONITOR_NAME),
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

func TestAccBigipLtmMonitor_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testMonitorsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_MONITOR_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckMonitorExists(TEST_MONITOR_NAME),
				),
				ResourceName:      TEST_MONITOR_NAME,
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
		return fmt.Errorf("Monitor %s was not created.", name)
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
				return fmt.Errorf("Monitor %s not destroyed.", name)
			}
		}
	}
	return nil
}
