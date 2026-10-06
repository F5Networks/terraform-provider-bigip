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

func TestDataSourceBigipNetRouteSchema(t *testing.T) {
	r := dataSourceBigipNetRoute()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	name, ok := r.Schema["name"]
	require.True(t, ok, "expected 'name' field to exist in schema")
	assert.True(t, name.Required)

	partition, ok := r.Schema["partition"]
	require.True(t, ok, "expected 'partition' field to exist in schema")
	assert.Equal(t, "Common", partition.Default)

	for _, field := range []string{"full_path", "network", "gw", "description", "tm_interface", "mtu", "blackhole"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "expected field '%s' to exist in schema", field)
		assert.True(t, s.Computed, "expected field '%s' to be computed", field)
	}
}

func TestDataSourceBigipNetRouteReadByNameSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route/~Common~external-route", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"external-route","partition":"Common","fullPath":"/Common/external-route","network":"10.10.10.0/24","gw":"11.1.1.2","description":"test route","mtu":0,"blackhole":false}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetRoute(), map[string]interface{}{
		"name": "external-route",
	})

	diags := dataSourceBigipNetRouteRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/external-route", d.Get("name"))
	assert.Equal(t, "/Common/external-route", d.Get("full_path"))
	assert.Equal(t, "10.10.10.0/24", d.Get("network"))
	assert.Equal(t, "11.1.1.2", d.Get("gw"))
	assert.Equal(t, "test route", d.Get("description"))
	assert.Equal(t, 0, d.Get("mtu"))
	assert.Equal(t, false, d.Get("blackhole"))

	assert.Equal(t, "/Common/external-route", d.Id())
}

func TestDataSourceBigipNetRouteReadByFullPathName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route/~Common~external-route", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"external-route","partition":"Common","fullPath":"/Common/external-route","network":"10.10.10.0/24","gw":"11.1.1.2"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetRoute(), map[string]interface{}{
		"name": "/Common/external-route",
	})

	diags := dataSourceBigipNetRouteRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/external-route", d.Get("full_path"))
}

func TestDataSourceBigipNetRouteReadWithPartition(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route/~TEST2~test-route", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"test-route","partition":"TEST2","fullPath":"/TEST2/test-route","network":"10.20.20.0/24","gw":"11.2.1.2"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetRoute(), map[string]interface{}{
		"name":      "test-route",
		"partition": "TEST2",
	})

	diags := dataSourceBigipNetRouteRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/TEST2/test-route", d.Get("full_path"))
	assert.Equal(t, "TEST2", d.Get("partition"))
}

func TestDataSourceBigipNetRouteReadTmInterfaceAndBlackhole(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route/~Common~blackhole-route", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"blackhole-route","partition":"Common","fullPath":"/Common/blackhole-route","network":"10.30.30.0/24","tmInterface":"/Common/tunnel1","blackhole":true,"mtu":1400}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetRoute(), map[string]interface{}{
		"name": "blackhole-route",
	})

	diags := dataSourceBigipNetRouteRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/tunnel1", d.Get("tm_interface"))
	assert.Equal(t, true, d.Get("blackhole"))
	assert.Equal(t, 1400, d.Get("mtu"))
}

func TestDataSourceBigipNetRouteReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route/~Common~missing", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"route not found: missing"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetRoute(), map[string]interface{}{
		"name": "missing",
	})

	diags := dataSourceBigipNetRouteRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the route does not exist")
	assert.Contains(t, diags[0].Summary, "route not found")
}

func TestDataSourceBigipNetRouteReadLookupError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route/~Common~external-route", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetRoute(), map[string]interface{}{
		"name": "external-route",
	})

	diags := dataSourceBigipNetRouteRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the route lookup call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving route")
}
