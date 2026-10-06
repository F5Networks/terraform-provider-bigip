/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TEST_SNATPOOL_NAME = fmt.Sprintf("/%s/test-snatpool", TestPartition)

var TEST_SNATPOOL_RESOURCE = `
resource "bigip_ltm_snatpool" "test-snatpool" {
  name = "` + TEST_SNATPOOL_NAME + `"
  members = ["/Common/191.1.1.1","/Common/194.2.2.2"]
}

`

func TestAccBigipLtmsnatpool_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksnatpoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SNATPOOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnatpoolExists(TEST_SNATPOOL_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_snatpool.test-snatpool", "name", TEST_SNATPOOL_NAME),
					// members is a TypeSet; checking a specific element's
					// value requires TestCheckTypeSetElemAttr (which hashes
					// internally and handles SDK internals correctly), not
					// a manually-computed schema.HashString index -- the
					// SDK no longer supports indexing into a TypeSet that
					// way ("likely indexes into TypeSet - This is
					// currently not possible in the SDK").
					resource.TestCheckTypeSetElemAttr("bigip_ltm_snatpool.test-snatpool", "members.*", "/Common/191.1.1.1"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_snatpool.test-snatpool", "members.*", "/Common/194.2.2.2"),
				),
			},
		},
	})
}

func TestAccBigipLtmsnatpool_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksnatpoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SNATPOOL_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnatpoolExists(TEST_SNATPOOL_NAME, true),
				),
				ResourceName:      TEST_SNATPOOL_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testChecksnatpoolExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.Snatpools(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("snatpool %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("snatpool %s still exists.", name)
		}
		return nil
	}
}

func testChecksnatpoolsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_snatpool" {
			continue
		}

		name := rs.Primary.ID
		snatpool, err := client.Snatpools(name)
		if err != nil {
			return err
		}
		if snatpool == nil {
			return fmt.Errorf("snatpool %s not destroyed.", name)
		}
	}
	return nil
}
