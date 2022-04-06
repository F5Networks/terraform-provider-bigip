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

var TEST_DEVICE_NAME = "test-device"

var TEST_DEVICE_RESOURCE = `
resource "bigip_cm_device" "test-device" {
            name = "` + TEST_DEVICE_NAME + `"
            configsync_ip = "2.2.2.2"
            mirror_ip = "10.10.10.10"
            mirror_secondary_ip = "11.11.11.11"
        }
`

func TestAccBigipCmDevice_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckdevicesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DEVICE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdeviceExists(TEST_DEVICE_NAME, true),
					resource.TestCheckResourceAttr("bigip_cm_device.test-device", "name", TEST_DEVICE_NAME),
					resource.TestCheckResourceAttr("bigip_cm_device.test-device", "configsync_ip", "2.2.2.2"),
					resource.TestCheckResourceAttr("bigip_cm_device.test-device", "mirror_ip", "10.10.10.10"),
					resource.TestCheckResourceAttr("bigip_cm_device.test-device", "mirror_secondary_ip", "11.11.11.11"),
				),
			},
		},
	})
}

func TestAccBigipCmDevice_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckdevicesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DEVICE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdeviceExists(TEST_DEVICE_NAME, true),
				),
				ResourceName:      TEST_DEVICE_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckdeviceExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		device, err := client.Devices(name)
		if err != nil {
			return err
		}
		if exists && device == nil {
			return fmt.Errorf("device %s was not created.", name)
		}
		if !exists && device != nil {
			return fmt.Errorf("device %s still exists.", name)
		}
		return nil
	}
}

func testCheckdevicesDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_cm_device" {
			continue
		}

		name := rs.Primary.ID
		device, err := client.Devices(name)
		if err != nil {
			return err
		}
		if device == nil {
			return fmt.Errorf("device %s not destroyed.", name)
		}
	}
	return nil
}
