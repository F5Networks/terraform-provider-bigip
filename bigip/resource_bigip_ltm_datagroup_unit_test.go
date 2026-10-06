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
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema
// ---------------------------------------------------------------------

func TestUnitResourceBigipLtmDataGroupSchema(t *testing.T) {
	r := resourceBigipLtmDataGroup()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.CreateContext)
	require.NotNil(t, r.ReadContext)
	require.NotNil(t, r.UpdateContext)
	require.NotNil(t, r.DeleteContext)
}

// ---------------------------------------------------------------------
// Create - internal data group
// ---------------------------------------------------------------------

func TestUnitResourceBigipLtmDataGroupCreate_Internal(t *testing.T) {
	name := "/Common/test-dg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/ltm/data-group/internal", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{"name":"test-dg","fullPath":"%s","type":"string","records":[{"name":"key1","data":"value1"}]}`, name)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": true,
		"record": []interface{}{
			map[string]interface{}{"name": "key1", "data": "value1"},
		},
	}, "")

	diags := resourceBigipLtmDataGroupCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())
}

func TestUnitResourceBigipLtmDataGroupCreate_InternalNoRecords(t *testing.T) {
	name := "/Common/test-dg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/ltm/data-group/internal", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"test-dg","fullPath":"%s","type":"string"}`, name)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": true,
	}, "")

	diags := resourceBigipLtmDataGroupCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())
}

func TestUnitResourceBigipLtmDataGroupCreate_InternalError(t *testing.T) {
	name := "/Common/test-dg"
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/ltm/data-group/internal", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": true,
	}, "")

	diags := resourceBigipLtmDataGroupCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "create failed")
}

// ---------------------------------------------------------------------
// Create - external data group (file-based)
// ---------------------------------------------------------------------

func TestUnitResourceBigipLtmDataGroupCreate_External(t *testing.T) {
	name := "/Common/test-dg-ext"
	mangledInternal := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)
	mangledExternal := "/mgmt/tm/ltm/data-group/external/" + MangleFullPath(name)

	tmpFile, err := os.CreateTemp(t.TempDir(), "test-dg-*.txt")
	require.NoError(t, err)
	_, err = tmpFile.WriteString("key1 := value1\n")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/data-group", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/data-group/external", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc(mangledInternal, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	mux.HandleFunc(mangledExternal, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"test-dg-ext","fullPath":"%s","type":"string"}`, name)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":        name,
		"type":        "string",
		"internal":    false,
		"records_src": tmpFile.Name(),
	}, "")

	diags := resourceBigipLtmDataGroupCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())
}

func TestUnitResourceBigipLtmDataGroupCreate_ExternalFileError(t *testing.T) {
	name := "/Common/test-dg-ext"
	_, client := NewUnitTestServer(t)

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":        name,
		"type":        "string",
		"internal":    false,
		"records_src": "/nonexistent/path/to/file.txt",
	}, "")

	diags := resourceBigipLtmDataGroupCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "error in reading file")
}

// ---------------------------------------------------------------------
// Read
// ---------------------------------------------------------------------

func TestUnitResourceBigipLtmDataGroupRead_Internal(t *testing.T) {
	name := "/Common/test-dg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"test-dg","fullPath":"%s","type":"string","records":[{"name":"key1","data":"value1"}]}`, name)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"type": "string",
	}, name)

	diags := resourceBigipLtmDataGroupRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "string", d.Get("type").(string))
}

func TestUnitResourceBigipLtmDataGroupRead_InternalError(t *testing.T) {
	name := "/Common/test-dg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"type": "string",
	}, name)

	diags := resourceBigipLtmDataGroupRead(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "read failed")
}

func TestUnitResourceBigipLtmDataGroupRead_FallbackToExternal(t *testing.T) {
	name := "/Common/test-dg-ext"
	mangledInternal := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)
	mangledExternal := "/mgmt/tm/ltm/data-group/external/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangledInternal, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	mux.HandleFunc(mangledExternal, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"test-dg-ext","fullPath":"%s","type":"string"}`, name)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"type": "string",
	}, name)

	diags := resourceBigipLtmDataGroupRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "string", d.Get("type").(string))
}

func TestUnitResourceBigipLtmDataGroupRead_ExternalError(t *testing.T) {
	name := "/Common/test-dg-ext"
	mangledInternal := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)
	mangledExternal := "/mgmt/tm/ltm/data-group/external/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangledInternal, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	mux.HandleFunc(mangledExternal, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"external read failed"}`)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
		"type": "string",
	}, name)

	diags := resourceBigipLtmDataGroupRead(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "external read failed")
}

// ---------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------

