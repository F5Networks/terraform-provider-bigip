/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
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

func TestResourceBigipNetVlanSchema(t *testing.T) {
	r := resourceBigipNetVlan()

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
	if r.Schema["mtu"].Default != 1500 {
		t.Errorf("Expected 'mtu' to default to 1500, got %v", r.Schema["mtu"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitNetVlanCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-vlan"
	mangled := "/mgmt/tm/net/vlan/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled+"/interfaces", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = fmt.Fprint(w, `{"name":"1.1"}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, `{"items":[{"name":"1.1","tagged":true}]}`)
		}
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","tag":100,"cmpHash":"default","mtu":1500}`, name, name)
		case http.MethodPut:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetVlan()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"tag":  100,
		"mtu":  1500,
		"interfaces": []interface{}{
			map[string]interface{}{
				"vlanport": "1.1",
				"tagged":   true,
			},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipNetVlanCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipNetVlanRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, 100, d.Get("tag").(int))

	diags = resourceBigipNetVlanUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipNetVlanDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitNetVlanCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetVlan()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/test-vlan",
	}, "")

	diags := resourceBigipNetVlanCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetVlanCreate_AddInterfaceError(t *testing.T) {
	name := "/Common/test-vlan"
	mangled := "/mgmt/tm/net/vlan/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/vlan", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled+"/interfaces", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"add interface failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetVlan()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"interfaces": []interface{}{
			map[string]interface{}{
				"vlanport": "1.1",
				"tagged":   true,
			},
		},
	}, "")

	diags := resourceBigipNetVlanCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetVlanRead_Error(t *testing.T) {
	name := "/Common/test-vlan"
	mangled := "/mgmt/tm/net/vlan/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetVlan()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipNetVlanRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetVlanRead_InterfacesError(t *testing.T) {
	name := "/Common/test-vlan"
	mangled := "/mgmt/tm/net/vlan/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","tag":100,"mtu":1500}`, name, name)
	})
	mux.HandleFunc(mangled+"/interfaces", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"interfaces read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetVlan()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipNetVlanRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetVlanUpdate_Error(t *testing.T) {
	name := "/Common/test-vlan"
	mangled := "/mgmt/tm/net/vlan/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetVlan()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipNetVlanUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetVlanDelete_Error(t *testing.T) {
	name := "/Common/test-vlan"
	mangled := "/mgmt/tm/net/vlan/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetVlan()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipNetVlanDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
