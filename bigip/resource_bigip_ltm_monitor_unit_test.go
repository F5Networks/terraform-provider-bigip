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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// monitorTypeURIs mirrors the list of monitor type URIs go-bigip's
// Monitors() iterates over, so tests can register empty-list handlers for
// every type and override just the one under test.
var monitorTypeURIs = []string{"http", "https", "icmp", "gateway-icmp", "tcp", "tcp-half-open", "ftp", "udp", "postgresql", "mysql", "mssql", "ldap", "smtp"}

// newMonitorTestMux registers empty {"items":[]} handlers for every monitor
// type endpoint under /mgmt/tm/ltm/monitor/<type>, except those present in
// overrides (which are registered with the supplied handler instead). This
// avoids duplicate-pattern panics from http.ServeMux when a test needs to
// customize the response for a specific monitor type.
func newMonitorTestMux(overrides ...map[string]http.HandlerFunc) *http.ServeMux {
	var ov map[string]http.HandlerFunc
	if len(overrides) > 0 {
		ov = overrides[0]
	}
	mux := http.NewServeMux()
	for _, mt := range monitorTypeURIs {
		path := "/mgmt/tm/ltm/monitor/" + mt
		if h, ok := ov[mt]; ok {
			mux.HandleFunc(path, h)
			continue
		}
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"items":[]}`)
		})
	}
	return mux
}

func ltmMonitorResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceBigipLtmMonitor().Schema, raw)
}

// ---------- Pure helper function tests ----------

func TestResourceBigipLtmMonitorValidateParent(t *testing.T) {
	validParents := []string{
		"/Common/udp", "/Common/postgresql", "/Common/mysql", "/Common/mssql",
		"/Common/http", "/Common/https", "/Common/icmp", "/Common/gateway_icmp",
		"/Common/tcp", "/Common/tcp_half_open", "/Common/ftp", "/Common/ldap",
		"/Common/smtp", "/Common/dns",
	}
	for _, p := range validParents {
		warnings, errs := validateParent(p, "parent")
		assert.Empty(t, warnings)
		assert.Empty(t, errs, "expected %s to be valid", p)
	}

	warnings, errs := validateParent("/Common/bogus", "parent")
	assert.Empty(t, warnings)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "parent must be one of")
}

func TestResourceBigipLtmMonitorMonitorParent(t *testing.T) {
	assert.Equal(t, "http", monitorParent("/Common/http"))
	assert.Equal(t, "gateway_icmp", monitorParent("/Common/gateway_icmp"))
	assert.Equal(t, "tcp_half_open", monitorParent("/Common/tcp_half_open"))
	assert.Equal(t, "custom", monitorParent("custom"))
}

func TestResourceBigipLtmMonitorGetLtmMonitorConfig(t *testing.T) {
	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":            "test-http-monitor",
		"parent":          "/Common/http",
		"interval":        5,
		"up_interval":     0,
		"timeout":         16,
		"send":            "GET /\\r\\n",
		"receive":         "200 OK",
		"receive_disable": "",
		"reverse":         "disabled",
		"transparent":     "disabled",
		"manual_resume":   "disabled",
		"ip_dscp":         0,
		"time_until_up":   0,
		"destination":     "*:*",
		"compatibility":   "enabled",
	})

	pss := &bigip.Monitor{}
	config := getLtmMonitorConfig(d, pss)
	assert.Equal(t, "/Common/http", config.ParentMonitor)
	assert.Equal(t, "GET /\\r\\n", config.SendString)
	assert.Equal(t, "200 OK", config.ReceiveString)
	assert.Equal(t, 5, config.Interval)
	assert.Equal(t, 16, config.Timeout)
	assert.Equal(t, "*:*", config.Destination)
}

func TestResourceBigipLtmMonitorGetLtmMonitorConfigCustomParent(t *testing.T) {
	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":          "test-monitor",
		"parent":        "/Common/http",
		"custom_parent": "/Common/my_custom_http",
	})
	pss := &bigip.Monitor{}
	config := getLtmMonitorConfig(d, pss)
	assert.Equal(t, "/Common/my_custom_http", config.ParentMonitor)
}

func TestResourceBigipLtmMonitorGetLtmMonitorConfigSMTPDomain(t *testing.T) {
	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "test-smtp-monitor",
		"parent": "/Common/smtp",
		"domain": "example.com",
	})
	pss := &bigip.Monitor{}
	config := getLtmMonitorConfig(d, pss)
	assert.Equal(t, "example.com", config.Domain)
}

func TestResourceBigipLtmMonitorGetLtmMonitorConfigDomainIgnoredForNonSMTP(t *testing.T) {
	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "test-http-monitor",
		"parent": "/Common/http",
		"domain": "example.com",
	})
	pss := &bigip.Monitor{}
	config := getLtmMonitorConfig(d, pss)
	assert.Empty(t, config.Domain)
}

// ---------- CRUD tests ----------

func newMonitorTestClient(serverURL string) *bigip.BigIP {
	return newDatasourceTestClient(serverURL)
}

func TestResourceBigipLtmMonitorCreateSuccess(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"http": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "POST":
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprint(w, `{"name":"test-http-monitor","fullPath":"/Common/test-http-monitor"}`)
			default:
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprint(w, `{"items":[{"name":"test-http-monitor","fullPath":"/Common/test-http-monitor","defaultsFrom":"/Common/http"}]}`)
			}
		},
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-http-monitor",
		"parent": "/Common/http",
	})

	diags := resourceBigipLtmMonitorCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/test-http-monitor", d.Id())
}

func TestResourceBigipLtmMonitorCreateError(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"http": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			_, _ = fmt.Fprint(w, `{"items":[]}`)
		},
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-http-monitor",
		"parent": "/Common/http",
	})

	diags := resourceBigipLtmMonitorCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmMonitorCreateGatewayICMP(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"gateway-icmp": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "POST":
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprint(w, `{"name":"test-gw-monitor","fullPath":"/Common/test-gw-monitor"}`)
			default:
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprint(w, `{"items":[{"name":"test-gw-monitor","fullPath":"/Common/test-gw-monitor","defaultsFrom":"/Common/gateway_icmp"}]}`)
			}
		},
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-gw-monitor",
		"parent": "/Common/gateway_icmp",
	})

	diags := resourceBigipLtmMonitorCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

func TestResourceBigipLtmMonitorCreateTCPHalfOpen(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"tcp-half-open": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "POST":
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprint(w, `{"name":"test-tcpho-monitor","fullPath":"/Common/test-tcpho-monitor"}`)
			default:
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprint(w, `{"items":[{"name":"test-tcpho-monitor","fullPath":"/Common/test-tcpho-monitor","defaultsFrom":"/Common/tcp_half_open"}]}`)
			}
		},
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-tcpho-monitor",
		"parent": "/Common/tcp_half_open",
	})

	diags := resourceBigipLtmMonitorCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

func TestResourceBigipLtmMonitorReadSuccess(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"https": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"items":[{"name":"test-https-monitor","fullPath":"/Common/test-https-monitor","defaultsFrom":"/Common/https","sslProfile":"clientssl","compatibility":"enabled"}]}`)
		},
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-https-monitor",
		"parent": "/Common/https",
	})
	d.SetId("/Common/test-https-monitor")

	diags := resourceBigipLtmMonitorRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "clientssl", d.Get("ssl_profile"))
	assert.Equal(t, "enabled", d.Get("compatibility"))
}

