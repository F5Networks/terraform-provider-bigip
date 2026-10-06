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

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func vsTestResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceBigipLtmVirtualServer().Schema, raw)
}

// ---------------------------------------------------------------------
// ltmVirtualServerAttrDefaults / getVirtualServerConfig (pure functions)
// ---------------------------------------------------------------------

func TestResourceBigipLtmVirtualServerAttrDefaults_IPv4(t *testing.T) {
	d := vsTestResourceData(t, map[string]interface{}{
		"name":        "/Common/test-vs",
		"destination": "10.10.10.10",
	})
	ltmVirtualServerAttrDefaults(d)
	assert.Equal(t, "255.255.255.255", d.Get("mask").(string))
	assert.Equal(t, "0.0.0.0/0", d.Get("source").(string))
}

func TestResourceBigipLtmVirtualServerAttrDefaults_IPv6(t *testing.T) {
	d := vsTestResourceData(t, map[string]interface{}{
		"name":        "/Common/test-vs",
		"destination": "2001:db8::1",
	})
	ltmVirtualServerAttrDefaults(d)
	assert.Equal(t, "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", d.Get("mask").(string))
	assert.Equal(t, "::/0", d.Get("source").(string))
}

func TestResourceBigipLtmVirtualServerAttrDefaults_AlreadySet(t *testing.T) {
	d := vsTestResourceData(t, map[string]interface{}{
		"name":        "/Common/test-vs",
		"destination": "10.10.10.10",
		"mask":        "255.255.255.0",
		"source":      "192.168.0.0/24",
	})
	ltmVirtualServerAttrDefaults(d)
	assert.Equal(t, "255.255.255.0", d.Get("mask").(string))
	assert.Equal(t, "192.168.0.0/24", d.Get("source").(string))
}

func TestResourceBigipLtmVirtualServerGetConfig_IPv4CIDRMask(t *testing.T) {
	d := vsTestResourceData(t, map[string]interface{}{
		"name":                           "/Common/test-vs",
		"destination":                    "10.10.10.10",
		"port":                           80,
		"mask":                           "24",
		"pool":                           "/Common/test-pool",
		"translate_port":                 "enabled",
		"translate_address":              "enabled",
		"source_port":                    "preserve",
		"firewall_enforced_policy":       "/Common/afm-policy",
		"source":                         "0.0.0.0/0",
		"profiles":                       []interface{}{"/Common/tcp"},
		"client_profiles":                []interface{}{"/Common/clientssl"},
		"server_profiles":                []interface{}{"/Common/serverssl"},
		"persistence_profiles":           []interface{}{"/Common/cookie", "/Common/source_addr"},
		"default_persistence_profile":    "/Common/cookie",
		"fallback_persistence_profile":   "/Common/dest_addr",
		"policies":                       []interface{}{"/Common/policy1"},
		"vlans":                          []interface{}{"/Common/external"},
		"irules":                         []interface{}{"/Common/my-rule"},
		"security_log_profiles":          []interface{}{"\"/Common/log-profile\""},
		"per_flow_request_access_policy": "/Common/policy",
		"description":                    "test description",
		"ip_protocol":                    "tcp",
		"trafficmatching_criteria":       "",
		"connection_limit":               100,
		"source_address_translation":     "snat",
		"snatpool":                       "/Common/my-snatpool",
		"vlans_enabled":                  true,
		"state":                          "enabled",
	})

	pss := &bigip.VirtualServer{Name: "test-vs"}
	config := getVirtualServerConfig(d, pss)

	assert.Equal(t, "/Common/test-vs", config.Name)
	assert.Equal(t, "/Common/test-pool", config.Pool)
	assert.Equal(t, "10.10.10.10:80", config.Destination)
	assert.Equal(t, "255.255.255.0", config.Mask)
	assert.Len(t, config.Profiles, 3)
	assert.Len(t, config.PersistenceProfiles, 2)
	assert.Equal(t, "/Common/dest_addr", config.FallbackPersistenceProfile)
	assert.Equal(t, []string{"/Common/policy1"}, config.Policies)
	assert.Equal(t, []string{"/Common/external"}, config.Vlans)
	assert.Equal(t, []string{"/Common/my-rule"}, config.Rules)
	assert.Equal(t, "snat", config.SourceAddressTranslation.Type)
	assert.Equal(t, "/Common/my-snatpool", config.SourceAddressTranslation.Pool)
	assert.True(t, config.VlansEnabled)
	assert.False(t, config.VlansDisabled)
	assert.True(t, config.Enabled)
	assert.False(t, config.Disabled)
	assert.Equal(t, 100, config.ConnectionLimit)
}

