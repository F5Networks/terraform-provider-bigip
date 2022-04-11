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
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_FASTHTTP_NAME = fmt.Sprintf("/%s/test-fasthttp", TEST_PARTITION)

var TEST_FASTHTTP_RESOURCE = `
resource "bigip_ltm_profile_fasthttp" "test-fasthttp" {
	name = "` + TEST_FASTHTTP_NAME + `"
	defaults_from = "/Common/fasthttp"
            idle_timeout = 0
            connpoolidle_timeoutoverride	= 0
            connpool_maxreuse = 0
            connpool_maxsize  = 0
            connpool_minsize = 0
            connpool_replenish = "enabled"
            connpool_step = 0
            maxheader_size = 0
}
`

func TestAccBigipLtmfasthttp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfasthttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FASTHTTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckfasthttpProfileExists(TEST_FASTHTTP_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "name", TEST_FASTHTTP_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "defaults_from", "/Common/fasthttp"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "idle_timeout", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "connpoolidle_timeoutoverride", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "connpool_maxreuse", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "connpool_maxsize", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "connpool_minsize", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "connpool_replenish", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "connpool_step", "0"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fasthttp.test-fasthttp", "maxheader_size", "0"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfilefasthttp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfasthttpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FASTHTTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckfasthttpProfileExists(TEST_FASTHTTP_NAME, true),
				),
				ResourceName:      TEST_FASTHTTP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckfasthttpProfileExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.GetFasthttp(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("fasthttp %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("fasthttp %s still exists.", name)
		}
		return nil
	}
}

func testCheckfasthttpsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_fasthttp" {
			continue
		}

		name := rs.Primary.ID
		fasthttp, err := client.GetFasthttp(name)
		if err != nil {
			return err
		}
		if fasthttp == nil {
			return fmt.Errorf("fasthttp %s not destroyed.", name)
		}
	}
	return nil
}
