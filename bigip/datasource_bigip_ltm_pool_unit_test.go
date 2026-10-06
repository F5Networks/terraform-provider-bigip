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

// ltmPoolDatasourceResourceData builds a *schema.ResourceData for the
// bigip_ltm_pool data source using schema.TestResourceDataRaw.
func ltmPoolDatasourceResourceData(t *testing.T, name, partition string) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipLtmPool(), map[string]interface{}{
		"name":      name,
		"partition": partition,
	})
}

// TestDataSourceBigipLtmPoolSchema exercises the schema definition of the
// bigip_ltm_pool data source without needing any BIG-IP connection.
func TestDataSourceBigipLtmPoolSchema(t *testing.T) {
	r := dataSourceBigipLtmPool()

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

	fullPathSchema, ok := r.Schema["full_path"]
	if !ok {
		t.Fatal("Expected field 'full_path' to exist in schema")
	}
	if !fullPathSchema.Computed {
		t.Error("Expected field 'full_path' to be Computed")
	}
	if fullPathSchema.Type != schema.TypeString {
		t.Errorf("Expected field 'full_path' to be TypeString, got %v", fullPathSchema.Type)
	}
}

// TestDataSourceBigipLtmPoolReadSuccess covers the Read happy path:
// GetPool succeeds and returns a non-nil pool, so name/partition/full_path
// are (re)populated from the API response and the resource ID is set to
// the pool's full path.
func TestDataSourceBigipLtmPoolReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"test-pool","partition":"Common","fullPath":"/Common/test-pool","loadBalancingMode":"round-robin"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := ltmPoolDatasourceResourceData(t, "test-pool", "Common")

	diags := dataSourceBigipLtmPoolRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/test-pool", d.Id())
	assert.Equal(t, "test-pool", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))
	assert.Equal(t, "/Common/test-pool", d.Get("full_path"))
}

// TestDataSourceBigipLtmPoolReadSubPartition covers Read with a pool that
// lives in a non-Common partition, ensuring the request path and returned
// fields are built/parsed correctly for that case too.
func TestDataSourceBigipLtmPoolReadSubPartition(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/pool/~MyPartition~sub-pool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"sub-pool","partition":"MyPartition","fullPath":"/MyPartition/sub-pool"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := ltmPoolDatasourceResourceData(t, "sub-pool", "MyPartition")

	diags := dataSourceBigipLtmPoolRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/MyPartition/sub-pool", d.Id())
	assert.Equal(t, "sub-pool", d.Get("name"))
	assert.Equal(t, "MyPartition", d.Get("partition"))
	assert.Equal(t, "/MyPartition/sub-pool", d.Get("full_path"))
}

// TestDataSourceBigipLtmPoolReadNotFound covers the branch where GetPool
// returns (nil, nil) -- i.e. the underlying API call succeeded with a 404
// and go-bigip.getForEntity reports the entity as absent rather than
// erroring. Read must surface this as an error diagnostic and leave the ID
// unset (it is explicitly cleared at the start of Read).
func TestDataSourceBigipLtmPoolReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~missing-pool", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := ltmPoolDatasourceResourceData(t, "missing-pool", "Common")
	d.SetId("preexisting-id")

	diags := dataSourceBigipLtmPoolRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when the pool is not found")
	assert.Contains(t, diags[0].Summary, "not found")
	assert.Empty(t, d.Id(), "resource ID should be cleared when the pool is not found")
}

// TestDataSourceBigipLtmPoolReadError covers the branch where GetPool
// itself returns an error (e.g. the BIG-IP is unreachable or returns a
// non-404 error response).
func TestDataSourceBigipLtmPoolReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~broken-pool", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := ltmPoolDatasourceResourceData(t, "broken-pool", "Common")
	d.SetId("preexisting-id")

	diags := dataSourceBigipLtmPoolRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when the API call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving pool")
	assert.Empty(t, d.Id(), "resource ID should be cleared when the API call fails")
}
