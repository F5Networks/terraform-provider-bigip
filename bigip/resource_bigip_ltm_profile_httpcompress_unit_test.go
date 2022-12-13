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

func TestAccBigipLtmProfilehttpcompressUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-httpcompress"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttpCompressInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfilehttpcompressUnitCreate(t *testing.T) {
	resourceName := "/Common/test-profile-httpcompress"
	httpcompressDefault := "/Common/httpcompression"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "%s","defaultsFrom":"%s", "contentTypeExclude":["nicecontentexclude.com"],"contentTypeInclude":["nicecontent.com"],"uriExclude":["f5.com"],"uriInclude":["cisco.com"]}`, resourceName, resourceName, httpcompressDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression/~Common~test-profile-httpcompress", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "%s","defaultsFrom":"%s", "contentTypeExclude":["nicecontentexclude.com"],"contentTypeInclude":["nicecontent.com"],"uriExclude":["f5.com"],"uriInclude":["cisco.com"]}`, resourceName, resourceName, httpcompressDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression/~Common~test-profile-httpcompress", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "%s","defaultsFrom":"%s", "contentTypeExclude":["nicecontentexclude.com"],"contentTypeInclude":["nicecontent.com"],"uriExclude":["f5.com","f5.net"],"uriInclude":["cisco.com"]}`, resourceName, resourceName, httpcompressDefault)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfilehttpCompressCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfilehttpCompressModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfilehttpcompressUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-profile-httpcompress"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testravi##) is invalid", http.StatusBadRequest)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttpCompressCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testravi##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmProfilehttpcompressUnitReadError(t *testing.T) {
	resourceName := "/Common/test-profile-httpcompress"
	httpcompressDefault := "/Common/httpcompression"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpcompressDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression/~Common~test-profile-httpcompress", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested HTTP Profile (/Common/test-profile-httpcompress) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttpCompressCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-httpcompress\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfilehttpCompressInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_httpcompress" "test-profile-httpcompress" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfilehttpCompressCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_httpcompress" "test-profile-httpcompress" {
  name    = "%s"
  defaults_from = "/Common/httpcompression"
  uri_exclude = ["f5.com"]
  uri_include = ["cisco.com"]
  content_type_include = ["nicecontent.com"]
  content_type_exclude = ["nicecontentexclude.com"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmProfilehttpCompressModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_httpcompress" "test-profile-httpcompress" {
  name    = "%s"
  defaults_from = "/Common/httpcompression"
  uri_exclude = ["f5.com","f5.net"]
  uri_include = ["cisco.com"]
  content_type_include = ["nicecontent.com"]
  content_type_exclude = ["nicecontentexclude.com"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
