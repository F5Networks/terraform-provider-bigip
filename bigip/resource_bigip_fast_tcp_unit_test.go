/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipFastTCPUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-http"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipFastTCPInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipFastTCPUnitCreate(t *testing.T) {
	resourceName := "fasttcpapp"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"code": 202,
			"requestId": 317,
			"message": [
				{
					"id": "create_id",
					"name": "bigip-fast-templates/tcp",
					"parameters": {
						"app_name": "%s",
						"enable_monitor": true,
						"enable_pool": true,
						"enable_snat": true,
						"make_monitor": true,
						"make_pool": true,
						"make_snatpool": false,
						"monitor_interval": 40,
						"pool_members": [
							{
								"connectionLimit": 10,
								"priorityGroup": 2,
								"serverAddresses": [
									"1.2.3.4",
									"4.5.4.7"
								],
								"servicePort": 80,
								"shareNodes": true
							}
						],
						"snat_automap": true,
						"tenant_name": "fasttcptenant",
						"virtual_address": "10.1.10.101",
						"virtual_port": 80
					}
				}
			]
		}
		`, resourceName)
	})

	mux.HandleFunc("/mgmt/shared/fast/tasks/create_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"id": "create_id",
			"code": 200,
			"message": "success",
			"name": "bigip-fast-templates/tcp",
			"parameters": {
				"app_name": "%s",
				"enable_monitor": true,
				"enable_pool": true,
				"enable_snat": true,
				"make_monitor": true,
				"make_pool": true,
				"make_snatpool": false,
				"monitor_interval": 40,
				"pool_members": [
					{
						"connectionLimit": 10,
						"priorityGroup": 2,
						"serverAddresses": [
							"1.2.3.4",
							"4.5.4.7"
						],
						"servicePort": 80,
						"shareNodes": true
					}
				],
				"snat_automap": true,
				"tenant_name": "fasttcptenant",
				"virtual_address": "10.1.10.101",
				"virtual_port": 80
			},
			"tenant": "fasttcptenant",
			"application": "%[1]s",
			"operation": "create"
		}
		`, resourceName)
	})

	mux.HandleFunc("/mgmt/shared/fast/tasks/delete_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
			"id": "delete_id",
    		"code": 200,
    		"message": "success"
		}`)
	})

	mux.HandleFunc("/mgmt/shared/fast/tasks/update_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
			"id": "update_id",
			"code": 200,
			"message": "success",
			"name": "bigip-fast-templates/tcp",
			"tenant": "fasttcptenant",
			"application": "%s",
			"operation": "update"
		}
		`, resourceName)
	})

	mux.HandleFunc("/mgmt/shared/fast/applications/fasttcptenant/fasttcpapp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, `
			{
				"constants": {
					"class": "Constants",
					"fast": {
						"template": "bigip-fast-templates/tcp",
						"setHash": "2d88c05f2b7ce83e595c42c780d51b1216c0cafcc027762f6f01990d2d43105a",
						"view": {
							"app_name": "%s",
							"enable_monitor": true,
							"enable_pool": true,
							"enable_snat": true,
							"make_monitor": true,
							"make_pool": true,
							"make_snatpool": false,
							"monitor_interval": 40,
							"pool_members": [
								{
									"connectionLimit": 10,
									"priorityGroup": 2,
									"serverAddresses": [
										"1.2.3.4",
										"4.5.4.7"
									],
									"servicePort": 80,
									"shareNodes": true
								}
							],
							"snat_automap": true,
							"tenant_name": "fasttcptenant",
							"virtual_address": "10.1.10.101",
							"virtual_port": 80
						}
					}
				}
			}
			`, resourceName)
		}
		if r.Method == "PATCH" {
			fmt.Fprintf(w, `
			{
				"code": 202,
				"requestId": 301,
				"message": [
					{
						"id": "create_id",
						"name": "bigip-fast-templates/udp",
						"parameters": {
							"app_name": "%s",
							"enable_asm_logging": false,
							"enable_fallback_persistence": false,
							"enable_monitor": true,
							"enable_persistence": false,
							"enable_pool": true,
							"enable_snat": true,
							"fastl4": false,
							"make_monitor": true,
							"make_pool": true,
							"make_snatpool": false,
							"monitor_interval": 2,
							"monitor_send_string": "GET / HTTP/1.1\r\nHost: example.com\r\nConnection: Close\r\n\r\n",
							"pool_members": [
								{
									"connectionLimit": 2,
									"priorityGroup": 2,
									"serverAddresses": [
										"19.20.39.40"
									],
									"servicePort": 443,
									"shareNodes": true
								}
							],
							"snat_automap": true,
							"tenant_name": "fasttcptenant",
							"virtual_address": "15.50.30.44",
							"virtual_port": 443,
							"vlans_allow": false,
							"vlans_enable": false
						}
					}
				]
			}`, resourceName)
		}
		if r.Method == "DELETE" {
			fmt.Fprintf(w, `{"id": "delete_id"}`)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipFastTCPCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipFastTCPModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testBigipFastTCPInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_fast_tcp_app" "test-profile-http" {
  application = "%s"
  tenant      = "fasttcptenant"
  invalidkey  = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipFastTCPCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_tcp_app" "%[1]s" {
  application = "%[1]s"
  tenant      = "fasttcptenant"
  virtual_server {
	ip   = "10.1.10.101"
	port = 80
  }
  pool_members {
	addresses = ["1.2.3.4", "4.5.4.7"]
	port = 80
	connection_limit = 10
	priority_group = 2
	share_nodes = true
  }
  monitor {
	interval = 40
  }
}
provider "bigip" {
  address  = "%[2]s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipFastTCPModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_tcp_app" "%[1]s" {
  application = "%[1]s"
  tenant      = "fasttcptenant"
  virtual_server {
	ip   = "10.1.10.101"
	port = 80
  }
  pool_members {
  	addresses = ["1.2.3.4", "4.5.4.7"]
	port = 80
	connection_limit = 10
	priority_group = 2
	share_nodes = true
  }
  monitor {
	interval = 40
  }
  slow_ramp_time = 2
}
provider "bigip" {
  address  = "%[2]s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
