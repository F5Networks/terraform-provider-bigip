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

func TestResourceBigipAuthRadiusServerSchema(t *testing.T) {
	r := resourceBigipAuthRadiusServer()
	require.NotNil(t, r.Schema)
	assert.True(t, r.Schema["secret"].Sensitive)
	assert.True(t, r.Schema["name"].ForceNew)
}

func TestUnitAuthRadiusServerCreateReadUpdateDelete(t *testing.T) {
	name := "radius-1"
	mangled := "/mgmt/tm/auth/radius-server/" + MangleFullPath(name)
	var sawPost, sawPut, sawDelete bool

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/auth/radius-server", func(w http.ResponseWriter, r *http.Request) {
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":"%s","server":"10.1.1.1","port":1812}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			sawPut = true
			_, _ = fmt.Fprintf(w, `{"name":"%s","server":"10.1.1.2","port":1812,"timeout":5}`, name)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","server":"10.1.1.2","port":1812,"timeout":5}`, name)
		case http.MethodDelete:
			sawDelete = true
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	r := resourceBigipAuthRadiusServer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":   name,
		"server": "10.1.1.1",
		"secret": "shh",
		"port":   1812,
	}, "")

	ctx := context.Background()
	diags := resourceBigipAuthRadiusServerCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.Equal(t, name, d.Id())

	diags = resourceBigipAuthRadiusServerRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "10.1.1.2", d.Get("server"))
	assert.Equal(t, 5, d.Get("timeout"))

	diags = resourceBigipAuthRadiusServerUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPut)

	diags = resourceBigipAuthRadiusServerDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitAuthRadiusServerReadNotFound(t *testing.T) {
	name := "missing-radius"
	mangled := "/mgmt/tm/auth/radius-server/" + MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipAuthRadiusServer()
	d := NewTestResourceData(t, r, map[string]interface{}{}, name)

	diags := resourceBigipAuthRadiusServerRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}

func TestUnitAuthRadiusCreateReadUpdateDelete(t *testing.T) {
	mangled := "/mgmt/tm/auth/radius/" + MangleFullPath(authRadiusFullPath)
	var sawPost, sawPut, sawDelete bool

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/auth/radius", func(w http.ResponseWriter, r *http.Request) {
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":"%s","servers":["/Common/radius-1"],"serviceType":"authenticate-only","retries":3}`, authRadiusFullPath)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			sawPut = true
			_, _ = fmt.Fprintf(w, `{"name":"%s","servers":["/Common/radius-1","/Common/radius-2"],"serviceType":"authenticate-only","retries":5}`, authRadiusFullPath)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","servers":["/Common/radius-1","/Common/radius-2"],"serviceType":"authenticate-only","retries":5}`, authRadiusFullPath)
		case http.MethodDelete:
			sawDelete = true
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	r := resourceBigipAuthRadius()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"servers": []interface{}{"radius-1"},
	}, "")

	ctx := context.Background()
	diags := resourceBigipAuthRadiusCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.Equal(t, authRadiusFullPath, d.Id())

	diags = resourceBigipAuthRadiusRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	servers := d.Get("servers").([]interface{})
	require.Len(t, servers, 2)
	assert.Equal(t, "radius-1", servers[0], "expected full-path prefix to be stripped")
	assert.Equal(t, "radius-2", servers[1])
	assert.Equal(t, 5, d.Get("retries"))

	diags = resourceBigipAuthRadiusUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPut)

	diags = resourceBigipAuthRadiusDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}
