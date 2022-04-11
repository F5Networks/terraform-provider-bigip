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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_HTTPCOMPRESS_NAME = fmt.Sprintf("/%s/test-httpcompress", TEST_PARTITION)

var TEST_HTTPCOMPRESS_RESOURCE = `
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
				Config: TEST_HTTPCOMPRESS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(TEST_HTTPCOMPRESS_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "name", "/Common/test-httpcompress"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress", "defaults_from", "/Common/httpcompression"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress",
						fmt.Sprintf("uri_exclude.%d", schema.HashString("f5.com")),
						"f5.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress",
						fmt.Sprintf("uri_include.%d", schema.HashString("cisco.com")),
						"cisco.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress",
						fmt.Sprintf("content_type_include.%d", schema.HashString("nicecontent.com")),
						"nicecontent.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_httpcompress.test-httpcompress",
						fmt.Sprintf("content_type_exclude.%d", schema.HashString("nicecontentexclude.com")),
						"nicecontentexclude.com"),
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
				Config: TEST_HTTPCOMPRESS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckHttpcompressExists(TEST_HTTPCOMPRESS_NAME, true),
				),
				ResourceName:      TEST_HTTPCOMPRESS_NAME,
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
			return fmt.Errorf("httpcompress %s was not created.", name)
		}
		if !exists && p == nil {
			return fmt.Errorf("httpcompress %s still exists.", name)
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
			return fmt.Errorf("httpcompress %s not destroyed.", name)
		}
	}
	return nil
}
