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

func TestResourceBigipLtmProfileHttpcompressSchema(t *testing.T) {
	r := resourceBigipLtmProfileHttpcompress()

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

func TestUnitLtmProfileHttpcompressCreateReadUpdateDelete(t *testing.T) {
	name := "test-httpcompress"
	created := false

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		created = true
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if !created {
				w.WriteHeader(http.StatusNotFound)
				_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
				return
			}
			_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/httpcompression","varyHeader":"enabled","cpuSaver":"enabled"}`, name, name)
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
	r := resourceBigipLtmProfileHttpcompress()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/httpcompression",
		"vary_header":   "enabled",
		"cpu_saver":     "enabled",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileHttpcompressCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileHttpcompressRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "enabled", d.Get("vary_header").(string))

	diags = resourceBigipLtmProfileHttpcompressUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileHttpcompressDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileHttpcompressCreate_AlreadyExists(t *testing.T) {
	// If GetHttpcompress's pre-check finds an object whose FullPath already
	// equals name, Create skips the POST entirely and just Reads.
	name := "test-httpcompress"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttpcompress()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, "")

	diags := resourceBigipLtmProfileHttpcompressCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())
}

func TestUnitLtmProfileHttpcompressCreate_Error(t *testing.T) {
	name := "test-httpcompress"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttpcompress()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, "")

	diags := resourceBigipLtmProfileHttpcompressCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileHttpcompressRead_Error(t *testing.T) {
	name := "test-httpcompress"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttpcompress()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileHttpcompressRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileHttpcompressUpdate_Error(t *testing.T) {
	name := "test-httpcompress"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttpcompress()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileHttpcompressUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileHttpcompressDelete_Error(t *testing.T) {
	name := "test-httpcompress"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http-compression/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttpcompress()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileHttpcompressDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getHTTPCompressProfileConfig
// ---------------------------------------------------------------------

func TestGetHTTPCompressProfileConfig(t *testing.T) {
	r := resourceBigipLtmProfileHttpcompress()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                   "test-httpcompress",
		"defaults_from":          "/Common/httpcompression",
		"uri_exclude":            []interface{}{"/exclude"},
		"uri_include":            []interface{}{"/include"},
		"content_type_include":   []interface{}{"text/html"},
		"content_type_exclude":   []interface{}{"image/jpeg"},
		"compression_buffersize": 8192,
		"gzip_compression_level": 6,
		"gzip_memory_level":      8192,
		"gzip_window_size":       16384,
		"keep_accept_encoding":   "enabled",
		"vary_header":            "disabled",
		"cpu_saver":              "disabled",
	}, "")

	cfg := getHTTPCompressProfileConfig(d, &bigip.Httpcompress{})
	require.Equal(t, "/Common/httpcompression", cfg.DefaultsFrom)
	require.Equal(t, []string{"/exclude"}, cfg.UriExclude)
	require.Equal(t, []string{"/include"}, cfg.UriInclude)
	require.Equal(t, []string{"text/html"}, cfg.ContentTypeInclude)
	require.Equal(t, []string{"image/jpeg"}, cfg.ContentTypeExclude)
	require.Equal(t, 8192, cfg.BufferSize)
	require.Equal(t, 6, cfg.GzipLevel)
	require.Equal(t, 8192, cfg.GzipMemoryLevel)
	require.Equal(t, 16384, cfg.GzipWindowSize)
	require.Equal(t, "enabled", cfg.KeepAcceptEncoding)
	require.Equal(t, "disabled", cfg.VaryHeader)
	require.Equal(t, "disabled", cfg.CPUSaver)
}
