package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_HTTPCOMPRESS_NAME = fmt.Sprintf("/%s/test-httpcompress", TEST_PARTITION)

var TEST_HTTPCOMPRESS_RESOURCE = `
resource "bigip_ltm_profile_httpcompress" "test-httpcompress"

        {
            name = "/Common/sanjose-httpcompress"
			defaults_from = "/Common/httpcompression"
            uri_exclude = "/ABCD"
            uri_include = "/XYZ"
        }
`

func TestBigipLtmProfileHttpcompress_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpcompresssDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_HTTPCOMPRESS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(TEST_HTTPCOMPRESS_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "name", "/Common/sanjose-httpcompress"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "defaults_from", "/Common/httpcompression"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "uri_exclude", "/ABCD"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "uri_include", "/XYZ"),
				),
			},
		},
	})
}

func TestBigipLtmProfileHttpcompress_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpcompresssDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_HTTPCOMPRESS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(TEST_HTTPCOMPRESS_NAME, true),
				),
				ResourceName:      TEST_HTTPCOMPRESS_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckHttpcompressExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.Httpcompress()
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("httpcompress %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("httpcompress %s was not created.", name)
		}
		return nil
	}
}

func testCheckHttpcompresssDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_httpcompress" {
			continue
		}

		name := rs.Primary.ID
		httpcompress, err := client.Httpcompress()
		if err != nil {
			return err
		}
		if httpcompress == nil {
			return fmt.Errorf("httpcompress %s was not created.", name)
		}
	}
	return nil
}
