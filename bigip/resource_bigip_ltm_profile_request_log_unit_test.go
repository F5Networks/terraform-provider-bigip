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

func TestResourceBigipLtmProfileRequestLogSchema(t *testing.T) {
	r := resourceBigipLtmProfileRequestLog()

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
	if r.Schema["defaults_from"].Default != "/Common/request-log" {
		t.Errorf("Expected 'defaults_from' to default to '/Common/request-log', got %v", r.Schema["defaults_from"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitLtmProfileRequestLogCreateReadUpdateDelete(t *testing.T) {
	name := "test-request-log"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/request-log", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/request-log/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","defaultsFrom":"/Common/request-log","description":"my request log","requestLogging":"enabled","requestLogPool":"/Common/logpool","requestLogErrorPool":"/Common/errpool","requestLogTemplate":"template","requestLogProtocol":"mds-udp","requestLogErrorProtocol":"mds-udp","responseLogProtocol":"mds-tcp","responseLogErrorProtocol":"mds-tcp","responseLogPool":"/Common/resppool","responseLogErrorPool":"/Common/resperrpool","proxyResponse":"enabled","proxyCloseOnError":"enabled","proxyRespondOnLoggingError":"enabled","responseLogging":"enabled","responseLogTemplate":"resptemplate","requestLogErrorTemplate":"errtemplate","responseLogErrorTemplate":"resperrtemplate"}`, name)
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
	r := resourceBigipLtmProfileRequestLog()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                         name,
		"defaults_from":                "/Common/request-log",
		"description":                  "my request log",
		"request_logging":              "enabled",
		"requestlog_pool":              "/Common/logpool",
		"requestlog_error_pool":        "/Common/errpool",
		"requestlog_template":          "template",
		"requestlog_protocol":          "mds-udp",
		"requestlog_error_protocol":    "mds-udp",
		"responselog_protocol":         "mds-tcp",
		"responselog_error_protocol":   "mds-tcp",
		"responselog_pool":             "/Common/resppool",
		"responselog_error_pool":       "/Common/resperrpool",
		"proxy_response":               "enabled",
		"proxyclose_on_error":          "enabled",
		"proxyrespond_on_loggingerror": "enabled",
		"response_logging":             "enabled",
		"responselog_template":         "resptemplate",
		"requestlog_error_template":    "errtemplate",
		"responselog_error_template":   "resperrtemplate",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileRequestLogCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileRequestLogRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "enabled", d.Get("request_logging").(string))

	diags = resourceBigipLtmProfileRequestLogUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileRequestLogDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileRequestLogCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/request-log", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipLtmProfileRequestLog()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-request-log",
	}, "")

	diags := resourceBigipLtmProfileRequestLogCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileRequestLogRead_Error(t *testing.T) {
	name := "test-request-log"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/request-log/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileRequestLog()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileRequestLogRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileRequestLogUpdate_Error(t *testing.T) {
	name := "test-request-log"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/request-log/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileRequestLog()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileRequestLogUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileRequestLogDelete_Error(t *testing.T) {
	name := "test-request-log"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/request-log/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileRequestLog()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileRequestLogDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getRequestLogProfileConfig
// ---------------------------------------------------------------------

func TestGetRequestLogProfileConfig(t *testing.T) {
	r := resourceBigipLtmProfileRequestLog()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                         "test-request-log",
		"defaults_from":                "/Common/request-log",
		"description":                  "my desc",
		"requestlog_pool":              "/Common/pool1",
		"requestlog_protocol":          "mds-udp",
		"requestlog_error_pool":        "/Common/pool2",
		"requestlog_error_protocol":    "mds-tcp",
		"responselog_pool":             "/Common/pool3",
		"responselog_protocol":         "mds-udp",
		"responselog_error_pool":       "/Common/pool4",
		"responselog_error_protocol":   "mds-tcp",
		"request_logging":              "enabled",
		"response_logging":             "enabled",
		"requestlog_template":          `hello "world"`,
		"requestlog_error_template":    "err",
		"responselog_template":         `resp "world"`,
		"responselog_error_template":   "resperr",
		"proxy_response":               "enabled",
		"proxyclose_on_error":          "enabled",
		"proxyrespond_on_loggingerror": "enabled",
	}, "")

	cfg := getRequestLogProfileConfig(d, &bigip.RequestLogProfile{})
	require.Equal(t, "/Common/request-log", cfg.DefaultsFrom)
	require.Equal(t, "my desc", cfg.Description)
	require.Equal(t, "/Common/pool1", cfg.RequestLogPool)
	require.Equal(t, "mds-udp", cfg.RequestLogProtocol)
	require.Equal(t, "/Common/pool2", cfg.RequestLogErrorPool)
	require.Equal(t, "mds-tcp", cfg.RequestLogErrorProtocol)
	require.Equal(t, "/Common/pool3", cfg.ResponseLogPool)
	require.Equal(t, "mds-udp", cfg.ResponseLogProtocol)
	require.Equal(t, "/Common/pool4", cfg.ResponseLogErrorPool)
	require.Equal(t, "mds-tcp", cfg.ResponseLogErrorProtocol)
	require.Equal(t, "enabled", cfg.RequestLogging)
	require.Equal(t, "enabled", cfg.ResponseLogging)
	require.Equal(t, `hello \"world\"`, cfg.RequestLogTemplate)
	require.Equal(t, "err", cfg.RequestLogErrorTemplate)
	require.Equal(t, `resp \"world\"`, cfg.ResponseLogTemplate)
	require.Equal(t, "resperr", cfg.ResponseLogErrorTemplate)
	require.Equal(t, "enabled", cfg.ProxyResponse)
	require.Equal(t, "enabled", cfg.ProxyCloseOnError)
	require.Equal(t, "enabled", cfg.ProxyRespondOnLoggingError)
}
