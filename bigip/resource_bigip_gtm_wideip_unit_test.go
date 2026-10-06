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

func TestResourceBigipGtmWideipSchema(t *testing.T) {
	r := resourceBigipGtmWideip()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.CreateContext == nil || r.ReadContext == nil || r.UpdateContext == nil || r.DeleteContext == nil {
		t.Fatal("Expected all CRUD contexts to be defined")
	}
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	for _, field := range []string{"name", "type"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "Expected field '%s' to exist in schema", field)
		assert.True(t, s.Required)
		assert.True(t, s.ForceNew)
	}
}

func TestParseWideIPIDValid(t *testing.T) {
	parts := parseWideIPID("a:/Common/testwideip.local")
	require.NotNil(t, parts)
	assert.Equal(t, "a", parts["type"])
	assert.Equal(t, "Common", parts["partition"])
	assert.Equal(t, "testwideip.local", parts["name"])
}

func TestParseWideIPIDNoPartition(t *testing.T) {
	parts := parseWideIPID("a:justaname")
	require.NotNil(t, parts)
	assert.Equal(t, "a", parts["type"])
	assert.Equal(t, "Common", parts["partition"])
	assert.Equal(t, "justaname", parts["name"])
}

func TestParseWideIPIDInvalid(t *testing.T) {
	parts := parseWideIPID("no-colon-here")
	assert.Nil(t, parts)
}

func TestUnitGtmWideipCreateReadUpdateDelete(t *testing.T) {
	fullPath := "/Common/test-wideip.local"
	mangled := "/mgmt/tm/gtm/wideip/a/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	var updated bool
	var sawPost, sawPut, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/wideip/a", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":"test-wideip.local","partition":"Common","fullPath":"%s"}`, fullPath)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			sawPut = true
			updated = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		case http.MethodGet:
			if updated {
				_, _ = fmt.Fprintf(w, `{"name":"test-wideip.local","partition":"Common","fullPath":"%s","description":"updated","enabled":true,"minimalResponse":"enabled","poolLbMode":"round-robin","pools":[{"name":"/Common/pool1","order":0,"ratio":1}],"aliases":["alt.example.com"]}`, fullPath)
			} else {
				_, _ = fmt.Fprintf(w, `{"name":"test-wideip.local","partition":"Common","fullPath":"%s"}`, fullPath)
			}
		case http.MethodDelete:
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmWideip()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "test-wideip.local",
		"type":      "a",
		"partition": "Common",
	}, "")

	ctx := context.Background()

	diags := resourceBigipGtmWideipCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.Equal(t, "a:"+fullPath, d.Id())
	assert.True(t, sawPut, "Create calls Update internally")

	diags = resourceBigipGtmWideipRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "test-wideip.local", d.Get("name"))
	assert.Equal(t, "updated", d.Get("description"))
	pools := d.Get("pools").([]interface{})
	require.Len(t, pools, 1)
	assert.Equal(t, "/Common/pool1", pools[0].(map[string]interface{})["name"])

	require.NoError(t, d.Set("description", "second update"))
	diags = resourceBigipGtmWideipUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipGtmWideipDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmWideipCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/wideip/a", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprint(w, `{"code":409,"message":"the requested object already exists"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmWideip()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "dup-wideip.local",
		"type":      "a",
		"partition": "Common",
	}, "")

	diags := resourceBigipGtmWideipCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmWideipReadWithLogVerbosity(t *testing.T) {
	fullPath := "/Common/verbose-wideip.local"
	mangled := "/mgmt/tm/gtm/wideip/a/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"verbose-wideip.local","partition":"Common","fullPath":"%s","loadBalancingDecisionLogVerbosity":["pool-selection","pool-traversal"]}`, fullPath)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmWideip()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "verbose-wideip.local",
		"type":      "a",
		"partition": "Common",
	}, "a:"+fullPath)

	diags := resourceBigipGtmWideipRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	verbosity := d.Get("load_balancing_decision_log_verbosity").(*schema.Set)
	assert.Equal(t, 2, verbosity.Len())
}

func TestUnitGtmWideipReadFallbackToState(t *testing.T) {
	// Exercises the "fallback to getting from state" branch in Read, which
	// triggers when d.Id() doesn't parse as a valid "type:/partition/name".
	fullPath := "/Common/fallback-wideip.local"
	mangled := "/mgmt/tm/gtm/wideip/a/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"fallback-wideip.local","partition":"Common","fullPath":"%s"}`, fullPath)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmWideip()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "fallback-wideip.local",
		"type":      "a",
		"partition": "Common",
	}, "not-a-valid-id")

	diags := resourceBigipGtmWideipRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "fallback-wideip.local", d.Get("name"))
}

func TestUnitGtmWideipReadNotFound(t *testing.T) {
	mangled := "/mgmt/tm/gtm/wideip/a/" + MangleFullPath("/Common/missing-wideip.local")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmWideip()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "missing-wideip.local",
		"type":      "a",
		"partition": "Common",
	}, "a:/Common/missing-wideip.local")

	// A 404 means the wide IP no longer exists on the device; Read clears
	// the resource ID rather than returning an error.
	diags := resourceBigipGtmWideipRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Equal(t, "", d.Id())
}

func TestUnitGtmWideipUpdateError(t *testing.T) {
	mangled := "/mgmt/tm/gtm/wideip/a/" + MangleFullPath("/Common/err-wideip.local")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmWideip()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "err-wideip.local",
		"type":      "a",
		"partition": "Common",
	}, "a:/Common/err-wideip.local")

	diags := resourceBigipGtmWideipUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmWideipDeleteError(t *testing.T) {
	mangled := "/mgmt/tm/gtm/wideip/a/" + MangleFullPath("/Common/del-err-wideip.local")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmWideip()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "del-err-wideip.local",
		"type":      "a",
		"partition": "Common",
	}, "a:/Common/del-err-wideip.local")

	diags := resourceBigipGtmWideipDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmWideipImport(t *testing.T) {
	mangled := "/mgmt/tm/gtm/wideip/a/" + MangleFullPath("/Common/import-wideip.local")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"import-wideip.local","partition":"Common","fullPath":"/Common/import-wideip.local"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmWideip()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "a:/Common/import-wideip.local")

	results, err := resourceBigipGtmWideipImport(context.Background(), d, client)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "import-wideip.local", results[0].Get("name"))
	assert.Equal(t, "a", results[0].Get("type"))
}

func TestUnitGtmWideipImportInvalid(t *testing.T) {
	r := resourceBigipGtmWideip()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "not-a-valid-id")

	_, err := resourceBigipGtmWideipImport(context.Background(), d, nil)
	require.Error(t, err)
}
