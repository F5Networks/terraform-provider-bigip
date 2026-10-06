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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipFastHTTPSAppSchema(t *testing.T) {
	r := resourceBigipFastHTTPSApp()

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
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	tenantSchema, ok := r.Schema["tenant"]
	if !ok {
		t.Fatal("Expected field 'tenant' to exist in schema")
	}
	if !tenantSchema.Required {
		t.Error("Expected field 'tenant' to be required")
	}
}

// ---------------------------------------------------------------------
// getFastHTTPSConfig (direct function tests, no HTTP server needed)
// ---------------------------------------------------------------------

func TestGetFastHTTPSConfig_Minimal(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "mytenant")
	assert.Contains(t, config, `"enable_tls_server":true`)
}

func TestGetFastHTTPSConfig_ExistingTlsServerProfile(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":                      "mytenant",
		"application":                 "myapp",
		"existing_tls_server_profile": "/Common/clientssl",
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "clientssl")
}

func TestGetFastHTTPSConfig_MakeTlsServerProfile(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"tls_server_profile": []interface{}{
			map[string]interface{}{
				"tls_cert_name": "/Common/my-cert",
				"tls_key_name":  "/Common/my-key",
			},
		},
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-cert")
	assert.Contains(t, config, "my-key")
	assert.Contains(t, config, `"make_tls_server_profile":true`)
}

func TestGetFastHTTPSConfig_ExistingTlsClientProfile(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":                      "mytenant",
		"application":                 "myapp",
		"existing_tls_client_profile": "/Common/serverssl",
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "serverssl")
	assert.Contains(t, config, `"enable_tls_client":true`)
}

func TestGetFastHTTPSConfig_MakeTlsClientProfile(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"tls_client_profile": []interface{}{
			map[string]interface{}{
				"tls_cert_name": "/Common/client-cert",
				"tls_key_name":  "/Common/client-key",
			},
		},
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "client-cert")
	assert.Contains(t, config, `"make_tls_client_profile":true`)
}

func TestGetFastHTTPSConfig_EndpointPolicyAndWaf(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":                       "mytenant",
		"application":                  "myapp",
		"endpoint_ltm_policy":          []interface{}{"/Common/my-policy"},
		"existing_waf_security_policy": "/Common/my-waf",
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-policy")
	assert.Contains(t, config, "my-waf")
}

func TestGetFastHTTPSConfig_MakeWafPolicy(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"waf_security_policy": []interface{}{
			map[string]interface{}{"enable": true},
		},
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, `"make_waf_policy":true`)
}

func TestGetFastHTTPSConfig_PersistenceAndFallback(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":               "mytenant",
		"application":          "myapp",
		"persistence_profile":  "/Common/my-persist",
		"fallback_persistence": "source-address",
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-persist")
	assert.Contains(t, config, "source-address")
}

func TestGetFastHTTPSConfig_PersistenceType(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":           "mytenant",
		"application":      "myapp",
		"persistence_type": "cookie",
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "cookie")
}

func TestGetFastHTTPSConfig_PoolMembersAndServiceDiscovery(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"pool_members": []interface{}{
			map[string]interface{}{
				"addresses": []interface{}{"192.0.2.1"},
				"port":      8443,
			},
		},
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "192.0.2.1")
}

func TestGetFastHTTPSConfig_SnatPoolAndAddresses(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":             "mytenant",
		"application":        "myapp",
		"existing_snat_pool": "/Common/my-snatpool",
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-snatpool")
}

func TestGetFastHTTPSConfig_SnatPoolAddresses(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":            "mytenant",
		"application":       "myapp",
		"snat_pool_address": []interface{}{"10.1.1.1"},
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "10.1.1.1")
}

func TestGetFastHTTPSConfig_MonitorWithAuth(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"monitor": []interface{}{
			map[string]interface{}{
				"monitor_auth": true,
				"username":     "admin",
				"password":     "secret",
				"interval":     10,
				"send_string":  "GET /",
				"response":     "200 OK",
			},
		},
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "admin")
}

func TestGetFastHTTPSConfig_ExistingMonitorAndLoadBalancing(t *testing.T) {
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":                "mytenant",
		"application":           "myapp",
		"existing_monitor":      "/Common/https",
		"load_balancing_mode":   "round-robin",
		"slow_ramp_time":        30,
		"security_log_profiles": []interface{}{"/Common/log-all"},
	}, "")

	config, err := getFastHTTPSConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "round-robin")
	assert.Contains(t, config, "log-all")
}

// ---------------------------------------------------------------------
// resourceBigipFastHTTPSAppCreate / Read / Update / Delete
// ---------------------------------------------------------------------

func fastHTTPSMockServer(t *testing.T, tenant, app string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			_, _ = fmt.Fprint(w, `{"message":[{"id":"task-1"}]}`)
		}
	})
	mux.HandleFunc("/mgmt/shared/fast/tasks/task-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"id":"task-1","code":200,"message":"success","tenant":"%s","application":"%s"}`, tenant, app)
	})
	mux.HandleFunc(fmt.Sprintf("/mgmt/shared/fast/applications/%s/%s", tenant, app), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"constants":{"fast":{"view":{"tenant_name":"%s","app_name":"%s","virtual_address":"10.0.0.1","virtual_port":443}}}}`, tenant, app)
		case http.MethodPatch:
			_, _ = fmt.Fprint(w, `{"message":[{"id":"task-1"}]}`)
		case http.MethodDelete:
			_, _ = fmt.Fprint(w, `{"id":"task-1"}`)
		}
	})

	return httptest.NewServer(mux)
}

func TestUnitFastHTTPSCreateReadUpdateDelete(t *testing.T) {
	tenant := "mytenant"
	app := "myapp"

	server := fastHTTPSMockServer(t, tenant, app)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call

	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      tenant,
		"application": app,
		"virtual_server": []interface{}{
			map[string]interface{}{"ip": "10.0.0.1", "port": 443},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipFastHTTPSAppCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, app, d.Id())
	assert.Equal(t, tenant, d.Get("tenant"))

	diags = resourceBigipFastHTTPSAppRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Contains(t, d.Get("fast_https_json").(string), "10.0.0.1")

	diags = resourceBigipFastHTTPSAppUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipFastHTTPSAppDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitFastHTTPSCreate_WithWaf(t *testing.T) {
	tenant := "mytenant"
	app := "wafapp"

	server := fastHTTPSMockServer(t, tenant, app)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":                       tenant,
		"application":                  app,
		"existing_waf_security_policy": "/Common/my-waf",
	}, "")

	diags := resourceBigipFastHTTPSAppCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
}

func TestUnitFastHTTPSCreate_PostError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"post failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "")

	diags := resourceBigipFastHTTPSAppCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastHTTPSRead_UnexpectedEndOfJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastHTTPSAppRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected nil diags for the unexpected-end-of-JSON branch")
	assert.Equal(t, "", d.Id())
}

func TestUnitFastHTTPSRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastHTTPSAppRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastHTTPSUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastHTTPSAppUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastHTTPSDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastHTTPSApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastHTTPSAppDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
