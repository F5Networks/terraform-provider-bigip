package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
	published_copy = "Drafts/test-policy"
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
