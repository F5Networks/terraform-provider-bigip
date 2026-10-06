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

func TestDataSourceBigipNetRoutesSchema(t *testing.T) {
	r := dataSourceBigipNetRoutes()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	routes, ok := r.Schema["routes"]
	require.True(t, ok, "expected 'routes' field to exist in schema")
	assert.True(t, routes.Computed)
}

func TestDataSourceBigipNetRoutesReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[
			{"name":"default-route","partition":"Common","fullPath":"/Common/default-route","network":"default","gw":"10.10.10.254"},
			{"name":"blackhole-route","partition":"Common","fullPath":"/Common/blackhole-route","network":"192.0.2.0/24","blackhole":true}
		]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetRoutes(), map[string]interface{}{})

	diags := dataSourceBigipNetRoutesRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	routes := d.Get("routes").([]interface{})
	require.Len(t, routes, 2)

	first := routes[0].(map[string]interface{})
	assert.Equal(t, "default-route", first["name"])
	assert.Equal(t, "default", first["network"])
	assert.Equal(t, "10.10.10.254", first["gw"])
	assert.Equal(t, false, first["blackhole"])

	second := routes[1].(map[string]interface{})
	assert.Equal(t, "192.0.2.0/24", second["network"])
	assert.Equal(t, true, second["blackhole"])
	assert.NotEmpty(t, d.Id())
}

func TestDataSourceBigipNetRoutesReadEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetRoutes(), map[string]interface{}{})

	diags := dataSourceBigipNetRoutesRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Get("routes").([]interface{}))
}

func TestDataSourceBigipNetRoutesReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetRoutes(), map[string]interface{}{})

	diags := dataSourceBigipNetRoutesRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the route list call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving routes")
}
