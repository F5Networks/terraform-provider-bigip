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

func TestResourceBigipFastUdpAppSchema(t *testing.T) {
	r := resourceBigipFastUdpApp()

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
}

// ---------------------------------------------------------------------
// getParamsConfigMapUdp (direct function tests, no HTTP server needed)
// ---------------------------------------------------------------------

func TestGetParamsConfigMapUdp_Minimal(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, "mytenant")
	assert.Contains(t, config, "myapp")
}

func TestGetParamsConfigMapUdp_Fastl4WithExistingProfile(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":           "mytenant",
		"application":      "myapp",
		"enable_fastl4":    true,
		"existing_profile": "/Common/fastL4",
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, `"fastl4":true`)
	assert.Contains(t, config, "fastL4")
}

func TestGetParamsConfigMapUdp_Fastl4MakeProfile(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":        "mytenant",
		"application":   "myapp",
		"enable_fastl4": true,
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, `"make_fastl4_profile":true`)
}

func TestGetParamsConfigMapUdp_ExistingUdpProfile(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":           "mytenant",
		"application":      "myapp",
		"existing_profile": "/Common/udp",
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, `"udp_profile_name":"/Common/udp"`)
}

func TestGetParamsConfigMapUdp_PersistenceFastl4(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":               "mytenant",
		"application":          "myapp",
		"enable_fastl4":        true,
		"existing_profile":     "/Common/fastL4",
		"persistence_profile":  "/Common/my-persist",
		"fallback_persistence": "source-address",
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, `"fastl4_persistence_profile":"/Common/my-persist"`)
	assert.Contains(t, config, "source-address")
}

func TestGetParamsConfigMapUdp_PersistenceUdp(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":           "mytenant",
		"application":      "myapp",
		"persistence_type": "destination-address",
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, `"persistence_type":"destination-address"`)
}

func TestGetParamsConfigMapUdp_SnatPoolAddresses(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":            "mytenant",
		"application":       "myapp",
		"snat_pool_address": []interface{}{"10.1.1.1"},
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, "10.1.1.1")
	assert.Contains(t, config, `"make_snatpool":true`)
}

func TestGetParamsConfigMapUdp_PoolMembers(t *testing.T) {
	r := resourceBigipFastUdpApp()
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

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, "192.0.2.1")
}

func TestGetParamsConfigMapUdp_MonitorAndVlans(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
		"monitor": []interface{}{
			map[string]interface{}{
				"interval":          10,
				"send_string":       "PING",
				"expected_response": "PONG",
			},
		},
		"irules":                []interface{}{"/Common/my-irule"},
		"vlans_allowed":         []interface{}{"/Common/external"},
		"security_log_profiles": []interface{}{"/Common/log-all"},
		"load_balancing_mode":   "round-robin",
		"slow_ramp_time":        30,
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, "PING")
	assert.Contains(t, config, "PONG")
	assert.Contains(t, config, "my-irule")
	assert.Contains(t, config, "external")
	assert.Contains(t, config, `"vlans_allow":true`)
	assert.Contains(t, config, "log-all")
	assert.Contains(t, config, "round-robin")
}

func TestGetParamsConfigMapUdp_VlansRejected(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":         "mytenant",
		"application":    "myapp",
		"vlans_rejected": []interface{}{"/Common/internal"},
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, "internal")
	assert.Contains(t, config, `"vlans_allow":false`)
}

func TestGetParamsConfigMapUdp_ExistingMonitorAndPool(t *testing.T) {
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":           "mytenant",
		"application":      "myapp",
		"existing_pool":    "/Common/my-pool",
		"existing_monitor": "/Common/udp",
	}, "")

	config, err := getParamsConfigMapUdp(d)
	require.NoError(t, err)
	assert.Contains(t, config, "my-pool")
	assert.Contains(t, config, `"monitor_name":"/Common/udp"`)
}

// ---------------------------------------------------------------------
// resourceBigipFastUdpAppCreate / Read / Update / Delete
// ---------------------------------------------------------------------

func fastUdpMockServer(t *testing.T, tenant, app string) *httptest.Server {
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

func TestUnitFastUdpCreateReadUpdateDelete(t *testing.T) {
	tenant := "mytenant"
	app := "myapp"

	server := fastUdpMockServer(t, tenant, app)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call

	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      tenant,
		"application": app,
		"virtual_server": []interface{}{
			map[string]interface{}{"ip": "10.0.0.1", "port": 8080},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipFastUdpAppCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, app, d.Id())
	assert.Equal(t, tenant, d.Get("tenant"))

	diags = resourceBigipFastUdpAppRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Contains(t, d.Get("fast_udp_json").(string), "10.0.0.1")

	diags = resourceBigipFastUdpAppUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipFastUdpAppDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitFastUdpCreate_PostError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"post failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "")

	diags := resourceBigipFastUdpAppCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastUdpRead_UnexpectedEndOfJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastUdpAppRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected nil diags for the unexpected-end-of-JSON branch")
	assert.Equal(t, "", d.Id())
}

func TestUnitFastUdpRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastUdpAppRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastUdpUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastUdpAppUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastUdpDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastUdpApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":      "mytenant",
		"application": "myapp",
	}, "myapp")

	diags := resourceBigipFastUdpAppDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
