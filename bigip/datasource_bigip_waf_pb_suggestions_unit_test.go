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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// wafPbResourceData builds a *schema.ResourceData for the bigip_waf_pb
// data source using schema.TestResourceDataRaw.
func wafPbResourceData(t *testing.T, policyName, partition string, minimumLearningScore int) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipWafPb(), map[string]interface{}{
		"policy_name":            policyName,
		"partition":              partition,
		"minimum_learning_score": minimumLearningScore,
	})
}

// TestDataSourceBigipWafPbSchema exercises the schema definition of the
// bigip_waf_pb data source without needing any BIG-IP connection.
func TestDataSourceBigipWafPbSchema(t *testing.T) {
	r := dataSourceBigipWafPb()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	requiredFields := []string{"policy_name", "partition"}
	for _, field := range requiredFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
		if s.Type != schema.TypeString {
			t.Errorf("Expected field '%s' to be TypeString, got %v", field, s.Type)
		}
	}

	scoreSchema, ok := r.Schema["minimum_learning_score"]
	if !ok {
		t.Fatal("Expected field 'minimum_learning_score' to exist in schema")
	}
	if !scoreSchema.Required {
		t.Error("Expected field 'minimum_learning_score' to be required")
	}
	if scoreSchema.Type != schema.TypeInt {
		t.Errorf("Expected field 'minimum_learning_score' to be TypeInt, got %v", scoreSchema.Type)
	}
	if scoreSchema.ValidateFunc == nil {
		t.Error("Expected field 'minimum_learning_score' to have a ValidateFunc")
	} else {
		_, errs := scoreSchema.ValidateFunc(0, "minimum_learning_score")
		assert.NotEmpty(t, errs, "expected validation error for a value below the allowed range (1-100)")
		_, errs = scoreSchema.ValidateFunc(101, "minimum_learning_score")
		assert.NotEmpty(t, errs, "expected validation error for a value above the allowed range (1-100)")
		_, errs = scoreSchema.ValidateFunc(50, "minimum_learning_score")
		assert.Empty(t, errs, "expected no validation error for a value within the allowed range (1-100)")
	}

	policyIDSchema, ok := r.Schema["policy_id"]
	if !ok {
		t.Fatal("Expected field 'policy_id' to exist in schema")
	}
	if !policyIDSchema.Optional {
		t.Error("Expected field 'policy_id' to be Optional")
	}
	if !policyIDSchema.Computed {
		t.Error("Expected field 'policy_id' to be Computed")
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

// TestDataSourceBigipWafPbReadSuccess covers the Read happy path:
// GetWafPolicyId, PostPbExport, and GetWafPbExportResult all succeed with
// an immediate COMPLETED status (no polling loop iterations), so
// policy_id/json get populated and the resource ID is set to the policy
// name.
func TestDataSourceBigipWafPbReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"test-policy","partition":"Common","id":"test-policy-id"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"pb-task-id","status":"COMPLETED"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions/pb-task-id", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"pb-task-id","status":"COMPLETED","result":{"suggestions":[]}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafPbResourceData(t, "test-policy", "Common", 50)

	diags := dataSourceBigipWafPbRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "test-policy", d.Id())
	assert.Equal(t, "test-policy-id", d.Get("policy_id"))
	assert.Contains(t, d.Get("json"), "suggestions")
}

// TestDataSourceBigipWafPbReadPolicyIdError covers the branch where
// GetWafPolicyId fails to resolve the policy (e.g. no matching policy on
// that partition).
func TestDataSourceBigipWafPbReadPolicyIdError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafPbResourceData(t, "missing-policy", "Common", 50)

	diags := dataSourceBigipWafPbRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the policy ID cannot be resolved")
	assert.Contains(t, diags[0].Summary, "error retrieving policy")
}

// TestDataSourceBigipWafPbReadExportError covers the branch where
// PostPbExport itself fails (e.g. a 500 from the export-suggestions
// endpoint).
func TestDataSourceBigipWafPbReadExportError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"test-policy","partition":"Common","id":"test-policy-id"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "export failed", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafPbResourceData(t, "test-policy", "Common", 50)

	diags := dataSourceBigipWafPbRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when PostPbExport fails")
	assert.Contains(t, diags[0].Summary, "error exporting pb suggestions")
}

// TestDataSourceBigipWafPbReadTaskFailure covers the branch where the
// export task itself completes with a FAILURE status (rather than an HTTP
// transport error), which Read must surface as an "export task failed"
// error.
func TestDataSourceBigipWafPbReadTaskFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"test-policy","partition":"Common","id":"test-policy-id"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"pb-task-id","status":"FAILURE"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions/pb-task-id", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"pb-task-id","status":"FAILURE"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafPbResourceData(t, "test-policy", "Common", 50)

	diags := dataSourceBigipWafPbRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the export task status is FAILURE")
	assert.Contains(t, diags[0].Summary, "export task failed")
}

// TestDataSourceBigipWafPbReadPolling covers the polling loop: the export
// task is still PENDING on the first status check (both the one right
// after PostPbExport and the first one inside the loop), and only reports
// COMPLETED on the second in-loop check. This exercises the loop body
// without ever reaching time.Sleep(3*time.Second): the loop breaks as soon
// as the newly fetched status is COMPLETED/FAILURE, before the sleep at
// the bottom of the loop runs.
func TestDataSourceBigipWafPbReadPolling(t *testing.T) {
	var callCount int
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"test-policy","partition":"Common","id":"test-policy-id"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"pb-task-id","status":"PENDING"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions/pb-task-id", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount < 2 {
			_, _ = fmt.Fprint(w, `{"id":"pb-task-id","status":"PENDING"}`)
			return
		}
		_, _ = fmt.Fprint(w, `{"id":"pb-task-id","status":"COMPLETED","result":{"suggestions":[]}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafPbResourceData(t, "test-policy", "Common", 50)

	diags := dataSourceBigipWafPbRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "test-policy", d.Id())
	assert.GreaterOrEqual(t, callCount, 2, "expected the polling loop to check status more than once")
}

// TestDataSourceBigipWafPbReadInitialStatusError covers the branch where
// the very first GetWafPbExportResult call (immediately after
// PostPbExport succeeds) fails with a transport/API error.
func TestDataSourceBigipWafPbReadInitialStatusError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/asm/policies/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[{"name":"test-policy","partition":"Common","id":"test-policy-id"}]}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"pb-task-id","status":"PENDING"}`)
	})
	mux.HandleFunc("/mgmt/tm/asm/tasks/export-suggestions/pb-task-id", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := wafPbResourceData(t, "test-policy", "Common", 50)

	diags := dataSourceBigipWafPbRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the initial GetWafPbExportResult call fails")
}
