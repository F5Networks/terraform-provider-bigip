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
	resourceName := "/Common/test-certificate"
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
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~Common~test-certificate", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~Common~test-certificate", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none","acceptXff": "enabled",}`, resourceName, httpDefault)
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
				Config:             testBigipSslCertificateModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipSslCertificateUnitReadError(t *testing.T) {
	resourceName := "/Common/test-certificate"
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
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~Common~test-certificate", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested HTTP Profile (/Common/test-certificate) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSslCertificateCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-certificate\\) was not found"),
			},
		},
	})
}

func TestAccBigipSslCertificateUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-certificate"
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
	//	certByte := []byte(`-----BEGIN CERTIFICATE-----
	//MIIDRjCCAi4CCQC6Dx6jDXj7dzANBgkqhkiG9w0BAQsFADBlMQswCQYDVQQGEwJJ
	//TjELMAkGA1UECAwCVFMxDDAKBgNVBAcMA0hZRDESMBAGA1UECgwJRWNvc3lzdGVt
	//MRIwEAYDVQQLDAlFY29zeXN0ZW0xEzARBgNVBAMMCnd3dy5mNS5jb20wHhcNMTkx
	//MTIwMDY0MTI4WhcNMjAxMTE5MDY0MTI4WjBlMQswCQYDVQQGEwJJTjELMAkGA1UE
	//CAwCVFMxDDAKBgNVBAcMA0hZRDESMBAGA1UECgwJRWNvc3lzdGVtMRIwEAYDVQQL
	//DAlFY29zeXN0ZW0xEzARBgNVBAMMCnd3dy5mNS5jb20wggEiMA0GCSqGSIb3DQEB
	//AQUAA4IBDwAwggEKAoIBAQDXSTUmCJBauE3DXb1YmDHFP/aTXzjQVBxbLUXvv9Vf
	//yxPvteH3l0RuxPJCOzTCpSArYJ5MDlxjH366MrsXJWjBVuucidWSFGDikmlvDEhW
	//Cb9KemK6300cD3hSwq0O7heY6klJ0VnLGNk1uuQdTwfPUM7ZRZzCP5TRiRls8Hi5
	//M4S/h1/9Pqf6j8/5pzwH5juoD+UeboWf9hIM5LYUDR+v/7+ymBvaAa6Jl9pUjAtH
	//yiN1swqWAMjGYYwbBpSrFqPLXaSZE/z8dLUZecI6ZMz+yA0Y9JZ3e4A7EDLsSvwd
	//y5q4mWBMsXzlhiX6c8wWBmhhwqZu3I4WA6ipUv+wWET5AgMBAAEwDQYJKoZIhvcN
	//AQELBQADggEBABsim7iVvVhL3RT4oA+sbvSDp1lDhiBS2eKcKqnIT0GSROoNpJIN
	//s3uUD5XUz9oBxLbD3p6uiDrfqvmKTBpbp7YJWYqGbcsG06J392DLTaC/6KPb4D/x
	//GSLpiyzYPP+YlbBp6VZXQbfx+GGr9UJx2E/Q0rmHVgUx0zFv3I+6rHGVKGA2E61X
	//8M2fsrkzCFCk8owrDHPV27vXXgUI6bAQNbcJpYb4BCv5eO3zjJFxI0ljWL9LpHDF
	//AJcu6l4kx4Jpo5lsExqC8QTctHRu2MIZM1MUdml+YyV1Rjb7W5WfL1vgbOVX7O1C
	//IN0JSq/C/zyaw90UuKb48HeO7aqrkNlmd/0=
	//-----END CERTIFICATE-----`)
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/~Common~test-certificate", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/octet-stream")
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		//_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/sys/file/ssl-cert", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testsslcert##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSslCertificateCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testsslcert##\\) is invalid"),
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
