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
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_VLAN_NAME = fmt.Sprintf("/%s/test-vlan", TEST_PARTITION)

var TEST_VLAN_RESOURCE = `
resource "bigip_net_vlan" "test-vlan" {
	name = "/Common/test-vlan"
	tag = 101
	interfaces {
		vlanport = 1.1
		tagged = true
	}
}
`

func TestAccBigipNetvlan_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckvlansDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_VLAN_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckvlanExists(TEST_VLAN_NAME, true),
					resource.TestCheckResourceAttr("bigip_net_vlan.test-vlan", "name", "/Common/test-vlan"),
					resource.TestCheckResourceAttr("bigip_net_vlan.test-vlan", "tag", "101"),
					resource.TestCheckResourceAttr("bigip_net_vlan.test-vlan", "interfaces.0.vlanport", "1.1"),
					resource.TestCheckResourceAttr("bigip_net_vlan.test-vlan", "interfaces.0.tagged", "true"),
				),
			},
		},
	})
}

func TestAccBigipNetvlan_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckvlansDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_VLAN_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckvlanExists(TEST_VLAN_NAME, true),
				),
				ResourceName:      TEST_VLAN_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckvlanExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.Vlan(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("vlan %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("vlan %s still exists.", name)
		}
		return nil
	}
}

func testCheckvlansDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	time.Sleep(2 * time.Second)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_net_vlan" {
			continue
		}

		name := rs.Primary.ID
		vlan, err := client.Vlan(name)
		if err != nil {
			return err
		}
		if vlan == nil {
			return fmt.Errorf("vlan %s not destroyed.", name)
		}
	}
	return nil
}
