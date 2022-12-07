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

func TestAccBigipLtmProfileDestAddrUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-ppdstaddr"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileDestAddrInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfileDestAddrUnitCreate(t *testing.T) {
	resourceName := "/Common/test-ppdstadr"
	dstAddrDefault := "/Common/dest_addr"
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
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, dstAddrDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr/~Common~test-ppdstadr", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, dstAddrDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr/~Common~test-ppdstadr", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Expected method 'PATCH', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, dstAddrDefault)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfileDestAddrCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfileDestAddrModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfileDestAddrUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-ppdstadr"
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
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testppdstaddr##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileDestAddrCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testppdstaddr##\\) is invalid"),
			},
		},
	})
}
func TestAccBigipLtmProfileDestAddrUnitReadError(t *testing.T) {
	resourceName := "/Common/test-ppdstadr"
	dstAddrDefault := "/Common/dest_addr"
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
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, dstAddrDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr/~Common~test-ppdstadr", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, dstAddrDefault)
		}
		http.Error(w, "The requested Persist DstAddr Profile (/Common/test-ppdstadr) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileDestAddrCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Persist DstAddr Profile \\(/Common/test-ppdstadr\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfileDestAddrInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_dstaddr" "test_ppdstaddr" {
  name       = "%s"
  defaults_from = "/Common/dest_addr"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfileDestAddrCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_dstaddr" "test_ppdstaddr" {
  name    = "%s"
  defaults_from = "/Common/dest_addr"
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

func testBigipLtmProfileDestAddrModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_persistence_profile_dstaddr" "test_ppdstaddr" {
  name    = "%s"
  defaults_from = "/Common/dest_addr"
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
