/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"testing"
)

var TEST_POLICY_NAME = "test-policy"

var TEST_POLICY_RESOURCE = `
resource "bigip_ltm_node" "test-node" {
	name = "` + TEST_NODE_NAME + `"
	address = "10.10.10.10"
	connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "default"
	rate_limit = "disabled"
}
resource "bigip_ltm_pool" "test-pool" {
	name = "` + TEST_POOL_NAME + `"
	monitors = ["/Common/http"]
	allow_nat = "yes"
	allow_snat = "yes"
	load_balancing_mode = "round-robin"
	depends_on = ["bigip_ltm_node.test-node"]
}
resource "bigip_ltm_pool_attachment" "test-pool_test-node" {
	pool = "` + TEST_POOL_NAME + `"
	node = "` + TEST_POOLNODE_NAMEPORT + `"
	depends_on = ["bigip_ltm_node.test-node", "bigip_ltm_pool.test-pool"]
}
resource "bigip_ltm_policy" "test-policy" {
	depends_on = ["bigip_ltm_pool.test-pool"]
	name = "` + TEST_POLICY_NAME + `"
	strategy = "/Common/first-match"
	requires = ["http"]
	published_copy = "Drafts/` + TEST_POLICY_NAME + `"
	controls = ["forwarding"]
	rule  {
	      name = "rule6"
		      action {
			      tm_name = "20"
			      forward = true
			      pool = "/Common/test-pool"
		      }
	}
}
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
`

func TestAccBigipLtmPolicy_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_POLICY_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TEST_POLICY_NAME, true),
					testCheckPolicyExists("http_to_https_redirect", true),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_POLICY_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TEST_POLICY_NAME, true),
				),
				ResourceName:      TEST_POLICY_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckPolicyExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		policy, err := client.GetPolicy(name)
		if err != nil {
			return fmt.Errorf("Error while fetching policy: %v", err)

		}
		if exists && policy == nil {
			return fmt.Errorf("Policy %s was not created.", name)
		}
		if !exists && policy != nil {
			return fmt.Errorf("Policy %s still exists.", name)
		}
		log.Printf("[DEBUG] Policy \"%s\" Created ", name)
		return nil
	}
}

func testCheckPolicysDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_policy" {
			continue
		}

		name := rs.Primary.ID
		policy, err := client.GetPolicy(name)

		if err != nil {
			return nil
		}
		if policy != nil {
			return fmt.Errorf("Policy %s not destroyed.", name)
		}
	}
	return nil
}
