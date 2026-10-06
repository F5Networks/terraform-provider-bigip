/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
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
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// clientSslResourceDataWithRawConfig builds a *schema.ResourceData for
// resourceBigipLtmProfileClientSsl whose GetRawConfig() actually returns
// rawConfig, unlike schema.TestResourceDataRaw (used by NewTestResourceData
// elsewhere in this file), which leaves GetRawConfig() null. This is needed
// for getClientSslConfig's tm_options handling specifically, which reads
// exclusively via d.GetRawConfig() (via the rawConfigAttrIsSet helper) to
// distinguish "the user set this in their .tf" from "Read backfilled this
// into state from an inherited/computed BIG-IP value" -- see the SFDC
// #01262589 fix. Thin wrapper around the shared
// NewTestResourceDataWithRawConfig helper (testing_helpers_test.go), which
// also backs policyResourceDataWithRawConfig in
// resource_bigip_ltm_policy_unit_test.go.
func clientSslResourceDataWithRawConfig(t *testing.T, raw map[string]interface{}, rawConfig cty.Value) *schema.ResourceData {
	t.Helper()
	return NewTestResourceDataWithRawConfig(t, resourceBigipLtmProfileClientSsl(), raw, rawConfig)
}

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipLtmProfileClientSslSchema(t *testing.T) {
	r := resourceBigipLtmProfileClientSsl()

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

	nameSchema, ok := r.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist in schema")
	}
	if !nameSchema.Required {
		t.Error("Expected field 'name' to be required")
	}
	if r.Schema["defaults_from"].Default != "/Common/clientssl" {
		t.Errorf("Expected 'defaults_from' to default to '/Common/clientssl', got %v", r.Schema["defaults_from"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitLtmProfileClientSslCreateReadUpdateDelete(t *testing.T) {
	name := "test-client-ssl"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","defaultsFrom":"/Common/clientssl","alertTimeout":"indefinite","allowNonSsl":"disabled","authenticate":"once","authenticateDepth":9,"cert":"/Common/default.crt","key":"/Common/default.key","chain":"none","ciphers":"DEFAULT","tmOptions":["Dont-Insert-Empty-Fragments"],"ocspStapling":"disabled","certKeyChain":[{"name":"default","cert":"/Common/default.crt","key":"/Common/default.key"}]}`, name)
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
	client.Teem = true // skip real telemetry network call
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                                name,
		"defaults_from":                       "/Common/clientssl",
		"ciphers":                             "DEFAULT",
		"alert_timeout":                       "indefinite",
		"allow_non_ssl":                       "disabled",
		"authenticate":                        "once",
		"authenticate_depth":                  9,
		"c3d_client_fallback_cert":            "none",
		"c3d_drop_unknown_ocsp_status":        "drop",
		"c3d_ocsp":                            "none",
		"ca_file":                             "none",
		"cache_size":                          262144,
		"cache_timeout":                       3600,
		"cert_extension_includes":             []interface{}{"basic-constraints"},
		"cert_life_span":                      30,
		"cert_lookup_by_ipaddr_port":          "disabled",
		"chain":                               "none",
		"key":                                 "/Common/default.key",
		"client_cert_ca":                      "none",
		"crl_file":                            "none",
		"allow_expired_crl":                   "disabled",
		"forward_proxy_bypass_default_action": "intercept",
		"generic_alert":                       "enabled",
		"handshake_timeout":                   "10",
		"inherit_cert_keychain":               "false",
		"mod_ssl_methods":                     "disabled",
		"mode":                                "enabled",
		"proxy_ca_cert":                       "none",
		"proxy_ca_key":                        "none",
		"peer_cert_mode":                      "ignore",
		"proxy_ca_passphrase":                 "secret",
		"proxy_ssl":                           "disabled",
		"proxy_ssl_passthrough":               "disabled",
		"renegotiate_period":                  "indefinite",
		"renegotiate_size":                    "indefinite",
		"renegotiation":                       "enabled",
		"retain_certificate":                  "true",
		"secure_renegotiation":                "require",
		"server_name":                         "example.com",
		"session_mirroring":                   "disabled",
		"session_ticket":                      "disabled",
		"sni_default":                         "false",
		"sni_require":                         "false",
		"ssl_c3d":                             "disabled",
		"ssl_forward_proxy":                   "disabled",
		"ssl_forward_proxy_bypass":            "disabled",
		"ssl_sign_hash":                       "any",
		"strict_resume":                       "enabled",
		"unclean_shutdown":                    "enabled",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileClientSSLCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileClientSSLRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "DEFAULT", d.Get("ciphers").(string))

	diags = resourceBigipLtmProfileClientSSLUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileClientSSLDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileClientSslRead_TmOptionsNone(t *testing.T) {
	name := "test-client-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","tmOptions":"none"}`, name)
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
	tmOptions := d.Get("tm_options").(*schema.Set)
	require.Equal(t, 0, tmOptions.Len())
}

func TestUnitLtmProfileClientSslRead_TmOptionsNil(t *testing.T) {
	// Exercises the nil branch directly: obj.TmOptions is untyped nil when
	// the tmOptions key is absent from BIG-IP's JSON response entirely
	// (e.g. never configured, no defaults_from to inherit from), since
	// go-bigip declares TmOptions as interface{} with an omitempty JSON tag
	// and encoding/json leaves an absent field at its Go zero value (nil)
	// rather than invoking UnmarshalJSON. Read must not call
	// reflect.TypeOf(obj.TmOptions).Kind() unguarded in this case, since
	// reflect.TypeOf(nil) returns a nil *rtype and .Kind() on it panics.
	name := "test-client-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	require.NotPanics(t, func() {
		diags := resourceBigipLtmProfileClientSSLRead(context.Background(), d, client)
		require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	})
	tmOptions := d.Get("tm_options").(*schema.Set)
	require.Equal(t, 0, tmOptions.Len())
}

func TestUnitLtmProfileClientSslRead_TmOptionsLegacyString(t *testing.T) {
	// Exercises the reflect.String branch: some BIG-IP versions return
	// tmOptions as a legacy brace-delimited string instead of a JSON array.
	name := "test-client-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","tmOptions":"{ Dont-Insert-Empty-Fragments }"}`, name)
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
}

func TestUnitLtmProfileClientSslCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-client-ssl",
	}, "")

	diags := resourceBigipLtmProfileClientSSLCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileClientSslRead_Error(t *testing.T) {
	name := "test-client-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileClientSSLRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileClientSslUpdate_Error(t *testing.T) {
	name := "test-client-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileClientSSLUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileClientSslDelete_Error(t *testing.T) {
	name := "test-client-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/client-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileClientSSLDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getClientSslConfig
// ---------------------------------------------------------------------

func TestGetClientSslConfig_NoCertKeyChain(t *testing.T) {
	raw := map[string]interface{}{
		"name":    "test-client-ssl",
		"cert":    "/Common/default.crt",
		"key":     "/Common/default.key",
		"chain":   "none",
		"ciphers": "DEFAULT",
	}
	// ciphers is read via rawConfigAttrIsSet (see getClientSslConfig's own
	// comment), so this must use the raw-config-aware helper -- plain
	// NewTestResourceData leaves GetRawConfig() null, which
	// rawConfigAttrIsSet treats as "not set in the user's .tf" regardless
	// of what's in raw/state.
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"ciphers":      cty.StringVal("DEFAULT"),
		"cipher_group": cty.NullVal(cty.String),
		"tm_options":   cty.NullVal(cty.Set(cty.String)),
	})
	d := clientSslResourceDataWithRawConfig(t, raw, rawConfig)

	cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
	require.Equal(t, "/Common/default.crt", cfg.Cert)
	require.Equal(t, "/Common/default.key", cfg.Key)
	require.Equal(t, "DEFAULT", cfg.Ciphers)
	require.Equal(t, "none", cfg.CipherGroup)
	require.Empty(t, cfg.CertKeyChain)
}

