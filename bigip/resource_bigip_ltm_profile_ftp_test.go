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

var TEST_FTP_NAME = fmt.Sprintf("/%s/test-ftp", TEST_PARTITION)

var TEST_FTP_RESOURCE = `
resource "bigip_ltm_profile_ftp" "test-ftp" {
            name = "/Common/sanjose-ftp-profile"
            defaults_from = "/Common/ftp"
            port  = 2020
            partition = "Common"
            description = "test-tftp-profile"
            security = "disabled"
            translate_extended = "enabled"
        }
`

func TestAccBigipLtmProfileFtp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFtpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckFtpExists(TEST_FTP_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_ftp.test-ftp", "name", "/Common/sanjose-ftp-profile"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_ftp.test-ftp", "defaults_from", "/Common/ftp"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_ftp.test-ftp", "port", "2020"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_ftp.test-ftp", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_ftp.test-ftp", "description", "test-tftp-profile"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_ftp.test-ftp", "security", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_ftp.test-ftp", "translate_extended", "enabled"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileFtp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFtpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckFtpExists(TEST_FTP_NAME, true),
				),
				ResourceName:      TEST_FTP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckFtpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetFtp(name)
		if err != nil {
			return err
		}
		if exists && p != nil {
			return fmt.Errorf("ftp %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("ftp %s still exists.", name)
		}
		return nil
	}
}

func testCheckFtpsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_ftp" {
			continue
		}

		name := rs.Primary.ID
		ftp, err := client.GetFtp(name)
		if err != nil {
			return err
		}
		if ftp != nil {
			return fmt.Errorf("ftp %s not destroyed.", name)
		}
	}
	return nil
}
