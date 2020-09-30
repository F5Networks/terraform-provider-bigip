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

var TestHttpName = fmt.Sprintf("/%s/test-http", TEST_PARTITION)

var TestHttpResource = `
resource "bigip_ltm_profile_http" "test-http" {
  name = "/Common/test-http"
  defaults_from = "/Common/http"
  description = "some http"
  fallback_host = "titanic"
  fallback_status_codes = ["400","500","300"]
}
`

func TestAccBigipLtmProfileHttpCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestHttpResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(TestHttpName, true),
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
func TestAccBigipLtmProfileHttpUpdateServerAgent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofilehttpDefaultConfig(TEST_PARTITION, TestHttpName, "http-profile-test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(TestHttpName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.http-profile-test", "name", TestHttpName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.http-profile-test", "defaults_from", "/Common/http"),
				),
			},
			{
				Config: testaccbigipltmprofilehttpUpdateServeragentConfig(TEST_PARTITION, TestHttpName, "http-profile-test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(TestHttpName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.http-profile-test", "name", TestHttpName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.http-profile-test", "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_http.http-profile-test", "server_agent_name", "myBIG-IP"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccBigipLtmHttpProfileImportConfig(),
			},
			{
				ResourceName:      "bigip_ltm_profile_http.test-http-profile",
				ImportStateId:     "/Common/test-http",
				ImportState:       true,
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

func testaccbigipltmprofilehttpDefaultConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http" "%[3]s" {
  name = "%[2]s"
  defaults_from = "/%[1]s/http"
}
`, partition, profileName, resourceName)
}

func testaccbigipltmprofilehttpUpdateServeragentConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http" "%[3]s" {
  name = "%[2]s"
  defaults_from = "/%[1]s/http"
  server_agent_name = "myBIG-IP"
}
`, partition, profileName, resourceName)
}

func testaccBigipLtmHttpProfileImportConfig() string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http" "test-http" {
  name = "%s"
  defaults_from = "/Common/http"
}`, "/Common/test-http")
}
