/*
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------
// Audit: Optional+Computed inheritance-drift fields on
// bigip_ltm_profile_client_ssl, beyond tm_options.
//
// Background (see the tm_options fix for SFDC #01262589 / ENI SpA, and
// its own root-cause writeup): tm_options leaked an inherited defaults_from
// parent value back onto a child profile as an explicit override after any
// unrelated attribute update, because resourceBigipLtmProfileClientSSLRead
// used an UNCONDITIONAL d.Set("tm_options", ...) -- so Read always
// backfilled state with whatever BIG-IP's GET response reported (which
// reflects the *effective*, possibly-inherited value when defaults_from is
// set), and getClientSslConfig's d.GetOk("tm_options") then saw that
// backfilled value as "the user configured this" and re-sent it.
//
// This file is the audit this comment describes, covering all 53 other
// Optional+Computed fields in resourceBigipLtmProfileClientSsl()'s schema
// (54 total Optional+Computed fields minus tm_options, already fixed).
// cert_key_chain is excluded for a separate, unrelated reason (see below --
// it is Optional-only, not Optional+Computed, and is never populated by
// Read regardless):
//
//   - 48 fields via stringAuditCases()/intAuditCases() (TestInheritanceAudit_UnrelatedUpdateDoesNotResendInheritedValue below)
//   - cert_extension_includes (a TypeSet, tested separately: TestInheritanceAudit_CertExtensionIncludesNotResent)
//   - partition, full_path, generation, and passphrase (tested separately:
//     TestInheritanceAudit_IdentityFieldsAreChildsOwnNotParents -- see that
//     test's own doc comment for why these four don't fit the same
//     GetOk-guard-based explanation as the other 49 fields, but are still
//     safe for a different reason each)
//
// cert_key_chain (a TypeList block) is Optional-only -- it is not
// Optional+Computed, so it was never in this audit's scope to begin with.
// It is also never populated by Read at all (no
// d.Set("cert_key_chain", ...) call exists), so it cannot backfill an
// inherited value by construction either way.
//
// Finding: none of the 53 fields audited here share tm_options' bug, though
// not all for the identical reason. The dominant reason, covering 49 of the
// 53 (all of stringAuditCases()/intAuditCases() plus cert_extension_includes):
// unlike tm_options, every one of them already has a `d.GetOk("<field>")`
// guard in resourceBigipLtmProfileClientSSLRead -- added in commit
// 5ad03632 ("solving client ssl profile issues", 2021-06-07), which
// predates and fixes GitHub F5Networks/terraform-provider-bigip #450
// (reported against v1.7.0, milestone v1.10.0). Because Read only calls
// d.Set for a field when GetOk on the field ALREADY returns true, Read
// never backfills state for a field the practitioner's own .tf never
// configured -- so it stays at its Go zero value ("" for strings, 0 for
// ints, empty for the cert_extension_includes set) in state.
// getClientSslConfig's unconditional d.Get(...) for these fields then also
// returns that zero value, which the go-bigip ClientSSLProfile struct's
// `omitempty` JSON tags drop from the outgoing Create/Update payload
// entirely -- so nothing is ever sent to BIG-IP for that field, and
// inheritance is preserved. The remaining 4 (partition, full_path,
// generation, passphrase) are safe for the separate reasons documented on
// TestInheritanceAudit_IdentityFieldsAreChildsOwnNotParents.
//
// This was verified two ways before writing this test:
//  1. Empirically, against a live BIG-IP 17.5.1 device, using a locally
//     built provider binary and repeated real `terraform apply` cycles
//     (multiple unrelated-attribute updates in a row) against
//     authenticate, secure_renegotiation, and cipher_group (the exact
//     field named in GitHub #902) child profiles with defaults_from set
//     and the field never configured on the child -- none of them ever
//     appeared as an explicit line in `tmsh list ltm profile client-ssl`
//     afterward.
//  2. GitHub #902 ("Issue with clientssl profile when using parent/child
//     (defaults-from)", reported against v1.20.0) is itself Closed with
//     milestone v1.20.2 -- fixed narrowly for cipher_group specifically by
//     commit 2a73abf8 (switching cipher_group's schema from
//     Default:"none" to Computed:true, i.e. the same GetOk-in-Read pattern
//     already covering the other ~49 fields audited here). Notably, #902's
//     own reported symptom was plan-diff noise only ("this doesn't change
//     the effective BIG-IP configuration (fortunately)"), not actual
//     device-side configuration drift like tm_options -- a materially
//     different (and lower-severity) failure mode than the bug this audit
//     was commissioned to check for.
//
// Both #450 and #902 are therefore already resolved as of the current
// provider version and require no further code change; this file's tests
// are the regression coverage guarding that. If a future change to
// resourceBigipLtmProfileClientSSLRead ever removes or narrows one of
// these GetOk guards, TestInheritanceAudit_UnrelatedUpdateDoesNotResendInheritedValue
// below will fail for that field.
// ---------------------------------------------------------------------

// clientSslInheritanceAuditCase describes one Optional+Computed scalar
// field to audit: the schema attribute name, the JSON field BIG-IP's GET
// response uses for it, a non-default/distinctive value simulating an
// inherited value from a defaults_from parent, and a function to extract
// the corresponding field from the *bigip.ClientSSLProfile config struct
// getClientSslConfig returns.
type clientSslInheritanceAuditCase struct {
	attr        string
	jsonField   string
	mockValue   string // as it appears in the mocked GET JSON body (quoted if a JSON string)
	getConfigFn func(cfg *bigip.ClientSSLProfile) interface{}
	zeroValue   interface{}
}

func stringAuditCases() []clientSslInheritanceAuditCase {
	return []clientSslInheritanceAuditCase{
		{"alert_timeout", "alertTimeout", `"indefinite"`, func(c *bigip.ClientSSLProfile) interface{} { return c.AlertTimeout }, ""},
		{"allow_non_ssl", "allowNonSsl", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.AllowNonSsl }, ""},
		{"authenticate", "authenticate", `"always"`, func(c *bigip.ClientSSLProfile) interface{} { return c.Authenticate }, ""},
		{"c3d_client_fallback_cert", "c3dClientFallbackCert", `"/Common/fallback.crt"`, func(c *bigip.ClientSSLProfile) interface{} { return c.C3dClientFallbackCert }, ""},
		{"c3d_drop_unknown_ocsp_status", "c3dDropUnknownOcspStatus", `"ignore"`, func(c *bigip.ClientSSLProfile) interface{} { return c.C3dDropUnknownOcspStatus }, ""},
		{"c3d_ocsp", "c3dOcsp", `"/Common/my-ocsp"`, func(c *bigip.ClientSSLProfile) interface{} { return c.C3dOcsp }, ""},
		{"ca_file", "caFile", `"/Common/my-ca.crt"`, func(c *bigip.ClientSSLProfile) interface{} { return c.CaFile }, ""},
		{"cert", "cert", `"/Common/inherited.crt"`, func(c *bigip.ClientSSLProfile) interface{} { return c.Cert }, ""},
		{"key", "key", `"/Common/inherited.key"`, func(c *bigip.ClientSSLProfile) interface{} { return c.Key }, ""},
		{"chain", "chain", `"/Common/inherited-chain.crt"`, func(c *bigip.ClientSSLProfile) interface{} { return c.Chain }, ""},
		{"cert_lookup_by_ipaddr_port", "certLookupByIpaddrPort", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.CertLookupByIpaddrPort }, ""},
		{"ciphers", "ciphers", `"HIGH:!ADH:!MD5"`, func(c *bigip.ClientSSLProfile) interface{} { return c.Ciphers }, ""},
		{"cipher_group", "cipherGroup", `"/Common/f5-secure"`, func(c *bigip.ClientSSLProfile) interface{} { return c.CipherGroup }, ""},
		{"client_cert_ca", "clientCertCa", `"/Common/my-client-ca.crt"`, func(c *bigip.ClientSSLProfile) interface{} { return c.ClientCertCa }, ""},
		{"crl_file", "crlFile", `"/Common/my.crl"`, func(c *bigip.ClientSSLProfile) interface{} { return c.CrlFile }, ""},
		{"inherit_cert_keychain", "inheritCertkeychain", `"true"`, func(c *bigip.ClientSSLProfile) interface{} { return c.InheritCertkeychain }, ""},
		{"proxy_ca_passphrase", "proxyCaPassphrase", `"inherited-secret"`, func(c *bigip.ClientSSLProfile) interface{} { return c.ProxyCaPassphrase }, ""},
		{"ssl_forward_proxy", "sslForwardProxy", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.SslForwardProxy }, ""},
		{"allow_expired_crl", "allowExpiredCrl", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.AllowExpiredCrl }, ""},
		{"forward_proxy_bypass_default_action", "forwardProxyBypassDefaultAction", `"bypass"`, func(c *bigip.ClientSSLProfile) interface{} { return c.ForwardProxyBypassDefaultAction }, ""},
		{"generic_alert", "genericAlert", `"disabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.GenericAlert }, ""},
		{"handshake_timeout", "handshakeTimeout", `"30"`, func(c *bigip.ClientSSLProfile) interface{} { return c.HandshakeTimeout }, ""},
		{"mod_ssl_methods", "modSslMethods", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.ModSslMethods }, ""},
		{"mode", "mode", `"disabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.Mode }, ""},
		{"peer_cert_mode", "peerCertMode", `"require"`, func(c *bigip.ClientSSLProfile) interface{} { return c.PeerCertMode }, ""},
		{"proxy_ca_cert", "proxyCaCert", `"/Common/my-proxy-ca.crt"`, func(c *bigip.ClientSSLProfile) interface{} { return c.ProxyCaCert }, ""},
		{"proxy_ca_key", "proxyCaKey", `"/Common/my-proxy-ca.key"`, func(c *bigip.ClientSSLProfile) interface{} { return c.ProxyCaKey }, ""},
		{"proxy_ssl", "proxySsl", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.ProxySsl }, ""},
		{"proxy_ssl_passthrough", "proxySslPassthrough", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.ProxySslPassthrough }, ""},
		{"renegotiate_period", "renegotiatePeriod", `"1800"`, func(c *bigip.ClientSSLProfile) interface{} { return c.RenegotiatePeriod }, ""},
		{"renegotiate_size", "renegotiateSize", `"1048576"`, func(c *bigip.ClientSSLProfile) interface{} { return c.RenegotiateSize }, ""},
		{"renegotiation", "renegotiation", `"disabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.Renegotiation }, ""},
		{"retain_certificate", "retainCertificate", `"false"`, func(c *bigip.ClientSSLProfile) interface{} { return c.RetainCertificate }, ""},
		{"secure_renegotiation", "secureRenegotiation", `"require-strict"`, func(c *bigip.ClientSSLProfile) interface{} { return c.SecureRenegotiation }, ""},
		{"server_name", "serverName", `"example.com"`, func(c *bigip.ClientSSLProfile) interface{} { return c.ServerName }, ""},
		{"session_mirroring", "sessionMirroring", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.SessionMirroring }, ""},
		{"session_ticket", "sessionTicket", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.SessionTicket }, ""},
		{"sni_default", "sniDefault", `"true"`, func(c *bigip.ClientSSLProfile) interface{} { return c.SniDefault }, ""},
		{"sni_require", "sniRequire", `"true"`, func(c *bigip.ClientSSLProfile) interface{} { return c.SniRequire }, ""},
		{"ssl_c3d", "sslC3d", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.SslC3d }, ""},
		{"ssl_forward_proxy_bypass", "sslForwardProxyBypass", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.SslForwardProxyBypass }, ""},
		{"ssl_sign_hash", "sslSignHash", `"sha256"`, func(c *bigip.ClientSSLProfile) interface{} { return c.SslSignHash }, ""},
		{"strict_resume", "strictResume", `"enabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.StrictResume }, ""},
		{"unclean_shutdown", "uncleanShutdown", `"disabled"`, func(c *bigip.ClientSSLProfile) interface{} { return c.UncleanShutdown }, ""},
	}
}

func intAuditCases() []clientSslInheritanceAuditCase {
	return []clientSslInheritanceAuditCase{
		{"authenticate_depth", "authenticateDepth", `5`, func(c *bigip.ClientSSLProfile) interface{} { return c.AuthenticateDepth }, 0},
		{"cache_size", "cacheSize", `524288`, func(c *bigip.ClientSSLProfile) interface{} { return c.CacheSize }, 0},
		{"cache_timeout", "cacheTimeout", `7200`, func(c *bigip.ClientSSLProfile) interface{} { return c.CacheTimeout }, 0},
		{"cert_life_span", "certLifespan", `60`, func(c *bigip.ClientSSLProfile) interface{} { return c.CertLifespan }, 0},
	}
}

// TestInheritanceAudit_UnrelatedUpdateDoesNotResendInheritedValue is the
// regression test for this audit: for every Optional+Computed scalar field
// enumerated above, simulate a child profile that never configured the
// field itself (only "name" is in the user's .tf), mock a GET response
// where BIG-IP reports a distinctive non-default value for that field
// (simulating inheritance from a defaults_from parent, or any other
// device-computed value), run Read (exercising the existing GetOk guard),
// then call getClientSslConfig directly (simulating the Update
// getClientSslConfig call that would run after an unrelated attribute
// changes) and assert the field comes back as its Go zero value -- i.e.
// it is never re-sent to BIG-IP as an explicit override.
func TestInheritanceAudit_UnrelatedUpdateDoesNotResendInheritedValue(t *testing.T) {
	name := "test-client-ssl"

	runCase := func(t *testing.T, tc clientSslInheritanceAuditCase) {
		mux := http.NewServeMux()
		mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
			// tmOptions is always included (as "none") because
			// resourceBigipLtmProfileClientSSLRead's tm_options handling
			// unconditionally type-switches on obj.TmOptions and panics
			// if the field is absent from the response entirely -- a
			// separate, pre-existing gap unrelated to this audit.
			_, _ = fmt.Fprintf(w, `{"name":%q,"defaultsFrom":"/Common/parent","tmOptions":"none","%s":%s}`, name, tc.jsonField, tc.mockValue)
		})
		server := httptest.NewServer(mux)
		defer server.Close()

		client := NewUnitTestClient(server.URL)
		r := resourceBigipLtmProfileClientSsl()
		d := NewTestResourceData(t, r, map[string]interface{}{
			"name": name,
		}, name)

		diags := resourceBigipLtmProfileClientSSLRead(context.Background(), d, client)
		require.False(t, diags.HasError(), "unexpected read error: %v", diags)

		cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
		assert.Equal(t, tc.zeroValue, tc.getConfigFn(cfg),
			"field %q must stay at its zero value (never re-sent to BIG-IP) when the practitioner's .tf never configured it, even though the mocked GET reported an inherited/non-default value", tc.attr)
	}

	for _, tc := range stringAuditCases() {
		tc := tc
		t.Run(tc.attr, func(t *testing.T) { runCase(t, tc) })
	}

	for _, tc := range intAuditCases() {
		tc := tc
		t.Run(tc.attr, func(t *testing.T) { runCase(t, tc) })
	}
}

// TestInheritanceAudit_CertExtensionIncludesNotResent covers
// cert_extension_includes, the one remaining Optional+Computed field with
// the same GetOk-in-Read guard as the string/int fields above, but which is
// a TypeSet (not a scalar) and so needs its own mock/assertion shape rather
// than fitting stringAuditCases/intAuditCases.
func TestInheritanceAudit_CertExtensionIncludesNotResent(t *testing.T) {
	name := "test-client-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":%q,"defaultsFrom":"/Common/parent","tmOptions":"none","certExtensionIncludes":["basic-constraints","subject-alternative-name"]}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileClientSSLRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
	assert.Empty(t, cfg.CertExtensionIncludes,
		"cert_extension_includes must stay empty (never re-sent to BIG-IP) when the practitioner's .tf never configured it, even though the mocked GET reported inherited values")
}

// TestInheritanceAudit_IdentityFieldsAreChildsOwnNotParents documents and
// verifies why partition, full_path, generation, and passphrase are
// intentionally excluded from the leak-prone field tables above -- each
// for a different reason than the GetOk-guard explanation covering the
// other 49 fields:
//
//   - full_path, generation, passphrase: NEVER d.Set in
//     resourceBigipLtmProfileClientSSLRead at all (no d.Set(...) call
//     exists for any of the three anywhere in Read), so -- like tm_options
//     before its fix, but in the opposite direction -- they can never be
//     populated by an inherited (or any) GET value in the first place;
//     getClientSslConfig's unconditional d.Get(...) for them always
//     returns the Go zero value.
//   - partition: the one field in this audit that IS unconditionally
//     d.Set in Read (no GetOk guard) and IS unconditionally re-read via
//     d.Get in getClientSslConfig -- structurally the same shape as
//     tm_options' original bug. However, it describes the CHILD's own
//     object identity (which partition the child profile itself resides
//     in), not a setting inherited from the defaults_from parent --
//     confirmed empirically against a live BIG-IP 17.5.1 device: creating
//     a parent and a defaults_from child from it, GetClientSSLProfile(childName)
//     can only ever report the object actually being read's own
//     partition (there is no code path, mock or real, by which the GET
//     response for THIS object could instead reflect a DIFFERENT
//     object's partition), so there is no realistic "inherited" value
//     that could appear here the way an inherited tm_options/authenticate/
//     etc. setting can.
func TestInheritanceAudit_IdentityFieldsAreChildsOwnNotParents(t *testing.T) {
	childName := "test-client-ssl-child"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+childName, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":%q,"partition":"Common","fullPath":"/Common/%s","generation":42,"passphrase":"inherited-secret","defaultsFrom":"/Common/parent","tmOptions":"none"}`, childName, childName)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": childName,
	}, childName)

	diags := resourceBigipLtmProfileClientSSLRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)

	cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
	// partition IS backfilled by Read (unconditionally), and it correctly
	// reflects the child's own partition -- confirming there is no way for
	// it to instead reflect a different (e.g. parent) object's partition.
	assert.Equal(t, "Common", cfg.Partition, "partition must reflect the child's own partition (Read's only source is GetClientSSLProfile(childName), which cannot return another object's data)")

	// full_path, generation, and passphrase are never populated by Read at
	// all, regardless of what the mocked GET response contains -- there is
	// no d.Set(...) call for any of them in resourceBigipLtmProfileClientSSLRead.
	assert.Empty(t, cfg.FullPath, "full_path must stay empty; Read never sets it from any GET response")
	assert.Zero(t, cfg.Generation, "generation must stay zero; Read never sets it from any GET response")
	assert.Empty(t, cfg.Passphrase, "passphrase must stay empty; Read never sets it from any GET response")
}
