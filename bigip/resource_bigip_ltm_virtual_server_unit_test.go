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

func TestAccBigipLtmVirtualServerUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-virtual-server"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmVirtualServerInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmVirtualServerUnitCreate(t *testing.T) {
	resourceName := "/Common/test-virtual-server"
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
	mux.HandleFunc("/mgmt/tm/ltm/virtual", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","description":"VirtualServer-test","destination":"/Common/10.255.255.254:9999","mask": "255.255.255.255","enabled":true,"fallbackPersistence":"/Common/dest_addr","ipProtocol":"tcp","pool":"","source":"0.0.0.0/0","sourceAddressTranslation":{"type":"automap"},"translateAddress":"enabled","translatePort":"enabled","vlansDisabled":true,"persist": [{"name": "hash","partition": "Common","tmDefault": "yes"}],"profiles":[{"name":"/Common/http","context":"all"},{"name":"/Common/tcp","context":"clientside"},{"name":"/Common/tcp-lan-optimized","context":"serverside"}],"policies":null}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-virtual-server", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","description":"VirtualServer-test","destination":"/Common/10.255.255.254:9999","mask": "255.255.255.255","enabled":true,"fallbackPersistence":"/Common/dest_addr","ipProtocol":"tcp","pool":"","source":"0.0.0.0/0","sourceAddressTranslation":{"type":"automap"},"translateAddress":"enabled","translatePort":"enabled","vlansDisabled":true,"persist": [{"name": "hash","partition": "Common","tmDefault": "yes"}],"profiles":[{"name":"/Common/http","context":"all"},{"name":"/Common/tcp","context":"clientside"},{"name":"/Common/tcp-lan-optimized","context":"serverside"}],"policies":null}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-virtual-server/profiles", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items": [{"name": "f5-tcp-progressive","partition": "Common","fullPath": "/Common/f5-tcp-progressive","context": "all"},{"name": "http","partition": "Common","fullPath": "/Common/http","context": "all"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-virtual-server/policies", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items": []}`)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-virtual-server", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","description":"VirtualServer-test","destination":"/Common/10.255.254.254:9999","mask": "255.255.255.255","enabled":true,"fallbackPersistence":"/Common/dest_addr","ipProtocol":"tcp","pool":"","source":"0.0.0.0/0","sourceAddressTranslation":{"type":"automap"},"translateAddress":"enabled","translatePort":"enabled","vlansDisabled":true,"persist": [{"name": "hash","partition": "Common","tmDefault": "yes"}],"profiles":[{"name":"/Common/http","context":"all"},{"name":"/Common/tcp","context":"clientside"},{"name":"/Common/tcp-lan-optimized","context":"serverside"}],"policies":null}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmVirtualServerCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmVirtualServerModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmVirtualServerUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-virtual-server"
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
	mux.HandleFunc("/mgmt/tm/ltm/virtual", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testvirtual##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmVirtualServerCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testvirtual##\\) is invalid"),
			},
		},
	})
}
func TestAccBigipLtmVirtualServerUnitReadError(t *testing.T) {
	resourceName := "/Common/test-virtual-server"
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
	mux.HandleFunc("/mgmt/tm/ltm/virtual", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","description":"VirtualServer-test","destination":"/Common/10.255.255.254:9999","mask": "255.255.255.255","enabled":true,"fallbackPersistence":"/Common/dest_addr","ipProtocol":"tcp","pool":"","source":"0.0.0.0/0","sourceAddressTranslation":{"type":"automap"},"translateAddress":"enabled","translatePort":"enabled","vlansDisabled":true,"persist": [{"name": "hash","partition": "Common","tmDefault": "yes"}],"profiles":[{"name":"/Common/http","context":"all"},{"name":"/Common/tcp","context":"clientside"},{"name":"/Common/tcp-lan-optimized","context":"serverside"}],"policies":null}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-virtual-server", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested Virtual server (/Common/test-virtual-server) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmVirtualServerCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Virtual server \\(/Common/test-virtual-server\\) was not found"),
			},
		},
	})
}

func testBigipLtmVirtualServerInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test-vs" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmVirtualServerCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test-vs" {
  name    = "%s"
  destination = "10.255.255.254"
  description = "VirtualServer-test"
  port = 9999
  mask = "255.255.255.255"
  source_address_translation = "automap"
  ip_protocol = "tcp"
  profiles = ["/Common/f5-tcp-progressive","/Common/http"]
  client_profiles = ["/Common/tcp"]
  server_profiles = ["/Common/tcp-lan-optimized"]
  default_persistence_profile = "/Common/hash"
  fallback_persistence_profile = "/Common/dest_addr"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmVirtualServerModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test-vs" {
  name    = "%s"
  destination = "10.255.254.254"
  description = "VirtualServer-test"
  port = 9999
  mask = "255.255.255.255"
  source_address_translation = "automap"
  ip_protocol = "tcp"
  profiles = ["/Common/f5-tcp-progressive","/Common/http"]
  client_profiles = ["/Common/tcp"]
  server_profiles = ["/Common/tcp-lan-optimized"]
  default_persistence_profile = "/Common/hash"
  fallback_persistence_profile = "/Common/dest_addr"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
