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

func TestAccBigipNetSelfipUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-net-selfip"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetSelfipInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipNetSelfipUnitCreate(t *testing.T) {
	resourceName := "/Common/test-net-selfip"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "%s","address":"11.1.1.1/24","trafficGroup":"traffic-group-local-only","vlan":"/Common/test-vlan","allowService": ["default"]}`, resourceName, resourceName)
		}
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/net/self/~Common~test-net-selfip", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "%s","address":"11.1.1.1/24","trafficGroup": "/Common/traffic-group-local-only","vlan":"/Common/test-vlan","allowService": ["default"]}`, resourceName, resourceName)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self/~Common~test-net-selfip", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "%s","address":"11.1.1.1/24","trafficGroup":"traffic-group-local-only","vlan":"/Common/test-vlan1","allowService": ["default"]}`, resourceName, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipNetSelfipCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipNetSelfipModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipNetSelfipUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-net-selfip"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			http.Error(w, "The requested object name (/Common/testnetselfip##) is invalid", http.StatusBadRequest)
		}
		_, _ = fmt.Fprintf(w, `{}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetSelfipCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testnetselfip##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipNetSelfipUnitReadError(t *testing.T) {
	resourceName := "/Common/test-net-selfip"
	httpDefault := "/Common/http"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
		}
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/net/self/~Common~test-net-selfip", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested Net Selfip (/Common/test-net-selfip) was not found", http.StatusNotFound)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetSelfipCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Net Selfip \\(/Common/test-net-selfip\\) was not found"),
			},
		},
	})
}

func testBigipNetSelfipInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_net_selfip" "test-selfip" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipNetSelfipCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_net_selfip" "test-selfip" {
  name    = "%s"
  ip = "11.1.1.1/24"
  vlan = "/Common/test-vlan"
  port_lockdown = ["default"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipNetSelfipModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_net_selfip" "test-selfip" {
  name    = "%s"
  ip = "11.1.1.1/24"
  vlan = "/Common/test-vlan1"
  port_lockdown = ["default"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
