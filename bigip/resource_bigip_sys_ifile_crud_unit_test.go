/*
Copyright 2024 F5 Networks Inc.
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

func TestResourceBigipSysIfileSchema(t *testing.T) {
	r := resourceBigipSysIfile()

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

	for _, field := range []string{"name", "content"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}
	if r.Schema["partition"].Default != "Common" {
		t.Errorf("Expected 'partition' to default to 'Common', got %v", r.Schema["partition"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitSysIfileCreateReadUpdateDelete(t *testing.T) {
	name := "testfile"
	fullPath := "/Common/" + name
	mangled := "/mgmt/tm/sys/file/ifile/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ifile", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, fullPath)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"%s","checksum":"abc123","size":42}`, name, fullPath)
		case http.MethodPut:
			_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, fullPath)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"partition": "Common",
		"content":   "fake-content",
	}, "")

	ctx := context.Background()

	diags := resourceBigipSysIfileCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, fullPath, d.Id())

	diags = resourceBigipSysIfileRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, 42, d.Get("size").(int))

	diags = resourceBigipSysIfileUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSysIfileDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitSysIfileCreate_UploadError(t *testing.T) {
	name := "testfile"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"upload failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"partition": "Common",
		"content":   "fake-content",
	}, "")

	diags := resourceBigipSysIfileCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysIfileCreate_CreateError(t *testing.T) {
	name := "testfile"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ifile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"partition": "Common",
		"content":   "fake-content",
	}, "")

	diags := resourceBigipSysIfileCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysIfileRead_NotFound(t *testing.T) {
	fullPath := "/Common/testfile"
	mangled := "/mgmt/tm/sys/file/ifile/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testfile",
		"partition": "Common",
		"content":   "fake-content",
	}, fullPath)

	diags := resourceBigipSysIfileRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysIfileUpdate_Error(t *testing.T) {
	fullPath := "/Common/testfile"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/testfile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"upload failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testfile",
		"partition": "Common",
		"content":   "fake-content",
	}, fullPath)

	diags := resourceBigipSysIfileUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysIfileDelete_Error(t *testing.T) {
	fullPath := "/Common/testfile"
	mangled := "/mgmt/tm/sys/file/ifile/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testfile",
		"partition": "Common",
		"content":   "fake-content",
	}, fullPath)

	diags := resourceBigipSysIfileDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// buildIFileFullPath
// ---------------------------------------------------------------------

func TestBuildIFileFullPath(t *testing.T) {
	require.Equal(t, "/Common/testfile", buildIFileFullPath("Common", "", "testfile"))
	require.Equal(t, "/Common/subdir/testfile", buildIFileFullPath("Common", "subdir", "testfile"))
}
