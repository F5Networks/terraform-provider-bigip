package bigip

import (
	"fmt"
	"strings"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TEST_GTM_SERVER_NAME = "test_gtm_server"
var TEST_GTM_SERVER_DATACENTER = "/Common/test_datacenter"

func TestAccBigipGtmServer_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmServerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmServerConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(TEST_GTM_SERVER_NAME, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "name", TEST_GTM_SERVER_NAME),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "datacenter", TEST_GTM_SERVER_DATACENTER),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "product", "bigip"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "virtual_server_discovery", "enabled"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "link_discovery", "disabled"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "addresses.#", "1"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "addresses.0.name", "10.10.10.10"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "addresses.0.translation", "none"),
				),
			},
		},
	})
}

func TestAccBigipGtmServer_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmServerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmServerConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(TEST_GTM_SERVER_NAME, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "addresses.#", "1"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "monitor", "/Common/bigip"),
				),
			},
			{
				Config: testAccBigipGtmServerConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(TEST_GTM_SERVER_NAME, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "addresses.#", "2"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "addresses.0.name", "10.10.10.10"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "addresses.1.name", "10.10.10.11"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "monitor", "/Common/http"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "virtual_server_discovery", "disabled"),
				),
			},
		},
	})
}

func TestAccBigipGtmServer_withDeviceName(t *testing.T) {
	serverName := "test_gtm_server_device"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmServerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmServerConfigWithDevice(serverName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(serverName, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-device", "addresses.0.name", "10.10.10.20"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-device", "addresses.0.device_name", "/Common/device1"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-device", "addresses.0.translation", "192.168.1.20"),
				),
			},
		},
	})
}

func TestAccBigipGtmServer_multipleAddresses(t *testing.T) {
	serverName := "test_gtm_server_multi"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmServerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmServerConfigMultiAddress(serverName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(serverName, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-multi", "addresses.#", "3"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-multi", "addresses.0.name", "10.10.10.30"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-multi", "addresses.1.name", "10.10.10.31"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-multi", "addresses.2.name", "10.10.10.32"),
				),
			},
		},
	})
}

func TestAccBigipGtmServer_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmServerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmServerConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(TEST_GTM_SERVER_NAME, true),
				),
			},
			{
				ResourceName:      "bigip_gtm_server.test-server",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckGtmServerExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		server, err := client.GetGtmserver(name)
		if err != nil {
			return err
		}
		if exists && server == nil {
			return fmt.Errorf("GTM Server %s does not exist", name)
		}
		if !exists && server != nil {
			return fmt.Errorf("GTM Server %s still exists", name)
		}
		return nil
	}
}

func testCheckGtmServerDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_gtm_server" {
			continue
		}

		name := rs.Primary.ID
		server, err := client.GetGtmserver(name)
		if err != nil {
			// Ignore "not found" errors - this is the expected state
			errStr := err.Error()
			if strings.Contains(errStr, "was not found") || strings.Contains(errStr, "01020036") {
				continue
			}
			return err
		}
		if server != nil {
			return fmt.Errorf("GTM Server %s still exists", name)
		}
	}
	return nil
}

func testAccBigipGtmServerConfig() string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "bigip"
  monitor    = "/Common/bigip"

  virtual_server_discovery = "enabled"
  link_discovery           = "disabled"

  addresses {
    name        = "10.10.10.10"
    translation = "none"
  }
}
`, TEST_GTM_SERVER_NAME)
}

func testAccBigipGtmServerConfigUpdated() string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "bigip"
  monitor    = "/Common/http"

  virtual_server_discovery = "disabled"
  link_discovery           = "disabled"

  addresses {
    name        = "10.10.10.10"
    translation = "none"
  }

  addresses {
    name        = "10.10.10.11"
    translation = "none"
  }
}
`, TEST_GTM_SERVER_NAME)
}

func testAccBigipGtmServerConfigWithDevice(serverName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server-device" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "bigip"

  virtual_server_discovery = "enabled"

  addresses {
    name        = "10.10.10.20"
    device_name = "/Common/device1"
    translation = "192.168.1.20"
  }
}
`, serverName)
}

func testAccBigipGtmServerConfigMultiAddress(serverName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server-multi" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "bigip"

  virtual_server_discovery = "enabled"

  addresses {
    name = "10.10.10.30"
  }

  addresses {
    name = "10.10.10.31"
  }

  addresses {
    name = "10.10.10.32"
  }
}
`, serverName)
}

