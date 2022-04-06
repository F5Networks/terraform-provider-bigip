/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_TUNNEL_NAME = "test-tunnel"

var TEST_TUNNEL_RESOURCE = `
resource "bigip_net_tunnel" "test_tunnel" {
    name = "` + TEST_TUNNEL_NAME + `"
    auto_last_hop     = "default"
    idle_timeout      = 300
//    if_index          = 464
    key               = 0
    local_address     = "192.16.81.240"
    mode              = "bidirectional"
    mtu               = 0
    profile           = "/Common/dslite"
    remote_address    = "any6"
    secondary_address = "any6"
    tos               = "preserve"
    transparent       = "disabled"
    use_pmtu          = "enabled"        
}
`
var TEST_TUNNEL_RESOURCE_UPDATE = `
resource "bigip_net_tunnel" "test_tunnel" {
    name = "` + TEST_TUNNEL_NAME + `"
    auto_last_hop     = "default"
    idle_timeout      = 300
//    if_index          = 464
    key               = 0
    local_address     = "192.16.81.240"
    mode              = "bidirectional"
    mtu               = 0
    profile           = "/Common/dslite"
    remote_address    = "any6"
    secondary_address = "any6"
    tos               = "preserve"
    transparent       = "enabled"
    use_pmtu          = "disabled"
}
`

func TestAccBigipNetTunnelCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testCheckBigipNetTunnelDestroyed),
		Steps: []resource.TestStep{
			{
				Config: TEST_TUNNEL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testBigipNetTunnelExists(TEST_TUNNEL_NAME, true),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "name", TEST_TUNNEL_NAME),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "auto_last_hop", "default"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "idle_timeout", "300"),
					//  resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "if_index", "464"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "key", "0"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "local_address", "192.16.81.240"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "mode", "bidirectional"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "mtu", "0"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "profile", "/Common/dslite"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "remote_address", "any6"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "secondary_address", "any6"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "tos", "preserve"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "transparent", "disabled"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "use_pmtu", "enabled"),
				),
			},
		},
	})

}
func TestAccBigipNetTunnelUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testCheckBigipNetTunnelDestroyed),
		Steps: []resource.TestStep{
			{
				Config: TEST_TUNNEL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "transparent", "disabled"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "use_pmtu", "enabled"),
				),
			},
			{
				Config: TEST_TUNNEL_RESOURCE_UPDATE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "transparent", "enabled"),
					resource.TestCheckResourceAttr("bigip_net_tunnel.test_tunnel", "use_pmtu", "disabled"),
				),
			},
		},
	})
}
func TestAccBigipNetTunnelImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckBigipNetTunnelDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_TUNNEL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testBigipNetTunnelExists(TEST_TUNNEL_NAME, true),
				),
				ResourceName:      TEST_TUNNEL_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testBigipNetTunnelExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		pp, err := client.GetTunnel(name)
		if err != nil {
			return err
		}
		if exists && pp == nil {
			return fmt.Errorf("Tunnel %s does not exist.", name)
		}
		if !exists && pp != nil {
			return fmt.Errorf("Tunnel %s exists.", name)
		}
		return nil
	}
}

func testCheckBigipNetTunnelDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_net_tunnel" {
			continue
		}

		name := rs.Primary.ID
		pp, err := client.GetTunnel(name)
		if err != nil {
			return err
		}

		if pp != nil {
			return fmt.Errorf("Tunnel %s not destroyed.", name)
		}
	}
	return nil
}
