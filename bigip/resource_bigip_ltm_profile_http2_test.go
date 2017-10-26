package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_HTTP2_NAME = fmt.Sprintf("/%s/test-http2", TEST_PARTITION)

var TEST_HTTP2_RESOURCE = `
resource "bigip_http2_profile" "test-http2"

        {
            name = "/Common/sanjose-http2"
						defaults_from = "/Common/http2"
            concurrent_streams_per_connection = 10
            connection_idle_timeout = 30
            activation_modes = ["alpn","npn"]
        }
`

func TestBigipLtmHttp2_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_HTTP2_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttp2Exists(TEST_HTTP2_NAME, true),
					resource.TestCheckResourceAttr("bigip_http2_profile.test-http2", "name", "/Common/sanjose-http2"),
					resource.TestCheckResourceAttr("bigip_http2_profile.test-http2", "defaults_from", "/Common/http2"),
					resource.TestCheckResourceAttr("bigip_http2_profile.test-http2", "concurrent_streams_per_connection", "10"),
					resource.TestCheckResourceAttr("bigip_http2_profile.test-http2", "connection_idle_timeout", "30"),
				),
			},
		},
	})
}

func TestBigipLtmHttp2_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttp2sDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
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
		p, err := client.Http2()
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("fastl4 ", name, " was not created.")
		}
		if !exists && p != nil {
			return fmt.Errorf("fastl4 ", name, " still exists.")
		}
		return nil
	}
}

func testCheckHttp2sDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_http2_profile" {
			continue
		}

		name := rs.Primary.ID
		http2, err := client.Http2()
		if err != nil {
			return err
		}
		if http2 == nil {
			return fmt.Errorf("fasthttp ", name, " not destroyed.")
		}
	}
	return nil
}
