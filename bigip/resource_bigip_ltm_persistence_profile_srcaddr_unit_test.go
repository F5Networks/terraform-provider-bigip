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

	"github.com/stretchr/testify/require"
)

func TestResourceBigipLtmPersistenceProfileSrcAddrSchema(t *testing.T) {
	r := resourceBigipLtmPersistenceProfileSrcAddr()

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

	for _, field := range []string{"name", "defaults_from"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}
}

func TestUnitLtmPersistenceProfileSrcAddrCreateReadUpdateDeleteWithTimeout(t *testing.T) {
	name := "/Common/test-srcaddr"
	mangled := MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/source_addr","matchAcrossPools":"disabled","matchAcrossServices":"disabled","matchAcrossVirtuals":"disabled","mirror":"disabled","overrideConnectionLimit":"disabled","timeout":"180","hashAlgorithm":"default","mapProxies":"enabled","mask":"255.255.255.255","appService":"none"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileSrcAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                name,
		"defaults_from":       "/Common/source_addr",
		"timeout":             180,
		"hash_algorithm":      "default",
		"map_proxies":         "enabled",
		"mask":                "255.255.255.255",
		"app_service":         "none",
		"mirror":              "enabled",
		"override_conn_limit": "enabled",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmPersistenceProfileSrcAddrCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmPersistenceProfileSrcAddrRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "/Common/source_addr", d.Get("defaults_from").(string))
	require.Equal(t, 180, d.Get("timeout").(int))
	require.Equal(t, "default", d.Get("hash_algorithm").(string))
	require.Equal(t, "enabled", d.Get("map_proxies").(string))
	require.Equal(t, "255.255.255.255", d.Get("mask").(string))

	diags = resourceBigipLtmPersistenceProfileSrcAddrUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmPersistenceProfileSrcAddrDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmPersistenceProfileSrcAddrUpdateNoTimeout(t *testing.T) {
	name := "/Common/test-srcaddr-noto"
	mangled := MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/source_addr","timeout":"0"}`, name)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileSrcAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/source_addr",
		"timeout":       0,
	}, name)

	diags := resourceBigipLtmPersistenceProfileSrcAddrUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitLtmPersistenceProfileSrcAddrCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileSrcAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/test-srcaddr-err",
		"defaults_from": "/Common/source_addr",
	}, "")

	diags := resourceBigipLtmPersistenceProfileSrcAddrCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileSrcAddrRead_Error(t *testing.T) {
	name := "/Common/test-srcaddr-rerr"
	mangled := MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileSrcAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/source_addr",
	}, name)

	diags := resourceBigipLtmPersistenceProfileSrcAddrRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// TestUnitLtmPersistenceProfileSrcAddrUpdate_ErrorWithDeleteSuccess exercises
// the Update error branch (with timeout set, so the ModifySourceAddrPersistenceProfile
// call is made) where the follow-up delete succeeds.
func TestUnitLtmPersistenceProfileSrcAddrUpdate_ErrorWithDeleteSuccess(t *testing.T) {
	name := "/Common/test-srcaddr-uerr"
	mangled := MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileSrcAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/source_addr",
		"timeout":       120,
	}, name)

	diags := resourceBigipLtmPersistenceProfileSrcAddrUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// TestUnitLtmPersistenceProfileSrcAddrUpdate_ErrorWithDeleteError exercises the
// Update error branch where the follow-up delete also fails, taking the
// errdel != nil return path (no-timeout variant).
func TestUnitLtmPersistenceProfileSrcAddrUpdate_ErrorWithDeleteError(t *testing.T) {
	name := "/Common/test-srcaddr-uderr"
	mangled := MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileSrcAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/source_addr",
		"timeout":       0,
	}, name)

	diags := resourceBigipLtmPersistenceProfileSrcAddrUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileSrcAddrDelete_Error(t *testing.T) {
	name := "/Common/test-srcaddr-derr"
	mangled := MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/source-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileSrcAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/source_addr",
	}, name)

	diags := resourceBigipLtmPersistenceProfileSrcAddrDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
