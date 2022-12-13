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

func TestAccBigipLtmPoolUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-pool"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmPoolInvalid(resourceName),
				ExpectError: regexp.MustCompile("Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmPoolUnitCreate(t *testing.T) {
	resourceName := "/Common/test-pool"
	lbmode := "round-robin"
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
	mux.HandleFunc("/mgmt/tm/ltm/pool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","loadBalancingMode":"%s"}`, resourceName, lbmode)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","loadBalancingMode":"%s","monitor":"/Common/icmp"}`, resourceName, lbmode)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","loadBalancingMode":"%s","monitor":"/Common/icmp"}`, resourceName, lbmode)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmPoolCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmPoolModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmPoolUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-pool"
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
	mux.HandleFunc("/mgmt/tm/ltm/pool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testpool##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmPoolCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testpool##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmPoolUnitReadError(t *testing.T) {
	resourceName := "/Common/test-pool"
	lbmode := "round-robin"
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
	mux.HandleFunc("/mgmt/tm/ltm/pool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","loadBalancingMode":"%s"}`, resourceName, lbmode)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","loadBalancingMode":"%s","monitor":"/Common/icmp"}`, resourceName, lbmode)
		}
		http.Error(w, "The requested pool (/Common/test-pool) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmPoolCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested pool \\(/Common/test-pool\\) was not found"),
			},
		},
	})
}

func testBigipLtmPoolInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool" "test-pool" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmPoolCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool" "test-pool" {
  name    = "%s"
  load_balancing_mode = "round-robin"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmPoolModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool" "test-pool" {
  name    = "%s"
  load_balancing_mode = "least-connections-member"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
