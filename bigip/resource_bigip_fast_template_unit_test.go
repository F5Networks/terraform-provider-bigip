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

func TestAccBigipFastTemplateUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-http"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipFastTemplateInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipFastTemplateUnitCreate(t *testing.T) {
	resourceName := "foo_template"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/~Common~foo_template.zip", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{}`)
	})

	mux.HandleFunc("/mgmt/shared/fast/templatesets", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{}`)
	})

	mux.HandleFunc("/mgmt/shared/fast/templatesets/~Common~foo_template", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
			"name": "/Common/foo_template",
			"hash": "abc73hd720djsh7"
		}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipFastTemplateCreate(resourceName, server.URL),
			},
			// {
			// 	Config:             testBigipFastTemplateModify(resourceName, server.URL),
			// 	ExpectNonEmptyPlan: true,
			// },
		},
	})
}

func TestAccBigipFastTemplateUnitCreateError(t *testing.T) {
	resourceName := "foo_template"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/~Common~foo_template.zip", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Create Page Not Found", 404)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipFastTemplateCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: Create Page Not Found"),
			},
		},
	})
}

func testBigipFastTemplateInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_fast_template" "foo_Template" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipFastTemplateCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_template" "foo_template" {
  name     = "/Common/%s"
  source   = "`+folder3+`/../examples/fast/foo_template.zip"
  md5_hash = filemd5("`+folder3+`/../examples/fast/foo_template.zip")
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipFastTemplateModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_fast_template" "foo_template" {
  name    = "/Common/%s"
  source   = "`+folder3+`/../examples/fast/foo_template_2.zip"
  md5_hash = filemd5("`+folder3+`/../examples/fast/foo_template_2.zip")
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
