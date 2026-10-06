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

func TestResourceBigipSysSnmpCRUDSchema(t *testing.T) {
	r := resourceBigipSysSnmp()

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
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitSysSnmpCreateReadUpdateDelete(t *testing.T) {
	sysContact := "admin@example.com"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/snmp", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch, http.MethodPut:
			_, _ = fmt.Fprintf(w, `{"sysContact":"%s"}`, sysContact)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"sysContact":"%s","sysLocation":"DataCenter1","allowedAddresses":["10.0.0.0/8"]}`, sysContact)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysSnmp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"sys_contact":      sysContact,
		"sys_location":     "DataCenter1",
		"allowedaddresses": []interface{}{"10.0.0.0/8"},
	}, "")

	ctx := context.Background()

	diags := resourceBigipSysSnmpCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, sysContact, d.Id())

	diags = resourceBigipSysSnmpRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "DataCenter1", d.Get("sys_location").(string))

	diags = resourceBigipSysSnmpUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSysSnmpDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
}

func TestUnitSysSnmpCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/snmp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysSnmp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"sys_contact": "admin@example.com",
	}, "")

	diags := resourceBigipSysSnmpCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysSnmpRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/snmp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysSnmp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"sys_contact": "admin@example.com",
	}, "admin@example.com")

	diags := resourceBigipSysSnmpRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitSysSnmpUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/snmp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipSysSnmp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"sys_contact": "admin@example.com",
	}, "admin@example.com")

	diags := resourceBigipSysSnmpUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}
