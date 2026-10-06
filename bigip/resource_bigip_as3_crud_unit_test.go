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

func TestResourceBigipAs3Schema(t *testing.T) {
	r := resourceBigipAs3()

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
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	as3JsonSchema, ok := r.Schema["as3_json"]
	if !ok {
		t.Fatal("Expected field 'as3_json' to exist in schema")
	}
	if !as3JsonSchema.Optional {
		t.Error("Expected field 'as3_json' to be optional")
	}
	assert.Equal(t, []string{"delete_apps"}, as3JsonSchema.ConflictsWith)

	deleteAppsSchema, ok := r.Schema["delete_apps"]
	if !ok {
		t.Fatal("Expected field 'delete_apps' to exist in schema")
	}
	assert.Equal(t, []string{"as3_json"}, deleteAppsSchema.ConflictsWith)
}

func TestResourceBigipAs3Importer(t *testing.T) {
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "mytenant")

	results, err := r.Importer.StateContext(context.Background(), d, nil)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "mytenant", results[0].Get("tenant_list"))
	assert.Equal(t, "mytenant", results[0].Get("tenant_filter"))
}

// ---------------------------------------------------------------------
// as3_json ValidateFunc
// ---------------------------------------------------------------------

func TestResourceBigipAs3ValidateFunc(t *testing.T) {
	r := resourceBigipAs3()
	validateFn := r.Schema["as3_json"].ValidateFunc
	require.NotNil(t, validateFn)

	tests := []struct {
		name        string
		json        string
		expectError bool
	}{
		{
			name:        "valid AS3 declaration",
			json:        `{"class":"AS3","declaration":{"class":"ADC"}}`,
			expectError: false,
		},
		{
			name:        "invalid JSON",
			json:        `{"class":`,
			expectError: true,
		},
		{
			name:        "wrong class",
			json:        `{"class":"NotAS3"}`,
			expectError: true,
		},
		{
			name:        "wrong declaration class",
			json:        `{"class":"AS3","declaration":{"class":"NotADC"}}`,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, errs := validateFn(tc.json, "as3_json")
			if tc.expectError {
				assert.NotEmpty(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

// ---------------------------------------------------------------------
// validateControlsParam
// ---------------------------------------------------------------------

func TestValidateControlsParam_Valid(t *testing.T) {
	controls := map[string]interface{}{
		"dry_run":        "yes",
		"trace":          "no",
		"trace_response": "yes",
		"log_level":      "debug",
		"user_agent":     "custom-agent",
	}
	diags := validateControlsParam(controls, nil)
	assert.False(t, diags.HasError(), "unexpected error: %v", diags)
}

func TestValidateControlsParam_InvalidType(t *testing.T) {
	diags := validateControlsParam("not-a-map", nil)
	assert.True(t, diags.HasError())
}

func TestValidateControlsParam_InvalidKey(t *testing.T) {
	controls := map[string]interface{}{
		"bogus_key": "value",
	}
	diags := validateControlsParam(controls, nil)
	assert.True(t, diags.HasError())
}

func TestValidateControlsParam_InvalidBooleanValue(t *testing.T) {
	for _, key := range []string{"dry_run", "trace", "trace_response"} {
		controls := map[string]interface{}{key: "maybe"}
		diags := validateControlsParam(controls, nil)
		assert.True(t, diags.HasError(), "expected error for key %s", key)
	}
}

func TestValidateControlsParam_InvalidLogLevel(t *testing.T) {
	controls := map[string]interface{}{"log_level": "bogus"}
	diags := validateControlsParam(controls, nil)
	assert.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// controlsQueraString
// ---------------------------------------------------------------------

func TestControlsQueraString(t *testing.T) {
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"controls": map[string]interface{}{
			"dry_run":        "yes",
			"trace":          "no",
			"trace_response": "yes",
			"log_level":      "debug",
		},
	}, "")

	query := controlsQueraString(d)
	assert.Contains(t, query, "controls.dryRun=true")
	assert.Contains(t, query, "controls.trace=false")
	assert.Contains(t, query, "controls.traceResponse=true")
	assert.Contains(t, query, "controls.logLevel=debug")
}

func TestControlsQueraString_Empty(t *testing.T) {
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "")

	query := controlsQueraString(d)
	assert.Equal(t, "", query)
}

// ---------------------------------------------------------------------
// contains / GenerateRandomString
// ---------------------------------------------------------------------

func TestContains(t *testing.T) {
	assert.True(t, contains([]string{"a", "b", "c"}, "b"))
	assert.False(t, contains([]string{"a", "b", "c"}, "z"))
	assert.False(t, contains(nil, "z"))
}

func TestGenerateRandomString(t *testing.T) {
	s, err := GenerateRandomString(10)
	require.NoError(t, err)
	assert.Len(t, s, 10)

	s2, err := GenerateRandomString(10)
	require.NoError(t, err)
	assert.NotEqual(t, s, s2, "expected two random strings to differ")
}

// ---------------------------------------------------------------------
// Shared mock server helpers for the AS3 CRUD unit tests
// ---------------------------------------------------------------------

const testAs3Declaration = `{"class":"AS3","action":"deploy","persist":true,"declaration":{"class":"ADC","schemaVersion":"3.0.0","mytenant":{"class":"Tenant","myapp":{"class":"Application"}}}}`

// testAs3RawADC is what the mock GET /mgmt/shared/appsvcs/declare/{tenant}
// endpoint should return: the bare ADC body (no AS3/action/persist wrapper),
// matching what a real BIG-IP returns from that endpoint. go-bigip's GetAs3
// re-wraps this into the full AS3 envelope itself
// (as3Json["declaration"] = adcJson), so returning the already-wrapped
// testAs3Declaration from this endpoint would double-wrap the response.
const testAs3RawADC = `{"class":"ADC","schemaVersion":"3.0.0","mytenant":{"class":"Tenant","myapp":{"class":"Application"}}}`

// as3TaskCompleteHandler returns a handler that always reports the AS3 async
// task as complete for the given tenant, so PostAs3Bigip's polling loop
// resolves on the first call without any real time.Sleep.
func as3TaskCompleteHandler(tenant string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"results":[{"code":200,"message":"success","tenant":"%s"}]}`, tenant)
	}
}

func as3MockServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/mgmt/shared/appsvcs/info", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"version":"3.50.0","release":"1","schemaCurrent":"3.50.0","schemaMinimum":"3.0.0"}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/settings", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"perAppDeploymentAllowed":false}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = fmt.Fprint(w, `{"id":"task-1"}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, testAs3RawADC)
		case http.MethodDelete:
			_, _ = fmt.Fprint(w, `{"id":"del-task-1"}`)
		}
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/task/task-1", as3TaskCompleteHandler("mytenant"))
	mux.HandleFunc("/mgmt/shared/appsvcs/task/del-task-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"results":[{"code":200,"message":"success","tenant":"mytenant"}]}`)
	})

	return httptest.NewServer(mux)
}

