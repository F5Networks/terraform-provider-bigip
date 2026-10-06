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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitLtmPersistenceProfileSSLCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-ssl-pp"
	mangled := "/mgmt/tm/ltm/persistence/ssl/~Common~test-ssl-pp"

	mux := http.NewServeMux()
	var sawPost, sawPatch, sawDelete bool

	mux.HandleFunc("/mgmt/tm/ltm/persistence/ssl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		sawPost = true
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			sawPatch = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{
				"name":"%s",
				"defaultsFrom":"/Common/ssl",
				"matchAcrossPools":"disabled",
				"matchAcrossServices":"disabled",
				"matchAcrossVirtuals":"disabled",
				"mirror":"disabled",
				"overrideConnLimit":"disabled",
				"timeout":"300"
			}`, name)
		case http.MethodDelete:
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileSSL()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                  name,
		"defaults_from":         "/Common/ssl",
		"match_across_pools":    "disabled",
		"match_across_services": "disabled",
		"match_across_virtuals": "disabled",
		"mirror":                "disabled",
		"timeout":               300,
		"override_conn_limit":   "disabled",
		"app_service":           "none",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmPersistenceProfileSSLCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.True(t, sawPatch, "Create calls Update internally")
	assert.Equal(t, name, d.Id())

	diags = resourceBigipLtmPersistenceProfileSSLRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "/Common/ssl", d.Get("defaults_from"))
	assert.Equal(t, 300, d.Get("timeout"))

	diags = resourceBigipLtmPersistenceProfileSSLUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmPersistenceProfileSSLDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitLtmPersistenceProfileSSLUpdateZeroTimeout(t *testing.T) {
	// Exercises the "timeout == 0" else branch in Update, which builds the
	// SSLPersistenceProfile without a Timeout field.
	name := "/Common/zero-timeout-ssl-pp"
	mangled := "/mgmt/tm/ltm/persistence/ssl/~Common~zero-timeout-ssl-pp"

	mux := http.NewServeMux()
	var sawPatch bool
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			sawPatch = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/ssl","timeout":"0"}`, name)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileSSL()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/ssl",
		"timeout":       0,
	}, name)

	diags := resourceBigipLtmPersistenceProfileSSLUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPatch)
}

func TestUnitLtmPersistenceProfileSSLCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/ssl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprint(w, `{"code":409,"message":"already exists"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileSSL()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/dup-ssl-pp",
		"defaults_from": "/Common/ssl",
	}, "")

	diags := resourceBigipLtmPersistenceProfileSSLCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileSSLReadError(t *testing.T) {
	mangled := "/mgmt/tm/ltm/persistence/ssl/~Common~broken-ssl-pp"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileSSL()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/broken-ssl-pp",
	}, "/Common/broken-ssl-pp")

	// getForEntity always returns a non-nil error on a non-2xx response
	// (including 404), so the "pp == nil" not-found branch in Read is
	// unreachable via the real client; this is a hard error path.
	diags := resourceBigipLtmPersistenceProfileSSLRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileSSLUpdateErrorThenDeleteSucceeds(t *testing.T) {
	mangled := "/mgmt/tm/ltm/persistence/ssl/~Common~err-ssl-pp"

	mux := http.NewServeMux()
	var sawDelete bool
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
		case http.MethodDelete:
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileSSL()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/err-ssl-pp",
		"defaults_from": "/Common/ssl",
		"timeout":       120,
	}, "/Common/err-ssl-pp")

	diags := resourceBigipLtmPersistenceProfileSSLUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
	assert.True(t, sawDelete, "Update should delete the profile when Modify fails")
}

func TestUnitLtmPersistenceProfileSSLUpdateZeroTimeoutErrorThenDeleteFails(t *testing.T) {
	mangled := "/mgmt/tm/ltm/persistence/ssl/~Common~err2-ssl-pp"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileSSL()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/err2-ssl-pp",
		"defaults_from": "/Common/ssl",
		"timeout":       0,
	}, "/Common/err2-ssl-pp")

	diags := resourceBigipLtmPersistenceProfileSSLUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileSSLDeleteError(t *testing.T) {
	mangled := "/mgmt/tm/ltm/persistence/ssl/~Common~del-err-ssl-pp"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileSSL()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/del-err-ssl-pp",
	}, "/Common/del-err-ssl-pp")

	diags := resourceBigipLtmPersistenceProfileSSLDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileSSLReadNonNumericTimeout(t *testing.T) {
	mangled := "/mgmt/tm/ltm/persistence/ssl/~Common~bad-timeout-ssl-pp"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
			"name":"/Common/bad-timeout-ssl-pp",
			"defaultsFrom":"/Common/ssl",
			"timeout":"not-a-number"
		}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileSSL()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "/Common/bad-timeout-ssl-pp")

	diags := resourceBigipLtmPersistenceProfileSSLRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "/Common/ssl", d.Get("defaults_from"))
	assert.Equal(t, 0, d.Get("timeout"))
}
