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

func TestResourceBigipLtmProfileHttpSchema(t *testing.T) {
	r := resourceBigipLtmProfileHttp()

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

func TestUnitLtmProfileHttpCreateReadUpdateDelete(t *testing.T) {
	name := "test-http-profile"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/http","proxyType":"reverse","acceptXff":"disabled","basicAuthRealm":"none","description":"my http profile","encryptCookieSecret":"","encryptCookies":[],"fallbackHost":"","fallbackStatusCodes":[],"headErase":"none","headInsert":"none","insertXforwardedFor":"disabled","lwsSeparator":"","oneconnectTransformations":"enabled","tmPartition":"Common","redirectRewrite":"none","requestChunking":"preserve","responseChunking":"selective","responseHeadersPermitted":[],"serverAgentName":"BigIP","viaHostName":"","viaRequest":"preserve","viaResponse":"preserve","xffAlternativeNames":[],"enforcement":{"maxHeaderCount":64,"maxHeaderSize":32768,"unknownMethod":"allow","knownMethods":["CONNECT","DELETE"]},"hsts":{"includeSubdomains":"enabled","maximumAge":16000000,"mode":"enabled","preload":"disabled"}}`, name)
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
	client.Teem = true // skip real telemetry network call
	r := resourceBigipLtmProfileHttp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":             name,
		"defaults_from":    "/Common/http",
		"description":      "my http profile",
		"basic_auth_realm": "none",
		"enforcement": []interface{}{
			map[string]interface{}{
				"known_methods":    []interface{}{"CONNECT", "DELETE"},
				"max_header_count": 64,
				"max_header_size":  32768,
				"unknown_method":   "allow",
			},
		},
		"http_strict_transport_security": []interface{}{
			map[string]interface{}{
				"include_subdomains": "enabled",
				"maximum_age":        16000000,
				"mode":               "enabled",
				"preload":            "disabled",
			},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileHttpCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileHttpRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "my http profile", d.Get("description").(string))

	diags = resourceBigipLtmProfileHttpUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileHttpDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileHttpCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipLtmProfileHttp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-http-profile",
	}, "")

	diags := resourceBigipLtmProfileHttpCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileHttpRead_Error(t *testing.T) {
	name := "test-http-profile"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileHttpRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileHttpUpdate_Error(t *testing.T) {
	name := "test-http-profile"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileHttpUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileHttpDelete_Error(t *testing.T) {
	name := "test-http-profile"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/http/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileHttp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileHttpDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getHttpProfileConfig
// ---------------------------------------------------------------------

func TestGetHttpProfileConfig(t *testing.T) {
	r := resourceBigipLtmProfileHttp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                       "test-http-profile",
		"app_service":                "my-app",
		"defaults_from":              "/Common/http",
		"accept_xff":                 "enabled",
		"basic_auth_realm":           "myrealm",
		"description":                "my desc",
		"encrypt_cookie_secret":      "secret",
		"encrypt_cookies":            []interface{}{"cookie1"},
		"fallback_host":              "fallback.example.com",
		"fallback_status_codes":      []interface{}{"500", "503"},
		"head_erase":                 "X-Erase",
		"head_insert":                "X-Insert: 1",
		"insert_xforwarded_for":      "enabled",
		"lws_separator":              "\\t",
		"lws_width":                  120,
		"oneconnect_transformations": "disabled",
		"tm_partition":               "Common",
		"proxy_type":                 "reverse",
		"redirect_rewrite":           "all",
		"request_chunking":           "unchunk",
		"response_chunking":          "rechunk",
		"response_headers_permitted": []interface{}{"X-Foo"},
		"server_agent_name":          "MyServer",
		"via_host_name":              "myhost",
		"via_request":                "append",
		"via_response":               "remove",
		"xff_alternative_names":      []interface{}{"X-Real-IP"},
		"enforcement": []interface{}{
			map[string]interface{}{
				"known_methods":    []interface{}{"GET", "POST"},
				"max_header_count": 32,
				"max_header_size":  16384,
				"unknown_method":   "reject",
			},
		},
		"http_strict_transport_security": []interface{}{
			map[string]interface{}{
				"include_subdomains": "disabled",
				"maximum_age":        3600,
				"mode":               "disabled",
				"preload":            "enabled",
			},
		},
	}, "")

	cfg := getHttpProfileConfig(d, &bigip.HttpProfile{})
	require.Equal(t, "my-app", cfg.AppService)
	require.Equal(t, "/Common/http", cfg.DefaultsFrom)
	require.Equal(t, "enabled", cfg.AcceptXff)
	require.Equal(t, "myrealm", cfg.BasicAuthRealm)
	require.Equal(t, "my desc", cfg.Description)
	require.Equal(t, "secret", cfg.EncryptCookieSecret)
	require.Equal(t, []string{"cookie1"}, cfg.EncryptCookies)
	require.Equal(t, "fallback.example.com", cfg.FallbackHost)
	require.ElementsMatch(t, []string{"500", "503"}, cfg.FallbackStatusCodes)
	require.Equal(t, "X-Erase", cfg.HeaderErase)
	require.Equal(t, "X-Insert: 1", cfg.HeaderInsert)
	require.Equal(t, "enabled", cfg.InsertXforwardedFor)
	require.Equal(t, "\\t", cfg.LwsSeparator)
	require.Equal(t, 120, cfg.LwsWidth)
	require.Equal(t, "disabled", cfg.OneconnectTransformations)
	require.Equal(t, "Common", cfg.TmPartition)
	require.Equal(t, "reverse", cfg.ProxyType)
	require.Equal(t, "all", cfg.RedirectRewrite)
	require.Equal(t, "unchunk", cfg.RequestChunking)
	require.Equal(t, "rechunk", cfg.ResponseChunking)
	require.Equal(t, "MyServer", cfg.ServerAgentName)
	require.Equal(t, "myhost", cfg.ViaHostName)
	require.Equal(t, "append", cfg.ViaRequest)
	require.Equal(t, "remove", cfg.ViaResponse)
	require.Equal(t, "disabled", cfg.Hsts.IncludeSubdomains)
	require.Equal(t, 3600, cfg.Hsts.MaximumAge)
	require.Equal(t, "disabled", cfg.Hsts.Mode)
	require.Equal(t, "enabled", cfg.Hsts.Preload)
	require.ElementsMatch(t, []string{"GET", "POST"}, cfg.Enforcement.KnownMethods)
	require.Equal(t, "reject", cfg.Enforcement.UnknownMethod)
	require.Equal(t, 32, cfg.Enforcement.MaxHeaderCount)
	require.Equal(t, 16384, cfg.Enforcement.MaxHeaderSize)
}

func TestGetHttpProfileConfig_NoFallbackHost(t *testing.T) {
	// Exercises the "else" branch: fallback_host unset -> FallbackHost
	// forced to empty string.
	r := resourceBigipLtmProfileHttp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-http-profile",
	}, "")

	cfg := getHttpProfileConfig(d, &bigip.HttpProfile{})
	require.Equal(t, "", cfg.FallbackHost)
}
