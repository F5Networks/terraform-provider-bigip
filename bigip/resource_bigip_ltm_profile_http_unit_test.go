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

func TestAccBigipLtmProfilehttpUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-http"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttpInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfilehttpUnitCreate(t *testing.T) {
	resourceName := "/Common/test-profile-http"
	httpDefault := "/Common/http"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-http", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:ltm:profile:http:httpstate",
    "name": "test-profile-http",
    "fullPath": "%s",
    "generation": 1,
    "selfLink": "https://localhost/mgmt/tm/ltm/profile/http/test-profile-http?ver=16.1.0",
    "acceptXff": "enabled",
    "appService": "none",
    "basicAuthRealm": "none",
    "defaultsFrom": "%s",
    "defaultsFromReference": {
        "link": "https://localhost/mgmt/tm/ltm/profile/http/~Common~http?ver=16.1.0"
    },
    "description": "none",
    "encryptCookies": [],
    "enforcement": {
        "allowWsHeaderName": "disabled",
        "excessClientHeaders": "reject",
        "excessServerHeaders": "reject",
        "knownMethods": [
            "CONNECT",
            "DELETE",
            "GET",
            "HEAD",
            "LOCK",
            "OPTIONS",
            "POST",
            "PROPFIND",
            "PUT",
            "TRACE",
            "UNLOCK"
        ],
        "maxHeaderCount": 64,
        "maxHeaderSize": 32768,
        "maxRequests": 0,
        "oversizeClientHeaders": "reject",
        "oversizeServerHeaders": "reject",
        "pipeline": "allow",
        "rfcCompliance": "disabled",
        "truncatedRedirects": "disabled",
        "unknownMethod": "allow"
    },
    "explicitProxy": {
        "badRequestMessage": "none",
        "badResponseMessage": "none",
        "connectErrorMessage": "none",
        "defaultConnectHandling": "deny",
        "dnsErrorMessage": "none",
        "dnsResolver": "none",
        "hostNames": [],
        "ipv6": "no",
        "routeDomain": "none",
        "tunnelName": "none",
        "tunnelOnAnyRequest": "no"
    },
    "fallbackHost": "none",
    "fallbackStatusCodes": [],
    "headerErase": "none",
    "headerInsert": "none",
    "hsts": {
        "includeSubdomains": "enabled",
        "maximumAge": 16070400,
        "mode": "disabled",
        "preload": "disabled"
    },
    "insertXforwardedFor": "disabled",
    "lwsSeparator": "none",
    "lwsWidth": 80,
    "oneconnectStatusReuse": "200 206",
    "oneconnectTransformations": "enabled",
    "proxyType": "reverse",
    "redirectRewrite": "none",
    "requestChunking": "sustain",
    "responseChunking": "sustain",
    "responseHeadersPermitted": [],
    "serverAgentName": "BigIP",
    "sflow": {
        "pollInterval": 0,
        "pollIntervalGlobal": "yes",
        "samplingRate": 0,
        "samplingRateGlobal": "yes"
    },
    "viaHostName": "none",
    "viaRequest": "preserve",
    "viaResponse": "preserve",
    "xffAlternativeNames": []
}`, resourceName, httpDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none","acceptXff": "enabled",}`, resourceName, httpDefault)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfilehttpCreate(resourceName, server.URL),
			},
			{
				Config: testBigipLtmProfilehttpModify(resourceName, server.URL),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfilehttpUnitReadError(t *testing.T) {
	resourceName := "/Common/test-profile-http"
	httpDefault := "/Common/http"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested HTTP Profile (/Common/test-profile-http) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttpCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-http\\) was not found"),
			},
		},
	})
}

func TestAccBigipLtmProfilehttpUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-profile-http"
	httpDefault := "/Common/http"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"/Common/testhttp##","defaultsFrom":"%s", "basicAuthRealm": "none"}`, httpDefault)
		http.Error(w, "The requested object name (/Common/testravi##) is invalid", http.StatusNotFound)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested HTTP Profile (/Common/test-profile-http) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfilehttpCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-http\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfilehttpInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http" "test-profile-http" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfilehttpCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http" "test-profile-http" {
  name    = "%s"
  encrypt_cookies = []
  fallback_host = "none"
  fallback_status_codes = []
  head_erase = "none"
  head_insert = "none"
  insert_xforwarded_for= "disabled"
  lws_separator = "none"
  oneconnect_transformations= "enabled"
  redirect_rewrite = "none"
  request_chunking = "sustain"
  response_chunking = "sustain"
  server_agent_name = "BigIP"
  via_host_name = "none"
  via_request = "preserve"
  via_response = "preserve"
  basic_auth_realm = "none"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmProfilehttpModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_http" "test-profile-http" {
  name    = "%s"
  accept_xff = "enabled"
  encrypt_cookies = []
  fallback_host = "none"
  fallback_status_codes = []
  head_erase = "none"
  head_insert = "none"
  insert_xforwarded_for= "disabled"
  lws_separator = "none"
  oneconnect_transformations= "enabled"
  redirect_rewrite = "none"
  request_chunking = "sustain"
  response_chunking = "sustain"
  server_agent_name = "BigIP"
  via_host_name = "none"
  via_request = "preserve"
  via_response = "preserve"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
