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

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipCmDeviceCRUDSchema(t *testing.T) {
	r := resourceBigipCmDevice()

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

	for _, field := range []string{"configsync_ip", "name"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitCmDeviceCreateReadUpdateDelete(t *testing.T) {
	name := "testdevice"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodPut:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		}
	})
	mux.HandleFunc("/mgmt/tm/cm/device/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","configsyncIp":"10.0.0.1","mirrorIp":"10.0.0.2","mirrorSecondaryIp":"10.0.0.3"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevice()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"configsync_ip": "10.0.0.1",
	}, "")

	ctx := context.Background()

	diags := resourceBigipCmDeviceCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipCmDeviceRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "10.0.0.1", d.Get("configsync_ip").(string))

	diags = resourceBigipCmDeviceUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipCmDeviceDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitCmDeviceCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevice()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "testdevice",
		"configsync_ip": "10.0.0.1",
	}, "")

	diags := resourceBigipCmDeviceCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitCmDeviceRead_Error(t *testing.T) {
	name := "testdevice"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevice()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"configsync_ip": "10.0.0.1",
	}, name)

	diags := resourceBigipCmDeviceRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitCmDeviceUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevice()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "testdevice",
		"configsync_ip": "10.0.0.1",
	}, "testdevice")

	diags := resourceBigipCmDeviceUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitCmDeviceDelete_Error(t *testing.T) {
	name := "testdevice"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/cm/device/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCmDevice()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"configsync_ip": "10.0.0.1",
	}, name)

	diags := resourceBigipCmDeviceDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
