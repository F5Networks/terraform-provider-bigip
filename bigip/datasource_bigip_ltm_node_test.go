/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigipLtmNode_basic(t *testing.T) {
	t.Parallel()
	resName := "bigip_ltm_node.NODETEST"
	dataSourceName := "data.bigip_ltm_node.NODETEST"
	var nodeName = "/Common/test-node"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNodeConfigBasic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(nodeName),
					resource.TestCheckResourceAttrPair(dataSourceName, "full_path", resName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "address", resName, "address"),
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
  name    = "%s"
  address = "192.168.30.1"
}

# We can't easily reference the node resource above because name includes the
# partition. Instead we have to split and pull out the separate pieces.
data "bigip_ltm_node" "NODETEST" {
  name      = split("/", bigip_ltm_node.NODETEST.name)[2]
  partition = split("/", bigip_ltm_node.NODETEST.name)[1]
}
`, nodeName)
}
