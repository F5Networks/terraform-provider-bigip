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

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipSslKeySchema(t *testing.T) {
	r := resourceBigipSslKey()

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
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	nameSchema, ok := r.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist in schema")
	}
	if !nameSchema.Required {
		t.Error("Expected field 'name' to be required")
	}
	if r.Schema["partition"].Default != "Common" {
		t.Errorf("Expected 'partition' to default to 'Common', got %v", r.Schema["partition"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitSslKeyCreateReadUpdateDelete(t *testing.T) {
	name := "testkey"
	mangled := "/mgmt/tm/sys/file/ssl-key/" + MangleFullPath("/Common/"+name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslKey()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       name,
		"content":    "fake-key-content",
		"partition":  "Common",
		"passphrase": "secret",
	}, "")

	ctx := context.Background()

	diags := resourceBigipSslKeyCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipSslKeyRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	diags = resourceBigipSslKeyUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSslKeyDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitSslKeyCreate_UploadError(t *testing.T) {
	name := "testkey"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"upload failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslKey()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-key-content",
		"partition": "Common",
	}, "")

	diags := resourceBigipSslKeyCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSslKeyCreate_AddKeyError(t *testing.T) {
	name := "testkey"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-key", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"add key failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslKey()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-key-content",
		"partition": "Common",
	}, "")

	diags := resourceBigipSslKeyCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSslKeyRead_NotFound(t *testing.T) {
	name := "testkey"
	mangled := "/mgmt/tm/sys/file/ssl-key/" + MangleFullPath("/Common/"+name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslKey()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-key-content",
		"partition": "Common",
	}, name)

	diags := resourceBigipSslKeyRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSslKeyUpdate_ModifyKeyError(t *testing.T) {
	name := "testkey"
	mangled := "/mgmt/tm/sys/file/ssl-key/" + MangleFullPath("/Common/"+name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"modify failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslKey()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-key-content",
		"partition": "Common",
	}, name)

	diags := resourceBigipSslKeyUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSslKeyDelete_Error(t *testing.T) {
	name := "testkey"
	mangled := "/mgmt/tm/sys/file/ssl-key/" + MangleFullPath("/Common/"+name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslKey()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-key-content",
		"partition": "Common",
	}, name)

	diags := resourceBigipSslKeyDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
