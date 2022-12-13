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

func TestAccBigipLtmProfileFtpUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-ftp"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileFtpInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfileFtpUnitCreate(t *testing.T) {
	resourceName := "/Common/test-profile-ftp"
	ftpDefault := "/Common/ftp"
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
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","allowFtps":"disabled","inheritParentProfile":"disabled","inheritVlanList":"disabled","port":20,"ftpsMode":"disallow","enforceTlsSessionReuse":"disabled","allowActiveMode":"enabled","translateExtended":"enabled"}`, resourceName, ftpDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp/~Common~test-profile-ftp", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","allowFtps":"disabled","inheritParentProfile":"disabled","inheritVlanList":"disabled","port":20,"ftpsMode":"disallow","enforceTlsSessionReuse":"disabled","allowActiveMode":"enabled","translateExtended":"enabled"}`, resourceName, ftpDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/~Common~test-profile-ftp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","translateExtended":"disabled"",}`, resourceName, ftpDefault)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfileFtpCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfileFtpModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfileFtpUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-profile-ftp"
	//ftpDefault := "/Common/ftp"
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
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testftpravi##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileFtpCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testftpravi##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmProfileFtpUnitReadError(t *testing.T) {
	resourceName := "/Common/test-profile-ftp"
	ftpDefault := "/Common/ftp"
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
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","allowFtps":"disabled","inheritParentProfile":"disabled","inheritVlanList":"disabled","port":20,"ftpsMode":"disallow","enforceTlsSessionReuse":"disabled","allowActiveMode":"enabled","translateExtended":"enabled"}`, resourceName, ftpDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/ftp/~Common~test-profile-ftp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","allowFtps":"disabled","inheritParentProfile":"disabled","inheritVlanList":"disabled","port":20,"ftpsMode":"disallow","enforceTlsSessionReuse":"disabled","allowActiveMode":"enabled","translateExtended":"enabled"}`, resourceName, ftpDefault)
		}
		http.Error(w, "The requested FTP Profile (/Common/test-profile-ftp) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileFtpCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested FTP Profile \\(/Common/test-profile-ftp\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfileFtpInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_ftp" "test-ftp" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfileFtpCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_ftp" "test-ftp" {
  name    = "%s"
  defaults_from = "/Common/ftp"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmProfileFtpModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_ftp" "test-ftp" {
  name    = "%s"
  defaults_from = "/Common/ftp"
  translate_extended = "disabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
