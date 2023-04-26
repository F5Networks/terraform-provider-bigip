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
	"strings"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestPolicyName = "/Common/test-policy"

var TestPolicyResource = `
resource "bigip_ltm_pool" "test-pool" {
	name = "` + TestPoolName + `"
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
	controls = ["forwarding"]
	rule  {
	      name = "rule6"
		      action {
			      forward    = true
				  connection = false
			      pool       = "/Common/test-pool"
		      }
	}
}
resource "bigip_ltm_policy" "test-policy-again" {
  name = "/Common/test-policy-again"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule  {
    name = "testrule"
    action {
      redirect   = true
	  connection = false
      location   = "tcl:https://[HTTP::host][HTTP::uri]"
      http_reply = true
    }
  }
}
`
var TestPolicyResource2 = `
resource "bigip_ltm_pool" "test-pool" {
        name = "` + TestPoolName + `"
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
//                            tm_name    = "20"
                              forward    = true
							  connection = false
                              pool       = "/Common/test-pool"
                      }
        }
}
`

var TestPolicyResource3 = `
resource "bigip_ltm_pool" "test-policy-pool" {
  name = "/Common/test-policy-pool"
}
resource "bigip_ltm_policy" "test-policy-rules" {
	depends_on = ["bigip_ltm_pool.test-policy-pool"]
	name       = "/Common/test-policy-rules"
	strategy   = "first-match"
	requires   = ["client-ssl"]
	controls   = ["forwarding"]
	rule {
	  name = "Rule-01"
	  condition {
		ssl_extension    = true
		server_name      = true
		ends_with        = true
		ssl_client_hello = true
		values = [
		  "domain1.net",
		  "domain2.nl"
		]
	  }
	  action {
		forward          = true
		connection       = false
		pool             = bigip_ltm_pool.test-policy-pool.name
		ssl_client_hello = true
	  }
	}
	rule {
	  name = "lastrule-deny"
	  action {
		shutdown         = true
		ssl_client_hello = true
	  }
	}
  }
`

var TestPolicyResource4 = `
resource "bigip_ltm_pool" "test-policy-pool" {
  name = "/Common/test-policy-pool"
}
resource "bigip_ltm_policy" "test-policy-rules" {
	depends_on = ["bigip_ltm_pool.test-policy-pool"]
	name       = "/Common/test-policy-rules"
	strategy   = "first-match"
	requires   = ["client-ssl"]
	controls   = ["forwarding"]
	rule {
	  name = "Rule-01"
	  condition {
		ssl_extension    = true
		server_name      = true
		ends_with        = true
		ssl_client_hello = true
		values = [
		  "domain1.net",
		  "domain2.nl"
		]
	  }
	  action {
		forward          = true
		connection       = false
		pool             = bigip_ltm_pool.test-policy-pool.name
		ssl_client_hello = true
	  }
	}
	rule {
	  name = "Rule-02"
	  condition {
		ssl_extension = true
		server_name   = true
		ends_with     = true
		ssl_client_hello = true
		values = [
		  "domain3.net",
		  "domain4.nl"
		]
	  }
	  action {
		forward          = true
		connection       = false
		pool             = bigip_ltm_pool.test-policy-pool.name
		ssl_client_hello = true
	  }
	}
	rule {
	  name = "lastrule-deny"
	  action {
		shutdown         = true
		ssl_client_hello = true
	  }
	}
  }
`

