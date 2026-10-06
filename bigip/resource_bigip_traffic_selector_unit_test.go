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

func TestResourceBigipTrafficselectorSchema(t *testing.T) {
	r := resourceBigipTrafficselector()

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

	for _, field := range []string{"name", "destination_address", "source_address"} {
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

func TestUnitTrafficselectorCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-traffic-selector"
	mangled := "/mgmt/tm/net/ipsec/traffic-selector/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","destinationAddress":"10.0.0.0/24","sourceAddress":"10.0.1.0/24","ipsecPolicy":"/Common/default-ipsec-policy","order":1,"destinationPort":0,"sourcePort":0,"direction":"both","ipProtocol":255}`, name)
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
	r := resourceBigipTrafficselector()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                name,
		"destination_address": "10.0.0.0/24",
		"source_address":      "10.0.1.0/24",
		"ipsec_policy":        "/Common/default-ipsec-policy",
	}, "")

	ctx := context.Background()

	diags := resourceBigipTrafficselectorCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipTrafficselectorRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "10.0.0.0/24", d.Get("destination_address").(string))

	diags = resourceBigipTrafficselectorUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipTrafficselectorDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitTrafficselectorCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/ipsec/traffic-selector", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipTrafficselector()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                "/Common/test-traffic-selector",
		"destination_address": "10.0.0.0/24",
		"source_address":      "10.0.1.0/24",
	}, "")

	diags := resourceBigipTrafficselectorCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitTrafficselectorRead_Error(t *testing.T) {
	name := "/Common/test-traffic-selector"
	mangled := "/mgmt/tm/net/ipsec/traffic-selector/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipTrafficselector()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                name,
		"destination_address": "10.0.0.0/24",
		"source_address":      "10.0.1.0/24",
	}, name)

	diags := resourceBigipTrafficselectorRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitTrafficselectorUpdate_Error(t *testing.T) {
	name := "/Common/test-traffic-selector"
	mangled := "/mgmt/tm/net/ipsec/traffic-selector/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipTrafficselector()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                name,
		"destination_address": "10.0.0.0/24",
		"source_address":      "10.0.1.0/24",
	}, name)

	diags := resourceBigipTrafficselectorUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitTrafficselectorDelete_Error(t *testing.T) {
	name := "/Common/test-traffic-selector"
	mangled := "/mgmt/tm/net/ipsec/traffic-selector/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipTrafficselector()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                name,
		"destination_address": "10.0.0.0/24",
		"source_address":      "10.0.1.0/24",
	}, name)

	diags := resourceBigipTrafficselectorDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getTrafficSelectorConfig
// ---------------------------------------------------------------------

func TestGetTrafficSelectorConfig(t *testing.T) {
	r := resourceBigipTrafficselector()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                "/Common/test-traffic-selector",
		"destination_address": "10.0.0.0/24",
		"destination_port":    443,
		"source_address":      "10.0.1.0/24",
		"source_port":         80,
		"direction":           "inbound",
		"ipsec_policy":        "/Common/my-policy",
		"order":               5,
		"ip_protocol":         6,
	}, "")

	cfg := getTrafficSelectorConfig(d, &bigip.TrafficSelector{})
	require.Equal(t, "10.0.0.0/24", cfg.DestinationAddress)
	require.Equal(t, 443, cfg.DestinationPort)
	require.Equal(t, "10.0.1.0/24", cfg.SourceAddress)
	require.Equal(t, 80, cfg.SourcePort)
	require.Equal(t, "inbound", cfg.Direction)
	require.Equal(t, "/Common/my-policy", cfg.IpsecPolicy)
	require.Equal(t, 5, cfg.Order)
	require.Equal(t, 6, cfg.IPProtocol)
}
