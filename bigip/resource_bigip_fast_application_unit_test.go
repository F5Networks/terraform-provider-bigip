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
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("mgmt/shared/fast/applications", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"code":202,"requestId":1,"message":[{"id":"dfa86058-78ca-4384-8cf0-fb5be8229c06","name":"examples/simple_http","parameters":{"application_name":"sample_app","server_addresses":["192.1.1.1","192.1.1.2"],"server_port":12584,"tenant_name":"sample_tenant","virtual_address":"10.1.1.1","virtual_port":8081}}],"task":"/mgmt/shared/fast/tasks/dfa86058-78ca-4384-8cf0-fb5be8229c06"}}`)
	})
	mux.HandleFunc("mgmt/shared/fast/tasks/dfa86058-78ca-4384-8cf0-fb5be8229c06", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
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

//func TestAccBigipFastAppUnitReadError(t *testing.T) {
//	resourceName := "examples/simple_http"
//	setup()
//	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//	})
//	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//		_, _ = fmt.Fprintf(w, `{}`)
//	})
//	mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		//_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName)
//	})
//	mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-http", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		http.Error(w, "The requested HTTP Profile (/Common/test-profile-http) was not found", http.StatusNotFound)
//	})
//
//	defer teardown()
//	resource.Test(t, resource.TestCase{
//		IsUnitTest: true,
//		Providers:  testProviders,
//		Steps: []resource.TestStep{
//			{
//				Config:      testBigipFastAppUnitCreate(resourceName, server.URL),
//				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-http\\) was not found"),
//			},
//		},
//	})
//}
//
//func TestAccBigipFastAppUnitCreateError(t *testing.T) {
//	resourceName := "examples/simple_http"
//	setup()
//	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//	})
//	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
//		_, _ = fmt.Fprintf(w, `{}`)
//	})
//	mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
//		//_, _ = fmt.Fprintf(w, `{"name":"/Common/testhttp##","defaultsFrom":"%s", "basicAuthRealm": "none"}`)
//		http.Error(w, "The requested object name (/Common/testravi##) is invalid", http.StatusNotFound)
//	})
//	mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-http", func(w http.ResponseWriter, r *http.Request) {
//		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
//		http.Error(w, "The requested HTTP Profile (/Common/test-profile-http) was not found", http.StatusNotFound)
//	})
//
//	defer teardown()
//	resource.Test(t, resource.TestCase{
//		IsUnitTest: true,
//		Providers:  testProviders,
//		Steps: []resource.TestStep{
//			{
//				Config:      testBigipFastAppUnitCreate(resourceName, server.URL),
//				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-http\\) was not found"),
//			},
//		},
//	})
//}

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
  fast_json = "${file("`+folder3+`/../examples/fast/new_fast_app.json")}"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
