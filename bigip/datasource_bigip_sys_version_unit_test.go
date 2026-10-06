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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataSourceBigipSysVersionSchema(t *testing.T) {
	r := dataSourceBigipSysVersion()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	for _, field := range []string{"version", "build", "product", "edition", "platform"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "expected field '%s' to exist in schema", field)
		assert.True(t, s.Computed, "expected field '%s' to be computed", field)
	}
}

const sysVersionResponseJSON = `{
	"kind": "tm:sys:version:versionstats",
	"selfLink": "https://localhost/mgmt/tm/sys/version?ver=17.5.1",
	"entries": {
		"https://localhost/mgmt/tm/sys/version/0": {
			"nestedStats": {
				"entries": {
					"Build": {"description": "0.0.7"},
					"Date": {"description": "Wed Jun 18 06:03:07 PDT 2025"},
					"Edition": {"description": "Final"},
					"Product": {"description": "BIG-IP"},
					"Title": {"description": "Main Package"},
					"Version": {"description": "17.5.1"}
				}
			}
		}
	}
}`

const sysHardwareResponseJSON = `{
	"kind": "tm:sys:hardware:hardwarestats",
	"selfLink": "https://localhost/mgmt/tm/sys/hardware?ver=17.5.1",
	"entries": {
		"https://localhost/mgmt/tm/sys/hardware/platform": {
			"nestedStats": {
				"entries": {
					"https://localhost/mgmt/tm/sys/hardware/platform/0": {
						"nestedStats": {
							"entries": {
								"baseMac": {"description": "fa:16:3e:b7:69:9d"},
								"biosRev": {"description": " "},
								"cloudName": {"description": " "},
								"hypervisorName": {"description": "OpenStack Nova"},
								"marketingName": {"description": "BIG-IP Virtual Edition"}
							}
						}
					}
				}
			}
		}
	}
}`

func TestDataSourceBigipSysVersionReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, sysVersionResponseJSON)
	})
	mux.HandleFunc("/mgmt/tm/sys/hardware", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, sysHardwareResponseJSON)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysVersion(), map[string]interface{}{})

	diags := dataSourceBigipSysVersionRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "17.5.1", d.Get("version"))
	assert.Equal(t, "0.0.7", d.Get("build"))
	assert.Equal(t, "BIG-IP", d.Get("product"))
	assert.Equal(t, "Final", d.Get("edition"))
	assert.Equal(t, "BIG-IP Virtual Edition", d.Get("platform"))

	assert.NotEmpty(t, d.Id())
}

// TestDataSourceBigipSysVersionReadHardwareUnavailable covers sys/hardware
// being unreachable/erroring: platform should come back empty but the read
// overall must still succeed, since platform is supplemental information
// layered on top of the primary sys/version read. The failure must not be
// swallowed silently -- it surfaces as a warning diagnostic so it isn't
// opaque to the practitioner (or to debugging) when platform unexpectedly
// comes back empty.
func TestDataSourceBigipSysVersionReadHardwareUnavailable(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, sysVersionResponseJSON)
	})
	mux.HandleFunc("/mgmt/tm/sys/hardware", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysVersion(), map[string]interface{}{})

	diags := dataSourceBigipSysVersionRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "sys/hardware failure must not fail the overall read")
	assert.Equal(t, "17.5.1", d.Get("version"))
	assert.Equal(t, "", d.Get("platform"))

	require.Len(t, diags, 1, "expected exactly one warning diagnostic about the platform lookup")
	assert.Equal(t, diag.Warning, diags[0].Severity)
	assert.Contains(t, diags[0].Detail, "error retrieving sys/hardware")
}

// TestDataSourceBigipSysVersionReadHardwareMissingPlatformEntry covers a
// sys/hardware response that succeeds but has no "platform" section at all
// -- platformMarketingName must not panic (Go map reads on a nil/missing
// key return the zero value, not a panic) and must return a reason so the
// caller can surface a diagnostic instead of a silent empty string.
func TestDataSourceBigipSysVersionReadHardwareMissingPlatformEntry(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, sysVersionResponseJSON)
	})
	mux.HandleFunc("/mgmt/tm/sys/hardware", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"entries":{}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysVersion(), map[string]interface{}{})

	diags := dataSourceBigipSysVersionRead(context.Background(), d, client)

	require.False(t, diags.HasError())
	assert.Equal(t, "", d.Get("platform"))
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Detail, "no \"platform\" entry")
}

// TestDataSourceBigipSysVersionReadHardwareMissingMarketingName covers a
// sys/hardware response with a "platform" section that is missing the
// "marketingName" field specifically (e.g. an unexpected/older field set),
// exercising the deepest level of the nested-map defensive checks.
func TestDataSourceBigipSysVersionReadHardwareMissingMarketingName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, sysVersionResponseJSON)
	})
	mux.HandleFunc("/mgmt/tm/sys/hardware", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{
			"entries": {
				"https://localhost/mgmt/tm/sys/hardware/platform": {
					"nestedStats": {
						"entries": {
							"https://localhost/mgmt/tm/sys/hardware/platform/0": {
								"nestedStats": {
									"entries": {
										"baseMac": {"description": "fa:16:3e:b7:69:9d"}
									}
								}
							}
						}
					}
				}
			}
		}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysVersion(), map[string]interface{}{})

	diags := dataSourceBigipSysVersionRead(context.Background(), d, client)

	require.False(t, diags.HasError())
	assert.Equal(t, "", d.Get("platform"))
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Detail, "no \"marketingName\" field")
}

func TestDataSourceBigipSysVersionReadEmptyEntries(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"entries":{}}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysVersion(), map[string]interface{}{})

	diags := dataSourceBigipSysVersionRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when sys/version has no entries")
	assert.Contains(t, diags[0].Summary, "system version response had no entries")
}

func TestDataSourceBigipSysVersionReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/version", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysVersion(), map[string]interface{}{})

	diags := dataSourceBigipSysVersionRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the sys/version lookup call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving system version")
}
