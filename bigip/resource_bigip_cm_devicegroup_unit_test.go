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

func TestBigipCmDeviceGroupUnitInvalid(t *testing.T) {
	resourceName := "test-devicegroup"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipCmDeviceGroupInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestBigipCmDeviceGroupUnitCreate(t *testing.T) {
	resourceName := "test-devicegroup"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/tm/cm/device-group/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" && r.Method != "DELETE" {
			assert.Fail(t, `request method is supposed to be "GET", "POST" or "DELETE", got %s`, r.Method)
		}
	})

	mux.HandleFunc(fmt.Sprintf("/mgmt/tm/cm/device-group/~Common~%s", resourceName),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"name":"/Common/%[1]s",
				"fullPath":"/Common/%[1]s",
				"type": "sync-only",
				"autoSync": "disabled",
				"description": "test description",
				"devicesReference": {
					"items": [
						{
							"name": "bigip15.com"
						}
					]
				}
			}`, resourceName)
		},
	)

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipCmDeviceGroupCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipCmDeviceGroupModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testBigipCmDeviceGroupInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_cm_devicegroup" "test-devicegroup" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipCmDeviceGroupCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_cm_devicegroup" "test-devicegroup" {
  name        = "/Common/%s"
  description = "test description"
  device {
	name = "bigip15.com"
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipCmDeviceGroupModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_cm_devicegroup" "test-devicegroup" {
  name               = "/Common/%s"
  description        = "test description 2"
  full_load_on_sync  = "false"
  incremental_config = 1444
  device {
	name = "bigip16.com"
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
