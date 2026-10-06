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

// iruleResourceData builds a *schema.ResourceData for the bigip_ltm_irule
// resource using schema.TestResourceDataRaw.
func iruleResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceBigipLtmIRule().Schema, raw)
}

// TestResourceBigipLtmIRuleCreateSuccess covers the Create success path:
// CreateIRule succeeds, the ID is set, and Read is invoked and succeeds.
func TestResourceBigipLtmIRuleCreateSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~test-irule", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"/Common/test-irule","fullPath":"/Common/test-irule","apiAnonymous":"when HTTP_REQUEST { }"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleResourceData(t, map[string]interface{}{
		"name":  "/Common/test-irule",
		"irule": "when HTTP_REQUEST { }",
	})

	diags := resourceBigipLtmIRuleCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/test-irule", d.Id())
	assert.Equal(t, "/Common/test-irule", d.Get("name"))
}

// TestResourceBigipLtmIRuleCreateError covers the Create path where
// CreateIRule returns an error.
func TestResourceBigipLtmIRuleCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleResourceData(t, map[string]interface{}{
		"name":  "/Common/bad-irule",
		"irule": "when HTTP_REQUEST { }",
	})

	diags := resourceBigipLtmIRuleCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when CreateIRule fails")
	assert.Contains(t, diags[0].Summary, "error creating iRule")
}

// TestResourceBigipLtmIRuleReadSuccess covers the Read success path.
func TestResourceBigipLtmIRuleReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~read-irule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"name":"/Common/read-irule","fullPath":"/Common/read-irule","apiAnonymous":"when HTTP_REQUEST { }"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleResourceData(t, map[string]interface{}{
		"name":  "/Common/read-irule",
		"irule": "when HTTP_REQUEST { }",
	})
	d.SetId("/Common/read-irule")

	diags := resourceBigipLtmIRuleRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/Common/read-irule", d.Get("name"))
	assert.Equal(t, "when HTTP_REQUEST { }", d.Get("irule"))
}

// TestResourceBigipLtmIRuleReadNotFound covers the Read branch where IRule
// returns a 404. A 404 means the iRule no longer exists on the device;
// Read clears the resource ID rather than returning an error.
func TestResourceBigipLtmIRuleReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~missing-irule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleResourceData(t, map[string]interface{}{
		"name":  "/Common/missing-irule",
		"irule": "when HTTP_REQUEST { }",
	})
	d.SetId("/Common/missing-irule")

	diags := resourceBigipLtmIRuleRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "expected no error: a 404 clears state instead of failing")
	require.Empty(t, d.Id())
}

// TestResourceBigipLtmIRuleReadError covers the Read branch where IRule
// itself returns an error.
func TestResourceBigipLtmIRuleReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~broken-irule", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleResourceData(t, map[string]interface{}{
		"name":  "/Common/broken-irule",
		"irule": "when HTTP_REQUEST { }",
	})
	d.SetId("/Common/broken-irule")

	diags := resourceBigipLtmIRuleRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when IRule fails")
	assert.Contains(t, diags[0].Summary, "error retrieving iRule")
}

// TestResourceBigipLtmIRuleUpdateSuccess covers the Update success path
// followed by a successful Read.
func TestResourceBigipLtmIRuleUpdateSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~update-irule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "PUT" || r.Method == "PATCH" {
			_, _ = fmt.Fprint(w, `{}`)
			return
		}
		_, _ = fmt.Fprint(w, `{"name":"/Common/update-irule","fullPath":"/Common/update-irule","apiAnonymous":"when HTTP_REQUEST { }"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleResourceData(t, map[string]interface{}{
		"name":  "/Common/update-irule",
		"irule": "when HTTP_REQUEST { }",
	})
	d.SetId("/Common/update-irule")

	diags := resourceBigipLtmIRuleUpdate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

// TestResourceBigipLtmIRuleUpdateError covers the Update path where
// ModifyIRule returns an error.
func TestResourceBigipLtmIRuleUpdateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~update-error-irule", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleResourceData(t, map[string]interface{}{
		"name":  "/Common/update-error-irule",
		"irule": "when HTTP_REQUEST { }",
	})
	d.SetId("/Common/update-error-irule")

	diags := resourceBigipLtmIRuleUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when ModifyIRule fails")
}

// TestResourceBigipLtmIRuleDeleteSuccess covers the Delete success path.
func TestResourceBigipLtmIRuleDeleteSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~delete-irule", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleResourceData(t, map[string]interface{}{
		"name":  "/Common/delete-irule",
		"irule": "when HTTP_REQUEST { }",
	})
	d.SetId("/Common/delete-irule")

	diags := resourceBigipLtmIRuleDelete(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Id())
}

// TestResourceBigipLtmIRuleDeleteError covers the Delete path where
// DeleteIRule returns an error.
func TestResourceBigipLtmIRuleDeleteError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~delete-error-irule", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := iruleResourceData(t, map[string]interface{}{
		"name":  "/Common/delete-error-irule",
		"irule": "when HTTP_REQUEST { }",
	})
	d.SetId("/Common/delete-error-irule")

	diags := resourceBigipLtmIRuleDelete(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when DeleteIRule fails")
}
