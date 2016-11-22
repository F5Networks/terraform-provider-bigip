package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
	"testing"
)

var TEST_POLICY_NAME = fmt.Sprintf("/%s/test-policy", TEST_PARTITION)
var TEST_RULE_NAME = fmt.Sprintf("/%s/test-rule", TEST_PARTITION)

var TEST_POLICY_RESOURCE = `
resource "bigip_ltm_pool" "test-pool" {
	name = "` + TEST_POOL_NAME + `"
	monitors = ["/Common/http"]
	allow_nat = true
	allow_snat = true
	load_balancing_mode = "round-robin"
}

resource "bigip_ltm_policy" "test-policy" {
	name = "` + TEST_POLICY_NAME + `"
	controls = ["forwarding"]
	requires = ["http"]
	rule {
		name = "` + TEST_RULE_NAME + `"
		condition {
        	        httpUri = true
                	startsWith = true
                	values = ["/foo", "/bar"]
                }

                condition {
                	httpMethod = true
                	values = ["GET"]
                }

                action {
                	forward = true
                	pool = "${bigip_ltm_pool.test-pool.name}"
                }
	}
}
`

func TestBigipLtmPolicy_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckPolicyDestroyed,
			testCheckPoolsDestroyed,
		),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_POLICY_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TEST_POLICY_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "name", TEST_POLICY_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy",
						fmt.Sprintf("controls.%d", schema.HashString("forwarding")),
						"forwarding"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy",
						fmt.Sprintf("requires.%d", schema.HashString("http")),
						"http"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "rule.0.name", TEST_RULE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "rule.0.condition.0.httpUri", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "rule.0.condition.0.startsWith", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "rule.0.condition.0.values.0", "/foo"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "rule.0.condition.0.values.1", "/bar"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "rule.0.condition.1.httpMethod", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "rule.0.condition.1.values.0", "GET"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "rule.0.action.0.forward", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "rule.0.action.0.pool", TEST_POOL_NAME),
				),
			},
		},
	})
}

func TestBigipLtmPolicy_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPolicyDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_POLICY_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPolicyExists(TEST_POLICY_NAME, true),
				),
				ResourceName:      TEST_POLICY_NAME,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckPolicyExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.GetPolicy(TEST_POLICY_NAME)
		if err != nil {
			return err
		}

		if p == nil {
			return fmt.Errorf("Policy %s not created.", TEST_POLICY_NAME)
		}

		return nil
	}
}

func testCheckPolicyDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_policy" {
			continue
		}

		name := rs.Primary.ID
		p, err := client.GetPolicy(name)
		if err != nil {
			return err
		}
		if p != nil {
			return fmt.Errorf("Virtual address %s not destroyed.", name)
		}
	}
	return nil
}
