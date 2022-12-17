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

func TestAccBigipLtmIruleUnitInvalid(t *testing.T) {
	resourceName := "test_irule"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmIruleInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmIruleUnitCreate(t *testing.T) {
	iruleName := "test_irule"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/ltm/rule", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		}
	})
	mux.HandleFunc(fmt.Sprintf("/mgmt/tm/ltm/rule/~Common~%s", iruleName), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"name":"/Common/%[1]s", "apiAnonymous":"test_%[1]s", "fullPath":"/Common/%[1]s"}`, iruleName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmIruleCreate(iruleName, server.URL),
			},
			{
				Config:             testBigipLtmIruleCreate(iruleName, server.URL),
				ExpectNonEmptyPlan: false,
			},
			{
				Config:             testBigipLtmIruleModify(iruleName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmIruleUnitCreateError(t *testing.T) {
	iruleName := "test_irule"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/ltm/rule", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Create Page Not Found", 404)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmIruleCreate(iruleName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: Create Page Not Found"),
			},
		},
	})
}

func testBigipLtmIruleInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_irule" "test-rule" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmIruleCreate(iruleName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_irule" "test-rule" {
  name    = "/Common/%[1]s"
  irule = "test_%[1]s"
}
provider "bigip" {
  address  = "%[2]s"
  username = ""
  password = ""
  login_ref = ""
}`, iruleName, url)
}

func testBigipLtmIruleModify(iruleName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_irule" "test-rule" {
  name    = "/Common/%[1]s"
  irule   = "test2_%[1]s"
}
provider "bigip" {
  address  = "%[2]s"
  username = ""
  password = ""
  login_ref = ""
}`, iruleName, url)
}
