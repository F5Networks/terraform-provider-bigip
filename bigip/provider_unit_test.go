/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// providerResourceData builds a *schema.ResourceData for the provider's own
// schema (Provider().Schema), pre-filled with sane defaults for every field
// providerConfigure reads, then overridden by raw.
func providerResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	base := map[string]interface{}{
		"address":                "",
		"port":                   "",
		"username":               "",
		"password":               "",
		"token_value":            "",
		"token_auth":             false,
		"validate_certs_disable": true,
		"trusted_cert_path":      "",
		"teem_disable":           false,
		"login_ref":              "tmos",
		"api_timeout":            60,
		"token_timeout":          1200,
		"api_retries":            1,
	}
	for k, v := range raw {
		base[k] = v
	}
	return schema.TestResourceDataRaw(t, Provider().Schema, base)
}

// ---------------------------------------------------------------------
// providerConfigure
// ---------------------------------------------------------------------

func TestProviderConfigure_SuccessInsecure(t *testing.T) {
	srv := newSelfIPServer(t)

	d := providerResourceData(t, map[string]interface{}{
		"address":                srv.URL,
		"username":               "admin",
		"password":               "secret",
		"validate_certs_disable": true,
		"teem_disable":           true,
	})

	cfg, diags := providerConfigure(d, "1.5.0")

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	client, ok := cfg.(*bigip.BigIP)
	require.True(t, ok, "expected *bigip.BigIP, got %T", cfg)
	assert.True(t, client.Transport.TLSClientConfig.InsecureSkipVerify)
	assert.True(t, client.Teem)
	assert.Contains(t, client.UserAgent, "Terraform/1.5.0")
	assert.Contains(t, client.UserAgent, "terraform-provider-bigip/")
}

func TestProviderConfigure_MissingTerraformVersionDefaultsCompatible(t *testing.T) {
	srv := newSelfIPServer(t)

	d := providerResourceData(t, map[string]interface{}{
		"address":                srv.URL,
		"username":               "admin",
		"password":               "secret",
		"validate_certs_disable": true,
	})

	// providerConfigure itself doesn't special-case an empty terraformVersion
	// (that fallback lives in Provider()'s ConfigureContextFunc closure), so
	// this just verifies whatever version string is passed ends up in the
	// UserAgent verbatim.
	cfg, diags := providerConfigure(d, "0.11+compatible")

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	client := cfg.(*bigip.BigIP)
	assert.Contains(t, client.UserAgent, "Terraform/0.11+compatible")
}

func TestProviderConfigure_TokenAuthSetsLoginReference(t *testing.T) {
	srv := newTLSSelfIPServer(t)
	certPath := writeServerCertPEM(t, srv)

	d := providerResourceData(t, map[string]interface{}{
		"address":                srv.URL,
		"username":               "admin",
		"password":               "secret",
		"validate_certs_disable": false,
		"trusted_cert_path":      certPath,
		"token_auth":             true,
		"login_ref":              "tmos",
	})

	cfg, diags := providerConfigure(d, "1.5.0")

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	client := cfg.(*bigip.BigIP)
	// A non-empty Token on the returned client confirms the token-auth path
	// (NewTokenSession, driven by config.LoginReference being set) actually
	// ran, rather than the plain NewSession/basic-auth path.
	assert.NotEmpty(t, client.Token)
	assert.False(t, client.Transport.TLSClientConfig.InsecureSkipVerify)
}

func TestProviderConfigure_CertVerifyEnabledMissingTrustedCertPath(t *testing.T) {
	d := providerResourceData(t, map[string]interface{}{
		"address":                "https://192.0.2.30",
		"username":               "admin",
		"password":               "secret",
		"validate_certs_disable": false,
		"trusted_cert_path":      "",
	})

	cfg, diags := providerConfigure(d, "1.5.0")

	require.True(t, diags.HasError())
	assert.Contains(t, diags[0].Summary, "valid Trust Certificate path not provided")
	assert.Nil(t, cfg)
}

func TestProviderConfigure_ClientError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	d := providerResourceData(t, map[string]interface{}{
		"address":                srv.URL,
		"username":               "admin",
		"password":               "secret",
		"validate_certs_disable": true,
	})

	cfg, diags := providerConfigure(d, "1.5.0")

	require.True(t, diags.HasError())
	// providerConfigure returns immediately on a Client() error, before ever
	// reaching the UserAgent/Teem/TLS setup block -- so even though Client()
	// itself returns a non-nil client on a validation failure, providerConfigure
	// passes it straight through unmodified (UserAgent left empty).
	require.NotNil(t, cfg)
	client := cfg.(*bigip.BigIP)
	assert.Empty(t, client.UserAgent)
}

// ---------------------------------------------------------------------
// Provider() top-level construction
// ---------------------------------------------------------------------

