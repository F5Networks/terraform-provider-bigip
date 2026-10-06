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

func TestResourceBigipFastTcpAppSchema(t *testing.T) {
	r := resourceBigipFastTcpApp()

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

	appSchema, ok := r.Schema["application"]
	if !ok {
		t.Fatal("Expected field 'application' to exist in schema")
	}
	if !appSchema.Required {
		t.Error("Expected field 'application' to be required")
	}
	if !appSchema.ForceNew {
		t.Error("Expected field 'application' to be ForceNew")
	}
}

// ---------------------------------------------------------------------
// getParamsConfigMap (direct function tests, no HTTP server needed)
// ---------------------------------------------------------------------

func TestGetParamsConfigMap_Minimal(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, "mytenant")
	assert.Contains(t, config, "myapp")
	assert.Contains(t, config, `"snat_automap":true`)
}

func TestGetParamsConfigMap_VirtualServer(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"virtual_server": []interface{}{
			map[string]interface{}{"ip": "10.0.0.1", "port": 8080},
		},
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, "10.0.0.1")
	assert.Contains(t, config, "8080")
}

func TestGetParamsConfigMap_SnatPool(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":             "mytenant",
		"application":        "myapp",
		"existing_snat_pool": "/Common/my-snatpool",
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-snatpool")
}

func TestGetParamsConfigMap_SnatPoolAddresses(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":            "mytenant",
		"application":       "myapp",
		"snat_pool_address": []interface{}{"10.1.1.1", "10.1.1.2"},
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, "10.1.1.1")
	assert.Contains(t, config, `"make_snatpool":true`)
}

func TestGetParamsConfigMap_PersistenceProfile(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":               "mytenant",
		"application":          "myapp",
		"persistence_profile":  "/Common/my-persist",
		"fallback_persistence": "source-address",
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-persist")
	assert.Contains(t, config, "source-address")
}

func TestGetParamsConfigMap_PersistenceType(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":           "mytenant",
		"application":      "myapp",
		"persistence_type": "destination-address",
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, "destination-address")
}

func TestGetParamsConfigMap_ExistingPool(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":        "mytenant",
		"application":   "myapp",
		"existing_pool": "/Common/my-pool",
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-pool")
	assert.Contains(t, config, `"make_pool":false`)
}

func TestGetParamsConfigMap_PoolMembers(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"pool_members": []interface{}{
			map[string]interface{}{
				"addresses":        []interface{}{"192.0.2.1", "192.0.2.2"},
				"port":             8080,
				"connection_limit": 100,
				"priority_group":   1,
				"share_nodes":      true,
			},
		},
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, "192.0.2.1")
	assert.Contains(t, config, `"make_pool":true`)
}

func TestGetParamsConfigMap_MonitorAndLoadBalancing(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":              "mytenant",
		"application":         "myapp",
		"existing_monitor":    "/Common/tcp",
		"load_balancing_mode": "round-robin",
		"slow_ramp_time":      30,
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, "round-robin")
	assert.Contains(t, config, `"monitor_name":"/Common/tcp"`)
}

func TestGetParamsConfigMap_GeneratedMonitor(t *testing.T) {
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"monitor": []interface{}{
			map[string]interface{}{"interval": 10},
		},
	}, "")

	config, err := getParamsConfigMap(d)
	require.NoError(t, err)
	assert.Contains(t, config, `"monitor_interval":10`)
	assert.Contains(t, config, `"make_monitor":true`)
}

// ---------------------------------------------------------------------
// resourceBigipFastTcpAppCreate / Read / Update / Delete
// ---------------------------------------------------------------------

func fastTcpMockServer(t *testing.T, tenant, app string) *httptest.Server {
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
			_, _ = fmt.Fprintf(w, `{"constants":{"fast":{"view":{"tenant_name":"%s","app_name":"%s","virtual_address":"10.0.0.1","virtual_port":8080}}}}`, tenant, app)
		case http.MethodPatch:
			_, _ = fmt.Fprint(w, `{"message":[{"id":"task-1"}]}`)
		case http.MethodDelete:
			_, _ = fmt.Fprint(w, `{"id":"task-1"}`)
		}
	})

	return httptest.NewServer(mux)
}

func TestUnitFastTcpCreateReadUpdateDelete(t *testing.T) {
	tenant := "mytenant"
	app := "myapp"

	server := fastTcpMockServer(t, tenant, app)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call

	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      tenant,
		"application": app,
		"virtual_server": []interface{}{
			map[string]interface{}{"ip": "10.0.0.1", "port": 8080},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipFastTcpAppCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, app, d.Id())
	assert.Equal(t, tenant, d.Get("tenant"))

	diags = resourceBigipFastTcpAppRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Contains(t, d.Get("fast_tcp_json").(string), "10.0.0.1")

	diags = resourceBigipFastTcpAppUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipFastTcpAppDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
}

func TestUnitFastTcpCreate_PostError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"post failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "")

	diags := resourceBigipFastTcpAppCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastTcpRead_AppNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/missing-app", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"Client Error: Could not find application mytenant/missing-app"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "missing-app",
	}, "missing-app")

	diags := resourceBigipFastTcpAppRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected the app-not-found branch to clear state without an error")
	assert.Equal(t, "", d.Id())
}

func TestUnitFastTcpRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastTcpAppRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastTcpUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastTcpAppUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastTcpDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastTcpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastTcpAppDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
