package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_FASTHTTP_NAME = fmt.Sprintf("/%s/test-fasthttp", TEST_PARTITION)

var TEST_FASTHTTP_RESOURCE = `
resource "bigip_fasthttp_profile" "test-fasthttp" {
	name = "` + TEST_FASTHTTP_NAME + `"
	defaults_from = "/Common/fasthttp"
            idle_timeout = 300
            connpoolidle_timeoutoverride	= 0
            connpool_maxreuse = 2
            connpool_maxsize  = 2048
            connpool_minsize = 0
            connpool_replenish = "enabled"
            connpool_step = 4
            forcehttp_10response = "disabled"
            maxheader_size = 32768
}
`

func TestBigipLtmfasthttp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfasthttpsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_FASTHTTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckfasthttpExists(TEST_FASTHTTP_NAME, true),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "name", TEST_FASTHTTP_NAME),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "defaults_from", "Common/fasthttp"),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "idle_timeout", "300"),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "connpoolidle_timeoutoverride", "0"),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "connpool_maxreuse", "2"),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "connpool_maxsize", "2048"),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "connpool_minsize", "0"),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "connpool_replenish", "enabled"),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "connpool_step", "4"),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "forcehttp_10response", "disabled"),
					resource.TestCheckResourceAttr("bigip_fasthttp_profile.test-fasthttp", "maxheader_size", "32768"),
				),
			},
		},
	})
}

func TestBigipLtmfasthttp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfasthttpsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_FASTHTTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckfasthttpExists(TEST_FASTHTTP_NAME, true),
				),
				ResourceName:      TEST_FASTHTTP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

//var TEST_FASTHTTP_IN_POOL_RESOURCE = `
//resource "bigip_ltm_pool" "test-pool" {
//	name = "` + TEST_POOL_NAME + `"
//  	load_balancing_mode = "round-robin"
//  	fasthttps = ["${formatlist("%s:80", bigip_fasthttp_profile.*.name)}"]
//  	allow_snat = false
//}
//`
//func TestBigipLtmfasthttp_removefasthttp(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAcctPreCheck(t)
//		},
//		Providers: testAccProviders,
//		CheckDestroy: testCheckfasthttpsDestroyed,
//		Steps: []resource.TestStep{
//			resource.TestStep{
//				Config: TEST_FASTHTTP_RESOURCE + TEST_FASTHTTP_IN_POOL_RESOURCE,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckfasthttpExists(TEST_FASTHTTP_NAME, true),
//					testCheckPoolExists(TEST_POOL_NAME, true),
//					testCheckPoolMember(TEST_POOL_NAME, TEST_FASTHTTP_NAME),
//				),
//			},
//			resource.TestStep{
//				Config: TEST_FASTHTTP_IN_POOL_RESOURCE,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckfasthttpExists(fmt.Sprintf("%s:%s", TEST_FASTHTTP_NAME, "80"), false),
//					testCheckEmptyPool(TEST_POOL_NAME),
//				),
//			},
//		},
//	})
//}

func testCheckfasthttpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		fasthttp, err := client.Fasthttp()
		if err != nil {
			return err
		}
		if exists && fasthttp == nil {
			return fmt.Errorf("fasthttp ", name, " was not created.")
		}
		if !exists && fasthttp != nil {
			return fmt.Errorf("fasthttp ", name, " still exists.")
		}
		return nil
	}
}

func testCheckfasthttpsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_fasthttp_profile" {
			continue
		}

		name := rs.Primary.ID
		fasthttp, err := client.Fasthttp()
		if err != nil {
			return err
		}
		if fasthttp != nil {
			return fmt.Errorf("fasthttp ", name, " not destroyed.")
		}
	}
	return nil
}
