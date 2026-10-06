/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// defaultTestConfigOptions returns ConfigOptions with sane values so that the
// go-bigip APICall retry loop actually executes (APICallRetries must be > 0)
// and requests don't hang forever.
func defaultTestConfigOptions() *bigip.ConfigOptions {
	return &bigip.ConfigOptions{
		APICallTimeout: 5 * time.Second,
		TokenTimeout:   20 * time.Minute,
		APICallRetries: 1,
	}
}

// newSelfIPServer spins up an httptest HTTP server that responds to the
// endpoint used by ValidateConnection() (SelfIPs -> /mgmt/tm/net/self).
func newSelfIPServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// newTLSSelfIPServer spins up an httptest TLS server that responds to the
// SelfIPs endpoint and also to the token authn endpoint.
func newTLSSelfIPServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	mux.HandleFunc("/mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// The returned timeout (1200s) must match defaultTestConfigOptions'
		// TokenTimeout (20m == 1200s). NewTokenSession only skips the follow-up
		// PATCH to /mgmt/shared/authz/tokens/... when these values are equal;
		// keeping them in sync avoids hitting that unregistered endpoint.
		_, _ = fmt.Fprint(w, `{"token":{"token":"unit-test-token"},"timeout":{"timeout":1200}}`)
	})
	srv := httptest.NewTLSServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// writeServerCertPEM writes the TLS server's certificate to a temp PEM file so
// that the CertVerifyDisable=false code path (os.ReadFile + AppendCertsFromPEM)
// succeeds and the server's cert is trusted.
func writeServerCertPEM(t *testing.T, srv *httptest.Server) string {
	t.Helper()
	require.NotNil(t, srv.TLS)
	require.NotEmpty(t, srv.TLS.Certificates)

	dir := t.TempDir()
	path := filepath.Join(dir, "trusted.pem")
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()

	for _, cert := range srv.TLS.Certificates {
		for _, der := range cert.Certificate {
			err = pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: der})
			require.NoError(t, err)
		}
	}
	return path
}

// TestClientNewSessionNoValidation exercises the plain NewSession branch where
// not all of Address/Username/Password are set, so ValidateConnection is
// skipped and the client is returned with a nil error.
func TestClientNewSessionNoValidation(t *testing.T) {
	config := &bigip.Config{
		Address:  "http://192.0.2.10",
		Username: "admin",
		// Password intentionally empty -> validation branch skipped.
		ConfigOptions: defaultTestConfigOptions(),
	}

	client, err := Client(config)

	assert.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "admin", client.User)
	assert.Empty(t, client.Token)
}

// TestClientNewSessionWithToken verifies that when a Token is provided (and no
// LoginReference), NewSession is used and the token is copied onto the client.
func TestClientNewSessionWithToken(t *testing.T) {
	config := &bigip.Config{
		Address:       "http://192.0.2.11",
		Token:         "preauth-token",
		ConfigOptions: defaultTestConfigOptions(),
	}

	client, err := Client(config)

	assert.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "preauth-token", client.Token)
}

// TestClientValidateConnectionSuccessInsecure covers the full happy path with
// CertVerifyDisable=true: NewSession + InsecureSkipVerify + successful
// ValidateConnection against the mock server.
func TestClientValidateConnectionSuccessInsecure(t *testing.T) {
	srv := newSelfIPServer(t)

	config := &bigip.Config{
		Address:           srv.URL,
		Username:          "admin",
		Password:          "secret",
		CertVerifyDisable: true,
		ConfigOptions:     defaultTestConfigOptions(),
	}

	client, err := Client(config)

	assert.NoError(t, err)
	require.NotNil(t, client)
	assert.True(t, client.Transport.TLSClientConfig.InsecureSkipVerify)
}

// TestClientValidateConnectionError covers the branch where all connection
// params are set but ValidateConnection fails (server returns an error status).
func TestClientValidateConnectionError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	config := &bigip.Config{
		Address:           srv.URL,
		Username:          "admin",
		Password:          "secret",
		CertVerifyDisable: true,
		ConfigOptions:     defaultTestConfigOptions(),
	}

	client, err := Client(config)

	assert.Error(t, err)
	// Even on validation failure the (non-nil) client is returned.
	assert.NotNil(t, client)
}

// TestClientTrustedCertificateSuccess covers CertVerifyDisable=false with a
// valid trusted certificate path: os.ReadFile succeeds, AppendCertsFromPEM
// succeeds, RootCAs is set, and ValidateConnection succeeds over TLS.
func TestClientTrustedCertificateSuccess(t *testing.T) {
	srv := newTLSSelfIPServer(t)
	certPath := writeServerCertPEM(t, srv)

	config := &bigip.Config{
		Address:            srv.URL,
		Username:           "admin",
		Password:           "secret",
		CertVerifyDisable:  false,
		TrustedCertificate: certPath,
		ConfigOptions:      defaultTestConfigOptions(),
	}

	client, err := Client(config)

	assert.NoError(t, err)
	require.NotNil(t, client)
	assert.False(t, client.Transport.TLSClientConfig.InsecureSkipVerify)
	require.NotNil(t, client.Transport.TLSClientConfig.RootCAs)
}

