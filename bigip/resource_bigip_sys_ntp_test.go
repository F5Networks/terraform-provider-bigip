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

var TEST_NTP_NAME = fmt.Sprintf("/%s/test-ntp", TEST_PARTITION)

var TEST_NTP_RESOURCE = `
resource "bigip_sys_ntp" "test-ntp" {
	description = "` + TEST_NTP_NAME + `"
	servers = ["10.10.10.10"]
	timezone = "America/Los_Angeles"
}
`

func TestAccBigipSysNtp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TEST_NTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckntpExists(TEST_NTP_NAME, true),
					resource.TestCheckResourceAttr("bigip_sys_ntp.test-ntp", "description", TEST_NTP_NAME),
					resource.TestCheckResourceAttr("bigip_sys_ntp.test-ntp", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("bigip_sys_ntp.test-ntp",
						fmt.Sprintf("servers.%d", schema.HashString("10.10.10.10")),
						"10.10.10.10"),
				),
			},
		},
	})
}

func TestAccBigipSysNtp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//	CheckDestroy: testCheckntpsDestroyed, ( No Delet API support)
		Steps: []resource.TestStep{
			{
				Config: TEST_NTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckntpExists(TEST_NTP_NAME, true),
				),
				ResourceName:      TEST_NTP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckntpExists(description string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		ntp, err := client.NTPs()
		if err != nil {
			return err
		}
		if exists && ntp == nil {
			return fmt.Errorf("ntp %s was not created.", description)

		}
		if !exists && ntp != nil {
			return fmt.Errorf("ntp %s still exists.", description)

		}
		return nil
	}
}
