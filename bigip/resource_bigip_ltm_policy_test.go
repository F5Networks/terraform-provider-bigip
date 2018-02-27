package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

//var TEST_POLICY_NAME = "/" + TEST_PARTITION + "/test-policy"
var TEST_POLICY_NAME = "test-policy"

//var TEST_POOL_NAME = fmt.Sprintf("/%s/test-pool", TEST_PARTITION)
//var TEST_POOLNODE_NAME = fmt.Sprintf("/%s/test-node", TEST_PARTITION)
//var TEST_POOLNODE_NAMEPORT = fmt.Sprintf("%s:443", TEST_POOLNODE_NAME)

var TEST_POLICY_RESOURCE = `
resource "bigip_ltm_node" "test-node" {
	name = "` + TEST_NODE_NAME + `"
	address = "10.10.10.10"
}
resource "bigip_ltm_pool" "test-pool" {
	name = "` + TEST_POOL_NAME + `"
	monitors = ["/Common/http"]
	allow_nat = "yes"
	allow_snat = "yes"
	load_balancing_mode = "round-robin"
	depends_on = ["bigip_ltm_node.test-node"]
	nodes = ["` + TEST_POOLNODE_NAMEPORT + `"]
}
resource "bigip_ltm_policy" "test-policy" {
	depends_on = ["bigip_ltm_pool.test-pool"]
 name = "` + TEST_POLICY_NAME + `"
 strategy = "first-match"
  requires = ["http"]
 published_copy = "Drafts/test-policy"
  controls = ["forwarding"]
  rule  {
  name = "rule6"

   action = {
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
					/*resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "name", TEST_POLICY_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "strategy", "first-match"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy",
						fmt.Sprintf("requires.%d", schema.HashString("http")),
						"http"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy", "published_copy", "Drafts/test-policy"),
					resource.TestCheckResourceAttr("bigip_ltm_policy.test-policy",
						fmt.Sprintf("controls.%d", schema.HashString("forwarding")),
						"forwarding"),  */
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