// ---------------------------------------------------------------------
// resourceBigipAs3Create / Read / Update / Delete
// ---------------------------------------------------------------------

func TestUnitAs3CreateReadUpdateDelete(t *testing.T) {
	server := as3MockServer(t)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call

	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json": testAs3Declaration,
	}, "")

	ctx := context.Background()

	diags := resourceBigipAs3Create(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, "mytenant", d.Id())
	assert.Equal(t, "mytenant", d.Get("tenant_list"))

	diags = resourceBigipAs3Read(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.NotEmpty(t, d.Get("as3_json"))

	diags = resourceBigipAs3Update(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipAs3Delete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitAs3Create_TenantFilterNotFound(t *testing.T) {
	server := as3MockServer(t)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json":      testAs3Declaration,
		"tenant_filter": "othertenant",
	}, "")

	diags := resourceBigipAs3Create(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAs3Create_CheckSettingError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/settings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json": testAs3Declaration,
	}, "")

	diags := resourceBigipAs3Create(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAs3Create_PostFails(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/info", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"version":"3.50.0","release":"1","schemaCurrent":"3.50.0","schemaMinimum":"3.0.0"}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/settings", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"perAppDeploymentAllowed":false}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"post failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json": testAs3Declaration,
	}, "")

	diags := resourceBigipAs3Create(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAs3Create_PerApplicationMode(t *testing.T) {
	perAppDecl := `{"myapp":{"class":"Application","serviceMain":{"class":"Service_HTTP"}}}`

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/settings", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"perAppDeploymentAllowed":true}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = fmt.Fprint(w, `{"id":"perapp-task-1"}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, perAppDecl)
		}
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/task/perapp-task-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"results":[{"code":200,"message":"success"}]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json":    perAppDecl,
		"tenant_name": "mytenant",
	}, "")

	diags := resourceBigipAs3Create(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, "mytenant", d.Get("tenant_list"))
	assert.Equal(t, true, d.Get("per_app_mode"))
}

