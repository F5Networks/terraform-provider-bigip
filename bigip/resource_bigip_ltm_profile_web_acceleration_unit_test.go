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

func TestResourceBigipLtmProfileWebAccelerationSchema(t *testing.T) {
	r := resourceBigipLtmProfileWebAcceleration()

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

func TestUnitLtmProfileWebAccelerationCreateReadUpdateDelete(t *testing.T) {
	name := "test-web-acceleration"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/web-acceleration", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/web-acceleration/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/webacceleration","cacheSize":100,"cacheMaxEntries":10,"cacheMaxAge":3600,"cacheObjectMinSize":500,"cacheObjectMaxSize":50000,"cacheUriExclude":["/exclude"],"cacheUriInclude":["/include"],"cacheUriIncludeOverride":["/override"],"cacheUriPinned":["/pinned"],"cacheClientCacheControlMode":"all","cacheAgingRate":9,"cacheInsertAgeHeader":"enabled"}`, name)
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
	r := resourceBigipLtmProfileWebAcceleration()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                            name,
		"defaults_from":                   "/Common/webacceleration",
		"cache_size":                      100,
		"cache_max_entries":               10,
		"cache_max_age":                   3600,
		"cache_object_min_size":           500,
		"cache_object_max_size":           50000,
		"cache_uri_exclude":               []interface{}{"/exclude"},
		"cache_uri_include":               []interface{}{"/include"},
		"cache_uri_include_override":      []interface{}{"/override"},
		"cache_uri_pinned":                []interface{}{"/pinned"},
		"cache_client_cache_control_mode": "all",
		"cache_aging_rate":                9,
		"cache_insert_age_header":         "enabled",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileWebAccelerationCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileWebAccelerationRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, 3600, d.Get("cache_max_age").(int))

	diags = resourceBigipLtmProfileWebAccelerationUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileWebAccelerationDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileWebAccelerationCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/web-acceleration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipLtmProfileWebAcceleration()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-web-acceleration",
	}, "")

	diags := resourceBigipLtmProfileWebAccelerationCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileWebAccelerationRead_Error(t *testing.T) {
	name := "test-web-acceleration"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/web-acceleration/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileWebAcceleration()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileWebAccelerationRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileWebAccelerationUpdate_Error(t *testing.T) {
	name := "test-web-acceleration"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/web-acceleration/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileWebAcceleration()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileWebAccelerationUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileWebAccelerationDelete_Error(t *testing.T) {
	name := "test-web-acceleration"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/web-acceleration/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileWebAcceleration()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileWebAccelerationDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getHttpProfileWebAccelerationConfig
// ---------------------------------------------------------------------

func TestGetHttpProfileWebAccelerationConfig(t *testing.T) {
	r := resourceBigipLtmProfileWebAcceleration()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                            "test-web-acceleration",
		"defaults_from":                   "/Common/webacceleration",
		"cache_size":                      200,
		"cache_max_entries":               10,
		"cache_max_age":                   7200,
		"cache_object_min_size":           500,
		"cache_object_max_size":           50000,
		"cache_uri_exclude":               []interface{}{"/exclude"},
		"cache_uri_include":               []interface{}{"/include"},
		"cache_uri_include_override":      []interface{}{"/override"},
		"cache_uri_pinned":                []interface{}{"/pinned"},
		"cache_client_cache_control_mode": "all",
		"cache_insert_age_header":         "enabled",
		"cache_aging_rate":                5,
	}, "")

	cfg := getHttpProfileWebAccelerationConfig(d, &bigip.WebAccelerationProfileService{})
	require.Equal(t, "/Common/webacceleration", cfg.DefaultsFrom)
	require.Equal(t, 200, cfg.CacheSize)
	require.Equal(t, 10, cfg.CacheMaxEntries)
	require.Equal(t, 7200, cfg.CacheMaxAge)
	require.Equal(t, 500, cfg.CacheObjectMinSize)
	require.Equal(t, 50000, cfg.CacheObjectMaxSize)
	require.Equal(t, []string{"/exclude"}, cfg.CacheUriExclude)
	require.Equal(t, []string{"/include"}, cfg.CacheUriInclude)
	require.Equal(t, []string{"/override"}, cfg.CacheUriIncludeOverride)
	require.Equal(t, []string{"/pinned"}, cfg.CacheUriPinned)
	require.Equal(t, "all", cfg.CacheClientCacheControlMode)
	require.Equal(t, "enabled", cfg.CacheInsertAgeHeader)
	require.Equal(t, 5, cfg.CacheAgingRate)
}
