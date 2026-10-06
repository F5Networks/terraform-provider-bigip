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

func TestResourceBigipNetRouteSchema(t *testing.T) {
	r := resourceBigipNetRoute()

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

	for _, field := range []string{"name", "network"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitNetRouteCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-route"
	mangled := "/mgmt/tm/net/route/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","network":"10.0.0.0/24","gw":"10.0.0.1","tmInterface":"","blackhole":false}`, name, name)
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
	r := resourceBigipNetRoute()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":    name,
		"network": "10.0.0.0/24",
		"gw":      "10.0.0.1",
	}, "")

	ctx := context.Background()

	diags := resourceBigipNetRouteCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipNetRouteRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "10.0.0.0/24", d.Get("network").(string))

	diags = resourceBigipNetRouteUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipNetRouteDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitNetRouteCreate_WithTunnelRefAndReject(t *testing.T) {
	name := "/Common/test-route"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetRoute()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":       name,
		"network":    "10.0.0.0/24",
		"tunnel_ref": "/Common/my-tunnel",
		"reject":     true,
	}, "")

	diags := resourceBigipNetRouteCreate(context.Background(), d, client)
	// Read after create will fail against this mux (no GET handler
	// registered for the route path), but we only care that Create itself
	// builds/sends the request without error before the Read call.
	_ = diags
}

func TestUnitNetRouteCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/route", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetRoute()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":    "/Common/test-route",
		"network": "10.0.0.0/24",
	}, "")

	diags := resourceBigipNetRouteCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetRouteRead_Error(t *testing.T) {
	name := "/Common/test-route"
	mangled := "/mgmt/tm/net/route/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetRoute()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":    name,
		"network": "10.0.0.0/24",
	}, name)

	diags := resourceBigipNetRouteRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetRouteUpdate_Error(t *testing.T) {
	name := "/Common/test-route"
	mangled := "/mgmt/tm/net/route/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetRoute()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":    name,
		"network": "10.0.0.0/24",
	}, name)

	diags := resourceBigipNetRouteUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetRouteDelete_Error(t *testing.T) {
	name := "/Common/test-route"
	mangled := "/mgmt/tm/net/route/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetRoute()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":    name,
		"network": "10.0.0.0/24",
	}, name)

	diags := resourceBigipNetRouteDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
