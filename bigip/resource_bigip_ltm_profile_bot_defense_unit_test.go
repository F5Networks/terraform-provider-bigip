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
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipLtmProfileBotDefenseSchema(t *testing.T) {
	r := resourceBigipLtmProfileBotDefense()

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
	if r.Schema["defaults_from"].Default != "/Common/bot-defense" {
		t.Errorf("Expected 'defaults_from' to default to '/Common/bot-defense', got %v", r.Schema["defaults_from"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitLtmProfileBotDefenseCreateReadUpdateDelete(t *testing.T) {
	name := "test-bot-defense"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/security/bot-defense/profile", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/security/bot-defense/profile/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/bot-defense","description":"my bot defense","template":"balanced","enforcementMode":"transparent"}`, name, name)
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
	r := resourceBigipLtmProfileBotDefense()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/bot-defense",
		"template":      "balanced",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileBotDefenseCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileBotDefenseRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "balanced", d.Get("template").(string))

	diags = resourceBigipLtmProfileBotDefenseUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileBotDefenseDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileBotDefenseCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/security/bot-defense/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileBotDefense()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-bot-defense",
	}, "")

	diags := resourceBigipLtmProfileBotDefenseCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileBotDefenseRead_Error(t *testing.T) {
	name := "test-bot-defense"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/security/bot-defense/profile/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileBotDefense()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileBotDefenseRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileBotDefenseUpdate_Error(t *testing.T) {
	name := "test-bot-defense"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/security/bot-defense/profile/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileBotDefense()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileBotDefenseUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileBotDefenseDelete_Error(t *testing.T) {
	name := "test-bot-defense"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/security/bot-defense/profile/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileBotDefense()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileBotDefenseDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getProfileBotDefenseConfig
// ---------------------------------------------------------------------

func TestGetProfileBotDefenseConfig(t *testing.T) {
	r := resourceBigipLtmProfileBotDefense()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":             "test-bot-defense",
		"defaults_from":    "/Common/bot-defense",
		"description":      "my bot defense",
		"template":         "strict",
		"enforcement_mode": "blocking",
	}, "")

	cfg := getProfileBotDefenseConfig(d, &bigip.BotDefenseProfile{})
	require.Equal(t, "test-bot-defense", cfg.Name)
	require.Equal(t, "/Common/bot-defense", cfg.DefaultsFrom)
	require.Equal(t, "my bot defense", cfg.Description)
	require.Equal(t, "strict", cfg.Template)
	require.Equal(t, "blocking", cfg.EnforcementMode)
}
