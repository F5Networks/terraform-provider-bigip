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

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipSysDnsSchema(t *testing.T) {
	r := resourceBigipSysDns()

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

	for _, field := range []string{"description", "name_servers"} {
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

func TestUnitSysDnsCreateReadUpdateDelete(t *testing.T) {
	description := "test dns"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			_, _ = fmt.Fprint(w, `{"description":"test dns"}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, `{"description":"test dns","nameServers":["8.8.8.8"],"numberOfDots":2,"search":["example.com"]}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysDns()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description":  description,
		"name_servers": []interface{}{"8.8.8.8"},
		"search":       []interface{}{"example.com"},
	}, "")

	ctx := context.Background()

	diags := resourceBigipSysDnsCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, description, d.Id())

	diags = resourceBigipSysDnsRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, 2, d.Get("number_of_dots").(int))

	diags = resourceBigipSysDnsUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSysDnsDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitSysDnsCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysDns()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description":  "test dns",
		"name_servers": []interface{}{"8.8.8.8"},
	}, "")

	diags := resourceBigipSysDnsCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysDnsRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysDns()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description":  "test dns",
		"name_servers": []interface{}{"8.8.8.8"},
	}, "test dns")

	diags := resourceBigipSysDnsRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysDnsUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysDns()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description":  "test dns",
		"name_servers": []interface{}{"8.8.8.8"},
	}, "test dns")

	diags := resourceBigipSysDnsUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysDnsDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/dns", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysDns()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description":  "test dns",
		"name_servers": []interface{}{"8.8.8.8"},
	}, "test dns")

	diags := resourceBigipSysDnsDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getSysDNSConfig
// ---------------------------------------------------------------------

func TestGetSysDNSConfig(t *testing.T) {
	r := resourceBigipSysDns()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description":    "test dns",
		"name_servers":   []interface{}{"8.8.8.8", "8.8.4.4"},
		"number_of_dots": 3,
		"search":         []interface{}{"example.com"},
	}, "")

	cfg := getSysDNSConfig(d, &bigip.DNS{})
	require.Equal(t, []string{"8.8.8.8", "8.8.4.4"}, cfg.NameServers)
	require.Equal(t, 3, cfg.NumberOfDots)
	require.Equal(t, []string{"example.com"}, cfg.Search)
}
