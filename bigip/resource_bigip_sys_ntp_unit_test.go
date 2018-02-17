package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
)

func testBigipSysNtpInvalid(resourceName string) string {
	return fmt.Sprintf(`
		resource "bigip_sys_ntp" "test-ntp" {
			description = "%s"
			servers = ["10.10.10.10"]
	    timezone = "America/Los_Angeles"
			invalidkey = "foo"
		}
		provider "bigip" {
			address = "10.10.10.1"
			username = "admin"
			password = "admin"
		}
	`, resourceName)
}

func TestAccBigipSysNtpInvalid(t *testing.T) {
	resourceName := "/Common/test-ntp"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysNtpInvalid(resourceName),
				ExpectError: regexp.MustCompile("invalid or unknown key: invalidkey"),
			},
		},
	})
}
