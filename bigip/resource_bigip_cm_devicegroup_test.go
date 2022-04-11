/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_DG_NAME = "test-devicegroup"

var TEST_DG_RESOURCE = `

resource "bigip_cm_device" "test-device" {
            name = "` + TEST_DEVICE_NAME + `"
            configsync_ip = "2.2.2.2"
            mirror_ip = "10.10.10.10"
            mirror_secondary_ip = "11.11.11.11"
        }
resource "bigip_cm_devicegroup" "test-devicegroup" {
            name = "` + TEST_DG_NAME + `"
	    partition = "Common"
	    description = "whatiknow"
            auto_sync = "disabled"
            full_load_on_sync = "false"
            type = "sync-only"
	    save_on_auto_sync = "false"
	    network_failover = "enabled"
	    incremental_config = 1024
	    device {
               name = "/Common/test-device"
             }
	    depends_on = ["bigip_cm_device.test-device"]
        }
`

func TestAccBigipCmDevicegroup_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckCmDevicegroupsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DG_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckCmDevicegroupExists(TEST_DG_NAME, true),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "name", TEST_DG_NAME),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "description", "whatiknow"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "auto_sync", "disabled"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "full_load_on_sync", "false"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "type", "sync-only"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "save_on_auto_sync", "false"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "network_failover", "enabled"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "incremental_config", "1024"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "device.0.name", "/Common/test-device"),
				),
			},
		},
	})
}

func TestAccBigipCmDevicegroup_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckCmDevicegroupsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DG_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckCmDevicegroupExists(TEST_DG_NAME, true),
				),
				ResourceName:      TEST_DG_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})

}

func testCheckCmDevicegroupExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		devicegroup, err := client.Devicegroups(name)
		if err != nil {
			return err
		}

		if exists && devicegroup == nil {
			return fmt.Errorf("devicegroup %s was not created.", name)
		}
		if !exists && devicegroup != nil {
			return fmt.Errorf("devicegroup %s still exists.", name)
		}

		return nil

	}
}

func testCheckCmDevicegroupsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_cm_devicegroup" {
			continue
		}
		name := rs.Primary.ID
		devicegroup, err := client.Devicegroups(name)
		log.Println("the state file is =================  ", rs.Type)
		if err != nil {
			return err
		}
		if devicegroup == nil {
			return fmt.Errorf("devicegroup %s not destroyed.", name)
		}

	}
	return nil
}
