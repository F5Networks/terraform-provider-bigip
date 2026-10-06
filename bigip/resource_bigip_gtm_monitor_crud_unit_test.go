/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newGtmMonitorTestClient returns a *bigip.BigIP configured to talk to the
// supplied httptest server, with retries/timeouts set so the go-bigip
// APICall retry loop actually executes exactly once per request.
func newGtmMonitorTestClient(serverURL string) *bigip.BigIP {
	return bigip.NewSession(&bigip.Config{
		Address:           serverURL,
		Username:          "admin",
		Password:          "admin",
		CertVerifyDisable: true,
		ConfigOptions: &bigip.ConfigOptions{
			APICallRetries: 1,
			APICallTimeout: 5 * time.Second,
		},
	})
}

// gtmMonitorTestResourceData builds a *schema.ResourceData for the given
// resource using schema.TestResourceDataRaw, then sets the Terraform ID to id
// (mimicking state as it would exist after a Create/Import).
func gtmMonitorTestResourceData(t *testing.T, r *schema.Resource, raw map[string]interface{}, id string) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, r.Schema, raw)
	if id != "" {
		d.SetId(id)
	}
	return d
}

// ---------------------------------------------------------------------
// GTM HTTP monitor
// ---------------------------------------------------------------------

func TestUnitGtmMonitorHttpCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-gtm-http-monitor"
	mangled := "/mgmt/tm/gtm/monitor/http/~Common~test-gtm-http-monitor"

	mux := http.NewServeMux()
	var updated, sawPUT, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/monitor/http", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/http","destination":"*:*","interval":30,"timeout":120,"probeTimeout":5,"ignoreDownResponse":"disabled","transparent":"disabled","reverse":"disabled","send":"GET /","recv":"200 OK"}`, name, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			sawPUT = true
			updated = true
			fallthrough
		case "GET":
			if updated {
				_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/http","destination":"10.1.1.100:8080","interval":60,"timeout":180,"probeTimeout":10,"ignoreDownResponse":"enabled","transparent":"enabled","reverse":"enabled","send":"GET /health","recv":"healthy"}`, name, name)
			} else {
				_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/http","destination":"*:*","interval":30,"timeout":120,"probeTimeout":5,"ignoreDownResponse":"disabled","transparent":"disabled","reverse":"disabled","send":"GET /","recv":"200 OK"}`, name, name)
			}
		case "DELETE":
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttp()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":                 name,
		"defaults_from":        "/Common/http",
		"destination":          "*:*",
		"interval":             30,
		"timeout":              120,
		"probe_timeout":        5,
		"ignore_down_response": "disabled",
		"transparent":          "disabled",
		"reverse":              "disabled",
		"send":                 "GET /",
		"receive":              "200 OK",
	})

	ctx := context.Background()

	// Create (which internally invokes Read)
	diags := resourceBigipGtmMonitorHttpCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, name, d.Id())
	assert.Equal(t, "*:*", d.Get("destination"))
	assert.Equal(t, "GET /", d.Get("send"))

	// Explicit Read
	diags = resourceBigipGtmMonitorHttpRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, name, d.Get("name"))

	// Update
	require.NoError(t, d.Set("destination", "10.1.1.100:8080"))
	require.NoError(t, d.Set("interval", 60))
	require.NoError(t, d.Set("timeout", 180))
	require.NoError(t, d.Set("probe_timeout", 10))
	require.NoError(t, d.Set("ignore_down_response", "enabled"))
	require.NoError(t, d.Set("transparent", "enabled"))
	require.NoError(t, d.Set("reverse", "enabled"))
	require.NoError(t, d.Set("send", "GET /health"))
	require.NoError(t, d.Set("receive", "healthy"))

	diags = resourceBigipGtmMonitorHttpUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPUT, "expected a PUT request during update")
	assert.Equal(t, "10.1.1.100:8080", d.Get("destination"))
	assert.Equal(t, 60, d.Get("interval"))
	assert.Equal(t, "healthy", d.Get("receive"))

	// Delete
	diags = resourceBigipGtmMonitorHttpDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete, "expected a DELETE request during delete")
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorHttpReadNotFound(t *testing.T) {
	name := "/Common/missing-http-monitor"
	mangled := "/mgmt/tm/gtm/monitor/http/~Common~missing-http-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttp()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorHttpRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error on not-found read: %v", diags)
	assert.Equal(t, "", d.Id(), "expected resource id to be cleared when not found")
}

func TestUnitGtmMonitorHttpCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/monitor/http", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "the requested object already exists", http.StatusConflict)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttp()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name": "/Common/dup-http-monitor",
	})

	diags := resourceBigipGtmMonitorHttpCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmMonitorHttpDeleteNotFound(t *testing.T) {
	name := "/Common/already-gone-http-monitor"
	mangled := "/mgmt/tm/gtm/monitor/http/~Common~already-gone-http-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttp()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorHttpDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error deleting an already-missing monitor: %v", diags)
	assert.Equal(t, "", d.Id())
}

// ---------------------------------------------------------------------
// GTM HTTPS monitor
// ---------------------------------------------------------------------

func TestUnitGtmMonitorHttpsCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-gtm-https-monitor"
	mangled := "/mgmt/tm/gtm/monitor/https/~Common~test-gtm-https-monitor"

	mux := http.NewServeMux()
	var updated, sawPUT, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/monitor/https", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/https","destination":"*:*","interval":30,"timeout":120,"probeTimeout":5,"ignoreDownResponse":"disabled","transparent":"disabled","reverse":"disabled","send":"GET /","recv":"200 OK","cipherlist":"DEFAULT:!EXPORT","compatibility":"enabled","sniServerName":"none"}`, name, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			sawPUT = true
			updated = true
			fallthrough
		case "GET":
			if updated {
				_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/https","destination":"10.1.1.100:8443","interval":45,"timeout":180,"probeTimeout":10,"ignoreDownResponse":"enabled","transparent":"enabled","reverse":"disabled","send":"GET /api/health","recv":"status: ok","cert":"/Common/cert.crt","key":"/Common/cert.key","cipherlist":"HIGH:!ADH:!MD5","compatibility":"disabled","sniServerName":"example.com"}`, name, name)
			} else {
				_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/https","destination":"*:*","interval":30,"timeout":120,"probeTimeout":5,"ignoreDownResponse":"disabled","transparent":"disabled","reverse":"disabled","send":"GET /","recv":"200 OK","cipherlist":"DEFAULT:!EXPORT","compatibility":"enabled","sniServerName":"none"}`, name, name)
			}
		case "DELETE":
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttps()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":                 name,
		"defaults_from":        "/Common/https",
		"destination":          "*:*",
		"interval":             30,
		"timeout":              120,
		"probe_timeout":        5,
		"ignore_down_response": "disabled",
		"transparent":          "disabled",
		"reverse":              "disabled",
		"send":                 "GET /",
		"receive":              "200 OK",
		"cipherlist":           "DEFAULT:!EXPORT",
		"compatibility":        "enabled",
		"sni_server_name":      "none",
	})

	ctx := context.Background()

	diags := resourceBigipGtmMonitorHttpsCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, name, d.Id())
	assert.Equal(t, "DEFAULT:!EXPORT", d.Get("cipherlist"))
	assert.Equal(t, "none", d.Get("sni_server_name"))

	diags = resourceBigipGtmMonitorHttpsRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	require.NoError(t, d.Set("destination", "10.1.1.100:8443"))
	require.NoError(t, d.Set("interval", 45))
	require.NoError(t, d.Set("cipherlist", "HIGH:!ADH:!MD5"))
	require.NoError(t, d.Set("compatibility", "disabled"))
	require.NoError(t, d.Set("cert", "/Common/cert.crt"))
	require.NoError(t, d.Set("key", "/Common/cert.key"))
	require.NoError(t, d.Set("sni_server_name", "example.com"))

	diags = resourceBigipGtmMonitorHttpsUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPUT, "expected a PUT request during update")
	assert.Equal(t, "10.1.1.100:8443", d.Get("destination"))
	assert.Equal(t, "HIGH:!ADH:!MD5", d.Get("cipherlist"))
	assert.Equal(t, "disabled", d.Get("compatibility"))
	assert.Equal(t, "/Common/cert.crt", d.Get("cert"))
	assert.Equal(t, "example.com", d.Get("sni_server_name"))

	diags = resourceBigipGtmMonitorHttpsDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete, "expected a DELETE request during delete")
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorHttpsReadNotFound(t *testing.T) {
	name := "/Common/missing-https-monitor"
	mangled := "/mgmt/tm/gtm/monitor/https/~Common~missing-https-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttps()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorHttpsRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error on not-found read: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorHttpsUpdateError(t *testing.T) {
	name := "/Common/test-gtm-https-monitor-err"
	mangled := "/mgmt/tm/gtm/monitor/https/~Common~test-gtm-https-monitor-err"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttps()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorHttpsUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmMonitorHttpsDeleteError(t *testing.T) {
	name := "/Common/test-gtm-https-monitor-delerr"
	mangled := "/mgmt/tm/gtm/monitor/https/~Common~test-gtm-https-monitor-delerr"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttps()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorHttpsDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmMonitorHttpsDeleteNotFound(t *testing.T) {
	name := "/Common/already-gone-https-monitor"
	mangled := "/mgmt/tm/gtm/monitor/https/~Common~already-gone-https-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttps()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorHttpsDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error deleting an already-missing monitor: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorHttpsCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/monitor/https", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "the requested object already exists", http.StatusConflict)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorHttps()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name": "/Common/dup-https-monitor",
	})

	diags := resourceBigipGtmMonitorHttpsCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// GTM TCP monitor
