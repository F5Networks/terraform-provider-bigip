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

func TestResourceBigipGtmTopologyRecordSchema(t *testing.T) {
	r := resourceBigipGtmTopologyRecord()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.CreateContext == nil || r.ReadContext == nil || r.UpdateContext == nil || r.DeleteContext == nil {
		t.Fatal("Expected all CRUD contexts to be defined")
	}
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	for _, field := range []string{"ldns", "server"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "Expected field '%s' to exist in schema", field)
		assert.True(t, s.Required)
		assert.True(t, s.ForceNew)
		assert.Equal(t, 1, s.MaxItems)
	}
}

func TestBuildTopologyNameAndParseRoundTrip(t *testing.T) {
	cases := []struct {
		name         string
		ldnsType     string
		ldnsValue    string
		ldnsNegate   bool
		serverType   string
		serverValue  string
		serverNegate bool
	}{
		{"simple", "region", "/Common/east-coast", false, "datacenter", "/Common/dc1", false},
		{"negated ldns", "country", "US", true, "pool", "/Common/my-pool", false},
		{"negated both", "state", "US/California", true, "continent", "NA", true},
		{"subnet with slash", "subnet", "10.0.0.0/8", false, "isp", "Comcast", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := resourceBigipGtmTopologyRecord()
			d := NewTestResourceData(t, r, map[string]interface{}{
				"ldns": []interface{}{
					map[string]interface{}{
						"match_type":   tc.ldnsType,
						"match_value":  tc.ldnsValue,
						"match_negate": tc.ldnsNegate,
					},
				},
				"server": []interface{}{
					map[string]interface{}{
						"match_type":   tc.serverType,
						"match_value":  tc.serverValue,
						"match_negate": tc.serverNegate,
					},
				},
			}, "")

			built := buildTopologyName(d)

			ldnsType, ldnsValue, ldnsNegate, serverType, serverValue, serverNegate, err := parseTopologyName(built)
			require.NoError(t, err)
			assert.Equal(t, tc.ldnsType, ldnsType)
			assert.Equal(t, tc.ldnsValue, ldnsValue)
			assert.Equal(t, tc.ldnsNegate, ldnsNegate)
			assert.Equal(t, tc.serverType, serverType)
			assert.Equal(t, tc.serverValue, serverValue)
			assert.Equal(t, tc.serverNegate, serverNegate)
		})
	}
}

func TestParseTopologyNameInvalid(t *testing.T) {
	_, _, _, _, _, _, err := parseTopologyName("this is not a valid topology name")
	require.Error(t, err)
}

func TestParseEndpointStringInvalid(t *testing.T) {
	_, _, _, err := parseEndpointString("nospacehere")
	require.Error(t, err)
}

