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

// serverSslResourceDataWithRawConfig builds a *schema.ResourceData for
// resourceBigipLtmProfileServerSsl whose GetRawConfig() actually returns
// rawConfig, unlike schema.TestResourceDataRaw (used by NewTestResourceData
// elsewhere in this file), which leaves GetRawConfig() null. This is needed
// for getServerSslConfig's tm_options handling specifically, which reads
// exclusively via d.GetRawConfig() (via the rawConfigAttrIsSet helper) to
// distinguish "the user set this in their .tf" from "Read backfilled this
// into state from an inherited/computed BIG-IP value" -- see the SFDC
// #01262589 fix (originally applied to bigip_ltm_profile_client_ssl; this
// mirrors the same fix for bigip_ltm_profile_server_ssl). Thin wrapper
// around the shared NewTestResourceDataWithRawConfig helper
// (testing_helpers_test.go), which also backs
// clientSslResourceDataWithRawConfig in
// resource_bigip_ltm_profile_ssl_client_unit_test.go and
// policyResourceDataWithRawConfig in resource_bigip_ltm_policy_unit_test.go.
func serverSslResourceDataWithRawConfig(t *testing.T, raw map[string]interface{}, rawConfig cty.Value) *schema.ResourceData {
	t.Helper()
	return NewTestResourceDataWithRawConfig(t, resourceBigipLtmProfileServerSsl(), raw, rawConfig)
}

// ---------------------------------------------------------------------
// Schema-shape assertion (no HTTP server needed)
// ---------------------------------------------------------------------

func TestResourceBigipLtmProfileServerSslSchema(t *testing.T) {
	r := resourceBigipLtmProfileServerSsl()

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
	if r.Schema["defaults_from"].Default != "/Common/serverssl" {
		t.Errorf("Expected 'defaults_from' to default to '/Common/serverssl', got %v", r.Schema["defaults_from"].Default)
	}
}

// ---------------------------------------------------------------------
// Direct CRUD-function unit tests (style 2)
// ---------------------------------------------------------------------

func TestUnitLtmProfileServerSslCreateReadUpdateDelete(t *testing.T) {
	name := "test-server-ssl"

	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl", func(w http.ResponseWriter, r *http.Request) {
		AssertRequestMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = fmt.Fprintf(w, `{"name":"%s","partition":"Common","defaultsFrom":"/Common/serverssl","alertTimeout":"indefinite","authenticate":"once","authenticateDepth":9,"c3dCaCert":"none","c3dCaKey":"none","cert":"/Common/default.crt","key":"/Common/default.key","chain":"none","ciphers":"DEFAULT","tmOptions":["Dont-Insert-Empty-Fragments"],"proxySsl":"disabled"}`, name)
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
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":          name,
		"defaults_from": "/Common/serverssl",
		"ciphers":       "DEFAULT",
	}, "")

	ctx := context.Background()

	diags := resourceBigipLtmProfileServerSslCreate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected create error: %v", diags)
	require.Equal(t, name, d.Id())

	diags = resourceBigipLtmProfileServerSslRead(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	require.Equal(t, "DEFAULT", d.Get("ciphers").(string))

	diags = resourceBigipLtmProfileServerSslUpdate(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected update error: %v", diags)

	diags = resourceBigipLtmProfileServerSslDelete(ctx, d, client)
	require.False(t, diags.HasError(), "unexpected delete error: %v", diags)
	require.Equal(t, "", d.Id())
}

func TestUnitLtmProfileServerSslRead_TmOptionsNone(t *testing.T) {
	name := "test-server-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","tmOptions":"none"}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileServerSslRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	tmOptions := d.Get("tm_options").(*schema.Set)
	require.Equal(t, 0, tmOptions.Len())
}

