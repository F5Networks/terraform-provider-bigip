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

func TestUnitBigipVcmpGuestInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-tcp"
	setup()
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipVcmpGuestInvalid(resourceName),
				ExpectError: regexp.MustCompile("An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestUnitBigipVcmpGuestCreate(t *testing.T) {
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
	// The client addresses the object by its mangled full name
	// (/Common/test-vcmp -> ~Common~test-vcmp). This single handler serves the
	// post-create read (GET) as well as the update read/PATCH in step 2.
	mux.HandleFunc("/mgmt/tm/vcmp/guest/~Common~test-vcmp", func(w http.ResponseWriter, r *http.Request) {
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
	// Destroy path: DeleteVcmpGuest DELETEs the mangled guest path (handled by
	// the handler above, which responds to any method), then — because
	// delete_virtual_disk defaults to true and virtual_disk is set — lists and
	// deletes the associated virtual disk.
	mux.HandleFunc("/mgmt/tm/vcmp/virtual-disk", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"test-vcmp.img"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/vcmp/virtual-disk/test-vcmp.img", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "Expected method 'DELETE', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipVcmpGuestCreate(resourceName),
				// The mock cannot round-trip every computed attribute exactly,
				// so a non-empty follow-up plan is expected for this unit test.
				ExpectNonEmptyPlan: true,
			},
			{
				Config:             testBigipVcmpGuestModify(resourceName),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestUnitBigipVcmpGuestCreateError(t *testing.T) {
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
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipVcmpGuestCreate(resourceName),
				ExpectError: regexp.MustCompile(`HTTP 400 :: The requested object name \(testguest##\) is invalid`),
			},
		},
	})
}
func TestUnitBigipVcmpGuestReadError(t *testing.T) {
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
	// The post-create read addresses the object by name and returns 404.
	mux.HandleFunc("/mgmt/tm/vcmp/guest/test-vcmp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested vCMP Guest (test-vcmp) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipVcmpGuestCreate(resourceName),
				ExpectError: regexp.MustCompile(`HTTP 404 :: The requested vCMP Guest \(test-vcmp\) was not found`),
			},
		},
	})
}

func TestUnitBigipVcmpGuestBridgedMissingAddress(t *testing.T) {
	// Covers Create's early-return validation branch: mgmt_network ==
	// "bridged" with no mgmt_address set fails before any HTTP call.
	resourceName := "test-vcmp"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipVcmpGuestBridgedNoAddress(resourceName),
				ExpectError: regexp.MustCompile("the mgmt_address must be provided if mgmt_network is set to bridged"),
			},
		},
	})
}

func TestUnitBigipVcmpGuestUpdate(t *testing.T) {
	resourceName := "/Common/test-vcmp"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	// currentImage tracks the mock server's "stored" initial_image across
	// requests, so GET reflects whatever the most recent POST/PATCH set --
	// otherwise a static GET response would make every step look like a
	// perpetual diff (or hide the PATCH ever having been issued).
	currentImage := "12.1.2.iso"
	patchCalls := 0
	mux.HandleFunc("/mgmt/tm/vcmp/guest", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","initialImage":"%s"}`, resourceName, resourceName, currentImage)
	})
	mux.HandleFunc("/mgmt/tm/vcmp/guest/~Common~test-vcmp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			patchCalls++
			currentImage = "12.1.3.iso"
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","initialImage":"%s"}`, resourceName, resourceName, currentImage)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipVcmpGuestInitialImage(resourceName, "12.1.2.iso"),
			},
			{
				Config: testBigipVcmpGuestInitialImage(resourceName, "12.1.3.iso"),
			},
		},
	})

	// Explicitly confirm resourceBigipVcmpGuestUpdate issued exactly one
	// PATCH (step 2's config change), rather than relying solely on the
	// implicit "plan converges" check to catch a silently no-op Update.
	assert.Equal(t, 1, patchCalls, "expected exactly one PATCH request from the update in step 2")
}

func TestUnitBigipVcmpGuestUpdateError(t *testing.T) {
	resourceName := "/Common/test-vcmp"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/vcmp/guest", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","initialImage":"12.1.2.iso"}`, resourceName, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/vcmp/guest/~Common~test-vcmp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			http.Error(w, "update failed", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","initialImage":"12.1.2.iso"}`, resourceName, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipVcmpGuestInitialImage(resourceName, "12.1.2.iso"),
			},
			{
				Config:      testBigipVcmpGuestInitialImage(resourceName, "12.1.3.iso"),
				ExpectError: regexp.MustCompile(`HTTP 500 :: update failed`),
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
`, resourceName)
}

func testBigipVcmpGuestCreate(resourceName string) string {
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
`, resourceName)
}

func testBigipVcmpGuestModify(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_vcmp_guest" "test-vcmp" {
  name  = "%s"
  state = "configured"
}
`, resourceName)
}

func testBigipVcmpGuestBridgedNoAddress(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_vcmp_guest" "test-vcmp" {
  name         = "%s"
  mgmt_network = "bridged"
}
`, resourceName)
}

func testBigipVcmpGuestInitialImage(resourceName, image string) string {
	return fmt.Sprintf(`
resource "bigip_vcmp_guest" "test-vcmp" {
  name          = "%s"
  initial_image = "%s"
}
`, resourceName, image)
}
