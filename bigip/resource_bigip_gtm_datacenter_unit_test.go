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

func TestResourceBigipGtmDatacenterSchema(t *testing.T) {
	r := resourceBigipGtmDatacenter()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.CreateContext == nil {
		t.Fatal("Expected CreateContext to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}
	if r.UpdateContext == nil {
		t.Fatal("Expected UpdateContext to be defined")
	}
	if r.DeleteContext == nil {
		t.Fatal("Expected DeleteContext to be defined")
	}
	if r.Importer == nil {
		t.Fatal("Expected Importer to be defined")
	}

	nameSchema, ok := r.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist in schema")
	}
	if !nameSchema.Required {
		t.Error("Expected field 'name' to be required")
	}
	if !nameSchema.ForceNew {
		t.Error("Expected field 'name' to be ForceNew")
	}

	partitionSchema, ok := r.Schema["partition"]
	if !ok {
		t.Fatal("Expected field 'partition' to exist in schema")
	}
	if partitionSchema.Default != "Common" {
		t.Errorf("Expected field 'partition' default to be 'Common', got %v", partitionSchema.Default)
	}
}

func TestResourceBigipGtmDatacenterProberFallbackValidation(t *testing.T) {
	r := resourceBigipGtmDatacenter()
	s := r.Schema["prober_fallback"]
	require.NotNil(t, s.ValidateFunc)

	for _, v := range []string{"any-available", "inside-datacenter", "outside-datacenter", "inherit", "pool"} {
		_, errs := s.ValidateFunc(v, "prober_fallback")
		assert.Empty(t, errs, "expected %q to be valid", v)
	}

	_, errs := s.ValidateFunc("bogus", "prober_fallback")
	assert.NotEmpty(t, errs, "expected 'bogus' to be invalid")
}

func TestResourceBigipGtmDatacenterProberPreferenceValidation(t *testing.T) {
	r := resourceBigipGtmDatacenter()
	s := r.Schema["prober_preference"]
	require.NotNil(t, s.ValidateFunc)

	for _, v := range []string{"inside-datacenter", "outside-datacenter", "inherit", "pool"} {
		_, errs := s.ValidateFunc(v, "prober_preference")
		assert.Empty(t, errs, "expected %q to be valid", v)
	}

	_, errs := s.ValidateFunc("bogus", "prober_preference")
	assert.NotEmpty(t, errs, "expected 'bogus' to be invalid")
}

func TestUnitGtmDatacenterCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-dc"
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath(name)

	mux := http.NewServeMux()
	var updated bool
	var sawPost, sawPut, sawDelete bool

	mux.HandleFunc("/mgmt/tm/gtm/datacenter", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		sawPost = true
		_, _ = fmt.Fprintf(w, `{"name":"test-dc","partition":"Common","fullPath":"%s","enabled":true}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			sawPut = true
			updated = true
			fallthrough
		case http.MethodGet:
			if updated {
				_, _ = fmt.Fprintf(w, `{"name":"test-dc","partition":"Common","fullPath":"%s","contact":"jane","description":"updated dc","enabled":true,"location":"NA","proberFallback":"pool","proberPreference":"pool"}`, name)
			} else {
				_, _ = fmt.Fprintf(w, `{"name":"test-dc","partition":"Common","fullPath":"%s","enabled":true}`, name)
			}
		case http.MethodDelete:
			sawDelete = true
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "test-dc",
		"partition": "Common",
		"enabled":   true,
	}, "")

	ctx := context.Background()

	diags := resourceBigipGtmDatacenterCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.True(t, sawPost, "expected a POST request during create")
	assert.Equal(t, name, d.Id())

	diags = resourceBigipGtmDatacenterRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "test-dc", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))

	require.NoError(t, d.Set("contact", "jane"))
	require.NoError(t, d.Set("description", "updated dc"))
	require.NoError(t, d.Set("location", "NA"))
	require.NoError(t, d.Set("prober_fallback", "pool"))
	require.NoError(t, d.Set("prober_preference", "pool"))

	diags = resourceBigipGtmDatacenterUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
	assert.True(t, sawPut, "expected a PUT request during update")
	assert.Equal(t, "jane", d.Get("contact"))
	assert.Equal(t, "updated dc", d.Get("description"))
	assert.Equal(t, "pool", d.Get("prober_fallback"))

	diags = resourceBigipGtmDatacenterDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.True(t, sawDelete, "expected a DELETE request during delete")
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmDatacenterCreateDisabled(t *testing.T) {
	name := "/Common/test-dc-disabled"
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/datacenter", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"test-dc-disabled","partition":"Common","fullPath":"%s"}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"test-dc-disabled","partition":"Common","fullPath":"%s","disabled":true}`, name)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "test-dc-disabled",
		"partition": "Common",
		"enabled":   false,
	}, "")

	ctx := context.Background()
	diags := resourceBigipGtmDatacenterCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Equal(t, false, d.Get("enabled"))
}

