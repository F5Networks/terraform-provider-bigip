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

func TestDataSourceBigipGtmServerSchema(t *testing.T) {
	r := dataSourceBigipGtmServer()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	nameSchema, ok := r.Schema["name"]
	require.True(t, ok, "Expected field 'name' to exist in schema")
	assert.True(t, nameSchema.Required)

	partitionSchema, ok := r.Schema["partition"]
	require.True(t, ok, "Expected field 'partition' to exist in schema")
	assert.Equal(t, "Common", partitionSchema.Default)

	for _, field := range []string{"datacenter", "description", "product", "enabled", "monitor", "addresses", "virtual_servers"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "Expected field '%s' to exist in schema", field)
		assert.True(t, s.Computed, "expected field '%s' to be computed", field)
	}
}

func TestDataSourceBigipGtmServerReadSuccess(t *testing.T) {
	mangled := "/mgmt/tm/gtm/server/" + MangleFullPath("/Common/test-server")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"name":"test-server","datacenter":"/Common/dc1","description":"a server","product":"bigip","enabled":true,"monitor":"/Common/bigip","virtualServerDiscovery":"enabled","linkDiscovery":"disabled","proberPreference":"inherit","proberFallback":"inherit","proberPool":"","addresses":[{"name":"10.1.1.1","deviceName":"dev1","translation":"none"}]}`)
	})
	mux.HandleFunc(mangled+"/virtual-servers", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[{"name":"vs1","destination":"10.1.1.1:80","enabled":true,"translationAddress":"none","translationPort":0,"monitor":"/Common/http"}]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := dataSourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "test-server",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipGtmServerRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "/Common/test-server", d.Id())
	assert.Equal(t, "test-server", d.Get("name"))
	assert.Equal(t, "/Common/dc1", d.Get("datacenter"))
	assert.Equal(t, "a server", d.Get("description"))
	assert.Equal(t, "bigip", d.Get("product"))
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, "/Common/bigip", d.Get("monitor"))

	addresses := d.Get("addresses").([]interface{})
	require.Len(t, addresses, 1)
	addr := addresses[0].(map[string]interface{})
	assert.Equal(t, "10.1.1.1", addr["name"])

	virtualServers := d.Get("virtual_servers").([]interface{})
	require.Len(t, virtualServers, 1)
	vs := virtualServers[0].(map[string]interface{})
	assert.Equal(t, "vs1", vs["name"])
	assert.Equal(t, "", vs["translation_address"], "expected 'none' translation_address to be normalized to empty string")
}

func TestDataSourceBigipGtmServerReadError(t *testing.T) {
	mangled := "/mgmt/tm/gtm/server/" + MangleFullPath("/Common/broken-server")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := dataSourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "broken-server",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipGtmServerRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestDataSourceBigipGtmServerReadNotFound(t *testing.T) {
	mangled := "/mgmt/tm/gtm/server/" + MangleFullPath("/Common/missing-server")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := dataSourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "missing-server",
		"partition": "Common",
	}, "")

	// GetGtmserver returns (nil, err) on a 404 (getForEntity always yields a
	// non-nil error on 404), so this always takes the error branch, never
	// the "server == nil" branch.
	diags := dataSourceBigipGtmServerRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestDataSourceBigipGtmServerReadNoAddressesOrVirtualServers(t *testing.T) {
	mangled := "/mgmt/tm/gtm/server/" + MangleFullPath("/Common/minimal-server")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"minimal-server","datacenter":"/Common/dc1"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := dataSourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "minimal-server",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipGtmServerRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "minimal-server", d.Get("name"))
}
