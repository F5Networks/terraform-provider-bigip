/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipLtmProfileTcpSchema(t *testing.T) {
	r := resourceBigipLtmProfileTcp()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.CreateContext == nil {
		t.Fatal("Expected CreateContext to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}
	if r.UpdateContext == nil {
		t.Fatal("Expected UpdateContext to be defined")
	}
	if r.DeleteContext == nil {
		t.Fatal("Expected DeleteContext to be defined")
	}

	nameSchema, ok := r.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist in schema")
	}
	if !nameSchema.Required {
		t.Error("Expected field 'name' to be required")
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitLtmProfileTcpCreateReadUpdateDelete(t *testing.T) {
	name := "test-tcp"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/tcp","idleTimeout":300,"closeWaitTimeout":5,"finWait_2Timeout":300,"finWaitTimeout":5,"congestionControl":"high-speed","delayedAcks":"enabled","nagle":"disabled","earlyRetransmit":"enabled","tailLossProbe":"enabled","initCwnd":10,"zeroWindowTimeout":20000,"sendBufferSize":131072,"receiveWindowSize":65535,"proxyBufferHigh":49152,"timeWaitRecycle":"enabled","verifiedAccept":"disabled","keepAliveInterval":1800,"deferredAccept":"disabled","fastOpen":"enabled"}`, name)
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileTcp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                          name,
		"defaults_from":                 "/Common/tcp",
		"idle_timeout":                  300,
		"close_wait_timeout":            5,
		"finwait_2timeout":              300,
		"finwait_timeout":               5,
		"congestion_control":            "high-speed",
		"delayed_acks":                  "enabled",
		"nagle":                         "disabled",
		"early_retransmit":              "enabled",
		"tailloss_probe":                "enabled",
		"initial_congestion_windowsize": 10,
		"zerowindow_timeout":            20000,
		"send_buffersize":               131072,
		"receive_windowsize":            65535,
		"proxybuffer_high":              49152,
		"timewait_recycle":              "enabled",
		"verified_accept":               "disabled",
		"keepalive_interval":            1800,
		"deferred_accept":               "disabled",
		"fast_open":                     "enabled",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileTcpCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileTcpRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, 300, d.Get("idle_timeout").(int))

	diags = resourceBigipLtmProfileTcpUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileTcpDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileTcpCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileTcp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-tcp",
	}, "")

	diags := resourceBigipLtmProfileTcpCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileTcpRead_Error(t *testing.T) {
	name := "test-tcp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileTcp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileTcpRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileTcpUpdate_Error(t *testing.T) {
	name := "test-tcp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileTcp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileTcpUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileTcpDelete_Error(t *testing.T) {
	name := "test-tcp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/tcp/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileTcp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileTcpDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getTCPProfileConfig
// ---------------------------------------------------------------------

func TestGetTCPProfileConfig(t *testing.T) {
	r := resourceBigipLtmProfileTcp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                          "test-tcp",
		"partition":                     "Common",
		"defaults_from":                 "/Common/tcp",
		"idle_timeout":                  100,
		"close_wait_timeout":            6,
		"finwait_2timeout":              200,
		"finwait_timeout":               7,
		"send_buffersize":               65536,
		"receive_windowsize":            32768,
		"proxybuffer_high":              8192,
		"zerowindow_timeout":            10000,
		"keepalive_interval":            900,
		"congestion_control":            "cdg",
		"initial_congestion_windowsize": 20,
		"delayed_acks":                  "disabled",
		"nagle":                         "auto",
		"early_retransmit":              "disabled",
		"tailloss_probe":                "disabled",
		"timewait_recycle":              "disabled",
		"verified_accept":               "enabled",
		"deferred_accept":               "enabled",
		"fast_open":                     "disabled",
	}, "")

	cfg := getTCPProfileConfig(d, &bigip.Tcp{})
	require.Equal(t, "Common", cfg.Partition)
	require.Equal(t, "/Common/tcp", cfg.DefaultsFrom)
	require.Equal(t, 100, cfg.IdleTimeout)
	require.Equal(t, 6, cfg.CloseWaitTimeout)
	require.Equal(t, 200, cfg.FinWait_2Timeout)
	require.Equal(t, 7, cfg.FinWaitTimeout)
	require.Equal(t, 65536, cfg.SendBufferSize)
	require.Equal(t, 32768, cfg.ReceiveWindowSize)
	require.Equal(t, 8192, cfg.ProxyBufferHigh)
	require.Equal(t, 10000, cfg.ZeroWindowTimeout)
	require.Equal(t, 900, cfg.KeepAliveInterval)
	require.Equal(t, "cdg", cfg.CongestionControl)
	require.Equal(t, 20, cfg.InitCwnd)
	require.Equal(t, "disabled", cfg.DelayedAcks)
	require.Equal(t, "auto", cfg.Nagle)
	require.Equal(t, "disabled", cfg.EarlyRetransmit)
	require.Equal(t, "disabled", cfg.TailLossProbe)
	require.Equal(t, "disabled", cfg.TimeWaitRecycle)
	require.Equal(t, "enabled", cfg.VerifiedAccept)
	require.Equal(t, "enabled", cfg.DeferredAccept)
	require.Equal(t, "disabled", cfg.FastOpen)
}
