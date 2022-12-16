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

func TestAccBigipFastUdpUnitInvalid(t *testing.T) {
	resourceName := "fastudpapp"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipFastUDPInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipFastUdpUnitCreate(t *testing.T) {
	resourceName := "fastudpapp"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
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
						"tenant_name": "fastudptenant",
						"virtual_address": "15.50.30.44",
						"virtual_port": 443,
						"vlans_allow": false,
						"vlans_enable": false
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
			"name": "bigip-fast-templates/udp",
			"parameters": {
				"app_name": "%[1]s",
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
				"tenant_name": "fastudptenant",
				"virtual_address": "15.50.30.44",
				"virtual_port": 443,
				"vlans_allow": false,
				"vlans_enable": false
			},
			"tenant": "fastudptenant",
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
			"name": "bigip-fast-templates/udp",
			"tenant": "fastudptenant",
			"application": "%s",
			"operation": "update"
		}
		`, resourceName)
	})

	mux.HandleFunc("/mgmt/shared/fast/applications/fastudptenant/fastudpapp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, `
			{
				"constants": {
					"class": "Constants",
					"fast": {
						"template": "bigip-fast-templates/udp",
						"setHash": "2d88c05f2b7ce83e595c42c780d51b1216c0cafcc027762f6f01990d2d43105a",
						"view": {
							"app_name": "%[1]s",
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
							"tenant_name": "fastudptenant",
							"virtual_address": "15.50.30.44",
							"virtual_port": 443,
							"vlans_allow": false,
							"vlans_enable": false
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
							"tenant_name": "fastudptenant",
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
				Config: testBigipFastUDPCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipFastUDPModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testBigipFastUDPInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_fast_udp_app" "%[1]s" {
  tenant      = "fastudptenant"
  application = "%[1]s"
  invalidkey  = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipFastUDPCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_udp_app" "%[1]s" {
  tenant      = "fastudptenant"
  application = "%[1]s"
  virtual_server {
	ip   = "15.50.30.44"
	port = 443
  }

  monitor {
	interval    = 2
	send_string = "GET / HTTP/1.1\r\nHost: example.com\r\nConnection: Close\r\n\r\n"
  }

  pool_members {
	addresses = ["19.20.39.40"]
	port = 443
	connection_limit = 2
	priority_group = 2
	share_nodes = true
  }
}
provider "bigip" {
  address  = "%[2]s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipFastUDPModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_udp_app" "%[1]s" {
  tenant      = "fastudptenant"
  application = "%[1]s"
  virtual_server {
	ip   = "15.50.30.44"
	port = 443
  }

  monitor {
	interval    = 2
	send_string = "GET / HTTP/1.1\r\nHost: example.com\r\nConnection: Close\r\n\r\n"
  }

  slow_ramp_time = 2

  pool_members {
	addresses = ["19.20.39.40"]
	port = 443
	connection_limit = 2
	priority_group = 2
	share_nodes = true
  }
}
provider "bigip" {
  address  = "%[2]s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