var TestPolicyResource5 = `
resource "bigip_ltm_policy" "test-policy-condition" {
  name     = "/Common/test-policy-condition"
  strategy = "first-match"
  requires = ["http"]
  rule {
	name = "replace_if_exists"
	action {
	  replace     = true
	  http_header = true
	  connection  = false
	  tm_name     = "X-Forwarded"
	  value       = "https"
	}
	condition {
	  http_header = true
	  tm_name     = "X-Forwarded"
	  exists      = true
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

func TestAccBigipLtmPolicy_create_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPolicyResource3,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists("/Common/test-policy-rules"),
				),
			},
			{
				Config: TestPolicyResource4,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists("/Common/test-policy-rules"),
				),
			},
			{
				Config: TestPolicyResource5,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists("/Common/test-policy-condition"),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_Issue132(t *testing.T) {
	t.Parallel()
	TestPolicyName = "/Common/test-policy-issue132"
	resName := fmt.Sprintf("%s.%s", "bigip_ltm_policy", "test-policy-issue132")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmPoolicyIssu132and133(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists(TestPolicyName),
					resource.TestCheckResourceAttr(resName, "name", TestPolicyName),
					resource.TestCheckResourceAttr(resName, "strategy", "first-match"),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_Issue132_a(t *testing.T) {
	t.Parallel()
	TestPolicyName = "/Common/test-policy-issue132-a"
	resName := fmt.Sprintf("%s.%s", "bigip_ltm_policy", "test-policy-issue132-a")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmpoolicyissu132A(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists(TestPolicyName),
					resource.TestCheckResourceAttr(resName, "name", TestPolicyName),
					resource.TestCheckResourceAttr(resName, "strategy", "all-match"),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_Issue132_b(t *testing.T) {
	t.Parallel()
	TestPolicyName = "/Common/test-policy-issue132-b"
	resName := fmt.Sprintf("%s.%s", "bigip_ltm_policy", "test-policy-issue132-b")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmpoolicyissu132B(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists(TestPolicyName),
					resource.TestCheckResourceAttr(resName, "name", TestPolicyName),
					resource.TestCheckResourceAttr(resName, "strategy", "best-match"),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_Issue132_c(t *testing.T) {
	t.Parallel()
	TestPolicyName = "/TEST/test-policy-issue132-c"
	resName := fmt.Sprintf("%s.%s", "bigip_ltm_policy", "test-policy-issue132-c")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmpoolicyissu132C(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists(TestPolicyName),
					resource.TestCheckResourceAttr(resName, "name", TestPolicyName),
					resource.TestCheckResourceAttr(resName, "strategy", "best-match"),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_Issue591(t *testing.T) {
	t.Parallel()
	TestPolicyName = "/Common/policy-issue-591"
	resName := fmt.Sprintf("%s.%s", "bigip_ltm_policy", "policy-issue-591")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmpoolicyissue591(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists(TestPolicyName),
					resource.TestCheckResourceAttr(resName, "name", TestPolicyName),
					resource.TestCheckResourceAttr(resName, "strategy", "first-match"),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_Issue634(t *testing.T) {
	t.Parallel()
	TestPolicyName = "/TEST/A1/test-policy-issue634"
	resName := fmt.Sprintf("%s.%s", "bigip_ltm_policy", "test-policy-issue634")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmpoolicyissue634(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists(TestPolicyName),
					resource.TestCheckResourceAttr(resName, "name", TestPolicyName),
					resource.TestCheckResourceAttr(resName, "strategy", "first-match"),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_Issue634_a(t *testing.T) {
	t.Parallel()
	TestPolicyName = "/TEST/test-policy-issue634-a"
	resName := fmt.Sprintf("%s.%s", "bigip_ltm_policy", "test-policy-issue634-a")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmpoolicyissue634a(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists(TestPolicyName),
					resource.TestCheckResourceAttr(resName, "name", TestPolicyName),
					resource.TestCheckResourceAttr(resName, "strategy", "first-match"),
				),
			},
		},
	})
}

func TestAccBigipLtmPolicy_Issue648(t *testing.T) {
	t.Parallel()
	TestPolicyName = "/Common/testpolicy-issue-648"
	resName := fmt.Sprintf("%s.%s", "bigip_ltm_policy", "testpolicy-issue-648")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmpoolicyissue648(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists(TestPolicyName),
					resource.TestCheckResourceAttr(resName, "name", TestPolicyName),
					resource.TestCheckResourceAttr(resName, "strategy", "first-match"),
				),
			},
		},
	})
}
func TestAccBigipLtmPolicyIssue737(t *testing.T) {
	t.Parallel()
	TestPolicyName = "/Common/testpolicy-issue-737"
	resName := fmt.Sprintf("%s.%s", "bigip_ltm_policy", "testpolicy-issue-737")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicysDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmpoolicyissue737(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TestPolicyName),
					testCheckPolicyExists(TestPolicyName),
					resource.TestCheckResourceAttr(resName, "name", TestPolicyName),
					resource.TestCheckResourceAttr(resName, "strategy", "first-match"),
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

func TestAccBigipLtmPolicyIssue794TestCases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: loadFixtureString("../examples/bigip_ltm_policy_issue794.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists("/Common/policy-issue-591"),
					testCheckPolicyExists("/Common/policy_issue794_tc1"),
					testCheckPolicyExists("/Common/policy_issue794_tc2"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.policy_issue794_tc1", "strategy", "first-match"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.policy_issue794_tc1", "name", "/Common/policy_issue794_tc1"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.policy_issue794_tc2", "strategy", "first-match"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.policy_issue794_tc2", "name", "/Common/policy_issue794_tc2"),
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
		polStr := strings.Split(name, "/")
		partition := strings.Join(polStr[:len(polStr)-1], "/")
		policyName := polStr[len(polStr)-1]
		// partition := match[1]
		// policyName := match[2]
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
		polStr := strings.Split(name, "/")
		partition := strings.Join(polStr[:len(polStr)-1], "/")
		policyName := polStr[len(polStr)-1]
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

func testaccbigipltmPoolicyIssu132and133() string {
	tfConfig := `
		resource "bigip_ltm_policy" "test-policy-issue132" {
		  name     = "/Common/test-policy-issue132"
		}`
	return tfConfig
}
func testaccbigipltmpoolicyissu132A() string {
	tfConfig := `
		resource "bigip_ltm_policy" "test-policy-issue132-a" {
		  name     = "/Common/test-policy-issue132-a"
          strategy = "/Common/all-match"
		}`
	return tfConfig
}

func testaccbigipltmpoolicyissu132B() string {
	tfConfig := `
		resource "bigip_ltm_policy" "test-policy-issue132-b" {
		  name     = "/Common/test-policy-issue132-b"
		  strategy = "/Common/best-match"
		}`
	return tfConfig
}

func testaccbigipltmpoolicyissu132C() string {
	tfConfig := `
		resource "bigip_ltm_policy" "test-policy-issue132-c" {
		  name     = "/TEST/test-policy-issue132-c"
		  strategy = "/Common/best-match"
		}`
	return tfConfig
}

func testaccbigipltmpoolicyissue591() string {
	tfConfig := `
		resource "bigip_ltm_pool" "k8s_prod" {
  			name = "/Common/k8prod_Pool"
		}
		resource "bigip_ltm_policy" "policy-issue-591" {
		  name     = "/Common/policy-issue-591"
		  strategy = "first-match"
		  requires = ["http"]
		  controls = ["forwarding"]
		  rule {
			name = "rule-issue591"
			condition {
			  index     = 0
			  http_host = true
			  contains  = true
			  values = [
				"domain1.net",
				"domain2.nl"
			  ]
			  request = true
			}
			condition {
			  http_uri    = true
			  path        = true
			  not         = true
			  starts_with = true
			  values      = ["/role-service"]
			  request     = true
			}
			action {
			  forward  = false
			  replace  = true
			  connection = false
			  http_uri = true
			  path     = "tcl:[string map {/role-service/ /} [HTTP::uri]]"
			  request  = true
			}
			action {
			  forward    = true
			  connection = false
			  pool       = bigip_ltm_pool.k8s_prod.name
			}
		  }
		}`
	return tfConfig
}

func testaccbigipltmpoolicyissue634() string {
	tfConfig := `
		resource "bigip_ltm_policy" "test-policy-issue634" {
		  name     = "/TEST/A1/test-policy-issue634"
		}`
	return tfConfig
}

func testaccbigipltmpoolicyissue634a() string {
	tfConfig := `
		resource "bigip_ltm_policy" "test-policy-issue634-a" {
		  name     = "/TEST/test-policy-issue634-a"
		}`
	return tfConfig
}
func testaccbigipltmpoolicyissue648() string {
	tfConfig := `
		resource "bigip_ltm_pool" "k8s_prod" {
  			name = "/Common/k8prod_Pool"
		}
		resource "bigip_ltm_policy" "testpolicy-issue-648" {
		  name     = "/Common/testpolicy-issue-648"
		  strategy = "first-match"
		  requires = ["tcp", "client-ssl"]
		  controls = ["forwarding"]
		  rule {
			name = "Rule-01"
			condition {
			  ssl_extension    = true
			  ssl_client_hello = true
			  server_name      = true
			  ends_with        = true
			  values = [
				"domain1.net",
				"domain2.nl"
			  ]
			}
			condition {
			  tcp             = true
			  matches         = true
			  address         = true
			  client_accepted = true
			  values = [
				"10.0.0.0/8",
				"20.0.0.0/8",
			  ]
			}
			action {
			  forward          = true
			  connection       = false
			  pool             = bigip_ltm_pool.k8s_prod.name
			  ssl_client_hello = true
			}
		  }
		
		  rule {
			name = "Rule-02"
			condition {
			  ssl_extension    = true
			  ssl_client_hello = true
			  server_name      = true
			  ends_with        = true
			  values = [
				"domain3.net",
				"domain4.nl"
			  ]
			}
			condition {
			  tcp             = true
			  matches         = true
			  address         = true
			  client_accepted = true
			  values = [
				"30.0.0.0/8",
				"40.0.0.0/8",
			  ]
			}
			action {
			  forward          = true
			  connection       = false
			  pool             = bigip_ltm_pool.k8s_prod.name
			  ssl_client_hello = true
			}
		  }
		  rule {
			name = "lastrule-deny"
			action {
			  shutdown         = true
			  ssl_client_hello = true
			}
		  }
		}`
	return tfConfig
}

func testaccbigipltmpoolicyissue737() string {
	tfConfig := `
