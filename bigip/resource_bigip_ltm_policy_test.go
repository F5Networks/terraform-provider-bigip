/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TestPolicyName = "/Common/test-policy"

var TestPolicyResource = `
resource "bigip_ltm_pool" "test-pool" {
	name = "` + TEST_POOL_NAME + `"
	monitors = ["/Common/http"]
	allow_nat = "yes"
	allow_snat = "yes"
	load_balancing_mode = "round-robin"
}
resource "bigip_ltm_policy" "test-policy" {
	depends_on = ["bigip_ltm_pool.test-pool"]
	name = "` + TestPolicyName + `"
	strategy = "first-match"
	requires = ["http"]
#	published_copy = "Drafts/` + TestPolicyName + `"
	controls = ["forwarding"]
	rule  {
	      name = "rule6"
		      action {
//			      tm_name = "20"
			      forward = true
			      pool = "/Common/test-pool"
		      }
	}
}
resource "bigip_ltm_policy" "test-policy-again" {
  name = "/Common/test-policy-again"
  strategy = "first-match"
  requires = ["http"]
 # published_copy = "Drafts/http_to_https_redirect"
  controls = ["forwarding"]
  rule  {
    name = "testrule"
    action {
  //    tm_name = "http_to_https_redirect"
      redirect = true
      location = "tcl:https://[HTTP::host][HTTP::uri]"
      http_reply = true
    }
  }
}
`
var TestPolicyResource2 = `
resource "bigip_ltm_pool" "test-pool" {
        name = "` + TEST_POOL_NAME + `"
        monitors = ["/Common/http"]
        allow_nat = "yes"
        allow_snat = "yes"
        description = "Test-Pool-Sample"
        load_balancing_mode = "round-robin"
        slow_ramp_time = "5"
        service_down_action = "reset"
        reselect_tries = "2"
}
resource "bigip_ltm_pool_attachment" "test-pool_test-node" {
	pool = bigip_ltm_pool.test-pool.name
	node = "` + poolMember + `"
    ratio                 = 2
    connection_limit      = 2
    connection_rate_limit = 2
    priority_group        = 2
    dynamic_ratio         = 3
}
resource "bigip_ltm_policy" "test-policy" {
        depends_on = ["bigip_ltm_pool.test-pool"]
        name = "` + TestPolicyName + `"
        strategy = "first-match"
        requires = ["http"]
#       published_copy = "Drafts/` + TestPolicyName + `"
        controls = ["forwarding"]
        rule  {
              name = "rule6"
                      action {
//                            tm_name = "20"
                              forward = true
                              pool = "/Common/test-pool"
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
				Config: TestPolicyResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists("/Common/test-policy-again"),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_create_newpoolbehavior(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPolicyResource2,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
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
				Config: TestPolicyResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
				),
				ResourceName:      TestPolicyName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigipLtmPolicy_import_newpoolbehavior(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPolicyResource2,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
				),
				ResourceName:      TestPolicyName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		re := regexp.MustCompile("/([a-zA-z0-9? ,_-]+)/([a-zA-z0-9? ,_-]+)")
		match := re.FindStringSubmatch(name)
		if match == nil {
			return fmt.Errorf("Failed to match regex in :%v ", name)
		}
		partition := match[1]
		policyName := match[2]
		policy, err := client.GetPolicy(policyName, partition)
		if err != nil {
			return fmt.Errorf("Error while fetching policy: %v ", err)
		}
		if policy == nil {
			return fmt.Errorf("Policy %s was not created ", name)
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
		re := regexp.MustCompile("/([a-zA-z0-9? ,_-]+)/([a-zA-z0-9? ,_-]+)")
		match := re.FindStringSubmatch(name)
		if match == nil {
			return fmt.Errorf("Failed to match regex :%v ", name)
		}
		partition := match[1]
		policyName := match[2]
		policy, err := client.GetPolicy(policyName, partition)
		if err != nil {
			return nil
		}
		if policy != nil {
			return fmt.Errorf("Policy %s not destroyed ", name)
		}
	}
	return nil
}
