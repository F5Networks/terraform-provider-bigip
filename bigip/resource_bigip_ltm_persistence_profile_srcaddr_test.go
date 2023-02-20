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

var TestPpsrcaddrName = fmt.Sprintf("/%s/test-ppsrcaddr", TestPartition)

var TestPpsrcaddrResource = `
resource "bigip_ltm_persistence_profile_srcaddr" "test_ppsrcaddr" {
	name = "` + TestPpsrcaddrName + `"
	defaults_from = "/Common/source_addr"
	match_across_pools = "enabled"
	match_across_services = "enabled"
	match_across_virtuals = "enabled"
	mirror = "enabled"
	timeout = 3600
	override_conn_limit = "enabled"
	hash_algorithm = "carp"
	map_proxies = "enabled"
	mask = "255.255.255.255"
	app_service = "none"
}

`

func TestAccBigipLtmPersistenceProfileSrcAddrCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testCheckBigipLtmPersistenceProfileSrcAddrDestroyed),
		Steps: []resource.TestStep{
			{
				Config: TestPpsrcaddrResource,
				Check: resource.ComposeTestCheckFunc(
					testBigipLtmPersistenceProfileSrcAddrExists(TestPpsrcaddrName, true),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "name", TestPpsrcaddrName),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "defaults_from", "/Common/source_addr"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "match_across_pools", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "match_across_services", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "match_across_virtuals", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "mirror", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "timeout", "3600"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "override_conn_limit", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "hash_algorithm", "carp"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "map_proxies", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_srcaddr.test_ppsrcaddr", "mask", "255.255.255.255"),
				),
			},
		},
	})

}

func TestAccBigipLtmPersistenceProfileSrcAddrImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckBigipLtmPersistenceProfileSrcAddrDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPpsrcaddrResource,
				Check: resource.ComposeTestCheckFunc(
					testBigipLtmPersistenceProfileSrcAddrExists(TestPpsrcaddrName, true),
				),
				ResourceName:      TestPpsrcaddrName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testBigipLtmPersistenceProfileSrcAddrExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		pp, err := client.GetSourceAddrPersistenceProfile(name)
		if err != nil {
			return err
		}
		if exists && pp == nil {
			return fmt.Errorf("Source Address Persistence Profile %s does not exist.", name)
		}
		if !exists && pp != nil {
			return fmt.Errorf("Source Address Persistence Profile %s exists.", name)
		}
		return nil
	}
}

func testCheckBigipLtmPersistenceProfileSrcAddrDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_persistence_profile_srcaddr" {
			continue
		}

		name := rs.Primary.ID
		pp, err := client.GetSourceAddrPersistenceProfile(name)
		if err != nil {
			return err
		}

		if pp != nil {
			return fmt.Errorf("Source Address Persistence Profile %s not destroyed.", name)
		}
	}
	return nil
}