func TestResourceBigipLtmVirtualServerGetConfig_IPv6Destination(t *testing.T) {
	d := vsTestResourceData(t, map[string]interface{}{
		"name":              "/Common/test-vs6",
		"destination":       "2001:db8::1",
		"port":              443,
		"mask":              "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"source":            "::/0",
		"translate_port":    "disabled",
		"translate_address": "disabled",
		"vlans_enabled":     false,
		"state":             "disabled",
	})
	pss := &bigip.VirtualServer{Name: "test-vs6"}
	config := getVirtualServerConfig(d, pss)
	assert.Equal(t, "2001:db8::1.443", config.Destination)
	assert.Equal(t, "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", config.Mask)
	assert.True(t, config.VlansDisabled)
	assert.True(t, config.Disabled)
	assert.False(t, config.Enabled)
}

func TestResourceBigipLtmVirtualServerGetConfig_IruleChangeNoValue(t *testing.T) {
	d := vsTestResourceData(t, map[string]interface{}{
		"name":        "/Common/test-vs",
		"destination": "10.10.10.10",
		"port":        80,
		"mask":        "255.255.255.255",
	})
	pss := &bigip.VirtualServer{Name: "test-vs"}
	config := getVirtualServerConfig(d, pss)
	// no irules set and no change -> rules stays nil
	assert.Nil(t, config.Rules)
}

// ---------------------------------------------------------------------
// CRUD tests using an httptest mux/server local to this file
// ---------------------------------------------------------------------

func vsNewTestClient(serverURL string) *bigip.BigIP {
	client := newDatasourceTestClient(serverURL)
	// Skip TEEM telemetry reporting in Create, which otherwise dereferences
	// UserAgent segments that don't exist for this synthetic test client.
	client.Teem = true
	return client
}

func TestResourceBigipLtmVirtualServerCreate_Success(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"name":"test-vs"}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"test-vs","destination":"/Common/10.10.10.10:80","mask":"255.255.255.255","source":"0.0.0.0/0","pool":"/Common/test-pool","enabled":true}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs/profiles", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[]}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs/policies", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[]}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name":        "/Common/test-vs",
		"destination": "10.10.10.10",
		"port":        80,
		"pool":        "/Common/test-pool",
	})

	diags := resourceBigipLtmVirtualServerCreate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "/Common/test-vs", d.Id())
}

func TestResourceBigipLtmVirtualServerCreate_Error(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"create failed"}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name":        "/Common/test-vs",
		"destination": "10.10.10.10",
		"port":        80,
	})

	diags := resourceBigipLtmVirtualServerCreate(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmVirtualServerRead_Success(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
			"name":"test-vs",
			"destination":"/Common/10.10.10.10:80",
			"mask":"255.255.255.255",
			"source":"0.0.0.0/0",
			"pool":"/Common/test-pool",
			"enabled":true,
			"ipProtocol":"tcp",
			"connectionLimit":0,
			"translateAddress":"enabled",
			"translatePort":"enabled",
			"sourceAddressTranslation":{"type":"automap"}
		}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs/profiles", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"name":"tcp","context":"all","fullPath":"/Common/tcp"},{"name":"clientssl","context":"clientside","fullPath":"/Common/clientssl"},{"name":"serverssl","context":"serverside","fullPath":"/Common/serverssl"}]}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs/policies", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[]}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name": "/Common/test-vs",
	})
	d.SetId("/Common/test-vs")

	diags := resourceBigipLtmVirtualServerRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "/Common/test-pool", d.Get("pool").(string))
	assert.Equal(t, "enabled", d.Get("state").(string))
	assert.Equal(t, 80, d.Get("port").(int))
}

