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
	var fqName = "google.com"
	var address = "192.168.1.1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAcctPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckNodesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNodeConfigBasic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testCheckNodeExists(nodeName, true),
					resource.TestCheckResourceAttr(resName, "name", nodeName),
					resource.TestCheckResourceAttr(resName, "address", address),
					resource.TestCheckResourceAttr(resName, "fqdn.0.name", fqName),
					resource.TestCheckResourceAttr(resName, "fqdn.0.address_family", "ipv4"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNodeConfigBasic(nodeName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_node" "NODETEST" {
	name = "%s"
    address = 192.168.1.1
	fqdn {
		name = google.com
		address_family = ipv4
	}
}`, nodeName)
}