func TestUnitAs3Create_DeleteAppsRedirect(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprint(w, testAs3RawADC)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant_name": "mytenant",
		"delete_apps": []interface{}{
			map[string]interface{}{
				"tenant_name": "mytenant",
				"apps":        []interface{}{"myapp"},
			},
		},
	}, "")

	diags := resourceBigipAs3Create(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.NotEqual(t, "", d.Id())
}

func TestUnitAs3Read_UnexpectedEndOfJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		// getForEntity will fail to unmarshal an empty body into adcJson,
		// surfacing "unexpected end of JSON input".
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "mytenant")

	diags := resourceBigipAs3Read(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected nil diags for the unexpected-end-of-JSON branch")
}

func TestUnitAs3Read_ByTaskId(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/task/task-99", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"declaration":{"mytenant":{"class":"Tenant"}}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"task_id": "task-99",
	}, "")

	diags := resourceBigipAs3Read(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.NotEmpty(t, d.Get("as3_json"))
}

func TestUnitAs3Read_ByTaskIdError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/task/task-err", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	// name resolves to d.Id(), which must be "" here so Read takes the
	// "else if task_id != nil" branch instead of calling GetAs3(name, ...).
	d := NewTestResourceData(t, r, map[string]interface{}{
		"task_id": "task-err",
	}, "")

	diags := resourceBigipAs3Read(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected nil diags (task-id error path returns nil, not an error)")
	assert.Equal(t, "", d.Get("as3_json"), "as3_json should remain unset when Getas3TaskResponse fails")
}

func TestUnitAs3Read_PerAppMode(t *testing.T) {
	perAppDecl := `{"myapp":{"class":"Application"},"schemaVersion":"3.50.0"}`
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, perAppDecl)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"application_list": "myapp",
		"per_app_mode":     true,
	}, "mytenant")

	diags := resourceBigipAs3Read(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Contains(t, d.Get("as3_json").(string), "myapp")
}

func TestUnitAs3Read_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "mytenant")

	diags := resourceBigipAs3Read(context.Background(), d, client)
	require.True(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}

func TestUnitAs3Update_TenantListChanged(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/settings", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"perAppDeploymentAllowed":false}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/info", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"version":"3.50.0"}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = fmt.Fprint(w, `{"id":"update-task-1"}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, testAs3RawADC)
		}
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/task/update-task-1", as3TaskCompleteHandler("mytenant"))
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/oldtenant", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			_, _ = fmt.Fprint(w, `{"id":"del-old-task"}`)
		}
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/task/del-old-task", as3TaskCompleteHandler("oldtenant"))
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json":    testAs3Declaration,
		"tenant_list": "oldtenant",
	}, "oldtenant")

	diags := resourceBigipAs3Update(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.Equal(t, "mytenant", d.Get("tenant_list"))
}

func TestUnitAs3Update_TenantFilterMismatchWarns(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/settings", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"perAppDeploymentAllowed":false}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/info", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"version":"3.50.0"}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = fmt.Fprint(w, `{"id":"update-task-2"}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, testAs3RawADC)
		}
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/task/update-task-2", as3TaskCompleteHandler("mytenant"))
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json":      testAs3Declaration,
		"tenant_list":   "mytenant",
		"tenant_filter": "nonexistenttenant",
	}, "mytenant")

	diags := resourceBigipAs3Update(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitAs3Update_CheckSettingError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/settings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json": testAs3Declaration,
	}, "mytenant")

	diags := resourceBigipAs3Update(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAs3Update_PerAppMode(t *testing.T) {
	// oldapp is only referenced via application_list in state (below); the
	// mock's DeletePerApplicationAs3Bigip handler for it exercises the
	// "app removed from curApplicationList" deletion branch in Update.
	newPerAppDecl := `{"newapp":{"class":"Application"}}`

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/settings", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"perAppDeploymentAllowed":true}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/oldapp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = fmt.Fprint(w, `{"id":"perapp-update-task"}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, newPerAppDecl)
		}
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/task/perapp-update-task", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"results":[{"code":200,"message":"success"}]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true

	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json":         newPerAppDecl,
		"application_list": "oldapp",
		"per_app_mode":     true,
	}, "mytenant")

	diags := resourceBigipAs3Update(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.Equal(t, "mytenant", d.Get("tenant_list"))
}

