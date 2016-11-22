package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
	"testing"
)

var TEST_VS_NAME = fmt.Sprintf("/%s/test-vs", TEST_PARTITION)

var TEST_VS_RESOURCE = TEST_IRULE_RESOURCE + TEST_POLICY_RESOURCE + `
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
	policies = ["${bigip_ltm_policy.test-policy.name}"]
}
`

func TestBigipLtmVS_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
			testCheckIRulesDestroyed,
			testCheckPolicyDestroyed,
		),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_VS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(TEST_VS_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", TEST_VS_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "10.255.255.254"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "9999"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source_address_translation", "automap"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("irules.%d", schema.HashString(TEST_IRULE_NAME)),
						TEST_IRULE_NAME),
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
						fmt.Sprintf("policies.%d", schema.HashString(TEST_POLICY_NAME)),
						TEST_POLICY_NAME),
				),
			},
		},
	})
}

func TestBigipLtmVS_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVSsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_VS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(TEST_VS_NAME, true),
				),
				ResourceName:      TEST_VS_NAME,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//TODO: test adding rules, profiles, policies, etc

func testCheckVSExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		vs, err := client.GetVirtualServer(name)
		if err != nil {
			return err
		}
		if exists && vs == nil {
			return fmt.Errorf("Virtual server ", name, " does not exist.")
		}
		if !exists && vs != nil {
			return fmt.Errorf("Virtual server ", name, " exists.")
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
			return fmt.Errorf("Virtual server ", name, " not destroyed.")
		}
	}
	return nil
}
