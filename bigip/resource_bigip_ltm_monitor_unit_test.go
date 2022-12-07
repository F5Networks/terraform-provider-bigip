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

func TestAccBigipLtmMonitorUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-monitor"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmMonitorInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmMonitorUnitCreate(t *testing.T) {
	resourceName := "/Common/test-monitor"
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
	mux.HandleFunc("/mgmt/tm/ltm/monitor/http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, httpDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/monitor/http", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, httpDefault)
	})
	//mux.HandleFunc("/mgmt/tm/ltm/monitor/http/~Common~test-monitor", func(w http.ResponseWriter, r *http.Request) {
	//	_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, httpDefault)
	//})
	//mux = http.NewServeMux()
	//mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-http", func(w http.ResponseWriter, r *http.Request) {
	//	assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
	//	_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none","acceptXff": "enabled",}`, resourceName, httpDefault)
	//})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmMonitorCreate(resourceName, server.URL),
			},
			//{
			//	Config:             testBigipLtmMonitorModify(resourceName, server.URL),
			//	ExpectNonEmptyPlan: true,
			//},
		},
	})
}

func TestAccBigipLtmMonitorUnitReadError(t *testing.T) {
	resourceName := "/Common/test-monitor"
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
				Config:      testBigipLtmMonitorCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-http\\) was not found"),
			},
		},
	})
}

func TestAccBigipLtmMonitorUnitCreateError(t *testing.T) {
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
				Config:      testBigipLtmMonitorCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-http\\) was not found"),
			},
		},
	})
}

func testBigipLtmMonitorInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_monitor" "test-monitor" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmMonitorCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_monitor" "test-monitor" {
  name    = "%s"
  parent = "/Common/http"
  send = "GET /some/path\r\n"
  timeout = 999
  interval = 998
  receive = "HTTP 1.1 302 Found"
  receive_disable = "HTTP/1.1 429"
  reverse = "disabled"
  transparent = "disabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmMonitorModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_monitor" "test-monitor" {
  name    = "%s"
  parent = "/Common/http"
  send = "GET /some/path\r\n"
  timeout = 999
  interval = 996
  receive = "HTTP 1.1 302 Found"
  receive_disable = "HTTP/1.1 429"
  reverse = "disabled"
  transparent = "disabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
