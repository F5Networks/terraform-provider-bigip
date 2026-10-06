/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceBigipGtmPoolSchema(t *testing.T) {
	r := resourceBigipGtmPool()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.CreateContext == nil || r.ReadContext == nil || r.UpdateContext == nil || r.DeleteContext == nil {
		t.Fatal("Expected all CRUD contexts to be defined")
	}
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	for _, field := range []string{"name", "type"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "Expected field '%s' to exist in schema", field)
		assert.True(t, s.Required)
		assert.True(t, s.ForceNew)
	}

	membersSchema, ok := r.Schema["members"]
	require.True(t, ok)
	assert.Equal(t, schema.TypeSet, membersSchema.Type)
}

func TestParseGtmPoolIDValid(t *testing.T) {
	parts := parseGtmPoolID("a:/Common/firstpool")
	require.NotNil(t, parts)
	assert.Equal(t, "a", parts["type"])
	assert.Equal(t, "Common", parts["partition"])
	assert.Equal(t, "firstpool", parts["name"])
}

func TestParseGtmPoolIDNoPartition(t *testing.T) {
	parts := parseGtmPoolID("a:justaname")
	require.NotNil(t, parts)
	assert.Equal(t, "Common", parts["partition"])
	assert.Equal(t, "justaname", parts["name"])
}

func TestParseGtmPoolIDInvalid(t *testing.T) {
	parts := parseGtmPoolID("no-colon-here")
	assert.Nil(t, parts)
}

func TestUnitGtmPoolCreateReadUpdateDelete(t *testing.T) {
	fullPath := "/Common/test-pool"
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	var updated bool
	var sawPost, sawPut, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/pool/a", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":"test-pool","partition":"Common","fullPath":"%s"}`, fullPath)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			sawPut = true
			updated = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		case http.MethodGet:
			if updated {
				_, _ = fmt.Fprintf(w, `{"name":"test-pool","partition":"Common","fullPath":"%s","loadBalancingMode":"topology","ttl":60}`, fullPath)
			} else {
				_, _ = fmt.Fprintf(w, `{"name":"test-pool","partition":"Common","fullPath":"%s"}`, fullPath)
			}
		case http.MethodDelete:
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})
	mux.HandleFunc(mangled+"/members", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[{"name":"server1:vs1","partition":"Common","enabled":true,"ratio":1,"monitor":"default"}]}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "test-pool",
		"type":      "a",
		"partition": "Common",
	}, "")

	ctx := context.Background()

	diags := resourceBigipGtmPoolCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost)
	assert.Equal(t, "a:"+fullPath, d.Id())
	assert.True(t, sawPut, "Create calls Update internally")

	diags = resourceBigipGtmPoolRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "test-pool", d.Get("name"))
	assert.Equal(t, "topology", d.Get("load_balancing_mode"))
	assert.Equal(t, 60, d.Get("ttl"))
	members := d.Get("members").(*schema.Set)
	require.Equal(t, 1, members.Len())
	member := members.List()[0].(map[string]interface{})
	assert.Equal(t, "/Common/server1:vs1", member["name"])

	require.NoError(t, d.Set("ttl", 90))
	diags = resourceBigipGtmPoolUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipGtmPoolDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete)
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmPoolReadMemberNameWithSubPath(t *testing.T) {
	fullPath := "/Common/subpath-pool"
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"subpath-pool","partition":"Common","fullPath":"%s"}`, fullPath)
	})
	mux.HandleFunc(mangled+"/members", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[{"name":"member1","partition":"Common","subPath":"sub","enabled":true},{"name":"member2","enabled":true}]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "subpath-pool",
		"type":      "a",
		"partition": "Common",
	}, "a:"+fullPath)

	diags := resourceBigipGtmPoolRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	members := d.Get("members").(*schema.Set)
	require.Equal(t, 2, members.Len())

	var names []string
	for _, m := range members.List() {
		names = append(names, m.(map[string]interface{})["name"].(string))
	}
	assert.Contains(t, names, "/Common/sub/member1")
	assert.Contains(t, names, "member2")
}

