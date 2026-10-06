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

func TestDataSourceBigipNetSelfipsSchema(t *testing.T) {
	r := dataSourceBigipNetSelfips()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	selfIPs, ok := r.Schema["self_ips"]
	require.True(t, ok, "expected 'self_ips' field to exist in schema")
	assert.True(t, selfIPs.Computed)
}

func TestDataSourceBigipNetSelfipsReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[
			{"name":"external-self","partition":"Common","fullPath":"/Common/external-self","address":"10.10.10.1/24","vlan":"/Common/external","trafficGroup":"/Common/traffic-group-local-only","floating":"false","unit":1,"allowService":"all"},
			{"name":"internal-self","partition":"Common","fullPath":"/Common/internal-self","address":"10.10.20.1/24","vlan":"/Common/internal","trafficGroup":"/Common/traffic-group-local-only","floating":"false","unit":1,"allowService":null},
			{"name":"restricted-self","partition":"Common","fullPath":"/Common/restricted-self","address":"10.10.30.1/24","vlan":"/Common/restricted","trafficGroup":"/Common/traffic-group-local-only","floating":"false","unit":1,"allowService":["tcp:443","tcp:22"]},
			{"name":"emptystring-self","partition":"Common","fullPath":"/Common/emptystring-self","address":"10.10.40.1/24","vlan":"/Common/emptystring","trafficGroup":"/Common/traffic-group-local-only","floating":"false","unit":1,"allowService":""}
		]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetSelfips(), map[string]interface{}{})

	diags := dataSourceBigipNetSelfipsRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	selfIPs := d.Get("self_ips").([]interface{})
	require.Len(t, selfIPs, 4)

	first := selfIPs[0].(map[string]interface{})
	assert.Equal(t, "external-self", first["name"])
	assert.Equal(t, "10.10.10.1/24", first["address"])
	assert.Equal(t, "/Common/external", first["vlan"])
	allowService := first["allow_service"].([]interface{})
	require.Len(t, allowService, 1)
	assert.Equal(t, "all", allowService[0])

	second := selfIPs[1].(map[string]interface{})
	allowServiceNone := second["allow_service"].([]interface{})
	require.Len(t, allowServiceNone, 1)
	assert.Equal(t, "none", allowServiceNone[0])

	third := selfIPs[2].(map[string]interface{})
	allowServiceSlice := third["allow_service"].([]interface{})
	require.Len(t, allowServiceSlice, 2)
	assert.Equal(t, "tcp:443", allowServiceSlice[0])
	assert.Equal(t, "tcp:22", allowServiceSlice[1])

	fourth := selfIPs[3].(map[string]interface{})
	allowServiceEmptyString := fourth["allow_service"].([]interface{})
	require.Len(t, allowServiceEmptyString, 1, "an empty-string allowService should normalize to [\"none\"], not [\"\"]")
	assert.Equal(t, "none", allowServiceEmptyString[0])

	assert.NotEmpty(t, d.Id())
}

func TestDataSourceBigipNetSelfipsReadEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetSelfips(), map[string]interface{}{})

	diags := dataSourceBigipNetSelfipsRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Get("self_ips").([]interface{}))
}

func TestDataSourceBigipNetSelfipsReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipNetSelfips(), map[string]interface{}{})

	diags := dataSourceBigipNetSelfipsRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the self IP list call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving self IPs")
}
