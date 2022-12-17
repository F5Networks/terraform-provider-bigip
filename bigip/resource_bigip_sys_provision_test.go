/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

var TestProvisionName = "afm"
var TestAsmProvisionName = "asm"
var TestGtmProvisionName = "gtm"
var TestApmProvisionName = "apm"
var TestAvrProvisionName = "avr"
var TestIlxProvisionName = "ilx"

var TestProvisionResource = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TestProvisionName + `"
 full_path  = "afm"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TestAsmProvisionResource = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TestAsmProvisionName + `"
 full_path  = "asm"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TestGtmProvisionResource = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TestGtmProvisionName + `"
 full_path  = "gtm"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TestApmProvisionResource = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TestApmProvisionName + `"
 full_path  = "apm"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TestAvrProvisionResource = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TestAvrProvisionName + `"
 full_path  = "avr"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TestIlxProvisionResource = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TestIlxProvisionName + `"
 full_path  = "ilx"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`

func TestAccBigipSysProvision_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TestProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "full_path", "afm"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "cpu_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "disk_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "level", "none"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "memory_ratio", "0"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAsmProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestAsmProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TestAsmProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "full_path", "asm"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "cpu_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "disk_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "level", "none"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "memory_ratio", "0"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestGtmProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestGtmProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TestGtmProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "full_path", "gtm"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "cpu_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "disk_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "level", "none"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "memory_ratio", "0"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestApmProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestApmProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TestApmProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "full_path", "apm"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "cpu_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "disk_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "level", "none"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "memory_ratio", "0"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAvrProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestAvrProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TestAvrProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "full_path", "avr"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "cpu_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "disk_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "level", "none"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "memory_ratio", "0"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestIlxProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestIlxProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TestIlxProvisionName),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "full_path", "ilx"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "cpu_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "disk_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "level", "none"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "memory_ratio", "0"),
				),
			},
		},
	})
}

func TestAccBigipSysProvision_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestProvisionName),
				),
				ResourceName:      TestProvisionName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAsmProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestAsmProvisionName),
				),
				ResourceName:      TestAsmProvisionName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestGtmProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestGtmProvisionName),
				),
				ResourceName:      TestGtmProvisionName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestApmProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestApmProvisionName),
				),
				ResourceName:      TestApmProvisionName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAvrProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestAvrProvisionName),
				),
				ResourceName:      TestAvrProvisionName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestIlxProvisionResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TestIlxProvisionName),
				),
				ResourceName:      TestIlxProvisionName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckprovisionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		provision, err := client.Provisions(name)
		if err != nil {
			return err
		}
		if provision == nil {
			return fmt.Errorf("provision %s was not created.", name)

		}

		return nil
	}
}
