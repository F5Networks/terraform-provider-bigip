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

func snatUnitResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceBigipLtmSnat().Schema, raw)
}

func snatDefaultRaw(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":          name,
		"partition":     "Common",
		"full_path":     name,
		"autolasthop":   "default",
		"mirror":        "disabled",
		"sourceport":    "preserve",
		"translation":   "",
		"snatpool":      "",
		"vlansdisabled": true,
		"vlans":         []interface{}{},
		"origins": []interface{}{
			map[string]interface{}{
				"name":        "10.10.10.0/24",
				"app_service": "",
			},
		},
	}
}

func TestResourceBigipLtmSnatCreate(t *testing.T) {
	name := "/Common/test-snat"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snat", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"test-snat","fullPath":"%s","autoLasthop":"default","mirror":"disabled","sourcePort":"preserve"}`, name)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatUnitResourceData(t, snatDefaultRaw(name))

	diags := resourceBigipLtmSnatCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, name, d.Id())
}

func TestResourceBigipLtmSnatCreate_error(t *testing.T) {
	name := "/Common/test-snat"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snat", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatUnitResourceData(t, snatDefaultRaw(name))

	diags := resourceBigipLtmSnatCreate(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmSnatRead(t *testing.T) {
	name := "/Common/test-snat"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"test-snat","fullPath":"%s","autoLasthop":"default","mirror":"disabled","sourcePort":"preserve","translation":"","snatpool":"","origins":[{"name":"10.10.10.0/24"}]}`, name)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatUnitResourceData(t, snatDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, name, d.Get("name"))
	assert.Equal(t, "default", d.Get("autolasthop"))
}

// A 404 from GetSnat means the SNAT no longer exists on the device; Read
// clears the resource ID rather than returning an error.
func TestResourceBigipLtmSnatRead_notFound(t *testing.T) {
	name := "/Common/test-snat"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatUnitResourceData(t, snatDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatRead(context.Background(), d, client)
	assert.False(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmSnatUpdate(t *testing.T) {
	name := "/Common/test-snat"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PATCH":
			_, _ = fmt.Fprintf(w, `{"name":"test-snat","fullPath":"%s"}`, name)
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"test-snat","fullPath":"%s","autoLasthop":"default"}`, name)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := snatUnitResourceData(t, snatDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmSnatUpdate_error(t *testing.T) {
	name := "/Common/test-snat"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatUnitResourceData(t, snatDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatUpdate(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmSnatDelete(t *testing.T) {
	name := "/Common/test-snat"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatUnitResourceData(t, snatDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmSnatDelete_error(t *testing.T) {
	name := "/Common/test-snat"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snat/~Common~test-snat", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatUnitResourceData(t, snatDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatDelete(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestDataToSnat(t *testing.T) {
	name := "/Common/test-snat"
	raw := snatDefaultRaw(name)
	raw["vlansdisabled"] = false
	raw["vlans"] = []interface{}{"/Common/vlan1"}
	d := snatUnitResourceData(t, raw)

	p := dataToSnat(name, d)
	assert.Equal(t, name, p.Name)
	assert.True(t, p.VlansEnabled)
	assert.False(t, p.VlansDisabled)
	assert.Len(t, p.Origins, 1)
}
