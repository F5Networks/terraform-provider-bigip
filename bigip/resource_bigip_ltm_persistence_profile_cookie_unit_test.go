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

func TestAccBigipLtmProfileCookieUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-ppcookie"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileCookieInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfileCookieUnitCreate(t *testing.T) {
	resourceName := "/Common/test-ppcookie"
	cookieDefault := "/Common/cookie"
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
	mux.HandleFunc("/mgmt/tm/ltm/persistence/cookie", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, cookieDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/persistence/cookie/~Common~test-ppcookie", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, cookieDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/cookie/~Common~test-ppcookie", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Expected method 'PATCH', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, cookieDefault)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfileCookieCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfileCookieModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfileCookieUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-ppcookie"
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
	mux.HandleFunc("/mgmt/tm/ltm/persistence/cookie", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testppcookie##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileCookieCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testppcookie##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmProfileCookieUnitReadError(t *testing.T) {
	resourceName := "/Common/test-ppcookie"
	cookieDefault := "/Common/cookie"
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
	mux.HandleFunc("/mgmt/tm/ltm/persistence/cookie", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, cookieDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/persistence/cookie/~Common~test-ppcookie", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, cookieDefault)
		}
		http.Error(w, "The requested Persist cookie Profile (/Common/test-test-ppcookie) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileCookieCreateReadError(resourceName, server.URL),
				ExpectError: regexp.MustCompile("The requested Persist cookie Profile \\(/Common/test-test-ppcookie\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfileCookieInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_cookie" "test_ppcookie" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfileCookieCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_cookie" "test_ppcookie" {
  name    = "%s"
  defaults_from = "/Common/cookie"
  match_across_pools = ""
  match_across_services = ""
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmProfileCookieModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_cookie" "test_ppcookie" {
  name    = "%s"
  defaults_from = "/Common/cookie"
  match_across_pools = "enabled"
  match_across_services = "enabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmProfileCookieCreateReadError(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_cookie" "test_ppcookie_read" {
  name    = "%s"
  defaults_from = "/Common/cookie"
  match_across_pools = ""
  match_across_services = ""
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
