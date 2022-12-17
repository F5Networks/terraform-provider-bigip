/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipLtmDatagroupUnitInvalid(t *testing.T) {
	resourceName := "test_dg"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmDatagroupInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmDatagroupUnitCreate(t *testing.T) {
	resourceName := "test_dg"
	cases := []struct {
		testHeading string
		internal    bool
	}{
		{"Creating internal Data Group", true},
		{"Creating external Data Group", false},
	}

	for _, tc := range cases {
		t.Run(tc.testHeading, func(t *testing.T) {
			setup()
			defer teardown()
			registerHandlers(t, resourceName)

			createCfg := testBigipLtmDatagroupCreateExternal(resourceName)
			modifyCfg := testBigipLtmDatagroupModifyExternal(resourceName)
			if tc.internal {
				createCfg = testBigipLtmDatagroupCreateInternal(resourceName)
				modifyCfg = testBigipLtmDatagroupModifyInternal(resourceName)
			}

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testProviders,
				Steps: []resource.TestStep{
					{
						Config: createCfg,
					},
					{
						Config:             modifyCfg,
						ExpectNonEmptyPlan: tc.internal,
					},
				},
			})
		})
	}
}

func registerHandlers(t *testing.T, resourceName string) {
	mux.HandleFunc("/mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/ltm/data-group/internal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var dg_json bigip.DataGroup
			json.NewDecoder(r.Body).Decode(&dg_json)
			log.Printf("Request body: %+v", dg_json)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/data-group/external", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var dg_json bigip.DataGroup
			json.NewDecoder(r.Body).Decode(&dg_json)
			log.Printf("Request body: %+v", dg_json)
		}
	})
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		version := `{"entries":{"https://localhost/mgmt/tm/cli/version/0":{"nestedStats":{"entries":{"active":{"description": "16.1.2.1"}}}}}}`
		io.WriteString(w, version)
	})
	mux.HandleFunc(fmt.Sprintf("/mgmt/tm/ltm/data-group/internal/~Common~%s", resourceName),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w,
				`{
					"name":"/Common/%[1]s",
					"fullPath":"/Common/%[1]s",
					"type":"string",
					"records":[{"name":"a","data":"1"},{"name":"b","data":"2"}]
				}`, resourceName,
			)
		},
	)
	mux.HandleFunc(fmt.Sprintf("/mgmt/shared/file-transfer/uploads/%s", resourceName),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
			assert.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
			fmt.Fprintf(w, "{}")
		},
	)
	mux.HandleFunc("/mgmt/tm/sys/file/data-group", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
	})
	mux.HandleFunc(fmt.Sprintf("/mgmt/tm/sys/file/data-group/~Common~%s", resourceName),
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" && r.Method != "POST" {
				assert.Fail(t, "Expected method 'POST' or 'DELETE', got %s", r.Method)
			}
		},
	)
	mux.HandleFunc(fmt.Sprintf("/mgmt/tm/ltm/data-group/external/~Common~%s", resourceName),
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				fmt.Fprintf(w,
					`{
					"name":"/Common/%[1]s",
					"fullPath":"/Common/%[1]s",
					"type":"string",
					"externalFileName":"/Common/records.txt"
				}`, resourceName,
				)
			}
			if r.Method == "PATCH" {
				fmt.Fprintf(w,
					`{
					"name":"/Common/%[1]s",
					"fullPath":"/Common/%[1]s",
					"type":"string",
					"externalFileName":"/Common/records2.txt"
				}`, resourceName,
				)
			}
		},
	)
}

func testBigipLtmDatagroupInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_datagroup" "test_dg" {
  name       = "/Common/%s"
  type       = "string"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmDatagroupCreateInternal(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_datagroup" "test_dg" {
  name    = "/Common/%s"
  type    = "string"
  internal = true
  record {
	name = "a"
    data = "1"
  }
  record {
	name = "b"
    data = "2"
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, server.URL)
}

func testBigipLtmDatagroupModifyInternal(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_datagroup" "test_dg" {
  name    = "/Common/%s"
  type    = "string"
  internal = true
  record {
	name = "a"
    data = "1"
  }
  record {
	name = "b"
    data = "3"
  }
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, server.URL)
}

func testBigipLtmDatagroupCreateExternal(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_datagroup" "test_dg" {
  name    = "/Common/%s"
  type    = "string"
  internal = false
  records_src = "/tmp/records.txt"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, server.URL)
}

func testBigipLtmDatagroupModifyExternal(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_datagroup" "test_dg" {
  name    = "/Common/%s"
  type    = "string"
  internal = false
  records_src = "/tmp/records2.txt"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, server.URL)
}
