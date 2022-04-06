/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
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

var TestHttpName = fmt.Sprintf("/%s/test-http", TEST_PARTITION)
var resHttpName = "bigip_ltm_profile_http"

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
					testCheckhttpExists(TestHttpName),
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
	t.Parallel()
	var instName = "test-http-Update-serveragent"
	var TestHttpName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttpName, "http-profile-test")
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
					testCheckhttpExists(TestHttpName),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttpName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
				),
			},
			{
				Config: testaccbigipltmprofilehttpUpdateServeragentConfig(TEST_PARTITION, TestHttpName, "http-profile-test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(TestHttpName),
					resource.TestCheckResourceAttr(resFullName, "name", TestHttpName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr(resFullName, "server_agent_name", "myBIG-IP"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpUpdateFallbackhost(t *testing.T) {
	t.Parallel()
	var instName = "test-http-Update-FallbackHost"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttpName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
				),
			},
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, "fallback_host"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr(resFullName, "fallback_host", "titanic"),
				),
			},
		},
	})
}
func TestAccBigipLtmProfileHttpUpdateBasicAuthRealm(t *testing.T) {
	t.Parallel()
	var instName = "test-http-Update-BasicAuthRealm"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttpName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
				),
			},
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, "basic_auth_realm"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr(resFullName, "basic_auth_realm", "titanic"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpUpdateHeaderErase(t *testing.T) {
	t.Parallel()
	var instName = "test-http-Update-headerErase"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttpName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
				),
			},
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, "head_erase"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr(resFullName, "head_erase", "titanic"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpUpdateDescription(t *testing.T) {
	t.Parallel()
	var instName = "test-http-Update-desciption"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttpName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
				),
			},
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, "description"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr(resFullName, "description", "my-http-profile"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpUpdateFallbackStatusCodes(t *testing.T) {
	t.Parallel()
	var instName = "test-http-Update-fallbackStatusCodes"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttpName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
				),
			},
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, "fallback_status_codes"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("fallback_status_codes.%d", schema.HashString("300")), "300"),
					resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("fallback_status_codes.%d", schema.HashString("500")), "500"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpUpdateHeaderInsert(t *testing.T) {
	t.Parallel()
	var instName = "test-http-Update-headerInsert"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttpName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
				),
			},
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, "head_insert"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr(resFullName, "head_insert", "X-Forwarded-IP: [expr { [IP::client_addr] }]"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpUpdateEncryptCookies(t *testing.T) {
	t.Parallel()
	var instName = "test-http-Update-encryptCookies"
	var instFullName = fmt.Sprintf("/%s/%s", TEST_PARTITION, instName)
	resFullName := fmt.Sprintf("%s.%s", resHttpName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
				),
			},
			{
				Config: testaccbigipltmprofilehttpUpdateParam(instName, "encrypt_cookies"),
				Check: resource.ComposeTestCheckFunc(
					testCheckhttpExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/http"),
					resource.TestCheckResourceAttr(resFullName, fmt.Sprintf("encrypt_cookies.%d", schema.HashString("peanutButter")), "peanutButter"),
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

func testCheckhttpExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetHttpProfile(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("http %s was not created", name)
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
			return fmt.Errorf("http %s not destroyed ", name)
		}
	}
	return nil
}

func testaccbigipltmprofilehttpDefaultConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http" "%[3]s" {
  name          = "%[2]s"
  defaults_from = "/%[1]s/http"
}
`, partition, profileName, resourceName)
}

func testaccbigipltmprofilehttpUpdateServeragentConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http" "%[3]s" {
  name              = "%[2]s"
  defaults_from     = "/%[1]s/http"
  server_agent_name = "myBIG-IP"
}
`, partition, profileName, resourceName)
}

func testaccBigipLtmHttpProfileImportConfig() string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http" "test-http" {
  name          = "%s"
  defaults_from = "/Common/http"
}
`, "/Common/test-http")
}

func testaccbigipltmprofilehttpUpdateParam(instName, updateParam string) string {
	resPrefix := fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			  name = "/Common/%[2]s"
			  //defaults_from = "/Common/http"`, resHttpName, instName)
	switch updateParam {
	case "fallback_host":
		resPrefix = fmt.Sprintf(`%s
			  fallback_host = "titanic"`, resPrefix)
	case "fallback_status_codes":
		resPrefix = fmt.Sprintf(`%s
			  fallback_host = "titanic"
			  fallback_status_codes = ["300","500"]`, resPrefix)
	case "encrypt_cookies":
		resPrefix = fmt.Sprintf(`%s
			  encrypt_cookies = ["peanutButter"]`, resPrefix)
	case "head_erase":
		resPrefix = fmt.Sprintf(`%s
			  head_erase = "titanic"`, resPrefix)
	case "description":
		resPrefix = fmt.Sprintf(`%s
			  description = "my-http-profile"`, resPrefix)
	case "head_insert":
		resPrefix = fmt.Sprintf(`%s
			  head_insert = "X-Forwarded-IP: [expr { [IP::client_addr] }]"`, resPrefix)
	case "insert_xforwarded_for":
		resPrefix = fmt.Sprintf(`%s
			  insert_xforwarded_for = 262100`, resPrefix)
	case "lws_separator":
		resPrefix = fmt.Sprintf(`%s
			  lws_separator = 2400`, resPrefix)
	case "oneconnect_transformations":
		resPrefix = fmt.Sprintf(`%s
			  oneconnect_transformations = 40`, resPrefix)
	case "proxy_type":
		resPrefix = fmt.Sprintf(`%s
			  proxy_type = 40`, resPrefix)
	case "redirect_rewrite":
		resPrefix = fmt.Sprintf(`%s
			  redirect_rewrite = "AES"`, resPrefix)
	case "basic_auth_realm":
		resPrefix = fmt.Sprintf(`%s
			  basic_auth_realm = "titanic"`, resPrefix)
	default:
	}
	return fmt.Sprintf(`%s
		}`, resPrefix)
}
