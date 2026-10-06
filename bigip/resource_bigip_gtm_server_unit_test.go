/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceBigipGtmServerSchema(t *testing.T) {
	r := resourceBigipGtmServer()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.CreateContext == nil || r.ReadContext == nil || r.UpdateContext == nil || r.DeleteContext == nil {
		t.Fatal("Expected all CRUD contexts to be defined")
	}
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	nameSchema, ok := r.Schema["name"]
	require.True(t, ok)
	assert.True(t, nameSchema.Required)
	assert.True(t, nameSchema.ForceNew)

	dcSchema, ok := r.Schema["datacenter"]
	require.True(t, ok)
	assert.True(t, dcSchema.Required)

	partitionSchema, ok := r.Schema["partition"]
	require.True(t, ok)
	assert.Equal(t, "Common", partitionSchema.Default)
}

// Note: unlike resource_bigip_gtm_datacenter.go and
// resource_bigip_gtm_topology_region.go, resource_bigip_gtm_server.go's
// Read/Update/Delete all operate on the bare "name" (d.Get("name")), not the
// mangled full path (d.Id()) -- so mock handlers below are registered at
// /mgmt/tm/gtm/server/<name> without a partition segment.

func TestUnitGtmServerCreateReadUpdateDelete(t *testing.T) {
	fullPath := "/Common/test-server"

	mux := http.NewServeMux()
	var updated bool
	var sawPost, sawPut, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/server", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		sawPost = true
		_, _ = fmt.Fprint(w, `{"name":"test-server","datacenter":"/Common/dc1"}`)
	})
	mux.HandleFunc("/mgmt/tm/gtm/server/test-server", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			sawPut = true
			updated = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		case http.MethodGet:
			if updated {
				_, _ = fmt.Fprint(w, `{"name":"test-server","datacenter":"/Common/dc1","description":"updated","product":"bigip","enabled":true,"monitor":"/Common/bigip","virtualServerDiscovery":"enabled","linkDiscovery":"disabled","proberPreference":"inherit","proberFallback":"inherit","exposeRouteDomains":"yes","iqAllowPath":"yes","iqAllowServiceCheck":"yes","iqAllowSnmp":"yes"}`)
			} else {
				_, _ = fmt.Fprint(w, `{"name":"test-server","datacenter":"/Common/dc1"}`)
			}
		case http.MethodDelete:
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})
	mux.HandleFunc("/mgmt/tm/gtm/server/test-server/virtual-servers", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       "test-server",
		"partition":  "Common",
		"datacenter": "/Common/dc1",
	}, "")

	ctx := context.Background()

	diags := resourceBigipGtmServerCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.Equal(t, fullPath, d.Id())
	assert.True(t, sawPut, "Create calls Update internally")

	diags = resourceBigipGtmServerRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "test-server", d.Get("name"))
	assert.Equal(t, "/Common/dc1", d.Get("datacenter"))
	assert.Equal(t, "updated", d.Get("description"))
	assert.Equal(t, true, d.Get("iq_allow_path"))

	require.NoError(t, d.Set("description", "second update"))
	diags = resourceBigipGtmServerUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipGtmServerDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmServerCreateWithAddressesAndMonitor(t *testing.T) {
	mux := http.NewServeMux()
	var receivedCreateBody string
	mux.HandleFunc("/mgmt/tm/gtm/server", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedCreateBody = string(body)
		_, _ = fmt.Fprint(w, `{"name":"srv-addr","datacenter":"/Common/dc1"}`)
	})
	mux.HandleFunc("/mgmt/tm/gtm/server/srv-addr", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _ = fmt.Fprint(w, `{"name":"srv-addr","datacenter":"/Common/dc1"}`)
	})
	mux.HandleFunc("/mgmt/tm/gtm/server/srv-addr/virtual-servers", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       "srv-addr",
		"partition":  "Common",
		"datacenter": "/Common/dc1",
		"monitor":    "/Common/bigip",
		"addresses": []interface{}{
			map[string]interface{}{"name": "10.1.1.1", "device_name": "dev1", "translation": "none"},
		},
	}, "")

	diags := resourceBigipGtmServerCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Contains(t, receivedCreateBody, "10.1.1.1")
	assert.Contains(t, receivedCreateBody, "/Common/bigip")
}

func TestUnitGtmServerCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/server", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprint(w, `{"code":409,"message":"the requested object already exists"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       "dup-server",
		"partition":  "Common",
		"datacenter": "/Common/dc1",
	}, "")

	diags := resourceBigipGtmServerCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmServerReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/server/broken-server", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       "broken-server",
		"partition":  "Common",
		"datacenter": "/Common/dc1",
	}, "/Common/broken-server")

	diags := resourceBigipGtmServerRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmServerReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/server/missing-server", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       "missing-server",
		"partition":  "Common",
		"datacenter": "/Common/dc1",
	}, "/Common/missing-server")

	// A 404 means the server no longer exists on the device; Read clears
	// the resource ID rather than returning an error.
	diags := resourceBigipGtmServerRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Equal(t, "", d.Id())
}

func TestUnitGtmServerReadFullPathFallback(t *testing.T) {
	// Exercises the else branch when the ID doesn't split into exactly
	// partition+name (a bare name with no leading slash).
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/server/bare-server", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"bare-server","datacenter":"/Common/dc1"}`)
	})
	mux.HandleFunc("/mgmt/tm/gtm/server/bare-server/virtual-servers", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       "bare-server",
		"datacenter": "/Common/dc1",
	}, "bare-server")

	diags := resourceBigipGtmServerRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "bare-server", d.Get("name"))
}

func TestUnitGtmServerUpdateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/server/err-server", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       "err-server",
		"partition":  "Common",
		"datacenter": "/Common/dc1",
	}, "/Common/err-server")

	diags := resourceBigipGtmServerUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmServerUpdateVirtualServersCreateModifyDelete(t *testing.T) {
	mux := http.NewServeMux()
	var sawCreateVS, sawModifyVS, sawDeleteVS bool

	mux.HandleFunc("/mgmt/tm/gtm/server/vs-server", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _ = fmt.Fprint(w, `{"name":"vs-server","datacenter":"/Common/dc1"}`)
	})
	mux.HandleFunc("/mgmt/tm/gtm/server/vs-server/virtual-servers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// existing-vs will be modified (it's also configured below);
			// stale-vs will be deleted (it's not configured below).
			_, _ = fmt.Fprint(w, `{"items":[{"name":"existing-vs","destination":"10.1.1.2:80"},{"name":"stale-vs","destination":"10.1.1.3:80"}]}`)
		case http.MethodPost:
			sawCreateVS = true
			_, _ = fmt.Fprint(w, `{"name":"new-vs"}`)
		}
	})
	mux.HandleFunc("/mgmt/tm/gtm/server/vs-server/virtual-servers/existing-vs", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		sawModifyVS = true
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/mgmt/tm/gtm/server/vs-server/virtual-servers/stale-vs", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		sawDeleteVS = true
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       "vs-server",
		"partition":  "Common",
		"datacenter": "/Common/dc1",
		"virtual_servers": []interface{}{
			map[string]interface{}{
				"name":        "new-vs",
				"destination": "10.1.1.1:80",
				"enabled":     true,
			},
			map[string]interface{}{
				"name":        "existing-vs",
				"destination": "10.1.1.2:80",
				"enabled":     true,
			},
		},
	}, "/Common/vs-server")
	// Force HasChange("virtual_servers") to be true: schema.TestResourceDataRaw
	// builds ResourceData from a diff against a nil old state, so any
	// non-zero-value field already reports HasChange == true here.

	diags := resourceBigipGtmServerUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawCreateVS, "expected new-vs to be created")
	assert.True(t, sawModifyVS, "expected existing-vs to be modified")
	assert.True(t, sawDeleteVS, "expected stale-vs to be deleted")
}

func TestUnitGtmServerDeleteError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/server/del-err-server", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "del-err-server",
	}, "/Common/del-err-server")

	diags := resourceBigipGtmServerDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmServerDeleteNotFoundIgnored(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/server/already-gone-server", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"the requested GTM Server (/Common/already-gone-server) was not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "already-gone-server",
	}, "/Common/already-gone-server")

	diags := resourceBigipGtmServerDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected 'was not found' delete error to be ignored: %v", diags)
	assert.Equal(t, "", d.Id())
}
