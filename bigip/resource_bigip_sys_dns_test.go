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

var TestDnsName = fmt.Sprintf("/%s/test-dns", TestPartition)

var TestDnsResource = `
resource "bigip_sys_dns" "test-dns" {
   description = "` + TestDnsName + `"
   name_servers = ["1.1.1.1"]
   search = ["f5.com"]
}
`

func TestAccBigipSysDNSCreateTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestDnsResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TestDnsName, true),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns", "description", TestDnsName),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns", "name_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("bigip_sys_dns.test-dns", "search.0", "f5.com"),
				),
			},
		},
	})
}

func TestAccBigipSysDNSCreateTC2(t *testing.T) {
	var TestDnsName = "test-dns-tc2"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getsysDNSConfigTC2(TestDnsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TestDnsName, true),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "description", TestDnsName),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "name_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "search.0", "f5.com"),
				),
			},
			{
				Config: getsysDNSConfigTC2(TestDnsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TestDnsName, true),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "description", TestDnsName),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "name_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "search.0", "f5.com"),
				),
			},
		},
	})
}

func TestAccBigipSysDNSCreateTC3(t *testing.T) {
	var TestDnsName = "test-dns-tc3"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getsysDNSConfigTC3(TestDnsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TestDnsName, true),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "description", TestDnsName),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "name_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "search.0", "f5.com"),
				),
			},
			{
				Config: getsysDNSConfigTC3a(TestDnsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TestDnsName, true),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "description", TestDnsName),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "name_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "name_servers.1", "2.2.2.2"),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "search.0", "f5.com"),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "search.1", "f5.net"),
				),
			},
		},
	})
}

func TestAccBigipSysDNSCreateTC4(t *testing.T) {
	var TestDnsName = "test-dns-tc4"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getsysDNSConfigTC4(TestDnsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TestDnsName, true),
					testCheckdnsExists("test-dns-tc4", false),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "description", TestDnsName),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "name_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "name_servers.1", "2.2.2.2"),
				),
			},
			{
				Config: getsysDNSConfigTC4(TestDnsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TestDnsName, true),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "description", TestDnsName),
					resource.TestCheckResourceAttr(fmt.Sprintf("bigip_sys_dns.%s", TestDnsName), "name_servers.0", "1.1.1.1"),
				),
			},
		},
	})
}
func TestAccBigipSysDNSImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestDnsResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TestDnsName, true),
				),
				ResourceName:      TestDnsName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckdnsExists(description string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		dns, err := client.DNSs()
		if err != nil {
			return err
		}
		if exists && dns == nil {
			return fmt.Errorf("dns %s was not created ", description)
		}
		if !exists && dns != nil {
			return fmt.Errorf("dns %s still exists ", description)
		}
		return nil
	}
}

func getsysDNSConfigTC2(sysDNSName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_dns" "test-dns-tc2" {
  description  = "%v"
  name_servers = ["1.1.1.1"]
  search       = ["f5.com"]
}
`, sysDNSName)
}

func getsysDNSConfigTC3(sysDNSName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_dns" "test-dns-tc3" {
  description  = "%v"
  name_servers = ["1.1.1.1"]
  search       = ["f5.com"]
}
`, sysDNSName)
}

func getsysDNSConfigTC3a(sysDNSName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_dns" "test-dns-tc3" {
  description  = "%v"
  name_servers = ["1.1.1.1", "2.2.2.2"]
  search       = ["f5.com", "f5.net"]
}
`, sysDNSName)
}

func getsysDNSConfigTC4(sysDNSName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_dns" "test-dns-tc4" {
  description  = "%v"
  name_servers = ["1.1.1.1", "2.2.2.2"]
}
`, sysDNSName)
}
