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

func TestAccBigipNetIkepeerUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-ike-peer"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetIkepeerInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipNetIkepeerUnitCreate(t *testing.T) {
	resourceName := "/Common/test-ike-peer"
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
	mux.HandleFunc("/mgmt/tm/net/ipsec/ike-peer", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","dpdDelay":30,"lifetime":1440,"generatePolicy":"off","mode":"main","myCertFile":"/Common/default.crt","myCertKeyFile":"/Common/default.key","myIdType":"address","natTraversal":"off","passive":"false","peersCertType":"none","peersIdType":"address","phase1AuthMethod":"rsa-signature","phase1EncryptAlgorithm":"3des","phase1HashAlgorithm":"sha256","phase1PerfectForwardSecrecy":"modp1024","prf":"sha256","proxySupport":"enabled","remoteAddress":"1.5.3.4","replayWindowSize":64,"state":"enabled","verifyCert":"false"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/net/ipsec/ike-peer/~Common~test-ike-peer", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","dpdDelay":30,"lifetime":1440,"generatePolicy":"off","mode":"main","myCertFile":"/Common/default.crt","myCertKeyFile":"/Common/default.key","myIdType":"address","natTraversal":"off","passive":"false","peersCertType":"none","peersIdType":"address","phase1AuthMethod":"rsa-signature","phase1EncryptAlgorithm":"3des","phase1HashAlgorithm":"sha256","phase1PerfectForwardSecrecy":"modp1024","prf":"sha256","proxySupport":"enabled","remoteAddress":"1.5.3.4","replayWindowSize":64,"state":"enabled","verifyCert":"false"}`, resourceName)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/ipsec/ike-peer/~Common~test-ike-peer", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","dpdDelay":32,"lifetime":1440,"generatePolicy":"off","mode":"main","myCertFile":"/Common/default.crt","myCertKeyFile":"/Common/default.key","myIdType":"address","natTraversal":"off","passive":"false","peersCertType":"none","peersIdType":"address","phase1AuthMethod":"rsa-signature","phase1EncryptAlgorithm":"3des","phase1HashAlgorithm":"sha256","phase1PerfectForwardSecrecy":"modp1024","prf":"sha256","proxySupport":"enabled","remoteAddress":"1.5.3.4","replayWindowSize":64,"state":"enabled","verifyCert":"false"}`, resourceName)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipNetIkepeerCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipNetIkepeerModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipNetIkepeerUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-ike-peer"
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
	mux.HandleFunc("/mgmt/tm/net/ipsec/ike-peer", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testikepeer##) is invalid", http.StatusBadRequest)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetIkepeerCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testikepeer##\\) is invalid"),
			},
		},
	})
}
func TestAccBigipNetIkepeerUnitReadError(t *testing.T) {
	resourceName := "/Common/test-ike-peer"
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
	mux.HandleFunc("/mgmt/tm/net/ipsec/ike-peer", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","dpdDelay":30,"lifetime":1440,"generatePolicy":"off","mode":"main","myCertFile":"/Common/default.crt","myCertKeyFile":"/Common/default.key","myIdType":"address","natTraversal":"off","passive":"false","peersCertType":"none","peersIdType":"address","phase1AuthMethod":"rsa-signature","phase1EncryptAlgorithm":"3des","phase1HashAlgorithm":"sha256","phase1PerfectForwardSecrecy":"modp1024","prf":"sha256","proxySupport":"enabled","remoteAddress":"1.5.3.4","replayWindowSize":64,"state":"enabled","verifyCert":"false"}`, resourceName)
	})
	mux.HandleFunc("/mgmt/tm/net/ipsec/ike-peer/~Common~test-ike-peer", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested IKE Peer (/Common/test-ike-peer) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipNetIkepeerCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested IKE Peer \\(/Common/test-ike-peer\\) was not found"),
			},
		},
	})
}

func testBigipNetIkepeerInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_net_ike_peer"  "test_ike_peer" {
  name       = "%s"
  remote_address                 = "1.5.3.4"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipNetIkepeerCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_net_ike_peer"  "test_ike_peer" {
  name    = "%s"
  dpd_delay                      = 30
  generate_policy                = "off"
  lifetime                       = 1440
  mode                           = "main"
  my_cert_file                   = "/Common/default.crt"
  my_cert_key_file               = "/Common/default.key"
  my_id_type                     = "address"
  nat_traversal                  = "off"
  passive                        = "false"
  peers_cert_type                = "none"
  peers_id_type                  = "address"
  phase1_auth_method             = "rsa-signature"
  phase1_encrypt_algorithm       = "3des"
  phase1_hash_algorithm          = "sha256"
  phase1_perfect_forward_secrecy = "modp1024"
  prf                            = "sha256"
  proxy_support                  = "enabled"
  remote_address                 = "1.5.3.4"
  replay_window_size             = 64
  state                          = "enabled"
  verify_cert                    = "false"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipNetIkepeerModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_net_ike_peer"  "test_ike_peer" {
  name    = "%s"
  dpd_delay                      = 32
  generate_policy                = "off"
  lifetime                       = 1440
  mode                           = "main"
  my_cert_file                   = "/Common/default.crt"
  my_cert_key_file               = "/Common/default.key"
  my_id_type                     = "address"
  nat_traversal                  = "off"
  passive                        = "false"
  peers_cert_type                = "none"
  peers_id_type                  = "address"
  phase1_auth_method             = "rsa-signature"
  phase1_encrypt_algorithm       = "3des"
  phase1_hash_algorithm          = "sha256"
  phase1_perfect_forward_secrecy = "modp1024"
  prf                            = "sha256"
  proxy_support                  = "enabled"
  remote_address                 = "1.5.3.4"
  replay_window_size             = 64
  state                          = "enabled"
  verify_cert                    = "false"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