func TestUnitResourceBigipLtmDataGroupUpdate_Internal(t *testing.T) {
	name := "/Common/test-dg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"entries":{"https://localhost/mgmt/tm/cli/version/0":{"nestedStats":{"entries":{"active":{"description":"BIG-IP_v17.1.0"}}}}}}`)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			_, _ = fmt.Fprint(w, `{}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-dg","fullPath":"%s","type":"string","records":[{"name":"key2","data":"value2"}]}`, name)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": true,
		"record": []interface{}{
			map[string]interface{}{"name": "key2", "data": "value2"},
		},
	}, name)

	diags := resourceBigipLtmDataGroupUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitResourceBigipLtmDataGroupUpdate_Internal12x(t *testing.T) {
	name := "/Common/test-dg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"entries":{"https://localhost/mgmt/tm/cli/version/0":{"nestedStats":{"entries":{"active":{"description":"12.1.2"}}}}}}`)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			_, _ = fmt.Fprint(w, `{}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-dg","fullPath":"%s","type":"string"}`, name)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": true,
	}, name)

	diags := resourceBigipLtmDataGroupUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitResourceBigipLtmDataGroupUpdate_VersionError(t *testing.T) {
	name := "/Common/test-dg"
	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"version failed"}`)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": true,
	}, name)

	diags := resourceBigipLtmDataGroupUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "version failed")
}

func TestUnitResourceBigipLtmDataGroupUpdate_ModifyError(t *testing.T) {
	name := "/Common/test-dg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/cli/version", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"entries":{"https://localhost/mgmt/tm/cli/version/0":{"nestedStats":{"entries":{"active":{"description":"BIG-IP_v17.1.0"}}}}}}`)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, `{"code":500,"message":"modify failed"}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-dg","fullPath":"%s","type":"string"}`, name)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": true,
	}, name)

	diags := resourceBigipLtmDataGroupUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "modify failed")
}

func TestUnitResourceBigipLtmDataGroupUpdate_External(t *testing.T) {
	name := "/Common/test-dg-ext"
	mangledInternal := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)
	mangledExternal := "/mgmt/tm/ltm/data-group/external/" + MangleFullPath(name)

	tmpFile, err := os.CreateTemp(t.TempDir(), "test-dg-*.txt")
	require.NoError(t, err)
	_, err = tmpFile.WriteString("key1 := value1\n")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/shared/file-transfer/uploads/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/file/data-group", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/data-group/external", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{}`)
	})
	mux.HandleFunc(mangledInternal, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"not found"}`)
	})
	mux.HandleFunc(mangledExternal, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"test-dg-ext","fullPath":"%s","type":"string"}`, name)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":        name,
		"type":        "string",
		"internal":    false,
		"records_src": tmpFile.Name(),
	}, name)

	diags := resourceBigipLtmDataGroupUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)
}

func TestUnitResourceBigipLtmDataGroupUpdate_ExternalFileError(t *testing.T) {
	name := "/Common/test-dg-ext"
	_, client := NewUnitTestServer(t)

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":        name,
		"type":        "string",
		"internal":    false,
		"records_src": "/nonexistent/path/to/file.txt",
	}, name)

	diags := resourceBigipLtmDataGroupUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "error in reading file")
}

// ---------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------

func TestUnitResourceBigipLtmDataGroupDelete_Internal(t *testing.T) {
	name := "/Common/test-dg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": true,
	}, name)

	diags := resourceBigipLtmDataGroupDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitResourceBigipLtmDataGroupDelete_InternalError(t *testing.T) {
	name := "/Common/test-dg"
	mangled := "/mgmt/tm/ltm/data-group/internal/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": true,
	}, name)

	diags := resourceBigipLtmDataGroupDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "delete failed")
}

func TestUnitResourceBigipLtmDataGroupDelete_External(t *testing.T) {
	name := "/Common/test-dg-ext"
	mangledExternal := "/mgmt/tm/ltm/data-group/external/" + MangleFullPath(name)
	mangledFile := "/mgmt/tm/sys/file/data-group/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangledExternal, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc(mangledFile, func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": false,
	}, name)

	diags := resourceBigipLtmDataGroupDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitResourceBigipLtmDataGroupDelete_ExternalDGError(t *testing.T) {
	name := "/Common/test-dg-ext"
	mangledExternal := "/mgmt/tm/ltm/data-group/external/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangledExternal, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete dg failed"}`)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": false,
	}, name)

	diags := resourceBigipLtmDataGroupDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "delete dg failed")
}

func TestUnitResourceBigipLtmDataGroupDelete_ExternalFileError(t *testing.T) {
	name := "/Common/test-dg-ext"
	mangledExternal := "/mgmt/tm/ltm/data-group/external/" + MangleFullPath(name)
	mangledFile := "/mgmt/tm/sys/file/data-group/" + MangleFullPath(name)

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc(mangledExternal, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc(mangledFile, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete file failed"}`)
	})

	r := resourceBigipLtmDataGroup()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":     name,
		"type":     "string",
		"internal": false,
	}, name)

	diags := resourceBigipLtmDataGroupDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "delete file failed")
}
