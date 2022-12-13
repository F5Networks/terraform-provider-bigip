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

func TestAccBigipLtmSnatUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-snat"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmSnatInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmSnatUnitCreate(t *testing.T) {
	resourceName := "/Common/test-snat"
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
	mux.HandleFunc("/mgmt/tm/ltm/snat", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","autoLasthop": "default","fullPath": "/Common/test-snat","description": "testsnat","mirror": "disabled","sourcePort": "preserve",
"translation": "/Common/1.1.1.1","vlansDisabled": true,"origins": [{"name": "0.0.0.0/0","listenerSyncookie": "enabled"}]}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","autoLasthop": "default","fullPath": "/Common/test-snat","description": "testsnat","mirror": "disabled","sourcePort": "preserve",
"translation": "/Common/1.1.1.1","vlansDisabled": true,"origins": [{"name": "2.2.2.2/32","listenerSyncookie": "enabled"},{"name": "3.3.3.3/32","listenerSyncookie": "enabled"}]}`, resourceName)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","autoLasthop": "default","fullPath": "/Common/test-snat","description": "testsnat","mirror": "disabled","sourcePort": "preserve",
"translation": "/Common/1.1.2.1","vlansDisabled": true,"origins": [{"name": "0.0.0.0/0","listenerSyncookie": "enabled"}]}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmSnatCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmSnatModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmSnatUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-snat"
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
	mux.HandleFunc("/mgmt/tm/ltm/snat", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testsnat##) is invalid", http.StatusBadRequest)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmSnatCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testsnat##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmSnatUnitReadError(t *testing.T) {
	resourceName := "/Common/test-snat"
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
	mux.HandleFunc("/mgmt/tm/ltm/snat", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","autoLasthop": "default","description": "testsnat","mirror": "disabled","sourcePort": "preserve",
"translation": "/Common/1.1.1.1","vlansDisabled": true,"origins": [{"name": "0.0.0.0/0","listenerSyncookie": "enabled"}]}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested Snat object (/Common/test-snat) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmSnatCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Snat object \\(/Common/test-snat\\) was not found"),
			},
		},
	})
}

func testBigipLtmSnatInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_snat" "test-snat" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmSnatCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_snat" "test-snat" {
  name    = "%s"
  translation = "/Common/1.1.1.1"
  origins { name = "2.2.2.2" }
  origins { name = "3.3.3.3" }
  vlansdisabled = true
  autolasthop = "default"
  mirror = "disabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmSnatModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_snat" "test-snat" {
  name    = "%s"
  translation = "/Common/1.1.2.1"
  origins { name = "2.2.2.2" }
  origins { name = "3.3.3.3" }
  vlansdisabled = true
  autolasthop = "default"
  mirror = "disabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