func TestGetClientSslConfig_WithCertKeyChain(t *testing.T) {
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-client-ssl",
		"cert_key_chain": []interface{}{
			map[string]interface{}{
				"name":       "default",
				"cert":       "/Common/default.crt",
				"key":        "/Common/default.key",
				"chain":      "none",
				"passphrase": "secret",
			},
		},
	}, "")

	cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
	require.Len(t, cfg.CertKeyChain, 1)
	require.Equal(t, "default", cfg.CertKeyChain[0].Name)
	require.Equal(t, "/Common/default.crt", cfg.CertKeyChain[0].Cert)
	// When cert_key_chain is set, the separate cert/key/chain/passphrase
	// fields are left unset on config.
	require.Empty(t, cfg.Cert)
}

func TestGetClientSslConfig_SslForwardProxyEnabled(t *testing.T) {
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                     "test-client-ssl",
		"ssl_forward_proxy":        "enabled",
		"ssl_forward_proxy_bypass": "",
	}, "")

	cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
	require.Equal(t, "/Common/default.crt", cfg.ProxyCaCert)
	require.Equal(t, "/Common/default.key", cfg.ProxyCaKey)
	require.Equal(t, "true", cfg.InheritCertkeychain)
	require.Equal(t, "disabled", cfg.SslForwardProxyBypass)
}