// ---------------------------------------------------------------------

func TestUnitGtmMonitorTcpCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-gtm-tcp-monitor"
	mangled := "/mgmt/tm/gtm/monitor/tcp/~Common~test-gtm-tcp-monitor"

	mux := http.NewServeMux()
	var sawPUT, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/monitor/tcp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/tcp","destination":"*:*","interval":30,"timeout":120,"probeTimeout":5,"ignoreDownResponse":"disabled","transparent":"disabled","reverse":"disabled"}`, name, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			sawPUT = true
			fallthrough
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/tcp","destination":"10.1.1.50:80","interval":15,"timeout":60,"probeTimeout":3,"ignoreDownResponse":"enabled","transparent":"enabled","reverse":"enabled","send":"PING","recv":"PONG"}`, name, name)
		case "DELETE":
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorTcp()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":                 name,
		"defaults_from":        "/Common/tcp",
		"destination":          "*:*",
		"interval":             30,
		"timeout":              120,
		"probe_timeout":        5,
		"ignore_down_response": "disabled",
		"transparent":          "disabled",
		"reverse":              "disabled",
	})

	ctx := context.Background()

	diags := resourceBigipGtmMonitorTcpCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, name, d.Id())

	diags = resourceBigipGtmMonitorTcpRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	require.NoError(t, d.Set("destination", "10.1.1.50:80"))
	require.NoError(t, d.Set("interval", 15))
	require.NoError(t, d.Set("send", "PING"))
	require.NoError(t, d.Set("receive", "PONG"))

	diags = resourceBigipGtmMonitorTcpUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPUT, "expected a PUT request during update")
	assert.Equal(t, "10.1.1.50:80", d.Get("destination"))
	assert.Equal(t, "PING", d.Get("send"))

	diags = resourceBigipGtmMonitorTcpDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete, "expected a DELETE request during delete")
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorTcpReadNotFound(t *testing.T) {
	name := "/Common/missing-tcp-monitor"
	mangled := "/mgmt/tm/gtm/monitor/tcp/~Common~missing-tcp-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorTcp()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorTcpRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error on not-found read: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorTcpReadError(t *testing.T) {
	name := "/Common/broken-tcp-monitor"
	mangled := "/mgmt/tm/gtm/monitor/tcp/~Common~broken-tcp-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorTcp()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorTcpRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmMonitorTcpMinimalDefaults(t *testing.T) {
	name := "/Common/test-gtm-tcp-monitor-minimal"
	mangled := "/mgmt/tm/gtm/monitor/tcp/~Common~test-gtm-tcp-monitor-minimal"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/monitor/tcp", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/tcp","destination":"*:*","interval":30,"timeout":120}`, name, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/tcp","destination":"*:*","interval":30,"timeout":120}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorTcp()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name": name,
	})

	diags := resourceBigipGtmMonitorTcpCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, "/Common/tcp", d.Get("defaults_from"))
	assert.Equal(t, "*:*", d.Get("destination"))
	assert.Equal(t, 30, d.Get("interval"))
	assert.Equal(t, 120, d.Get("timeout"))
}

func TestUnitGtmMonitorTcpDeleteError(t *testing.T) {
	name := "/Common/test-gtm-tcp-monitor-delerr"
	mangled := "/mgmt/tm/gtm/monitor/tcp/~Common~test-gtm-tcp-monitor-delerr"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorTcp()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorTcpDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmMonitorTcpDeleteNotFound(t *testing.T) {
	name := "/Common/already-gone-tcp-monitor"
	mangled := "/mgmt/tm/gtm/monitor/tcp/~Common~already-gone-tcp-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorTcp()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorTcpDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error deleting an already-missing monitor: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorTcpCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/monitor/tcp", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "the requested object already exists", http.StatusConflict)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorTcp()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name": "/Common/dup-tcp-monitor",
	})

	diags := resourceBigipGtmMonitorTcpCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmMonitorTcpUpdateError(t *testing.T) {
	name := "/Common/test-gtm-tcp-monitor-updateerr"
	mangled := "/mgmt/tm/gtm/monitor/tcp/~Common~test-gtm-tcp-monitor-updateerr"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorTcp()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorTcpUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// GTM PostgreSQL monitor
// ---------------------------------------------------------------------

func TestUnitGtmMonitorPostgresqlCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-gtm-postgresql-monitor"
	mangled := "/mgmt/tm/gtm/monitor/postgresql/~Common~test-gtm-postgresql-monitor"

	mux := http.NewServeMux()
	var updated, sawPUT, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/monitor/postgresql", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/postgresql","destination":"*:5432","interval":30,"timeout":91,"probeTimeout":5,"ignoreDownResponse":"disabled","database":"testdb","username":"testuser","count":"1","debug":"no"}`, name, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			sawPUT = true
			updated = true
			fallthrough
		case "GET":
			if updated {
				_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/postgresql","destination":"192.168.1.50:5432","interval":45,"timeout":91,"probeTimeout":10,"ignoreDownResponse":"enabled","database":"production","username":"produser","count":"1","debug":"yes"}`, name, name)
			} else {
				_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/postgresql","destination":"*:5432","interval":30,"timeout":91,"probeTimeout":5,"ignoreDownResponse":"disabled","database":"testdb","username":"testuser","count":"1","debug":"no"}`, name, name)
			}
		case "DELETE":
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorPostgresql()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":                 name,
		"defaults_from":        "/Common/postgresql",
		"destination":          "*:5432",
		"interval":             30,
		"timeout":              91,
		"probe_timeout":        5,
		"ignore_down_response": "disabled",
		"database":             "testdb",
		"username":             "testuser",
		"password":             "secretpw",
		"debug":                "no",
	})

	ctx := context.Background()

	diags := resourceBigipGtmMonitorPostgresqlCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, name, d.Id())
	assert.Equal(t, "testdb", d.Get("database"))

	diags = resourceBigipGtmMonitorPostgresqlRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	// Password must never be set back from Read (sensitive/not returned by API).
	assert.Equal(t, "secretpw", d.Get("password"))

	require.NoError(t, d.Set("destination", "192.168.1.50:5432"))
	require.NoError(t, d.Set("interval", 45))
	require.NoError(t, d.Set("database", "production"))
	require.NoError(t, d.Set("username", "produser"))
	require.NoError(t, d.Set("debug", "yes"))

	diags = resourceBigipGtmMonitorPostgresqlUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPUT, "expected a PUT request during update")
	assert.Equal(t, "192.168.1.50:5432", d.Get("destination"))
	assert.Equal(t, "production", d.Get("database"))
	assert.Equal(t, "yes", d.Get("debug"))

	diags = resourceBigipGtmMonitorPostgresqlDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete, "expected a DELETE request during delete")
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorPostgresqlReadNotFound(t *testing.T) {
	name := "/Common/missing-postgresql-monitor"
	mangled := "/mgmt/tm/gtm/monitor/postgresql/~Common~missing-postgresql-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorPostgresql()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorPostgresqlRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error on not-found read: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorPostgresqlCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/monitor/postgresql", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "the requested object already exists", http.StatusConflict)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorPostgresql()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name": "/Common/dup-postgresql-monitor",
	})

	diags := resourceBigipGtmMonitorPostgresqlCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmMonitorPostgresqlDeleteError(t *testing.T) {
	name := "/Common/test-gtm-postgresql-monitor-delerr"
	mangled := "/mgmt/tm/gtm/monitor/postgresql/~Common~test-gtm-postgresql-monitor-delerr"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorPostgresql()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorPostgresqlDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmMonitorPostgresqlDeleteNotFound(t *testing.T) {
	name := "/Common/already-gone-postgresql-monitor"
	mangled := "/mgmt/tm/gtm/monitor/postgresql/~Common~already-gone-postgresql-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorPostgresql()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorPostgresqlDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error deleting an already-missing monitor: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorPostgresqlUpdateError(t *testing.T) {
	name := "/Common/test-gtm-postgresql-monitor-updateerr"
	mangled := "/mgmt/tm/gtm/monitor/postgresql/~Common~test-gtm-postgresql-monitor-updateerr"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorPostgresql()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorPostgresqlUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// GTM BIG-IP monitor
// ---------------------------------------------------------------------

func TestUnitGtmMonitorBigipCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-gtm-bigip-monitor"
	mangled := "/mgmt/tm/gtm/monitor/bigip/~Common~test-gtm-bigip-monitor"

	mux := http.NewServeMux()
	var updated, sawPUT, sawDelete bool
	var lastCreateBody string

	mux.HandleFunc("/mgmt/tm/gtm/monitor/bigip", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		body, _ := io.ReadAll(r.Body)
		lastCreateBody = string(body)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/bigip","destination":"*:*","interval":30,"timeout":90,"ignoreDownResponse":"disabled","aggregateDynamicRatios":"none"}`, name, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			sawPUT = true
			updated = true
			fallthrough
		case "GET":
			if updated {
				_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/bigip","destination":"10.1.1.60:443","interval":20,"timeout":95,"ignoreDownResponse":"enabled","aggregateDynamicRatios":"average-nodes"}`, name, name)
			} else {
				_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","defaultsFrom":"/Common/bigip","destination":"*:*","interval":30,"timeout":90,"ignoreDownResponse":"disabled","aggregateDynamicRatios":"none"}`, name, name)
			}
		case "DELETE":
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorBigip()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":                 name,
		"defaults_from":        "/Common/bigip",
		"destination":          "*:*",
		"interval":             30,
		"timeout":              90,
		"ignore_down_response": "disabled",
		"aggregation_type":     "none",
	})

	ctx := context.Background()

	diags := resourceBigipGtmMonitorBigipCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, name, d.Id())
	assert.Contains(t, lastCreateBody, `"destination":"*:*"`)

	diags = resourceBigipGtmMonitorBigipRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	require.NoError(t, d.Set("destination", "10.1.1.60:443"))
	require.NoError(t, d.Set("interval", 20))
	require.NoError(t, d.Set("timeout", 95))
	require.NoError(t, d.Set("aggregation_type", "average-nodes"))

	diags = resourceBigipGtmMonitorBigipUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPUT, "expected a PUT request during update")
	assert.Equal(t, "10.1.1.60:443", d.Get("destination"))
	assert.Equal(t, "average-nodes", d.Get("aggregation_type"))

	diags = resourceBigipGtmMonitorBigipDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete, "expected a DELETE request during delete")
	assert.Equal(t, "", d.Id())
}

