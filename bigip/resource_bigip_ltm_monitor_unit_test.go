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
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipLtmMonitorUnitInvalid(t *testing.T) {
	resourceName := "test-monitor"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmMonitorInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmMonitorUnitCreate(t *testing.T) {
	monitors := []string{"http", "https", "icmp", "gateway-icmp", "tcp", "tcp-half-open", "ftp", "udp", "postgresql", "mysql", "mssql", "ldap"}
	resourceName := "test-monitor"
	httpDefault := "/Common/http"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	for _, name := range monitors {
		mux.HandleFunc(fmt.Sprintf("/mgmt/tm/ltm/monitor/%s", name), func(w http.ResponseWriter, r *http.Request) {
			reqUrl := r.URL.String()
			if strings.HasSuffix(reqUrl, "http") {
				fmt.Fprintf(w, `{
					"items": [
						{
							"name":"/Common/%[1]s",
							"fullPath":"/Common/%[1]s",
							"defaultsFrom":"%[2]s"
						}
					]
				}`, resourceName, httpDefault)
			} else {
				fmt.Fprintf(w, `{"items":[]}`)
			}
		})
	}

	mux.HandleFunc(fmt.Sprintf("/mgmt/tm/ltm/monitor/http/~Common~%s", resourceName),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
			"name":"/Common/%[1]s",
			"fullPath":"/Common/%[1]s,
			"defaultsFrom": "%[2]s"
		}`, resourceName, httpDefault)
		})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmMonitorCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmMonitorModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmMonitorUnitCreateError(t *testing.T) {
	resourceName := "test-monitor"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})

	mux.HandleFunc("/mgmt/tm/ltm/monitor/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Create Page Not Found", 404)
	})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmMonitorCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: Create Page Not Found"),
			},
		},
	})
}

func testBigipLtmMonitorInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_monitor" "test-monitor" {
  name       = "/Common/%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmMonitorCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_monitor" "test-monitor" {
  name    = "/Common/%s"
  parent = "/Common/http"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmMonitorModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_monitor" "test-monitor" {
  name    = "/Common/%s"
  parent = "/Common/http"
  send = "GET /some/path\r\n"
  timeout = 999
  interval = 996
  receive = "HTTP 1.1 302 Found"
  receive_disable = "HTTP/1.1 429"
  reverse = "disabled"
  transparent = "disabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
