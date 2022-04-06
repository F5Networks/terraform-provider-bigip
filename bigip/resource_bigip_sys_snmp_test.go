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

var TEST_SNMP_NAME = fmt.Sprintf("/%s/test-snmp", TEST_PARTITION)

var TEST_SNMP_RESOURCE = `
resource "bigip_sys_snmp" "test-snmp" {
  sys_contact = "NetOPsAdmin s.shitole@f5.com"
  sys_location = "SeattleHQ"
  allowedaddresses = ["202.10.10.2"]
}
`

func TestAccBigipSyssnmp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TEST_SNMP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnmpExists(TEST_SNMP_NAME, true),
					resource.TestCheckResourceAttr("bigip_sys_snmp.test-snmp", "sys_contact", "NetOPsAdmin s.shitole@f5.com"),
					resource.TestCheckResourceAttr("bigip_sys_snmp.test-snmp", "sys_location", "SeattleHQ"),
					resource.TestCheckResourceAttr("bigip_sys_snmp.test-snmp",
						fmt.Sprintf("allowedaddresses.%d", schema.HashString("202.10.10.2")),
						"202.10.10.2"),
				),
			},
		},
	})
}

func TestAccBigipSyssnmp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TEST_SNMP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testChecksnmpExists(TEST_SNMP_NAME, true),
				),
				ResourceName:      TEST_SNMP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testChecksnmpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.SNMPs()
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("snmp %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("snmp %s still exists.", name)
		}
		return nil
	}
}
