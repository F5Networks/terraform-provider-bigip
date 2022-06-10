/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"regexp"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_SELFIP_NAME = fmt.Sprintf("/%s/test-selfip", TEST_PARTITION)
var TEST_FLOAT_SELFIP_NAME = fmt.Sprintf("/%s/test-float-selfip", TEST_PARTITION)

var TEST_SELFIP_RESOURCE = `
resource "bigip_net_vlan" "test-vlan" {
  name = "` + TEST_VLAN_NAME + `"
  tag = 101
  interfaces {
    vlanport = 1.1
    tagged = true
  }
}
resource "bigip_net_selfip" "test-selfip" {
  name = "` + TEST_SELFIP_NAME + `"
  ip = "11.1.1.1/24"
  vlan = "/Common/test-vlan"
  depends_on = ["bigip_net_vlan.test-vlan"]
}
resource "bigip_net_selfip" "test-float-selfip" {
  name = "` + TEST_FLOAT_SELFIP_NAME + `"
  ip = "11.1.1.2/24"
  traffic_group = "traffic-group-1"
  vlan = "/Common/test-vlan"
  depends_on = ["bigip_net_selfip.test-selfip"]
}
`

func TestAccBigipNetselfip_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckselfipsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SELFIP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME),
					testCheckselfipExists(TEST_FLOAT_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "name", TEST_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "ip", "11.1.1.1/24"),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "vlan", TEST_VLAN_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-float-selfip", "name", TEST_FLOAT_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-float-selfip", "ip", "11.1.1.2/24"),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-float-selfip", "vlan", TEST_VLAN_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-float-selfip", "traffic_group", "traffic-group-1"),
				),
			},
		},
	})
}

func TestAccBigipNetselfip_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckselfipsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SELFIP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME),
				),
				ResourceName:      TEST_SELFIP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckselfipsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_SELFIP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_FLOAT_SELFIP_NAME),
				),
				ResourceName:      TEST_FLOAT_SELFIP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigipNetselfipPortlockdown(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckselfipsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccselfipPortLockdownParam("all"),
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "port_lockdown.0", "all"),
				),
			},
			{
				Config: testaccselfipPortLockdownParam("protocol"),
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "port_lockdown.0", "egp:0"),
				),
			},
			{
				Config: testaccselfipPortLockdownParam("tcp"),
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "port_lockdown.0", "tcp:4040"),
				),
			},
			{
				Config: testaccselfipPortLockdownParam("udp"),
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "port_lockdown.0", "udp:4040"),
				),
			},
			{
				Config: testaccselfipPortLockdownParam("default"),
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "port_lockdown.0", "default"),
				),
			},
			{
				Config: testaccselfipPortLockdownParam("custom_default"),
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "port_lockdown.0", "default"),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "port_lockdown.1", "tcp:4040"),
				),
			},
			{
				Config: testaccselfipPortLockdownParam("none"),
				Check: resource.ComposeTestCheckFunc(
					testCheckselfipExists(TEST_SELFIP_NAME),
					resource.TestCheckResourceAttr("bigip_net_selfip.test-selfip", "port_lockdown.0", "none"),
				),
			},
		},
	})
}

func TestAccBigipNetselfipRouteDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckselfipsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccselfipRouteDomain("10.11.12.13%0/24"),
				Check:  testCheckselfipExists(TEST_SELFIP_NAME),
			},
			{
				Config:             testaccselfipRouteDomain("10.11.12.13/24"),
				Check:              testCheckselfipExists(TEST_SELFIP_NAME),
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile("Expected a non-empty plan, but got an empty plan!"),
			},
			{
				Config: testaccselfipRouteDomain("10.11.12.13%0/24"),
				Check:  testCheckselfipExists(TEST_SELFIP_NAME),
			},
			{
				Config:             testaccselfipRouteDomain("10.11.12.13%0/24"),
				Check:              testCheckselfipExists(TEST_SELFIP_NAME),
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile("Expected a non-empty plan, but got an empty plan!"),
			},
		},
	})
}

func testaccselfipRouteDomain(ip string) string {
	resPrefix := `
	resource "bigip_net_vlan" "test-vlan" {
      name = "` + TEST_VLAN_NAME + `"
	  tag = 101
	  interfaces {
		vlanport = 1.1
		tagged = true
	  }
	}
	resource "bigip_net_selfip" "test-selfip" {
	  name = "/Common/test-selfip"
	  ip   = "%s"
	  vlan = "/Common/test-vlan"
	  depends_on = ["bigip_net_vlan.test-vlan"]
	}
	`
	return fmt.Sprintf(resPrefix, ip)
}

func testaccselfipPortLockdownParam(portLockdown string) string {
	resPrefix := `
	resource "bigip_net_vlan" "test-vlan" {
      name = "` + TEST_VLAN_NAME + `"
	  tag = 101
	  interfaces {
		vlanport = 1.1
		tagged = true
	  }
	}
	resource "bigip_net_selfip" "test-selfip" {
	  name = "/Common/test-selfip"
		ip   = "11.1.1.1/24"
		vlan = "/Common/test-vlan"
		depends_on = ["bigip_net_vlan.test-vlan"]
	`
	switch portLockdown {
	case "all":
		resPrefix = fmt.Sprintf(
			`%s
			port_lockdown = ["all"]
			`,
			resPrefix,
		)
	case "protocol":
		resPrefix = fmt.Sprintf(
			`%s
			port_lockdown = ["egp:0"]
			`,
			resPrefix,
		)
	case "tcp":
		resPrefix = fmt.Sprintf(
			`%s
			port_lockdown = ["tcp:4040"]
			`,
			resPrefix,
		)
	case "udp":
		resPrefix = fmt.Sprintf(
			`%s
			port_lockdown = ["udp:4040"]
			`,
			resPrefix,
		)
	case "default":
		resPrefix = fmt.Sprintf(
			`%s
			port_lockdown = ["default"]
			`,
			resPrefix,
		)
	case "custom_default":
		resPrefix = fmt.Sprintf(
			`%s
			port_lockdown = ["default", "tcp:4040"]
			`,
			resPrefix,
		)
	case "none":
		resPrefix = fmt.Sprintf(
			`%s
			port_lockdown = ["none"]
			`,
			resPrefix,
		)
	default:
	}
	return fmt.Sprintf(`%s
		}`, resPrefix)
}

func testCheckselfipExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.SelfIP(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("selfip %s was not created.", name)
		}

		return nil
	}
}

func testCheckselfipsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_net_selfip" {
			continue
		}

		name := rs.Primary.ID
		selfip, err := client.SelfIP(name)
		if err != nil {
			return err
		}
		if selfip == nil {
			return fmt.Errorf("selfip %s not destroyed.", name)
		}
	}
	return nil
}
