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

func TestUnitLtmPersistenceProfileCookieCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-cookie-pp"
	mangled := "/mgmt/tm/ltm/persistence/cookie/~Common~test-cookie-pp"

	mux := http.NewServeMux()
	var sawPost, sawPatch, sawDelete bool

	mux.HandleFunc("/mgmt/tm/ltm/persistence/cookie", func(w http.ResponseWriter, r *http.Request) {
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
				"defaultsFrom":"/Common/cookie",
				"matchAcrossPools":"disabled",
				"matchAcrossServices":"disabled",
				"matchAcrossVirtuals":"disabled",
				"mirror":"disabled",
				"overrideConnLimit":"disabled",
				"httponly":"enabled",
				"expiration":"0",
				"alwaysSend":"disabled",
				"hashLength":0,
				"hashOffset":0,
				"method":"insert",
				"timeout":"180",
				"appService":"none",
				"cookieEncryption":"disabled",
				"cookieEncryptionPassphrase":"secret",
				"cookieName":"my-cookie"
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

	r := resourceBigipLtmPersistenceProfileCookie()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                         name,
		"defaults_from":                "/Common/cookie",
		"match_across_pools":           "disabled",
		"match_across_services":        "disabled",
		"match_across_virtuals":        "disabled",
		"mirror":                       "disabled",
		"method":                       "insert",
		"timeout":                      180,
		"override_conn_limit":          "disabled",
		"always_send":                  "disabled",
		"cookie_encryption":            "disabled",
		"cookie_encryption_passphrase": "secret",
		"cookie_name":                  "my-cookie",
		"expiration":                   "0",
		"hash_length":                  0,
		"hash_offset":                  0,
		"httponly":                     "enabled",
		"app_service":                  "none",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmPersistenceProfileCookieCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.True(t, sawPatch, "Create calls Update internally")
	assert.Equal(t, name, d.Id())

	diags = resourceBigipLtmPersistenceProfileCookieRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "/Common/cookie", d.Get("defaults_from"))
	assert.Equal(t, "my-cookie", d.Get("cookie_name"))
	assert.Equal(t, "insert", d.Get("method"))
	assert.Equal(t, 180, d.Get("timeout"))

	diags = resourceBigipLtmPersistenceProfileCookieUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmPersistenceProfileCookieDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitLtmPersistenceProfileCookieCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/persistence/cookie", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprint(w, `{"code":409,"message":"already exists"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileCookie()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/dup-cookie-pp",
		"defaults_from": "/Common/cookie",
	}, "")

	diags := resourceBigipLtmPersistenceProfileCookieCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileCookieReadError(t *testing.T) {
	mangled := "/mgmt/tm/ltm/persistence/cookie/~Common~broken-cookie-pp"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileCookie()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/broken-cookie-pp",
	}, "/Common/broken-cookie-pp")

	// getForEntity always returns a non-nil error on a non-2xx response
	// (including 404), so the "pp == nil" not-found branch in Read is
	// unreachable via the real client; this is a hard error path.
	diags := resourceBigipLtmPersistenceProfileCookieRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileCookieUpdateErrorThenDeleteSucceeds(t *testing.T) {
	mangled := "/mgmt/tm/ltm/persistence/cookie/~Common~err-cookie-pp"

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

	r := resourceBigipLtmPersistenceProfileCookie()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/err-cookie-pp",
		"defaults_from": "/Common/cookie",
	}, "/Common/err-cookie-pp")

	diags := resourceBigipLtmPersistenceProfileCookieUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
	assert.True(t, sawDelete, "Update should delete the profile when Modify fails")
}

func TestUnitLtmPersistenceProfileCookieUpdateErrorThenDeleteFails(t *testing.T) {
	mangled := "/mgmt/tm/ltm/persistence/cookie/~Common~err2-cookie-pp"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileCookie()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/err2-cookie-pp",
		"defaults_from": "/Common/cookie",
	}, "/Common/err2-cookie-pp")

	diags := resourceBigipLtmPersistenceProfileCookieUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileCookieDeleteError(t *testing.T) {
	mangled := "/mgmt/tm/ltm/persistence/cookie/~Common~del-err-cookie-pp"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileCookie()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/del-err-cookie-pp",
	}, "/Common/del-err-cookie-pp")

	diags := resourceBigipLtmPersistenceProfileCookieDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmPersistenceProfileCookieReadNoOptionalGetOk(t *testing.T) {
	// Exercises the Read branches where app_service, cookie_encryption,
	// cookie_encryption_passphrase, and cookie_name are absent from state
	// (d.GetOk returns false), so those Set calls are skipped.
	mangled := "/mgmt/tm/ltm/persistence/cookie/~Common~minimal-cookie-pp"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
			"name":"/Common/minimal-cookie-pp",
			"defaultsFrom":"/Common/cookie",
			"timeout":"not-a-number"
		}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipLtmPersistenceProfileCookie()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "/Common/minimal-cookie-pp")

	diags := resourceBigipLtmPersistenceProfileCookieRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "/Common/cookie", d.Get("defaults_from"))
	// timeout was non-numeric, so strconv.Atoi failed and Set was skipped;
	// the field keeps its zero value.
	assert.Equal(t, 0, d.Get("timeout"))
}
