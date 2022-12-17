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

func TestAccBigipIPsecProfileUnitInvalid(t *testing.T) {
	resourceName := "test-ipsec-profile"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipIPsecProfileInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipIPsecProfileUnitCreate(t *testing.T) {
	resourceName := "test-ipsec-profile"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/tm/net/tunnels/ipsec", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
	})

	mux.HandleFunc("/mgmt/tm/net/tunnels/ipsec/~Common~test-ipsec-profile", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"name":"%[1]s", "fullPath": "/Common/%[1]s", "defaultsFrom": "/Common/ipsec"}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipIPsecProfileCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipIPsecProfileModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipIPsecProfileUnitReadError(t *testing.T) {
	resourceName := "test-ipsec-profile"
	setup()

	mux.HandleFunc("/mgmt/tm/net/tunnels/ipsec/~Common~test-ipsec-profile", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"message": "The requested IPsec Tunnel Profile (/Common/test-ipsec-profiler) was not found."}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipIPsecProfileCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: 404 page not found"),
			},
		},
	})
}

func TestAccBigipIPsecProfileUnitCreateError(t *testing.T) {
	resourceName := "test-ipsec-profile"
	setup()

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipIPsecProfileCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: 404 page not found"),
			},
		},
	})
}

func testBigipIPsecProfileInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ipsec_profile" "test-profile-profile" {
  name       = "/Common/%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipIPsecProfileCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ipsec_profile" "test-ipsec-profile" {
  name = "/Common/%s"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipIPsecProfileModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ipsec_profile" "test-ipsec-profile" {
  name        = "/Common/%s"
  description = "mytestipsecprofile"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
