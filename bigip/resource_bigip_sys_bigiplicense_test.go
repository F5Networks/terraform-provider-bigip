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

func TestAccBigipsysLicenseUnitInvalid(t *testing.T) {
	regKey := "testkey.key"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysLicenseInvalid(regKey),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipsysLicenseUnitCreate(t *testing.T) {
	regKey := "Z6298-13163-33476-87907-5067310"
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
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"command": "install", "options":[ { "registration-key": "%s" } ]}`, regKey)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipSysLicenseCreate(regKey, server.URL),
			},
			{
				Config: testBigipSysLicenseModify(regKey, server.URL),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipsysLicenseUnitReadError(t *testing.T) {
	regKey := "Z6298-13163-33476-87907-5067310"
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
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.Error(w, "License Not installed", http.StatusNotFound)
		}
		_, _ = fmt.Fprintf(w, `{"command": "install", "options":[ { "registration-key": "%s" } ]}`, regKey)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysLicenseCreate(regKey, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: License Not installed"),
			},
		},
	})
}

func TestAccBigipsysLicenseUnitCreateError(t *testing.T) {
	regKey := "Z6298-13163-33476-87907-5067310"
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
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysLicenseCreate(regKey, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: Bad Request"),
			},
		},
	})
}

func testBigipSysLicenseInvalid(regKey string) string {
	return fmt.Sprintf(`
resource "bigip_sys_bigiplicense" "test-sys-license" {
  command       = "install"
  registration_key = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, regKey)
}

func testBigipSysLicenseCreate(regKey, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_bigiplicense" "test-sys-license" {
  command    = "install"
  registration_key = "%s"
  timeout = 3
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, regKey, url)
}

func testBigipSysLicenseModify(regKey, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_bigiplicense" "test-sys-license" {
  command    = "install"
  registration_key = "%s"
  timeout = 4
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, regKey, url)
}