func TestResourceBigipLtmMonitorReadNotFoundInList(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"http": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"items":[{"name":"other-monitor","fullPath":"/Common/other-monitor","defaultsFrom":"/Common/http"}]}`)
		},
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/missing-monitor",
		"parent": "/Common/http",
	})
	d.SetId("/Common/missing-monitor")

	diags := resourceBigipLtmMonitorRead(context.Background(), d, client)
	require.True(t, diags.HasError())
	assert.Contains(t, diags[0].Summary, "Couldn't find LTM Monitor")
}

func TestResourceBigipLtmMonitorReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/monitor/http", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-http-monitor",
		"parent": "/Common/http",
	})
	d.SetId("/Common/test-http-monitor")

	diags := resourceBigipLtmMonitorRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmMonitorUpdateSuccess(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"http": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"items":[{"name":"test-http-monitor","fullPath":"/Common/test-http-monitor","defaultsFrom":"/Common/http"}]}`)
		},
	})
	mux.HandleFunc("/mgmt/tm/ltm/monitor/http/~Common~test-http-monitor", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"test-http-monitor","fullPath":"/Common/test-http-monitor"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-http-monitor",
		"parent": "/Common/http",
	})
	d.SetId("/Common/test-http-monitor")

	diags := resourceBigipLtmMonitorUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

func TestResourceBigipLtmMonitorUpdateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/monitor/http/~Common~test-http-monitor", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-http-monitor",
		"parent": "/Common/http",
	})
	d.SetId("/Common/test-http-monitor")

	diags := resourceBigipLtmMonitorUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmMonitorUpdateGatewayICMP(t *testing.T) {
	mux := newMonitorTestMux(map[string]http.HandlerFunc{
		"gateway-icmp": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"items":[{"name":"test-gw-monitor","fullPath":"/Common/test-gw-monitor","defaultsFrom":"/Common/gateway_icmp"}]}`)
		},
	})
	mux.HandleFunc("/mgmt/tm/ltm/monitor/gateway-icmp/~Common~test-gw-monitor", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"test-gw-monitor","fullPath":"/Common/test-gw-monitor"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-gw-monitor",
		"parent": "/Common/gateway_icmp",
	})
	d.SetId("/Common/test-gw-monitor")

	diags := resourceBigipLtmMonitorUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

func TestResourceBigipLtmMonitorDeleteSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/monitor/http/~Common~test-http-monitor", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-http-monitor",
		"parent": "/Common/http",
	})
	d.SetId("/Common/test-http-monitor")

	diags := resourceBigipLtmMonitorDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Id())
}

func TestResourceBigipLtmMonitorDeleteError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/monitor/http/~Common~test-http-monitor", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-http-monitor",
		"parent": "/Common/http",
	})
	d.SetId("/Common/test-http-monitor")

	diags := resourceBigipLtmMonitorDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmMonitorDeleteGatewayICMP(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/monitor/gateway-icmp/~Common~test-gw-monitor", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-gw-monitor",
		"parent": "/Common/gateway_icmp",
	})
	d.SetId("/Common/test-gw-monitor")

	diags := resourceBigipLtmMonitorDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

func TestResourceBigipLtmMonitorDeleteTCPHalfOpen(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/monitor/tcp-half-open/~Common~test-tcpho-monitor", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newMonitorTestClient(server.URL)

	d := ltmMonitorResourceData(t, map[string]interface{}{
		"name":   "/Common/test-tcpho-monitor",
		"parent": "/Common/tcp_half_open",
	})
	d.SetId("/Common/test-tcpho-monitor")

	diags := resourceBigipLtmMonitorDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}
