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

func TestAccBigipSysSNMPTrapsUnitInvalid(t *testing.T) {
	resourceName := "snmptraps"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysSnmpTrapsInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipSysSNMPTrapsUnitCreate(t *testing.T) {
	resourceName := "snmptraps"
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
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","community":"f5community","description":"Setup snmp traps","host":"195.10.10.1","port":111}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/snmptraps", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","community":"f5community","description":"Setup snmp traps","host":"195.10.10.1","port":111}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipSysSnmpTrapsCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipSysSnmpTrapsModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipSysSNMPTrapsUnitReadError(t *testing.T) {
	resourceName := "snmptraps"
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
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","community":"f5community","description":"Setup snmp traps","host":"195.10.10.1","port":111}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/snmptraps", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "The requested SNMP (/Common/test-snmp-trap) was not found", http.StatusNotFound)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysSnmpTrapsCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested SNMP \\(/Common/test-snmp-trap\\) was not found"),
			},
		},
	})
}

func TestAccBigipSysSNMPTrapsUnitCreateError(t *testing.T) {
	resourceName := "snmptraps"
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
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "The requested object name (/Common/testsnmp##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysSnmpTrapsCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testsnmp##\\) is invalid"),
			},
		},
	})
}

func testBigipSysSnmpTrapsInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_sys_snmp_traps" "test-snmp-trap" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipSysSnmpTrapsCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_snmp_traps" "test-snmp-trap" {
  name = "%s"
  community   = "f5community"
  host        = "195.10.10.1"
  description = "Setup snmp traps"
  port        = 111
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipSysSnmpTrapsModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_sys_snmp_traps" "test-snmp-trap" {
  name = "%s"
  community   = "f5community"
  host        = "195.10.11.1"
  description = "Setup snmp traps"
  port        = 111
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
