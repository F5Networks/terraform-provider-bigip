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

func TestResourceBigipSysOcspSchema(t *testing.T) {
	r := resourceBigipSysOcsp()

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
	if r.Schema["concurrent_connections_limit"].Default != 50 {
		t.Errorf("Expected 'concurrent_connections_limit' to default to 50, got %v", r.Schema["concurrent_connections_limit"].Default)
	}
	if r.Schema["sign_hash"].Default != "sha256" {
		t.Errorf("Expected 'sign_hash' to default to 'sha256', got %v", r.Schema["sign_hash"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitSysOcspCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-ocsp"
	fqdn := "~Common~test-ocsp"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/crypto/cert-validator/ocsp", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"test-ocsp","fullPath":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/sys/crypto/cert-validator/ocsp/"+fqdn, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"test-ocsp","fullPath":"%s","proxyServerPool":"/Common/my-pool","routeDomain":"0","responderUrl":"http://ocsp.example.com","concurrentConnectionsLimit":50,"connectionTimeout":8,"clockSkew":300,"statusAge":0,"strictRespCertCheck":"enabled","cacheTimeout":"indefinite","cacheErrorTimeout":3600,"signHash":"sha256"}`, name)
		case http.MethodPut:
			_, _ = fmt.Fprintf(w, `{"name":"test-ocsp","fullPath":"%s"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":              name,
		"proxy_server_pool": "/Common/my-pool",
		"responder_url":     "http://ocsp.example.com",
	}, "")

	ctx := context.Background()

	diags := resourceBigipSysOcspCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipSysOcspRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "/Common/my-pool", d.Get("proxy_server_pool").(string))

	diags = resourceBigipSysOcspUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSysOcspDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitSysOcspCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/crypto/cert-validator/ocsp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/test-ocsp",
	}, "")

	diags := resourceBigipSysOcspCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysOcspRead_InvalidID(t *testing.T) {
	client := NewUnitTestClient("http://localhost")
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "invalid-id-no-slash",
	}, "invalid-id-no-slash")

	diags := resourceBigipSysOcspRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysOcspRead_Error(t *testing.T) {
	fqdn := "~Common~test-ocsp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/crypto/cert-validator/ocsp/"+fqdn, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/test-ocsp",
	}, "/Common/test-ocsp")

	diags := resourceBigipSysOcspRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysOcspUpdate_InvalidID(t *testing.T) {
	client := NewUnitTestClient("http://localhost")
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "invalid-id-no-slash",
	}, "invalid-id-no-slash")

	diags := resourceBigipSysOcspUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysOcspUpdate_Error(t *testing.T) {
	fqdn := "~Common~test-ocsp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/crypto/cert-validator/ocsp/"+fqdn, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/test-ocsp",
	}, "/Common/test-ocsp")

	diags := resourceBigipSysOcspUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysOcspDelete_InvalidID(t *testing.T) {
	client := NewUnitTestClient("http://localhost")
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "invalid-id-no-slash",
	}, "invalid-id-no-slash")

	diags := resourceBigipSysOcspDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysOcspDelete_Error(t *testing.T) {
	fqdn := "~Common~test-ocsp"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/crypto/cert-validator/ocsp/"+fqdn, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/test-ocsp",
	}, "/Common/test-ocsp")

	diags := resourceBigipSysOcspDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// populateOcspConfig / setOcspStateData
// ---------------------------------------------------------------------

func TestPopulateOcspConfig(t *testing.T) {
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                         "/Common/test-ocsp",
		"proxy_server_pool":            "/Common/my-pool",
		"dns_resolver":                 "/Common/my-resolver",
		"route_domain":                 "0",
		"concurrent_connections_limit": 100,
		"responder_url":                "http://ocsp.example.com",
		"connection_timeout":           10,
		"clock_skew":                   120,
		"status_age":                   5,
		"strict_resp_cert_check":       "disabled",
		"cache_timeout":                "300",
		"cache_error_timeout":          1800,
		"signer_cert":                  "/Common/signer-cert",
		"signer_key":                   "/Common/signer-key",
		"passphrase":                   "secret",
		"sign_hash":                    "sha1",
	}, "")

	ocsp := &bigip.OCSP{Name: "test-ocsp"}
	populateOcspConfig(ocsp, d)

	require.Equal(t, "/Common/my-pool", ocsp.ProxyServerPool)
	require.Equal(t, "/Common/my-resolver", ocsp.DnsResolver)
	require.Equal(t, "0", ocsp.RouteDomain)
	require.Equal(t, int64(100), ocsp.ConcurrentConnectionsLimit)
	require.Equal(t, "http://ocsp.example.com", ocsp.ResponderUrl)
	require.Equal(t, int64(10), ocsp.ConnectionTimeout)
	require.Equal(t, int64(120), ocsp.ClockSkew)
	require.Equal(t, int64(5), ocsp.StatusAge)
	require.Equal(t, "disabled", ocsp.StrictRespCertCheck)
	require.Equal(t, "300", ocsp.CacheTimeout)
	require.Equal(t, int64(1800), ocsp.CacheErrorTimeout)
	require.Equal(t, "/Common/signer-cert", ocsp.SignerCert)
	require.Equal(t, "/Common/signer-key", ocsp.SignerKey)
	require.Equal(t, "secret", ocsp.Passphrase)
	require.Equal(t, "sha1", ocsp.SignHash)
}

func TestSetOcspStateData_DnsResolverBranch(t *testing.T) {
	r := resourceBigipSysOcsp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/test-ocsp",
	}, "")

	ocsp := &bigip.OCSP{
		Name:                       "test-ocsp",
		DnsResolver:                "/Common/my-resolver",
		RouteDomain:                "0",
		ResponderUrl:               "http://ocsp.example.com",
		TrustedResponders:          "/Common/trusted",
		SignerCert:                 "/Common/signer-cert",
		SignerKey:                  "/Common/signer-key",
		ConcurrentConnectionsLimit: 50,
		ClockSkew:                  300,
		StatusAge:                  0,
		CacheTimeout:               "indefinite",
		CacheErrorTimeout:          3600,
		ConnectionTimeout:          8,
		StrictRespCertCheck:        "enabled",
		SignHash:                   "sha256",
	}
	setOcspStateData(d, ocsp)

	require.Equal(t, "/Common/my-resolver", d.Get("dns_resolver").(string))
	require.Equal(t, "", d.Get("proxy_server_pool").(string))
	require.Equal(t, "/Common/trusted", d.Get("trusted_responders").(string))
	require.Equal(t, "/Common/signer-cert", d.Get("signer_cert").(string))
	require.Equal(t, "/Common/signer-key", d.Get("signer_key").(string))
}
