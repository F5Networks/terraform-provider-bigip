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

func TestAccBigipFastHttpUnitInvalid(t *testing.T) {
	resourceName := "fasthttpapp"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipFastHttpInvalid(resourceName),
				ExpectError: regexp.MustCompile(`config is invalid: Unsupported argument: An argument named "invalid_key" is not expected here.`),
			},
		},
	})
}

func TestAccBigipFastHttpUnitCreate(t *testing.T) {
	resourceName := "fasthttpapp"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		{
			"code": 202,
			"requestId": 171,
			"message": [
				{
					"id": "create_id",
					"name": "bigip-fast-templates/http",
					"parameters": {
						"app_name": "%s",
						"enable_asm_logging": false,
						"enable_monitor": true,
						"enable_pool": false,
						"enable_snat": true,
						"enable_tls_client": false,
						"enable_tls_server": false,
						"enable_waf_policy": false,
						"make_monitor": false,
						"make_pool": false,
						"make_snatpool": false,
						"make_tls_client_profile": false,
						"make_tls_server_profile": false,
						"make_waf_policy": false,
						"monitor_credentials": false,
						"monitor_name_http": "/Common/http",
						"snat_automap": true,
						"tenant_name": "fasthttptenant",
						"virtual_address": "10.30.30.44",
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
				"enable_monitor": true,
				"enable_pool": false,
				"enable_snat": true,
				"enable_tls_client": false,
				"enable_tls_server": false,
				"enable_waf_policy": false,
				"make_monitor": false,
				"make_pool": false,
				"make_snatpool": false,
				"make_tls_client_profile": false,
				"make_tls_server_profile": false,
				"make_waf_policy": false,
				"monitor_credentials": false,
				"monitor_name_http": "/Common/http",
				"snat_automap": true,
				"tenant_name": "fasthttptenant",
				"virtual_address": "10.30.30.44",
				"virtual_port": 443
			},
			"tenant": "fasthttptenant",
			"application": "fasthttpapp",
			"operation": "create",
			"timestamp": "2022-12-14T13:59:36.656Z",
			"host": "localhost",
			"_links": {
				"self": "/mgmt/shared/fast/tasks/create_id"
			}
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
			"tenant": "fasthttptenant",
			"application": "%s",
			"operation": "update"
		}
		`, resourceName)
	})

	mux.HandleFunc("/mgmt/shared/fast/applications/fasthttptenant/fasthttpapp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, `
			{
				"constants": {
					"class": "Constants",
					"fast": {
						"template": "bigip-fast-templates/http",
						"setHash": "2d88c05f2b7ce83e595c42c780d51b1216c0cafcc027762f6f01990d2d43105a",
						"view": {
							"app_name": "%s",
							"enable_asm_logging": false,
							"enable_monitor": true,
							"enable_pool": false,
							"enable_snat": true,
							"enable_tls_client": false,
							"enable_tls_server": false,
							"enable_waf_policy": false,
							"make_monitor": true,
							"make_pool": false,
							"make_snatpool": false,
							"make_tls_client_profile": false,
							"make_tls_server_profile": false,
							"make_waf_policy": false,
							"monitor_credentials": true,
							"monitor_interval": 2,
							"monitor_name_http": "/Common/http",
							"monitor_passphrase": "x5$ie02",
							"monitor_username": "abc",
							"snat_automap": true,
							"tenant_name": "fasthttptenant",
							"virtual_address": "10.30.30.44",
							"virtual_port": 443
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
						"parameters": {
							"app_name": "%s",
							"enable_asm_logging": false,
							"enable_monitor": true,
							"enable_pool": false,
							"enable_snat": true,
							"enable_tls_client": false,
							"enable_tls_server": false,
							"enable_waf_policy": false,
							"make_monitor": true,
							"make_pool": false,
							"make_snatpool": false,
							"make_tls_client_profile": false,
							"make_tls_server_profile": false,
							"make_waf_policy": false,
							"monitor_credentials": true,
							"monitor_name_http": "/Common/http",
							"snat_automap": true,
							"tenant_name": "fasthttptenant",
							"virtual_address": "10.30.30.44",
							"virtual_port": 443,
							"monitor_interval": 2,
							"monitor_passphrase": "x5$ie02",
							"monitor_username": "abc"
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
				Config: testBigipFastHttpCreate(resourceName, server.URL),
			},
			{
				Config: testBigipFastHttpModify(resourceName, server.URL),
			},
		},
	})
}

func testBigipFastHttpInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_fast_http_app" "fasthttpapp" {
  tenant      = "fasthttptenant"
  application = "%s"
  invalid_key = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipFastHttpCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_http_app" "fasthttpapp" {
  tenant      = "fasthttptenant"
  application = "%s"
  virtual_server {
    ip   = "10.30.30.44"
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

func testBigipFastHttpModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_http_app" "fasthttpapp" {
  tenant      = "fasthttptenant"
  application = "%s"
  virtual_server {
	ip   = "10.30.30.44"
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
