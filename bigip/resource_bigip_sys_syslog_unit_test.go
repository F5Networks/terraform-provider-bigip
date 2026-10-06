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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceBigipSysSyslogSchema(t *testing.T) {
	r := resourceBigipSysSyslog()
	require.NotNil(t, r.Schema)
	require.NotNil(t, r.CreateContext)
	require.NotNil(t, r.ReadContext)
	require.NotNil(t, r.UpdateContext)
	require.NotNil(t, r.DeleteContext)
	require.NotNil(t, r.Importer)
	_, ok := r.Schema["remote_servers"]
	require.True(t, ok, "expected 'remote_servers' to exist in schema")
}

func TestUnitSysSyslogCreateReadUpdateDelete(t *testing.T) {
	var patchedBody string
	var patchCount int

	mux, client := NewUnitTestServer(t)
	mux.HandleFunc("/mgmt/tm/sys/syslog", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			patchCount++
			b, _ := io.ReadAll(r.Body)
			patchedBody = string(b)
			_, _ = fmt.Fprint(w, `{"authPrivFrom":"notice","authPrivTo":"emerg","consoleLog":"enabled","remoteServers":[{"name":"/Common/remote1","host":"10.1.1.1","localIp":"none","remotePort":514}]}`)
		case http.MethodGet:
			_, _ = fmt.Fprint(w, `{"authPrivFrom":"notice","authPrivTo":"emerg","consoleLog":"enabled","remoteServers":[{"name":"/Common/remote1","host":"10.1.1.1","localIp":"none","remotePort":514}]}`)
		}
	})

	r := resourceBigipSysSyslog()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"auth_priv_from": "notice",
		"auth_priv_to":   "emerg",
		"console_log":    "enabled",
		"remote_servers": []interface{}{
			map[string]interface{}{"name": "remote1", "host": "10.1.1.1", "remote_port": 514},
		},
	}, "")

	ctx := context.Background()

	diags := resourceBigipSysSyslogCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	assert.Contains(t, patchedBody, "notice")
	assert.Equal(t, "syslog", d.Id())

	diags = resourceBigipSysSyslogRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	remoteServers := d.Get("remote_servers").([]interface{})
	require.Len(t, remoteServers, 1)
	rs := remoteServers[0].(map[string]interface{})
	assert.Equal(t, "remote1", rs["name"], "expected full-path prefix to be stripped")
	assert.Equal(t, "", rs["local_ip"], "expected BIG-IP's \"none\" placeholder to normalize to empty string")

	diags = resourceBigipSysSyslogUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipSysSyslogDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	assert.Equal(t, 3, patchCount, "expected create+update+delete to each PATCH sys/syslog")
	assert.Equal(t, "", d.Id())
}

func TestUnitSysSyslogReadNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/sys/syslog", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":404,"message":"sys/syslog not found"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client := NewUnitTestClient(server.URL)

	r := resourceBigipSysSyslog()
	d := NewTestResourceData(t, r, map[string]interface{}{}, "syslog")

	diags := resourceBigipSysSyslogRead(context.Background(), d, client)
	require.False(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}
