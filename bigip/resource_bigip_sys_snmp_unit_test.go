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

func TestAccBigipSysSNMPUnitInvalid(t *testing.T) {
	resourceName := "NetOPsAdmin@f5.com"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysSnmpInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipSysSNMPUnitCreate(t *testing.T) {
	resourceName := "NetOPsAdmin@f5.com"
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
	mux.HandleFunc("/mgmt/tm/sys/snmp", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"sysContact":"%s","sysLocation":"SeattleHQ","allowedAddresses":["202.10.10.2"]}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipSysSnmpCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipSysSnmpModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipSysSNMPUnitReadError(t *testing.T) {
	resourceName := "NetOPsAdmin@f5.com"
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
	mux.HandleFunc("/mgmt/tm/sys/snmp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.Error(w, "The requested SNMP (/Common/test-profile-http) was not found", http.StatusNotFound)
		}
		_, _ = fmt.Fprintf(w, `{"sysContact":"%s","sysLocation":"SeattleHQ","allowedAddresses":["202.10.10.2"]}`, resourceName)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysSnmpCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested SNMP \\(/Common/test-profile-http\\) was not found"),
			},
		},
	})
}

func TestAccBigipSysSNMPUnitCreateError(t *testing.T) {
	resourceName := "NetOPsAdmin@f5.com"
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
	mux.HandleFunc("/mgmt/tm/sys/snmp", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "The requested object name (/Common/testsnmp##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysSnmpCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testsnmp##\\) is invalid"),
			},
		},
	})
}

func testBigipSysSnmpInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_snmp" "test-snmp" {
  sys_contact       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipSysSnmpCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_snmp" "test-snmp" {
  sys_contact = "%s"
  sys_location = "SeattleHQ"
  allowedaddresses = ["202.10.10.2"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipSysSnmpModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_snmp" "test-snmp" {
  sys_contact = "%s"
  sys_location = "SeattleHQ"
  allowedaddresses = ["202.10.10.22"]
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
