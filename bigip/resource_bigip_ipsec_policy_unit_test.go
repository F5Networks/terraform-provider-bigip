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

func TestAccBigipIPSecPolicyUnitInvalid(t *testing.T) {
	resourceName := "test-ipsec-policy"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipIPSecPolicyInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipIPSecPolicyUnitCreate(t *testing.T) {
	resourceName := "test-ipsec-policy"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/tm/net/ipsec/ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
	})

	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector/~Common~test-ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, `{}`)
	})

	mux.HandleFunc("/mgmt/tm/net/ipsec/ipsec-policy/~Common~test-ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"name": "/Common/%[1]s", "fullPath": "/Common/%[1]s"}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipIPSecPolicyCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipIPSecPolicyModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// func TestAccBigipIPSecPolicyUnitExistsError(t *testing.T) {
// 	resourceName := "test-ipsec-policy"
// 	setup()
// 	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
// 		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
// 		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
// 	})

// 	mux.HandleFunc("/mgmt/tm/net/ipsec/ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
// 		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
// 	})

// 	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector/~Common~test-ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
// 		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
// 		fmt.Fprintf(w, "")
// 	})

// 	mux.HandleFunc("/mgmt/tm/net/ipsec/ipsec-policy/~Common~test-ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, `{"name": "/Common/%[1]s", "fullPath": "/Common/%[1]s"}`, resourceName)
// 	})

// 	defer teardown()
// 	resource.Test(t, resource.TestCase{
// 		IsUnitTest: true,
// 		Providers:  testProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testBigipIPSecPolicyCreate(resourceName, server.URL),
// 			},
// 			{
// 				Config:             testBigipIPSecPolicyModify(resourceName, server.URL),
// 				ExpectNonEmptyPlan: true,
// 			},
// 		},
// 	})
// }

func TestAccBigipIPSecPolicyUnitCreateError(t *testing.T) {
	resourceName := "test-ipsec-policy"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/tm/net/ipsec/ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Create Page Not Found", 404)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipIPSecPolicyCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: Create Page Not Found"),
			},
		},
	})
}

func TestAccBigipIPSecPolicyUnitModifyError(t *testing.T) {
	resourceName := "test-ipsec-policy"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/tm/net/ipsec/ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
	})

	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector/~Common~test-ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		fmt.Fprintf(w, `{}`)
	})

	mux.HandleFunc("/mgmt/tm/net/ipsec/ipsec-policy/~Common~test-ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			http.Error(w, "Modify Page Not Found", 404)
		} else {
			fmt.Fprintf(w, `{"name": "/Common/%[1]s", "fullPath": "/Common/%[1]s"}`, resourceName)
		}
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipIPSecPolicyCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipIPSecPolicyModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile("HTTP 404 :: Modify Page Not Found"),
			},
		},
	})
}

func testBigipIPSecPolicyInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ipsec_policy" "test-ipsec-policy" {
  name       = "/Common/%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipIPSecPolicyCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ipsec_policy" "test-ipsec-policy" {
  name    = "/Common/%s"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipIPSecPolicyModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ipsec_policy" "test-ipsec-policy" {
  name    = "/Common/%s"
  description = "test ipsec policy"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
