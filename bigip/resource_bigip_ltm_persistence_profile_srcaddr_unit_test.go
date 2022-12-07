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

func TestAccBigipLtmProfileSrcAddrUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-ppsrcaddr"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileSrcAddrInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfileSrcAddrUnitCreate(t *testing.T) {
	resourceName := "/Common/test-ppsrcaddr"
	srcppAddr := "/Common/source_addr"
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
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, srcppAddr)
	})
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr/~Common~test-ppsrcaddr", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, srcppAddr)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr/~Common~test-ppsrcaddr", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Expected method 'PATCH', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, srcppAddr)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfileSrcAddrCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfileSrcAddrModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfileSrcAddrUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-ppsrcaddr"
	//srcppAddr := "/Common/source_addr"
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
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testppsrcaddr##) is invalid", http.StatusBadRequest)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileSrcAddrCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testppsrcaddr##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmProfileSrcAddrUnitReadError(t *testing.T) {
	resourceName := "/Common/test-ppsrcaddr"
	srcppAddr := "/Common/source_addr"
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
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, srcppAddr)
	})
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr/~Common~test-ppsrcaddr", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, srcppAddr)
		}
		http.Error(w, "The requested Persist SrcAddr Profile (/Common/test-ppsrcadr) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileSrcAddrCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Persist SrcAddr Profile \\(/Common/test-ppsrcadr\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfileSrcAddrInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_srcaddr" "test_ppsrcaddr" {
  name       = "%s"
  defaults_from = "/Common/source_addr"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfileSrcAddrCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_srcaddr" "test_ppsrcaddr" {
  name    = "%s"
  defaults_from = "/Common/source_addr"
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

func testBigipLtmProfileSrcAddrModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_srcaddr" "test_ppsrcaddr" {
  name    = "%s"
  defaults_from = "/Common/source_addr"
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
