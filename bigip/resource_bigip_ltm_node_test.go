package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
	"testing"
)

var TEST_NODE_NAME = fmt.Sprintf("/%s/test-node", TEST_PARTITION)

var TEST_NODE_RESOURCE = `
resource "bigip_ltm_node" "test-node" {
	name = "` + TEST_NODE_NAME + `"
	address = "10.10.10.10"
}
`

func TestBigipLtmNode_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_NODE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TEST_NODE_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "name", TEST_NODE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_node.test-node", "address", "10.10.10.10"),
				),
			},
		},
	})
}

func TestBigipLtmNode_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_NODE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(TEST_NODE_NAME, true),
				),
				ResourceName:      TEST_NODE_NAME,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//var TEST_NODE_IN_POOL_RESOURCE = `
//resource "bigip_ltm_pool" "test-pool" {
//	name = "` + TEST_POOL_NAME + `"
//  	load_balancing_mode = "round-robin"
//  	nodes = ["${formatlist("%s:80", bigip_ltm_node.*.name)}"]
//  	allow_snat = false
//}
//`
//func TestBigipLtmNode_removeNode(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAcctPreCheck(t)
//		},
//		Providers: testAccProviders,
//		CheckDestroy: testCheckNodesDestroyed,
//		Steps: []resource.TestStep{
//			resource.TestStep{
//				Config: TEST_NODE_RESOURCE + TEST_NODE_IN_POOL_RESOURCE,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckNodeExists(TEST_NODE_NAME, true),
//					testCheckPoolExists(TEST_POOL_NAME, true),
//					testCheckPoolMember(TEST_POOL_NAME, TEST_NODE_NAME),
//				),
//			},
//			resource.TestStep{
//				Config: TEST_NODE_IN_POOL_RESOURCE,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckNodeExists(fmt.Sprintf("%s:%s", TEST_NODE_NAME, "80"), false),
//					testCheckEmptyPool(TEST_POOL_NAME),
//				),
//			},
//		},
//	})
//}

func testCheckNodeExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		node, err := client.GetNode(name)
		if err != nil {
			return err
		}
		if exists && node == nil {
			return fmt.Errorf("Node ", name, " was not created.")
		}
		if !exists && node != nil {
			return fmt.Errorf("Node ", name, " still exists.")
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
			return fmt.Errorf("Node ", name, " not destroyed.")
		}
	}
	return nil
}
