package bigip

import (
	"fmt"
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
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "virtual_server_discovery", "true"),
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
					resource.TestCheckResourceAttr("bigip_gtm_server.test-server", "virtual_server_discovery", "false"),
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
  name = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "bigip"
  monitor    = "/Common/bigip"

  virtual_server_discovery = true
  link_discovery          = "disabled"

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
  name = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "bigip"
  monitor    = "/Common/http"

  virtual_server_discovery = false
  link_discovery          = "disabled"

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
  name = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server-device" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "bigip"

  virtual_server_discovery = true

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
  name = "test_datacenter"
  partition = "Common"
}

resource "bigip_gtm_server" "test-server-multi" {
  name       = "%s"
  partition  = "Common"
  datacenter = bigip_gtm_datacenter.test-datacenter.id
  product    = "bigip"

  virtual_server_discovery = true

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
