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

var TestHttpcompressName = fmt.Sprintf("/%s/test-httpcompress", TestPartition)

var TestHttpcompressResource = `
resource "bigip_ltm_profile_httpcompress" "test-httpcompress" {
            name = "/Common/test-httpcompress"
	    defaults_from = "/Common/httpcompression"
            uri_exclude = ["f5.com"]
            uri_include = ["cisco.com"]
	    content_type_include = ["nicecontent.com"]
	    content_type_exclude = ["nicecontentexclude.com"]
        }
`

func TestAccBigipLtmProfileHttpcompress_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpcompresssDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestHttpcompressResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(TestHttpcompressName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "name", "/Common/test-httpcompress"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "defaults_from", "/Common/httpcompression"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "uri_exclude.*", "f5.com"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "uri_include.*", "cisco.com"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "content_type_include.*", "nicecontent.com"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "content_type_exclude.*", "nicecontentexclude.com"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpcompressTC1(t *testing.T) {
	profileHttpComprsName := fmt.Sprintf("/%s/%s", "Common", "test_httpcompress_profiletc1")
	httpsTenantName = "fast_https_tenanttc1"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpcompresssDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getProfileHttpComprsConfig(profileHttpComprsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(profileHttpComprsName, true),
					testCheckHttpcompressExists("/Common/xxx_tets_compre", false),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "name", profileHttpComprsName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "defaults_from", "/Common/httpcompression"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "compression_buffersize", "4090"),
				),
			},
			{
				Config: getProfileHttpComprsConfig(profileHttpComprsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(profileHttpComprsName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "name", profileHttpComprsName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "defaults_from", "/Common/httpcompression"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "compression_buffersize", "4090"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpcompressTC2(t *testing.T) {
	profileHttpComprsName := fmt.Sprintf("/%s/%s", "Common", "test_httpcompress_profiletc2")
	httpsTenantName = "fast_https_tenanttc2"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpcompresssDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getProfileHttpComprsTC2Config(profileHttpComprsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(profileHttpComprsName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "name", profileHttpComprsName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "defaults_from", "/Common/httpcompression"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "gzip_compression_level", "2"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "gzip_memory_level", "32768"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "gzip_window_size", "32768"),
				),
			},
			{
				Config: getProfileHttpComprsTC2Config(profileHttpComprsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(profileHttpComprsName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "name", profileHttpComprsName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "defaults_from", "/Common/httpcompression"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "gzip_compression_level", "2"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "gzip_memory_level", "32768"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test_httpcomprs_profile", "gzip_window_size", "32768"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileHttpcompress_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckHttpcompresssDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestHttpcompressResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(TestHttpcompressName, true),
				),
				ResourceName:      TestHttpcompressName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckHttpcompressExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetHttpcompress(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("httpcompress %s was not created. ", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("httpcompress %s still exists. ", name)
		}
		return nil
	}
}

func testCheckHttpcompresssDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_httpcompress" {
			continue
		}

		name := rs.Primary.ID
		httpcompress, err := client.GetHttpcompress(name)
		if err != nil {
			return err
		}
		if httpcompress != nil {
			return fmt.Errorf("httpcompress %s not destroyed. ", name)
		}
	}
	return nil
}

func getProfileHttpComprsConfig(profileName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_httpcompress" "test_httpcomprs_profile" {
  name                   = "%v"
  compression_buffersize = 4090
}
`, profileName)
}

func getProfileHttpComprsTC2Config(profileName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_httpcompress" "test_httpcomprs_profile" {
  name                   = "%v"
  gzip_compression_level = 2
  gzip_memory_level      = 32768
  gzip_window_size       = 32768
}
`, profileName)
}
