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

func TestAccBigipNetTunnelUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-net-tunnel"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetTunnelInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipNetTunnelUnitCreate(t *testing.T) {
	resourceName := "/Common/test-net-tunnel"
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
	mux.HandleFunc("/mgmt/tm/net/tunnels/tunnel", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","autoLasthop":"default","idleTimeout":300,"localAddress":"192.16.81.240","mode":"bidirectional","profile":"/Common/dslite","remoteAddress":"any6","secondaryAddress":"any6","tos":"preserve","transparent":"disabled","usePmtu":"enabled"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/net/tunnels/tunnel/~Common~test-net-tunnel", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","autoLasthop":"default","idleTimeout":300,"localAddress":"192.16.81.240","mode":"bidirectional","profile":"/Common/dslite","remoteAddress":"any6","secondaryAddress":"any6","tos":"preserve","transparent":"disabled","usePmtu":"enabled"}`, resourceName)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/tunnels/tunnel/~Common~test-net-tunnel", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","autoLasthop":"default","idleTimeout":301,"localAddress":"192.16.81.240","mode":"bidirectional","profile":"/Common/dslite","remoteAddress":"any6","secondaryAddress":"any6","tos":"preserve","transparent":"disabled","usePmtu":"enabled"}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipNetTunnelCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipNetTunnelModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipNetTunnelUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-net-tunnel"
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
	mux.HandleFunc("/mgmt/tm/net/tunnels/tunnel", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testnettunnel##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetTunnelCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testnettunnel##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipNetTunnelUnitReadError(t *testing.T) {
	resourceName := "/Common/test-net-tunnel"
	httpDefault := "/Common/http"
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
	mux.HandleFunc("/mgmt/tm/net/tunnels/tunnel", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, httpDefault)
	})
	mux.HandleFunc("/mgmt/tm/net/tunnels/tunnel/~Common~test-net-tunnel", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested Net Tunnel (/Common/test-net-tunnel) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetTunnelCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Net Tunnel \\(/Common/test-net-tunnel\\) was not found"),
			},
		},
	})
}

func testBigipNetTunnelInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_net_tunnel" "test_tunnel" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipNetTunnelCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_net_tunnel" "test_tunnel" {
  name    = "%s"
  auto_last_hop     = "default"
  idle_timeout      = 300
  key               = 0
  local_address     = "192.16.81.240"
  mode              = "bidirectional"
  mtu               = 0
  profile           = "/Common/dslite"
  remote_address    = "any6"
  secondary_address = "any6"
  tos               = "preserve"
  transparent       = "disabled"
  use_pmtu          = "enabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipNetTunnelModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_net_tunnel" "test_tunnel" {
  name    = "%s"
  auto_last_hop     = "default"
  idle_timeout      = 301
  key               = 0
  local_address     = "192.16.81.240"
  mode              = "bidirectional"
  mtu               = 0
  profile           = "/Common/dslite"
  remote_address    = "any6"
  secondary_address = "any6"
  tos               = "preserve"
  transparent       = "disabled"
  use_pmtu          = "enabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
