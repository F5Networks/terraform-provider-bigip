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

var TestPpsslName = fmt.Sprintf("/%s/test-ppssl", TestPartition)

var TestPpsslResource = `
resource "bigip_ltm_persistence_profile_ssl" "test_ppssl" {
	name = "` + TestPpsslName + `"
	defaults_from = "/Common/ssl"
	match_across_pools = "enabled"
	match_across_services = "enabled"
	match_across_virtuals = "enabled"
	mirror = "enabled"
	timeout = 3600
	override_conn_limit = "enabled"
}

`

func TestAccBigipLtmPersistenceProfileSSLCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testCheckBigipLtmPersistenceProfileSSLDestroyed),
		Steps: []resource.TestStep{
			{
				Config: TestPpsslResource,
				Check: resource.ComposeTestCheckFunc(
					testBigipLtmPersistenceProfileSSLExists(TestPpsslName, true),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_ssl.test_ppssl", "name", TestPpsslName),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_ssl.test_ppssl", "defaults_from", "/Common/ssl"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_ssl.test_ppssl", "match_across_pools", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_ssl.test_ppssl", "match_across_services", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_ssl.test_ppssl", "match_across_virtuals", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_ssl.test_ppssl", "mirror", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_ssl.test_ppssl", "timeout", "3600"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_ssl.test_ppssl", "override_conn_limit", "enabled"),
				),
			},
		},
	})

}

func TestAccBigipLtmPersistenceProfileSSLImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckBigipLtmPersistenceProfileSSLDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPpsslResource,
				Check: resource.ComposeTestCheckFunc(
					testBigipLtmPersistenceProfileSSLExists(TestPpsslName, true),
				),
				ResourceName:      TestPpsslName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testBigipLtmPersistenceProfileSSLExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		pp, err := client.GetSSLPersistenceProfile(name)
		if err != nil {
			return err
		}
		if exists && pp == nil {
			return fmt.Errorf("SSL Persistence Profile %s does not exist.", name)
		}
		if !exists && pp != nil {
			return fmt.Errorf("SSL Persistence Profile %s exists.", name)
		}
		return nil
	}
}

func testCheckBigipLtmPersistenceProfileSSLDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_persistence_profile_ssl" {
			continue
		}

		name := rs.Primary.ID
		pp, err := client.GetSourceAddrPersistenceProfile(name)
		if err != nil {
			return err
		}

		if pp != nil {
			return fmt.Errorf("SSL Persistence Profile %s not destroyed.", name)
		}
	}
	return nil
}
