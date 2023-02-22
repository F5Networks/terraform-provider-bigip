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

var TestDatagroupName = "/" + TestPartition + "/test-datagroup"

var TestDatagroupStringResource = `
        resource "bigip_ltm_datagroup" "test-datagroup-string" {
                name = "` + TestDatagroupName + `"
                type = "string"
                record {
                        name = "test-name1"
                        data = "test-data1"
                }
                record {
                        name = "test-name2"
                }
        }`

var TestDatagroupIpResource = `
	resource "bigip_ltm_datagroup" "test-datagroup-ip" {
		name = "` + TestDatagroupName + `"
		type = "ip"
		record {
			name = "3.3.3.3/32"
			data = "1.1.1.1"
		}
		record {
			name = "2.2.2.2/32"
		}
	}`

var TestDatagroupIntegerResource = `
        resource "bigip_ltm_datagroup" "test-datagroup-integer" {
                name = "` + TestDatagroupName + `"
                type = "integer"
                record {
                        name = "1"
                        data = "2"
                }
                record {
                        name = "3"
                }
        }`

func TestAccBigipLtmDataGroup_Create_TypeString(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestDatagroupStringResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TestDatagroupName),
				),
			},
		},
	})
}
func TestAccBigipLtmDataGroup_Create_TypeIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestDatagroupIpResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TestDatagroupName),
				),
			},
		},
	})
}

func TestAccBigipLtmDataGroup_Create_TypeInteger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestDatagroupIntegerResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TestDatagroupName),
				),
			},
		},
	})
}

func TestAccBigipLtmDataGroup_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestDatagroupStringResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TestDatagroupName),
				),
				ResourceName:      TestDatagroupName,
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
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestDatagroupIpResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TestDatagroupName),
				),
				ResourceName:      TestDatagroupName,
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
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestDatagroupIntegerResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TestDatagroupName),
				),
				ResourceName:      TestDatagroupName,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckDataGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		datagroup, err := client.GetInternalDataGroup(name)
		if err != nil {
			return fmt.Errorf("Error while fetching Data Group: %v ", err)

		}

		datagroupName := fmt.Sprintf("/%s/%s", datagroup.Partition, datagroup.Name)
		if datagroupName != name {
			return fmt.Errorf("Data Group name does not match. Expecting %s got %s ", name, datagroupName)
		}
		return nil
	}
}

func testCheckDataGroupDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_datagroup" {
			continue
		}

		name := rs.Primary.ID
		datagroup, err := client.GetInternalDataGroup(name)

		if err != nil {
			return nil
		}
		if datagroup != nil {
			return fmt.Errorf("Data Group %s not destroyed.", name)
		}
	}
	return nil
}
