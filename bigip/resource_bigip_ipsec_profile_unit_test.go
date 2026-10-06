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

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipIpsecProfileSchema(t *testing.T) {
	r := resourceBigipIpsecProfile()

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
	if r.Schema["parent_profile"].Default != "/Common/ipsec" {
		t.Errorf("Expected 'parent_profile' to default to '/Common/ipsec', got %v", r.Schema["parent_profile"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitIpsecProfileCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-ipsec-profile"
	mangled := "/mgmt/tm/net/tunnels/ipsec/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/tunnels/ipsec", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/ipsec","trafficSelector":"/Common/ts1","description":"desc"}`, name)
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
	r := resourceBigipIpsecProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":             name,
		"parent_profile":   "/Common/ipsec",
		"traffic_selector": "/Common/ts1",
		"description":      "desc",
	}, "")

	ctx := context.Background()

	diags := resourceBigipIpsecProfileCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipIpsecProfileRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "/Common/ipsec", d.Get("parent_profile").(string))

	diags = resourceBigipIpsecProfileUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipIpsecProfileDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitIpsecProfileCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/tunnels/ipsec", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipIpsecProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/test-ipsec-profile",
	}, "")

	diags := resourceBigipIpsecProfileCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitIpsecProfileRead_Error(t *testing.T) {
	name := "/Common/test-ipsec-profile"
	mangled := "/mgmt/tm/net/tunnels/ipsec/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipIpsecProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipIpsecProfileRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitIpsecProfileUpdate_Error(t *testing.T) {
	name := "/Common/test-ipsec-profile"
	mangled := "/mgmt/tm/net/tunnels/ipsec/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipIpsecProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipIpsecProfileUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitIpsecProfileDelete_Error(t *testing.T) {
	name := "/Common/test-ipsec-profile"
	mangled := "/mgmt/tm/net/tunnels/ipsec/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipIpsecProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipIpsecProfileDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getIPSecProfileConfig
// ---------------------------------------------------------------------

func TestGetIPSecProfileConfig(t *testing.T) {
	r := resourceBigipIpsecProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":             "/Common/test-ipsec-profile",
		"parent_profile":   "/Common/custom-parent",
		"traffic_selector": "/Common/ts1",
		"description":      "my desc",
	}, "")

	cfg := getIPSecProfileConfig(d, &bigip.IPSecProfile{})
	require.Equal(t, "/Common/custom-parent", cfg.DefaultsFrom)
	require.Equal(t, "/Common/ts1", cfg.TrafficSelector)
	require.Equal(t, "my desc", cfg.Description)
}
