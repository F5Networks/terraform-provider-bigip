package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_NTP_NAME = fmt.Sprintf("/%s/test-ntp", TEST_PARTITION)

var TEST_NTP_RESOURCE = `
resource "bigip_ntp" "test-ntp" {
	description = "` + TEST_NTP_NAME + `"
	servers = ["10.10.10.10"]
	timezone = "America/Los_Angeles"
}
`

func TestBigipLtmNtp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//CheckDestroy: testCheckntpsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_NTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckntpExists(TEST_NTP_NAME, true),
					resource.TestCheckResourceAttr("bigip_ntp.test-ntp", "description", TEST_NTP_NAME),
					resource.TestCheckResourceAttr("bigip_ntp.test-ntp", "servers", "[10.10.10.10]"),
					resource.TestCheckResourceAttr("bigip_ntp.test-ntp", "timezone", "America/Los_Angeles"),
				),
			},
		},
	})
}

func TestBigipLtmNtp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//	CheckDestroy: testCheckntpsDestroyed, ( No Delet API support)
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_NTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckntpExists(TEST_NTP_NAME, true),
				),
				ResourceName:      TEST_NTP_NAME,
				ImportState:       false,
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

func testCheckntpExists(description string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		ntp, err := client.NTPs()
		if err != nil {
			return err
		}
		if exists && ntp == nil {
			return fmt.Errorf("ntp ", description, " was not created.")

		}
		if !exists && ntp != nil {
			return fmt.Errorf("ntp ", description, " still exists.")

		}
		return nil
	}
}

func testCheckntpsDestroyed(s *terraform.State) error {
	/* client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ntp" {
			continue
		}

		description := rs.Primary.ID
		ntp, err := client.NTPs()
		if err != nil {
			return err
		}
		if ntp != nil {
			return fmt.Errorf("ntp ", description, " not destroyed.")

		}
	}*/
	return nil
}
