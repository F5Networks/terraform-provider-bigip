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

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TEST_DATAGROUP_NAME = "/" + TEST_PARTITION + "/test-datagroup"

var TEST_DATAGROUP_STRING_RESOURCE = `
        resource "bigip_ltm_datagroup" "test-datagroup-string" {
                name = "` + TEST_DATAGROUP_NAME + `"
                type = "string"

                record {
                        name = "test-name1"
                        data = "test-data1"
                }
                record {
                        name = "test-name2"
                }
        }`

var TEST_DATAGROUP_IP_RESOURCE = `
	resource "bigip_ltm_datagroup" "test-datagroup-ip" {
		name = "` + TEST_DATAGROUP_NAME + `"
		type = "ip"

		record {
			name = "3.3.3.3/32"
			data = "1.1.1.1"
		}
		record {
			name = "2.2.2.2/32"
		}
	}`

var TEST_DATAGROUP_INTEGER_RESOURCE = `
        resource "bigip_ltm_datagroup" "test-datagroup-integer" {
                name = "` + TEST_DATAGROUP_NAME + `"
                type = "integer"

                record {
                        name = "1"
                        data = "2"
                }
                record {
                        name = "3"
                }
        }`

func TestAccBigipLtmDataGroup_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DATAGROUP_STRING_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TEST_DATAGROUP_NAME),
				),
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
				Config: TEST_DATAGROUP_IP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TEST_DATAGROUP_NAME),
				),
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
				Config: TEST_DATAGROUP_INTEGER_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TEST_DATAGROUP_NAME),
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
				Config: TEST_DATAGROUP_STRING_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TEST_DATAGROUP_NAME),
				),
				ResourceName:      TEST_DATAGROUP_NAME,
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
				Config: TEST_DATAGROUP_IP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TEST_DATAGROUP_NAME),
				),
				ResourceName:      TEST_DATAGROUP_NAME,
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
				Config: TEST_DATAGROUP_INTEGER_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TEST_DATAGROUP_NAME),
				),
				ResourceName:      TEST_DATAGROUP_NAME,
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
			return fmt.Errorf("Error while fetching Data Group: %v", err)

		}

		datagroup_name := fmt.Sprintf("/%s/%s", datagroup.Partition, datagroup.Name)
		if datagroup_name != name {
			return fmt.Errorf("Data Group name does not match. Expecting %s got %s.", name, datagroup_name)
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