func TestProvider_InternalValidateAndConfigureContextFunc(t *testing.T) {
	p := Provider()
	require.NoError(t, p.InternalValidate())
	require.NotNil(t, p.ConfigureContextFunc)

	srv := newSelfIPServer(t)
	d := providerResourceData(t, map[string]interface{}{
		"address":                srv.URL,
		"username":               "admin",
		"password":               "secret",
		"validate_certs_disable": true,
	})

	cfg, diags := p.ConfigureContextFunc(nil, d) //nolint:staticcheck // context arg unused by providerConfigure
	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	require.NotNil(t, cfg)
}

// ---------------------------------------------------------------------
// Small pure helper functions
// ---------------------------------------------------------------------

func TestMakeStringList(t *testing.T) {
	in := []string{"a", "b", "c"}
	out := makeStringList(&in)
	require.Equal(t, []interface{}{"a", "b", "c"}, out)
}

func TestMakeStringSet(t *testing.T) {
	in := []string{"a", "b"}
	out := makeStringSet(&in)
	require.Equal(t, 2, out.Len())
	require.True(t, out.Contains("a"))
	require.True(t, out.Contains("b"))
}

func TestListToStringSlice(t *testing.T) {
	in := []interface{}{"x", "y"}
	out := listToStringSlice(in)
	require.Equal(t, []string{"x", "y"}, out)
}

func TestListToIntSlice(t *testing.T) {
	in := []interface{}{1, 2, 3}
	out := listToIntSlice(in)
	require.Equal(t, []int{1, 2, 3}, out)
}

func TestSetToStringSlice(t *testing.T) {
	s := schema.NewSet(schema.HashString, []interface{}{"a", "b"})
	out := setToStringSlice(s)
	require.ElementsMatch(t, []string{"a", "b"}, out)
}

func TestSetToInterfaceSlice(t *testing.T) {
	s := schema.NewSet(schema.HashString, []interface{}{"a", "b"})
	out := setToInterfaceSlice(s)
	require.Len(t, out, 2)
}

func TestGetVersion(t *testing.T) {
	require.Equal(t, ProviderVersion, getVersion())
}

func TestHashForState(t *testing.T) {
	require.Equal(t, "", hashForState(""))
	require.NotEmpty(t, hashForState("some-value"))
	// Leading/trailing whitespace is trimmed before hashing.
	require.Equal(t, hashForState("some-value"), hashForState("  some-value  "))
}

func TestToCamelCase(t *testing.T) {
	require.Equal(t, "FooBar", toCamelCase("foo_bar"))
	require.Equal(t, "Foo", toCamelCase("foo"))
}

func TestToSnakeCase(t *testing.T) {
	require.Equal(t, "foo_bar", toSnakeCase("FooBar"))
	require.Equal(t, "foo", toSnakeCase("Foo"))
}

// ---------------------------------------------------------------------
// mapEntity
// ---------------------------------------------------------------------

type mapEntityTestStruct struct {
	Name       string
	MaxLength  int
	Tags       []string
	CustomName string
}

func TestMapEntity_DirectFieldMatch(t *testing.T) {
	obj := &mapEntityTestStruct{}
	mapEntity(map[string]interface{}{
		"Name": "test",
	}, obj)
	require.Equal(t, "test", obj.Name)
}

func TestMapEntity_SliceField(t *testing.T) {
	obj := &mapEntityTestStruct{}
	mapEntity(map[string]interface{}{
		"Tags": []interface{}{"a", "b"},
	}, obj)
	require.Equal(t, []string{"a", "b"}, obj.Tags)
}

func TestMapEntity_SnakeCaseFallback(t *testing.T) {
	obj := &mapEntityTestStruct{}
	mapEntity(map[string]interface{}{
		"custom_name": "converted",
	}, obj)
	require.Equal(t, "converted", obj.CustomName)
}

// ---------------------------------------------------------------------
// ctyValIsSet / ctyObjectToMap
// ---------------------------------------------------------------------

func TestCtyValIsSet(t *testing.T) {
	require.False(t, ctyValIsSet(cty.NullVal(cty.String)))
	require.False(t, ctyValIsSet(cty.UnknownVal(cty.String)))
	require.True(t, ctyValIsSet(cty.StringVal("value")))
}

func TestCtyObjectToMap_NotSet(t *testing.T) {
	require.Nil(t, ctyObjectToMap(cty.NullVal(cty.EmptyObject), nil))
}

func TestCtyObjectToMap_WrongType(t *testing.T) {
	require.Nil(t, ctyObjectToMap(cty.StringVal("not-an-object"), nil))
}

