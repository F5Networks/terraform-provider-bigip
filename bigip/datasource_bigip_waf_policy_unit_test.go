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

// wafPolicyResourceData builds a *schema.ResourceData for the
// bigip_waf_policy data source using schema.TestResourceDataRaw.
func wafPolicyResourceData(t *testing.T, policyID string) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipWafPolicy(), map[string]interface{}{
		"policy_id": policyID,
	})
}

// TestDataSourceBigipWafPolicySchema exercises the schema definition of the
// bigip_waf_policy data source without needing any BIG-IP connection.
func TestDataSourceBigipWafPolicySchema(t *testing.T) {
	r := dataSourceBigipWafPolicy()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	policyIDSchema, ok := r.Schema["policy_id"]
	if !ok {
		t.Fatal("Expected field 'policy_id' to exist in schema")
	}
	if !policyIDSchema.Required {
		t.Error("Expected field 'policy_id' to be required")
	}
	if policyIDSchema.Type != schema.TypeString {
		t.Errorf("Expected field 'policy_id' to be TypeString, got %v", policyIDSchema.Type)
	}

	policyJSONSchema, ok := r.Schema["policy_json"]
	if !ok {
		t.Fatal("Expected field 'policy_json' to exist in schema")
	}
	if !policyJSONSchema.Optional {
		t.Error("Expected field 'policy_json' to be Optional")
	}
	if !policyJSONSchema.Computed {
		t.Error("Expected field 'policy_json' to be Computed")
	}
	if policyJSONSchema.Type != schema.TypeString {
		t.Errorf("Expected field 'policy_json' to be TypeString, got %v", policyJSONSchema.Type)
	}
}

// TestDataSourceBigipWafPolicyReadSuccess covers the Read happy path:
// GetWafPolicy and ExportPolicy both succeed, so policy_id/policy_json are
// populated and the resource ID is set to the policy's ID.
func TestDataSourceBigipWafPolicyReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/test-policy-id", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"test-policy-id","name":"test-policy","partition":"Common","fullPath":"/Common/test-policy"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-policy", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"export-task-id","status":"COMPLETED"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-policy/export-task-id", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"export-task-id","status":"COMPLETED","result":{"file":"{\"policy\":{\"id\":\"test-policy-id\",\"name\":\"test-policy\"}}"}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafPolicyResourceData(t, "test-policy-id")

	diags := dataSourceBigipWafPolicyRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "test-policy-id", d.Id())
	assert.Equal(t, "test-policy-id", d.Get("policy_id"))
	assert.Contains(t, d.Get("policy_json"), "test-policy")
}

// TestDataSourceBigipWafPolicyReadGetPolicyError covers the branch where
// GetWafPolicy itself returns an error (e.g. the BIG-IP is unreachable or
// the policy doesn't exist).
func TestDataSourceBigipWafPolicyReadGetPolicyError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/broken-policy-id", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafPolicyResourceData(t, "broken-policy-id")
	d.SetId("preexisting-id")

	diags := dataSourceBigipWafPolicyRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when GetWafPolicy fails")
	assert.Contains(t, diags[0].Summary, "error retrieving waf policy")
	assert.Empty(t, d.Id(), "resource ID should be cleared when the API call fails")
}

// TestDataSourceBigipWafPolicyReadExportError covers the branch where
// GetWafPolicy succeeds but ExportPolicy fails (e.g. the export task POST
// itself errors out).
func TestDataSourceBigipWafPolicyReadExportError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/test-policy-id", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"test-policy-id","name":"test-policy"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-policy", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "export failed", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafPolicyResourceData(t, "test-policy-id")
	d.SetId("preexisting-id")

	diags := dataSourceBigipWafPolicyRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when ExportPolicy fails")
	assert.Contains(t, diags[0].Summary, "error Exporting waf policy")
	assert.Empty(t, d.Id(), "resource ID should be cleared when the export call fails")
}
