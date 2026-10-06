/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
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

func TestResourceBigipLtmProfileOneconnectSchema(t *testing.T) {
	r := resourceBigipLtmProfileOneconnect()

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
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitLtmProfileOneconnectCreateReadUpdateDelete(t *testing.T) {
	name := "test-oneconnect"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","defaultsFrom":"/Common/oneconnect","sharePools":"enabled","sourceMask":"255.255.255.255","limitType":"none","maxAge":3600,"maxReuse":1000,"maxSize":1000,"idleTimeoutOverride":"disabled"}`, name)
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
	r := resourceBigipLtmProfileOneconnect()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                  name,
		"partition":             "Common",
		"defaults_from":         "/Common/oneconnect",
		"share_pools":           "enabled",
		"source_mask":           "255.255.255.255",
		"limit_type":            "none",
		"max_age":               3600,
		"max_reuse":             1000,
		"max_size":              1000,
		"idle_timeout_override": "disabled",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileOneconnectCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileOneconnectRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, 3600, d.Get("max_age").(int))

	diags = resourceBigipLtmProfileOneconnectUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileOneconnectDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileOneconnectRead_NotFound(t *testing.T) {
	name := "test-oneconnect"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileOneconnect()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileOneconnectRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
}

func TestUnitLtmProfileOneconnectCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileOneconnect()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-oneconnect",
	}, "")

	diags := resourceBigipLtmProfileOneconnectCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileOneconnectRead_Error(t *testing.T) {
	name := "test-oneconnect"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileOneconnect()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileOneconnectRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileOneconnectUpdate_Error(t *testing.T) {
	name := "test-oneconnect"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileOneconnect()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileOneconnectUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileOneconnectDelete_Error(t *testing.T) {
	name := "test-oneconnect"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/one-connect/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileOneconnect()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileOneconnectDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
