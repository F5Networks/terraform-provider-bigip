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

func TestDataSourceBigipNetVlanSchema(t *testing.T) {
	r := dataSourceBigipNetVlan()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	name, ok := r.Schema["name"]
	require.True(t, ok, "expected 'name' field to exist in schema")
	assert.True(t, name.Optional)
	assert.ElementsMatch(t, []string{"name", "vlan_id"}, name.ExactlyOneOf)

	vlanID, ok := r.Schema["vlan_id"]
	require.True(t, ok, "expected 'vlan_id' field to exist in schema")
	assert.True(t, vlanID.Optional)
	assert.ElementsMatch(t, []string{"name", "vlan_id"}, vlanID.ExactlyOneOf)

	taggedInterfaces, ok := r.Schema["tagged_interfaces"]
	require.True(t, ok, "expected 'tagged_interfaces' field to exist in schema")
	assert.True(t, taggedInterfaces.Computed)

	untaggedInterfaces, ok := r.Schema["untagged_interfaces"]
	require.True(t, ok, "expected 'untagged_interfaces' field to exist in schema")
	assert.True(t, untaggedInterfaces.Computed)
}

func TestDataSourceBigipNetVlanReadByNameSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~external", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"external","partition":"Common","fullPath":"/Common/external","tag":100,"mtu":1500,"cmpHash":"default"}`)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~external/interfaces", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"1.1","tagged":false},{"name":"1.2","tagged":true}]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlan(), map[string]interface{}{
		"name": "external",
	})

	diags := dataSourceBigipNetVlanRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/external", d.Get("name"))
	assert.Equal(t, "/Common/external", d.Get("full_path"))
	assert.Equal(t, 100, d.Get("vlan_id"))

	tagged := d.Get("tagged_interfaces").([]interface{})
	require.Len(t, tagged, 1)
	assert.Equal(t, "1.2", tagged[0])

	untagged := d.Get("untagged_interfaces").([]interface{})
	require.Len(t, untagged, 1)
	assert.Equal(t, "1.1", untagged[0])

	assert.Equal(t, "/Common/external", d.Id())
}

func TestDataSourceBigipNetVlanReadByFullPathName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~external", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"external","partition":"Common","fullPath":"/Common/external","tag":100}`)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~external/interfaces", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlan(), map[string]interface{}{
		"name": "/Common/external",
	})

	diags := dataSourceBigipNetVlanRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/external", d.Get("full_path"))
}

func TestDataSourceBigipNetVlanReadByVlanIDSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[
			{"name":"external","partition":"Common","fullPath":"/Common/external","tag":100},
			{"name":"internal","partition":"Common","fullPath":"/Common/internal","tag":200}
		]}`)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~internal/interfaces", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"1.3","tagged":true}]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlan(), map[string]interface{}{
		"vlan_id": 200,
	})

	diags := dataSourceBigipNetVlanRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/internal", d.Get("full_path"))
	assert.Equal(t, 200, d.Get("vlan_id"))

	tagged := d.Get("tagged_interfaces").([]interface{})
	require.Len(t, tagged, 1)
	assert.Equal(t, "1.3", tagged[0])
	assert.Equal(t, "/Common/internal", d.Id())
}

func TestDataSourceBigipNetVlanReadByVlanIDNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"external","partition":"Common","fullPath":"/Common/external","tag":100}]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlan(), map[string]interface{}{
		"vlan_id": 999,
	})

	diags := dataSourceBigipNetVlanRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when no VLAN matches the given vlan_id")
	assert.Contains(t, diags[0].Summary, "VLAN not found")
}

// TestDataSourceBigipNetVlanReadByNameNotFound covers a 404 response for a
// name-based lookup. go-bigip's Vlan() returns (nil, nil) for a missing
// object, and the Read function surfaces the same "VLAN not found" message
// the vlan_id path uses, keeping the error consistent across both lookup
// methods.
func TestDataSourceBigipNetVlanReadByNameNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~missing", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlan(), map[string]interface{}{
		"name": "missing",
	})

	diags := dataSourceBigipNetVlanRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the VLAN does not exist")
	assert.Contains(t, diags[0].Summary, "VLAN not found")
}

func TestDataSourceBigipNetVlanReadByNameLookupError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~external", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlan(), map[string]interface{}{
		"name": "external",
	})

	diags := dataSourceBigipNetVlanRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the VLAN lookup call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving VLAN")
}

func TestDataSourceBigipNetVlanReadByVlanIDListError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlan(), map[string]interface{}{
		"vlan_id": 100,
	})

	diags := dataSourceBigipNetVlanRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the VLAN list call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving VLANs")
}

func TestDataSourceBigipNetVlanReadInterfacesError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~external", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"external","partition":"Common","fullPath":"/Common/external","tag":100}`)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~external/interfaces", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlan(), map[string]interface{}{
		"name": "external",
	})

	diags := dataSourceBigipNetVlanRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the interfaces lookup fails")
	assert.Contains(t, diags[0].Summary, "error retrieving interfaces for VLAN")
}
