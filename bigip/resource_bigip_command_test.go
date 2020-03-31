/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var TestCommandResource = `
resource "bigip_command" "test-command" {
  name      = "command1"
  commands   = ["show sys version"]
}
`
var testCmd = "show sys version"
var testCmdResult = "\nSys::Version\nMain Package\n  Product     BIG-IP\n  Version     14.1.0.3\n  Build       0.0.6\n  Edition     Point Release 3\n  Date        Mon Mar 25 17:15:27 PDT 2019\n\n"

func TestAccBigipCommand_run(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestCommandResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_command.test-command", "commands.0", testCmd),
					resource.TestCheckResourceAttr("bigip_command.test-command", "command_result.0", testCmdResult),
				),
			},
		},
	})
}