func TestUnitGtmPoolReadFromImportID(t *testing.T) {
	// Exercises the ID-parsing fallback in Read: name/type/partition are
	// empty in state (as after schema.TestResourceDataRaw with no values),
	// so Read must derive them from d.Id().
	fullPath := "/Common/id-parsed-pool"
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"id-parsed-pool","partition":"Common","fullPath":"%s"}`, fullPath)
	})
	mux.HandleFunc(mangled+"/members", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "a:"+fullPath)

	diags := resourceBigipGtmPoolRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "id-parsed-pool", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))
}

func TestUnitGtmPoolCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/pool/a", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprint(w, `{"code":409,"message":"the requested object already exists"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "dup-pool",
		"type":      "a",
		"partition": "Common",
	}, "")

	diags := resourceBigipGtmPoolCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmPoolReadError(t *testing.T) {
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath("/Common/broken-pool")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "broken-pool",
		"type":      "a",
		"partition": "Common",
	}, "a:/Common/broken-pool")

	diags := resourceBigipGtmPoolRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmPoolReadNotFound(t *testing.T) {
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath("/Common/missing-pool")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "missing-pool",
		"type":      "a",
		"partition": "Common",
	}, "a:/Common/missing-pool")

	// A 404 means the pool no longer exists on the device; Read clears the
	// resource ID rather than returning an error.
	diags := resourceBigipGtmPoolRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Equal(t, "", d.Id())
}

func TestUnitGtmPoolUpdateWithMembers(t *testing.T) {
	fullPath := "/Common/members-pool"
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath(fullPath)

	mux := http.NewServeMux()
	var receivedBody string
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			body, _ := io.ReadAll(r.Body)
			receivedBody = string(body)
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"members-pool","partition":"Common","fullPath":"%s"}`, fullPath)
	})
	mux.HandleFunc(mangled+"/members", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "members-pool",
		"type":      "a",
		"partition": "Common",
		"members": []interface{}{
			map[string]interface{}{"name": "server1:vs1", "enabled": true, "ratio": 2},
		},
	}, "a:"+fullPath)

	diags := resourceBigipGtmPoolUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.Contains(t, receivedBody, "server1:vs1")
}

func TestUnitGtmPoolUpdateError(t *testing.T) {
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath("/Common/err-pool")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "err-pool",
		"type":      "a",
		"partition": "Common",
	}, "a:/Common/err-pool")

	diags := resourceBigipGtmPoolUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmPoolDeleteError(t *testing.T) {
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath("/Common/del-err-pool")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "del-err-pool",
		"type":      "a",
		"partition": "Common",
	}, "a:/Common/del-err-pool")

	diags := resourceBigipGtmPoolDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmPoolImport(t *testing.T) {
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath("/Common/import-pool")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"import-pool","partition":"Common","fullPath":"/Common/import-pool"}`)
	})
	mux.HandleFunc(mangled+"/members", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "a:/Common/import-pool")

	results, err := resourceBigipGtmPoolImport(context.Background(), d, client)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "import-pool", results[0].Get("name"))
	assert.Equal(t, "a", results[0].Get("type"))
}

func TestUnitGtmPoolImportInvalid(t *testing.T) {
	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "not-a-valid-id")

	_, err := resourceBigipGtmPoolImport(context.Background(), d, nil)
	require.Error(t, err)
}

func TestUnitGtmPoolImportReadError(t *testing.T) {
	mangled := "/mgmt/tm/gtm/pool/a/" + MangleFullPath("/Common/import-err-pool")

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmPool()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "a:/Common/import-err-pool")

	_, err := resourceBigipGtmPoolImport(context.Background(), d, client)
	require.Error(t, err)
}
