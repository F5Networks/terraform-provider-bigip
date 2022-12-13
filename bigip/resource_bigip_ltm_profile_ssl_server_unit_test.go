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

func TestAccBigipLtmProfileSslServerUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-server-ssl"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileSslServerInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfileSslServerUnitCreate(t *testing.T) {
	resourceName := "/Common/test-profile-server-ssl"
	serversslDefault := "/Common/serverssl"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","alertTimeout": "indefinite","allowExpiredCrl": "disabled",
"appService": "none",
"authenticate": "once",
"authenticateDepth": 9,
"authenticateName": "none",
"bypassOnClientCertFail": "disabled",
"bypassOnHandshakeAlert": "disabled",
"c3dCaCert": "none",
"c3dCaKey": "none",
"c3dCertExtensionCustomOids": [],
"c3dCertExtensionIncludes": [
"basic-constraints",
"extended-key-usage",
"key-usage",
"subject-alternative-name"
],
"c3dCertLifespan": 24,
"caFile": "none",
"cacheSize": 262144,
"cacheTimeout": 3600,
"cert": "/Common/default.crt",
"key": "/Common/default.key",
"chain": "none",
"cipherGroup": "none",
"ciphers": "DEFAULT",
"crl": "none",
"crlFile": "none",
"data_0rtt": "disabled",
"description": "none",
"expireCertResponseControl": "drop",
"genericAlert": "enabled",
"handshakeTimeout": "10",
"key": "/Common/default.key",
"logPublisher": "/Common/sys-ssl-publisher",
"maxActiveHandshakes": "indefinite",
"modSslMethods": "disabled",
"mode": "enabled",
"ocsp": "none",
"tmOptions": "{ dont-insert-empty-fragments no-tlsv1.3 no-dtlsv1.2 }",
"peerCertMode": "ignore",
"proxySsl": "disabled",
"proxySslPassthrough": "disabled",
"renegotiatePeriod": "indefinite",
"renegotiateSize": "indefinite",
"renegotiation": "enabled",
"retainCertificate": "true",
"revokedCertStatusResponseControl": "drop",
"secureRenegotiation": "require-strict",
"serverName": "none",
"sessionMirroring": "disabled",
"sessionTicket": "disabled",
"sniDefault": "false",
"sniRequire": "false",
"sslC3d": "disabled",
"sslForwardProxy": "disabled",
"sslForwardProxyBypass": "disabled",
"sslForwardProxyVerifiedHandshake": "disabled",
"sslSignHash": "any",
"strictResume": "disabled",
"uncleanShutdown": "enabled",
"unknownCertStatusResponseControl": "ignore",
"untrustedCertResponseControl": "drop"}`, resourceName, serversslDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/~Common~test-profile-server-ssl", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","alertTimeout": "indefinite",
"allowExpiredCrl": "disabled",
"appService": "none",
"authenticate": "once",
"authenticateDepth": 9,
"authenticateName": "none",
"bypassOnClientCertFail": "disabled",
"bypassOnHandshakeAlert": "disabled",
"c3dCaCert": "none",
"c3dCaKey": "none",
"c3dCertExtensionCustomOids": [],
"c3dCertExtensionIncludes": [
"basic-constraints",
"extended-key-usage",
"key-usage",
"subject-alternative-name"
],
"c3dCertLifespan": 24,
"caFile": "none",
"cacheSize": 262144,
"cacheTimeout": 3600,
"cert": "/Common/default.crt",
"key": "/Common/default.key",
"chain": "none",
"cipherGroup": "none",
"ciphers": "DEFAULT",
"crl": "none",
"crlFile": "none",
"data_0rtt": "disabled",
"description": "none",
"expireCertResponseControl": "drop",
"genericAlert": "enabled",
"handshakeTimeout": "10",
"logPublisher": "/Common/sys-ssl-publisher",
"maxActiveHandshakes": "indefinite",
"modSslMethods": "disabled",
"mode": "enabled",
"ocsp": "none",
"tmOptions": "{ dont-insert-empty-fragments no-tlsv1.3 no-dtlsv1.2 }",
"peerCertMode": "ignore",
"proxySsl": "disabled",
"proxySslPassthrough": "disabled",
"renegotiatePeriod": "indefinite",
"renegotiateSize": "indefinite",
"renegotiation": "enabled",
"retainCertificate": "true",
"revokedCertStatusResponseControl": "drop",
"secureRenegotiation": "require-strict",
"serverName": "none",
"sessionMirroring": "disabled",
"sessionTicket": "disabled",
"sniDefault": "false",
"sniRequire": "false",
"sslC3d": "disabled",
"sslForwardProxy": "disabled",
"sslForwardProxyBypass": "disabled",
"sslForwardProxyVerifiedHandshake": "disabled",
"sslSignHash": "any",
"strictResume": "disabled",
"uncleanShutdown": "enabled",
"unknownCertStatusResponseControl": "ignore",
"untrustedCertResponseControl": "drop"}`, resourceName, serversslDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/~Common~test-profile-server-ssl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","alertTimeout": "indefinite",
"allowExpiredCrl": "disabled",
"appService": "none",
"authenticate": "always",
"authenticateDepth": 9,
"authenticateName": "none",
"bypassOnClientCertFail": "disabled",
"bypassOnHandshakeAlert": "disabled",
"c3dCaCert": "none",
"c3dCaKey": "none",
"c3dCertExtensionCustomOids": [],
"c3dCertExtensionIncludes": [
"basic-constraints",
"extended-key-usage",
"key-usage",
"subject-alternative-name"
],
"c3dCertLifespan": 24,
"caFile": "none",
"cacheSize": 262144,
"cacheTimeout": 3600,
"cert": "/Common/default.crt",
"key": "/Common/default.key",
"chain": "none",
"cipherGroup": "none",
"ciphers": "DEFAULT",
"crl": "none",
"crlFile": "none",
"data_0rtt": "disabled",
"description": "none",
"expireCertResponseControl": "drop",
"genericAlert": "enabled",
"handshakeTimeout": "10",
"logPublisher": "/Common/sys-ssl-publisher",
"maxActiveHandshakes": "indefinite",
"modSslMethods": "disabled",
"mode": "enabled",
"ocsp": "none",
"tmOptions": "{ dont-insert-empty-fragments no-tlsv1.3 no-dtlsv1.2 }",
"peerCertMode": "ignore",
"proxySsl": "disabled",
"proxySslPassthrough": "disabled",
"renegotiatePeriod": "indefinite",
"renegotiateSize": "indefinite",
"renegotiation": "enabled",
"retainCertificate": "true",
"revokedCertStatusResponseControl": "drop",
"secureRenegotiation": "require-strict",
"serverName": "none",
"sessionMirroring": "disabled",
"sessionTicket": "disabled",
"sniDefault": "false",
"sniRequire": "false",
"sslC3d": "disabled",
"sslForwardProxy": "disabled",
"sslForwardProxyBypass": "disabled",
"sslForwardProxyVerifiedHandshake": "disabled",
"sslSignHash": "any",
"strictResume": "disabled",
"uncleanShutdown": "enabled",
"unknownCertStatusResponseControl": "ignore",
"untrustedCertResponseControl": "drop"}`, resourceName, serversslDefault)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfileSslServerCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfileSslServerModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfileSslServerUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-profile-server-ssl"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testserverssl##) is invalid", http.StatusBadRequest)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileSslServerCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testserverssl##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmProfileSslServerUnitReadError(t *testing.T) {
	resourceName := "/Common/test-profile-server-ssl"
	serversslDefault := "/Common/serverssl"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, serversslDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/~Common~test-profile-server-ssl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested Server-SSL Profile (/Common/test-profile-server-ssl) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileSslServerCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Server-SSL Profile \\(/Common/test-profile-server-ssl\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfileSslServerInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_server_ssl" "test-profile-server-ssl" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfileSslServerCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_server_ssl" "test-profile-server-ssl" {
  name    = "%s"
  defaults_from = "/Common/serverssl"
  authenticate  = "once"
  ciphers       = "DEFAULT"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmProfileSslServerModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_server_ssl" "test-profile-server-ssl" {
  name    = "%s"
  defaults_from = "/Common/serverssl"
  authenticate  = "always"
  ciphers       = "DEFAULT"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
