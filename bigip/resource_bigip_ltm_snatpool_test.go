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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TestSnatpoolName = fmt.Sprintf("/%s/test-snatpool", TEST_PARTITION)

var TestSnatpoolResource = `
resource "bigip_ltm_snatpool" "test-snatpool" {
  name = "` + TestSnatpoolName + `"
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
				Config: TestSnatpoolResource,
				Check: resource.ComposeTestCheckFunc(
					testChecksnatpoolExists(TestSnatpoolName, true),
					resource.TestCheckResourceAttr("bigip_ltm_snatpool.test-snatpool", "name", TestSnatpoolName),
					resource.TestCheckResourceAttr("bigip_ltm_snatpool.test-snatpool",
						fmt.Sprintf("members.%d", schema.HashString("/Common/191.1.1.1")),
						"/Common/191.1.1.1"),
					resource.TestCheckResourceAttr("bigip_ltm_snatpool.test-snatpool",
						fmt.Sprintf("members.%d", schema.HashString("/Common/194.2.2.2")),
						"/Common/194.2.2.2"),
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
				Config: TestSnatpoolResource,
				Check: resource.ComposeTestCheckFunc(
					testChecksnatpoolExists(TestSnatpoolName, true),
				),
				ResourceName:      TestSnatpoolName,
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
