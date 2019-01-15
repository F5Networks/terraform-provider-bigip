package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_MONITOR_NAME = fmt.Sprintf("/%s/test-monitor", TEST_PARTITION)
var TEST_HTTPS_MONITOR_NAME = fmt.Sprintf("/%s/test-https-monitor", TEST_PARTITION)

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
