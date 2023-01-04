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

func TestAccBigipSslCertKeyUnitInvalid(t *testing.T) {
	resourceName := "testkey.key"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipCertKeyInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipSslCertKeyUnitCreate(t *testing.T) {
	resourceName := "testkey.key"
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
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/testkey.key", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Range", "[0-1195/1196]")
		r.Header.Add("Content-Type", "[application/octet-stream]")
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","sourcePath":"file:///var/config/rest/downloads/%s"}`, resourceName, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key/~Common~testkey.key", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "name": "testkey.key",
    "partition": "Common",
    "fullPath": "/Common/testkey.key",
    "checksum": "SHA1:1704:3daeb88d01b0efbd98a73284941a0a5c6306b6b5",
    "createTime": "2021-07-16T05:46:24Z",
    "createdBy": "root",
    "curveName": "none",
    "keySize": 2048,
    "keyType": "rsa-private",
    "lastUpdateTime": "2021-07-16T05:46:24Z",
    "mode": 33184,
    "revision": 1,
    "securityType": "normal",
    "size": 1704,
    "systemPath": "/config/ssl/ssl.key/default.key",
    "updatedBy": "root"}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipCertKeyCreate(resourceName, server.URL),
			},
			{
				Config: testBigipCertKeyModify(resourceName, server.URL),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipSslCertKeyUnitReadError(t *testing.T) {
	resourceName := "testkey.key"
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
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/testkey.key", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Range", "[0-1195/1196]")
		r.Header.Add("Content-Type", "[application/octet-stream]")
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","sourcePath":"file:///var/config/rest/downloads/%s"}`, resourceName, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key/~Common~testkey.key", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "Requested Cert Key (testkey.key) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipCertKeyCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: Requested Cert Key \\(testkey.key\\) was not found"),
			},
		},
	})
}

func TestAccBigipSslCertKeyUnitCreateError(t *testing.T) {
	resourceName := "testkey.key"
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
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/testkey.key", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Range", "[0-1195/1196]")
		r.Header.Add("Content-Type", "[application/octet-stream]")
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipCertKeyCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: Bad Request"),
			},
		},
	})
}

func testBigipCertKeyInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ssl_key" "test-cert-key" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipCertKeyCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ssl_key" "test-cert-key" {
  name    = "%s"
  content = "${file("`+folder1+`/../examples/serverkey.key")}"
  partition = "Common"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipCertKeyModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ssl_key" "test-cert-key" {
  name    = "%s"
  content = "${file("`+folder1+`/../examples/serverkey2.key")}"
  partition = "Common"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
