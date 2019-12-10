/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_HTTP_NAME = fmt.Sprintf("/%s/test-http", TEST_PARTITION)

var TEST_HTTP_RESOURCE = `
resource "bigip_ltm_profile_http" "test-http" {
  name = "/Common/test-http"
  defaults_from = "/Common/http"
  description = "some http"
  fallback_host = "titanic"
  fallback_status_codes = ["400","500","300"]
}
`

func TestAccBigipLtmProfilehttp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_HTTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(TEST_HTTP_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.test-http", "name", "/Common/test-http"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.test-http", "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.test-http", "description", "some http"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.test-http", "fallback_host", "titanic"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.test-http",
						fmt.Sprintf("fallback_status_codes.%d", schema.HashString("400")),
						"400"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.test-http",
						fmt.Sprintf("fallback_status_codes.%d", schema.HashString("500")),
						"500"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.test-http",
						fmt.Sprintf("fallback_status_codes.%d", schema.HashString("300")),
						"300"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfilehttp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_HTTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(TEST_HTTP_NAME, true),
				),
				ResourceName:      TEST_HTTP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckhttpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetHttpProfile(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("http %s was not created.", name)
		}
		if !exists && p == nil {
			return fmt.Errorf("http %s still exists.", name)
		}
		return nil
	}
}

func testCheckHttpsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_http" {
			continue
		}

		name := rs.Primary.ID
		http, err := client.GetHttpProfile(name)
		if err != nil {
			return err
		}
		if http != nil {
			return fmt.Errorf("http %s not destroyed.", name)
		}
	}
	return nil
}