// TestClientTrustedCertificateInvalidPath covers CertVerifyDisable=false where
// os.ReadFile fails because the trusted certificate path does not exist.
func TestClientTrustedCertificateInvalidPath(t *testing.T) {
	config := &bigip.Config{
		Address:            "https://192.0.2.20",
		Username:           "admin",
		Password:           "secret",
		CertVerifyDisable:  false,
		TrustedCertificate: filepath.Join(t.TempDir(), "does-not-exist.pem"),
		ConfigOptions:      defaultTestConfigOptions(),
	}

	client, err := Client(config)

	require.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "provide Valid Trusted certificate path")
}

// TestClientTrustedCertificateNoCertsAppended covers the branch where
// os.ReadFile succeeds but AppendCertsFromPEM returns false (the file exists
// but contains no valid PEM certificates). Execution continues and validation
// still runs (and fails to connect, which is fine for this branch).
func TestClientTrustedCertificateNoCertsAppended(t *testing.T) {
	dir := t.TempDir()
	badCert := filepath.Join(dir, "not-a-cert.pem")
	require.NoError(t, os.WriteFile(badCert, []byte("not a valid pem certificate"), 0600))

	// Point at a closed server address so ValidateConnection fails fast.
	srv := httptest.NewServer(http.NewServeMux())
	addr := srv.URL
	srv.Close()

	config := &bigip.Config{
		Address:            addr,
		Username:           "admin",
		Password:           "secret",
		CertVerifyDisable:  false,
		TrustedCertificate: badCert,
		ConfigOptions:      defaultTestConfigOptions(),
	}

	client, err := Client(config)

	// AppendCertsFromPEM failed (logged), then ValidateConnection errored.
	assert.Error(t, err)
	require.NotNil(t, client)
	// AppendCertsFromPEM returned false for the malformed PEM (logged as "no
	// certs appended"), but RootCAs is set unconditionally, so it is still
	// non-nil (falling back to whatever system roots were loaded).
	require.NotNil(t, client.Transport.TLSClientConfig.RootCAs)
}

// TestClientTokenSessionSuccess covers the NewTokenSession branch: LoginReference
// set, Token empty, Address set. The mock TLS server returns an auth token and
// SelfIPs, so the full path (token session + validation) succeeds.
func TestClientTokenSessionSuccess(t *testing.T) {
	srv := newTLSSelfIPServer(t)

	config := &bigip.Config{
		Address:           srv.URL,
		Username:          "admin",
		Password:          "secret",
		LoginReference:    "tmos",
		CertVerifyDisable: true,
		ConfigOptions:     defaultTestConfigOptions(),
	}

	client, err := Client(config)

	assert.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "unit-test-token", client.Token)
}

// TestClientTokenSessionError covers the NewTokenSession error branch where the
// token session cannot be established (server unreachable), returning nil,err.
func TestClientTokenSessionError(t *testing.T) {
	// Create then immediately close a server to get an unreachable address.
	srv := httptest.NewServer(http.NewServeMux())
	addr := srv.URL
	srv.Close()

	config := &bigip.Config{
		Address:           addr,
		Username:          "admin",
		Password:          "secret",
		LoginReference:    "tmos",
		CertVerifyDisable: true,
		ConfigOptions:     defaultTestConfigOptions(),
	}

	client, err := Client(config)

	require.Error(t, err)
	assert.Nil(t, client)
}

// TestConfigStructInitialization validates that the bigip.Config struct fields
// are correctly consumed by Client() and mapped onto the resulting BigIP
// session (connection parameter validation).
func TestConfigStructInitialization(t *testing.T) {
	config := &bigip.Config{
		Address:  "http://bigip.example.com",
		Port:     "8443",
		Username: "myuser",
		// Password intentionally empty so the ValidateConnection branch (which
		// requires Address+Username+Password) is skipped and no network call is
		// attempted; we only verify config->session field mapping here.
		Token:             "mytoken",
		CertVerifyDisable: true,
		LoginReference:    "",
		ConfigOptions:     defaultTestConfigOptions(),
	}

	client, err := Client(config)

	assert.NoError(t, err)
	require.NotNil(t, client)

	// Host is derived from Address (+ Port) by NewSession.
	assert.True(t, strings.HasPrefix(client.Host, "http://bigip.example.com"))
	assert.Contains(t, client.Host, "8443")
	assert.Equal(t, "myuser", client.User)
	// Token was provided and copied onto the client.
	assert.Equal(t, "mytoken", client.Token)
}

// TestClientTLSConfigDefaults is a small sanity check ensuring a fresh session
// exposes a usable *tls.Config on its transport (used elsewhere in Client()).
func TestClientTLSConfigDefaults(t *testing.T) {
	config := &bigip.Config{
		Address:       "http://192.0.2.30",
		ConfigOptions: defaultTestConfigOptions(),
	}

	client, err := Client(config)

	assert.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.Transport)
	var _ *tls.Config = client.Transport.TLSClientConfig
}
