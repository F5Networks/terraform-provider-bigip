/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipCmDeviceUnitInvalid(t *testing.T) {
	resourceName := "test-device"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipCmDeviceInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipCmDeviceUnitCreate(t *testing.T) {
	resourceName := "test-device"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/tm/cm/device", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("------------Request Method: %s-------------", r.Method)
	})
	mux.HandleFunc(fmt.Sprintf("/mgmt/tm/cm/device/~Common~%s", resourceName), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			`{
				"name": "/Common/%[1]s",
				"configsyncIp": "2.2.2.2",
				"fullPath": "/Common/%[1]s"
			}`,
			resourceName,
		)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipCmDeviceCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipCmDeviceModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipCmDeviceUnitCreateError(t *testing.T) {
	resourceName := "test-device"
	setup()

	mux.HandleFunc("/mgmt/tm/cm/device", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Create Page Not Found", 404)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipCmDeviceCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile(`HTTP 404 :: Create Page Not Found`),
			},
		},
	})
}

func testBigipCmDeviceInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_cm_device" "test-device" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipCmDeviceCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_cm_device" "test-device" {
  name    = "/Common/%s"
  configsync_ip = "2.2.2.2"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipCmDeviceModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_cm_device" "test-device" {
  name    = "/Common/%s"
  configsync_ip = "2.2.2.3"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
