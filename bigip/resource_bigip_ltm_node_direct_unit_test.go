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

// nodeResourceData builds a *schema.ResourceData for the bigip_ltm_node
// resource using schema.TestResourceDataRaw.
func nodeResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceBigipLtmNode().Schema, raw)
}

// TestResourceBigipLtmNodeDirectCreate covers the Create success path for a
// plain (non-FQDN) node: the Exists check (GetNode) returns nil (not yet
// created), so AddNode is invoked, followed by a Read.
func TestResourceBigipLtmNodeDirectCreate(t *testing.T) {
	mux := http.NewServeMux()
	created := false
	mux.HandleFunc("/mgmt/tm/ltm/node", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		created = true
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"/Common/test-node","address":"10.10.10.10"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~test-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !created && r.Method == "GET" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
			return
		}
		_, _ = fmt.Fprint(w, `{"name":"/Common/test-node","address":"10.10.10.10","session":"user-enabled","monitor":"/Common/icmp"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/test-node",
		"address": "10.10.10.10",
		"monitor": "/Common/icmp",
	})

	diags := resourceBigipLtmNodeCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, created, "expected AddNode to be called")
	assert.Equal(t, "/Common/test-node", d.Id())
}

// TestResourceBigipLtmNodeDirectCreateFQDN covers the Create path for an
// FQDN node (address does not match the IP/IPv6 regex), exercising the
// FQDN branch of Create.
func TestResourceBigipLtmNodeDirectCreateFQDN(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"fqdn-node"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~fqdn-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			// alternate: first GET (exists check) returns not found,
			// second GET (post-create read) returns the node.
		}
		_, _ = fmt.Fprint(w, `{"name":"/Common/fqdn-node","fqdn":{"tmName":"example.com","interval":"3600"}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/fqdn-node",
		"address": "example.com",
		"fqdn": []interface{}{
			map[string]interface{}{
				"interval": "3600",
			},
		},
	})

	diags := resourceBigipLtmNodeCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/fqdn-node", d.Id())
}

// TestResourceBigipLtmNodeDirectCreateAlreadyExists covers the Create path
// where the Exists check finds the node already present, so AddNode must
// not be called (it goes straight to Read).
func TestResourceBigipLtmNodeDirectCreateAlreadyExists(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node", func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("AddNode should not be called when the node already exists")
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~existing-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"/Common/existing-node","address":"10.10.10.11","session":"user-enabled"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/existing-node",
		"address": "10.10.10.11",
	})

	diags := resourceBigipLtmNodeCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/existing-node", d.Id())
}

// TestResourceBigipLtmNodeDirectCreateAddNodeError covers the Create path
// where AddNode returns an error: the ID must be cleared and an error
// diagnostic returned.
func TestResourceBigipLtmNodeDirectCreateAddNodeError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~bad-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/bad-node",
		"address": "10.10.10.12",
	})

	diags := resourceBigipLtmNodeCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when AddNode fails")
	assert.Empty(t, d.Id())
}

// TestResourceBigipLtmNodeDirectReadSuccess covers the Read success path,
// including the fqdn.Set branch when the "fqdn" key is present in state.
func TestResourceBigipLtmNodeDirectReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~read-node", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"/Common/read-node","address":"10.10.10.13","session":"user-enabled","monitor":"/Common/icmp ","connectionLimit":5,"dynamicRatio":2,"ratio":3,"description":"desc"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/read-node",
		"address": "10.10.10.13",
		"fqdn": []interface{}{
			map[string]interface{}{},
		},
	})
	d.SetId("/Common/read-node")

	diags := resourceBigipLtmNodeRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "10.10.10.13", d.Get("address"))
	assert.Equal(t, "user-enabled", d.Get("session"))
	assert.Equal(t, "/Common/icmp", d.Get("monitor"))
	assert.Equal(t, 5, d.Get("connection_limit"))
	assert.Equal(t, "desc", d.Get("description"))
}

// TestResourceBigipLtmNodeDirectReadSessionDisabled covers the branch where
// node.Session is neither monitor-enabled nor user-enabled, so session is
// set to user-disabled.
func TestResourceBigipLtmNodeDirectReadSessionDisabled(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~disabled-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"/Common/disabled-node","address":"10.10.10.14","session":"user-disabled"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/disabled-node",
		"address": "10.10.10.14",
	})
	d.SetId("/Common/disabled-node")

	diags := resourceBigipLtmNodeRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "user-disabled", d.Get("session"))
}

