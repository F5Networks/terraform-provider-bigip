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

func TestDataSourceBigipNetVlansSchema(t *testing.T) {
	r := dataSourceBigipNetVlans()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	vlans, ok := r.Schema["vlans"]
	require.True(t, ok, "expected 'vlans' field to exist in schema")
	assert.True(t, vlans.Computed)
}

func TestDataSourceBigipNetVlansReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[
			{"name":"external","partition":"Common","fullPath":"/Common/external","tag":100,"mtu":1500,"cmpHash":"default"},
			{"name":"internal","partition":"Common","fullPath":"/Common/internal","tag":200,"mtu":1500,"cmpHash":"src-ip"}
		]}`)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~external/interfaces", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"1.1","tagged":false}]}`)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~internal/interfaces", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"1.2","tagged":true}]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlans(), map[string]interface{}{})

	diags := dataSourceBigipNetVlansRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	vlans := d.Get("vlans").([]interface{})
	require.Len(t, vlans, 2)

	first := vlans[0].(map[string]interface{})
	assert.Equal(t, "external", first["name"])
	assert.Equal(t, "/Common/external", first["full_path"])
	assert.Equal(t, 100, first["tag"])
	ifaces := first["interfaces"].([]interface{})
	require.Len(t, ifaces, 1)
	iface := ifaces[0].(map[string]interface{})
	assert.Equal(t, "1.1", iface["name"])
	assert.Equal(t, false, iface["tagged"])

	second := vlans[1].(map[string]interface{})
	assert.Equal(t, "internal", second["name"])
	assert.Equal(t, "src-ip", second["cmp_hash"])
	assert.NotEmpty(t, d.Id())
}

func TestDataSourceBigipNetVlansReadEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlans(), map[string]interface{}{})

	diags := dataSourceBigipNetVlansRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Get("vlans").([]interface{}))
}

func TestDataSourceBigipNetVlansReadListError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlans(), map[string]interface{}{})

	diags := dataSourceBigipNetVlansRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the VLAN list call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving VLANs")
}

func TestDataSourceBigipNetVlansReadInterfacesError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"external","partition":"Common","fullPath":"/Common/external","tag":100}]}`)
	})
	mux.HandleFunc("/mgmt/tm/net/vlan/~Common~external/interfaces", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetVlans(), map[string]interface{}{})

	diags := dataSourceBigipNetVlansRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when an interfaces lookup fails")
	assert.Contains(t, diags[0].Summary, "error retrieving interfaces for VLAN")
}
