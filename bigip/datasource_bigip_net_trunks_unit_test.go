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

func TestDataSourceBigipNetTrunksSchema(t *testing.T) {
	d := dataSourceBigipNetTrunks()
	require.NotNil(t, d.Schema)
	require.NotNil(t, d.ReadContext)
	_, ok := d.Schema["trunks"]
	require.True(t, ok, "expected 'trunks' to exist in schema")
}

func TestUnitDataSourceNetTrunksRead(t *testing.T) {
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/net/trunk", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[
			{"name":"lag1","fullPath":"lag1","interfaces":["1.1","1.2"],"lacp":"enabled","lacpMode":"active","lacpTimeout":"long","distributionHash":"src-dst-ipport","linkSelectPolicy":"auto","bandwidth":20000,"id":1,"stp":"enabled","type":"lacp","workingMbrCount":2}
		]}`)
	})

	d := dataSourceBigipNetTrunks()
	rd := NewTestResourceData(t, d, map[string]interface{}{}, "")

	diags := dataSourceBigipNetTrunksRead(context.Background(), rd, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	trunks := rd.Get("trunks").([]interface{})
	require.Len(t, trunks, 1)
	trunk := trunks[0].(map[string]interface{})
	assert.Equal(t, "lag1", trunk["name"])
	assert.Equal(t, "active", trunk["lacp_mode"])
	assert.Equal(t, "long", trunk["lacp_timeout"])
	members := trunk["interfaces"].([]interface{})
	require.Len(t, members, 2)
	assert.Equal(t, "1.1", members[0])
	assert.Equal(t, "1.2", members[1])
	assert.Equal(t, 2, trunk["working_member_count"])

	assert.NotEmpty(t, rd.Id())
}

func TestUnitDataSourceNetTrunksReadEmpty(t *testing.T) {
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/net/trunk", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})

	d := dataSourceBigipNetTrunks()
	rd := NewTestResourceData(t, d, map[string]interface{}{}, "")

	diags := dataSourceBigipNetTrunksRead(context.Background(), rd, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Empty(t, rd.Get("trunks").([]interface{}))
}

func TestUnitDataSourceNetTrunksReadError(t *testing.T) {
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/net/trunk", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})

	d := dataSourceBigipNetTrunks()
	rd := NewTestResourceData(t, d, map[string]interface{}{}, "")

	diags := dataSourceBigipNetTrunksRead(context.Background(), rd, client)
	require.True(t, diags.HasError())
}
