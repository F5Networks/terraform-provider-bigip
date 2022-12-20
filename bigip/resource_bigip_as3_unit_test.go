/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipAS3UnitInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipAs3configInvalid(),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipAS3UnitCreate(t *testing.T) {
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/shared/appsvcs/info", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		io.WriteString(w, `{
			"version": "3.41.0",
			"release": "1",
			"schemaCurrent": "3.41.0",
			"schemaMinimum": "3.0.0"
		}`)
	})

	mux.HandleFunc("/mgmt/shared/appsvcs/declare/Sample_02,Sample_01?async=true",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
			io.WriteString(w, `
			{
				"id": "create_id",
				"results": [
					{
						"message": "Declaration successfully submitted",
						"tenant": "",
						"host": "",
						"runTime": 0,
						"code": 0
					}
				]
			}
			`)
		})

	mux.HandleFunc("/mgmt/shared/appsvcs/declare/Sample_02,Sample_01",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" || r.Method == "DELETE" {
				io.WriteString(w, `
			{
				"id": "create_id",
				"results": [
					{
						"message": "Declaration successfully submitted",
						"tenant": "",
						"host": "",
						"runTime": 0,
						"code": 0
					}
				]
			}
			`)
			}
			if r.Method == "GET" {
				fmt.Fprintf(w, `{
					"action":"deploy",
					"class":"AS3",
					"declaration":{
						"Sample_01":{
							"Application_1":{
								"class":"Application",
								"serviceMain":{
									"class":"Service_HTTP",
									"pool":"web_pool1",
									"virtualAddresses":["10.1.2.12"]
								},
								"template":"http",
								"web_pool1":{
									"class":"Pool",
									"members":[
										{
											"serverAddresses":["192.1.1.102","192.1.1.112"],
											"servicePort":80
										}
									],
									"monitors":["http"]
								}
							},
							"class":"Tenant",
							"defaultRouteDomain":0
						},
						"Sample_02":{
							"Application_2":{
								"class":"Application",
								"serviceMain":{
									"class":"Service_HTTP",
									"pool":"web_pool2",
									"virtualAddresses":["10.0.2.12"]
								},
								"template":"http",
								"web_pool2":{
									"class":"Pool",
									"members":[
										{
											"serverAddresses":["193.0.1.101","193.0.1.112"],
											"servicePort":80
										}
									],
									"monitors":["http"]
								}
							},
							"class":"Tenant",
							"defaultRouteDomain":0
						},
					"class":"ADC",
					"id":"example-declaration-01",
					"label":"Sample 1",
					"remark":"Simple HTTP application with round robin pool",
					"schemaVersion":"3.0.0"
				},
				"persist":true
				}`)
			}
		},
	)

	mux.HandleFunc("/mgmt/shared/appsvcs/task/create_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
			"id": "create_id",
			"results": [
				{
					"code": 200,
					"message": "success",
					"lineCount": 25,
					"host": "localhost",
					"tenant": "Sample_01",
					"runTime": 1284
				},
				{
					"code": 200,
					"message": "success",
					"lineCount": 25,
					"host": "localhost",
					"tenant": "Sample_02",
					"runTime": 1015
				}
			]
		}`)
	})

	mux.HandleFunc("/mgmt/shared/appsvcs/task/delete_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
			"id": "delete_id",
			"results": [
				{
					"code": 200,
					"message": "success",
					"lineCount": 28,
					"host": "localhost",
					"tenant": "Sample_01",
					"runTime": 1779
				},
				{
					"code": 200,
					"message": "success",
					"lineCount": 28,
					"host": "localhost",
					"tenant": "Sample_02",
					"runTime": 1670
				}
			]
		}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:             testBigipAs3configCreate(server.URL),
				ExpectNonEmptyPlan: true,
			},
			{
				Config:             testBigipAs3configModify(server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipAS3UnitReadError(t *testing.T) {
	resourceName := "/Common/test-profile-http"
	httpDefault := "/Common/http"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested HTTP Profile (/Common/test-profile-http) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttpCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-http\\) was not found"),
			},
		},
	})
}

func TestAccBigipAS3UnitCreateError(t *testing.T) {
	resourceName := "/Common/test-profile-http"
	httpDefault := "/Common/http"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"/Common/testhttp##","defaultsFrom":"%s", "basicAuthRealm": "none"}`, httpDefault)
		http.Error(w, "The requested object name (/Common/testravi##) is invalid", http.StatusNotFound)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested HTTP Profile (/Common/test-profile-http) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttpCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-http\\) was not found"),
			},
		},
	})
}

func testBigipAs3configInvalid() string {
	return fmt.Sprintf(`
resource "bigip_as3" "test-as3" {
  as3_json       = "{}"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`)
}

func testBigipAs3configCreate(url string) string {
	return fmt.Sprintf(`
resource "bigip_as3" "test-as3" {
  as3_json = "${file("%s/../examples/as3/as3_example1.json")}"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, folder3, url)
}

func testBigipAs3configModify(url string) string {
	return fmt.Sprintf(`
resource "bigip_as3" "test-as3" {
  as3_json = "${file("%s/../examples/as3/as3_example1.json")}"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, folder3, url)
}
