/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccDataSourceBigipNetTrunk_basic exercises the data source against a
// real device. There is no bigip_net_trunk resource, so a fixture trunk is
// created/destroyed via bigip_command tmsh invocations (mirroring the
// bigip_net_trunks acceptance test's approach of relying on real physical
// interfaces).
func TestAccDataSourceBigipNetTrunk_basic(t *testing.T) {
	trunkName := "test-trunk-ds"
	var freeInterface string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			freeInterface = testAcctPreCheckFreeInterface(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTrunkDestroyed(trunkName),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipNetTrunkConfig(trunkName, freeInterface),
				Check: resource.ComposeTestCheckFunc(
					testCheckTrunkExists(trunkName),
					resource.TestCheckResourceAttr("data.bigip_net_trunk.test", "name", trunkName),
					resource.TestCheckResourceAttr("data.bigip_net_trunk.test", "interfaces.#", "1"),
					resource.TestCheckResourceAttrPtr("data.bigip_net_trunk.test", "interfaces.0", &freeInterface),
					resource.TestCheckResourceAttrSet("data.bigip_net_trunk.test", "lacp_mode"),
					resource.TestCheckResourceAttrSet("data.bigip_net_trunk.test", "lacp_timeout"),
				),
				// The "when = destroy" bigip_command resource leaves its
				// Optional+Computed command_result unset until its Delete
				// runs (Create no-ops when when != "apply"), so Terraform
				// always sees a non-empty plan on the post-apply refresh.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDataSourceBigipNetTrunkConfig(trunkName, iface string) string {
	return fmt.Sprintf(`
resource "bigip_command" "create_trunk" {
  commands = ["create net trunk %[1]s interfaces add { %[2]s }"]
  when     = "apply"
}

resource "bigip_command" "delete_trunk" {
  commands   = ["delete net trunk %[1]s"]
  when       = "destroy"
  depends_on = [bigip_command.create_trunk]
}

data "bigip_net_trunk" "test" {
  name       = "%[1]s"
  depends_on = [bigip_command.create_trunk]
}
`, trunkName, iface)
}

func testCheckTrunkExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		trunks, err := client.Trunks()
		if err != nil {
			return err
		}
		for _, trunk := range trunks.Trunks {
			if trunk.Name == name {
				return nil
			}
		}
		return fmt.Errorf("trunk %s was not created", name)
	}
}

func testCheckTrunkDestroyed(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		trunks, err := client.Trunks()
		if err != nil {
			return err
		}
		for _, trunk := range trunks.Trunks {
			if trunk.Name == name {
				return fmt.Errorf("trunk %s not destroyed", name)
			}
		}
		return nil
	}
}