resource "bigip_ltm_policy" "policy" {
  name     = "/Common/f5-policy"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding", "asm"]

  ### IMPORTANT! Please do not change rules order as it matters!
  # Cluster rules
  dynamic "rule" {
    for_each = { for r in local.rules : r.name => r }

    content {
      name = "${rule.value.name}-https"

      action {
        connection = false
        forward    = true
        asm        = false
        disable    = false
        pool       = "/Common/${rule.value.pool}"
      }
      condition {
        http_host = true
        host      = true
        ends_with = true
        request   = true
        values    = ["${rule.value.matching}"]
      }
      action {
        asm        = true
        connection = false
        enable     = rule.value.enable
        disable    = rule.value.disable
        select     = rule.value.select
        policy     = rule.value.policy
      }
    }
  }
  rule {
    name = "default-rule"
    action {
      asm        = true
      connection = false
      disable    = true
      request    = true
    }
  }
}

locals {
  rules = [
    {
      name = "rule1"
      pool = "pool1"
      matching = "example.com"
      enable = true
      select = false
      disable = false
      policy = "/Common/f5-waf-profile"
    },
    {
      name = "rule2"
      pool = "pool2"
      matching = "example2.com"
      enable = false
      disable = true
      select = false
      policy = ""
    },
    {
      name = "rule3"
      pool = "pool3"
      matching = "example3.com"
      enable = true
      disable = false
      select = false
      policy = "/Common/f5-waf-profile"
    }
  ]
}
`
	return tfConfig
}
