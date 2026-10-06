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

func TestResourceBigipAuthLdapSchema(t *testing.T) {
	r := resourceBigipAuthLdap()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.CreateContext)
	require.NotNil(t, r.ReadContext)
	require.NotNil(t, r.UpdateContext)
	require.NotNil(t, r.DeleteContext)
	require.NotNil(t, r.Importer)
	assert.True(t, r.Schema["bind_pw"].Sensitive)
}

func TestUnitAuthLdapCreateReadUpdateDelete(t *testing.T) {
	mangled := "/mgmt/tm/auth/ldap/" + MangleFullPath(authLdapFullPath)
	var sawPost, sawPut, sawDelete bool

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/auth/ldap", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":"%s","servers":["ldap1.example.com"],"port":389}`, authLdapFullPath)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			sawPut = true
			_, _ = fmt.Fprintf(w, `{"name":"%s","servers":["ldap1.example.com"],"port":636,"ssl":"enabled"}`, authLdapFullPath)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","servers":["ldap1.example.com"],"port":636,"ssl":"enabled"}`, authLdapFullPath)
		case http.MethodDelete:
			sawDelete = true
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	r := resourceBigipAuthLdap()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"servers": []interface{}{"ldap1.example.com"},
		"port":    389,
	}, "")

	ctx := context.Background()

	diags := resourceBigipAuthLdapCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.Equal(t, authLdapFullPath, d.Id())

	diags = resourceBigipAuthLdapRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, 636, d.Get("port"))

	require.NoError(t, d.Set("ssl", "enabled"))
	diags = resourceBigipAuthLdapUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPut)

	diags = resourceBigipAuthLdapDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitAuthLdapReadNotFound(t *testing.T) {
	mangled := "/mgmt/tm/auth/ldap/" + MangleFullPath(authLdapFullPath)
	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"01020036:3: The requested LDAP Authentication Configuration (/Common/system-auth) was not found."}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipAuthLdap()
	d := NewTestResourceData(t, r, map[string]interface{}{}, authLdapFullPath)

	diags := resourceBigipAuthLdapRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}