func TestResourceBigipLtmVirtualServerRead_NotFound(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, `{"code":404,"message":"not found"}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name": "/Common/test-vs",
	})
	d.SetId("/Common/test-vs")

	diags := resourceBigipLtmVirtualServerRead(context.Background(), d, client)
	// A 404 means the virtual server no longer exists on the device; Read
	// clears the resource ID rather than returning an error.
	assert.False(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmVirtualServerRead_Error(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"boom"}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name": "/Common/test-vs",
	})
	d.SetId("/Common/test-vs")

	diags := resourceBigipLtmVirtualServerRead(context.Background(), d, client)
	assert.True(t, diags.HasError())
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmVirtualServerUpdate_Success(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			_, _ = fmt.Fprintf(w, `{"name":"test-vs"}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-vs","destination":"/Common/10.10.10.10:80","mask":"255.255.255.255","source":"0.0.0.0/0","pool":"/Common/test-pool","enabled":true}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs/profiles", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[]}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs/policies", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[]}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name":        "/Common/test-vs",
		"destination": "10.10.10.10",
		"port":        80,
		"pool":        "/Common/test-pool",
	})
	d.SetId("/Common/test-vs")

	diags := resourceBigipLtmVirtualServerUpdate(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestResourceBigipLtmVirtualServerUpdate_Error(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"update failed"}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name":        "/Common/test-vs",
		"destination": "10.10.10.10",
		"port":        80,
	})
	d.SetId("/Common/test-vs")

	diags := resourceBigipLtmVirtualServerUpdate(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmVirtualServerDelete_Success(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name": "/Common/test-vs",
	})
	d.SetId("/Common/test-vs")

	diags := resourceBigipLtmVirtualServerDelete(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Id())
}

func TestResourceBigipLtmVirtualServerDelete_Error(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"code":500,"message":"delete failed"}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name": "/Common/test-vs",
	})
	d.SetId("/Common/test-vs")

	diags := resourceBigipLtmVirtualServerDelete(context.Background(), d, client)
	assert.True(t, diags.HasError())
}

func TestResourceBigipLtmVirtualServerRead_IPv6Destination(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs6", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
			"name":"test-vs6",
			"destination":"/Common/2001:db8::1.443",
			"mask":"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			"source":"::/0",
			"pool":"/Common/test-pool",
			"enabled":true
		}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs6/profiles", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[]}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs6/policies", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[]}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name": "/Common/test-vs6",
	})
	d.SetId("/Common/test-vs6")

	diags := resourceBigipLtmVirtualServerRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, 443, d.Get("port").(int))
}

func TestResourceBigipLtmVirtualServerRead_MaskAny(t *testing.T) {
	localMux := http.NewServeMux()
	localServer := httptest.NewServer(localMux)
	defer localServer.Close()

	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{
			"name":"test-vs",
			"destination":"/Common/10.10.10.10:80",
			"mask":"any",
			"source":"0.0.0.0/0",
			"pool":"/Common/test-pool",
			"enabled":false
		}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs/profiles", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[]}`)
	})
	localMux.HandleFunc("/mgmt/tm/ltm/virtual/~Common~test-vs/policies", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[]}`)
	})

	client := vsNewTestClient(localServer.URL)
	d := vsTestResourceData(t, map[string]interface{}{
		"name": "/Common/test-vs",
	})
	d.SetId("/Common/test-vs")

	diags := resourceBigipLtmVirtualServerRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "0.0.0.0", d.Get("mask").(string))
	assert.Equal(t, "disabled", d.Get("state").(string))
}
