/*
Copyright 2019 F5 Networks Inc.
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

func TestResourceBigipNetTunnelSchema(t *testing.T) {
	r := resourceBigipNetTunnel()

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

	for _, field := range []string{"name", "local_address", "profile"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitNetTunnelCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-tunnel"
	mangled := "/mgmt/tm/net/tunnels/tunnel/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/tunnels/tunnel", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","localAddress":"10.0.0.1","profile":"gre","mtu":1500}`, name)
		case http.MethodPut:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetTunnel()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"local_address": "10.0.0.1",
		"profile":       "gre",
	}, "")

	ctx := context.Background()

	diags := resourceBigipNetTunnelCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipNetTunnelRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "10.0.0.1", d.Get("local_address").(string))

	diags = resourceBigipNetTunnelUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipNetTunnelDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitNetTunnelCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/tunnels/tunnel", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetTunnel()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/test-tunnel",
		"local_address": "10.0.0.1",
		"profile":       "gre",
	}, "")

	diags := resourceBigipNetTunnelCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetTunnelRead_Error(t *testing.T) {
	name := "/Common/test-tunnel"
	mangled := "/mgmt/tm/net/tunnels/tunnel/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetTunnel()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"local_address": "10.0.0.1",
		"profile":       "gre",
	}, name)

	diags := resourceBigipNetTunnelRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetTunnelUpdate_Error(t *testing.T) {
	name := "/Common/test-tunnel"
	mangled := "/mgmt/tm/net/tunnels/tunnel/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetTunnel()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"local_address": "10.0.0.1",
		"profile":       "gre",
	}, name)

	diags := resourceBigipNetTunnelUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetTunnelDelete_Error(t *testing.T) {
	name := "/Common/test-tunnel"
	mangled := "/mgmt/tm/net/tunnels/tunnel/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetTunnel()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"local_address": "10.0.0.1",
		"profile":       "gre",
	}, name)

	diags := resourceBigipNetTunnelDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getConfig
// ---------------------------------------------------------------------

func TestGetConfig_Tunnel(t *testing.T) {
	r := resourceBigipNetTunnel()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":              "/Common/test-tunnel",
		"app_service":       "my-app",
		"auto_last_hop":     "default",
		"description":       "my tunnel",
		"local_address":     "10.0.0.1",
		"mode":              "bidirectional",
		"partition":         "Common",
		"profile":           "gre",
		"remote_address":    "10.0.0.2",
		"secondary_address": "10.0.0.3",
		"tos":               "preserve",
		"traffic_group":     "/Common/traffic-group-1",
		"transparent":       "disabled",
		"use_pmtu":          "enabled",
		"idle_timeout":      300,
		"key":               5,
		"mtu":               1400,
	}, "")

	cfg := getConfig(d, &bigip.Tunnel{})
	require.Equal(t, "my-app", cfg.AppService)
	require.Equal(t, "default", cfg.AutoLasthop)
	require.Equal(t, "my tunnel", cfg.Description)
	require.Equal(t, "10.0.0.1", cfg.LocalAddress)
	require.Equal(t, "gre", cfg.Profile)
	require.Equal(t, 300, cfg.IdleTimeout)
	require.Equal(t, 5, cfg.Key)
	require.Equal(t, "bidirectional", cfg.Mode)
	require.Equal(t, 1400, cfg.Mtu)
	require.Equal(t, "Common", cfg.Partition)
	require.Equal(t, "10.0.0.2", cfg.RemoteAddress)
	require.Equal(t, "10.0.0.3", cfg.SecondaryAddress)
	require.Equal(t, "preserve", cfg.Tos)
	require.Equal(t, "/Common/traffic-group-1", cfg.TrafficGroup)
	require.Equal(t, "disabled", cfg.Transparent)
	require.Equal(t, "enabled", cfg.UsePmtu)
}
