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

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipSysIappSchema(t *testing.T) {
	r := resourceBigipSysIapp()

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

	for _, field := range []string{"jsonfile", "name"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}
	if r.Schema["partition"].Default != "Common" {
		t.Errorf("Expected 'partition' to default to 'Common', got %v", r.Schema["partition"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitSysIappCreateReadUpdateDelete(t *testing.T) {
	name := "testapp"
	iappPath := "/mgmt/tm/sys/application/service/~Common~" + name + ".app~" + name

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/application/service", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(iappPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","deviceGroup":"none","inheritedDevicegroup":"true","inheritedTrafficGroup":"true","strictUpdates":"enabled","templateModified":"no","trafficGroup":"/Common/traffic-group-1"}`, name)
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
	r := resourceBigipSysIapp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"jsonfile":  `{"name":"testapp"}`,
		"partition": "Common",
	}, "")

	ctx := context.Background()

	diags := resourceBigipSysIappCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipSysIappRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "enabled", d.Get("strict_updates").(string))

	diags = resourceBigipSysIappUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSysIappDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitSysIappCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/application/service", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIapp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testapp",
		"jsonfile":  `{"name":"testapp"}`,
		"partition": "Common",
	}, "")

	diags := resourceBigipSysIappCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysIappRead_Error(t *testing.T) {
	name := "testapp"
	iappPath := "/mgmt/tm/sys/application/service/~Common~" + name + ".app~" + name

	mux := http.NewServeMux()
	mux.HandleFunc(iappPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIapp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"jsonfile":  `{"name":"testapp"}`,
		"partition": "Common",
	}, name)

	diags := resourceBigipSysIappRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysIappUpdate_Error(t *testing.T) {
	name := "testapp"
	iappPath := "/mgmt/tm/sys/application/service/~Common~" + name + ".app~" + name

	mux := http.NewServeMux()
	mux.HandleFunc(iappPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIapp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"jsonfile":  `{"name":"testapp"}`,
		"partition": "Common",
	}, name)

	diags := resourceBigipSysIappUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysIappDelete_Error(t *testing.T) {
	name := "testapp"
	iappPath := "/mgmt/tm/sys/application/service/~Common~" + name + ".app~" + name

	mux := http.NewServeMux()
	mux.HandleFunc(iappPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysIapp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      name,
		"jsonfile":  `{"name":"testapp"}`,
		"partition": "Common",
	}, name)

	diags := resourceBigipSysIappDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// dataToIapp
// ---------------------------------------------------------------------

func TestDataToIapp(t *testing.T) {
	r := resourceBigipSysIapp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                    "testapp",
		"jsonfile":                `{"name":"testapp","description":"from-json"}`,
		"partition":               "Common",
		"execute_action":          "definition",
		"template_modified":       "yes",
		"strict_updates":          "disabled",
		"description":             "overridden",
		"inherited_devicegroup":   "false",
		"inherited_traffic_group": "false",
	}, "")

	p := dataToIapp(d)
	require.Equal(t, "definition", p.ExecuteAction)
	require.Equal(t, "Common", p.Partition)
	require.Equal(t, "yes", p.TemplateModified)
	require.Equal(t, "disabled", p.StrictUpdates)
	require.Equal(t, "overridden", p.Description)
	require.Equal(t, "false", p.InheritedDevicegroup)
	require.Equal(t, "false", p.InheritedTrafficGroup)
}

func TestDataToIapp_InvalidJSON(t *testing.T) {
	// Malformed jsonfile: json.Unmarshal fails and dataToIapp just logs the
	// error (via fmt.Println) rather than returning it, so the rest of the
	// GetOk-driven fields are still applied to a zero-value Iapp.
	r := resourceBigipSysIapp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testapp",
		"jsonfile":  `not-json`,
		"partition": "Common",
	}, "")

	p := dataToIapp(d)
	require.Equal(t, "Common", p.Partition)
}
