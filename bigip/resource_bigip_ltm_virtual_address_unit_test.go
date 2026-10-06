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

func vaUnitResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceBigipLtmVirtualAddress().Schema, raw)
}

func vaDefaultRaw(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":            name,
		"arp":             true,
		"auto_delete":     true,
		"conn_limit":      0,
		"enabled":         true,
		"icmp_echo":       "enabled",
		"advertize_route": "disabled",
		"traffic_group":   "/Common/traffic-group-1",
	}
}

func TestResourceBigipLtmVirtualAddressCreate(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
			return
		}
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"test-va","fullPath":"%s","arp":"enabled","autoDelete":"true","connectionLimit":0,"enabled":"yes","icmpEcho":"enabled","routeAdvertisement":"disabled","trafficGroup":"/Common/traffic-group-1"}]}`, name)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))

	diags := resourceBigipLtmVirtualAddressCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, name, d.Id())
}

func TestResourceBigipLtmVirtualAddressCreate_error(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))

	diags := resourceBigipLtmVirtualAddressCreate(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmVirtualAddressRead(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"test-va","fullPath":"%s","arp":"enabled","autoDelete":"true","connectionLimit":5,"enabled":"yes","icmpEcho":"enabled","routeAdvertisement":"disabled","trafficGroup":"/Common/traffic-group-1"}]}`, name)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmVirtualAddressRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, name, d.Get("name"))
	assert.Equal(t, 5, d.Get("conn_limit"))
}

func TestResourceBigipLtmVirtualAddressRead_listError(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmVirtualAddressRead(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmVirtualAddressRead_notFound(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[{"name":"other","fullPath":"/Common/other"}]}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmVirtualAddressRead(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmVirtualAddressUpdate(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"test-va","fullPath":"%s","arp":"enabled","autoDelete":"true","connectionLimit":0,"enabled":"yes","icmpEcho":"enabled","routeAdvertisement":"disabled","trafficGroup":"/Common/traffic-group-1"}]}`, name)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmVirtualAddressUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmVirtualAddressUpdate_error(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmVirtualAddressUpdate(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmVirtualAddressExists(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"test-va","fullPath":"%s"}]}`, name)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	exists, err := resourceBigipLtmVirtualAddressExists(d, client)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestResourceBigipLtmVirtualAddressExists_notFound(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[{"name":"other","fullPath":"/Common/other"}]}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	exists, err := resourceBigipLtmVirtualAddressExists(d, client)
	require.NoError(t, err)
	assert.False(t, exists)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmVirtualAddressExists_error(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	exists, err := resourceBigipLtmVirtualAddressExists(d, client)
	assert.Error(t, err)
	assert.False(t, exists)
}

func TestResourceBigipLtmVirtualAddressDelete(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"test-va","fullPath":"%s"}]}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/virtual-address/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmVirtualAddressDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmVirtualAddressDelete_notExists(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[{"name":"other","fullPath":"/Common/other"}]}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmVirtualAddressDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmVirtualAddressDelete_error(t *testing.T) {
	name := "/Common/test-va"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/mgmt/tm/ltm/virtual-address", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"test-va","fullPath":"%s"}]}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/virtual-address/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"boom"}`)
	})

	client := newDatasourceTestClient(server.URL)
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	diags := resourceBigipLtmVirtualAddressDelete(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestModifyNameForRouteDomain(t *testing.T) {
	assert.Equal(t, "/Common/10.10.10.10%251", modifyNameForRouteDomain("/Common/10.10.10.10%1"))
	assert.Equal(t, "noSlash", modifyNameForRouteDomain("noSlash"))
}

func TestHydrateVirtualAddress(t *testing.T) {
	name := "/Common/test-va"
	d := vaUnitResourceData(t, vaDefaultRaw(name))
	d.SetId(name)

	va := hydrateVirtualAddress(d)
	assert.Equal(t, name, va.Name)
	assert.True(t, va.ARP)
	assert.Equal(t, "enabled", va.ICMPEcho)
}
