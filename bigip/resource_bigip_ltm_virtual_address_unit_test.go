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

func TestAccBigipLtmVirtualAddressUnitInvalid(t *testing.T) {
	resourceName := "/Common/10.1.2.13"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmVirtualAddressInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmVirtualAddressUnitCreate(t *testing.T) {
	resourceName := "/Common/10.1.2.13"
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
	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "/Common/10.1.2.13","arp": "enabled","autoDelete": "true","connectionLimit": 0,"enabled": "yes","floating": "enabled","icmpEcho": "enabled","inheritedTrafficGroup": "true","mask": "255.255.255.255","routeAdvertisement": "disabled","serverScope": "any","spanning": "disabled","trafficGroup": "/Common/traffic-group-1"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/virtual-address/~Common~10.1.2.13", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "/Common/10.1.2.13","arp": "enabled","autoDelete": "true","connectionLimit": 0,"enabled": "yes","floating": "enabled","icmpEcho": "enabled","inheritedTrafficGroup": "true","mask": "255.255.255.255","routeAdvertisement": "disabled","serverScope": "any","spanning": "disabled","trafficGroup": "/Common/traffic-group-1"}`, resourceName)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/virtual-address/~Common~10.1.2.13", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		//if r.Method == "GET" {
		//	_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "/Common/10.1.2.13","arp": "enabled","autoDelete": "true","connectionLimit": 0,"enabled": "yes","floating": "enabled","icmpEcho": "enabled","inheritedTrafficGroup": "true","mask": "255.255.255.255","routeAdvertisement": "disabled","serverScope": "any","spanning": "disabled","trafficGroup": "/Common/traffic-group-1"}`, resourceName)
		//}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "/Common/10.1.2.13","arp": "enabled","autoDelete": "true","connectionLimit": 1,"enabled": "yes","floating": "enabled","icmpEcho": "enabled","inheritedTrafficGroup": "true","mask": "255.255.255.255","routeAdvertisement": "disabled","serverScope": "any","spanning": "disabled","trafficGroup": "/Common/traffic-group-1"}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmVirtualAddressCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmVirtualAddressModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmVirtualAddressUnitCreateError(t *testing.T) {
	resourceName := "/Common/10.1.2.13"
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
	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testvirtualaddress##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmVirtualAddressCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testvirtualaddress##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmVirtualAddressUnitReadError(t *testing.T) {
	resourceName := "/Common/10.1.2.13"
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
	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath": "/Common/10.1.2.13","arp": "enabled","autoDelete": "true","connectionLimit": 0,"enabled": "yes","floating": "enabled","icmpEcho": "enabled","inheritedTrafficGroup": "true","mask": "255.255.255.255","routeAdvertisement": "disabled","serverScope": "any","spanning": "disabled","trafficGroup": "/Common/traffic-group-1"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/virtual-address/~Common~10.1.2.13", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested Virtual-Address (/Common/10.1.2.13) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmVirtualAddressCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Virtual-Address \\(/Common/10.1.2.13\\) was not found"),
			},
		},
	})
}

func testBigipLtmVirtualAddressInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_address" "test-va" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmVirtualAddressCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_address" "test-va" {
  name    = "%s"
  arp  = true
  conn_limit  = 0
  icmp_echo  = "enabled"
  advertize_route  = "disabled"
  traffic_group= "/Common/traffic-group-1"
  auto_delete = true
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmVirtualAddressModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_address" "test-va" {
  name    = "%s"
  arp  = true
  conn_limit  = 1
  icmp_echo  = "enabled"
  advertize_route  = "disabled"
  traffic_group= "/Common/traffic-group-1"
  auto_delete = true
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
