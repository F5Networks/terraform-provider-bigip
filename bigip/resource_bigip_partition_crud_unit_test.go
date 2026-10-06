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

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipPartitionSchema(t *testing.T) {
	r := resourceBigipPartition()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.CreateContext == nil {
		t.Fatal("Expected CreateContext to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}
	if r.UpdateContext == nil {
		t.Fatal("Expected UpdateContext to be defined")
	}
	if r.DeleteContext == nil {
		t.Fatal("Expected DeleteContext to be defined")
	}

	nameSchema, ok := r.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist in schema")
	}
	if !nameSchema.Required {
		t.Error("Expected field 'name' to be required")
	}
	if !nameSchema.ForceNew {
		t.Error("Expected field 'name' to be ForceNew")
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitPartitionCreateReadUpdateDelete(t *testing.T) {
	name := "testpartition"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/auth/partition", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/sys/folder/~"+name, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPatch)
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/auth/partition/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultRouteDomain":5}`, name)
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipPartition()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":            name,
		"description":     "my partition",
		"route_domain_id": 5,
	}, "")

	ctx := context.Background()

	diags := resourceBigipPartitionCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipPartitionRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, 5, d.Get("route_domain_id").(int))

	diags = resourceBigipPartitionUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipPartitionDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
}

func TestUnitPartitionCreate_NoDescription(t *testing.T) {
	// description == "" means Create skips the ModifyFolderDescription call
	// entirely.
	name := "testpartition"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/auth/partition", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/auth/partition/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultRouteDomain":0}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipPartition()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, "")

	diags := resourceBigipPartitionCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
}

func TestUnitPartitionCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/auth/partition", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipPartition()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "testpartition",
	}, "")

	diags := resourceBigipPartitionCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitPartitionCreate_DescriptionError(t *testing.T) {
	name := "testpartition"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/auth/partition", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/sys/folder/~"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"description failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipPartition()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":        name,
		"description": "my partition",
	}, "")

	diags := resourceBigipPartitionCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitPartitionRead_Error(t *testing.T) {
	name := "testpartition"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/auth/partition/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipPartition()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipPartitionRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitPartitionUpdate_RouteDomainError(t *testing.T) {
	name := "testpartition"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/auth/partition/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipPartition()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":            name,
		"route_domain_id": 5,
	}, name)

	diags := resourceBigipPartitionUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitPartitionUpdate_DescriptionError(t *testing.T) {
	name := "testpartition"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/folder/~"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"description failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipPartition()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":        name,
		"description": "my partition",
	}, name)

	diags := resourceBigipPartitionUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitPartitionUpdate_NoChanges(t *testing.T) {
	// Neither route_domain_id nor description changed from their zero
	// values -- Update should skip both API calls and just Read.
	name := "testpartition"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/auth/partition/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","defaultRouteDomain":0}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipPartition()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipPartitionUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitPartitionDelete_Error(t *testing.T) {
	name := "testpartition"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/auth/partition/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipPartition()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipPartitionDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
