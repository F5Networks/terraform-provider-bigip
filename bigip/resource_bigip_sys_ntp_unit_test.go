/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testBigipSysNtpInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_ntp" "test-ntp" {
  description = "%s"
  servers     = ["10.10.10.10"]
  timezone    = "America/Los_Angeles"
  invalidkey  = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxxx"
  password = "xxxxx"
}
	`, resourceName)
}

func TestAccBigipSysNtpInvalid(t *testing.T) {
	resourceName := "/Common/test-ntp"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysNtpInvalid(resourceName),
				ExpectError: regexp.MustCompile("Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}
