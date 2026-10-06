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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestDataSourceBigipLtmDataGroupSchema(t *testing.T) {
	r := dataSourceBigipLtmDataGroup()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	for _, field := range []string{"name", "partition"} {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}
}

// ---------------------------------------------------------------------
// Direct Read-function unit tests
// ---------------------------------------------------------------------

func TestUnitLtmDataGroupRead(t *testing.T) {
	name := "/Common/testdg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"testdg","partition":"Common","fullPath":"%s","type":"string","records":[{"name":"key1","data":"value1"},{"name":"key2","data":"value2"}]}`, name)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := dataSourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testdg",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipLtmDataGroupRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, name, d.Id())
	require.Equal(t, "string", d.Get("type").(string))
	records := d.Get("record").(*schema.Set)
	require.Equal(t, 2, records.Len())
}

// A 404 from GetInternalDataGroup means the data group does not exist on the
// device. dataSourceBigipLtmDataGroupRead treats that as "not found" (no
// error, empty ID) rather than a hard failure.
func TestUnitLtmDataGroupRead_NotFound(t *testing.T) {
	name := "/Common/testdg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := dataSourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testdg",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipLtmDataGroupRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	require.Empty(t, d.Id())
}

func TestUnitLtmDataGroupRead_Error(t *testing.T) {
	name := "/Common/testdg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := dataSourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testdg",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipLtmDataGroupRead(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "read failed")
}

func TestUnitLtmDataGroupRead_NoRecords(t *testing.T) {
	name := "/Common/testdg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"testdg","partition":"Common","fullPath":"%s","type":"string"}`, name)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := dataSourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":      "testdg",
		"partition": "Common",
	}, "")

	diags := dataSourceBigipLtmDataGroupRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, name, d.Id())
	records := d.Get("record").(*schema.Set)
	require.Equal(t, 0, records.Len(), "expected no records when the API response omits the records field")
}
