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

func testPoolResourceData(t *testing.T, name string) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resourceBigipLtmPool().Schema, map[string]interface{}{
		"name":                   name,
		"monitors":               []interface{}{"/Common/http"},
		"allow_nat":              "yes",
		"allow_snat":             "yes",
		"description":            "test pool",
		"load_balancing_mode":    "round-robin",
		"minimum_active_members": 1,
		"slow_ramp_time":         10,
		"service_down_action":    "none",
		"reselect_tries":         3,
	})
	d.SetId(name)
	return d
}

func TestResourceBigipLtmPoolCreate_success(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"

	mux.HandleFunc("/mgmt/tm/ltm/pool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"test-pool","monitor":"/Common/http"}`)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	client := newDatasourceTestClient(server.URL)
	client.Teem = true // skip real telemetry network call
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, name, d.Id())
}

func TestResourceBigipLtmPoolCreate_createError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"

	mux.HandleFunc("/mgmt/tm/ltm/pool", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"create failed"}`)
	})

	client := newDatasourceTestClient(server.URL)
	client.Teem = true
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPoolRead_success(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"test-pool","allowNat":"yes","allowSnat":"yes","loadBalancingMode":"round-robin","slowRampTime":10,"minActiveMembers":1,"serviceDownAction":"none","reselectTries":3,"description":"test pool","monitor":"/Common/http "}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "yes", d.Get("allow_nat"))
	assert.Equal(t, "round-robin", d.Get("load_balancing_mode"))
}

func TestResourceBigipLtmPoolRead_notFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, `{"code":404,"message":"object not found"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmPoolRead_otherError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"internal server error occurred"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPoolUpdate_success(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			_, _ = fmt.Fprintf(w, `{"name":"test-pool"}`)
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"test-pool","monitor":"/Common/http"}`)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmPoolUpdate_modifyErrorDeleteSuccess(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"code":500,"message":"modify failed"}`)
		case "DELETE":
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPoolUpdate_modifyErrorDeleteError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"code":500,"message":"modify failed"}`)
		case "DELETE":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"code":500,"message":"delete failed"}`)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestResourceBigipLtmPoolDelete_success(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmPoolDelete_error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	name := "/Common/test-pool"
	mux.HandleFunc("/mgmt/tm/ltm/pool/~Common~test-pool", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"delete failed"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := testPoolResourceData(t, name)

	diags := resourceBigipLtmPoolDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}
