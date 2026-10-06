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

// iruleDatasourceResourceData builds a *schema.ResourceData for the
// bigip_ltm_irule data source using schema.TestResourceDataRaw.
func iruleDatasourceResourceData(t *testing.T, name, partition string) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipLtmIrule(), map[string]interface{}{
		"name":      name,
		"partition": partition,
	})
}

// TestDataSourceBigipLtmIruleSchema exercises the schema definition of the
// bigip_ltm_irule data source without needing any BIG-IP connection.
func TestDataSourceBigipLtmIruleSchema(t *testing.T) {
	r := dataSourceBigipLtmIrule()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	requiredFields := []string{"name", "partition"}
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
		if s.Description == "" {
			t.Errorf("Expected field '%s' to have a description", field)
		}
	}

	iruleSchema, ok := r.Schema["irule"]
	if !ok {
		t.Fatal("Expected field 'irule' to exist in schema")
	}
	if iruleSchema.Required {
		t.Error("Expected field 'irule' to be optional")
	}
	if iruleSchema.Type != schema.TypeString {
		t.Errorf("Expected field 'irule' to be TypeString, got %v", iruleSchema.Type)
	}
	if iruleSchema.StateFunc == nil {
		t.Error("Expected field 'irule' to have a StateFunc")
	} else {
		trimmed := iruleSchema.StateFunc("  some irule body  \n")
		assert.Equal(t, "some irule body", trimmed, "StateFunc should trim surrounding whitespace")
	}

	if r.ReadContext == nil {
		t.Error("Expected ReadContext to be defined")
	}
}

// TestDataSourceBigipLtmIruleReadSuccess covers the Read happy path: IRule()
// succeeds and returns a non-nil iRule, so the resource fields/ID are
// populated from the API response.
func TestDataSourceBigipLtmIruleReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~test-irule", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"test-irule","partition":"Common","fullPath":"/Common/test-irule","apiAnonymous":"when HTTP_REQUEST {\n  log local0. \"test\"\n}"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleDatasourceResourceData(t, "test-irule", "Common")

	diags := dataSourceBigipLtmIruleRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/test-irule", d.Id())
	assert.Equal(t, "/Common/test-irule", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))
	assert.Equal(t, "when HTTP_REQUEST {\n  log local0. \"test\"\n}", d.Get("irule"))
}

// TestDataSourceBigipLtmIruleReadError covers the branch where IRule() itself
// returns an error (e.g. the BIG-IP is unreachable or returns a
// non-JSON/error payload).
func TestDataSourceBigipLtmIruleReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~broken-irule", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleDatasourceResourceData(t, "broken-irule", "Common")

	diags := dataSourceBigipLtmIruleRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when the API call fails")
	assert.Contains(t, diags[0].Summary, "Error retrieving iRule")
}
