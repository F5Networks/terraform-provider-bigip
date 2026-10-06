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

func snatpoolUnitResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceBigipLtmSnatpool().Schema, raw)
}

func snatpoolDefaultRaw(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":    name,
		"members": []interface{}{"10.10.10.1", "10.10.10.2"},
	}
}

func TestResourceBigipLtmSnatpoolCreate(t *testing.T) {
	name := "/Common/test-snatpool"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snatpool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/snatpool/~Common~test-snatpool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","members":["10.10.10.1","10.10.10.2"]}`, name)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatpoolUnitResourceData(t, snatpoolDefaultRaw(name))

	diags := resourceBigipLtmSnatpoolCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, name, d.Id())
}

func TestResourceBigipLtmSnatpoolCreate_error(t *testing.T) {
	name := "/Common/test-snatpool"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snatpool", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatpoolUnitResourceData(t, snatpoolDefaultRaw(name))

	diags := resourceBigipLtmSnatpoolCreate(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmSnatpoolRead(t *testing.T) {
	name := "/Common/test-snatpool"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snatpool/~Common~test-snatpool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","members":["10.10.10.1"]}`, name)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatpoolUnitResourceData(t, snatpoolDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatpoolRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, name, d.Get("name"))
}

// A 404 from GetSnatPool means the SNAT pool no longer exists on the
// device; Read clears the resource ID rather than returning an error.
func TestResourceBigipLtmSnatpoolRead_notFound(t *testing.T) {
	name := "/Common/test-snatpool"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snatpool/~Common~test-snatpool", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatpoolUnitResourceData(t, snatpoolDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatpoolRead(context.Background(), d, client)
	assert.False(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmSnatpoolUpdate(t *testing.T) {
	name := "/Common/test-snatpool"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snatpool/~Common~test-snatpool", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case "GET":
			_, _ = fmt.Fprintf(w, `{"name":"%s","members":["10.10.10.1"]}`, name)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	client := newDatasourceTestClient(server.URL)
	d := snatpoolUnitResourceData(t, snatpoolDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatpoolUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmSnatpoolUpdate_error(t *testing.T) {
	name := "/Common/test-snatpool"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snatpool/~Common~test-snatpool", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatpoolUnitResourceData(t, snatpoolDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatpoolUpdate(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmSnatpoolDelete(t *testing.T) {
	name := "/Common/test-snatpool"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snatpool/~Common~test-snatpool", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatpoolUnitResourceData(t, snatpoolDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatpoolDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmSnatpoolDelete_error(t *testing.T) {
	name := "/Common/test-snatpool"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/snatpool/~Common~test-snatpool", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := snatpoolUnitResourceData(t, snatpoolDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmSnatpoolDelete(context.Background(), d, client)
	assert.True(t, diags.HasError())
}