func TestUnitGtmTopologyRecordCreateReadUpdateDelete(t *testing.T) {
	topologyName := "ldns: region /Common/east-coast server: datacenter /Common/dc1"

	// The topology name (BIG-IP's identifier for this resource) contains
	// spaces and a colon, e.g. "ldns: region ... server: ...". Go 1.22+'s
	// http.ServeMux treats a colon in a registered pattern as a "METHOD
	// pattern" separator, so registering the exact literal path panics.
	// Register a prefix pattern instead and dispatch on the full path here.
	mux := http.NewServeMux()
	var updated, sawPost, sawPut, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/topology", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":%q,"order":0,"score":1}`, topologyName)
	})
	mux.HandleFunc("/mgmt/tm/gtm/topology/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			sawPut = true
			updated = true
			fallthrough
		case http.MethodGet:
			if updated {
				_, _ = fmt.Fprintf(w, `{"name":%q,"description":"updated desc","order":5,"score":10}`, topologyName)
			} else {
				_, _ = fmt.Fprintf(w, `{"name":%q,"order":0,"score":1}`, topologyName)
			}
		case http.MethodDelete:
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRecord()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"ldns": []interface{}{
			map[string]interface{}{"match_type": "region", "match_value": "/Common/east-coast", "match_negate": false},
		},
		"server": []interface{}{
			map[string]interface{}{"match_type": "datacenter", "match_value": "/Common/dc1", "match_negate": false},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipGtmTopologyRecordCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost, "expected a POST request during create")
	assert.Equal(t, topologyName, d.Id())

	require.NoError(t, d.Set("description", "updated desc"))
	require.NoError(t, d.Set("order", 5))
	require.NoError(t, d.Set("score", 10))

	diags = resourceBigipGtmTopologyRecordUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPut, "expected a PUT request during update")
	assert.Equal(t, "updated desc", d.Get("description"))
	assert.Equal(t, 5, d.Get("order"))
	assert.Equal(t, 10, d.Get("score"))

	diags = resourceBigipGtmTopologyRecordRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	diags = resourceBigipGtmTopologyRecordDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete, "expected a DELETE request during delete")
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmTopologyRecordCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/topology", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprint(w, `{"code":409,"message":"the requested object already exists"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRecord()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"ldns": []interface{}{
			map[string]interface{}{"match_type": "region", "match_value": "/Common/dup", "match_negate": false},
		},
		"server": []interface{}{
			map[string]interface{}{"match_type": "datacenter", "match_value": "/Common/dc1", "match_negate": false},
		},
	}, "")

	diags := resourceBigipGtmTopologyRecordCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmTopologyRecordReadNotFound(t *testing.T) {
	topologyName := "ldns: region /Common/missing server: datacenter /Common/dc1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/topology/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRecord()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"ldns": []interface{}{
			map[string]interface{}{"match_type": "region", "match_value": "/Common/missing", "match_negate": false},
		},
		"server": []interface{}{
			map[string]interface{}{"match_type": "datacenter", "match_value": "/Common/dc1", "match_negate": false},
		},
	}, topologyName)

	// A 404 means the topology record no longer exists on the device; Read
	// clears the resource ID rather than returning an error.
	diags := resourceBigipGtmTopologyRecordRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Equal(t, "", d.Id())
}

func TestUnitGtmTopologyRecordReadParseError(t *testing.T) {
	// The mock returns a name that doesn't match the "ldns: ... server: ..."
	// format, so parseTopologyName fails inside Read.
	topologyName := "ldns: region /Common/east-coast server: datacenter /Common/dc1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/topology/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"not a valid topology string","order":0,"score":1}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRecord()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"ldns": []interface{}{
			map[string]interface{}{"match_type": "region", "match_value": "/Common/east-coast", "match_negate": false},
		},
		"server": []interface{}{
			map[string]interface{}{"match_type": "datacenter", "match_value": "/Common/dc1", "match_negate": false},
		},
	}, topologyName)

	diags := resourceBigipGtmTopologyRecordRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmTopologyRecordUpdateError(t *testing.T) {
	topologyName := "ldns: region /Common/err server: datacenter /Common/dc1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/topology/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRecord()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"ldns": []interface{}{
			map[string]interface{}{"match_type": "region", "match_value": "/Common/err", "match_negate": false},
		},
		"server": []interface{}{
			map[string]interface{}{"match_type": "datacenter", "match_value": "/Common/dc1", "match_negate": false},
		},
	}, topologyName)

	diags := resourceBigipGtmTopologyRecordUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmTopologyRecordDeleteError(t *testing.T) {
	topologyName := "ldns: region /Common/delerr server: datacenter /Common/dc1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/topology/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRecord()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"ldns": []interface{}{
			map[string]interface{}{"match_type": "region", "match_value": "/Common/delerr", "match_negate": false},
		},
		"server": []interface{}{
			map[string]interface{}{"match_type": "datacenter", "match_value": "/Common/dc1", "match_negate": false},
		},
	}, topologyName)

	diags := resourceBigipGtmTopologyRecordDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmTopologyRecordImport(t *testing.T) {
	r := resourceBigipGtmTopologyRecord()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "ldns: region /Common/east-coast server: datacenter /Common/dc1")

	results, err := resourceBigipGtmTopologyRecordImport(context.Background(), d, nil)
	require.NoError(t, err)
	require.Len(t, results, 1)
}

func TestUnitGtmTopologyRecordImportInvalid(t *testing.T) {
	r := resourceBigipGtmTopologyRecord()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "not-a-valid-topology-name")

	_, err := resourceBigipGtmTopologyRecordImport(context.Background(), d, nil)
	require.Error(t, err)
}
