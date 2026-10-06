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

func TestResourceBigipSysProvisionSchema(t *testing.T) {
	r := resourceBigipSysProvision()

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

	nameSchema, ok := r.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist in schema")
	}
	if !nameSchema.Required {
		t.Error("Expected field 'name' to be required")
	}
	if r.Schema["level"].Default != "nominal" {
		t.Errorf("Expected 'level' to default to 'nominal', got %v", r.Schema["level"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitSysProvisionCreateReadUpdateDelete(t *testing.T) {
	// ProvisionModule/Provisions only issue an HTTP call for a handful of
	// module names (asm, afm, gtm, apm, avr, ilx); "asm" is used here to
	// exercise the real request path.
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			_, _ = fmt.Fprint(w, `{"name":"asm"}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, `{"name":"asm","fullPath":"asm","cpuRatio":0,"diskRatio":0,"level":"nominal","memoryRatio":0}`)
		}
	})
	// Create/Update call waitForProvisionReady after ProvisionModule
	// succeeds, which (in addition to polling client.Provisions, already
	// handled by the handler above) also probes a brand-new
	// /mgmt/shared/authn/login via freshLoginSucceeds -- without a
	// handler for it here, that probe 404s on every attempt and
	// waitForProvisionReady retries for its full 10-minute timeout before
	// giving up, turning this fast unit test into a multi-minute hang.
	mux.HandleFunc("/mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"token":{"token":"test-token"}}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysProvision()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":  "asm",
		"level": "nominal",
	}, "")

	ctx := context.Background()

	diags := resourceBigipSysProvisionCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, "asm", d.Id())

	diags = resourceBigipSysProvisionRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "nominal", d.Get("level").(string))

	diags = resourceBigipSysProvisionUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSysProvisionDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
}

func TestUnitSysProvisionCreate_UnsupportedModuleNoop(t *testing.T) {
	// A module name outside the special-cased list (e.g. "ltm") results in
	// ProvisionModule/Provisions performing no HTTP call at all and
	// returning nil error / a zero-value Provision -- Create should still
	// succeed.
	client := NewUnitTestClient("http://127.0.0.1:0")
	r := resourceBigipSysProvision()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":  "ltm",
		"level": "nominal",
	}, "")

	diags := resourceBigipSysProvisionCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, "ltm", d.Id())
}

func TestUnitSysProvisionCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysProvision()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "asm",
	}, "")

	diags := resourceBigipSysProvisionCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysProvisionRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysProvision()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "asm",
	}, "asm")

	diags := resourceBigipSysProvisionRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysProvisionUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysProvision()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "asm",
	}, "asm")

	diags := resourceBigipSysProvisionUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getsysProvisionConfig
// ---------------------------------------------------------------------

func TestGetsysProvisionConfig(t *testing.T) {
	r := resourceBigipSysProvision()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":         "asm",
		"full_path":    "asm",
		"cpu_ratio":    5,
		"disk_ratio":   6,
		"level":        "minimum",
		"memory_ratio": 7,
	}, "")

	cfg := getsysProvisionConfig(d, &bigip.Provision{})
	require.Equal(t, "asm", cfg.FullPath)
	require.Equal(t, 5, cfg.CpuRatio)
	require.Equal(t, 6, cfg.DiskRatio)
	require.Equal(t, "minimum", cfg.Level)
	require.Equal(t, 7, cfg.MemoryRatio)
}
