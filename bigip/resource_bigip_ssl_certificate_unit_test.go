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

func TestResourceBigipSslCertificateSchema(t *testing.T) {
	r := resourceBigipSslCertificate()

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
	if r.Schema["partition"].Default != "Common" {
		t.Errorf("Expected 'partition' to default to 'Common', got %v", r.Schema["partition"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitSslCertificateCreateReadUpdateDelete(t *testing.T) {
	name := "testcert"
	mangled := "/mgmt/tm/sys/file/ssl-cert/" + MangleFullPath("/Common/"+name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
	})
	certHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	}
	mux.HandleFunc(mangled, certHandler)
	// ModifyCertificate (used by UpdateCertificate) appends a literal
	// "?expandSubcollections=true" path segment (joined with "/", not a
	// real query string), producing a subpath under the mangled cert path.
	mux.HandleFunc(mangled+"/", certHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call
	r := resourceBigipSslCertificate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-cert-content",
		"partition": "Common",
	}, "")

	ctx := context.Background()

	diags := resourceBigipSslCertificateCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipSslCertificateRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	diags = resourceBigipSslCertificateUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSslCertificateDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitSslCertificateCreate_WithOcspAndMonitoring(t *testing.T) {
	name := "testcert"
	mangled := "/mgmt/tm/sys/file/ssl-cert/" + MangleFullPath("/Common/"+name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
	})
	certHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
		}
	}
	mux.HandleFunc(mangled, certHandler)
	mux.HandleFunc(mangled+"/", certHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipSslCertificate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":            name,
		"content":         "fake-cert-content",
		"partition":       "Common",
		"monitoring_type": "certificate",
		"issuer_cert":     "/Common/issuer",
		"ocsp":            "/Common/my-ocsp",
	}, "")

	diags := resourceBigipSslCertificateCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
}

func TestUnitSslCertificateCreate_UploadError(t *testing.T) {
	name := "testcert"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"upload failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipSslCertificate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-cert-content",
		"partition": "Common",
	}, "")

	diags := resourceBigipSslCertificateCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSslCertificateRead_NotFound(t *testing.T) {
	name := "testcert"
	mangled := "/mgmt/tm/sys/file/ssl-cert/" + MangleFullPath("/Common/"+name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslCertificate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-cert-content",
		"partition": "Common",
	}, name)

	diags := resourceBigipSslCertificateRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSslCertificateRead_NoPartitionAndMonitoringType(t *testing.T) {
	// name lacking a "/" prefix and no partition exercises the "must be in
	// full_path format" warning branch (a log only, not an error); the mock
	// certificate response includes cert_validation_options so the
	// monitoring_type d.Set branch is also exercised.
	name := "/Common/testcert"
	mangled := "/mgmt/tm/sys/file/ssl-cert/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"testcert","partition":"Common","fullPath":"%s","certValidationOptions":["certificate"]}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslCertificate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testcert",
		"content":   "fake-cert-content",
		"partition": "",
	}, name)

	diags := resourceBigipSslCertificateRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "certificate", d.Get("monitoring_type").(string))
}

func TestUnitSslCertificateUpdate_Error(t *testing.T) {
	name := "testcert"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"upload failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslCertificate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-cert-content",
		"partition": "Common",
	}, name)

	diags := resourceBigipSslCertificateUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSslCertificateUpdate_WithOcspAndMonitoring(t *testing.T) {
	name := "testcert"
	mangled := "/mgmt/tm/sys/file/ssl-cert/" + MangleFullPath("/Common/"+name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"remainingByteCount":0,"usedChunks":{},"totalByteCount":10,"localFilePath":"","temporaryFilePath":"","generation":0,"lastUpdateMicros":0}`)
	})
	certHandler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","fullPath":"/Common/%s"}`, name, name)
	}
	mux.HandleFunc(mangled, certHandler)
	mux.HandleFunc(mangled+"/", certHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslCertificate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":            name,
		"content":         "fake-cert-content",
		"partition":       "Common",
		"monitoring_type": "certificate",
		"issuer_cert":     "/Common/issuer",
		"ocsp":            "/Common/my-ocsp",
	}, name)

	diags := resourceBigipSslCertificateUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitSslCertificateDelete_Error(t *testing.T) {
	name := "testcert"
	mangled := "/mgmt/tm/sys/file/ssl-cert/" + MangleFullPath("/Common/"+name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSslCertificate()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"content":   "fake-cert-content",
		"partition": "Common",
	}, name)

	diags := resourceBigipSslCertificateDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
