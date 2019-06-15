package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_VS_NAME = fmt.Sprintf("/%s/test-vs", TEST_PARTITION)

var TEST_VS_RESOURCE = TEST_IRULE_RESOURCE + `


resource "bigip_ltm_virtual_server" "test-vs" {
	name = "` + TEST_VS_NAME + `"
	destination = "10.255.255.254"
	port = 9999
	mask = "255.255.255.255"
	source_address_translation = "automap"
	ip_protocol = "tcp"
	irules = ["${bigip_ltm_irule.test-rule.name}"]
	profiles = ["/Common/http"]
	client_profiles = ["/Common/tcp"]
	server_profiles = ["/Common/tcp-lan-optimized"]
	persistence_profiles = ["/Common/source_addr"]
	fallback_persistence_profile = "/Common/dest_addr"

}
`

var TEST_VS6_NAME = fmt.Sprintf("/%s/test-vs6", TEST_PARTITION)

var TEST_VS6_RESOURCE = TEST_IRULE_RESOURCE + `


resource "bigip_ltm_virtual_server" "test-vs" {
	name = "` + TEST_VS6_NAME + `"
  destination = "fe80::11"
	port = 9999
	source_address_translation = "automap"
	ip_protocol = "tcp"
	irules = ["${bigip_ltm_irule.test-rule.name}"]
	profiles = ["/Common/http"]
	client_profiles = ["/Common/tcp"]
	server_profiles = ["/Common/tcp-lan-optimized"]
	persistence_profiles = ["/Common/source_addr"]
	fallback_persistence_profile = "/Common/dest_addr"
}
`

func TestAccBigipLtmVS_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
			testCheckIRulesDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: TEST_VS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(TEST_VS_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", TEST_VS_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "10.255.255.254"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "9999"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source_address_translation", "automap"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "irules.0", TEST_IRULE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/http")),
						"/Common/http"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("client_profiles.%d", schema.HashString("/Common/tcp")),
						"/Common/tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("server_profiles.%d", schema.HashString("/Common/tcp-lan-optimized")),
						"/Common/tcp-lan-optimized"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("persistence_profiles.%d", schema.HashString("/Common/source_addr")),
						"/Common/source_addr"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "fallback_persistence_profile", "/Common/dest_addr"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
			testCheckIRulesDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: TEST_VS6_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(TEST_VS6_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", TEST_VS6_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "fe80::11"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "9999"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source", "::/0"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source_address_translation", "automap"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "irules.0", TEST_IRULE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/http")),
						"/Common/http"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("client_profiles.%d", schema.HashString("/Common/tcp")),
						"/Common/tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("server_profiles.%d", schema.HashString("/Common/tcp-lan-optimized")),
						"/Common/tcp-lan-optimized"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("persistence_profiles.%d", schema.HashString("/Common/source_addr")),
						"/Common/source_addr"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "fallback_persistence_profile", "/Common/dest_addr"),
				),
			},
		},
	})
}

func TestAccBigipLtmVS_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVSsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_VS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(TEST_VS_NAME, true),
				),
				ResourceName:      TEST_VS_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVSsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_VS6_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(TEST_VS6_NAME, true),
				),
				ResourceName:      TEST_VS6_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckVSExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		vs, err := client.GetVirtualServer(name)
		if err != nil {
			return err
		}
		if exists && vs == nil {
			return fmt.Errorf("Virtual server %s does not exist.", name)
		}
		if !exists && vs != nil {
			return fmt.Errorf("Virtual server %s exists.", name)
		}
		return nil
	}
}

func testCheckVSsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_virtual_address" {
			continue
		}

		name := rs.Primary.ID
		vs, err := client.GetVirtualServer(name)
		if err != nil {
			return err
		}
		if vs != nil {
			return fmt.Errorf("Virtual server %s not destroyed.", name)
		}
	}
	return nil
}
