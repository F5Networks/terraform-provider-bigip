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

func TestAccBigipLtmProfilehttp2UnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-http2"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttp2Invalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfilehttp2UnitCreate(t *testing.T) {
	resourceName := "/Common/test-profile-http2"
	http2Default := "/Common/http2"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","concurrentStreamsPerConnection":30,"connectionIdleTimeout":100,"headerTableSize":4092,"activationModes":["always"],"enforceTlsRequirements":"enabled","frameSize":2021,"includeContentLength":"enabled","insertHeader":"disabled","receiveWindow":31,"writeSize":16380}`, resourceName, http2Default)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2/~Common~test-profile-http2", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","concurrentStreamsPerConnection":30,"connectionIdleTimeout":100,"headerTableSize":4092,"activationModes":["always"],"enforceTlsRequirements":"enabled","frameSize":2021,"includeContentLength":"enabled","insertHeader":"disabled","receiveWindow":31,"writeSize":16380}`, resourceName, http2Default)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2/~Common~test-profile-http2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","concurrentStreamsPerConnection":33,"connectionIdleTimeout":100,"headerTableSize":4092,"activationModes":["always"],"enforceTlsRequirements":"enabled","frameSize":2021,"includeContentLength":"enabled","insertHeader":"disabled","receiveWindow":31,"writeSize":16380}`, resourceName, http2Default)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfilehttp2Create(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfilehttp2Modify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfilehttp2UnitCreateError(t *testing.T) {
	resourceName := "/Common/test-profile-http2"
	//http2Default := "/Common/http2"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testhttp2ravi##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttp2Create(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testhttp2ravi##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmProfilehttp2UnitReadError(t *testing.T) {
	resourceName := "/Common/test-profile-http2"
	http2Default := "/Common/http2"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, http2Default)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2/~Common~test-profile-http2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested HTTP Profile (/Common/test-profile-http2) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttp2Create(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-http2\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfilehttp2Invalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http2" "test-profile-http2" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfilehttp2Create(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http2" "test-profile-http2" {
  name    = "%s"
  frame_size                        = 2021
  receive_window                    = 31
  write_size                        = 16380
  header_table_size                 = 4092
  include_content_length            = "enabled"
  enforce_tls_requirements          = "enabled"
  insert_header                     = "disabled"
  concurrent_streams_per_connection = 30
  connection_idle_timeout           = 100
  activation_modes                  = ["always"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmProfilehttp2Modify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http2" "test-profile-http2" {
  name    = "%s"
  frame_size                        = 2021
  receive_window                    = 31
  write_size                        = 16380
  header_table_size                 = 4092
  include_content_length            = "enabled"
  enforce_tls_requirements          = "enabled"
  insert_header                     = "disabled"
  concurrent_streams_per_connection = 40
  connection_idle_timeout           = 100
  activation_modes                  = ["always"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
