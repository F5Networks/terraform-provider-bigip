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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_VS_NAME = fmt.Sprintf("/%s/test-vs", TEST_PARTITION)

var TEST_VS_RESOURCE = TEST_IRULE_RESOURCE + `


resource "bigip_ltm_policy" "http_to_https_redirect" {
  name = "http_to_https_redirect"
  strategy = "/Common/first-match"
  requires = ["http"]
  published_copy = "Drafts/http_to_https_redirect"
  controls = ["forwarding"]
  rule  {
    name = "http_to_https_redirect_rule"
    action {
      tm_name = "http_to_https_redirect"
      redirect = true
      location = "tcl:https://[HTTP::host][HTTP::uri]"
      http_reply = true
    }
  }
}
resource "bigip_ltm_virtual_server" "test-vs" {
	name = "` + TEST_VS_NAME + `"
	destination = "10.255.255.254"
	description = "VirtualServer-test"
	port = 9999
	mask = "255.255.255.255"
	source_address_translation = "automap"
	ip_protocol = "tcp"
	irules = ["${bigip_ltm_irule.test-rule.name}"]
	profiles = ["/Common/http"]
	client_profiles = ["/Common/tcp"]
	server_profiles = ["/Common/tcp-lan-optimized"]
	persistence_profiles = ["/Common/source_addr","/Common/hash"]
	default_persistence_profile = "/Common/hash"
	fallback_persistence_profile = "/Common/dest_addr"
        policies = ["${bigip_ltm_policy.http_to_https_redirect.name}"]

}
`
var TEST_VS6_NAME = fmt.Sprintf("/%s/test-vs6", TEST_PARTITION)

var TEST_VS6_RESOURCE = TEST_IRULE_RESOURCE + `
resource "bigip_ltm_virtual_server" "test-vs" {
	name = "` + TEST_VS6_NAME + `"
        destination = "fe80::11"
	description = "VirtualServer-test"
	port = 9999
	source_address_translation = "automap"
	ip_protocol = "tcp"
	irules = ["${bigip_ltm_irule.test-rule.name}"]
	profiles = ["/Common/http"]
	client_profiles = ["/Common/tcp"]
	server_profiles = ["/Common/tcp-lan-optimized"]
	persistence_profiles = ["/Common/source_addr", "/Common/hash"]
	default_persistence_profile = "/Common/hash"
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
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "enabled"),
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
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", "VirtualServer-test"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source_address_translation", "automap"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "irules.0", TEST_IRULE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "default_persistence_profile", "/Common/hash"),
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
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("persistence_profiles.%d", schema.HashString("/Common/hash")),
						"/Common/hash"),
					//resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
					//		fmt.Sprintf("policies.%d", schema.HashString("http_to_https_redirect")),
					//	"http_to_https_redirect"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "fallback_persistence_profile", "/Common/dest_addr"),
				),
			},
		},
	})
}

func TestAccBigipLtmVS_create_Defaultstate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testVSCreateDefaultstate("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample", true),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")), "/Common/fastL4"),
				),
			},
		},
	})
}

func TestAccBigipLtmVS_Modify_stateDisabledtoEnabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testVSCreatestatedisabled("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample", true),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
						"/Common/fastL4"),
				),
			},
			{
				Config: testVSCreatestateenabled("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample", true),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
						"/Common/fastL4"),
				),
			},
		},
	})
}

func TestAccBigipLtmVS_Modify_stateEnabledtoDisabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testVSCreatestateenabled("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample", true),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
						"/Common/fastL4"),
				),
			},
			{
				Config: testVSCreatestatedisabled("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample", true),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
						"/Common/fastL4"),
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
func testVSCreatestatedisabled(vs_name string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test-vs" {
        name = "/Common/%s"
        destination = "192.168.50.1"
        port = 800
        mask = "255.255.255.255"
        state = "disabled"
}`, vs_name)
}

func testVSCreatestateenabled(vs_name string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test-vs" {
        name = "/Common/%s"
        destination = "192.168.50.1"
        port = 800
        mask = "255.255.255.255"
        state = "enabled"
}`, vs_name)
}

func testVSCreateDefaultstate(vs_name string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test-vs" {
        name = "/Common/%s"
        destination = "192.168.50.1"
        port = 800
        mask = "255.255.255.255"
}`, vs_name)
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