func TestCtyObjectToMap_AllTypes(t *testing.T) {
	val := cty.ObjectVal(map[string]cty.Value{
		"str":       cty.StringVal("hello"),
		"boolTrue":  cty.True,
		"intNum":    cty.NumberIntVal(42),
		"floatNum":  cty.NumberFloatVal(3.14),
		"strList":   cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
		"nullField": cty.NullVal(cty.String),
	})

	result := ctyObjectToMap(val, nil)

	require.Equal(t, "hello", result["str"])
	require.Equal(t, true, result["boolTrue"])
	require.Equal(t, 42, result["intNum"])
	require.InDelta(t, 3.14, result["floatNum"].(float64), 0.0001)
	require.Equal(t, []interface{}{"a", "b"}, result["strList"])
	require.NotContains(t, result, "nullField")
}

func TestCtyObjectToMap_NullFieldUsesSchemaDefault(t *testing.T) {
	val := cty.ObjectVal(map[string]cty.Value{
		"withDefault": cty.NullVal(cty.String),
	})
	schemaMap := map[string]*schema.Schema{
		"withDefault": {
			Type:    schema.TypeString,
			Default: "default-value",
		},
	}

	result := ctyObjectToMap(val, schemaMap)

	require.Equal(t, "default-value", result["withDefault"])
}

func TestCtyObjectToMap_MapType(t *testing.T) {
	val := cty.MapVal(map[string]cty.Value{
		"key1": cty.StringVal("value1"),
	})

	result := ctyObjectToMap(val, nil)
	require.Equal(t, "value1", result["key1"])
}

func TestCtyObjectToMap_UnhandledType(t *testing.T) {
	// A list of a non-string element type isn't handled by any case, so it's
	// logged and simply omitted from the result rather than causing an error.
	val := cty.ObjectVal(map[string]cty.Value{
		"numList": cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2)}),
	})

	result := ctyObjectToMap(val, nil)
	require.NotContains(t, result, "numList")
}

// ---------------------------------------------------------------------
// rawConfigAttrIsSet
// ---------------------------------------------------------------------

func TestRawConfigAttrIsSet_TopLevelAttr(t *testing.T) {
	// The supported case: attr is a single top-level attribute name.
	setRawConfig := cty.ObjectVal(map[string]cty.Value{
		"tm_options": cty.SetVal([]cty.Value{cty.StringVal("no-tlsv1.3")}),
	})
	d := clientSslResourceDataWithRawConfig(t, map[string]interface{}{
		"name":       "test-client-ssl",
		"tm_options": []interface{}{"no-tlsv1.3"},
	}, setRawConfig)
	require.True(t, rawConfigAttrIsSet(d, "tm_options"))

	nullRawConfig := cty.ObjectVal(map[string]cty.Value{
		"tm_options": cty.NullVal(cty.Set(cty.String)),
	})
	d = clientSslResourceDataWithRawConfig(t, map[string]interface{}{
		"name": "test-client-ssl",
	}, nullRawConfig)
	require.False(t, rawConfigAttrIsSet(d, "tm_options"))
}

func TestRawConfigAttrIsSet_RawConfigNull(t *testing.T) {
	// d.GetRawConfig() itself null (the whole config, not just one null
	// attribute) -- e.g. what schema.TestResourceDataRaw/NewTestResourceData
	// produce, and what happens for real on terraform destroy. Must
	// short-circuit to false rather than reach GetAttr at all.
	r := resourceBigipLtmProfileClientSsl()
	d := NewTestResourceData(t, r, map[string]interface{}{
		"name": "test-client-ssl",
	}, "")
	require.True(t, d.GetRawConfig().IsNull())
	require.False(t, rawConfigAttrIsSet(d, "tm_options"))
}

// TestRawConfigAttrIsSet_DottedPathPanics locks in the LIMITATION documented
// on rawConfigAttrIsSet: attr must be a single top-level attribute name.
// cty.Value.GetAttr takes a literal attribute name, not a path expression --
// passing a dotted string does not "just return false" the way an absent
// top-level name might be assumed to; it panics with "value has no
// attribute of that name", since cty never treats "." as a path separator
// within GetAttr. Any future caller tempted to pass a dotted/nested path
// (e.g. "block.0.nested_attr") needs to extend this helper first (see its
// doc comment) rather than assume it degrades gracefully.
func TestRawConfigAttrIsSet_DottedPathPanics(t *testing.T) {
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"tm_options": cty.SetVal([]cty.Value{cty.StringVal("no-tlsv1.3")}),
	})
	d := clientSslResourceDataWithRawConfig(t, map[string]interface{}{
		"name":       "test-client-ssl",
		"tm_options": []interface{}{"no-tlsv1.3"},
	}, rawConfig)

	require.Panics(t, func() {
		rawConfigAttrIsSet(d, "tm_options.0.nested")
	}, "rawConfigAttrIsSet is documented as top-level-only; a dotted path must panic via cty.GetAttr, not silently return false")
}
