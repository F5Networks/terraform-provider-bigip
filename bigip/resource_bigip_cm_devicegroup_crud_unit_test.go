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

func TestResourceBigipCmDevicegroupSchema(t *testing.T) {
	r := resourceBigipCmDevicegroup()

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
	if r.Schema["type"].Default != "sync-only" {
		t.Errorf("Expected 'type' to default to 'sync-only', got %v", r.Schema["type"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitCmDevicegroupCreateReadUpdateDelete(t *testing.T) {
	name := "testdg"
	deviceName := "device1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device-group", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/cm/device-group/"+name+"/devices/"+deviceName, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, deviceName)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})
	mux.HandleFunc("/mgmt/tm/cm/device-group/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","autoSync":"disabled","type":"sync-only","fullLoadOnSync":"false","saveOnAutoSync":"false","networkFailover":"enabled","incrementalConfigSyncSizeMax":1024}`, name)
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
	r := resourceBigipCmDevicegroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"device": []interface{}{
			map[string]interface{}{
				"name": deviceName,
			},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipCmDevicegroupCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipCmDevicegroupRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "sync-only", d.Get("type").(string))

	diags = resourceBigipCmDevicegroupUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipCmDevicegroupDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitCmDevicegroupCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device-group", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevicegroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "testdg",
	}, "")

	diags := resourceBigipCmDevicegroupCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitCmDevicegroupRead_DevicesLookupErrorIsLoggedNotFailed(t *testing.T) {
	// DevicegroupsDevices errors inside the Read loop are only logged, not
	// propagated -- Read should still succeed as long as the subsequent
	// Devicegroups(name) call succeeds.
	name := "testdg"
	deviceName := "device1"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device-group/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","autoSync":"disabled","type":"sync-only"}`, name)
	})
	// No handler registered for the devices sub-path, so DevicegroupsDevices
	// gets a 404 and returns an error that Read merely logs.
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevicegroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"device": []interface{}{
			map[string]interface{}{
				"name": deviceName,
			},
		},
	}, name)

	diags := resourceBigipCmDevicegroupRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
}

func TestUnitCmDevicegroupRead_Error(t *testing.T) {
	name := "testdg"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device-group/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevicegroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipCmDevicegroupRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitCmDevicegroupUpdate_Error(t *testing.T) {
	name := "testdg"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device-group/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevicegroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipCmDevicegroupUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitCmDevicegroupDelete_DeviceRemovalError(t *testing.T) {
	name := "testdg"
	deviceName := "device1"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device-group/"+name+"/devices/"+deviceName, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"remove device failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevicegroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"device": []interface{}{
			map[string]interface{}{
				"name": deviceName,
			},
		},
	}, name)

	diags := resourceBigipCmDevicegroupDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitCmDevicegroupDelete_Error(t *testing.T) {
	name := "testdg"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device-group/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevicegroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipCmDevicegroupDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// dataToDevicegroup / DevicegroupToData
// ---------------------------------------------------------------------

func TestDataToDevicegroup(t *testing.T) {
	r := resourceBigipCmDevicegroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":               "testdg",
		"partition":          "Common",
		"auto_sync":          "enabled",
		"description":        "my dg",
		"type":               "sync-failover",
		"full_load_on_sync":  "true",
		"save_on_auto_sync":  "true",
		"network_failover":   "disabled",
		"incremental_config": 2048,
		"device": []interface{}{
			map[string]interface{}{"name": "device1"},
			map[string]interface{}{"name": "device2"},
		},
	}, "")

	p := dataToDevicegroup("testdg", d)
	require.Equal(t, "testdg", p.Name)
	require.Equal(t, "Common", p.Partition)
	require.Equal(t, "enabled", p.AutoSync)
	require.Equal(t, "my dg", p.Description)
	require.Equal(t, "sync-failover", p.Type)
	require.Equal(t, "true", p.FullLoadOnSync)
	require.Equal(t, "true", p.SaveOnAutoSync)
	require.Equal(t, "disabled", p.NetworkFailover)
	require.Equal(t, 2048, p.IncrementalConfigSyncSizeMax)
	require.Len(t, p.Deviceb, 2)
	require.Equal(t, "device1", p.Deviceb[0].Name)
	require.Equal(t, "device2", p.Deviceb[1].Name)
}

func TestDevicegroupToData(t *testing.T) {
	// Note: DevicegroupToData's per-index "device.N.name" / ".set_sync_leader"
	// d.Set calls are not valid schema.TestResourceDataRaw paths for list
	// elements (the SDK only supports setting a full list, not an indexed
	// element), so this test intentionally leaves Deviceb empty to exercise
	// the top-level field assignments without hitting that loop.
	r := resourceBigipCmDevicegroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "testdg",
	}, "")

	p := &bigip.Devicegroup{
		Name:                         "testdg",
		Partition:                    "Common",
		AutoSync:                     "enabled",
		Description:                  "my dg",
		Type:                         "sync-failover",
		FullLoadOnSync:               "true",
		SaveOnAutoSync:               "true",
		NetworkFailover:              "disabled",
		IncrementalConfigSyncSizeMax: 2048,
	}

	err := DevicegroupToData(p, d)
	require.NoError(t, err)
	require.Equal(t, "testdg", d.Get("name").(string))
	require.Equal(t, "enabled", d.Get("auto_sync").(string))
	require.Equal(t, 2048, d.Get("incremental_config").(int))
}
