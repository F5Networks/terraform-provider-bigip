/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestTcpName = fmt.Sprintf("/%s/test-tcp", TestPartition)

var TestTcpResource = `
resource "bigip_ltm_profile_tcp" "test-tcp" {
  name               = "/Common/sanjose-tcp-wan-profile"
  defaults_from      = "/Common/tcp-wan-optimized"
  idle_timeout       = 300
  close_wait_timeout = 5
  finwait_2timeout   = 5
  finwait_timeout    = 300
  keepalive_interval = 1700
  deferred_accept    = "enabled"
  fast_open          = "enabled"
}
`

func TestAccBigipLtmProfileTcp_create(t *testing.T) {
	TestTcpName = fmt.Sprintf("/%s/%s", "Common", "sanjose-tcp-wan-profile")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestTcpResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(TestTcpName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "name", "/Common/sanjose-tcp-wan-profile"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "defaults_from", "/Common/tcp-wan-optimized"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "idle_timeout", "300"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "close_wait_timeout", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "finwait_2timeout", "5"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "finwait_timeout", "300"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "keepalive_interval", "1700"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "deferred_accept", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test-tcp", "fast_open", "enabled"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileTcpTC1(t *testing.T) {
	profileTcpName := fmt.Sprintf("/%s/%s", "Common", "test_tcp_profiletc1")
	httpsTenantName = "fast_https_tenanttc1"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getProfileTCPConfig(profileTcpName),
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(profileTcpName, true),
					testCheckTcpExists(TestTcpName, false),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "name", profileTcpName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "initial_congestion_windowsize", "20"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "congestion_control", "cdg"),
				),
			},
			{
				Config: getProfileTCPConfig(profileTcpName),
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(profileTcpName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "name", profileTcpName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "initial_congestion_windowsize", "20"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "congestion_control", "cdg"),
				),
			},
			{
				Config: getProfileTCPConfig2(profileTcpName),
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(profileTcpName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "name", profileTcpName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "initial_congestion_windowsize", "30"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "congestion_control", "bbr"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileTcpTC2(t *testing.T) {
	profileTcpName := fmt.Sprintf("/%s/%s", "Common", "test_tcp_profiletc2")
	httpsTenantName = "fast_https_tenanttc1"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getProfileTCPConfigDefault(profileTcpName),
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(profileTcpName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "name", profileTcpName),
				),
			},
			{
				Config: getProfileTCPConfigTC2(profileTcpName),
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(profileTcpName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "name", profileTcpName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "delayed_acks", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "nagle", "auto"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "early_retransmit", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "tailloss_probe", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "timewait_recycle", "disabled"),
				),
			},
			{
				Config: getProfileTCPConfigTC2(profileTcpName),
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(profileTcpName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "name", profileTcpName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "delayed_acks", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "nagle", "auto"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "early_retransmit", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "tailloss_probe", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "timewait_recycle", "disabled"),
				),
			},
			{
				Config: getProfileTCPConfigTC2Modify(profileTcpName),
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(profileTcpName, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "name", profileTcpName),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "delayed_acks", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "nagle", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "early_retransmit", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "tailloss_probe", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_tcp.test_tcp_profile", "timewait_recycle", "enabled"),
				),
			},
		},
	})
}

func TestAccBigipLtmProfileTcp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestTcpResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(TestTcpName, true),
				),
				ResourceName:      TestTcpName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckTcpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetTcp(name)
		if err != nil {
			return err
		}
		log.Printf("P:%+v", p)
		if exists && p == nil {
			return fmt.Errorf("tcp %s was not created ", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("tcp %s still exists ", name)
		}
		return nil
	}
}

func testCheckTcpsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_tcp" {
			continue
		}

		name := rs.Primary.ID
		tcp, err := client.GetTcp(name)
		if err != nil {
			return err
		}
		if tcp != nil {
			return fmt.Errorf("tcp %s not destroyed ", name)
		}
	}
	return nil
}

func getProfileTCPConfigDefault(profileName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_tcp" "test_tcp_profile" {
  name = "%v"
}
`, profileName)
}

func getProfileTCPConfig(profileName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_tcp" "test_tcp_profile" {
  name                          = "%v"
  congestion_control            = "cdg"
  initial_congestion_windowsize = 20
}
`, profileName)
}

func getProfileTCPConfig2(profileName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_tcp" "test_tcp_profile" {
  name                          = "%v"
  congestion_control            = "bbr"
  initial_congestion_windowsize = 30
}
`, profileName)
}

func getProfileTCPConfigTC2(profileName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_tcp" "test_tcp_profile" {
  name             = "%v"
  delayed_acks     = "disabled"
  nagle            = "auto"
  early_retransmit = "disabled"
  tailloss_probe   = "disabled"
  timewait_recycle = "disabled"
}
`, profileName)
}
func getProfileTCPConfigTC2Modify(profileName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_tcp" "test_tcp_profile" {
  name             = "%v"
  delayed_acks     = "enabled"
  nagle            = "disabled"
  early_retransmit = "enabled"
  tailloss_probe   = "enabled"
  timewait_recycle = "enabled"
}
`, profileName)
}
