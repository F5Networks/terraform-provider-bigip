/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2024 F5 Networks Inc.
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipLtmRewriteProfileUriRulesSchema(t *testing.T) {
	r := resourceBigipLtmRewriteProfileUriRules()

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

	for _, field := range []string{"profile_name", "rule_name", "client", "server"} {
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

func uriRulesResourceData(t *testing.T, id string) *schema.ResourceData {
	r := resourceBigipLtmRewriteProfileUriRules()
	return NewTestResourceData(t, r, map[string]interface{}{
		"profile_name": "/Common/test-rewrite",
		"rule_name":    "rule1",
		"rule_type":    "both",
		"client": []interface{}{
			map[string]interface{}{
				"host":   "www.foo.com",
				"path":   "/",
				"scheme": "https",
				"port":   "none",
			},
		},
		"server": []interface{}{
			map[string]interface{}{
				"host":   "www.bar.com",
				"path":   "/",
				"scheme": "https",
				"port":   "none",
			},
		},
	}, id)
}

func TestUnitLtmRewriteProfileUriRuleCreateReadDelete(t *testing.T) {
	profileMangled := MangleFullPath("/Common/test-rewrite")
	ruleName := "rule1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+profileMangled+"/uri-rules", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, ruleName)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+profileMangled+"/uri-rules/"+ruleName, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","type":"both","client":{"host":"www.foo.com","path":"/","scheme":"https","port":"none"},"server":{"host":"www.bar.com","path":"/","scheme":"https","port":"none"}}`, ruleName)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := uriRulesResourceData(t, "")

	ctx := context.Background()

	diags := resourceBigipLtmRewriteProfileUriRuleCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, ruleName, d.Id())

	diags = resourceBigipLtmProfileRewriteUriRuleRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "both", d.Get("rule_type").(string))

	diags = resourceBigipLtmProfileRewriteUriRuleDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmRewriteProfileUriRuleUpdate(t *testing.T) {
	profileMangled := MangleFullPath("/Common/test-rewrite")
	ruleName := "rule1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+profileMangled+"/uri-rules/"+ruleName, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, ruleName)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","type":"both","client":{"host":"www.foo.com","path":"/","scheme":"https","port":"none"},"server":{"host":"www.bar.com","path":"/","scheme":"https","port":"none"}}`, ruleName)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := uriRulesResourceData(t, ruleName)

	diags := resourceBigipLtmProfileRewriteUriRuleUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitLtmRewriteProfileUriRuleCreate_Error(t *testing.T) {
	profileMangled := MangleFullPath("/Common/test-rewrite")
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+profileMangled+"/uri-rules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := uriRulesResourceData(t, "")

	diags := resourceBigipLtmRewriteProfileUriRuleCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileRewriteUriRuleRead_Error(t *testing.T) {
	profileMangled := MangleFullPath("/Common/test-rewrite")
	ruleName := "rule1"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+profileMangled+"/uri-rules/"+ruleName, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := uriRulesResourceData(t, ruleName)

	diags := resourceBigipLtmProfileRewriteUriRuleRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileRewriteUriRuleUpdate_Error(t *testing.T) {
	profileMangled := MangleFullPath("/Common/test-rewrite")
	ruleName := "rule1"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+profileMangled+"/uri-rules/"+ruleName, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := uriRulesResourceData(t, ruleName)

	diags := resourceBigipLtmProfileRewriteUriRuleUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileRewriteUriRuleDelete_Error(t *testing.T) {
	profileMangled := MangleFullPath("/Common/test-rewrite")
	ruleName := "rule1"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+profileMangled+"/uri-rules/"+ruleName, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	d := uriRulesResourceData(t, ruleName)

	diags := resourceBigipLtmProfileRewriteUriRuleDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getUriRulesConfig / setUriRulesData
// ---------------------------------------------------------------------

func TestGetUriRulesConfig(t *testing.T) {
	d := uriRulesResourceData(t, "")

	cfg := getUriRulesConfig(d, &bigip.RewriteProfileUriRule{})
	require.Equal(t, "both", cfg.Type)
	require.Equal(t, "www.foo.com", cfg.Client.Host)
	require.Equal(t, "https", cfg.Client.Scheme)
	require.Equal(t, "www.bar.com", cfg.Server.Host)
	require.Equal(t, "https", cfg.Server.Scheme)
}

func TestSetUriRulesData(t *testing.T) {
	d := uriRulesResourceData(t, "rule1")

	data := &bigip.RewriteProfileUriRule{
		Name: "rule1",
		Type: "request",
		Client: bigip.RewriteProfileUrlClSrv{
			Host:   "www.foo.com",
			Path:   "/a",
			Scheme: "http",
			Port:   "8080",
		},
		Server: bigip.RewriteProfileUrlClSrv{
			Host:   "www.bar.com",
			Path:   "/b",
			Scheme: "http",
			Port:   "8081",
		},
	}

	setUriRulesData(d, data)
	require.Equal(t, "rule1", d.Get("rule_name").(string))
	require.Equal(t, "request", d.Get("rule_type").(string))
}
