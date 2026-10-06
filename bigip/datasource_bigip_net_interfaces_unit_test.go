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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataSourceBigipNetInterfacesSchema(t *testing.T) {
	d := dataSourceBigipNetInterfaces()
	require.NotNil(t, d.Schema)
	require.NotNil(t, d.ReadContext)
	_, ok := d.Schema["interfaces"]
	require.True(t, ok, "expected 'interfaces' to exist in schema")
}

func TestUnitDataSourceNetInterfacesRead(t *testing.T) {
	var sawStatsRequests int
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/net/interface", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[
			{"name":"1.1","fullPath":"1.1","enabled":true,"mediaActive":"none","mediaFixed":"10000T-FD","mediaMax":"auto","mediaSfp":"auto","macAddress":"fa:16:3e:3f:ef:f8","mtu":1500,"bundle":"not-supported","ifIndex":48,"flowControl":"tx-rx","lldpAdmin":"txonly","stp":"enabled","stpLinkType":"auto","preferPort":"sfp"},
			{"name":"mgmt","fullPath":"mgmt","enabled":true,"mediaActive":"100TX-FD","mtu":9000}
		]}`)
	})
	// A single bulk stats request should cover every interface -- verified
	// below via sawStatsRequests, so this data source doesn't regress back
	// into one HTTP round trip per interface.
	mux.HandleFunc("/mgmt/tm/net/interface/stats", func(w http.ResponseWriter, r *http.Request) {
		sawStatsRequests++
		_, _ = fmt.Fprint(w, `{"entries":{
			"https://localhost/mgmt/tm/net/interface/1.1/stats":{"nestedStats":{"entries":{"tmName":{"description":"1.1"},"status":{"description":"uninit"},"mediaActive":{"description":"none"}}}},
			"https://localhost/mgmt/tm/net/interface/mgmt/stats":{"nestedStats":{"entries":{"tmName":{"description":"mgmt"},"status":{"description":"up"},"mediaActive":{"description":"100TX-FD"}}}}
		}}`)
	})

	d := dataSourceBigipNetInterfaces()
	rd := NewTestResourceData(t, d, map[string]interface{}{}, "")

	diags := dataSourceBigipNetInterfacesRead(context.Background(), rd, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, 1, sawStatsRequests, "expected exactly one bulk stats request regardless of interface count")

	ifaces := rd.Get("interfaces").([]interface{})
	require.Len(t, ifaces, 2)

	byName := map[string]map[string]interface{}{}
	for _, raw := range ifaces {
		m := raw.(map[string]interface{})
		byName[m["name"].(string)] = m
	}

	i11 := byName["1.1"]
	require.NotNil(t, i11)
	assert.Equal(t, "uninit", i11["status"])
	assert.Equal(t, "none", i11["media_active"])
	assert.Equal(t, "10000T-FD", i11["media_fixed"])
	assert.Equal(t, true, i11["enabled"])

	mgmt := byName["mgmt"]
	require.NotNil(t, mgmt)
	assert.Equal(t, "up", mgmt["status"])
	assert.Equal(t, "100TX-FD", mgmt["media_active"])

	assert.NotEmpty(t, rd.Id())
}

func TestUnitDataSourceNetInterfacesReadStatsError(t *testing.T) {
	// If the bulk stats request fails, the data source should still
	// succeed overall and just leave every interface's status empty,
	// rather than failing the whole read.
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/net/interface", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[{"name":"1.1","fullPath":"1.1","enabled":true}]}`)
	})
	mux.HandleFunc("/mgmt/tm/net/interface/stats", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})

	d := dataSourceBigipNetInterfaces()
	rd := NewTestResourceData(t, d, map[string]interface{}{}, "")

	diags := dataSourceBigipNetInterfacesRead(context.Background(), rd, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	ifaces := rd.Get("interfaces").([]interface{})
	require.Len(t, ifaces, 1)
	assert.Equal(t, "", ifaces[0].(map[string]interface{})["status"])
}

func TestUnitDataSourceNetInterfacesReadError(t *testing.T) {
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/net/interface", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})

	d := dataSourceBigipNetInterfaces()
	rd := NewTestResourceData(t, d, map[string]interface{}{}, "")

	diags := dataSourceBigipNetInterfacesRead(context.Background(), rd, client)
	require.True(t, diags.HasError())
}
