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

func TestAccBigipLtmProfileOneconnectUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-oneconnect"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileOneconnectInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfileOneconnectUnitCreate(t *testing.T) {
	resourceName := "/Common/test-profile-oneconnect"
	oneconnectDefault := "/Common/oneconnect"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","idleTimeoutOverride":"disabled","maxAge":3600,"maxReuse":1000,"maxSize":1000,"sharePools":"disabled"}`, resourceName, oneconnectDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect/~Common~test-profile-oneconnect", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","idleTimeoutOverride":"disabled","maxAge":3600,"maxReuse":1000,"maxSize":1000,"sharePools":"disabled"}`, resourceName, oneconnectDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect/~Common~test-profile-oneconnect", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","idleTimeoutOverride":"disabled","maxAge":3500,"maxReuse":1000,"maxSize":1000,"sharePools":"disabled"}`, resourceName, oneconnectDefault)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfileOneconnectCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfileOneconnectModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfileOneconnectUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-profile-oneconnect"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testoneconnect##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileOneconnectCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testoneconnect##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmProfileOneconnectUnitReadError(t *testing.T) {
	resourceName := "/Common/test-profile-oneconnect"
	oneconnectDefault := "/Common/oneconnect"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, oneconnectDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect/~Common~test-profile-oneconnect", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested Oneconnect Profile (/Common/test-profile-oneconnect) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileOneconnectCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Oneconnect Profile \\(/Common/test-profile-oneconnect\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfileOneconnectInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_oneconnect" "test-profile-oneconnect" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfileOneconnectCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_oneconnect" "test-profile-oneconnect" {
  name    = "%s"
  defaults_from = "/Common/oneconnect"
  idle_timeout_override = "disabled"
  max_age = 3600
  max_reuse = 1000
  max_size = 1000
  share_pools = "disabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmProfileOneconnectModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_oneconnect" "test-profile-oneconnect" {
  name    = "%s"
  defaults_from = "/Common/oneconnect"
  idle_timeout_override = "disabled"
  max_age = 3500
  max_reuse = 1000
  max_size = 1000
  share_pools = "disabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