func TestUnitAs3Update_PerAppModeNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/settings", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"perAppDeploymentAllowed":false}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json":     testAs3Declaration,
		"per_app_mode": true,
	}, "mytenant")

	diags := resourceBigipAs3Update(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAs3Update_DeleteAppsRedirect(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = fmt.Fprint(w, testAs3RawADC)
		}
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant_name": "mytenant",
		"delete_apps": []interface{}{
			map[string]interface{}{
				"tenant_name": "mytenant",
				"apps":        []interface{}{"myapp"},
			},
		},
	}, "mytenant")

	diags := resourceBigipAs3Update(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitAs3Delete_Success(t *testing.T) {
	server := as3MockServer(t)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json": testAs3Declaration,
	}, "mytenant")

	diags := resourceBigipAs3Delete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitAs3Delete_DryRunSkips(t *testing.T) {
	client := NewUnitTestClient("http://127.0.0.1:0")
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"controls": map[string]interface{}{
			"dry_run": "yes",
		},
	}, "mytenant")

	diags := resourceBigipAs3Delete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitAs3Delete_PerAppMode(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"application_list": "myapp",
		"per_app_mode":     true,
	}, "mytenant")

	diags := resourceBigipAs3Delete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitAs3Delete_PerAppModeError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"application_list": "myapp",
		"per_app_mode":     true,
	}, "mytenant")

	diags := resourceBigipAs3Delete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAs3Delete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"as3_json": testAs3Declaration,
	}, "mytenant")

	diags := resourceBigipAs3Delete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitAs3Delete_DeleteAppsRedirect(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant_name": "mytenant",
		"delete_apps": []interface{}{
			map[string]interface{}{
				"tenant_name": "mytenant",
				"apps":        []interface{}{"myapp"},
			},
		},
	}, "mytenant")

	diags := resourceBigipAs3Delete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
}

// ---------------------------------------------------------------------
// handleDeleteApps
// ---------------------------------------------------------------------

func TestHandleDeleteApps_TenantNotFoundSkips(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/missingtenant", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant_name": "missingtenant",
		"delete_apps": []interface{}{
			map[string]interface{}{
				"tenant_name": "missingtenant",
				"apps":        []interface{}{"myapp"},
			},
		},
	}, "")

	diags := handleDeleteApps(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected tenant-not-found to be skipped, not errored")
	assert.NotEqual(t, "", d.Id())
}

func TestHandleDeleteApps_DeleteError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, testAs3RawADC)
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/mytenant/applications/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipAs3()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant_name": "mytenant",
		"delete_apps": []interface{}{
			map[string]interface{}{
				"tenant_name": "mytenant",
				"apps":        []interface{}{"myapp"},
			},
		},
	}, "")

	diags := handleDeleteApps(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// stripAS3Metadata (also covered indirectly via the existing
// resource_bigip_as3_unit_test.go DiffSuppressFunc tests, but exercised
// directly here for completeness)
// ---------------------------------------------------------------------

func TestStripAS3Metadata_Direct(t *testing.T) {
	jsonRef := map[string]interface{}{
		"persist": true,
		"declaration": map[string]interface{}{
			"updateMode":    "selective",
			"schemaVersion": "3.0.0",
			"id":            "abc",
			"label":         "l",
			"remark":        "r",
			"Common":        map[string]interface{}{},
			"MyTenant":      map[string]interface{}{},
		},
	}

	stripAS3Metadata(jsonRef, true)

	_, hasPersist := jsonRef["persist"]
	assert.False(t, hasPersist)

	decl := jsonRef["declaration"].(map[string]interface{})
	for _, k := range []string{"updateMode", "schemaVersion", "id", "label", "remark", "Common"} {
		_, ok := decl[k]
		assert.False(t, ok, "expected key %q to be stripped", k)
	}
	_, hasTenant := decl["MyTenant"]
	assert.True(t, hasTenant, "expected user-defined tenant to remain")
}

func TestStripAS3Metadata_NoDeclaration(t *testing.T) {
	// Exercises the "declaration key missing/not a map" branch: nothing
	// should panic and persist should still be stripped.
	jsonRef := map[string]interface{}{
		"persist": true,
	}
	stripAS3Metadata(jsonRef, false)
	_, hasPersist := jsonRef["persist"]
	assert.False(t, hasPersist)
}
