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

var TestPpcookieName = fmt.Sprintf("/%s/test-ppcookie", TestPartition)

var TestPpcookieResource = `
resource "bigip_ltm_persistence_profile_cookie" "test_ppcookie" {
	name = "` + TestPpcookieName + `"
	defaults_from = "/Common/cookie"
	match_across_pools = "enabled"
	match_across_services = "enabled"
	match_across_virtuals = "enabled"
	timeout = 3600
	override_conn_limit = "enabled"
	always_send = "enabled"
	cookie_encryption = "disabled"
	cookie_name = "ham"
	expiration = "1:0:0"
	hash_length = 0
	app_service = "none"
}

`

func TestAccBigipLtmPersistenceProfileCookieCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testCheckBigipLtmPersistenceProfileCookieDestroyed),
		Steps: []resource.TestStep{
			{
				Config: TestPpcookieResource,
				Check: resource.ComposeTestCheckFunc(
					testBigipLtmPersistenceProfileCookieExists(TestPpcookieName, true),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "name", TestPpcookieName),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "defaults_from", "/Common/cookie"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "match_across_pools", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "match_across_services", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "match_across_virtuals", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "timeout", "3600"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "override_conn_limit", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "always_send", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "cookie_encryption", "disabled"),
					// unable to validate since value is encrypted
					// resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "cookie_encryption_passphrase", "iloveham"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "cookie_name", "ham"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "expiration", "1:0:0"),
					resource.TestCheckResourceAttr("bigip_ltm_persistence_profile_cookie.test_ppcookie", "hash_length", "0"),
				),
			},
		},
	})

}

func TestAccBigipLtmPersistenceProfileCookieImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckBigipLtmPersistenceProfileCookieDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPpcookieResource,
				Check: resource.ComposeTestCheckFunc(
					testBigipLtmPersistenceProfileCookieExists(TestPpcookieName, true),
				),
				ResourceName:      TestPpcookieName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testBigipLtmPersistenceProfileCookieExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		pp, err := client.GetCookiePersistenceProfile(name)
		if err != nil {
			return err
		}
		if exists && pp == nil {
			return fmt.Errorf("Cookie Persistence Profile %s does not exist.", name)
		}
		if !exists && pp != nil {
			return fmt.Errorf("Cookie Persistence Profile %s exists.", name)
		}
		return nil
	}
}

func testCheckBigipLtmPersistenceProfileCookieDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_persistence_profile_cookie" {
			continue
		}

		name := rs.Primary.ID
		pp, err := client.GetSourceAddrPersistenceProfile(name)
		if err != nil {
			return err
		}

		if pp != nil {
			return fmt.Errorf("Cookie Persistence Profile %s not destroyed.", name)
		}
	}
	return nil
}
