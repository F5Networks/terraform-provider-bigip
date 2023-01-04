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

func TestAccBigipSysProvisionUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-sys-provision"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysProvisionInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipSysProvisionUnitCreate(t *testing.T) {
	resourceName := "afm"
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
	mux.HandleFunc("/mgmt/tm/sys/provision/afm", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"afm","level":"none"}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipSysProvisionCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipSysProvisionModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipSysProvisionUnitReadError(t *testing.T) {
	resourceName := "afm"
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
	mux.HandleFunc("/mgmt/tm/sys/provision/afm", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.Error(w, "Status not found", http.StatusNotFound)
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"afm","level":"none"}`, resourceName)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysProvisionCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: Status not found"),
			},
		},
	})
}

func TestAccBigipSysProvisionUnitCreateError(t *testing.T) {
	resourceName := "afm"
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
	mux.HandleFunc("/mgmt/tm/sys/provision/afm", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"afm","level":"none"}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysProvisionCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: Bad Request"),
			},
		},
	})
}

func testBigipSysProvisionInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_provision" "test-sys-provision" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipSysProvisionCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_provision" "test-sys-provision" {
  name    = "%s"
  full_path  = "afm"
  cpu_ratio = 0
  disk_ratio = 0
  level = "none"
  memory_ratio = 0
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipSysProvisionModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_provision" "test-sys-provision" {
  name    = "%s"
  full_path  = "afm"
  cpu_ratio = 1
  disk_ratio = 0
  level = "none"
  memory_ratio = 0
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
