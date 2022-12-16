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

func TestAccBigipFastHttpsUnitInvalid(t *testing.T) {
	resourceName := "fasthttpsapp"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipFastHttpsInvalid(resourceName),
				ExpectError: regexp.MustCompile(`config is invalid: Unsupported argument: An argument named "invalid_key" is not expected here.`),
			},
		},
	})
}

func TestAccBigipFastHttpsUnitCreate(t *testing.T) {
	resourceName := "fasthttpsapp"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"code": 202,
			"requestId": 210,
			"message": [
				{
					"id": "create_id",
					"name": "bigip-fast-templates/http",
					"parameters": {
						"app_name": "%s",
						"enable_asm_logging": false,
						"enable_pool": false,
						"enable_snat": true,
						"enable_tls_client": false,
						"enable_tls_server": true,
						"enable_waf_policy": false,
						"load_balancing_mode": "least-connections-member",
						"make_monitor": false,
						"make_pool": false,
						"make_snatpool": false,
						"make_tls_client_profile": false,
						"make_tls_server_profile": false,
						"make_waf_policy": false,
						"monitor_credentials": false,
						"snat_automap": true,
						"tenant_name": "fasthttpstenant",
						"virtual_address": "15.50.30.44",
						"virtual_port": 443
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
			"name": "bigip-fast-templates/http",
			"parameters": {
				"app_name": "%s",
				"enable_asm_logging": false,
				"enable_pool": false,
				"enable_snat": true,
				"enable_tls_client": false,
				"enable_tls_server": true,
				"enable_waf_policy": false,
				"load_balancing_mode": "least-connections-member",
				"make_monitor": false,
				"make_pool": false,
				"make_snatpool": false,
				"make_tls_client_profile": false,
				"make_tls_server_profile": false,
				"make_waf_policy": false,
				"monitor_credentials": false,
				"snat_automap": true,
				"tenant_name": "fasthttpstenant",
				"virtual_address": "15.50.30.44",
				"virtual_port": 443
			},
			"tenant": "fasthttpstenant",
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
			"name": "bigip-fast-templates/http",
			"tenant": "fasthttpstenant",
			"application": "%s",
			"operation": "update"
		}
		`, resourceName)
	})

	mux.HandleFunc("/mgmt/shared/fast/applications/fasthttpstenant/fasthttpsapp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, `
			{
				"constants": {
					"class": "Constants",
					"fast": {
						"template": "bigip-fast-templates/http",
						"setHash": "2d88c05f2b7ce83e595c42c780d51b1216c0cafcc027762f6f01990d2d43105a",
						"view": {
							"app_name":"%s",
							"enable_asm_logging":false,
							"enable_pool":false,
							"enable_snat":true,
							"enable_tls_client":false,
							"enable_tls_server":true,
							"enable_waf_policy":false,
							"load_balancing_mode":"least-connections-member",
							"make_monitor":false,
							"make_pool":false,
							"make_snatpool":false,
							"make_tls_client_profile":false,
							"make_tls_server_profile":false,
							"make_waf_policy":false,
							"monitor_credentials":false,
							"snat_automap":true,
							"tenant_name":"fasthttpstenant",
							"virtual_address":"15.50.30.44",
							"virtual_port":443
						}
					}
				}
			}
			`, resourceName)
		}
		if r.Method == "PATCH" {
			fmt.Fprintf(w, `
			{
				"code": 200,
				"requestId": 178,
				"message": [
					{
						"id": "update_id",
						"name": "bigip-fast-templates/http",
						"parameters":{
							"app_name":"%s",
							"enable_asm_logging":false,
							"enable_pool":true,
							"enable_snat":true,
							"enable_tls_client":false,
							"enable_tls_server":true,
							"enable_waf_policy":false,
							"load_balancing_mode":"least-connections-member",
							"make_monitor":true,
							"make_pool":true,
							"make_snatpool":false,
							"make_tls_client_profile":false,
							"make_tls_server_profile":false,
							"make_waf_policy":false,
							"monitor_credentials":true,
							"snat_automap":true,
							"tenant_name":"fasthttpstenant",
							"virtual_address":"15.50.30.44",
							"virtual_port":443,
							"enable_monitor":true,
							"monitor_interval":2,
							"monitor_passphrase":"x5$ie02",
							"monitor_username":"abc",
							"pool_members":[
								{
									"connectionLimit":2,
									"priorityGroup":2,
									"serverAddresses":["19.20.39.40"],
									"servicePort":443,
									"shareNodes":true
								}
							]
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
				Config: testBigipFastHttpsCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipFastHttpsModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testBigipFastHttpsInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_fast_https_app" "fasthttpsapp" {
  tenant      = "fasthttpstenant"
  application = "%s"
  invalid_key = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipFastHttpsCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_https_app" "fasthttpsapp" {
  tenant      = "fasthttpstenant"
  application = "%s"
  virtual_server {
    ip   = "15.50.30.44"
	port = 443
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipFastHttpsModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_https_app" "fasthttpsapp" {
  tenant      = "fasthttpstenant"
  application = "%s"
  virtual_server {
    ip   = "15.50.30.44"
	port = 443
  }
  monitor {
	monitor_auth = true
	username     = "abc"
	password     = "x5$ie02"
	interval     = 2
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
