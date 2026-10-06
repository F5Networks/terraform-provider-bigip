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

func TestResourceBigipLtmProfileFastl4Schema(t *testing.T) {
	r := resourceBigipLtmProfileFastl4()

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

func TestUnitLtmProfileFastl4CreateReadUpdateDelete(t *testing.T) {
	name := "test-fastl4"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/fastl4", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/fastl4/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/fastL4","clientTimeout":30,"explicitFlowMigration":"enabled","ipTosToClient":"pass-through","ipTosToServer":"pass-through","hardwareSynCookie":"disabled","idleTimeout":"300","keepAliveInterval":"disabled","tcpHandshakeTimeout":"5","looseInitialization":"disabled","looseClose":"disabled","lateBinding":"enabled","receiveWindowSize":0}`, name)
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
	r := resourceBigipLtmProfileFastl4()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                   name,
		"defaults_from":          "/Common/fastL4",
		"late_binding":           "enabled",
		"explicitflow_migration": "enabled",
		"client_timeout":         30,
		"iptos_toclient":         "pass-through",
		"iptos_toserver":         "pass-through",
		"hardware_syncookie":     "disabled",
		"idle_timeout":           "300",
		"keepalive_interval":     "disabled",
		"tcp_handshake_timeout":  "5",
		"loose_initiation":       "disabled",
		"loose_close":            "disabled",
		"receive_windowsize":     0,
	}, "")

	ctx := context.Background()

	diags := resourceBigipProfileLtmFastl4Create(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileFastl4Read(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "enabled", d.Get("late_binding").(string))

	diags = resourceBigipLtmProfileFastl4Update(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileFastl4Delete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileFastl4Create_ExplicitFlowValidationError(t *testing.T) {
	// explicitflow_migration=enabled requires late_binding=enabled; this
	// should fail before any HTTP call is made.
	client := NewUnitTestClient("http://127.0.0.1:0")
	r := resourceBigipLtmProfileFastl4()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                   "test-fastl4",
		"explicitflow_migration": "enabled",
		"late_binding":           "disabled",
	}, "")

	diags := resourceBigipProfileLtmFastl4Create(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFastl4Update_ExplicitFlowValidationError(t *testing.T) {
	client := NewUnitTestClient("http://127.0.0.1:0")
	r := resourceBigipLtmProfileFastl4()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                   "test-fastl4",
		"explicitflow_migration": "enabled",
		"late_binding":           "disabled",
	}, "test-fastl4")

	diags := resourceBigipLtmProfileFastl4Update(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFastl4Create_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/fastl4", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFastl4()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-fastl4",
	}, "")

	diags := resourceBigipProfileLtmFastl4Create(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFastl4Read_Error(t *testing.T) {
	name := "test-fastl4"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/fastl4/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFastl4()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileFastl4Read(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFastl4Update_Error(t *testing.T) {
	name := "test-fastl4"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/fastl4/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFastl4()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileFastl4Update(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileFastl4Delete_Error(t *testing.T) {
	name := "test-fastl4"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/fastl4/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileFastl4()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileFastl4Delete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getFastL4ProfileConfig
// ---------------------------------------------------------------------

func TestGetFastL4ProfileConfig(t *testing.T) {
	r := resourceBigipLtmProfileFastl4()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                   "test-fastl4",
		"defaults_from":          "/Common/fastL4",
		"client_timeout":         60,
		"late_binding":           "enabled",
		"explicitflow_migration": "enabled",
		"hardware_syncookie":     "enabled",
		"idle_timeout":           "600",
		"iptos_toclient":         "mimic",
		"iptos_toserver":         "mimic",
		"keepalive_interval":     "30",
		"tcp_handshake_timeout":  "10",
		"loose_initiation":       "enabled",
		"loose_close":            "enabled",
		"receive_windowsize":     4096,
	}, "")

	cfg := getFastL4ProfileConfig(d, &bigip.Fastl4{})
	require.Equal(t, "/Common/fastL4", cfg.DefaultsFrom)
	require.Equal(t, 60, cfg.ClientTimeout)
	require.Equal(t, "enabled", cfg.LateBinding)
	require.Equal(t, "enabled", cfg.ExplicitFlowMigration)
	require.Equal(t, "enabled", cfg.HardwareSynCookie)
	require.Equal(t, "600", cfg.IdleTimeout)
	require.Equal(t, "mimic", cfg.IpTosToClient)
	require.Equal(t, "mimic", cfg.IpTosToServer)
	require.Equal(t, "30", cfg.KeepAliveInterval)
	require.Equal(t, "10", cfg.TCPHandshakeTimeout)
	require.Equal(t, "enabled", cfg.LooseInitialization)
	require.Equal(t, "enabled", cfg.LooseClose)
	require.Equal(t, 4096, cfg.ReceiveWindowSize)
}
