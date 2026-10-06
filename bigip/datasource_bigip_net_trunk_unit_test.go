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

func TestDataSourceBigipNetTrunkSchema(t *testing.T) {
	r := dataSourceBigipNetTrunk()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	name, ok := r.Schema["name"]
	require.True(t, ok, "expected 'name' field to exist in schema")
	assert.True(t, name.Required)

	for _, field := range []string{"interfaces", "lacp", "lacp_mode", "lacp_timeout", "distribution_hash", "link_select_policy", "bandwidth", "trunk_id", "stp", "type", "working_member_count"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "expected field '%s' to exist in schema", field)
		assert.True(t, s.Computed, "expected field '%s' to be computed", field)
	}
}

func TestDataSourceBigipNetTrunkReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/trunk/lag1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"lag1","fullPath":"lag1","interfaces":["1.1","1.2"],"lacp":"enabled","lacpMode":"active","lacpTimeout":"long","distributionHash":"src-dst-ipport","linkSelectPolicy":"auto","bandwidth":20000,"id":1,"stp":"enabled","type":"lacp","workingMbrCount":2}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetTrunk(), map[string]interface{}{
		"name": "lag1",
	})

	diags := dataSourceBigipNetTrunkRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "lag1", d.Get("name"))
	assert.Equal(t, "active", d.Get("lacp_mode"))
	assert.Equal(t, "long", d.Get("lacp_timeout"))
	assert.Equal(t, "enabled", d.Get("lacp"))
	assert.Equal(t, "src-dst-ipport", d.Get("distribution_hash"))
	assert.Equal(t, "auto", d.Get("link_select_policy"))
	assert.Equal(t, 20000, d.Get("bandwidth"))
	assert.Equal(t, 1, d.Get("trunk_id"))
	assert.Equal(t, "enabled", d.Get("stp"))
	assert.Equal(t, "lacp", d.Get("type"))
	assert.Equal(t, 2, d.Get("working_member_count"))

	members := d.Get("interfaces").([]interface{})
	require.Len(t, members, 2)
	assert.Equal(t, "1.1", members[0])
	assert.Equal(t, "1.2", members[1])

	assert.Equal(t, "lag1", d.Id())
}

func TestDataSourceBigipNetTrunkReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/trunk/missing", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"01020036:3: The requested Trunk (missing) was not found."}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetTrunk(), map[string]interface{}{
		"name": "missing",
	})

	diags := dataSourceBigipNetTrunkRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the trunk does not exist")
	assert.Contains(t, diags[0].Summary, "trunk not found")
}

func TestDataSourceBigipNetTrunkReadLookupError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/trunk/lag1", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetTrunk(), map[string]interface{}{
		"name": "lag1",
	})

	diags := dataSourceBigipNetTrunkRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the trunk lookup call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving trunk")
}
