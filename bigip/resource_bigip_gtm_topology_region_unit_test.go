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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceBigipGtmTopologyRegionSchema(t *testing.T) {
	r := resourceBigipGtmTopologyRegion()

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
	require.True(t, ok, "Expected field 'name' to exist in schema")
	assert.True(t, nameSchema.Required)
	assert.True(t, nameSchema.ForceNew)

	membersSchema, ok := r.Schema["members"]
	require.True(t, ok, "Expected field 'members' to exist in schema")
	assert.Equal(t, schema.TypeSet, membersSchema.Type)
}

func TestUnitGtmTopologyRegionCreateReadUpdateDelete(t *testing.T) {
	fullPath := "/Common/test-region"
	mangled := "/mgmt/tm/gtm/region/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	var updated, sawPost, sawPut, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/region", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":"test-region","partition":"Common","fullPath":"%s"}`, fullPath)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			sawPut = true
			updated = true
			fallthrough
		case http.MethodGet:
			if updated {
				_, _ = fmt.Fprintf(w, `{"name":"test-region","partition":"Common","fullPath":"%s","regionMembers":[{"name":"subnet 10.0.0.0/8"},{"name":"country US"}]}`, fullPath)
			} else {
				_, _ = fmt.Fprintf(w, `{"name":"test-region","partition":"Common","fullPath":"%s"}`, fullPath)
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

	r := resourceBigipGtmTopologyRegion()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "test-region",
		"partition": "Common",
	}, "")

	ctx := context.Background()

	diags := resourceBigipGtmTopologyRegionCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.Equal(t, fullPath, d.Id())
	assert.Equal(t, "test-region", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))

	diags = resourceBigipGtmTopologyRegionUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPut)
	// Update's internal Read now sees the post-update response with members.
	members := d.Get("members").(*schema.Set)
	assert.Equal(t, 2, members.Len())

	diags = resourceBigipGtmTopologyRegionRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	members = d.Get("members").(*schema.Set)
	assert.Equal(t, 2, members.Len())

	diags = resourceBigipGtmTopologyRegionDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmTopologyRegionCreateWithMembers(t *testing.T) {
	fullPath := "/Common/region-with-members"
	mangled := "/mgmt/tm/gtm/region/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	var receivedBody string
	mux.HandleFunc("/mgmt/tm/gtm/region", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		_, _ = fmt.Fprintf(w, `{"name":"region-with-members","partition":"Common","fullPath":"%s"}`, fullPath)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"region-with-members","partition":"Common","fullPath":"%s","regionMembers":[{"name":"datacenter /Common/dc1"}]}`, fullPath)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRegion()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "region-with-members",
		"partition": "Common",
		"members": []interface{}{
			map[string]interface{}{"name": "datacenter /Common/dc1"},
		},
	}, "")

	diags := resourceBigipGtmTopologyRegionCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Contains(t, receivedBody, "datacenter /Common/dc1")
}

func TestUnitGtmTopologyRegionCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/region", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprint(w, `{"code":409,"message":"the requested object already exists"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRegion()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "dup-region",
		"partition": "Common",
	}, "")

	diags := resourceBigipGtmTopologyRegionCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmTopologyRegionReadNotFound(t *testing.T) {
	fullPath := "/Common/missing-region"
	mangled := "/mgmt/tm/gtm/region/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRegion()
	d := NewTestResourceData(t, r, map[string]interface{}{"name": "missing-region", "partition": "Common"}, fullPath)

	// A 404 means the topology region no longer exists on the device; Read
	// clears the resource ID rather than returning an error.
	diags := resourceBigipGtmTopologyRegionRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Equal(t, "", d.Id())
}

func TestUnitGtmTopologyRegionUpdateError(t *testing.T) {
	fullPath := "/Common/region-err"
	mangled := "/mgmt/tm/gtm/region/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRegion()
	d := NewTestResourceData(t, r, map[string]interface{}{"name": "region-err", "partition": "Common"}, fullPath)

	diags := resourceBigipGtmTopologyRegionUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmTopologyRegionDeleteError(t *testing.T) {
	fullPath := "/Common/region-del-err"
	mangled := "/mgmt/tm/gtm/region/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRegion()
	d := NewTestResourceData(t, r, map[string]interface{}{"name": "region-del-err", "partition": "Common"}, fullPath)

	diags := resourceBigipGtmTopologyRegionDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmTopologyRegionReadFullPathFallback(t *testing.T) {
	name := "bare-region"
	mangled := "/mgmt/tm/gtm/region/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"bare-region","partition":"Common"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmTopologyRegion()
	d := NewTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmTopologyRegionRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "bare-region", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))
}