func TestUnitLtmProfileServerSslRead_TmOptionsNil(t *testing.T) {
	// Exercises the nil branch directly: obj.TmOptions is untyped nil when
	// the tmOptions key is absent from BIG-IP's JSON response entirely
	// (e.g. never configured, no defaults_from to inherit from), since
	// go-bigip declares TmOptions as interface{} with an omitempty JSON tag
	// and encoding/json leaves an absent field at its Go zero value (nil)
	// rather than invoking UnmarshalJSON. Read must not call
	// reflect.TypeOf(obj.TmOptions).Kind() unguarded in this case, since
	// reflect.TypeOf(nil) returns a nil *rtype and .Kind() on it panics.
	name := "test-server-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s"}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	require.NotPanics(t, func() {
		diags := resourceBigipLtmProfileServerSslRead(context.Background(), d, client)
		require.False(t, diags.HasError(), "unexpected read error: %v", diags)
	})
	tmOptions := d.Get("tm_options").(*schema.Set)
	require.Equal(t, 0, tmOptions.Len())
}

func TestUnitLtmProfileServerSslRead_TmOptionsLegacyString(t *testing.T) {
	// Exercises the reflect.String branch: some BIG-IP versions return
	// tmOptions as a legacy brace-delimited string instead of a JSON array.
	name := "test-server-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","tmOptions":"{ Dont-Insert-Empty-Fragments }"}`, name)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileServerSslRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "unexpected read error: %v", diags)
}

func TestUnitLtmProfileServerSslCreate_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"create failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	client.Teem = true
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-server-ssl",
	}, "")

	diags := resourceBigipLtmProfileServerSslCreate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileServerSslRead_Error(t *testing.T) {
	name := "test-server-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"read failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileServerSslRead(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileServerSslUpdate_Error(t *testing.T) {
	name := "test-server-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"update failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileServerSslUpdate(context.Background(), d, client)
	require.True(t, diags.HasError())
}

func TestUnitLtmProfileServerSslDelete_Error(t *testing.T) {
	name := "test-server-ssl"
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/ltm/profile/server-ssl/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"code":500,"message":"delete failed"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewUnitTestClient(server.URL)
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": name,
	}, name)

	diags := resourceBigipLtmProfileServerSslDelete(context.Background(), d, client)
	require.True(t, diags.HasError())
}

// ---------------------------------------------------------------------
// getServerSslConfig
// ---------------------------------------------------------------------

func TestGetServerSslConfig_Defaults(t *testing.T) {
	raw := map[string]interface{}{
		"name":    "test-server-ssl",
		"ciphers": "DEFAULT",
	}
	// ciphers/cipher_group are read via rawConfigAttrIsSet (see
	// getServerSslConfig's own comment), so this must use the raw-config-
	// aware helper -- plain NewTestResourceData leaves GetRawConfig()
	// null, which rawConfigAttrIsSet treats as "not set in the user's
	// .tf" regardless of what's in raw/state.
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"ciphers":      cty.StringVal("DEFAULT"),
		"cipher_group": cty.NullVal(cty.String),
		"tm_options":   cty.NullVal(cty.Set(cty.String)),
	})
	d := serverSslResourceDataWithRawConfig(t, raw, rawConfig)

	cfg := getServerSslConfig(d, &bigip.ServerSSLProfile{})
	require.Equal(t, "DEFAULT", cfg.Ciphers)
	require.Equal(t, "none", cfg.CipherGroup)
	require.Equal(t, "none", cfg.C3dCaCert)
	require.Equal(t, "none", cfg.C3dCaKey)
}

func TestGetServerSslConfig_SslForwardProxyEnabled(t *testing.T) {
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":                     "test-server-ssl",
		"ssl_forward_proxy":        "enabled",
		"ssl_forward_proxy_bypass": "",
	}, "")

	cfg := getServerSslConfig(d, &bigip.ServerSSLProfile{})
	require.Equal(t, "/Common/default.crt", cfg.ProxyCaCert)
	require.Equal(t, "/Common/default.key", cfg.ProxyCaKey)
	require.Equal(t, "disabled", cfg.SslForwardProxyBypass)
}

func TestGetServerSslConfig_SslC3dEnabledDefaults(t *testing.T) {
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name":        "test-server-ssl",
		"ssl_c3d":     "enabled",
		"c3d_ca_cert": "none",
		"c3d_ca_key":  "none",
	}, "")

	cfg := getServerSslConfig(d, &bigip.ServerSSLProfile{})
	require.Equal(t, "/Common/default.crt", cfg.C3dCaCert)
	require.Equal(t, "/Common/default.key", cfg.C3dCaKey)
	require.Equal(t, []string{"basic-constraints", "extended-key-usage", "key-usage", "subject-alternative-name"}, cfg.C3dCertExtensionIncludes)
}

func TestGetServerSslConfig_CipherGroupSet(t *testing.T) {
	raw := map[string]interface{}{
		"name":         "test-server-ssl",
		"cipher_group": "/Common/my-cipher-group",
	}
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"cipher_group": cty.StringVal("/Common/my-cipher-group"),
		"ciphers":      cty.NullVal(cty.String),
		"tm_options":   cty.NullVal(cty.Set(cty.String)),
	})
	d := serverSslResourceDataWithRawConfig(t, raw, rawConfig)

	cfg := getServerSslConfig(d, &bigip.ServerSSLProfile{})
	require.Equal(t, "/Common/my-cipher-group", cfg.CipherGroup)
	require.Equal(t, "none", cfg.Ciphers)
}

func TestGetServerSslConfig_TmOptionsAndC3dLists(t *testing.T) {
	raw := map[string]interface{}{
		"name":                           "test-server-ssl",
		"tm_options":                     []interface{}{"Dont-Insert-Empty-Fragments"},
		"c3d_cert_extension_custom_oids": []interface{}{"1.2.3.4"},
		"c3d_cert_extension_includes":    []interface{}{"basic-constraints"},
	}
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"tm_options":   cty.SetVal([]cty.Value{cty.StringVal("Dont-Insert-Empty-Fragments")}),
		"ciphers":      cty.NullVal(cty.String),
		"cipher_group": cty.NullVal(cty.String),
	})
	d := serverSslResourceDataWithRawConfig(t, raw, rawConfig)

	cfg := getServerSslConfig(d, &bigip.ServerSSLProfile{})
	require.Equal(t, []string{"Dont-Insert-Empty-Fragments"}, cfg.TmOptions)
	require.Equal(t, []string{"1.2.3.4"}, cfg.C3dCertExtensionCustomOids)
	require.Equal(t, []string{"basic-constraints"}, cfg.C3dCertExtensionIncludes)
}

// TestGetServerSslConfig_TmOptionsNotInRawConfigIsNotResent is the
// regression test for SFDC #01262589's bigip_ltm_profile_server_ssl analog:
// a child profile whose own .tf never sets tm_options (raw config attribute
// is null) must NOT have tm_options included in the outgoing Create/Update
// payload, even if prior state already has a non-empty tm_options value
// backfilled by Read from an inherited defaults_from value (which is
// exactly what happens after any unrelated attribute update triggers
// another Read/Update cycle -- see getServerSslConfig's tm_options handling
// and getClientSslConfig's identical comment on
// resource_bigip_ltm_profile_ssl_client.go for the full root-cause
// explanation). Before the fix, this scenario incorrectly set
// cfg.TmOptions to the backfilled value, causing BIG-IP to receive an
// explicit "options" line and permanently breaking inheritance for that
// attribute on the child profile -- reproduced directly against the
// pre-fix code (confirmed to fail before the getServerSslConfig fix,
// passes after).
func TestGetServerSslConfig_TmOptionsNotInRawConfigIsNotResent(t *testing.T) {
	raw := map[string]interface{}{
		"name": "test-server-ssl",
	}
	// The user's current .tf does not mention tm_options at all.
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"tm_options":   cty.NullVal(cty.Set(cty.String)),
		"ciphers":      cty.NullVal(cty.String),
		"cipher_group": cty.NullVal(cty.String),
	})
	d := serverSslResourceDataWithRawConfig(t, raw, rawConfig)

	// Simulate a prior Read having backfilled state with the parent
	// profile's inherited options, as happens whenever defaults_from is
	// set and BIG-IP's GET response reflects the effective (inherited)
	// tmOptions rather than only locally-configured ones.
	require.NoError(t, d.Set("tm_options", []interface{}{"no-ssl", "no-tlsv1.1"}))
	_, getOkResult := d.GetOk("tm_options")
	require.True(t, getOkResult, "GetOk is expected to see the backfilled state value -- this is the bug trigger getServerSslConfig must not rely on")

	cfg := getServerSslConfig(d, &bigip.ServerSSLProfile{})
	require.Empty(t, cfg.TmOptions, "tm_options must not be re-sent to BIG-IP when the user's .tf never set it, even if state has a backfilled inherited value")
}

// TestGetServerSslConfig_TmOptionsExplicitEmptyClearsInherited covers the
// opposite edge of the same rawConfigAttrIsSet guard: a practitioner who
// writes tm_options = [] explicitly in their .tf -- rather than omitting
// the attribute entirely -- is expressing clear intent to override any
// inherited tm_options with an empty set, not merely declining to set the
// attribute. rawConfigAttrIsSet must return true here (the attribute is
// present and known, just empty), and getServerSslConfig must still assign
// that empty slice to config.TmOptions so it is sent to BIG-IP as an
// explicit "options {}" clear -- as opposed to
// TestGetServerSslConfig_TmOptionsNotInRawConfigIsNotResent, where
// tm_options is absent from the raw config entirely (null, not an empty
// set) and must NOT be sent at all.
func TestGetServerSslConfig_TmOptionsExplicitEmptyClearsInherited(t *testing.T) {
	raw := map[string]interface{}{
		"name": "test-server-ssl",
	}
	// The user's current .tf explicitly sets tm_options = [], distinct from
	// omitting the attribute (which cty represents as a null set, per
	// TestGetServerSslConfig_TmOptionsNotInRawConfigIsNotResent above).
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"tm_options":   cty.SetValEmpty(cty.String),
		"ciphers":      cty.NullVal(cty.String),
		"cipher_group": cty.NullVal(cty.String),
	})
	d := serverSslResourceDataWithRawConfig(t, raw, rawConfig)
	require.True(t, rawConfigAttrIsSet(d, "tm_options"), "an explicit empty set in raw config is set/known, not null")

	cfg := getServerSslConfig(d, &bigip.ServerSSLProfile{})
	require.NotNil(t, cfg.TmOptions, "explicit tm_options = [] must still be sent to BIG-IP as an explicit clear, not silently dropped")
	require.Empty(t, cfg.TmOptions, "explicit tm_options = [] must clear to an empty list")
}

// TestGetServerSslConfig_RawConfigNilIsNotResent covers the other half of
// the rawConfigAttrIsSet guard: d.GetRawConfig() itself being the SDK-null
// value (rawConfig.IsNull() true for the whole config, not just a null
// tm_options attribute within an otherwise-present object). This is exactly
// what schema.TestResourceDataRaw (and thus NewTestResourceData) produces,
// and is also what happens for real on terraform destroy, where there is no
// current .tf config to read at all. rawConfigAttrIsSet must short-circuit
// on this before ever calling cty.Value.GetAttr, which panics on a null
// object value rather than returning null.
func TestGetServerSslConfig_RawConfigNilIsNotResent(t *testing.T) {
	r := resourceBigipLtmProfileServerSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-server-ssl",
	}, "")
	require.True(t, d.GetRawConfig().IsNull(), "schema.TestResourceDataRaw is expected to leave GetRawConfig() null")

	require.False(t, rawConfigAttrIsSet(d, "tm_options"))

	// Simulate a prior Read having backfilled state with an inherited
	// value, as in TestGetServerSslConfig_TmOptionsNotInRawConfigIsNotResent,
	// to prove the nil-rawConfig case is guarded independently of state.
	require.NoError(t, d.Set("tm_options", []interface{}{"no-ssl", "no-tlsv1.1"}))
	_, getOkResult := d.GetOk("tm_options")
	require.True(t, getOkResult, "GetOk is expected to see the backfilled state value -- this is the bug trigger getServerSslConfig must not rely on")

	require.NotPanics(t, func() {
		cfg := getServerSslConfig(d, &bigip.ServerSSLProfile{})
		require.Empty(t, cfg.TmOptions, "tm_options must not be re-sent to BIG-IP when GetRawConfig() is null")
	})
}
