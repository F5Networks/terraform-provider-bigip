/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
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
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipNetSelfIPSchema(t *testing.T) {
	r := resourceBigipNetSelfIP()

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

	for _, field := range []string{"name", "ip", "vlan"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}
	if r.Schema["traffic_group"].Default != "traffic-group-local-only" {
		t.Errorf("Expected 'traffic_group' to default to 'traffic-group-local-only', got %v", r.Schema["traffic_group"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitNetSelfIPCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-selfip"
	mangled := "/mgmt/tm/net/self/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","vlan":"/Common/vlan1","address":"10.0.0.1/24","trafficGroup":"/Common/traffic-group-local-only","allowService":["default"]}`, name, name)
		case http.MethodPut:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetSelfIP()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"ip":   "10.0.0.1/24",
		"vlan": "/Common/vlan1",
	}, "")

	ctx := context.Background()

	diags := resourceBigipNetSelfIPCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipNetSelfIPRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "/Common/vlan1", d.Get("vlan").(string))
	require.Equal(t, "traffic-group-local-only", d.Get("traffic_group").(string))

	diags = resourceBigipNetSelfIPUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipNetSelfIPDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitNetSelfIPRead_AllowServiceNil(t *testing.T) {
	name := "/Common/test-selfip"
	mangled := "/mgmt/tm/net/self/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","fullPath":"%s","vlan":"/Common/vlan1","address":"10.0.0.1/24","trafficGroup":"/Common/traffic-group-local-only"}`, name, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetSelfIP()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"ip":   "10.0.0.1/24",
		"vlan": "/Common/vlan1",
	}, name)

	diags := resourceBigipNetSelfIPRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	portLockdown := d.Get("port_lockdown").([]interface{})
	require.Len(t, portLockdown, 1)
	require.Equal(t, "none", portLockdown[0].(string))
}

func TestUnitNetSelfIPCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetSelfIP()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/test-selfip",
		"ip":   "10.0.0.1/24",
		"vlan": "/Common/vlan1",
	}, "")

	diags := resourceBigipNetSelfIPCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetSelfIPRead_Error(t *testing.T) {
	name := "/Common/test-selfip"
	mangled := "/mgmt/tm/net/self/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetSelfIP()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"ip":   "10.0.0.1/24",
		"vlan": "/Common/vlan1",
	}, name)

	diags := resourceBigipNetSelfIPRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetSelfIPUpdate_Error(t *testing.T) {
	name := "/Common/test-selfip"
	mangled := "/mgmt/tm/net/self/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetSelfIP()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"ip":   "10.0.0.1/24",
		"vlan": "/Common/vlan1",
	}, name)

	diags := resourceBigipNetSelfIPUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetSelfIPDelete_Error(t *testing.T) {
	name := "/Common/test-selfip"
	mangled := "/mgmt/tm/net/self/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetSelfIP()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"ip":   "10.0.0.1/24",
		"vlan": "/Common/vlan1",
	}, name)

	diags := resourceBigipNetSelfIPDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getNetSelfIPConfig
// ---------------------------------------------------------------------

func TestGetNetSelfIPConfig_PortLockdownVariants(t *testing.T) {
	r := resourceBigipNetSelfIP()

	t.Run("all", func(t *testing.T) {
		d := NewTestResourceData(t, r, map[string]interface{}{
			"name":          "/Common/test-selfip",
			"ip":            "10.0.0.1/24",
			"vlan":          "/Common/vlan1",
			"port_lockdown": []interface{}{"all"},
		}, "")
		cfg := getNetSelfIPConfig(d, &bigip.SelfIP{})
		require.Equal(t, "all", cfg.AllowService)
	})

	t.Run("none", func(t *testing.T) {
		d := NewTestResourceData(t, r, map[string]interface{}{
			"name":          "/Common/test-selfip",
			"ip":            "10.0.0.1/24",
			"vlan":          "/Common/vlan1",
			"port_lockdown": []interface{}{"none"},
		}, "")
		cfg := getNetSelfIPConfig(d, &bigip.SelfIP{})
		require.Nil(t, cfg.AllowService)
	})

	t.Run("explicit list", func(t *testing.T) {
		d := NewTestResourceData(t, r, map[string]interface{}{
			"name":          "/Common/test-selfip",
			"ip":            "10.0.0.1/24",
			"vlan":          "/Common/vlan1",
			"port_lockdown": []interface{}{"tcp:80", "tcp:443"},
		}, "")
		cfg := getNetSelfIPConfig(d, &bigip.SelfIP{})
		require.NotNil(t, cfg.AllowService)
	})

	t.Run("empty", func(t *testing.T) {
		d := NewTestResourceData(t, r, map[string]interface{}{
			"name": "/Common/test-selfip",
			"ip":   "10.0.0.1/24",
			"vlan": "/Common/vlan1",
		}, "")
		cfg := getNetSelfIPConfig(d, &bigip.SelfIP{})
		require.Nil(t, cfg.AllowService)
	})
}
