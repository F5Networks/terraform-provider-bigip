/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccBigipLtmNode_basic(t *testing.T) {
	t.Parallel()
	resName := "bigip_ltm_node.NODETEST"
	var nodeName = "test-node"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNodeConfigBasic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(nodeName, true),
					resource.TestCheckResourceAttr(resName, "name", "/Common/test-node"),
					resource.TestCheckResourceAttr(resName, "address", "192.168.30.1"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckNodeConfigBasic(nodeName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_node" "NODETEST" {
	name = "/%s/%s"
    address = "192.168.30.1"
}`, "Common", nodeName)
}
