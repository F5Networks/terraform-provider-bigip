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
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipNetIkePeerSchema(t *testing.T) {
	r := resourceBigipNetIkePeer()

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

	for _, field := range []string{"name", "remote_address"} {
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
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitNetIkePeerCreateReadUpdateDelete(t *testing.T) {
	name := "/Common/test-ike-peer"
	mangled := "/mgmt/tm/net/ipsec/ike-peer/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/ipsec/ike-peer", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","remoteAddress":"10.0.0.5","presharedKey":"","presharedKeyEncrypted":"$M$enc","version":["v2"],"trafficSelector":["/Common/ts1"],"dpdDelay":10,"lifetime":1440,"replayWindowSize":64}`, name)
		case http.MethodPatch:
			_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetIkePeer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":             name,
		"remote_address":   "10.0.0.5",
		"version":          []interface{}{"v2"},
		"traffic_selector": []interface{}{"/Common/ts1"},
	}, "")

	ctx := context.Background()

	diags := resourceBigipNetIkePeerCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipNetIkePeerRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "10.0.0.5", d.Get("remote_address").(string))

	diags = resourceBigipNetIkePeerUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipNetIkePeerDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitNetIkePeerRead_PresharedKeySet(t *testing.T) {
	// Exercises the "ikepeer.PresharedKey != \"\" && d.Get(\"preshared_key\")
	// != \"\"" branch: the mock response includes a non-empty presharedKey,
	// and prior state also has a non-empty preshared_key.
	name := "/Common/test-ike-peer"
	mangled := "/mgmt/tm/net/ipsec/ike-peer/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","remoteAddress":"10.0.0.5","presharedKey":"newsecret"}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetIkePeer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":           name,
		"remote_address": "10.0.0.5",
		"preshared_key":  "oldsecret",
	}, name)

	diags := resourceBigipNetIkePeerRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "newsecret", d.Get("preshared_key").(string))
}

func TestUnitNetIkePeerCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/ipsec/ike-peer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetIkePeer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":           "/Common/test-ike-peer",
		"remote_address": "10.0.0.5",
	}, "")

	diags := resourceBigipNetIkePeerCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetIkePeerRead_Error(t *testing.T) {
	name := "/Common/test-ike-peer"
	mangled := "/mgmt/tm/net/ipsec/ike-peer/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetIkePeer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":           name,
		"remote_address": "10.0.0.5",
	}, name)

	diags := resourceBigipNetIkePeerRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetIkePeerUpdate_Error(t *testing.T) {
	name := "/Common/test-ike-peer"
	mangled := "/mgmt/tm/net/ipsec/ike-peer/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetIkePeer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":           name,
		"remote_address": "10.0.0.5",
	}, name)

	diags := resourceBigipNetIkePeerUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitNetIkePeerDelete_Error(t *testing.T) {
	name := "/Common/test-ike-peer"
	mangled := "/mgmt/tm/net/ipsec/ike-peer/" + MangleFullPath(name)

	mux := http.NewServeMux()
	mux.HandleFunc(mangled, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipNetIkePeer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":           name,
		"remote_address": "10.0.0.5",
	}, name)

	diags := resourceBigipNetIkePeerDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getIkeConfig
// ---------------------------------------------------------------------

func TestGetIkeConfig(t *testing.T) {
	r := resourceBigipNetIkePeer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                           "/Common/test-ike-peer",
		"app_service":                    "my-app",
		"ca_cert_file":                   "/Common/ca.crt",
		"crl_file":                       "/Common/crl.pem",
		"description":                    "my ike peer",
		"dpd_delay":                      10,
		"lifetime":                       1440,
		"generate_policy":                "off",
		"mode":                           "main",
		"my_cert_file":                   "/Common/my.crt",
		"my_cert_key_file":               "/Common/my.key",
		"my_cert_key_passphrase":         "secret",
		"my_id_type":                     "address",
		"my_id_value":                    "10.0.0.1",
		"nat_traversal":                  "off",
		"passive":                        "false",
		"peers_cert_file":                "/Common/peer.crt",
		"peers_cert_type":                "certfile",
		"peers_id_type":                  "address",
		"peers_id_value":                 "10.0.0.5",
		"phase1_auth_method":             "rsa-signature",
		"phase1_encrypt_algorithm":       "aes256",
		"phase1_hash_algorithm":          "sha256",
		"phase1_perfect_forward_secrecy": "modp2048",
		"preshared_key":                  "secretkey",
		"prf":                            "hmac-sha256",
		"proxy_support":                  "enabled",
		"remote_address":                 "10.0.0.5",
		"replay_window_size":             64,
		"state":                          "enabled",
		"traffic_selector":               []interface{}{"/Common/ts1", "/Common/ts2"},
		"verify_cert":                    "true",
		"version":                        []interface{}{"v1", "v2"},
	}, "")

	cfg := getIkeConfig(d, &bigip.IkePeer{})
	require.Equal(t, "my-app", cfg.AppService)
	require.Equal(t, "/Common/ca.crt", cfg.CaCertFile)
	require.Equal(t, "/Common/crl.pem", cfg.CrlFile)
	require.Equal(t, 10, cfg.DpdDelay)
	require.Equal(t, 1440, cfg.Lifetime)
	require.Equal(t, "my ike peer", cfg.Description)
	require.Equal(t, "off", cfg.GeneratePolicy)
	require.Equal(t, "main", cfg.Mode)
	require.Equal(t, "/Common/my.crt", cfg.MyCertFile)
	require.Equal(t, "/Common/my.key", cfg.MyCertKeyFile)
	require.Equal(t, "secret", cfg.MyCertKeyPassphrase)
	require.Equal(t, "address", cfg.MyIdType)
	require.Equal(t, "10.0.0.1", cfg.MyIdValue)
	require.Equal(t, "off", cfg.NatTraversal)
	require.Equal(t, "false", cfg.Passive)
	require.Equal(t, "/Common/peer.crt", cfg.PeersCertFile)
	require.Equal(t, "certfile", cfg.PeersCertType)
	require.Equal(t, "address", cfg.PeersIdType)
	require.Equal(t, "10.0.0.5", cfg.PeersIdValue)
	require.Equal(t, "rsa-signature", cfg.Phase1AuthMethod)
	require.Equal(t, "aes256", cfg.Phase1EncryptAlgorithm)
	require.Equal(t, "sha256", cfg.Phase1HashAlgorithm)
	require.Equal(t, "modp2048", cfg.Phase1PerfectForwardSecrecy)
	require.Equal(t, "secretkey", cfg.PresharedKey)
	require.Equal(t, "hmac-sha256", cfg.Prf)
	require.Equal(t, "enabled", cfg.ProxySupport)
	require.Equal(t, "10.0.0.5", cfg.RemoteAddress)
	require.Equal(t, 64, cfg.ReplayWindowSize)
	require.Equal(t, "enabled", cfg.State)
	require.Equal(t, []string{"/Common/ts1", "/Common/ts2"}, cfg.TrafficSelector)
	require.Equal(t, "true", cfg.VerifyCert)
	require.Equal(t, []string{"v1", "v2"}, cfg.Version)
}

func TestGetIkeConfig_NoVersionOrTrafficSelector(t *testing.T) {
	// Exercises the "!ok" branches of GetOk("version")/GetOk("traffic_selector")
	// when neither is set in config.
	r := resourceBigipNetIkePeer()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":           "/Common/test-ike-peer",
		"remote_address": "10.0.0.5",
	}, "")

	cfg := getIkeConfig(d, &bigip.IkePeer{})
	require.Empty(t, cfg.Version)
	require.Empty(t, cfg.TrafficSelector)
}
