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

func TestDataSourceBigipNetSelfSchema(t *testing.T) {
	r := dataSourceBigipNetSelf()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	name, ok := r.Schema["name"]
	require.True(t, ok, "expected 'name' field to exist in schema")
	assert.True(t, name.Required)

	partition, ok := r.Schema["partition"]
	require.True(t, ok, "expected 'partition' field to exist in schema")
	assert.Equal(t, "Common", partition.Default)

	for _, field := range []string{"full_path", "address", "vlan", "traffic_group", "floating", "unit", "allow_service"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "expected field '%s' to exist in schema", field)
		assert.True(t, s.Computed, "expected field '%s' to be computed", field)
	}
}

func TestDataSourceBigipNetSelfReadByNameSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self/~Common~external-self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"external-self","partition":"Common","fullPath":"/Common/external-self","address":"10.10.10.1/24","vlan":"/Common/external","trafficGroup":"/Common/traffic-group-local-only","floating":"false","unit":1,"allowService":"all"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetSelf(), map[string]interface{}{
		"name": "external-self",
	})

	diags := dataSourceBigipNetSelfRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/external-self", d.Get("name"))
	assert.Equal(t, "/Common/external-self", d.Get("full_path"))
	assert.Equal(t, "10.10.10.1/24", d.Get("address"))
	assert.Equal(t, "/Common/external", d.Get("vlan"))
	assert.Equal(t, "/Common/traffic-group-local-only", d.Get("traffic_group"))
	assert.Equal(t, "false", d.Get("floating"))
	assert.Equal(t, 1, d.Get("unit"))

	allowService := d.Get("allow_service").([]interface{})
	require.Len(t, allowService, 1)
	assert.Equal(t, "all", allowService[0])

	assert.Equal(t, "/Common/external-self", d.Id())
}

func TestDataSourceBigipNetSelfReadByFullPathName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self/~Common~external-self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"external-self","partition":"Common","fullPath":"/Common/external-self","address":"10.10.10.1/24","vlan":"/Common/external","allowService":null}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetSelf(), map[string]interface{}{
		"name": "/Common/external-self",
	})

	diags := dataSourceBigipNetSelfRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/external-self", d.Get("full_path"))

	allowService := d.Get("allow_service").([]interface{})
	require.Len(t, allowService, 1)
	assert.Equal(t, "none", allowService[0])
}

func TestDataSourceBigipNetSelfReadAllowServiceList(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self/~Common~restricted-self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"restricted-self","partition":"Common","fullPath":"/Common/restricted-self","address":"10.10.30.1/24","vlan":"/Common/restricted","allowService":["tcp:443","tcp:22"]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetSelf(), map[string]interface{}{
		"name": "restricted-self",
	})

	diags := dataSourceBigipNetSelfRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	allowService := d.Get("allow_service").([]interface{})
	require.Len(t, allowService, 2)
	assert.Equal(t, "tcp:443", allowService[0])
	assert.Equal(t, "tcp:22", allowService[1])
}

func TestDataSourceBigipNetSelfReadAllowServiceEmptyString(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self/~Common~emptystring-self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"emptystring-self","partition":"Common","fullPath":"/Common/emptystring-self","address":"10.10.40.1/24","vlan":"/Common/emptystring","allowService":""}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetSelf(), map[string]interface{}{
		"name": "emptystring-self",
	})

	diags := dataSourceBigipNetSelfRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	allowService := d.Get("allow_service").([]interface{})
	require.Len(t, allowService, 1, "an empty-string allowService should normalize to [\"none\"], not [\"\"]")
	assert.Equal(t, "none", allowService[0])
}

func TestDataSourceBigipNetSelfReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self/~Common~missing", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetSelf(), map[string]interface{}{
		"name": "missing",
	})

	diags := dataSourceBigipNetSelfRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the self IP does not exist")
	assert.Contains(t, diags[0].Summary, "not found")
}

func TestDataSourceBigipNetSelfReadLookupError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self/~Common~external-self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetSelf(), map[string]interface{}{
		"name": "external-self",
	})

	diags := dataSourceBigipNetSelfRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the self IP lookup call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving self IP")
}
