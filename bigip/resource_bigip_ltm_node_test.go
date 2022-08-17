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

var TestNodeName = fmt.Sprintf("/%s/test-node", TEST_PARTITION)
var TestV6NodeName = fmt.Sprintf("/%s/test-v6-node", TEST_PARTITION)
var TestFqdnNodeName = fmt.Sprintf("/%s/test-fqdn-node", TEST_PARTITION)

var resNodeName = "bigip_ltm_node"

type UpdateParam struct {
	key   string
	value string
}

var TestNodeResource = `
resource "bigip_ltm_node" "test-node" {
	name = "` + TestNodeName + `"
	address = "192.168.30.1"
	connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "/Common/icmp"
	rate_limit = "disabled"
	state = "user-up"
	ratio = "91"
}
`

var TestV6NodeResource = `
resource "bigip_ltm_node" "test-node" {
	name = "` + TestV6NodeName + `"
	address = "fe80::10"
	connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "default"
	rate_limit = "disabled"
	state = "user-up"
}
`

var TestFqdnNodeResource = `
resource "bigip_ltm_node" "test-fqdn-node" {
	name = "` + TestFqdnNodeName + `"
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

func TestAccBigipLtmNode_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestNodeResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TestNodeName),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "name", TestNodeName),
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
}

func TestAccBigipLtmNode_V6create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestV6NodeResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TestV6NodeName),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "name", TestV6NodeName),
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
}

func TestAccBigipLtmNode_FqdnCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestFqdnNodeResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TestFqdnNodeName),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-fqdn-node", "name", TestFqdnNodeName),
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
func TestAccBigipLtmNodeUpdateMonitor(t *testing.T) {
	t.Parallel()
	var instName = "test-node-monitor"
	var TestNodeName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resNodeName, instName)
	var moni UpdateParam
	var moni2 UpdateParam
	moni.key = "monitor"
	moni.value = "default"
	moni2.key = "monitor"
	moni2.value = "/Common/none"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmNodeUpdateParam(instName, moni),
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TestNodeName),
					resource.TestCheckResourceAttr(resFullName, "name", TestNodeName),
					resource.TestCheckResourceAttr(resFullName, "monitor", "default"),
				),
			},
			{
				Config: testaccbigipltmNodeUpdateParam(instName, moni2),
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TestNodeName),
					resource.TestCheckResourceAttr(resFullName, "name", TestNodeName),
					resource.TestCheckResourceAttr(resFullName, "monitor", "/Common/none"),
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
				Config: TestNodeResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TestNodeName),
				),
				ResourceName:      TestNodeName,
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
				Config: TestV6NodeResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TestV6NodeName),
				),
				ResourceName:      TestV6NodeName,
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
				Config: TestFqdnNodeResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TestFqdnNodeName),
				),
				ResourceName:      TestFqdnNodeName,
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
			return fmt.Errorf("Node %s was not created ", name)
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
			return fmt.Errorf("Node %s not destroyed ", name)
		}
	}
	return nil
}

func testaccbigipltmNodeUpdateParam(instName string, updateParam UpdateParam) string {
	resPrefix := fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			  name = "/Common/%[2]s"
              address = "192.168.100.100"
		`, resNodeName, instName)
	switch updateParam.key {
	case "monitor":
		resPrefix = fmt.Sprintf(`%s
			  monitor = "%s"`, resPrefix, updateParam.value)
	default:
	}
	return fmt.Sprintf(`%s
		}`, resPrefix)
}
