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

func TestResourceBigipAuthTacacsSchema(t *testing.T) {
	r := resourceBigipAuthTacacs()
	require.NotNil(t, r.Schema)
	assert.True(t, r.Schema["secret"].Sensitive)
}

func TestUnitAuthTacacsCreateReadUpdateDelete(t *testing.T) {
	mangled := "/mgmt/tm/auth/tacacs/" + MangleFullPath(authTacacsFullPath)
	var sawPost, sawPut, sawDelete bool

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/auth/tacacs", func(w http.ResponseWriter, r *http.Request) {
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":"%s","servers":["10.2.2.1"],"service":"ppp","protocol":"ip"}`, authTacacsFullPath)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			sawPut = true
			_, _ = fmt.Fprintf(w, `{"name":"%s","servers":["10.2.2.1","10.2.2.2"],"service":"ppp","protocol":"ip","encryption":"enabled","accounting":"send-to-all-servers"}`, authTacacsFullPath)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","servers":["10.2.2.1","10.2.2.2"],"service":"ppp","protocol":"ip","encryption":"enabled","accounting":"send-to-all-servers"}`, authTacacsFullPath)
		case http.MethodDelete:
			sawDelete = true
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	r := resourceBigipAuthTacacs()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"servers": []interface{}{"10.2.2.1"},
		"secret":  "shh",
	}, "")

	ctx := context.Background()
	diags := resourceBigipAuthTacacsCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.Equal(t, authTacacsFullPath, d.Id())

	diags = resourceBigipAuthTacacsRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	servers := d.Get("servers").([]interface{})
	require.Len(t, servers, 2)
	assert.Equal(t, "send-to-all-servers", d.Get("accounting"))

	diags = resourceBigipAuthTacacsUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPut)

	diags = resourceBigipAuthTacacsDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitAuthTacacsReadNotFound(t *testing.T) {
	mangled := "/mgmt/tm/auth/tacacs/" + MangleFullPath(authTacacsFullPath)
	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipAuthTacacs()
	d := NewTestResourceData(t, r, map[string]interface{}{}, authTacacsFullPath)

	diags := resourceBigipAuthTacacsRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}