func TestUnitGtmDatacenterCreateError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/gtm/datacenter", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprint(w, `{"code":409,"message":"the requested object already exists"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "dup-dc",
		"partition": "Common",
	}, "")

	diags := resourceBigipGtmDatacenterCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmDatacenterReadNotFound(t *testing.T) {
	// Proof-of-concept use of NewUnitTestServer: replaces the usual
	// mux := http.NewServeMux() / server := httptest.NewServer(mux) /
	// defer server.Close() / client := NewUnitTestClient(server.URL)
	// four-liner repeated across the GTM unit test files with a single call.
	name := "/Common/missing-dc"
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})

	r := resourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{"name": "missing-dc", "partition": "Common"}, name)

	// A 404 means the datacenter no longer exists on the device; Read
	// clears the resource ID rather than returning an error.
	diags := resourceBigipGtmDatacenterRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}

func TestUnitGtmDatacenterUpdateError(t *testing.T) {
	name := "/Common/test-dc-err"
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{"name": "test-dc-err", "partition": "Common"}, name)
	require.NoError(t, d.Set("contact", "someone"))

	diags := resourceBigipGtmDatacenterUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmDatacenterDeleteError(t *testing.T) {
	name := "/Common/test-dc-del-err"
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"internal error"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{"name": "test-dc-del-err", "partition": "Common"}, name)

	diags := resourceBigipGtmDatacenterDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitGtmDatacenterReadFullPathFallback(t *testing.T) {
	// Exercises the "else" branch of the fullPath-parsing logic in Read,
	// which only triggers when the ID doesn't split into exactly
	// partition+name (e.g. a bare name with no leading slash at all).
	name := "bare-name-dc"
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"name":"bare-name-dc","partition":"Common","contact":"c","description":"d","location":"l","proberFallback":"any-available","proberPreference":"inside-datacenter","enabled":true}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{"name": name}, name)

	diags := resourceBigipGtmDatacenterRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "bare-name-dc", d.Get("name"))
	assert.Equal(t, "Common", d.Get("partition"))
}

func TestUnitGtmDatacenterImport(t *testing.T) {
	// bigip_gtm_datacenter uses the plain schema.ImportStatePassthroughContext
	// importer (no custom ID-parsing logic to exercise, unlike e.g. the pool
	// and wideip resources' "type:fullPath" composite IDs), so the meaningful
	// behavior to verify is: the imported ID passes through unchanged, and a
	// subsequent Read using that ID populates state from the device response.
	name := "/Common/import-dc"
	mangled := "/mgmt/tm/gtm/datacenter/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"import-dc","partition":"Common","fullPath":"%s","contact":"jane","description":"imported dc","enabled":true,"location":"NA","proberFallback":"pool","proberPreference":"pool"}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipGtmDatacenter()
	d := NewTestResourceData(t, r, map[string]interface{}{}, name)

	results, err := r.Importer.StateContext(context.Background(), d, client)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, name, results[0].Id())

	diags := resourceBigipGtmDatacenterRead(context.Background(), results[0], client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	assert.Equal(t, "import-dc", results[0].Get("name"))
	assert.Equal(t, "Common", results[0].Get("partition"))
	assert.Equal(t, "jane", results[0].Get("contact"))
	assert.Equal(t, "pool", results[0].Get("prober_fallback"))
}