func TestAccBigipGtmServer_withVirtualServers(t *testing.T) {
	serverName := "test_gtm_generic_vs"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmServerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmServerConfigWithVirtualServers(serverName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(serverName, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "name", serverName),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "product", "generic-host"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_server_discovery", "disabled"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.#", "2"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.0.name", "web_http"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.0.destination", "10.20.30.40:80"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.0.enabled", "true"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.1.name", "web_https"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.1.destination", "10.20.30.40:443"),
				),
			},
		},
	})
}

func TestAccBigipGtmServer_virtualServersUpdate(t *testing.T) {
	serverName := "test_gtm_vs_update"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmServerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmServerConfigWithVirtualServers(serverName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(serverName, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.#", "2"),
				),
			},
			{
				Config: testAccBigipGtmServerConfigWithVirtualServersUpdated(serverName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(serverName, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.#", "3"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.2.name", "api_service"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs", "virtual_servers.2.destination", "10.20.30.40:8080"),
				),
			},
		},
	})
}

func TestAccBigipGtmServer_virtualServersWithLimits(t *testing.T) {
	serverName := "test_gtm_vs_limits"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmServerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmServerConfigWithVSLimits(serverName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(serverName, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-limits", "virtual_servers.#", "1"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-limits", "virtual_servers.0.name", "limited_vs"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-limits", "virtual_servers.0.limit_max_connections", "10000"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-limits", "virtual_servers.0.limit_max_connections_status", "enabled"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-limits", "virtual_servers.0.limit_max_bps", "1000000"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-limits", "virtual_servers.0.limit_max_bps_status", "enabled"),
				),
			},
		},
	})
}

func TestAccBigipGtmServer_virtualServersWithTranslation(t *testing.T) {
	serverName := "test_gtm_vs_nat"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckGtmServerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipGtmServerConfigWithVSTranslation(serverName),
				Check: resource.ComposeTestCheckFunc(
					testCheckGtmServerExists(serverName, true),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-nat", "virtual_servers.#", "1"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-nat", "virtual_servers.0.name", "nat_vs"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-nat", "virtual_servers.0.translation_address", "203.0.113.100"),
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server-vs-nat", "virtual_servers.0.translation_port", "8080"),
				),
			},
		},
	})
}

func testAccBigipGtmServerConfigWithVirtualServers(serverName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server-vs" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "generic-host"

  virtual_server_discovery = "disabled"

  addresses {
    name = "10.20.30.40"
  }

  virtual_servers {
    name        = "web_http"
    destination = "10.20.30.40:80"
    enabled     = true
  }

  virtual_servers {
    name        = "web_https"
    destination = "10.20.30.40:443"
    enabled     = true
  }
}
`, serverName)
}

func testAccBigipGtmServerConfigWithVirtualServersUpdated(serverName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server-vs" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "generic-host"

  virtual_server_discovery = "disabled"

  addresses {
    name = "10.20.30.40"
  }

  virtual_servers {
    name        = "web_http"
    destination = "10.20.30.40:80"
    enabled     = true
  }

  virtual_servers {
    name        = "web_https"
    destination = "10.20.30.40:443"
    enabled     = true
  }

  virtual_servers {
    name        = "api_service"
    destination = "10.20.30.40:8080"
    enabled     = true
  }
}
`, serverName)
}

func testAccBigipGtmServerConfigWithVSLimits(serverName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server-vs-limits" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "generic-host"

  virtual_server_discovery = "disabled"

  addresses {
    name = "10.20.30.50"
  }

  virtual_servers {
    name                         = "limited_vs"
    destination                  = "10.20.30.50:8080"
    enabled                      = true
    limit_max_connections        = 10000
    limit_max_connections_status = "enabled"
    limit_max_bps                = 1000000
    limit_max_bps_status         = "enabled"
    limit_max_pps                = 5000
    limit_max_pps_status         = "enabled"
  }
}
`, serverName)
}

func testAccBigipGtmServerConfigWithVSTranslation(serverName string) string {
	return fmt.Sprintf(`
resource "bigip_gtm_datacenter" "test-datacenter" {
  name      = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server-vs-nat" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "generic-host"

  virtual_server_discovery = "disabled"

  addresses {
    name        = "192.168.1.100"
    translation = "203.0.113.100"
  }

  virtual_servers {
    name                = "nat_vs"
    destination         = "192.168.1.100:80"
    enabled             = true
    translation_address = "203.0.113.100"
    translation_port    = 8080
  }
}
`, serverName)
}
