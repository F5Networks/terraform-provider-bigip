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

var TEST_IAPP_NAME = "/" + TEST_PARTITION + "/test-iapp"

var TEST_IAPP_RESOURCE = `
	resource "bigip_sys_iapp" "test-iapp" {
		name = "test-iapp"
		jsonfile = <<EOF
		{
			"inheritedDevicegroup": "true",
			"inheritedTrafficGroup": "true",
			"name": "test-iapp",
			"partition": "Common",
			"strictUpdates": "enabled",
			"tables": [
					{
							"name": "basic__snatpool_members"
					},
					{
							"name": "net__snatpool_members"
					},
					{
							"name": "optimizations__hosts"
					},
					{
							"columnNames": [
									"name"
							],
							"name": "pool__hosts",
							"rows": [
									{
											"row": [
													"f5.cisco.com"
											]
									}
							]
					},
					{
							"columnNames": [
									"addr",
									"port",
									"connection_limit"
							],
							"name": "pool__members",
							"rows": [
									{
											"row": [
													"10.0.2.167",
													"80",
													"0"
											]
									},
									{
											"row": [
													"10.0.2.168",
													"80",
													"0"
											]
									}
							]
					},
					{
							"name": "server_pools__servers"
					}
			],
			"template": "/Common/f5.http",
			"templateModified": "no",
			"templateReference": {
					"link": "https://localhost/mgmt/tm/sys/application/template/~Common~f5.http?ver=13.0.0"
			},
		 
			"variables": [
					{
							"encrypted": "no",
							"name": "client__http_compression",
							"value": "/#create_new#"
					},
					{
							"encrypted": "no",
							"name": "monitor__monitor",
							"value": "/Common/http"
					},
					{
							"encrypted": "no",
							"name": "net__client_mode",
							"value": "wan"
					},
					{
							"encrypted": "no",
							"name": "net__server_mode",
							"value": "lan"
					},
					{
							"encrypted": "no",
							"name": "net__v13_tcp",
							"value": "warn"
					},
					{
							"encrypted": "no",
							"name": "pool__addr",
							"value": "10.0.1.100"
					},
					{
							"encrypted": "no",
							"name": "pool__pool_to_use",
							"value": "/#create_new#"
					},
					{
							"encrypted": "no",
							"name": "pool__port",
							"value": "80"
					},
					{
							"encrypted": "no",
							"name": "ssl__mode",
							"value": "no_ssl"
					},
					{
							"encrypted": "no",
							"name": "ssl_encryption_questions__advanced",
							"value": "no"
					},
					{
							"encrypted": "no",
							"name": "ssl_encryption_questions__help",
							"value": "hide"
					}
			]
	}
	 

EOF
	}`

func TestAccBigipSysIapp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckIappDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_IAPP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckIappExists(TEST_IAPP_NAME),
				),
			},
		},
	})
}

func TestAccBigipSysIapp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckIappDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_IAPP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckIappExists(TEST_IAPP_NAME),
				),
				ResourceName:      TEST_IAPP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckIappExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		jsonfile, err := client.Iapp(name)
		log.Println(" I am here in Exists !!!!!!!!!!!!", name)
		if err != nil {
			return fmt.Errorf("Error while fetching iapp: %v", err)

		}
		body := s.RootModule().Resources["bigip_sys_iapp.test-iapp"].Primary.Attributes["name"]
		if jsonfile.Name == body {
			return fmt.Errorf("jsonfile  body does not match. Expecting %s got %s.", body, jsonfile.Name)
		}

		jsonfile_name := fmt.Sprintf("/%s/%s", jsonfile.Partition, jsonfile.Name)
		if jsonfile_name == name {
			return fmt.Errorf("Jsonfile name does not match. Expecting %s got %s.", name, jsonfile_name)
		}
		return nil
	}
}

func testCheckIappDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_sys_iapp" {
			continue
		}

		name := rs.Primary.ID
		log.Println(" I am in Destroy function currently +++++++++++++++++++++++++++ ", name)

		// Join three strings into one.
		jsonfile, err := client.Iapp(name)

		if err != nil {
			return nil
		}

		if jsonfile == nil {
			return fmt.Errorf("Iapp %s not destroyed.", name)
		}
	}
	return nil
}
