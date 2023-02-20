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

var TestFastl4Name = fmt.Sprintf("/%s/test-fastl4", TestPartition)

var TestFastl4Resource = `
resource "bigip_ltm_profile_fastl4" "test-fastl4" {
            name = "` + TestFastl4Name + `"
            partition = "Common"
            defaults_from = "/Common/fastL4"
			client_timeout = 40
			idle_timeout = "200"
            explicitflow_migration = "enabled"
			late_binding = "enabled"
            hardware_syncookie = "enabled"
            iptos_toclient = "pass-through"
            iptos_toserver = "pass-through"
            keepalive_interval = "disabled"
 }
`

func TestAccBigipLtmProfileFastl4TC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfastl4sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestFastl4Resource,
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(TestFastl4Name, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "name", TestFastl4Name),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "defaults_from", "/Common/fastL4"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "client_timeout", "40"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "explicitflow_migration", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "late_binding", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "hardware_syncookie", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "idle_timeout", "200"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "hardware_syncookie", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "iptos_toclient", "pass-through"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "iptos_toserver", "pass-through"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test-fastl4", "keepalive_interval", "disabled"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileFastl4TC2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfastl4sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmprofilefastl4DefaultConfig(TestPartition, TestFastl4Name, "fastl4profileParent"),
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(TestFastl4Name, true),
					testCheckfastl4Exists("/Common/tetsfastl44", false),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.fastl4profileParent", "name", TestFastl4Name),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.fastl4profileParent", "defaults_from", "/Common/fastL4"),
				),
			},
			{
				Config: testaccbigipltmprofilefastl4UpdateIdletimeoutConfig(TestPartition, TestFastl4Name, "fastl4profileParent"),
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(TestFastl4Name, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.fastl4profileParent", "name", TestFastl4Name),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.fastl4profileParent", "defaults_from", "/Common/fastL4"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.fastl4profileParent", "idle_timeout", "307"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileFastl4TC3(t *testing.T) {
	profileFastL4Name := fmt.Sprintf("/%s/%s", "Common", "test_fastl4_profiletc3")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getProfileFastl4ConfigTC3(profileFastL4Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(profileFastL4Name, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc3", "name", profileFastL4Name),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc3", "defaults_from", "/Common/fastL4"),
				),
			},
			{
				Config: getProfileFastl4ConfigTC3(profileFastL4Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(profileFastL4Name, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc3", "name", profileFastL4Name),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc3", "defaults_from", "/Common/fastL4"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileFastl4TC4(t *testing.T) {
	profileFastL4Name := fmt.Sprintf("/%s/%s", "Common", "test_fastl4_profiletc4")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getProfileFastl4ConfigTC4(profileFastL4Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(profileFastL4Name, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "name", profileFastL4Name),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "defaults_from", "/Common/fastL4"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "idle_timeout", "200"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "late_binding", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "loose_close", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "loose_initiation", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "explicitflow_migration", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "keepalive_interval", "150"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "tcp_handshake_timeout", "100"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "receive_windowsize", "100"),
				),
			},
			{
				Config: getProfileFastl4ConfigTC4(profileFastL4Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(profileFastL4Name, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "name", profileFastL4Name),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "defaults_from", "/Common/fastL4"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "idle_timeout", "200"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "late_binding", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "loose_close", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "loose_initiation", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "explicitflow_migration", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "keepalive_interval", "150"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "tcp_handshake_timeout", "100"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_fastl4.test_fastl4_profile_tc4", "receive_windowsize", "100"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileFastl4_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfastl4sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestFastl4Resource,
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(TestFastl4Name, true),
				),
				ResourceName:      TestFastl4Name,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckfastl4Exists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetFastl4(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("fastl4 %s was not created. ", name)
		}
		if !exists && p == nil {
			return fmt.Errorf("fastl4 %s was still exist. ", name)
		}
		return nil
	}
}

func testCheckfastl4sDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_fastl4" {
			continue
		}

		name := rs.Primary.ID
		fastl4, err := client.GetFastl4(name)
		if err != nil {
			return err
		}
		if fastl4 != nil {
			return fmt.Errorf("fastl4 %s not destroyed. ", name)
		}
	}
	return nil
}

func testaccbigipltmprofilefastl4DefaultConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_fastl4" "%[3]s" {
  name          = "%[2]s"
  defaults_from = "/%[1]s/fastL4"
}
`, partition, profileName, resourceName)
}

func testaccbigipltmprofilefastl4UpdateIdletimeoutConfig(partition, profileName, resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_fastl4" "%[3]s" {
  name          = "%[2]s"
  defaults_from = "/%[1]s/fastL4"
  idle_timeout  = "307"
}
`, partition, profileName, resourceName)
}

func getProfileFastl4ConfigTC3(profileName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_fastl4" "test_fastl4_profile_tc3" {
  name          = "%v"
  defaults_from = "/Common/fastL4"
}
`, profileName)
}

func getProfileFastl4ConfigTC4(profileName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_fastl4" "test_fastl4_profile_tc4" {
  name                   = "%v"
  defaults_from          = "/Common/fastL4"
  idle_timeout           = "200"
  explicitflow_migration = "disabled"
  late_binding           = "enabled"
  loose_close            = "enabled"
  loose_initiation       = "enabled"
  keepalive_interval     = "150"
  tcp_handshake_timeout  = "100"
  receive_windowsize     = 100
}
`, profileName)
}
