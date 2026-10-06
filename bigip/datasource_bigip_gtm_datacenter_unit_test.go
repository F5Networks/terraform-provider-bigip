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

func TestDataSourceBigipGtmDatacenterSchema(t *testing.T) {
	r := dataSourceBigipGtmDatacenter()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	for _, field := range []string{"name", "partition"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}

	for _, field := range []string{"description", "contact", "enabled", "location", "prober_fallback", "prober_preference"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Computed {
			t.Errorf("Expected field '%s' to be computed", field)
		}
	}
}

func TestDataSourceBigipGtmDatacenterReadSuccess(t *testing.T) {
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath("/Common/test-dc")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"name":"test-dc","partition":"Common","contact":"jane","description":"a datacenter","location":"NA","proberFallback":"any-available","proberPreference":"inside-datacenter"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := dataSourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "test-dc",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipGtmDatacenterRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "/Common/test-dc", d.Id())
	assert.Equal(t, "test-dc", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))
	assert.Equal(t, "jane", d.Get("contact"))
	assert.Equal(t, "a datacenter", d.Get("description"))
	assert.Equal(t, "NA", d.Get("location"))
	assert.Equal(t, "any-available", d.Get("prober_fallback"))
	assert.Equal(t, "inside-datacenter", d.Get("prober_preference"))
	assert.Equal(t, true, d.Get("enabled"))
}

func TestDataSourceBigipGtmDatacenterReadDisabled(t *testing.T) {
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath("/Common/disabled-dc")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"disabled-dc","partition":"Common","disabled":true}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := dataSourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "disabled-dc",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipGtmDatacenterRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, false, d.Get("enabled"))
}

func TestDataSourceBigipGtmDatacenterReadError(t *testing.T) {
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath("/Common/broken-dc")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := dataSourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "broken-dc",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipGtmDatacenterRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestDataSourceBigipGtmDatacenterReadNotFound(t *testing.T) {
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath("/Common/missing-dc")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := dataSourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "missing-dc",
		"partition": "Common",
	}, "")

	// Same as the resource: a 404 always yields a non-nil error from
	// getForEntity, so the "datacenter == nil" not-found branch never
	// actually executes via the real client -- the error branch (returning
	// diag.FromErr(err)) is what fires.
	diags := dataSourceBigipGtmDatacenterRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}
