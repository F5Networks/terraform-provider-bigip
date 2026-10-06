/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2024 F5 Networks Inc.
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

func TestResourceBigipLtmRewriteProfileSchema(t *testing.T) {
	r := resourceBigipLtmRewriteProfile()

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
	if !r.Schema["rewrite_mode"].Required {
		t.Error("Expected field 'rewrite_mode' to be required")
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitLtmRewriteProfileCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-rewrite"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+MangleFullPath(name), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"test-rewrite","fullPath":"%s","rewriteMode":"portal","defaultsFrom":"/Common/rewrite","javaCaFile":"none","javaCrl":"none","javaSigner":"none","javaSignKey":"none","clientCachingType":"cache-img-css-js","splitTunneling":"false","rewriteList":["/foo"],"bypassList":["/bar"],"request":{"insertXforwardedFor":"enabled","insertXforwardedHost":"enabled","insertXforwardedProto":"enabled","rewriteHeaders":"enabled"},"response":{"rewriteContent":"enabled","rewriteHeaders":"enabled"},"setCookieRules":[{"name":"rule1","client":{"domain":"client.example.com","path":"/"},"server":{"domain":"server.example.com","path":"/"}}]}`, name)
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
	r := resourceBigipLtmRewriteProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"rewrite_mode":  "portal",
		"defaults_from": "/Common/rewrite",
		"ca_file":       "none",
		"crl_file":      "none",
		"signing_cert":  "none",
		"signing_key":   "none",
		"cache_type":    "cache-img-css-js",
		"rewrite_list":  []interface{}{"/foo"},
		"bypass_list":   []interface{}{"/bar"},
		"request": []interface{}{
			map[string]interface{}{
				"insert_xfwd_for":      "enabled",
				"insert_xfwd_host":     "enabled",
				"insert_xfwd_protocol": "enabled",
				"rewrite_headers":      "enabled",
			},
		},
		"response": []interface{}{
			map[string]interface{}{
				"rewrite_content": "enabled",
				"rewrite_headers": "enabled",
			},
		},
		"cookie_rules": []interface{}{
			map[string]interface{}{
				"rule_name":     "rule1",
				"client_domain": "client.example.com",
				"client_path":   "/",
				"server_domain": "server.example.com",
				"server_path":   "/",
			},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmRewriteProfileCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileRewriteRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "portal", d.Get("rewrite_mode").(string))

	diags = resourceBigipLtmProfileRewriteUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileRewriteDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmRewriteProfileCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmRewriteProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":         "/Common/test-rewrite",
		"rewrite_mode": "portal",
	}, "")

	diags := resourceBigipLtmRewriteProfileCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileRewriteRead_Error(t *testing.T) {
	name := "/Common/test-rewrite"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+MangleFullPath(name), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmRewriteProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":         name,
		"rewrite_mode": "portal",
	}, name)

	diags := resourceBigipLtmProfileRewriteRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileRewriteUpdate_Error(t *testing.T) {
	name := "/Common/test-rewrite"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+MangleFullPath(name), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmRewriteProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":         name,
		"rewrite_mode": "portal",
	}, name)

	diags := resourceBigipLtmProfileRewriteUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileRewriteDelete_Error(t *testing.T) {
	name := "/Common/test-rewrite"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/rewrite/"+MangleFullPath(name), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmRewriteProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":         name,
		"rewrite_mode": "portal",
	}, name)

	diags := resourceBigipLtmProfileRewriteDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getRewriteProfileConfig / setRewriteProfileData
// ---------------------------------------------------------------------

func TestGetRewriteProfileConfig(t *testing.T) {
	r := resourceBigipLtmRewriteProfile()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "/Common/test-rewrite",
		"rewrite_mode":  "uri-translation",
		"defaults_from": "/Common/rewrite",
		"ca_file":       "/Common/ca.crt",
		"crl_file":      "/Common/crl.pem",
		"signing_cert":  "/Common/signer.crt",
		"signing_key":   "/Common/signer.key",
		"cache_type":    "cache-all",
		"rewrite_list":  []interface{}{"/a", "/b"},
		"bypass_list":   []interface{}{"/c"},
		"request": []interface{}{
			map[string]interface{}{
				"insert_xfwd_for":      "enabled",
				"insert_xfwd_host":     "disabled",
				"insert_xfwd_protocol": "enabled",
				"rewrite_headers":      "disabled",
			},
		},
		"response": []interface{}{
			map[string]interface{}{
				"rewrite_content": "disabled",
				"rewrite_headers": "enabled",
			},
		},
		"cookie_rules": []interface{}{
			map[string]interface{}{
				"rule_name":     "rule1",
				"client_domain": "client.example.com",
				"client_path":   "/",
				"server_domain": "server.example.com",
				"server_path":   "/",
			},
		},
	}, "")

	cfg := getRewriteProfileConfig(d, &bigip.RewriteProfile{})
	require.Equal(t, "uri-translation", cfg.Mode)
	require.Equal(t, "/Common/ca.crt", cfg.CaFile)
	require.Equal(t, "cache-all", cfg.CachingType)
	require.Equal(t, []string{"/a", "/b"}, cfg.RewriteList)
	require.Equal(t, []string{"/c"}, cfg.BypassList)
	require.Equal(t, "enabled", cfg.Request.XfwdFor)
	require.Equal(t, "disabled", cfg.Response.RewriteContent)
	require.Len(t, cfg.Cookies, 1)
	require.Equal(t, "rule1", cfg.Cookies[0].Name)
	require.Equal(t, "client.example.com", cfg.Cookies[0].Client.Domain)
}
