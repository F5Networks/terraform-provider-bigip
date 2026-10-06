/*
Copyright 2021 F5 Networks Inc.
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

func TestResourceBigipIpsecPolicySchema(t *testing.T) {
	r := resourceBigipIpsecPolicy()

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
	if !nameSchema.ForceNew {
		t.Error("Expected field 'name' to be ForceNew")
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitIpsecPolicyCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-ipsec-policy"
	mangled := "/mgmt/tm/net/ipsec/ipsec-policy/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/ipsec/ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","protocol":"esp","mode":"tunnel","ikePhase2EncryptAlgorithm":"aes256","ikePhase2AuthAlgorithm":"sha256","ikePhase2Lifetime":120,"ikePhase2LifetimeKilobytes":0,"ikePhase2PerfectForwardSecrecy":"modp2048","ipcomp":"none"}`, name)
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
	client.Teem = true // skip real telemetry network call
	r := resourceBigipIpsecPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                    name,
		"protocol":                "esp",
		"mode":                    "tunnel",
		"encrypt_algorithm":       "aes256",
		"auth_algorithm":          "sha256",
		"lifetime":                120,
		"perfect_forward_secrecy": "modp2048",
		"ipcomp":                  "none",
	}, "")

	ctx := context.Background()

	diags := resourceBigipIpsecPolicyCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipIpsecPolicyRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	diags = resourceBigipIpsecPolicyUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipIpsecPolicyDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitIpsecPolicyCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/ipsec/ipsec-policy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipIpsecPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "/Common/test-ipsec-policy",
	}, "")

	diags := resourceBigipIpsecPolicyCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitIpsecPolicyRead_Error(t *testing.T) {
	name := "/Common/test-ipsec-policy"
	mangled := "/mgmt/tm/net/ipsec/ipsec-policy/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipIpsecPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipIpsecPolicyRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitIpsecPolicyUpdate_Error(t *testing.T) {
	name := "/Common/test-ipsec-policy"
	mangled := "/mgmt/tm/net/ipsec/ipsec-policy/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipIpsecPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipIpsecPolicyUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitIpsecPolicyDelete_Error(t *testing.T) {
	name := "/Common/test-ipsec-policy"
	mangled := "/mgmt/tm/net/ipsec/ipsec-policy/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipIpsecPolicy()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipIpsecPolicyDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
