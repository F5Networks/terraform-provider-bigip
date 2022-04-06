/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_PROVISION_NAME = "afm"
var TEST_ASM_PROVISION_NAME = "asm"
var TEST_GTM_PROVISION_NAME = "gtm"
var TEST_APM_PROVISION_NAME = "apm"
var TEST_AVR_PROVISION_NAME = "avr"
var TEST_ILX_PROVISION_NAME = "ilx"

var TEST_PROVISION_RESOURCE = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TEST_PROVISION_NAME + `"
 full_path  = "afm"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TEST_ASM_PROVISION_RESOURCE = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TEST_ASM_PROVISION_NAME + `"
 full_path  = "asm"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TEST_GTM_PROVISION_RESOURCE = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TEST_GTM_PROVISION_NAME + `"
 full_path  = "gtm"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TEST_APM_PROVISION_RESOURCE = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TEST_APM_PROVISION_NAME + `"
 full_path  = "apm"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TEST_AVR_PROVISION_RESOURCE = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TEST_AVR_PROVISION_NAME + `"
 full_path  = "avr"
 cpu_ratio = 0
 disk_ratio = 0
 level = "none"
 memory_ratio = 0
}
`
var TEST_ILX_PROVISION_RESOURCE = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TEST_ILX_PROVISION_NAME + `"
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
				Config: TEST_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_PROVISION_NAME),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TEST_PROVISION_NAME),
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
				Config: TEST_ASM_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_ASM_PROVISION_NAME),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TEST_ASM_PROVISION_NAME),
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
				Config: TEST_GTM_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_GTM_PROVISION_NAME),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TEST_GTM_PROVISION_NAME),
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
				Config: TEST_APM_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_APM_PROVISION_NAME),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TEST_APM_PROVISION_NAME),
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
				Config: TEST_AVR_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_AVR_PROVISION_NAME),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TEST_AVR_PROVISION_NAME),
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
				Config: TEST_ILX_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_ILX_PROVISION_NAME),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TEST_ILX_PROVISION_NAME),
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
				Config: TEST_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_PROVISION_NAME),
				),
				ResourceName:      TEST_PROVISION_NAME,
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
				Config: TEST_ASM_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_ASM_PROVISION_NAME),
				),
				ResourceName:      TEST_ASM_PROVISION_NAME,
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
				Config: TEST_GTM_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_GTM_PROVISION_NAME),
				),
				ResourceName:      TEST_GTM_PROVISION_NAME,
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
				Config: TEST_APM_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_APM_PROVISION_NAME),
				),
				ResourceName:      TEST_APM_PROVISION_NAME,
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
				Config: TEST_AVR_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_AVR_PROVISION_NAME),
				),
				ResourceName:      TEST_AVR_PROVISION_NAME,
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
				Config: TEST_ILX_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_ILX_PROVISION_NAME),
				),
				ResourceName:      TEST_ILX_PROVISION_NAME,
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
