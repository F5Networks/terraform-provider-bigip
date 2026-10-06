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

// TestDataSourceBigipLtmNodeReadSuccess covers the Read happy path where
// GetNode succeeds and returns a non-FQDN node, so address/name/partition
// and the other computed fields get populated from the API response.
func TestDataSourceBigipLtmNodeReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~test-node", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"test-node","partition":"Common","fullPath":"/Common/test-node","address":"10.10.10.10","connectionLimit":10,"dynamicRatio":1,"monitor":"/Common/icmp","rateLimit":"disabled","ratio":1,"state":"user-up","session":"user-enabled"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipLtmNode(), map[string]interface{}{
		"name":      "test-node",
		"partition": "Common",
	})

	diags := dataSourceBigipLtmNodeRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "test-node", d.Id())
	assert.Equal(t, "test-node", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))
	assert.Equal(t, "10.10.10.10", d.Get("address"))
	assert.Equal(t, 10, d.Get("connection_limit"))
	assert.Equal(t, 1, d.Get("dynamic_ratio"))
	assert.Equal(t, "/Common/icmp", d.Get("monitor"))
}

// TestDataSourceBigipLtmNodeReadFQDN covers Read when the node is an FQDN
// node so address is populated from node.FQDN.Name instead of node.Address.
func TestDataSourceBigipLtmNodeReadFQDN(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~fqdn-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"fqdn-node","partition":"Common","fullPath":"/Common/fqdn-node","fqdn":{"tmName":"example.com","interval":"3600","downInterval":5,"autopopulate":"enabled","addressFamily":"ipv4"}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipLtmNode(), map[string]interface{}{
		"name":      "fqdn-node",
		"partition": "Common",
	})

	diags := dataSourceBigipLtmNodeRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "fqdn-node", d.Id())
	assert.Equal(t, "example.com", d.Get("address"))
}

// TestDataSourceBigipLtmNodeReadNotFound covers the branch where GetNode
// returns a 404. A 404 means the node does not exist on the device;
// dataSourceBigipLtmNodeRead treats that as "not found" (no error, empty ID)
// rather than a hard failure.
func TestDataSourceBigipLtmNodeReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~missing-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipLtmNode(), map[string]interface{}{
		"name":      "missing-node",
		"partition": "Common",
	})
	d.SetId("preexisting-id")

	diags := dataSourceBigipLtmNodeRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "expected no error: a 404 clears state instead of failing")
	require.Empty(t, d.Id())
}

// TestDataSourceBigipLtmNodeReadError covers the branch where GetNode
// itself returns an error (e.g. non-404 server error).
func TestDataSourceBigipLtmNodeReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~broken-node", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipLtmNode(), map[string]interface{}{
		"name":      "broken-node",
		"partition": "Common",
	})
	d.SetId("preexisting-id")

	diags := dataSourceBigipLtmNodeRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when the API call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving node")
}
