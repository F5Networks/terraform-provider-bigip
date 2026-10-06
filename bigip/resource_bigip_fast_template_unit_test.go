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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipFastTemplateSchema(t *testing.T) {
	r := resourceBigipFastTemplate()

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

	sourceSchema, ok := r.Schema["source"]
	if !ok {
		t.Fatal("Expected field 'source' to exist in schema")
	}
	if !sourceSchema.Required {
		t.Error("Expected field 'source' to be required")
	}
	if !sourceSchema.ForceNew {
		t.Error("Expected field 'source' to be ForceNew")
	}

	md5Schema, ok := r.Schema["md5_hash"]
	if !ok {
		t.Fatal("Expected field 'md5_hash' to exist in schema")
	}
	if !md5Schema.Required {
		t.Error("Expected field 'md5_hash' to be required")
	}
}

// ---------------------------------------------------------------------
// resourceBigipFastCreate / Read / Update / Delete
// ---------------------------------------------------------------------

// fastTemplateZipPath returns the path to the real example FAST template
// archive committed under examples/fast, used as the "source" file so
// os.OpenFile in resourceBigipFastCreate has a real file to read.
func fastTemplateZipPath(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	require.NoError(t, err)
	path := filepath.Join(wd, "..", "examples", "fast", "foo_template.zip")
	_, err = os.Stat(path)
	require.NoError(t, err, "expected example FAST template zip to exist at %s", path)
	return path
}

func TestUnitFastTemplateCreateReadDelete(t *testing.T) {
	templateName := "foo_template"

	mux := http.NewServeMux()
	var sawUpload, sawAddTemplateSet, sawGet, sawDelete bool

	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/foo_template.zip", func(w http.ResponseWriter, r *http.Request) {
		sawUpload = true
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":644,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/shared/fast/templatesets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sawAddTemplateSet = true
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, templateName)
		}
	})
	mux.HandleFunc("/mgmt/shared/fast/templatesets/"+templateName, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sawGet = true
			_, _ = fmt.Fprintf(w, `{"name":"%s","hash":"abc123","supported":true,"enabled":true}`, templateName)
		case http.MethodDelete:
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call

	r := resourceBigipFastTemplate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"source":   fastTemplateZipPath(t),
		"md5_hash": "89011331d11ac8bac2a1ad3235f38c80",
	}, "")

	ctx := context.Background()

	diags := resourceBigipFastCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawUpload, "expected the template zip to be uploaded")
	assert.True(t, sawAddTemplateSet, "expected a POST to add the template set")
	assert.True(t, sawGet, "expected a GET during the tail Read")
	assert.Equal(t, templateName, d.Id())
	assert.Equal(t, templateName, d.Get("name"))

	diags = resourceBigipFastRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, templateName, d.Get("name"))

	diags = resourceBigipFastDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete, "expected a DELETE request during delete")
	assert.Equal(t, "", d.Id())
}

func TestUnitFastTemplateCreate_ExplicitName(t *testing.T) {
	// When "name" is set explicitly it takes precedence over the basename of
	// "source" (with the .zip suffix stripped).
	templateName := "my-custom-name"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/my-custom-name.zip", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":644,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/shared/fast/templatesets", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, templateName)
	})
	mux.HandleFunc("/mgmt/shared/fast/templatesets/"+templateName, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, templateName)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipFastTemplate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     templateName,
		"source":   fastTemplateZipPath(t),
		"md5_hash": "89011331d11ac8bac2a1ad3235f38c80",
	}, "")

	diags := resourceBigipFastCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, templateName, d.Id())
}

func TestUnitFastTemplateCreate_FileNotFound(t *testing.T) {
	client := NewUnitTestClient("http://127.0.0.1:0")
	r := resourceBigipFastTemplate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"source":   "/nonexistent/path/does-not-exist.zip",
		"md5_hash": "deadbeef",
	}, "")

	// The file-open error returns before any client method is called, so
	// the client here is never actually used.
	diags := resourceBigipFastCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastTemplateCreate_UploadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/foo_template.zip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"upload failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastTemplate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"source":   fastTemplateZipPath(t),
		"md5_hash": "89011331d11ac8bac2a1ad3235f38c80",
	}, "")

	diags := resourceBigipFastCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastTemplateCreate_AddTemplateSetError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/foo_template.zip", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":644,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/shared/fast/templatesets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"add template set failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastTemplate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"source":   fastTemplateZipPath(t),
		"md5_hash": "89011331d11ac8bac2a1ad3235f38c80",
	}, "")

	diags := resourceBigipFastCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// A 404 from GetTemplateSet during Read means the template set no longer
// exists on the device. Per standard Terraform provider convention this is
// not an error: Read clears the resource ID so the next plan recreates it.
func TestUnitFastTemplateRead_NotFoundClearsId(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/templatesets/missing-template", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastTemplate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"source":   fastTemplateZipPath(t),
		"md5_hash": "89011331d11ac8bac2a1ad3235f38c80",
	}, "missing-template")

	diags := resourceBigipFastRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Empty(t, d.Id())
}

func TestUnitFastTemplateUpdate_DelegatesToCreate(t *testing.T) {
	// resourceBigipFastUpdate is a direct passthrough to resourceBigipFastCreate.
	templateName := "foo_template"

	mux := http.NewServeMux()
	var sawUpload bool
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/foo_template.zip", func(w http.ResponseWriter, r *http.Request) {
		sawUpload = true
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":644,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/shared/fast/templatesets", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, templateName)
	})
	mux.HandleFunc("/mgmt/shared/fast/templatesets/"+templateName, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, templateName)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipFastTemplate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"source":   fastTemplateZipPath(t),
		"md5_hash": "89011331d11ac8bac2a1ad3235f38c80",
	}, templateName)

	diags := resourceBigipFastUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawUpload, "expected update to re-upload the template via Create")
}

func TestUnitFastTemplateDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/templatesets/bad-template", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastTemplate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"source":   fastTemplateZipPath(t),
		"md5_hash": "89011331d11ac8bac2a1ad3235f38c80",
	}, "bad-template")

	diags := resourceBigipFastDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