// TestUnitGtmMonitorBigipDestinationNormalization verifies that a bare
// wildcard address destination ("*", with no port) is normalized to "*:*"
// before being sent to the API, exercising normalizeGtmBigipDestination from
// within Create. (A bare specific IP address such as "10.1.1.100" would
// normalize to "10.1.1.100:*", which validateGtmBigipDestination correctly
// rejects — see TestUnitGtmMonitorBigipInvalidDestination.)
func TestUnitGtmMonitorBigipDestinationNormalization(t *testing.T) {
	name := "/Common/test-gtm-bigip-monitor-norm"
	mangled := "/mgmt/tm/gtm/monitor/bigip/~Common~test-gtm-bigip-monitor-norm"

	mux := http.NewServeMux()
	var lastCreateBody string

	mux.HandleFunc("/mgmt/tm/gtm/monitor/bigip", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		lastCreateBody = string(body)
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","destination":"*:*"}`, name, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","destination":"*:*"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorBigip()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":        name,
		"destination": "*",
	})

	diags := resourceBigipGtmMonitorBigipCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Contains(t, lastCreateBody, `"destination":"*:*"`)
}

// TestUnitGtmMonitorBigipInvalidDestination verifies Create/Update reject a
// destination with a specific IP and wildcard port, per BIG-IP API rules.
func TestUnitGtmMonitorBigipInvalidDestination(t *testing.T) {
	r := resourceBigipGtmMonitorBigip()

	t.Run("create", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
			"name":        "/Common/test-gtm-bigip-invalid",
			"destination": "10.1.1.50:*",
		})
		client := newGtmMonitorTestClient("http://127.0.0.1:1")
		diags := resourceBigipGtmMonitorBigipCreate(context.Background(), d, client)
		require.True(t, diags.HasError())
	})

	t.Run("update", func(t *testing.T) {
		d := gtmMonitorTestResourceData(t, r, map[string]interface{}{
			"name":        "/Common/test-gtm-bigip-invalid",
			"destination": "10.1.1.50:*",
		}, "/Common/test-gtm-bigip-invalid")
		client := newGtmMonitorTestClient("http://127.0.0.1:1")
		diags := resourceBigipGtmMonitorBigipUpdate(context.Background(), d, client)
		require.True(t, diags.HasError())
	})
}

