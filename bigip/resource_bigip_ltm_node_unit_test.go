package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
)

func testBigipLtmNodeInvalid(resourceName string) string {
	return fmt.Sprintf(`
		resource "bigip_ltm_node" "test-node" {
			name = "%s"
			address = "10.10.10.10"
	        invalidkey = "foo"
		}
		provider "bigip" {
			address = "10.10.10.1"
			username = "admin"
			password = "admin"
		}
	`, resourceName)
}

func TestBigipLtmNodeInvalid(t *testing.T) {
	resourceName := "/Common/test-node"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmNodeInvalid(resourceName),
				ExpectError: regexp.MustCompile("invalid or unknown key: invalidkey"),
			},
		},
	})
}
