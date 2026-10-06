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

func TestResourceBigipLtmProfileHttp2Schema(t *testing.T) {
	r := resourceBigipLtmProfileHttp2()

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

func TestUnitLtmProfileHttp2CreateReadUpdateDelete(t *testing.T) {
	name := "test-http2"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/http2","concurrentStreamsPerConnection":10,"connectionIdleTimeout":300,"headerTableSize":4096,"enforceTlsRequirements":"enabled","frameSize":2048,"receiveWindow":32,"writeSize":16384,"insertHeader":"disabled","insertHeaderName":"X-HTTP2","activationModes":["alpn"]}`, name)
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
	client.Teem = true // skip real telemetry network call
	r := resourceBigipLtmProfileHttp2()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                              name,
		"defaults_from":                     "/Common/http2",
		"concurrent_streams_per_connection": 10,
		"connection_idle_timeout":           300,
		"header_table_size":                 4096,
		"enforce_tls_requirements":          "enabled",
		"frame_size":                        2048,
		"receive_window":                    32,
		"write_size":                        16384,
		"insert_header":                     "disabled",
		"insert_header_name":                "X-HTTP2",
		"activation_modes":                  []interface{}{"alpn"},
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileHttp2Create(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileHttp2Read(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, 4096, d.Get("header_table_size").(int))

	diags = resourceBigipLtmProfileHttp2Update(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileHttp2Delete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileHttp2Create_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipLtmProfileHttp2()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-http2",
	}, "")

	diags := resourceBigipLtmProfileHttp2Create(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileHttp2Read_Error(t *testing.T) {
	name := "test-http2"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttp2()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileHttp2Read(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileHttp2Update_Error(t *testing.T) {
	name := "test-http2"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttp2()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileHttp2Update(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileHttp2Delete_Error(t *testing.T) {
	name := "test-http2"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http2/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttp2()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileHttp2Delete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getHttp2ProfileConfig
// ---------------------------------------------------------------------

func TestGetHttp2ProfileConfig(t *testing.T) {
	r := resourceBigipLtmProfileHttp2()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                              "test-http2",
		"defaults_from":                     "/Common/http2",
		"concurrent_streams_per_connection": 20,
		"connection_idle_timeout":           600,
		"header_table_size":                 8192,
		"activation_modes":                  []interface{}{"alpn", "npn"},
		"enforce_tls_requirements":          "disabled",
		"frame_size":                        4096,
		"include_content_length":            "enabled",
		"insert_header":                     "enabled",
		"insert_header_name":                "X-Test",
		"receive_window":                    64,
		"write_size":                        32768,
	}, "")

	cfg := getHttp2ProfileConfig(d, &bigip.Http2{})
	require.Equal(t, "/Common/http2", cfg.DefaultsFrom)
	require.Equal(t, 20, cfg.ConcurrentStreamsPerConnection)
	require.Equal(t, 600, cfg.ConnectionIdleTimeout)
	require.Equal(t, 8192, cfg.HeaderTableSize)
	require.ElementsMatch(t, []string{"alpn", "npn"}, cfg.ActivationModes)
	require.Equal(t, "disabled", cfg.EnforceTLSRequirements)
	require.Equal(t, 4096, cfg.FrameSize)
	require.Equal(t, "enabled", cfg.IncludeContentLength)
	require.Equal(t, "enabled", cfg.InsertHeader)
	require.Equal(t, "X-Test", cfg.InsertHeaderName)
	require.Equal(t, 64, cfg.ReceiveWindow)
	require.Equal(t, 32768, cfg.WriteSize)
}