func TestUnitGtmMonitorBigipReadNotFound(t *testing.T) {
	name := "/Common/missing-bigip-monitor"
	mangled := "/mgmt/tm/gtm/monitor/bigip/~Common~missing-bigip-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorBigip()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorBigipRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error on not-found read: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorBigipDeleteNotFound(t *testing.T) {
	name := "/Common/already-gone-bigip-monitor"
	mangled := "/mgmt/tm/gtm/monitor/bigip/~Common~already-gone-bigip-monitor"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorBigip()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorBigipDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected no error deleting an already-missing monitor: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmMonitorBigipUpdateError(t *testing.T) {
	name := "/Common/test-gtm-bigip-monitor-err"
	mangled := "/mgmt/tm/gtm/monitor/bigip/~Common~test-gtm-bigip-monitor-err"

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newGtmMonitorTestClient(server.URL)

	r := resourceBigipGtmMonitorBigip()
	d := gtmMonitorTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmMonitorBigipUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// TestUnitValidateGtmBigipDestinationMalformed exercises the "no colon"
// branch of validateGtmBigipDestination, which is not covered by the
// existing table-driven schema unit tests.
func TestUnitValidateGtmBigipDestinationMalformed(t *testing.T) {
	err := validateGtmBigipDestination("not-a-valid-destination")
	require.Error(t, err)
}