func TestGetClientSslConfig_CipherGroupSet(t *testing.T) {
	raw := map[string]interface{}{
		"name":         "test-client-ssl",
		"cipher_group": "/Common/my-cipher-group",
	}
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"cipher_group": cty.StringVal("/Common/my-cipher-group"),
		"ciphers":      cty.NullVal(cty.String),
		"tm_options":   cty.NullVal(cty.Set(cty.String)),
	})
	d := clientSslResourceDataWithRawConfig(t, raw, rawConfig)

	cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
	require.Equal(t, "/Common/my-cipher-group", cfg.CipherGroup)
	require.Equal(t, "none", cfg.Ciphers)
}

func TestGetClientSslConfig_ProxyCaNone(t *testing.T) {
	// proxy_ca_cert/proxy_ca_key == "none" should leave config.ProxyCaCert/Key
	// untouched (zero-value), per the "if proxyCaCert != none" guards.
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          "test-client-ssl",
		"proxy_ca_cert": "none",
		"proxy_ca_key":  "none",
	}, "")

	cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
	require.Empty(t, cfg.ProxyCaCert)
	require.Empty(t, cfg.ProxyCaKey)
}

func TestGetClientSslConfig_TmOptionsAndCertExtensionIncludes(t *testing.T) {
	raw := map[string]interface{}{
		"name":                    "test-client-ssl",
		"tm_options":              []interface{}{"Dont-Insert-Empty-Fragments"},
		"cert_extension_includes": []interface{}{"basic-constraints"},
	}
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"tm_options":   cty.SetVal([]cty.Value{cty.StringVal("Dont-Insert-Empty-Fragments")}),
		"ciphers":      cty.NullVal(cty.String),
		"cipher_group": cty.NullVal(cty.String),
	})
	d := clientSslResourceDataWithRawConfig(t, raw, rawConfig)

	cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
	require.Equal(t, []string{"Dont-Insert-Empty-Fragments"}, cfg.TmOptions)
	require.Equal(t, []string{"basic-constraints"}, cfg.CertExtensionIncludes)
}

// TestGetClientSslConfig_TmOptionsNotInRawConfigIsNotResent is the
// regression test for SFDC #01262589: a child profile whose own .tf never
// sets tm_options (raw config attribute is null) must NOT have tm_options
// included in the outgoing Create/Update payload, even if prior state
// already has a non-empty tm_options value backfilled by Read from an
// inherited defaults_from value (which is exactly what happens after any
// unrelated attribute update triggers another Read/Update cycle -- see
// getClientSslConfig's tm_options handling and its comment for the full
// root-cause explanation). Before the fix, this scenario incorrectly set
// cfg.TmOptions to the backfilled value, causing BIG-IP to receive an
// explicit "options" line and permanently breaking inheritance for that
// attribute on the child profile.
func TestGetClientSslConfig_TmOptionsNotInRawConfigIsNotResent(t *testing.T) {
	raw := map[string]interface{}{
		"name": "test-client-ssl",
	}
	// The user's current .tf does not mention tm_options at all.
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"tm_options":   cty.NullVal(cty.Set(cty.String)),
		"ciphers":      cty.NullVal(cty.String),
		"cipher_group": cty.NullVal(cty.String),
	})
	d := clientSslResourceDataWithRawConfig(t, raw, rawConfig)

	// Simulate a prior Read having backfilled state with the parent
	// profile's inherited options, as happens whenever defaults_from is
	// set and BIG-IP's GET response reflects the effective (inherited)
	// tmOptions rather than only locally-configured ones.
	require.NoError(t, d.Set("tm_options", []interface{}{"no-ssl", "no-tlsv1.1"}))
	_, getOkResult := d.GetOk("tm_options")
	require.True(t, getOkResult, "GetOk is expected to see the backfilled state value -- this is the bug trigger getClientSslConfig must not rely on")

	cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
	require.Empty(t, cfg.TmOptions, "tm_options must not be re-sent to BIG-IP when the user's .tf never set it, even if state has a backfilled inherited value")
}

