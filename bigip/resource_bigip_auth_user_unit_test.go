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

func TestResourceBigipAuthUserSchema(t *testing.T) {
	r := resourceBigipAuthUser()
	require.NotNil(t, r.Schema)
	assert.True(t, r.Schema["password"].Sensitive)
	assert.True(t, r.Schema["name"].ForceNew)
	assert.Equal(t, 1, r.Schema["partition_access"].MinItems)
}

func TestUnitAuthUserCreateReadUpdateDelete(t *testing.T) {
	name := "test-user"
	mangled := "/mgmt/tm/auth/user/" + MangleFullPath(name)
	var sawPost, sawPut, sawDelete bool

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/auth/user", func(w http.ResponseWriter, r *http.Request) {
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":"%s","shell":"tmsh","partitionAccess":[{"name":"Common","role":"guest"}]}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			sawPut = true
			_, _ = fmt.Fprintf(w, `{"name":"%s","description":"updated","shell":"tmsh","partitionAccess":[{"name":"Common","role":"manager"}]}`, name)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","description":"updated","shell":"tmsh","partitionAccess":[{"name":"Common","role":"manager"}]}`, name)
		case http.MethodDelete:
			sawDelete = true
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	r := resourceBigipAuthUser()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"password": "secretpw",
		"shell":    "tmsh",
		"partition_access": []interface{}{
			map[string]interface{}{"partition": "Common", "role": "guest"},
		},
	}, "")

	ctx := context.Background()
	diags := resourceBigipAuthUserCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.Equal(t, name, d.Id())

	diags = resourceBigipAuthUserRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "updated", d.Get("description"))
	access := d.Get("partition_access").([]interface{})
	require.Len(t, access, 1)
	entry := access[0].(map[string]interface{})
	assert.Equal(t, "Common", entry["partition"])
	assert.Equal(t, "manager", entry["role"])

	diags = resourceBigipAuthUserUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPut)

	diags = resourceBigipAuthUserDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitAuthUserReadNotFound(t *testing.T) {
	name := "missing-user"
	mangled := "/mgmt/tm/auth/user/" + MangleFullPath(name)
	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"01020036:3: The requested user (missing-user) was not found."}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipAuthUser()
	d := NewTestResourceData(t, r, map[string]interface{}{}, name)

	diags := resourceBigipAuthUserRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}
