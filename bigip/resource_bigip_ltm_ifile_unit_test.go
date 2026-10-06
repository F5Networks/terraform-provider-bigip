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
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipLtmIfileSchema(t *testing.T) {
	r := resourceBigipLtmIfile()

	require.NotNil(t, r.Schema)
	require.NotNil(t, r.CreateContext)
	require.NotNil(t, r.ReadContext)
	require.NotNil(t, r.UpdateContext)
	require.NotNil(t, r.DeleteContext)

	for _, field := range []string{"name", "file_name"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "expected field %q to exist in schema", field)
		require.True(t, s.Required, "expected field %q to be required", field)
	}
	require.Equal(t, "Common", r.Schema["partition"].Default)
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests
// ---------------------------------------------------------------------

func TestUnitLtmIfileCreateReadUpdateDelete(t *testing.T) {
	name := "testfile"
	fullPath := "/Common/" + name
	mangled := "/mgmt/tm/ltm/ifile/" + MangleFullPath(fullPath)

	mux, client := NewUnitTestServer(t)

	mux.HandleFunc("/mgmt/tm/ltm/ifile", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, fullPath)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"%s","fileName":"/Common/my-sys-ifile"}`, name, fullPath)
		case http.MethodPut:
			_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, fullPath)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	r := resourceBigipLtmIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"partition": "Common",
		"file_name": "/Common/my-sys-ifile",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmIfileCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, fullPath, d.Id())
	require.Equal(t, "/Common/my-sys-ifile", d.Get("file_name"))
	require.Equal(t, fullPath, d.Get("full_path"))

	// Read
	diags = resourceBigipLtmIfileRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, name, d.Get("name"))
	require.Equal(t, "Common", d.Get("partition"))
	require.Equal(t, "/Common/my-sys-ifile", d.Get("file_name"))
	require.Equal(t, fullPath, d.Get("full_path"))

	// Update
	require.NoError(t, d.Set("file_name", "/Common/other-sys-ifile"))
	diags = resourceBigipLtmIfileUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	// Delete
	diags = resourceBigipLtmIfileDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmIfileCreate_SubPath(t *testing.T) {
	name := "testfile"
	subPath := "subdir"
	fullPath := "/Common/subdir/testfile"

	mux, client := NewUnitTestServer(t)

	mux.HandleFunc("/mgmt/tm/ltm/ifile", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, fullPath)
	})
	mux.HandleFunc("/mgmt/tm/ltm/ifile/"+MangleFullPath(fullPath), func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"%s","fileName":"/Common/my-sys-ifile"}`, name, fullPath)
	})

	r := resourceBigipLtmIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"partition": "Common",
		"sub_path":  subPath,
		"file_name": "/Common/my-sys-ifile",
	}, "")

	diags := resourceBigipLtmIfileCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, fullPath, d.Id())
}

func TestUnitLtmIfileCreate_Error(t *testing.T) {
	mux, client := NewUnitTestServer(t)

	mux.HandleFunc("/mgmt/tm/ltm/ifile", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})

	r := resourceBigipLtmIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testfile",
		"partition": "Common",
		"file_name": "/Common/my-sys-ifile",
	}, "")

	diags := resourceBigipLtmIfileCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmIfileRead_NotFound(t *testing.T) {
	name := "testfile"
	fullPath := "/Common/" + name
	mangled := "/mgmt/tm/ltm/ifile/" + MangleFullPath(fullPath)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})

	r := resourceBigipLtmIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"partition": "Common",
		"file_name": "/Common/my-sys-ifile",
	}, fullPath)

	// A 404 means the iFile no longer exists on the device; Read clears the
	// resource ID rather than returning an error.
	diags := resourceBigipLtmIfileRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Equal(t, "", d.Id())
}

func TestUnitLtmIfileUpdate_Error(t *testing.T) {
	name := "testfile"
	fullPath := "/Common/" + name
	mangled := "/mgmt/tm/ltm/ifile/" + MangleFullPath(fullPath)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})

	r := resourceBigipLtmIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"partition": "Common",
		"file_name": "/Common/my-sys-ifile",
	}, fullPath)

	diags := resourceBigipLtmIfileUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmIfileDelete_Error(t *testing.T) {
	name := "testfile"
	fullPath := "/Common/" + name
	mangled := "/mgmt/tm/ltm/ifile/" + MangleFullPath(fullPath)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})

	r := resourceBigipLtmIfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"partition": "Common",
		"file_name": "/Common/my-sys-ifile",
	}, fullPath)

	diags := resourceBigipLtmIfileDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Equal(t, fullPath, d.Id())
}
