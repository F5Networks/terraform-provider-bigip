/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"regexp"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestWebAccelerationName = fmt.Sprintf("/%s/test", TestPartition)
var resWebAccelerationName = "bigip_ltm_profile_web_acceleration"

var TestWebAccelerationResource = `
resource "bigip_ltm_profile_web_acceleration" "web_acceleration" {
  name = "/Common/web_acceleration"
  defaults_from = "/Common/webacceleration"
}
`

func TestAccBigipLtmWebAccelerationProfileCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestWebAccelerationResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(TestWebAccelerationName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_web_acceleration.web_acceleration", "name", "/Common/web_acceleration"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_web_acceleration.web_acceleration", "defaults_from", "/Common/webacceleration"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileCreateFail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestWebAccelerationResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(TestWebAccelerationName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_web_acceleration.web_acceleration", "name", "/Common/web_acceleration"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_web_acceleration.web_acceleration", "defaults_from", "/Common/web_acceleration"),
				),
				ExpectError: regexp.MustCompile("Attribute 'defaults_from' expected"),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheSize(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheSize"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_size"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckResourceAttr(resFullName, "cache_size", "100"),
				),
			},
		},
	})
}
func TestAccBigipLtmWebAccelerationProfileUpdateCacheMaxEntries(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheMaxEntries"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_max_entries"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckResourceAttr(resFullName, "cache_max_entries", "201"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheMaxAge(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheMaxAge"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_max_age"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckResourceAttr(resFullName, "cache_max_age", "301"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheObjectMinSize(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheObjectMinSize"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_object_min_size"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckResourceAttr(resFullName, "cache_object_min_size", "501"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheObjectMaxSize(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheObjectMaxSize"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_object_max_size"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckResourceAttr(resFullName, "cache_object_max_size", "502"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheUriExclude(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheUriExclude"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_uri_exclude"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckTypeSetElemAttr(resFullName, "cache_uri_exclude.*", "exclude"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheUriInclude(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheUriInclude"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_uri_include"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckTypeSetElemAttr(resFullName, "cache_uri_include.*", "include"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheUriIncludeOverride(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheUriIncludeOverride"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_uri_include_override"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckTypeSetElemAttr(resFullName, "cache_uri_include_override.*", "include_override"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheUriPinned(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheUriPinned"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_uri_pinned"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckTypeSetElemAttr(resFullName, "cache_uri_pinned.*", "pinned"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheClientCacheControlMode(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheClientCacheControlMode"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_client_cache_control_mode"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckResourceAttr(resFullName, "cache_client_cache_control_mode", "all"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheInsertAgeHeader(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheInsertAgeHeader"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_insert_age_header"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckResourceAttr(resFullName, "cache_insert_age_header", "disabled"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileUpdateCacheAgingRate(t *testing.T) {
	t.Parallel()
	var instName = "web_acceleration-Update-CacheAgingRate"
	var instFullName = fmt.Sprintf("/%s/%s", TestPartition, instName)
	resFullName := fmt.Sprintf("%s.%s", resWebAccelerationName, instName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
				),
			},
			{
				Config: testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, "cache_aging_rate"),
				Check: resource.ComposeTestCheckFunc(
					testCheckWebAccelerationExists(instFullName),
					resource.TestCheckResourceAttr(resFullName, "name", instFullName),
					resource.TestCheckResourceAttr(resFullName, "defaults_from", "/Common/webacceleration"),
					resource.TestCheckResourceAttr(resFullName, "cache_aging_rate", "9"),
				),
			},
		},
	})
}

func TestAccBigipLtmWebAccelerationProfileImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckWebAccelerationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccBigipLtmWebAccelerationProfileImportConfig(),
			},
			{
				ResourceName:      "bigip_ltm_profile_web_acceleration.test-web-acceleration",
				ImportStateId:     "/Common/test-web-acceleration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckWebAccelerationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetWebAccelerationProfile(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("web acceleration %s was not created", name)
		}

		return nil
	}
}

func testCheckWebAccelerationDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_web_acceleration" {
			continue
		}

		name := rs.Primary.ID
		http, err := client.GetWebAccelerationProfile(name)
		if err != nil {
			return err
		}
		if http != nil {
			return fmt.Errorf("web acceleration %s not destroyed ", name)
		}
	}
	return nil
}

func testaccBigipLtmWebAccelerationProfileImportConfig() string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_web_acceleration" "test-web-acceleration" {
  name          = "%s"
  defaults_from = "/Common/webacceleration"
}
`, "/Common/test-web-acceleration")
}

func testAccBigipLtmWebAccelerationProfileDefaultConfig(instName, updateParam string) string {
	resPrefix := fmt.Sprintf(`
		resource "%[1]s" "%[2]s" {
			  name = "/Common/%[2]s"
			  defaults_from = "/Common/webacceleration"`, resWebAccelerationName, instName)
	switch updateParam {
	case "cache_size":
		resPrefix = fmt.Sprintf(`%s
			  cache_size = 100`, resPrefix)
	case "cache_max_entries":
		resPrefix = fmt.Sprintf(`%s
			  cache_max_entries = 201`, resPrefix)
	case "cache_max_age":
		resPrefix = fmt.Sprintf(`%s
			  cache_max_age = 301`, resPrefix)
	case "cache_object_min_size":
		resPrefix = fmt.Sprintf(`%s
		cache_object_min_size = 501`, resPrefix)
	case "cache_object_max_size":
		resPrefix = fmt.Sprintf(`%s
		cache_object_max_size = 502`, resPrefix)
	case "cache_uri_exclude":
		resPrefix = fmt.Sprintf(`%s
			  cache_uri_exclude = ["exclude"]`, resPrefix)
	case "cache_uri_include":
		resPrefix = fmt.Sprintf(`%s
			  cache_uri_include =  ["include"]`, resPrefix)
	case "cache_uri_include_override":
		resPrefix = fmt.Sprintf(`%s
			  cache_uri_include_override = ["include_override"]`, resPrefix)
	case "cache_uri_pinned":
		resPrefix = fmt.Sprintf(`%s
			  cache_uri_pinned = ["pinned"]`, resPrefix)
	case "cache_client_cache_control_mode":
		resPrefix = fmt.Sprintf(`%s
			  cache_client_cache_control_mode = "all"`, resPrefix)
	case "cache_insert_age_header":
		resPrefix = fmt.Sprintf(`%s
			  cache_insert_age_header = "disabled"`, resPrefix)
	case "cache_aging_rate":
		resPrefix = fmt.Sprintf(`%s
			  cache_aging_rate = 9`, resPrefix)
	default:
	}
	return fmt.Sprintf(`%s
		}`, resPrefix)
}
