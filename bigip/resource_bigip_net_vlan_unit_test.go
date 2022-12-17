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

func TestAccBigipNetVlanUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-net-vlan"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetVlanInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipNetVlanUnitCreate(t *testing.T) {
	resourceName := "/Common/test-net-vlan"
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
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","sflow":{},"tag":101}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~test-net-vlan", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "%s","sflow":{},"tag":101,"mtu": 1500,"sourceChecking": "disabled"}`, resourceName, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~test-net-vlan/interfaces", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"1.1","tagged":true}]}`)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~test-net-vlan", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "%s","sflow":{},"tag":102}`, resourceName, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipNetVlanCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipNetVlanModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipNetVlanUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-net-vlan"
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
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testnetvlan##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetVlanCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testnetvlan##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipNetVlanUnitReadError(t *testing.T) {
	resourceName := "/Common/test-net-vlan"
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
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","sflow":{},"tag":101}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~test-net-vlan/interfaces", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"name":"1.1","tagged":true}`)
		}
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~test-net-vlan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.Error(w, "The requested Net VLAN (/Common/test-net-vlan) was not found", http.StatusNotFound)
		}
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetVlanCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Net VLAN \\(/Common/test-net-vlan\\) was not found"),
			},
		},
	})
}

func testBigipNetVlanInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-vlan" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipNetVlanCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-vlan" {
  name    = "%s"
  tag = 101
  interfaces {
    vlanport = 1.1
    tagged = true
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipNetVlanModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_net_vlan" "test-vlan" {
  name    = "%s"
  tag = 102
  interfaces {
    vlanport = 1.1
    tagged = true
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
