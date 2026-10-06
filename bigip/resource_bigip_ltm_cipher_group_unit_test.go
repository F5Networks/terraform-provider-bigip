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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceBigipLtmCipherGroupCreate(t *testing.T) {
	name := "/Common/test-cipher-group"
	mangled := "/mgmt/tm/ltm/cipher/group/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	client.Teem = true // skip real telemetry network call
	mux.HandleFunc("/mgmt/tm/ltm/cipher/group", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{"name":"test-cipher-group","fullPath":"%s","ordering":"default","allow":[{"name":"f5-default","partition":"Common"}],"require":[]}`, name)
	})

	r := resourceBigipLtmCipherGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, "")

	diags := resourceBigipLtmCipherGroupCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())
	require.Equal(t, "default", d.Get("ordering").(string))
}

func TestUnitResourceBigipLtmCipherGroupCreate_Error(t *testing.T) {
	name := "/Common/test-cipher-group"
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/ltm/cipher/group", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})

	r := resourceBigipLtmCipherGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, "")

	diags := resourceBigipLtmCipherGroupCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "create failed")
}

func TestUnitResourceBigipLtmCipherGroupRead(t *testing.T) {
	name := "/Common/test-cipher-group"
	mangled := "/mgmt/tm/ltm/cipher/group/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{"name":"test-cipher-group","fullPath":"%s","ordering":"speed","allow":[{"name":"f5-default","partition":"Common"}],"require":[{"name":"f5-secure","partition":"Common"}]}`, name)
	})

	r := resourceBigipLtmCipherGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmCipherGroupRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "speed", d.Get("ordering").(string))
	require.Equal(t, 1, d.Get("allow").(*schema.Set).Len())
	require.Equal(t, 1, d.Get("require").(*schema.Set).Len())
}

// A 404 means the cipher group no longer exists on the device; Read clears
// the resource ID rather than returning an error.
func TestUnitResourceBigipLtmCipherGroupRead_NotFound(t *testing.T) {
	name := "/Common/test-cipher-group"
	mangled := "/mgmt/tm/ltm/cipher/group/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})

	r := resourceBigipLtmCipherGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmCipherGroupRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Equal(t, "", d.Id())
}

func TestUnitResourceBigipLtmCipherGroupUpdate(t *testing.T) {
	name := "/Common/test-cipher-group"
	mangled := "/mgmt/tm/ltm/cipher/group/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			_, _ = fmt.Fprint(w, `{}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-cipher-group","fullPath":"%s","ordering":"strength","allow":[],"require":[]}`, name)
	})

	r := resourceBigipLtmCipherGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"ordering": "strength",
		"require":  []interface{}{"/Common/f5-secure"},
	}, name)

	diags := resourceBigipLtmCipherGroupUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	require.Equal(t, "strength", d.Get("ordering").(string))
}

func TestUnitResourceBigipLtmCipherGroupUpdate_Error(t *testing.T) {
	name := "/Common/test-cipher-group"
	mangled := "/mgmt/tm/ltm/cipher/group/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})

	r := resourceBigipLtmCipherGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmCipherGroupUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "update failed")
}

func TestUnitResourceBigipLtmCipherGroupDelete(t *testing.T) {
	name := "/Common/test-cipher-group"
	mangled := "/mgmt/tm/ltm/cipher/group/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
	})

	r := resourceBigipLtmCipherGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmCipherGroupDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitResourceBigipLtmCipherGroupDelete_Error(t *testing.T) {
	name := "/Common/test-cipher-group"
	mangled := "/mgmt/tm/ltm/cipher/group/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})

	r := resourceBigipLtmCipherGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmCipherGroupDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "delete failed")
}

func TestUnitResourceBigipLtmCipherGroupSchema(t *testing.T) {
	r := resourceBigipLtmCipherGroup()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.CreateContext)
	require.NotNil(t, r.ReadContext)
	require.NotNil(t, r.UpdateContext)
	require.NotNil(t, r.DeleteContext)
}
