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

func TestDataSourceBigipSysDnsSchema(t *testing.T) {
	r := dataSourceBigipSysDns()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	for _, field := range []string{"description", "name_servers", "search", "number_of_dots"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "expected field '%s' to exist in schema", field)
		assert.True(t, s.Computed, "expected field '%s' to be computed", field)
	}
}

func TestDataSourceBigipSysDnsReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"description":"configured-by-dhcp","nameServers":["172.27.1.1","8.8.8.8"],"numberOfDots":2,"search":["example.com","internal"]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysDns(), map[string]interface{}{})

	diags := dataSourceBigipSysDnsRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "configured-by-dhcp", d.Get("description"))
	assert.Equal(t, 2, d.Get("number_of_dots"))

	nameServers := d.Get("name_servers").([]interface{})
	require.Len(t, nameServers, 2)
	assert.Equal(t, "172.27.1.1", nameServers[0])
	assert.Equal(t, "8.8.8.8", nameServers[1])

	search := d.Get("search").([]interface{})
	require.Len(t, search, 2)
	assert.Equal(t, "example.com", search[0])
	assert.Equal(t, "internal", search[1])

	assert.NotEmpty(t, d.Id())
}

// TestDataSourceBigipSysDnsReadEmpty covers a response with no description
// (BIG-IP omits the field entirely when unset, rather than returning an
// empty string). Read must not set an empty ID in this case -- it falls
// back to hashing the struct instead of using dns.Description directly.
func TestDataSourceBigipSysDnsReadEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysDns(), map[string]interface{}{})

	diags := dataSourceBigipSysDnsRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Get("name_servers").([]interface{}))
	assert.Empty(t, d.Get("search").([]interface{}))
	assert.NotEmpty(t, d.Id(), "ID must not be empty even when description is unset")
}

func TestDataSourceBigipSysDnsReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysDns(), map[string]interface{}{})

	diags := dataSourceBigipSysDnsRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the DNS lookup call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving system DNS configuration")
}
