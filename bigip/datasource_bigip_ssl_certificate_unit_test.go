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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sslCertificateDatasourceResourceData builds a *schema.ResourceData for the
// bigip_ssl_certificate data source using schema.TestResourceDataRaw.
func sslCertificateDatasourceResourceData(t *testing.T, name, partition string) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipSslCertificate(), map[string]interface{}{
		"name":      name,
		"partition": partition,
	})
}

// TestDataSourceBigipSslCertificateSchema exercises the schema definition of
// the bigip_ssl_certificate data source without needing any BIG-IP
// connection.
func TestDataSourceBigipSslCertificateSchema(t *testing.T) {
	r := dataSourceBigipSslCertificate()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	requiredFields := []string{"name", "partition"}
	for _, field := range requiredFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
		if s.Type != schema.TypeString {
			t.Errorf("Expected field '%s' to be TypeString, got %v", field, s.Type)
		}
	}

	certSchema, ok := r.Schema["certificate"]
	if !ok {
		t.Fatal("Expected field 'certificate' to exist in schema")
	}
	if !certSchema.Computed {
		t.Error("Expected field 'certificate' to be Computed")
	}
	if certSchema.Type != schema.TypeString {
		t.Errorf("Expected field 'certificate' to be TypeString, got %v", certSchema.Type)
	}
}

// TestDataSourceBigipSslCertificateReadSuccess covers the Read happy path:
// GetCertificate succeeds and returns a non-nil certificate, so
// name/partition are (re)populated from the API response and the resource
// ID is set to the certificate's full path.
//
// Note: unlike name/partition, the "certificate" schema field is never
// populated by Read (it's Computed but Read has no d.Set call for it), so
// this test asserts it remains empty rather than asserting a value that
// the implementation doesn't actually set.
func TestDataSourceBigipSslCertificateReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~Common~test-cert", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"test-cert","partition":"Common","fullPath":"/Common/test-cert","issuer":"/Common/ca-cert"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := sslCertificateDatasourceResourceData(t, "test-cert", "Common")

	diags := dataSourceBigipSslCertificateRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/test-cert", d.Id())
	assert.Equal(t, "test-cert", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))
	assert.Empty(t, d.Get("certificate"), "Read does not populate the certificate field")
}

// TestDataSourceBigipSslCertificateReadSubPartition covers Read with a
// certificate that lives in a non-Common partition, ensuring the request
// path and returned fields are built/parsed correctly for that case too.
func TestDataSourceBigipSslCertificateReadSubPartition(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~MyPartition~sub-cert", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"sub-cert","partition":"MyPartition","fullPath":"/MyPartition/sub-cert"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := sslCertificateDatasourceResourceData(t, "sub-cert", "MyPartition")

	diags := dataSourceBigipSslCertificateRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/MyPartition/sub-cert", d.Id())
	assert.Equal(t, "sub-cert", d.Get("name"))
	assert.Equal(t, "MyPartition", d.Get("partition"))
}

// TestDataSourceBigipSslCertificateReadNotFound covers the branch where
// GetCertificate returns (nil, nil) -- i.e. the underlying API call
// succeeded with a 404 and go-bigip.getForEntity reports the entity as
// absent rather than erroring. Read must surface this as an error
// diagnostic and leave the ID unset (it is explicitly cleared at the start
// of Read).
func TestDataSourceBigipSslCertificateReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~Common~missing-cert", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := sslCertificateDatasourceResourceData(t, "missing-cert", "Common")
	d.SetId("preexisting-id")

	diags := dataSourceBigipSslCertificateRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when the certificate is not found")
	assert.Contains(t, diags[0].Summary, "not found")
	assert.Empty(t, d.Id(), "resource ID should be cleared when the certificate is not found")
}

// TestDataSourceBigipSslCertificateReadError covers the branch where
// GetCertificate itself returns an error (e.g. the BIG-IP is unreachable or
// returns a non-404 error response).
func TestDataSourceBigipSslCertificateReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/file/ssl-cert/~Common~broken-cert", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := sslCertificateDatasourceResourceData(t, "broken-cert", "Common")
	d.SetId("preexisting-id")

	diags := dataSourceBigipSslCertificateRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when the API call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving certificate")
	assert.Empty(t, d.Id(), "resource ID should be cleared when the API call fails")
}
