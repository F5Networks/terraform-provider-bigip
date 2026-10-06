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

	"github.com/stretchr/testify/require"
)

func TestUnitResourceBigipLtmCipherRuleCreate(t *testing.T) {
	name := "/Common/test-cipher-rule"
	mangled := "/mgmt/tm/ltm/cipher/rule/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	client.Teem = true // skip real telemetry network call
	mux.HandleFunc("/mgmt/tm/ltm/cipher/rule", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = fmt.Fprintf(w, `{"name":"test-cipher-rule","fullPath":"%s","cipher":"ECDHE-RSA-AES128-GCM-SHA256","dhGroups":"DEFAULT","signatureAlgorithms":"DEFAULT"}`, name)
			return
		}
	})

	r := resourceBigipLtmCipherRule()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":   name,
		"cipher": "ECDHE-RSA-AES128-GCM-SHA256",
	}, "")

	diags := resourceBigipLtmCipherRuleCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())
	require.Equal(t, "DEFAULT", d.Get("dh_groups").(string))
}

func TestUnitResourceBigipLtmCipherRuleCreate_Error(t *testing.T) {
	name := "/Common/test-cipher-rule"
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/ltm/cipher/rule", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})

	r := resourceBigipLtmCipherRule()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":   name,
		"cipher": "ECDHE-RSA-AES128-GCM-SHA256",
	}, "")

	diags := resourceBigipLtmCipherRuleCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "create failed")
}

func TestUnitResourceBigipLtmCipherRuleRead(t *testing.T) {
	name := "/Common/test-cipher-rule"
	mangled := "/mgmt/tm/ltm/cipher/rule/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{"name":"test-cipher-rule","fullPath":"%s","cipher":"DEFAULT","dhGroups":"DEFAULT","signatureAlgorithms":"DEFAULT"}`, name)
	})

	r := resourceBigipLtmCipherRule()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":   name,
		"cipher": "DEFAULT",
	}, name)

	diags := resourceBigipLtmCipherRuleRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "DEFAULT", d.Get("cipher").(string))
}

// A 404 means the cipher rule no longer exists on the device; Read clears
// the resource ID rather than returning an error.
func TestUnitResourceBigipLtmCipherRuleRead_NotFound(t *testing.T) {
	name := "/Common/test-cipher-rule"
	mangled := "/mgmt/tm/ltm/cipher/rule/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})

	r := resourceBigipLtmCipherRule()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":   name,
		"cipher": "DEFAULT",
	}, name)

	diags := resourceBigipLtmCipherRuleRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Equal(t, "", d.Id())
}

func TestUnitResourceBigipLtmCipherRuleUpdate(t *testing.T) {
	name := "/Common/test-cipher-rule"
	mangled := "/mgmt/tm/ltm/cipher/rule/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			_, _ = fmt.Fprint(w, `{}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-cipher-rule","fullPath":"%s","cipher":"DEFAULT:!SSLv3","dhGroups":"DEFAULT","signatureAlgorithms":"DEFAULT"}`, name)
	})

	r := resourceBigipLtmCipherRule()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":   name,
		"cipher": "DEFAULT:!SSLv3",
	}, name)

	diags := resourceBigipLtmCipherRuleUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	require.Equal(t, "DEFAULT:!SSLv3", d.Get("cipher").(string))
}

func TestUnitResourceBigipLtmCipherRuleUpdate_Error(t *testing.T) {
	name := "/Common/test-cipher-rule"
	mangled := "/mgmt/tm/ltm/cipher/rule/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})

	r := resourceBigipLtmCipherRule()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":   name,
		"cipher": "DEFAULT",
	}, name)

	diags := resourceBigipLtmCipherRuleUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "update failed")
}

func TestUnitResourceBigipLtmCipherRuleDelete(t *testing.T) {
	name := "/Common/test-cipher-rule"
	mangled := "/mgmt/tm/ltm/cipher/rule/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
	})

	r := resourceBigipLtmCipherRule()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":   name,
		"cipher": "DEFAULT",
	}, name)

	diags := resourceBigipLtmCipherRuleDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitResourceBigipLtmCipherRuleDelete_Error(t *testing.T) {
	name := "/Common/test-cipher-rule"
	mangled := "/mgmt/tm/ltm/cipher/rule/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})

	r := resourceBigipLtmCipherRule()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":   name,
		"cipher": "DEFAULT",
	}, name)

	diags := resourceBigipLtmCipherRuleDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "delete failed")
}

func TestUnitResourceBigipLtmCipherRuleSchema(t *testing.T) {
	r := resourceBigipLtmCipherRule()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.CreateContext)
	require.NotNil(t, r.ReadContext)
	require.NotNil(t, r.UpdateContext)
	require.NotNil(t, r.DeleteContext)
}
