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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestDataSourceBigipAs3DeviceInformationSchema(t *testing.T) {
	r := dataSourceBigipAs3()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	tenantSchema, ok := r.Schema["tenant"]
	if !ok {
		t.Fatal("Expected field 'tenant' to exist in schema")
	}
	if !tenantSchema.Required {
		t.Error("Expected field 'tenant' to be required")
	}

	appsSchema, ok := r.Schema["applications"]
	if !ok {
		t.Fatal("Expected field 'applications' to exist in schema")
	}
	if !appsSchema.Optional {
		t.Error("Expected field 'applications' to be optional")
	}

	as3JsonSchema, ok := r.Schema["as3_json"]
	if !ok {
		t.Fatal("Expected field 'as3_json' to exist in schema")
	}
	if !as3JsonSchema.Computed {
		t.Error("Expected field 'as3_json' to be computed")
	}
}

// ---------------------------------------------------------------------
// extractApplications
// ---------------------------------------------------------------------

func TestExtractApplications_Empty(t *testing.T) {
	r := dataSourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant": "test-tenant",
	}, "")

	apps := extractApplications(d)
	assert.Empty(t, apps)
}

func TestExtractApplications_WithValues(t *testing.T) {
	r := dataSourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":       "test-tenant",
		"applications": []interface{}{"app1", "app2"},
	}, "")

	apps := extractApplications(d)
	assert.Equal(t, []string{"app1", "app2"}, apps)
}

// ---------------------------------------------------------------------
// validateAndNormalizeTenant
// ---------------------------------------------------------------------

func TestValidateAndNormalizeTenant(t *testing.T) {
	tests := []struct {
		name        string
		tenant      string
		expected    string
		expectError bool
	}{
		{name: "simple tenant name", tenant: "mytenant", expected: "mytenant", expectError: false},
		{name: "tenant with partition prefix", tenant: "/Common/mytenant", expected: "Common~mytenant", expectError: false},
		{name: "empty tenant", tenant: "", expected: "", expectError: true},
		{name: "whitespace-only tenant", tenant: "   ", expected: "", expectError: true},
		{name: "invalid characters", tenant: "bad$tenant!", expected: "", expectError: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validateAndNormalizeTenant(tc.tenant)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

// ---------------------------------------------------------------------
// fetchAs3Configuration (uses a locally-owned httptest.Server)
// ---------------------------------------------------------------------

func TestFetchAs3Configuration_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{"class":"ADC","mytenant":{"class":"Tenant"}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	resp, err := fetchAs3Configuration(client, "mytenant", nil)
	require.NoError(t, err)
	assert.Contains(t, resp, "mytenant")
}

func TestFetchAs3Configuration_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/missing-tenant", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"HTTP 404 :: Not Found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	_, err := fetchAs3Configuration(client, "missing-tenant", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found on BIG-IP system")
}

func TestFetchAs3Configuration_Unauthorized(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/unauth-tenant", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, `{"code":401,"message":"HTTP 401 :: Unauthorized"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	_, err := fetchAs3Configuration(client, "unauth-tenant", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized access")
}

func TestFetchAs3Configuration_OtherError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/err-tenant", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	_, err := fetchAs3Configuration(client, "err-tenant", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch AS3 configuration")
}

// ---------------------------------------------------------------------
// filterAs3JSON
// ---------------------------------------------------------------------

func TestFilterAs3JSON_NoFiltering(t *testing.T) {
	as3Resp := `{"declaration":{"mytenant":{"class":"Tenant"}}}`
	result, err := filterAs3JSON(as3Resp, nil)
	require.NoError(t, err)
	assert.Equal(t, as3Resp, result)
}

func TestFilterAs3JSON_WithFiltering(t *testing.T) {
	as3Resp := `{"declaration":{"mytenant":{"class":"Tenant","app1":{"class":"Application"},"app2":{"class":"Application"}}}}`
	result, err := filterAs3JSON(as3Resp, []string{"app1"})
	require.NoError(t, err)
	assert.Contains(t, result, "app1")
	assert.NotContains(t, result, "app2")
}

func TestFilterAs3JSON_AppNotFound(t *testing.T) {
	as3Resp := `{"declaration":{"mytenant":{"class":"Tenant"}}}`
	result, err := filterAs3JSON(as3Resp, []string{"missing-app"})
	require.NoError(t, err)
	assert.Equal(t, "{}", result)
}

func TestFilterAs3JSON_InvalidTenantEntry(t *testing.T) {
	as3Resp := `{"declaration":{"mytenant":"not-an-object"}}`
	result, err := filterAs3JSON(as3Resp, []string{"app1"})
	require.NoError(t, err)
	assert.Equal(t, "{}", result)
}

func TestFilterAs3JSON_MalformedJSON(t *testing.T) {
	_, err := filterAs3JSON(`{"declaration":`, []string{"app1"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse AS3 JSON response")
}

func TestFilterAs3JSON_MissingDeclaration(t *testing.T) {
	_, err := filterAs3JSON(`{"class":"ADC"}`, []string{"app1"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing 'declaration'")
}

// ---------------------------------------------------------------------
// dataSourceBigipAs3Read (full read function against a mock server)
// ---------------------------------------------------------------------

func TestUnitDataSourceAs3Read_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"class":"ADC","mytenant":{"class":"Tenant"}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := dataSourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant": "mytenant",
	}, "")

	diags := dataSourceBigipAs3Read(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "mytenant", d.Id())
	assert.Contains(t, d.Get("as3_json").(string), "mytenant")
}

func TestUnitDataSourceAs3Read_InvalidTenant(t *testing.T) {
	client := NewUnitTestClient("http://127.0.0.1:0")
	r := dataSourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant": "",
	}, "")

	diags := dataSourceBigipAs3Read(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDataSourceAs3Read_FetchError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/missing-tenant", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"HTTP 404 :: Not Found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := dataSourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant": "missing-tenant",
	}, "")

	diags := dataSourceBigipAs3Read(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitDataSourceAs3Read_WithApplicationFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"mytenant":{"class":"Tenant","app1":{"class":"Application"}}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := dataSourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant":       "mytenant",
		"applications": []interface{}{"app1"},
	}, "")

	diags := dataSourceBigipAs3Read(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Contains(t, d.Get("as3_json").(string), "app1")
}
