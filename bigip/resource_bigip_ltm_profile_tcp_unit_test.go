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

func TestAccBigipLtmProfileTCPUnitInvalid(t *testing.T) {
	resourceName := "/Common/test-profile-tcp"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileTCPInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipLtmProfileTCPUnitCreate(t *testing.T) {
	resourceName := "/Common/test-profile-tcp"
	tcpDefault := "/Common/tcp"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","abc": "enabled",
"ackOnPush": "enabled",
"appService": "none",
"autoProxyBufferSize": "disabled",
"autoReceiveWindowSize": "disabled",
"autoSendBufferSize": "disabled",
"closeWaitTimeout": 5,
"cmetricsCache": "enabled",
"cmetricsCacheTimeout": 0,
"congestionControl": "high-speed",
"defaultsFrom": "/Common/tcp",
"deferredAccept": "disabled",
"delayWindowControl": "disabled",
"delayedAcks": "enabled",
"description": "none",
"dsack": "disabled",
"earlyRetransmit": "enabled",
"ecn": "enabled",
"enhancedLossRecovery": "enabled",
"fastOpen": "enabled",
"fastOpenCookieExpiration": 21600,
"finWait_2Timeout": 300,
"finWaitTimeout": 5,
"hardwareSynCookie": "enabled",
"idleTimeout": 300,
"initCwnd": 10,
"initRwnd": 10,
"ipDfMode": "pmtu",
"ipTosToClient": "0",
"ipTtlMode": "proxy",
"ipTtlV4": 255,
"ipTtlV6": 64,
"keepAliveInterval": 1800,
"limitedTransmit": "enabled",
"linkQosToClient": "0",
"maxRetrans": 8,
"maxSegmentSize": 1460,
"md5Signature": "disabled",
"minimumRto": 1000,
"mptcp": "disabled",
"mptcpCsum": "disabled",
"mptcpCsumVerify": "disabled",
"mptcpDebug": "disabled",
"mptcpFallback": "reset",
"mptcpFastjoin": "disabled",
"mptcpIdleTimeout": 300,
"mptcpJoinMax": 5,
"mptcpMakeafterbreak": "disabled",
"mptcpNojoindssack": "disabled",
"mptcpRtomax": 5,
"mptcpRxmitmin": 1000,
"mptcpSubflowmax": 6,
"mptcpTimeout": 3600,
"nagle": "disabled",
"pktLossIgnoreBurst": 0,
"pktLossIgnoreRate": 0,
"proxyBufferHigh": 65535,
"proxyBufferLow": 32768,
"proxyMss": "enabled",
"proxyOptions": "disabled",
"pushFlag": "default",
"ratePace": "enabled",
"ratePaceMaxRate": 0,
"receiveWindowSize": 65535,
"resetOnTimeout": "enabled",
"rexmtThresh": 3,
"selectiveAcks": "enabled",
"selectiveNack": "disabled",
"sendBufferSize": 131072,
"slowStart": "enabled",
"synCookieEnable": "enabled",
"synCookieWhitelist": "disabled",
"synMaxRetrans": 3,
"synRtoBase": 3000,
"tailLossProbe": "enabled",
"tcpOptions": "none",
"timeWaitRecycle": "enabled",
"timeWaitTimeout": "2000",
"timestamps": "enabled",
"verifiedAccept": "disabled",
"zeroWindowTimeout": 20000}
`, resourceName, tcpDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp/~Common~test-profile-tcp", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","abc": "enabled",
"ackOnPush": "enabled",
"appService": "none",
"autoProxyBufferSize": "disabled",
"autoReceiveWindowSize": "disabled",
"autoSendBufferSize": "disabled",
"closeWaitTimeout": 5,
"cmetricsCache": "enabled",
"cmetricsCacheTimeout": 0,
"congestionControl": "high-speed",
"defaultsFrom": "/Common/tcp",
"deferredAccept": "disabled",
"delayWindowControl": "disabled",
"delayedAcks": "enabled",
"description": "none",
"dsack": "disabled",
"earlyRetransmit": "enabled",
"ecn": "enabled",
"enhancedLossRecovery": "enabled",
"fastOpen": "enabled",
"fastOpenCookieExpiration": 21600,
"finWait_2Timeout": 300,
"finWaitTimeout": 5,
"hardwareSynCookie": "enabled",
"idleTimeout": 300,
"initCwnd": 10,
"initRwnd": 10,
"ipDfMode": "pmtu",
"ipTosToClient": "0",
"ipTtlMode": "proxy",
"ipTtlV4": 255,
"ipTtlV6": 64,
"keepAliveInterval": 1800,
"limitedTransmit": "enabled",
"linkQosToClient": "0",
"maxRetrans": 8,
"maxSegmentSize": 1460,
"md5Signature": "disabled",
"minimumRto": 1000,
"mptcp": "disabled",
"mptcpCsum": "disabled",
"mptcpCsumVerify": "disabled",
"mptcpDebug": "disabled",
"mptcpFallback": "reset",
"mptcpFastjoin": "disabled",
"mptcpIdleTimeout": 300,
"mptcpJoinMax": 5,
"mptcpMakeafterbreak": "disabled",
"mptcpNojoindssack": "disabled",
"mptcpRtomax": 5,
"mptcpRxmitmin": 1000,
"mptcpSubflowmax": 6,
"mptcpTimeout": 3600,
"nagle": "disabled",
"pktLossIgnoreBurst": 0,
"pktLossIgnoreRate": 0,
"proxyBufferHigh": 65535,
"proxyBufferLow": 32768,
"proxyMss": "enabled",
"proxyOptions": "disabled",
"pushFlag": "default",
"ratePace": "enabled",
"ratePaceMaxRate": 0,
"receiveWindowSize": 65535,
"resetOnTimeout": "enabled",
"rexmtThresh": 3,
"selectiveAcks": "enabled",
"selectiveNack": "disabled",
"sendBufferSize": 131072,
"slowStart": "enabled",
"synCookieEnable": "enabled",
"synCookieWhitelist": "disabled",
"synMaxRetrans": 3,
"synRtoBase": 3000,
"tailLossProbe": "enabled",
"tcpOptions": "none",
"timeWaitRecycle": "enabled",
"timeWaitTimeout": "2000",
"timestamps": "enabled",
"verifiedAccept": "disabled",
"zeroWindowTimeout": 20000}`, resourceName, tcpDefault)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp/~Common~test-profile-tcp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Expected method 'PUT', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s","abc": "enabled",
"ackOnPush": "enabled",
"appService": "none",
"autoProxyBufferSize": "disabled",
"autoReceiveWindowSize": "disabled",
"autoSendBufferSize": "disabled",
"closeWaitTimeout": 10,
"cmetricsCache": "enabled",
"cmetricsCacheTimeout": 0,
"congestionControl": "high-speed",
"defaultsFrom": "/Common/tcp",
"deferredAccept": "disabled",
"delayWindowControl": "disabled",
"delayedAcks": "enabled",
"description": "none",
"dsack": "disabled",
"earlyRetransmit": "enabled",
"ecn": "enabled",
"enhancedLossRecovery": "enabled",
"fastOpen": "enabled",
"fastOpenCookieExpiration": 21600,
"finWait_2Timeout": 300,
"finWaitTimeout": 5,
"hardwareSynCookie": "enabled",
"idleTimeout": 300,
"initCwnd": 10,
"initRwnd": 10,
"ipDfMode": "pmtu",
"ipTosToClient": "0",
"ipTtlMode": "proxy",
"ipTtlV4": 255,
"ipTtlV6": 64,
"keepAliveInterval": 1800,
"limitedTransmit": "enabled",
"linkQosToClient": "0",
"maxRetrans": 8,
"maxSegmentSize": 1460,
"md5Signature": "disabled",
"minimumRto": 1000,
"mptcp": "disabled",
"mptcpCsum": "disabled",
"mptcpCsumVerify": "disabled",
"mptcpDebug": "disabled",
"mptcpFallback": "reset",
"mptcpFastjoin": "disabled",
"mptcpIdleTimeout": 300,
"mptcpJoinMax": 5,
"mptcpMakeafterbreak": "disabled",
"mptcpNojoindssack": "disabled",
"mptcpRtomax": 5,
"mptcpRxmitmin": 1000,
"mptcpSubflowmax": 6,
"mptcpTimeout": 3600,
"nagle": "disabled",
"pktLossIgnoreBurst": 0,
"pktLossIgnoreRate": 0,
"proxyBufferHigh": 65535,
"proxyBufferLow": 32768,
"proxyMss": "enabled",
"proxyOptions": "disabled",
"pushFlag": "default",
"ratePace": "enabled",
"ratePaceMaxRate": 0,
"receiveWindowSize": 65535,
"resetOnTimeout": "enabled",
"rexmtThresh": 3,
"selectiveAcks": "enabled",
"selectiveNack": "disabled",
"sendBufferSize": 131072,
"slowStart": "enabled",
"synCookieEnable": "enabled",
"synCookieWhitelist": "disabled",
"synMaxRetrans": 3,
"synRtoBase": 3000,
"tailLossProbe": "enabled",
"tcpOptions": "none",
"timeWaitRecycle": "enabled",
"timeWaitTimeout": "2000",
"timestamps": "enabled",
"verifiedAccept": "disabled",
"zeroWindowTimeout": 20000}`, resourceName, tcpDefault)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmProfileTCPCreate(resourceName, server.URL),
			},
			{
				Config:             testBigipLtmProfileTCPModify(resourceName, server.URL),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigipLtmProfileTCPUnitCreateError(t *testing.T) {
	resourceName := "/Common/test-profile-tcp"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, "The requested object name (/Common/testprofiletcp##) is invalid", http.StatusBadRequest)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileTCPCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 400 :: The requested object name \\(/Common/testprofiletcp##\\) is invalid"),
			},
		},
	})
}
func TestAccBigipLtmProfileTCPUnitReadError(t *testing.T) {
	resourceName := "/Common/test-profile-tcp"
	tcpDefault := "/Common/tcp"
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
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"%s", "basicAuthRealm": "none"}`, resourceName, tcpDefault)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp/~Common~test-profile-tcp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "The requested HTTP Profile (/Common/test-profile-tcp) was not found", http.StatusNotFound)
	})

	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmProfileTCPCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("HTTP 404 :: The requested HTTP Profile \\(/Common/test-profile-tcp\\) was not found"),
			},
		},
	})
}

func testBigipLtmProfileTCPInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_tcp" "test-tcp" {
  name       = "%s"
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`, resourceName)
}

func testBigipLtmProfileTCPCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_tcp" "test-tcp" {
  name    = "%s"
  defaults_from = "/Common/tcp"
  close_wait_timeout = 5
  idle_timeout = 300
  finwait_2timeout = 300
  finwait_timeout = 5
  congestion_control = "high-speed"
  delayed_acks = "enabled"
  nagle = "disabled"
  early_retransmit = "enabled"
  tailloss_probe = "enabled"
  initial_congestion_windowsize = 10 
  zerowindow_timeout = 20000
  send_buffersize = 131072
  receive_windowsize = 65535
  proxybuffer_high = 65535
  timewait_recycle = "enabled"
  verified_accept = "disabled"
  keepalive_interval = 1800
  deferred_accept  = "disabled"
  fast_open  = "enabled"
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}

func testBigipLtmProfileTCPModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_profile_tcp" "test-tcp" {
  name    = "%s"
  close_wait_timeout = 10
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, resourceName, url)
}
