/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipVcmpGuestUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-tcp"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipVcmpGuestInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipVcmpGuestUnitCreate(t *testing.T) {
	resourceName := "/Common/test-vcmp"
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
	mux.HandleFunc("/mgmt/tm/vcmp/guest", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{
"kind": "tm:vcmp:guest:gueststate", "name": "%s", "fullPath": "test-vcmp",
"generation": 145,
"selfLink": "https://localhost/mgmt/tm/vcmp/guest/test-vcmp?ver=12.1.2",
"allowedSlots": [
	1,
	2
],
"coresPerSlot": 2,
"hostname": "localhost.localdomain",
"initialImage": "12.1.2.iso",
"managementGw": "none",
"managementIp": "10.1.1.1/24",
"managementNetwork": "bridged",
"minSlots": 1,
"slots": 1,
"sslMode": "shared",
"state": "configured",
"virtualDisk": "test-vcmp.img"
}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/vcmp/guest/test-vcmp", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:vcmp:guest:gueststate",
    "name": "%s",
    "fullPath": "test-vcmp",
    "generation": 133,
    "selfLink": "https://localhost/mgmt/tm/vcmp/guest/test-vcmp?ver=12.1.2",
    "allowedSlots": [
        1,
        2
    ],
    "assignedSlots": [
        1
    ],
    "coresPerSlot": 2,
    "hostname": "localhost.localdomain",
    "initialImage": "12.1.2.iso",
    "managementGw": "none",
    "managementIp": "10.1.1.1/24",
    "managementNetwork": "bridged",
    "minSlots": 1,
    "slots": 1,
    "sslMode": "shared",
    "state": "provisioned",
    "virtualDisk": "test-vcmp.img"
}`, resourceName)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/vcmp/guest/test-vcmp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Expected method 'PATCH', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{
    "kind": "tm:vcmp:guest:gueststate",
    "name": "%s",
    "fullPath": "test-vcmp",
    "generation": 146,
    "selfLink": "https://localhost/mgmt/tm/vcmp/guest/test-vcmp?ver=12.1.2",
    "allowedSlots": [
        1,
        2
    ],
    "assignedSlots": [
        1
    ],
    "coresPerSlot": 2,
    "hostname": "localhost.localdomain",
    "initialImage": "12.1.2.iso",
    "managementGw": "none",
    "managementIp": "10.1.1.1/24",
    "managementNetwork": "bridged",
    "minSlots": 1,
    "slots": 1,
    "sslMode": "shared",
    "state": "provisioned",
    "virtualDisk": "test-vcmp.img"
}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipVcmpGuestCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipVcmpGuestModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipVcmpGuestUnitCreateError(t *testing.T) {
	resourceName := "test-vcmp"
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
	mux.HandleFunc("/mgmt/tm/vcmp/guest", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (testguest##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipVcmpGuestCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile(`HTTP 400 :: The requested object name (testguest##) is invalid`),
			},
		},
	})
}
func TestAccBigipVcmpGuestUnitReadError(t *testing.T) {
	resourceName := "test-vcmp"
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
	mux.HandleFunc("/mgmt/tm/vcmp/guest", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/vcmp/guest/test-vcmp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested vCMP Guest (test-vcmp) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipVcmpGuestCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile(`HTTP 404 :: The requested vCMP Guest (test-vcmp) was not found`),
			},
		},
	})
}

func testBigipVcmpGuestInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_vcmp_guest" "test-vcmp" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}
`, resourceName)
}

func testBigipVcmpGuestCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_vcmp_guest" "test-vcmp" {
  name                = "%s"
  initial_image       = "12.1.2.iso"
  mgmt_network        = "bridged"
  mgmt_address        = "10.1.1.1/24"
  mgmt_route          = "none"
  state               = "provisioned"
  cores_per_slot      = 2
  number_of_slots     = 1
  min_number_of_slots = 1
  vlans               = ["/Common/testvlan"]
}
provider "bigip" {
  address   = "%s"
  username  = ""
  password  = ""
  login_ref = ""
}
`, resourceName, url)
}

func testBigipVcmpGuestModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_vcmp_guest" "test-vcmp" {
  name  = "%s"
  state = "configured"
}
provider "bigip" {
  address   = "%s"
  username  = ""
  password  = ""
  login_ref = ""
}
`, resourceName, url)
}
