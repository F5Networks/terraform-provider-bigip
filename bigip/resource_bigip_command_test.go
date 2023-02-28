/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var TestCommandResource = `
resource "bigip_command" "test-command" {
  commands   = ["show sys version"]
}
`
var testCmd = "show sys version"

func TestAccBigipCommand_run(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestCommandResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_command.test-command", "commands.0", testCmd),
					resource.TestMatchResourceAttr("bigip_command.test-command", "command_result.0", regexp.MustCompile("^\nSys::Version\nMain Package\n {2}Product {5}BIG-IP\n {2}Version")),
				),
			},
		},
	})
}
