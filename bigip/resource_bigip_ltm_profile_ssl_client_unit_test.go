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

func TestAccBigipLtmProfileSslClientUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-clientSsl"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileSslClientInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfileSslClientUnitCreate(t *testing.T) {
	resourceName := "/Common/test-profile-clientSsl"
	clientsslDefault := "/Common/clientssl"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `
{"name":"%s","defaultsFrom":"%s","alertTimeout": "indefinite",
"allowDynamicRecordSizing": "disabled",
"allowExpiredCrl": "disabled",
"allowNonSsl": "disabled",
"appService": "none",
"authenticate": "once",
"authenticateDepth": 9,
"bypassOnClientCertFail": "disabled",
"bypassOnHandshakeAlert": "disabled",
"c3dClientFallbackCert": "none",
"c3dDropUnknownOcspStatus": "drop",
"c3dOcsp": "none",
"caFile": "none",
"cacheSize": 262144,
"cacheTimeout": 3600,
"cert": "/Common/default.crt",
"key": "/Common/default.key",
"certExtensionIncludes": [
"basic-constraints",
"subject-alternative-name"
],
"certLifespan": 30,
"certLookupByIpaddrPort": "disabled",
"chain": "none",
"cipherGroup": "none",
"ciphers": "DEFAULT",
"clientCertCa": "none",
"crl": "none",
"crlFile": "none",
"data_0rtt": "disabled",
"description": "none",
"destinationIpBlacklist": "none",
"destinationIpWhitelist": "none",
"forwardProxyBypassDefaultAction": "intercept",
"genericAlert": "enabled",
"handshakeTimeout": "10",
"helloExtensionIncludes": [],
"hostnameBlacklist": "none",
"hostnameWhitelist": "none",
"inheritCaCertkeychain": "false",
"inheritCertkeychain": "false",
"logPublisher": "/Common/sys-ssl-publisher",
"maxActiveHandshakes": "indefinite",
"maxAggregateRenegotiationPerMinute": "indefinite",
"maxRenegotiationsPerMinute": 5,
"maximumRecordSize": 16384,
"modSslMethods": "disabled",
"mode": "enabled",
"notifyCertStatusToVirtualServer": "disabled",
"ocspStapling": "disabled",
"tmOptions": "{ dont-insert-empty-fragments no-tlsv1.3 no-dtlsv1.2 }",
"peerCertMode": "ignore",
"peerNoRenegotiateTimeout": "10",
"proxyCaCert": "none",
"proxyCaKey": "none",
"proxySsl": "disabled",
"proxySslPassthrough": "disabled",
"renegotiateMaxRecordDelay": "indefinite",
"renegotiatePeriod": "indefinite",
"renegotiateSize": "indefinite",
"renegotiation": "enabled",
"retainCertificate": "true",
"secureRenegotiation": "require",
"serverName": "none",
"sessionMirroring": "disabled",
"sessionTicket": "disabled",
"sessionTicketTimeout": 0,
"sniDefault": "false",
"sniRequire": "false",
"sourceIpBlacklist": "none",
"sourceIpWhitelist": "none",
"sslC3d": "disabled",
"sslForwardProxy": "disabled",
"sslForwardProxyBypass": "disabled",
"sslForwardProxyVerifiedHandshake": "disabled",
"sslSignHash": "any",
"strictResume": "disabled",
"uncleanShutdown": "enabled",
"certKeyChain": [
{
"name": "default",
"appService": "none",
"cert": "/Common/default.crt",
"chain": "none",
"key": "/Common/default.key",
"usage": "SERVER"
}
]
}`, resourceName, clientsslDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/~Common~test-profile-clientSsl", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","alertTimeout": "indefinite",
"allowDynamicRecordSizing": "disabled",
"allowExpiredCrl": "disabled",
"allowNonSsl": "disabled",
"appService": "none",
"authenticate": "once",
"authenticateDepth": 9,
"bypassOnClientCertFail": "disabled",
"bypassOnHandshakeAlert": "disabled",
"c3dClientFallbackCert": "none",
"c3dDropUnknownOcspStatus": "drop",
"c3dOcsp": "none",
"caFile": "none",
"cacheSize": 262144,
"cacheTimeout": 3600,
"cert": "/Common/default.crt",
"key": "/Common/default.key",
"certExtensionIncludes": [
"basic-constraints",
"subject-alternative-name"
],
"certLifespan": 30,
"certLookupByIpaddrPort": "disabled",
"chain": "none",
"cipherGroup": "none",
"ciphers": "DEFAULT",
"clientCertCa": "none",
"crl": "none",
"crlFile": "none",
"data_0rtt": "disabled",
"description": "none",
"destinationIpBlacklist": "none",
"destinationIpWhitelist": "none",
"forwardProxyBypassDefaultAction": "intercept",
"genericAlert": "enabled",
"handshakeTimeout": "10",
"helloExtensionIncludes": [],
"hostnameBlacklist": "none",
"hostnameWhitelist": "none",
"inheritCaCertkeychain": "false",
"inheritCertkeychain": "false",
"logPublisher": "/Common/sys-ssl-publisher",
"maxActiveHandshakes": "indefinite",
"maxAggregateRenegotiationPerMinute": "indefinite",
"maxRenegotiationsPerMinute": 5,
"maximumRecordSize": 16384,
"modSslMethods": "disabled",
"mode": "enabled",
"notifyCertStatusToVirtualServer": "disabled",
"ocspStapling": "disabled",
"tmOptions": "{ dont-insert-empty-fragments no-tlsv1.3 no-dtlsv1.2 }",
"peerCertMode": "ignore",
"peerNoRenegotiateTimeout": "10",
"proxyCaCert": "none",
"proxyCaKey": "none",
"proxySsl": "disabled",
"proxySslPassthrough": "disabled",
"renegotiateMaxRecordDelay": "indefinite",
"renegotiatePeriod": "indefinite",
"renegotiateSize": "indefinite",
"renegotiation": "enabled",
"retainCertificate": "true",
"secureRenegotiation": "require",
"serverName": "none",
"sessionMirroring": "disabled",
"sessionTicket": "disabled",
"sessionTicketTimeout": 0,
"sniDefault": "false",
"sniRequire": "false",
"sourceIpBlacklist": "none",
"sourceIpWhitelist": "none",
"sslC3d": "disabled",
"sslForwardProxy": "disabled",
"sslForwardProxyBypass": "disabled",
"sslForwardProxyVerifiedHandshake": "disabled",
"sslSignHash": "any",
"strictResume": "disabled",
"uncleanShutdown": "enabled",
"certKeyChain": [
{
"name": "default",
"appService": "none",
"cert": "/Common/default.crt",
"chain": "none",
"key": "/Common/default.key",
"usage": "SERVER"
}
]
}`, resourceName, clientsslDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/~Common~test-profile-clientSsl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","alertTimeout": "indefinite",
"allowDynamicRecordSizing": "disabled",
"allowExpiredCrl": "disabled",
"allowNonSsl": "disabled",
"appService": "none",
"authenticate": "always",
"authenticateDepth": 9,
"bypassOnClientCertFail": "disabled",
"bypassOnHandshakeAlert": "disabled",
"c3dClientFallbackCert": "none",
"c3dDropUnknownOcspStatus": "drop",
"c3dOcsp": "none",
"caFile": "none",
"cacheSize": 262144,
"cacheTimeout": 3600,
"cert": "/Common/default.crt",
"certExtensionIncludes": [
"basic-constraints",
"subject-alternative-name"
],
"certLifespan": 30,
"certLookupByIpaddrPort": "disabled",
"chain": "none",
"cipherGroup": "none",
"ciphers": "DEFAULT",
"clientCertCa": "none",
"crl": "none",
"crlFile": "none",
"data_0rtt": "disabled",
"description": "none",
"destinationIpBlacklist": "none",
"destinationIpWhitelist": "none",
"forwardProxyBypassDefaultAction": "intercept",
"genericAlert": "enabled",
"handshakeTimeout": "10",
"helloExtensionIncludes": [],
"hostnameBlacklist": "none",
"hostnameWhitelist": "none",
"inheritCaCertkeychain": "false",
"inheritCertkeychain": "false",
"key": "/Common/default.key",
"logPublisher": "/Common/sys-ssl-publisher",
"maxActiveHandshakes": "indefinite",
"maxAggregateRenegotiationPerMinute": "indefinite",
"maxRenegotiationsPerMinute": 5,
"maximumRecordSize": 16384,
"modSslMethods": "disabled",
"mode": "enabled",
"notifyCertStatusToVirtualServer": "disabled",
"ocspStapling": "disabled",
"tmOptions": "{ dont-insert-empty-fragments no-tlsv1.3 no-dtlsv1.2 }",
"peerCertMode": "ignore",
"peerNoRenegotiateTimeout": "10",
"proxyCaCert": "none",
"proxyCaKey": "none",
"proxySsl": "disabled",
"proxySslPassthrough": "disabled",
"renegotiateMaxRecordDelay": "indefinite",
"renegotiatePeriod": "indefinite",
"renegotiateSize": "indefinite",
"renegotiation": "enabled",
"retainCertificate": "true",
"secureRenegotiation": "require",
"serverName": "none",
"sessionMirroring": "disabled",
"sessionTicket": "disabled",
"sessionTicketTimeout": 0,
"sniDefault": "false",
"sniRequire": "false",
"sourceIpBlacklist": "none",
"sourceIpWhitelist": "none",
"sslC3d": "disabled",
"sslForwardProxy": "disabled",
"sslForwardProxyBypass": "disabled",
"sslForwardProxyVerifiedHandshake": "disabled",
"sslSignHash": "any",
"strictResume": "disabled",
"uncleanShutdown": "enabled",
"certKeyChain": [
{
"name": "default",
"appService": "none",
"cert": "/Common/default.crt",
"chain": "none",
"key": "/Common/default.key",
"usage": "SERVER"
}
]
}`, resourceName, clientsslDefault)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfileSslClientCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfileSslClientModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfileSslClientUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-profile-clientSsl"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testclientssl##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileSslClientCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testclientssl##\\) is invalid"),
			},
		},
	})
}

func TestAccBigipLtmProfileSslClientUnitReadError(t *testing.T) {
	resourceName := "/Common/test-profile-clientSsl"
	clientsslDefault := "/Common/clientssl"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s"}`, resourceName, clientsslDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/~Common~test-profile-clientSsl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested Client-SSL Profile (/Common/test-profile-clientSsl) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileSslClientCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested Client-SSL Profile \\(/Common/test-profile-clientSsl\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfileSslClientInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_client_ssl" "test-profile-clientSsl" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfileSslClientCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_client_ssl" "test-profile-clientSsl" {
  name    = "%s"
  defaults_from = "/Common/clientssl"
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

func testBigipLtmProfileSslClientModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_client_ssl" "test-profile-clientSsl" {
  name    = "%s"
  defaults_from = "/Common/clientssl"
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
