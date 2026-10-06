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

func ltmMonitorDatasourceResourceData(t *testing.T, name, partition string) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipLtmMonitor(), map[string]interface{}{
		"name":      name,
		"partition": partition,
	})
}

func TestDataSourceBigipLtmMonitorReadSuccess(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"http": func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"items":[{"name":"test-http-monitor","fullPath":"/Common/test-http-monitor","defaultsFrom":"/Common/http","interval":5,"timeout":16}]}`)
		},
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := ltmMonitorDatasourceResourceData(t, "test-http-monitor", "Common")

	diags := dataSourceBigipLtmMonitorRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/test-http-monitor", d.Id())
	assert.Equal(t, "/Common/http", d.Get("defaults_from"))
	assert.Equal(t, 5, d.Get("interval"))
	assert.Equal(t, 16, d.Get("timeout"))
}

func TestDataSourceBigipLtmMonitorReadNotFound(t *testing.T) {
	mux := newMonitorTestMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := ltmMonitorDatasourceResourceData(t, "missing-monitor", "Common")

	diags := dataSourceBigipLtmMonitorRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Id())
}

func TestDataSourceBigipLtmMonitorReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/monitor/http", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := ltmMonitorDatasourceResourceData(t, "test-http-monitor", "Common")

	diags := dataSourceBigipLtmMonitorRead(context.Background(), d, client)

	require.True(t, diags.HasError())
}

func TestDataSourceBigipLtmMonitorReadSubPartition(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"ldap": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"items":[{"name":"sub-monitor","fullPath":"/MyPartition/sub-monitor","defaultsFrom":"/Common/ldap","base":"dc=example,dc=com","filter":"objectClass=*"}]}`)
		},
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := ltmMonitorDatasourceResourceData(t, "sub-monitor", "MyPartition")

	diags := dataSourceBigipLtmMonitorRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/MyPartition/sub-monitor", d.Id())
	assert.Equal(t, "dc=example,dc=com", d.Get("base"))
	assert.Equal(t, "objectClass=*", d.Get("filter"))
}
