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

func TestResourceBigipFastAppSchema(t *testing.T) {
	r := resourceBigipFastApp()

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
	if r.Exists == nil {
		t.Fatal("Expected Exists to be defined")
	}
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	fastJsonSchema, ok := r.Schema["fast_json"]
	if !ok {
		t.Fatal("Expected field 'fast_json' to exist in schema")
	}
	if !fastJsonSchema.Required {
		t.Error("Expected field 'fast_json' to be required")
	}
	if fastJsonSchema.DiffSuppressFunc == nil {
		t.Fatal("Expected field 'fast_json' to have a DiffSuppressFunc")
	}
}

// ---------------------------------------------------------------------
// fast_json DiffSuppressFunc
// ---------------------------------------------------------------------

func TestFastAppDiffSuppressFunc(t *testing.T) {
	r := resourceBigipFastApp()
	diffSuppress := r.Schema["fast_json"].DiffSuppressFunc
	d := NewTestResourceData(t, r, map[string]interface{}{}, "")

	tests := []struct {
		name     string
		old, new string
		expected bool
	}{
		{
			name:     "identical JSON",
			old:      `{"a":1,"b":2}`,
			new:      `{"a":1,"b":2}`,
			expected: true,
		},
		{
			name:     "old has extra keys not present in new -- equal after pruning old down to new's keys",
			old:      `{"a":1,"b":2}`,
			new:      `{"a":1}`,
			expected: true,
		},
		{
			name:     "old key value changed",
			old:      `{"a":1}`,
			new:      `{"a":2}`,
			expected: false,
		},
		{
			name:     "malformed JSON on both sides unmarshals to empty maps -- treated as equal",
			old:      `{`,
			new:      `{`,
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := diffSuppress("fast_json", tc.old, tc.new, d)
			assert.Equal(t, tc.expected, got)
		})
	}
}

// ---------------------------------------------------------------------
// resourceBigipFastAppCreate / Read / Update / Delete / Exists
// ---------------------------------------------------------------------

const testFastAppJSON = `{"tenant_name":"mytenant"}`

func fastAppMockServer(t *testing.T, tenant, app string) (*httptest.Server, *http.ServeMux) {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			_, _ = fmt.Fprint(w, `{"message":[{"id":"task-1"}]}`)
		}
	})
	mux.HandleFunc("/mgmt/shared/fast/tasks/task-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"id":"task-1","code":200,"message":"success","tenant":"%s","application":"%s"}`, tenant, app)
	})
	mux.HandleFunc(fmt.Sprintf("/mgmt/shared/fast/applications/%s/%s", tenant, app), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprint(w, `{"constants":{"fast":{"view":{"tenant_name":"mytenant"}}}}`)
		case http.MethodPatch:
			_, _ = fmt.Fprint(w, `{"message":[{"id":"task-1"}]}`)
		case http.MethodDelete:
			_, _ = fmt.Fprint(w, `{"id":"task-1"}`)
		}
	})

	server := httptest.NewServer(mux)
	return server, mux
}

func TestUnitFastAppCreateReadUpdateDelete(t *testing.T) {
	tenant := "mytenant"
	app := "myapp"

	server, _ := fastAppMockServer(t, tenant, app)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call

	r := resourceBigipFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"fast_json": testFastAppJSON,
		"template":  "bigip-fast-templates/http",
	}, "")

	ctx := context.Background()

	diags := resourceBigipFastAppCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, tenant, d.Get("tenant"))
	assert.Equal(t, app, d.Get("application"))
	assert.Equal(t, app, d.Id())

	diags = resourceBigipFastAppRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Contains(t, d.Get("fast_json").(string), "mytenant")

	exists, err := resourceBigipFastAppExists(d, client)
	require.NoError(t, err)
	assert.True(t, exists)

	diags = resourceBigipFastAppUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipFastAppDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, "", d.Id())
}

func TestUnitFastAppCreate_PostError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"post failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"fast_json": testFastAppJSON,
		"template":  "bigip-fast-templates/http",
	}, "")

	diags := resourceBigipFastAppCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastAppRead_UnexpectedEndOfJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		// Empty body triggers "unexpected end of JSON input" from getForEntity's
		// json.Unmarshal, which the resource's Read treats as "not found"
		// rather than a hard error.
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"fast_json": testFastAppJSON,
		"tenant":    "mytenant",
	}, "myapp")

	diags := resourceBigipFastAppRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "expected nil diags for the unexpected-end-of-JSON branch")
	assert.Equal(t, "", d.Id())
}

func TestUnitFastAppRead_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant": "mytenant",
	}, "myapp")

	diags := resourceBigipFastAppRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastAppExists_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant": "mytenant",
	}, "myapp")

	exists, err := resourceBigipFastAppExists(d, client)
	require.Error(t, err)
	assert.False(t, exists)
}

func TestUnitFastAppExists_UnexpectedEndOfJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant": "mytenant",
	}, "myapp")

	exists, err := resourceBigipFastAppExists(d, client)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestUnitFastAppUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"fast_json": testFastAppJSON,
		"tenant":    "mytenant",
	}, "myapp")

	diags := resourceBigipFastAppUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitFastAppDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/fast/applications/mytenant/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipFastApp()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"tenant": "mytenant",
	}, "myapp")

	diags := resourceBigipFastAppDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
