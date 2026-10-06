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

func TestDataSourceBigipSysNtpSchema(t *testing.T) {
	r := dataSourceBigipSysNtp()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.ReadContext)

	for _, field := range []string{"description", "ntp_servers", "timezone"} {
		s, ok := r.Schema[field]
		require.True(t, ok, "expected field '%s' to exist in schema", field)
		assert.True(t, s.Computed, "expected field '%s' to be computed", field)
	}
}

func TestDataSourceBigipSysNtpReadSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/ntp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"description":"test ntp","servers":["0.pool.ntp.org","1.pool.ntp.org"],"timezone":"America/Los_Angeles"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysNtp(), map[string]interface{}{})

	diags := dataSourceBigipSysNtpRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "test ntp", d.Get("description"))
	assert.Equal(t, "America/Los_Angeles", d.Get("timezone"))

	ntpServers := d.Get("ntp_servers").([]interface{})
	require.Len(t, ntpServers, 2)
	assert.Equal(t, "0.pool.ntp.org", ntpServers[0])
	assert.Equal(t, "1.pool.ntp.org", ntpServers[1])

	assert.Equal(t, "test ntp", d.Id())
}

// TestDataSourceBigipSysNtpReadEmpty covers a response with no description
// (BIG-IP omits the field entirely when unset, rather than returning an
// empty string). Read must not set an empty ID in this case -- it falls
// back to hashing the struct instead of using ntp.Description directly.
func TestDataSourceBigipSysNtpReadEmpty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/ntp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"timezone":"America/Los_Angeles"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysNtp(), map[string]interface{}{})

	diags := dataSourceBigipSysNtpRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Empty(t, d.Get("ntp_servers").([]interface{}))
	assert.Equal(t, "America/Los_Angeles", d.Get("timezone"))
	assert.NotEmpty(t, d.Id(), "ID must not be empty even when description is unset")
}

func TestDataSourceBigipSysNtpReadError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/ntp", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := newDatasourceTestClient(server.URL)

	d := newDatasourceTestResourceData(t, dataSourceBigipSysNtp(), map[string]interface{}{})

	diags := dataSourceBigipSysNtpRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error when the NTP lookup call fails")
	assert.Contains(t, diags[0].Summary, "error retrieving system NTP configuration")
}
