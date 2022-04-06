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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_NODE_NAME = fmt.Sprintf("/%s/test-node", TEST_PARTITION)
var TEST_V6_NODE_NAME = fmt.Sprintf("/%s/test-v6-node", TEST_PARTITION)
var TEST_FQDN_NODE_NAME = fmt.Sprintf("/%s/test-fqdn-node", TEST_PARTITION)

var TEST_NODE_RESOURCE = `
resource "bigip_ltm_node" "test-node" {
	name = "` + TEST_NODE_NAME + `"
	address = "192.168.30.1"
	connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "/Common/icmp"
	rate_limit = "disabled"
	state = "user-up"
	ratio = "91"
}
`

var TEST_V6_NODE_RESOURCE = `
resource "bigip_ltm_node" "test-node" {
	name = "` + TEST_V6_NODE_NAME + `"
	address = "fe80::10"
	connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "default"
	rate_limit = "disabled"
	state = "user-up"
}
`

var TEST_FQDN_NODE_RESOURCE = `
resource "bigip_ltm_node" "test-fqdn-node" {
	name = "` + TEST_FQDN_NODE_NAME + `"
	address = "f5.com"
	connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "default"
	rate_limit = "disabled"
	fqdn { interval = "3000"}
	state = "user-up"
	ratio = "19"
}
`

func TestAccBigipLtmNode_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_NODE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TEST_NODE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "name", TEST_NODE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "address", "192.168.30.1"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "connection_limit", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "dynamic_ratio", "1"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "monitor", "/Common/icmp"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "rate_limit", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "state", "user-up"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "session", "user-enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "ratio", "91"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_V6_NODE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TEST_V6_NODE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "name", TEST_V6_NODE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "address", "fe80::10"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "connection_limit", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "dynamic_ratio", "1"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "monitor", "default"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "rate_limit", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "state", "user-up"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FQDN_NODE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TEST_FQDN_NODE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "name", TEST_FQDN_NODE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "address", "f5.com"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "connection_limit", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "dynamic_ratio", "1"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "monitor", "default"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "rate_limit", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "state", "user-up"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "fqdn.0.interval", "3000"),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "ratio", "19"),
				),
			},
		},
	})
}

func TestAccBigipLtmNode_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_NODE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TEST_NODE_NAME),
				),
				ResourceName:      TEST_NODE_NAME,
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
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_V6_NODE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TEST_V6_NODE_NAME),
				),
				ResourceName:      TEST_V6_NODE_NAME,
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
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FQDN_NODE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TEST_FQDN_NODE_NAME),
				),
				ResourceName:      TEST_FQDN_NODE_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckNodeExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		node, err := client.GetNode(name)
		if err != nil {
			return err
		}
		if node == nil {
			return fmt.Errorf("Node %s was not created.", name)
		}

		return nil
	}
}

func testCheckNodesDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_node" {
			continue
		}

		name := rs.Primary.ID
		node, err := client.GetNode(name)
		if err != nil {
			return err
		}
		if node != nil {
			return fmt.Errorf("Node %s not destroyed.", name)
		}
	}
	return nil
}
