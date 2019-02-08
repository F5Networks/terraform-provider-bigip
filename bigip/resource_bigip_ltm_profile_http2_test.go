package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_HTTP2_NAME = fmt.Sprintf("/%s/test-http2", TEST_PARTITION)

var TEST_HTTP2_RESOURCE = `
resource "bigip_ltm_profile_http2" "test-http2"

        {
            name = "/Common/test-http2"
	    defaults_from = "/Common/http2"
            concurrent_streams_per_connection = 10
            connection_idle_timeout = 30
            activation_modes = ["alpn"]
        }
`

func TestAccBigipLtmProfileHttp2_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_HTTP2_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TEST_HTTP2_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http2.test-http2", "name", "/Common/test-http2"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http2.test-http2", "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http2.test-http2", "concurrent_streams_per_connection", "10"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http2.test-http2", "connection_idle_timeout", "30"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttp2_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_HTTP2_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TEST_HTTP2_NAME, true),
				),
				ResourceName:      TEST_HTTP2_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckHttp2Exists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetHttp2(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("http2 %s was not created.", name)
		}
		if !exists && p == nil {
			return fmt.Errorf("http2 %s still exists.", name)
		}
		return nil
	}
}

func testCheckHttp2sDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_http2" {
			continue
		}

		name := rs.Primary.ID
		http2, err := client.GetHttp2(name)
		if err != nil {
			return err
		}
		if http2 != nil {
			return fmt.Errorf("http2 %s not destroyed.", name)
		}
	}
	return nil
}
