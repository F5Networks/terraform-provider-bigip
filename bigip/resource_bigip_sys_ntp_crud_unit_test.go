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

func TestResourceBigipSysNtpCRUDSchema(t *testing.T) {
	r := resourceBigipSysNtp()

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

	for _, field := range []string{"description", "servers"} {
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

func TestUnitSysNtpCreateReadUpdateDelete(t *testing.T) {
	description := "test ntp"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/ntp", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			_, _ = fmt.Fprint(w, `{"description":"test ntp"}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, `{"description":"test ntp","servers":["10.10.10.10"],"timezone":"America/Los_Angeles"}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysNtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description": description,
		"servers":     []interface{}{"10.10.10.10"},
		"timezone":    "America/Los_Angeles",
	}, "")

	ctx := context.Background()

	diags := resourceBigipSysNtpCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, description, d.Id())

	diags = resourceBigipSysNtpRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "America/Los_Angeles", d.Get("timezone").(string))

	diags = resourceBigipSysNtpUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSysNtpDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitSysNtpCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/ntp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysNtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description": "test ntp",
		"servers":     []interface{}{"10.10.10.10"},
	}, "")

	diags := resourceBigipSysNtpCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysNtpRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/ntp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysNtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description": "test ntp",
		"servers":     []interface{}{"10.10.10.10"},
	}, "test ntp")

	diags := resourceBigipSysNtpRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysNtpUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/ntp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysNtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description": "test ntp",
		"servers":     []interface{}{"10.10.10.10"},
	}, "test ntp")

	diags := resourceBigipSysNtpUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysNtpDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/ntp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysNtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description": "test ntp",
		"servers":     []interface{}{"10.10.10.10"},
	}, "test ntp")

	diags := resourceBigipSysNtpDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getSysNTPConfig
// ---------------------------------------------------------------------

func TestGetSysNTPConfig(t *testing.T) {
	r := resourceBigipSysNtp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"description": "test ntp",
		"servers":     []interface{}{"10.10.10.10", "10.10.10.11"},
		"timezone":    "UTC",
	}, "")

	cfg := getSysNTPConfig(d, &bigip.NTP{})
	require.Equal(t, []string{"10.10.10.10", "10.10.10.11"}, cfg.Servers)
	require.Equal(t, "UTC", cfg.Timezone)
}
