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

func TestResourceBigipLtmPersistenceProfileDstAddrSchema(t *testing.T) {
	r := resourceBigipLtmPersistenceProfileDstAddr()

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

func TestUnitLtmPersistenceProfileDstAddrCreateReadUpdateDeleteWithTimeout(t *testing.T) {
	name := "/Common/test-dstaddr"
	mangled := MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		}
	})
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/dest_addr","matchAcrossPools":"disabled","matchAcrossServices":"disabled","matchAcrossVirtuals":"disabled","mirror":"disabled","overrideConnectionLimit":"disabled","timeout":"180","hashAlgorithm":"default","mask":"255.255.255.255","appService":"none"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileDstAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                name,
		"defaults_from":       "/Common/dest_addr",
		"timeout":             180,
		"hash_algorithm":      "default",
		"mask":                "255.255.255.255",
		"app_service":         "none",
		"mirror":              "enabled",
		"override_conn_limit": "enabled",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmPersistenceProfileDstAddrCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmPersistenceProfileDstAddrRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "/Common/dest_addr", d.Get("defaults_from").(string))
	require.Equal(t, 180, d.Get("timeout").(int))
	require.Equal(t, "default", d.Get("hash_algorithm").(string))
	require.Equal(t, "255.255.255.255", d.Get("mask").(string))

	diags = resourceBigipLtmPersistenceProfileDstAddrUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmPersistenceProfileDstAddrDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmPersistenceProfileDstAddrUpdateNoTimeout(t *testing.T) {
	name := "/Common/test-dstaddr-noto"
	mangled := MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/dest_addr","timeout":"0"}`, name)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileDstAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/dest_addr",
		"timeout":       0,
	}, name)

	diags := resourceBigipLtmPersistenceProfileDstAddrUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitLtmPersistenceProfileDstAddrCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileDstAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/test-dstaddr-err",
		"defaults_from": "/Common/dest_addr",
	}, "")

	diags := resourceBigipLtmPersistenceProfileDstAddrCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileDstAddrRead_Error(t *testing.T) {
	name := "/Common/test-dstaddr-rerr"
	mangled := MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileDstAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/dest_addr",
	}, name)

	diags := resourceBigipLtmPersistenceProfileDstAddrRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// TestUnitLtmPersistenceProfileDstAddrUpdate_ErrorWithDeleteSuccess exercises
// the Update error branch (with timeout set, so the ModifyDestAddrPersistenceProfile
// call is made) where the follow-up delete succeeds.
func TestUnitLtmPersistenceProfileDstAddrUpdate_ErrorWithDeleteSuccess(t *testing.T) {
	name := "/Common/test-dstaddr-uerr"
	mangled := MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
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
	r := resourceBigipLtmPersistenceProfileDstAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/dest_addr",
		"timeout":       120,
	}, name)

	diags := resourceBigipLtmPersistenceProfileDstAddrUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// TestUnitLtmPersistenceProfileDstAddrUpdate_ErrorWithDeleteError exercises the
// Update error branch where the follow-up delete also fails, taking the
// errdel != nil return path (no-timeout variant).
func TestUnitLtmPersistenceProfileDstAddrUpdate_ErrorWithDeleteError(t *testing.T) {
	name := "/Common/test-dstaddr-uderr"
	mangled := MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileDstAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/dest_addr",
		"timeout":       0,
	}, name)

	diags := resourceBigipLtmPersistenceProfileDstAddrUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileDstAddrDelete_Error(t *testing.T) {
	name := "/Common/test-dstaddr-derr"
	mangled := MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/dest-addr/"+mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmPersistenceProfileDstAddr()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/dest_addr",
	}, name)

	diags := resourceBigipLtmPersistenceProfileDstAddrDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
