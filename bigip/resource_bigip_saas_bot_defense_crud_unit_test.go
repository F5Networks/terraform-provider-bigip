/*
Copyright 2024 F5 Networks Inc.
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipSaasBotDefenseProfileSchema(t *testing.T) {
	r := resourceBigipSaasBotDefenseProfile()

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

	for _, field := range []string{"name", "application_id", "tenant_id", "api_key", "shape_protection_pool", "ssl_profile", "protected_endpoints"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}
	if r.Schema["defaults_from"].Default != "/Common/bd" {
		t.Errorf("Expected 'defaults_from' to default to '/Common/bd', got %v", r.Schema["defaults_from"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func testBotDefenseProfileResourceData(t *testing.T, name, id string) *schema.ResourceData {
	r := resourceBigipSaasBotDefenseProfile()
	return NewTestResourceData(t, r, map[string]interface{}{
		"name":                  name,
		"application_id":        "app-id",
		"tenant_id":             "tenant-id",
		"api_key":               "api-key",
		"shape_protection_pool": "/Common/cs1.pool",
		"ssl_profile":           "/Common/cloud-service-default-ssl",
		"protected_endpoints": []interface{}{
			map[string]interface{}{
				"name":              "pe1",
				"host":              "abc.com",
				"endpoint":          "/login",
				"post":              "enabled",
				"put":               "disabled",
				"mitigation_action": "",
			},
		},
	}, id)
}

func TestUnitSaasBotDefenseProfileCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/bd-test"
	mangled := "/mgmt/tm/saas/bd/profile/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/saas/bd/profile", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"bd-test","fullPath":"%s","defaultsFrom":"/Common/bd","tenantId":"tenant-id","apiKey":"api-key","shapeProtectionPool":"/Common/cs1.pool","sslProfile":"/Common/cloud-service-default-ssl","protectedEndpointsReference":{"items":[{"name":"pe1","host":"abc.com","endpoint":"/login","post":"enabled"}]}}`, name)
	}
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})
	mux.HandleFunc(mangled+"/", getHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSaasBotDefenseProfile()
	d := testBotDefenseProfileResourceData(t, name, "")

	ctx := context.Background()

	diags := resourceBigipSaasBotDefenseProfileCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipSaasBotDefenseProfileRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "tenant-id", d.Get("tenant_id").(string))

	diags = resourceBigipSaasBotDefenseProfileUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSaasBotDefenseProfileDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
	_ = r
}

func TestUnitSaasBotDefenseProfileCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/saas/bd/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := testBotDefenseProfileResourceData(t, "/Common/bd-test", "")

	diags := resourceBigipSaasBotDefenseProfileCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSaasBotDefenseProfileRead_Error(t *testing.T) {
	name := "/Common/bd-test"
	mangled := "/mgmt/tm/saas/bd/profile/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled+"/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := testBotDefenseProfileResourceData(t, name, name)

	diags := resourceBigipSaasBotDefenseProfileRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSaasBotDefenseProfileUpdate_Error(t *testing.T) {
	name := "/Common/bd-test"
	mangled := "/mgmt/tm/saas/bd/profile/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := testBotDefenseProfileResourceData(t, name, name)

	diags := resourceBigipSaasBotDefenseProfileUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSaasBotDefenseProfileDelete_Error(t *testing.T) {
	name := "/Common/bd-test"
	mangled := "/mgmt/tm/saas/bd/profile/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := testBotDefenseProfileResourceData(t, name, name)

	diags := resourceBigipSaasBotDefenseProfileDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getSaasBotDefenseProfileConfig / flattenProtectedEndpointsReference
// ---------------------------------------------------------------------

func TestGetSaasBotDefenseProfileConfig(t *testing.T) {
	d := testBotDefenseProfileResourceData(t, "/Common/bd-test", "")

	cfg := getSaasBotDefenseProfileConfig(d, &bigip.SaasBotDefenseProfile{})
	require.Equal(t, "/Common/bd-test", cfg.Name)
	require.Equal(t, "tenant-id", cfg.TenantId)
	require.Equal(t, "api-key", cfg.ApiKey)
	require.Equal(t, "/Common/cs1.pool", cfg.ShapeProtectionPool)
	require.Equal(t, "/Common/cloud-service-default-ssl", cfg.SslProfile)
	require.Len(t, cfg.ProtectedEndpointsReference.Items, 1)
	require.Equal(t, "pe1", cfg.ProtectedEndpointsReference.Items[0].Name)
	require.Equal(t, "abc.com", cfg.ProtectedEndpointsReference.Items[0].Host)
}

func TestFlattenProtectedEndpointsReference(t *testing.T) {
	items := []bigip.ProtectedEndpoint{
		{Name: "pe1", Host: "abc.com", Endpoint: "/login", Post: "enabled", Put: "disabled", MitigationAction: "block"},
	}

	flattened := flattenProtectedEndpointsReference(items)
	require.Len(t, flattened, 1)
	m := flattened[0].(map[string]interface{})
	require.Equal(t, "pe1", m["name"])
	require.Equal(t, "abc.com", m["host"])
	require.Equal(t, "block", m["mitigation_action"])
}

func TestFlattenProtectedEndpointsReference_Empty(t *testing.T) {
	var items []bigip.ProtectedEndpoint
	flattened := flattenProtectedEndpointsReference(items)
	require.Empty(t, flattened)
}
