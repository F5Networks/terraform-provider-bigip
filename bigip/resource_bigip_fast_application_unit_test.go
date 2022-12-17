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

func TestAccBigipFastAppUnitInvalid(t *testing.T) {
	resourceName := "examples/simple_http"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipFastAppUnitInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipFastAppUnitCreate(t *testing.T) {
	resourceName := "examples/simple_http"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
			"code":202,
			"requestId":1,
			"message":[
				{
					"id":"dfa86058-78ca-4384-8cf0-fb5be8229c06",
					"name":"examples/simple_http",
					"parameters":{
						"application_name":"sample_app",
						"server_addresses":["192.1.1.1","192.1.1.2"],
						"server_port":12584,
						"tenant_name":"sample_tenant",
						"virtual_address":"10.1.1.1",
						"virtual_port":8081
					}
				}
			],
			"task":"/mgmt/shared/fast/tasks/dfa86058-78ca-4384-8cf0-fb5be8229c06"
		}
	`)
	})
	mux.HandleFunc("/mgmt/shared/fast/tasks/dfa86058-78ca-4384-8cf0-fb5be8229c06", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "id": "dfa86058-78ca-4384-8cf0-fb5be8229c06",
    "code": 200,
    "message": "success",
    "name": "examples/simple_http",
    "parameters": {
        "application_name": "sample_app",
        "server_addresses": [
            "192.1.1.1",
            "192.1.1.2"
        ],
        "server_port": 12584,
        "tenant_name": "sample_tenant",
        "virtual_address": "10.1.1.1",
        "virtual_port": 8081
    },
    "tenant": "sample_tenant",
    "application": "sample_app",
    "operation": "create",
    "timestamp": "2022-12-06T11:14:24.830Z",
    "host": "localhost",
    "_links": {
        "self": "/mgmt/shared/fast/tasks/dfa86058-78ca-4384-8cf0-fb5be8229c06"
    }}`)
	})

	mux.HandleFunc("/mgmt/shared/fast/applications/sample_tenant/sample_app", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, `{
				"constants": {
					"fast": {
						"template": "examples/simple_http",
						"view": {
							"tenant_name": "sample_tenant",
							"application_name": "sample_app",
							  "virtual_port": 8081,
							  "virtual_address": "10.1.1.1",
							  "server_port": 12584,
							"server_addresses": ["192.1.1.1","192.1.1.2"]
						}
					}
				}
			}`)
		}
		if r.Method == "PATCH" {
			fmt.Fprintf(w,
				`{
					"code":202,
					"requestId":1,
					"message":[
						{
							"id":"dfa86058-78ca-4384-8cf0-fb5be8229c06",
							"name":"examples/simple_http",
							"parameters":{
								"application_name":"sample_app",
								"server_addresses":["192.1.1.1","192.1.1.2", "192.1.1.3"],
								"server_port":12584,
								"tenant_name":"sample_tenant",
								"virtual_address":"10.1.1.1",
								"virtual_port":8081
							}
						}
					],
					"task":"/mgmt/shared/fast/tasks/dfa86058-78ca-4384-8cf0-fb5be8229c06"
				}`,
			)
		}
		if r.Method == "DELETE" {
			fmt.Fprintf(w, `{"id": "dfa86058-78ca-4384-8cf0-fb5be8229c06"}`)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipFastAppUnitCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipFastAppUnitModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipFastAppUnitCreateError(t *testing.T) {
	resourceName := "examples/simple_http"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Create Page Not Found", 404)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipFastAppUnitCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: Create Page Not Found"),
			},
		},
	})
}

func testBigipFastAppUnitInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_fast_application"  "foo-app" {
  template       = "%s"
  fast_json = "{}"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipFastAppUnitCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_application"  "foo-app" {
  template    = "%s"
  fast_json = "${file("`+folder3+`/../examples/fast/new_fast_app.json")}"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipFastAppUnitModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_application"  "foo-app" {
  template    = "%s"
  fast_json = "${file("`+folder3+`/../examples/fast/new_fast_app_2.json")}"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
