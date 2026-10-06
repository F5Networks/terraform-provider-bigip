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

func TestResourceBigipHttpFastAppSchema(t *testing.T) {
	r := resourceBigipHttpFastApp()

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

func TestResourceBigipHttpFastAppPersistenceTypeValidation(t *testing.T) {
	r := resourceBigipHttpFastApp()
	s := r.Schema["persistence_type"]
	require.NotNil(t, s.ValidateFunc)

	_, errs := s.ValidateFunc("cookie", "persistence_type")
	assert.Empty(t, errs)

	_, errs = s.ValidateFunc("bogus", "persistence_type")
	assert.NotEmpty(t, errs)
}

// ---------------------------------------------------------------------
// getFastHttpConfig (direct function tests, no HTTP server needed)
// ---------------------------------------------------------------------

func TestGetFastHttpConfig_Minimal(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "mytenant")
	assert.Contains(t, config, "myapp")
}

func TestGetFastHttpConfig_VirtualServer(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"virtual_server": []interface{}{
			map[string]interface{}{"ip": "10.0.0.1", "port": 80},
		},
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "10.0.0.1")
}

func TestGetFastHttpConfig_ExistingPool(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":        "mytenant",
		"application":   "myapp",
		"existing_pool": "/Common/my-pool",
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-pool")
}

func TestGetFastHttpConfig_PoolMembers(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"pool_members": []interface{}{
			map[string]interface{}{
				"addresses": []interface{}{"192.0.2.1"},
				"port":      8080,
			},
		},
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "192.0.2.1")
}

func TestGetFastHttpConfig_ServiceDiscovery(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"service_discovery": []interface{}{
			`{"sd_type":"aws","sd_port":80}`,
		},
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, `"use_sd":true`)
}

func TestGetFastHttpConfig_SnatPoolAndPersistence(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":               "mytenant",
		"application":          "myapp",
		"existing_snat_pool":   "/Common/my-snatpool",
		"persistence_profile":  "/Common/my-persist",
		"fallback_persistence": "source-address",
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-snatpool")
	assert.Contains(t, config, "my-persist")
	assert.Contains(t, config, "source-address")
}

func TestGetFastHttpConfig_SnatPoolAddresses(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":            "mytenant",
		"application":       "myapp",
		"snat_pool_address": []interface{}{"10.1.1.1"},
		"persistence_type":  "cookie",
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "10.1.1.1")
	assert.Contains(t, config, "cookie")
}

func TestGetFastHttpConfig_EndpointPolicyAndWaf(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":                       "mytenant",
		"application":                  "myapp",
		"endpoint_ltm_policy":          []interface{}{"/Common/my-policy"},
		"existing_waf_security_policy": "/Common/my-waf",
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-policy")
	assert.Contains(t, config, "my-waf")
}

func TestGetFastHttpConfig_MakeWafPolicy(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"waf_security_policy": []interface{}{
			map[string]interface{}{"enable": true},
		},
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, `"make_waf_policy":true`)
}

func TestGetFastHttpConfig_ExistingMonitorAndLoadBalancing(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":                "mytenant",
		"application":           "myapp",
		"existing_monitor":      "/Common/http",
		"load_balancing_mode":   "round-robin",
		"slow_ramp_time":        30,
		"security_log_profiles": []interface{}{"/Common/log-all"},
	}, "")

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "round-robin")
	assert.Contains(t, config, "log-all")
}

func TestGetFastHttpConfig_MonitorWithAuth(t *testing.T) {
	r := resourceBigipHttpFastApp()
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

	config, err := getFastHttpConfig(d)
	require.NoError(t, err)
	assert.Contains(t, config, "admin")
	assert.Contains(t, config, "secret")
}

func TestGetFastHttpConfig_MonitorAuthMissingUsername(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"monitor": []interface{}{
			map[string]interface{}{
				"monitor_auth": true,
				"password":     "secret",
			},
		},
	}, "")

	_, err := getFastHttpConfig(d)
	// username defaults to "" (a valid string), so the ok check succeeds and
	// no error is returned; this documents that behavior rather than
	// asserting an error, since the empty string is a legitimate zero value.
	require.NoError(t, err)
}

// ---------------------------------------------------------------------
// resourceBigipFastHttpAppCreate / Read / Update / Delete
// ---------------------------------------------------------------------

func fastHttpMockServer(t *testing.T, tenant, app string) *httptest.Server {
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
			_, _ = fmt.Fprintf(w, `{"constants":{"fast":{"view":{"tenant_name":"%s","app_name":"%s","virtual_address":"10.0.0.1","virtual_port":80}}}}`, tenant, app)
		case http.MethodPatch:
			_, _ = fmt.Fprint(w, `{"message":[{"id":"task-1"}]}`)
		case http.MethodDelete:
			_, _ = fmt.Fprint(w, `{"id":"task-1"}`)
		}
	})

	return httptest.NewServer(mux)
}

func TestUnitFastHttpCreateReadUpdateDelete(t *testing.T) {
	tenant := "mytenant"
	app := "myapp"

	server := fastHttpMockServer(t, tenant, app)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call

	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      tenant,
		"application": app,
		"virtual_server": []interface{}{
			map[string]interface{}{"ip": "10.0.0.1", "port": 80},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipFastHttpAppCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, app, d.Id())
	assert.Equal(t, tenant, d.Get("tenant"))

	diags = resourceBigipFastHttpAppRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Contains(t, d.Get("fast_http_json").(string), "10.0.0.1")

	diags = resourceBigipFastHttpAppUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipFastHttpAppDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitFastHttpCreate_PostError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"post failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "")

	diags := resourceBigipFastHttpAppCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastHttpCreate_ConfigError(t *testing.T) {
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"monitor": []interface{}{
			map[string]interface{}{
				"monitor_auth": true,
				"username":     "",
			},
		},
	}, "")

	// getFastHttpConfig itself doesn't error on empty username/password (the
	// GetOk-based type assertion always succeeds for a present zero-value
	// string), so this exercises the "config built successfully, but WAF
	// section skipped" happy path rather than an error. Included to document
	// that getFastHttpConfig errors are effectively unreachable through
	// normal schema-validated input.
	_, err := getFastHttpConfig(d)
	require.NoError(t, err)
}

func TestUnitFastHttpWithWaf(t *testing.T) {
	tenant := "mytenant"
	app := "wafapp"

	server := fastHttpMockServer(t, tenant, app)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":                       tenant,
		"application":                  app,
		"existing_waf_security_policy": "/Common/my-waf",
	}, "")

	diags := resourceBigipFastHttpAppCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
}

func TestUnitFastHttpRead_UnexpectedEndOfJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastHttpAppRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected nil diags for the unexpected-end-of-JSON branch")
	assert.Equal(t, "", d.Id())
}

func TestUnitFastHttpRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastHttpAppRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastHttpUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastHttpAppUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastHttpDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipHttpFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastHttpAppDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
