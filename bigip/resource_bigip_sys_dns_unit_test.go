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

func TestAccBigipSysDNSUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-sys-dns"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysDNSInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipSysDNSUnitCreate(t *testing.T) {
	resourceName := "/Common/test-sys-dns"
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
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"description":"%s","nameServers":["1.1.1.1"],"search":["f5.com"]}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipSysDNSCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipSysDNSModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipSysDNSUnitReadError(t *testing.T) {
	resourceName := "/Common/test-sys-dns"
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
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.Error(w, "The requested Sys DNS (/Common/test-sys-dns) was not found", http.StatusNotFound)
		}
		_, _ = fmt.Fprintf(w, `{"description":"%s","nameServers":["1.1.1.1"],"search":["f5.com"]}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysDNSCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Sys DNS \\(/Common/test-sys-dns\\) was not found"),
			},
		},
	})
}

func TestAccBigipSysDNSUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-sys-dns"
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
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "The requested object name (/Common/testsysdns##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysDNSCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testsysdns##\\) is invalid"),
			},
		},
	})
}

func testBigipSysDNSInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_dns" "test-dns" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipSysDNSCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_dns" "test-dns" {
  description    = "%s"
  name_servers = ["1.1.1.1"]
  search = ["f5.com"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipSysDNSModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_dns" "test-dns" {
  description    = "%s"
  name_servers = ["1.1.2.1"]
  search = ["f5.com"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