// TestResourceBigipLtmNodeDirectReadFQDNAddress covers the Read branch
// where node.FQDN.Name is set, so address is populated from it.
func TestResourceBigipLtmNodeDirectReadFQDNAddress(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~fqdn-read-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"/Common/fqdn-read-node","fqdn":{"tmName":"example.com"},"session":"user-enabled"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/fqdn-read-node",
		"address": "example.com",
	})
	d.SetId("/Common/fqdn-read-node")

	diags := resourceBigipLtmNodeRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "example.com", d.Get("address"))
}

// TestResourceBigipLtmNodeDirectReadNotFound covers the Read branch where
// GetNode returns a 404. A 404 means the node no longer exists on the
// device; Read clears the resource ID rather than returning an error.
func TestResourceBigipLtmNodeDirectReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~missing-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/missing-node",
		"address": "10.10.10.15",
	})
	d.SetId("/Common/missing-node")

	diags := resourceBigipLtmNodeRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "expected no error: a 404 clears state instead of failing")
	require.Empty(t, d.Id())
}

// TestResourceBigipLtmNodeDirectReadError covers the Read branch where
// GetNode returns an error.
func TestResourceBigipLtmNodeDirectReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~broken-node", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/broken-node",
		"address": "10.10.10.16",
	})
	d.SetId("/Common/broken-node")

	diags := resourceBigipLtmNodeRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when GetNode fails")
}

// TestResourceBigipLtmNodeDirectExists covers the Exists helper's true,
// false, and error branches.
func TestResourceBigipLtmNodeDirectExists(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~exists-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"/Common/exists-node","address":"10.10.10.17"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~missing-exists-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~error-exists-node", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	dExists := nodeResourceData(t, map[string]interface{}{"name": "/Common/exists-node", "address": "10.10.10.17"})
	dExists.SetId("/Common/exists-node")
	ok, err := resourceBigipLtmNodeExists(dExists, client)
	require.NoError(t, err)
	assert.True(t, ok)

	dMissing := nodeResourceData(t, map[string]interface{}{"name": "/Common/missing-exists-node", "address": "10.10.10.18"})
	dMissing.SetId("/Common/missing-exists-node")
	ok, err = resourceBigipLtmNodeExists(dMissing, client)
	require.NoError(t, err, "a 404 means the node doesn't exist, not an error")
	assert.False(t, ok)

	dError := nodeResourceData(t, map[string]interface{}{"name": "/Common/error-exists-node", "address": "10.10.10.19"})
	dError.SetId("/Common/error-exists-node")
	ok, err = resourceBigipLtmNodeExists(dError, client)
	require.Error(t, err)
	assert.False(t, ok)
}

// TestResourceBigipLtmNodeDirectUpdate covers the Update success path for
// a plain node address (regex matches, so Address is set on the update
// payload) followed by a successful Read.
func TestResourceBigipLtmNodeDirectUpdate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~update-node", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "PUT" || r.Method == "PATCH" {
			_, _ = fmt.Fprint(w, `{"name":"/Common/update-node","address":"10.10.10.20"}`)
			return
		}
		_, _ = fmt.Fprint(w, `{"name":"/Common/update-node","address":"10.10.10.20","session":"user-enabled"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/update-node",
		"address": "10.10.10.20",
		"monitor": "/Common/icmp",
	})
	d.SetId("/Common/update-node")

	diags := resourceBigipLtmNodeUpdate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

// TestResourceBigipLtmNodeDirectUpdateError covers the Update path where
// ModifyNode returns an error.
func TestResourceBigipLtmNodeDirectUpdateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~update-error-node", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/update-error-node",
		"address": "10.10.10.21",
	})
	d.SetId("/Common/update-error-node")

	diags := resourceBigipLtmNodeUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when ModifyNode fails")
}

// TestResourceBigipLtmNodeDirectDelete covers the Delete success path.
func TestResourceBigipLtmNodeDirectDelete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~delete-node", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/delete-node",
		"address": "10.10.10.22",
	})
	d.SetId("/Common/delete-node")

	diags := resourceBigipLtmNodeDelete(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Id())
}

// TestResourceBigipLtmNodeDirectDeleteError covers the Delete path where
// DeleteNode returns an error.
func TestResourceBigipLtmNodeDirectDeleteError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~delete-error-node", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := nodeResourceData(t, map[string]interface{}{
		"name":    "/Common/delete-error-node",
		"address": "10.10.10.23",
	})
	d.SetId("/Common/delete-error-node")

	diags := resourceBigipLtmNodeDelete(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when DeleteNode fails")
}
