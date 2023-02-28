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

var TestPpdstaddrName = fmt.Sprintf("/%s/test-ppdstaddr", TestPartition)

var TestPpdstaddrResource = `
resource "bigip_ltm_persistence_profile_dstaddr" "test_ppdstaddr" {
	name = "` + TestPpdstaddrName + `"
	defaults_from = "/Common/dest_addr"
	match_across_pools = "enabled"
	match_across_services = "enabled"
	match_across_virtuals = "enabled"
	mirror = "enabled"
	timeout = 3600
	override_conn_limit = "enabled"
	hash_algorithm = "carp"
	mask = "255.255.255.255"
	app_service = "none"
}

`

func TestAccBigipLtmPersistenceProfileDstAddrCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testCheckBigipLtmPersistenceProfileDstAddrDestroyed),
		Steps: []resource.TestStep{
			{
				Config: TestPpdstaddrResource,
				Check: resource.ComposeTestCheckFunc(
					testBigipLtmPersistenceProfileDstAddrExists(TestPpdstaddrName, true),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "name", TestPpdstaddrName),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "defaults_from", "/Common/dest_addr"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "match_across_pools", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "match_across_services", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "match_across_virtuals", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "mirror", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "timeout", "3600"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "override_conn_limit", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "hash_algorithm", "carp"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_dstaddr.test_ppdstaddr", "mask", "255.255.255.255"),
				),
			},
		},
	})

}

func TestAccBigipLtmPersistenceProfileDstAddrImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckBigipLtmPersistenceProfileDstAddrDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPpdstaddrResource,
				Check: resource.ComposeTestCheckFunc(
					testBigipLtmPersistenceProfileDstAddrExists(TestPpdstaddrName, true),
				),
				ResourceName:      TestPpdstaddrName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testBigipLtmPersistenceProfileDstAddrExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		pp, err := client.GetDestAddrPersistenceProfile(name)
		if err != nil {
			return err
		}
		if exists && pp == nil {
			return fmt.Errorf("Destination Address Persistence Profile %s does not exist.", name)
		}
		if !exists && pp != nil {
			return fmt.Errorf("Destination Address Persistence Profile %s exists.", name)
		}
		return nil
	}
}

func testCheckBigipLtmPersistenceProfileDstAddrDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_persistence_profile_dstaddr" {
			continue
		}

		name := rs.Primary.ID
		pp, err := client.GetDestAddrPersistenceProfile(name)
		if err != nil {
			return err
		}

		if pp != nil {
			return fmt.Errorf("Destination Address Persistence Profile %s not destroyed.", name)
		}
	}
	return nil
}