// TestGetClientSslConfig_RawConfigNilIsNotResent covers the other half of
// the rawConfigAttrIsSet guard: d.GetRawConfig() itself being the SDK-null
// value (rawConfig.IsNull() true for the whole config, not just a null
// tm_options attribute within an otherwise-present object). This is exactly
// what schema.TestResourceDataRaw (and thus NewTestResourceData) produces,
// and is also what happens for real on terraform destroy, where there is no
// current .tf config to read at all. rawConfigAttrIsSet must short-circuit
// on this before ever calling cty.Value.GetAttr, which panics on a null
// object value rather than returning null.
func TestGetClientSslConfig_RawConfigNilIsNotResent(t *testing.T) {
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-client-ssl",
	}, "")
	require.True(t, d.GetRawConfig().IsNull(), "schema.TestResourceDataRaw is expected to leave GetRawConfig() null")

	require.False(t, rawConfigAttrIsSet(d, "tm_options"))

	// Simulate a prior Read having backfilled state with an inherited
	// value, as in TestGetClientSslConfig_TmOptionsNotInRawConfigIsNotResent,
	// to prove the nil-rawConfig case is guarded independently of state.
	require.NoError(t, d.Set("tm_options", []interface{}{"no-ssl", "no-tlsv1.1"}))
	_, getOkResult := d.GetOk("tm_options")
	require.True(t, getOkResult, "GetOk is expected to see the backfilled state value -- this is the bug trigger getClientSslConfig must not rely on")

	require.NotPanics(t, func() {
		cfg := getClientSslConfig(d, &bigip.ClientSSLProfile{})
		require.Empty(t, cfg.TmOptions, "tm_options must not be re-sent to BIG-IP when GetRawConfig() is null")
	})
}

// TestExplicitOptionsLineRegexp locks in explicitOptionsLineRegexp's
// tolerance for tmsh-list formatting variance (indentation depth,
// brace-spacing) used by testCheckClientSslNoExplicitOptions
// (resource_bigip_ltm_profile_ssl_client_test.go), and its rejection of
// unrelated lines that merely contain the substring "options" elsewhere in
// the object dump.
func TestExplicitOptionsLineRegexp(t *testing.T) {
	matchCases := map[string]string{
		"no indentation, space before brace": "options { dont-insert-empty-fragments }",
		"tmsh-typical 4-space indentation":   "    options {\n        dont-insert-empty-fragments\n    }",
		"deeper indentation":                 "        options {",
		"no space before brace":              "options{ dont-insert-empty-fragments }",
		"tabs instead of spaces":             "\toptions {",
		"mid-dump, other lines around it":    "ltm profile client-ssl /Common/child {\n    cert /Common/default.crt\n    options { dont-insert-empty-fragments }\n    key /Common/default.key\n}",
	}
	for desc, sample := range matchCases {
		t.Run("match/"+desc, func(t *testing.T) {
			require.True(t, explicitOptionsLineRegexp.MatchString(sample), "expected to match tmsh output:\n%s", sample)
		})
	}

	noMatchCases := map[string]string{
		"no options line at all":                     "ltm profile client-ssl /Common/child {\n    cert /Common/default.crt\n    key /Common/default.key\n}",
		"options as part of an unrelated identifier": "    proxy-options none",
		"options substring not at start of line":     "    # options { this is a comment, not tmsh syntax }",
		"options-like key without opening brace":     "    options none",
	}
	for desc, sample := range noMatchCases {
		t.Run("no-match/"+desc, func(t *testing.T) {
			require.False(t, explicitOptionsLineRegexp.MatchString(sample), "expected not to match tmsh output:\n%s", sample)
		})
	}
}
