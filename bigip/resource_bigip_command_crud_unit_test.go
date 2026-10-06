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

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipCommandSchema(t *testing.T) {
	r := resourceBigipCommand()

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

	commandsSchema, ok := r.Schema["commands"]
	if !ok {
		t.Fatal("Expected field 'commands' to exist in schema")
	}
	if !commandsSchema.Required {
		t.Error("Expected field 'commands' to be required")
	}
	if r.Schema["when"].Default != "apply" {
		t.Errorf("Expected 'when' to default to 'apply', got %v", r.Schema["when"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitCommandCreate_Apply(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/util/bash", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{"commandResult":"ok"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true // skip real telemetry network call
	r := resourceBigipCommand()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"when":     "apply",
		"commands": []interface{}{"show sys version"},
	}, "")

	diags := resourceBigipCommandCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, "apply", d.Id())
	results := d.Get("command_result").([]interface{})
	require.Len(t, results, 1)
	require.Equal(t, "ok", results[0].(string))
}

func TestUnitCommandCreate_WhenDestroySkipsCommands(t *testing.T) {
	// "when" != "apply" (e.g. "destroy") means Create doesn't run any
	// commands at all -- no HTTP call should be made.
	client := NewUnitTestClient("http://127.0.0.1:0")
	client.Teem = true
	r := resourceBigipCommand()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"when":     "destroy",
		"commands": []interface{}{"show sys version"},
	}, "")

	diags := resourceBigipCommandCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, "destroy", d.Id())
}

func TestUnitCommandCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/util/bash", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"command failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipCommand()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"when":     "apply",
		"commands": []interface{}{"show sys version"},
	}, "")

	diags := resourceBigipCommandCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitCommandRead_Noop(t *testing.T) {
	client := NewUnitTestClient("http://127.0.0.1:0")
	r := resourceBigipCommand()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"when":     "apply",
		"commands": []interface{}{"show sys version"},
	}, "apply")

	diags := resourceBigipCommandRead(context.Background(), d, client)
	require.False(t, diags.HasError())
}

func TestUnitCommandUpdate_Apply(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/util/bash", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commandResult":"updated"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCommand()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"when":     "apply",
		"commands": []interface{}{"show sys version"},
	}, "apply")

	diags := resourceBigipCommandUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	results := d.Get("command_result").([]interface{})
	require.Len(t, results, 1)
	require.Equal(t, "updated", results[0].(string))
}

func TestUnitCommandUpdate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/util/bash", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"command failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCommand()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"when":     "apply",
		"commands": []interface{}{"show sys version"},
	}, "apply")

	diags := resourceBigipCommandUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitCommandDelete_Destroy(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/util/bash", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"commandResult":"deleted"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCommand()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"when":     "destroy",
		"commands": []interface{}{"show sys version"},
	}, "destroy")

	diags := resourceBigipCommandDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitCommandDelete_WhenApplySkipsCommands(t *testing.T) {
	client := NewUnitTestClient("http://127.0.0.1:0")
	r := resourceBigipCommand()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"when":     "apply",
		"commands": []interface{}{"show sys version"},
	}, "apply")

	diags := resourceBigipCommandDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitCommandDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/util/bash", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipCommand()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"when":     "destroy",
		"commands": []interface{}{"show sys version"},
	}, "destroy")

	diags := resourceBigipCommandDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
