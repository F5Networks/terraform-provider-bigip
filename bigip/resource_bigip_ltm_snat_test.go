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

var TEST_SNAT_NAME = fmt.Sprintf("/%s/test-snat", TestPartition)

var TEST_SNAT_RESOURCE = `
resource "bigip_ltm_snat" "test-snat" {
 name = "` + TEST_SNAT_NAME + `"
 translation = "/Common/136.1.1.1"
 origins { name = "2.2.2.2" }
 origins { name = "3.3.3.3" }
 vlansdisabled = true
 autolasthop = "default"
 mirror = "disabled"
 partition = "Common"
 full_path = "/Common/test-snat"
} `

func TestAccBigipLtmsnat_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksnatsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SNAT_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckSnatExists(TEST_SNAT_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_snat.test-snat", "name", TEST_SNAT_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_snat.test-snat", "translation", "/Common/136.1.1.1"),
					resource.TestCheckResourceAttr("bigip_ltm_snat.test-snat", "origins.0.name", "2.2.2.2"),
					resource.TestCheckResourceAttr("bigip_ltm_snat.test-snat", "origins.1.name", "3.3.3.3"),
					resource.TestCheckResourceAttr("bigip_ltm_snat.test-snat", "vlansdisabled", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_snat.test-snat", "autolasthop", "default"),
					resource.TestCheckResourceAttr("bigip_ltm_snat.test-snat", "mirror", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_snat.test-snat", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ltm_snat.test-snat", "full_path", "/Common/test-snat"),
				),
			},
		},
	})
}

func TestAccBigipLtmsnat_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testChecksnatsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SNAT_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckSnatExists(TEST_SNAT_NAME, true),
				),
				ResourceName:      TEST_SNAT_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckSnatExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetSnat(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("Snat %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("Snat %s still exists.", name)
		}
		return nil
	}
}

func testChecksnatsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_snat" {
			continue
		}

		name := rs.Primary.ID
		snat, err := client.GetSnat(name)
		if err != nil {
			return err
		}
		if snat != nil {
			return fmt.Errorf("Snat %s not destroyed.", name)
		}
	}
	return nil
}
