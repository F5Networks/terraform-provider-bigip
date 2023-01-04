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

func TestAccBigipSslCertificateUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-certificate"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSslCertificateInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipSslCertificateUnitCreate(t *testing.T) {
	resourceName := "test-certificate.crt"
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
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/test-certificate.crt", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Range", "[0-1195/1196]")
		r.Header.Add("Content-Type", "[application/octet-stream]")
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","sourcePath":"file:///var/config/rest/downloads/%s"}`, resourceName, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~Common~test-certificate.crt", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "name": "test-certificate.crt",
    "partition": "Common",
    "fullPath": "/Common/test-certificate.crt",
    "certificateKeyCurveName": "none",
    "certificateKeySize": 2048,
    "checksum": "SHA1:1338:47796744f4adf6ef012ffe47063de1f6bb86fb03",
    "createTime": "2021-07-16T05:46:24Z",
    "createdBy": "root",
    "email": "root@localhost.localdomain",
    "expirationDate": 1941774384,
    "expirationString": "Jul 14 05:46:24 2031 GMT",
    "fingerprint": "SHA256/AA:6B:8D:F5:64:5F:9F:D7:5D:73:D3:34:18:19:3C:B6:A8:98:7C:99:A4:49:33:B8:26:30:5D:36:89:13:E5:50",
    "isBundle": "false",
    "issuer": "emailAddress=root@localhost.localdomain,CN=localhost.localdomain,OU=IT,O=MyCompany,L=Seattle,ST=WA,C=US",
    "keyType": "rsa-public",
    "lastUpdateTime": "2021-07-16T05:46:24Z",
    "mode": 33188,
    "revision": 1,
    "serialNumber": "364081584",
    "size": 1338,
    "subject": "emailAddress=root@localhost.localdomain,CN=localhost.localdomain,OU=IT,O=MyCompany,L=Seattle,ST=WA,C=US",
    "systemPath": "/config/ssl/ssl.crt/default.crt",
    "updatedBy": "root",
    "version": 3}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipSslCertificateCreate(resourceName, server.URL),
			},
			{
				Config: testBigipSslCertificateModify(resourceName, server.URL),
			},
		},
	})
}

func TestAccBigipSslCertificateUnitReadError(t *testing.T) {
	resourceName := "test-certificate.crt"
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
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/test-certificate.crt", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Range", "[0-1195/1196]")
		r.Header.Add("Content-Type", "[application/octet-stream]")
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","sourcePath":"file:///var/config/rest/downloads/%s"}`, resourceName, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~Common~test-certificate.crt", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "Requested Cert (/Common/test-certificate.crt) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSslCertificateCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: Requested Cert \\(/Common/test-certificate.crt\\) was not found"),
			},
		},
	})
}

func TestAccBigipSslCertificateUnitCreateError(t *testing.T) {
	resourceName := "test-certificate.crt"
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
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/test-certificate.crt", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Range", "[0-1195/1196]")
		r.Header.Add("Content-Type", "[application/octet-stream]")
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSslCertificateCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: Bad Request"),
			},
		},
	})
}

func testBigipSslCertificateInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ssl_certificate" "test-cert" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipSslCertificateCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ssl_certificate" "test-cert" {
  name    = "%s"
  content = "${file("`+folder+`/../examples/servercert.crt")}"
  partition = "Common"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipSslCertificateModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ssl_certificate" "test-cert" {
  name    = "%s"
  content = "${file("`+folder+`/../examples/servercert2.crt")}"
  partition = "Common"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
