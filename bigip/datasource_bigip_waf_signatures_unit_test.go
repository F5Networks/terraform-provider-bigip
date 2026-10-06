/*
Copyright 2022 F5 Networks Inc.
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// wafSignaturesResourceData builds a *schema.ResourceData for the
// bigip_waf_signatures data source using schema.TestResourceDataRaw.
func wafSignaturesResourceData(t *testing.T, signatureID int) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipWafSignatures(), map[string]interface{}{
		"signature_id": signatureID,
	})
}

// TestDataSourceBigipWafSignaturesSchema exercises the schema definition of
// the bigip_waf_signatures data source without needing any BIG-IP
// connection.
func TestDataSourceBigipWafSignaturesSchema(t *testing.T) {
	r := dataSourceBigipWafSignatures()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	sigIDSchema, ok := r.Schema["signature_id"]
	if !ok {
		t.Fatal("Expected field 'signature_id' to exist in schema")
	}
	if !sigIDSchema.Required {
		t.Error("Expected field 'signature_id' to be required")
	}
	if sigIDSchema.Type != schema.TypeInt {
		t.Errorf("Expected field 'signature_id' to be TypeInt, got %v", sigIDSchema.Type)
	}

	computedOptionalStringFields := []string{"name", "system_signature_id", "tag", "type", "accuracy", "risk"}
	for _, field := range computedOptionalStringFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Optional {
			t.Errorf("Expected field '%s' to be Optional", field)
		}
		if !s.Computed {
			t.Errorf("Expected field '%s' to be Computed", field)
		}
		if s.Type != schema.TypeString {
			t.Errorf("Expected field '%s' to be TypeString, got %v", field, s.Type)
		}
	}

	boolFields := []string{"perform_staging", "enabled"}
	for _, field := range boolFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Optional {
			t.Errorf("Expected field '%s' to be Optional", field)
		}
		if s.Type != schema.TypeBool {
			t.Errorf("Expected field '%s' to be TypeBool, got %v", field, s.Type)
		}
	}

	descriptionSchema, ok := r.Schema["description"]
	if !ok {
		t.Fatal("Expected field 'description' to exist in schema")
	}
	if !descriptionSchema.Optional {
		t.Error("Expected field 'description' to be Optional")
	}
	if descriptionSchema.Computed {
		t.Error("Expected field 'description' to not be Computed")
	}

	jsonSchema, ok := r.Schema["json"]
	if !ok {
		t.Fatal("Expected field 'json' to exist in schema")
	}
	if !jsonSchema.Computed {
		t.Error("Expected field 'json' to be Computed")
	}
	if jsonSchema.Type != schema.TypeString {
		t.Errorf("Expected field 'json' to be TypeString, got %v", jsonSchema.Type)
	}
}

// TestDataSourceBigipWafSignatureReadSuccess covers the Read happy path:
// asm is provisioned and GetWafSignature returns a single matching
// signature, so all fields are populated, the json field is built from the
// resource data, and the resource ID is set to the signature ID.
func TestDataSourceBigipWafSignatureReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"asm","level":"nominal"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/signatures/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"test-signature","id":"sig-resource-id","description":"a test signature","signatureId":1234,"signatureType":"request","accuracy":"high","risk":"high"}]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafSignaturesResourceData(t, 1234)

	diags := dataSourceBigipWafSignatureRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "1234", d.Id())
	assert.Equal(t, "test-signature", d.Get("name"))
	assert.Equal(t, "sig-resource-id", d.Get("system_signature_id"))
	assert.Equal(t, "request", d.Get("type"))
	assert.Equal(t, "high", d.Get("accuracy"))
	assert.Equal(t, "high", d.Get("risk"))
	assert.Contains(t, d.Get("json"), "1234")
}

// TestDataSourceBigipWafSignatureReadNotProvisioned covers the branch where
// asm is not provisioned on the device: Read must fail fast with an error
// before attempting to look up the signature.
func TestDataSourceBigipWafSignatureReadNotProvisioned(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"asm","level":"none"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafSignaturesResourceData(t, 1234)

	diags := dataSourceBigipWafSignatureRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when asm is not provisioned")
	assert.Contains(t, diags[0].Summary, "asm Module is not provisioned")
}

// TestDataSourceBigipWafSignatureReadProvisionsError covers the branch
// where the Provisions call itself fails (e.g. the BIG-IP is unreachable).
func TestDataSourceBigipWafSignatureReadProvisionsError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafSignaturesResourceData(t, 1234)

	diags := dataSourceBigipWafSignatureRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when Provisions fails")
}

// TestDataSourceBigipWafSignatureReadGetSignatureError covers the branch
// where GetWafSignature itself fails after asm is confirmed provisioned.
func TestDataSourceBigipWafSignatureReadGetSignatureError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"asm","level":"nominal"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/signatures/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafSignaturesResourceData(t, 1234)

	diags := dataSourceBigipWafSignatureRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when GetWafSignature fails")
	assert.Contains(t, diags[0].Summary, "error retrieving signature")
}

// TestDataSourceBigipWafSignatureReadNotFound covers the branch where the
// filter query returns zero matching signatures: Read should skip the
// d.Set calls for the API-populated fields (leaving them at their zero
// values) but still succeed and build the json payload from
// signature_id/enabled/perform_staging.
func TestDataSourceBigipWafSignatureReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/provision/asm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"asm","level":"nominal"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/signatures/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafSignaturesResourceData(t, 9999)

	diags := dataSourceBigipWafSignatureRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "9999", d.Id())
	assert.Empty(t, d.Get("name"), "expected name to remain unset when the signature is not found")
	assert.Contains(t, d.Get("json"), "9999")
}
