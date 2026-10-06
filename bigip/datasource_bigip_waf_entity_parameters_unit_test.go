/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// wafEntityParameterResourceData builds a *schema.ResourceData for the
// bigip_waf_entity_parameter data source using schema.TestResourceDataRaw.
func wafEntityParameterResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipWafEntityParameter(), raw)
}

// TestDataSourceBigipWafEntityParameterSchema exercises the schema
// definition of the bigip_waf_entity_parameter data source without needing
// any BIG-IP connection.
func TestDataSourceBigipWafEntityParameterSchema(t *testing.T) {
	r := dataSourceBigipWafEntityParameter()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	nameSchema, ok := r.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist in schema")
	}
	if !nameSchema.Required {
		t.Error("Expected field 'name' to be required")
	}
	if nameSchema.Type != schema.TypeString {
		t.Errorf("Expected field 'name' to be TypeString, got %v", nameSchema.Type)
	}

	levelSchema, ok := r.Schema["level"]
	if !ok {
		t.Fatal("Expected field 'level' to exist in schema")
	}
	if levelSchema.ValidateFunc == nil {
		t.Error("Expected field 'level' to have a ValidateFunc")
	} else {
		_, errs := levelSchema.ValidateFunc("invalid-level", "level")
		assert.NotEmpty(t, errs, "expected validation error for a level outside {global, url, flow}")
		_, errs = levelSchema.ValidateFunc("url", "level")
		assert.Empty(t, errs, "expected no validation error for 'url'")
	}

	boolFields := []string{
		"allow_empty_type", "allow_repeated_parameter_name", "attack_signatures_check",
		"check_max_value_length", "check_min_value_length", "enable_regular_expression",
		"is_base64", "is_cookie", "is_header", "mandatory",
		"metachars_on_parameter_value_check", "perform_staging", "sensitive_parameter",
	}
	for _, field := range boolFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Optional {
			t.Errorf("Expected field '%s' to be Optional", field)
		}
		if s.Type != schema.TypeBool {
			t.Errorf("Expected field '%s' to be TypeBool, got %v", field, s.Type)
		}
	}

	intFields := []string{"max_value_length", "min_value_length"}
	for _, field := range intFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if s.Type != schema.TypeInt {
			t.Errorf("Expected field '%s' to be TypeInt, got %v", field, s.Type)
		}
	}

	urlSchema, ok := r.Schema["url"]
	if !ok {
		t.Fatal("Expected field 'url' to exist in schema")
	}
	if urlSchema.Type != schema.TypeList {
		t.Errorf("Expected field 'url' to be TypeList, got %v", urlSchema.Type)
	}
	if urlSchema.MaxItems != 1 {
		t.Errorf("Expected field 'url' to have MaxItems 1, got %d", urlSchema.MaxItems)
	}

	sigOverridesSchema, ok := r.Schema["signature_overrides_disable"]
	if !ok {
		t.Fatal("Expected field 'signature_overrides_disable' to exist in schema")
	}
	if sigOverridesSchema.Type != schema.TypeList {
		t.Errorf("Expected field 'signature_overrides_disable' to be TypeList, got %v", sigOverridesSchema.Type)
	}

	jsonSchema, ok := r.Schema["json"]
	if !ok {
		t.Fatal("Expected field 'json' to exist in schema")
	}
	if !jsonSchema.Optional {
		t.Error("Expected field 'json' to be Optional")
	}
	if !jsonSchema.Computed {
		t.Error("Expected field 'json' to be Computed")
	}
}

// TestDataSourceBigipWafEntityParameterReadMinimal covers the Read happy
// path with only the required "name" field set: Read never calls the
// BIG-IP, so it must succeed purely from schema defaults, set the resource
// ID to the parameter name, and build a JSON payload reflecting those
// defaults.
func TestDataSourceBigipWafEntityParameterReadMinimal(t *testing.T) {
	d := wafEntityParameterResourceData(t, map[string]interface{}{
		"name": "test-param",
	})

	diags := dataSourceBigipWafEntityParameterRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "test-param", d.Id())
	jsonStr := d.Get("json").(string)
	assert.Contains(t, jsonStr, `"name":"test-param"`)
}

// TestDataSourceBigipWafEntityParameterReadLevelUrlMissingUrl covers the
// branch where level is set to "url" but no "url" block is supplied: Read
// must fail with a descriptive error before calling getEPConfig.
func TestDataSourceBigipWafEntityParameterReadLevelUrlMissingUrl(t *testing.T) {
	d := wafEntityParameterResourceData(t, map[string]interface{}{
		"name":  "test-param",
		"level": "url",
	})

	diags := dataSourceBigipWafEntityParameterRead(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error when level is 'url' but no url block is set")
	assert.Contains(t, diags[0].Summary, "url object must be specificed")
}

// TestDataSourceBigipWafEntityParameterReadLevelUrlWithUrl covers the
// branch where level is "url" and a url block is supplied: Read must
// succeed and getEPConfig must populate the nested URL fields in the JSON
// output.
func TestDataSourceBigipWafEntityParameterReadLevelUrlWithUrl(t *testing.T) {
	d := wafEntityParameterResourceData(t, map[string]interface{}{
		"name":  "test-param",
		"level": "url",
		"url": []interface{}{
			map[string]interface{}{
				"name":     "/some-url",
				"method":   "GET",
				"protocol": "http",
				"type":     "wildcard",
			},
		},
	})

	diags := dataSourceBigipWafEntityParameterRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	jsonStr := d.Get("json").(string)
	assert.Contains(t, jsonStr, `"level":"url"`)
	assert.Contains(t, jsonStr, `"url":{`)
	assert.Contains(t, jsonStr, `/some-url`)
}

// TestDataSourceBigipWafEntityParameterReadFull covers Read/getEPConfig
// with the full set of optional scalar fields populated (excluding url,
// covered separately above), to exercise each of getEPConfig's field
// assignment branches.
func TestDataSourceBigipWafEntityParameterReadFull(t *testing.T) {
	d := wafEntityParameterResourceData(t, map[string]interface{}{
		"name":                               "full-param",
		"description":                        "a full parameter",
		"type":                               "explicit",
		"value_type":                         "static",
		"allow_repeated_parameter_name":      true,
		"attack_signatures_check":            true,
		"check_max_value_length":             true,
		"check_min_value_length":             true,
		"max_value_length":                   100,
		"min_value_length":                   1,
		"data_type":                          "alpha-numeric",
		"enable_regular_expression":          true,
		"is_base64":                          true,
		"is_cookie":                          true,
		"is_header":                          true,
		"level":                              "global",
		"mandatory":                          true,
		"metachars_on_parameter_value_check": true,
		"parameter_location":                 "query",
		"perform_staging":                    true,
		"sensitive_parameter":                true,
		"signature_overrides_disable":        []interface{}{111, 222},
	})

	diags := dataSourceBigipWafEntityParameterRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	jsonStr := d.Get("json").(string)
	assert.Contains(t, jsonStr, `"description":"a full parameter"`)
	assert.Contains(t, jsonStr, `"type":"explicit"`)
	assert.Contains(t, jsonStr, `"valueType":"static"`)
	assert.Contains(t, jsonStr, `"allowRepeatedParameterName":true`)
	assert.Contains(t, jsonStr, `"attackSignaturesCheck":true`)
	assert.Contains(t, jsonStr, `"checkMaxValueLength":true`)
	assert.Contains(t, jsonStr, `"checkMinValueLength":true`)
	assert.Contains(t, jsonStr, `"maximumLength":100`)
	assert.Contains(t, jsonStr, `"minimumLength":1`)
	assert.Contains(t, jsonStr, `"dataType":"alpha-numeric"`)
	assert.Contains(t, jsonStr, `"enableRegularExpression":true`)
	assert.Contains(t, jsonStr, `"isBase64":true`)
	assert.Contains(t, jsonStr, `"isCookie":true`)
	assert.Contains(t, jsonStr, `"isHeader":true`)
	assert.Contains(t, jsonStr, `"level":"global"`)
	assert.Contains(t, jsonStr, `"mandatory":true`)
	assert.Contains(t, jsonStr, `"metacharsOnParameterValueCheck":true`)
	assert.Contains(t, jsonStr, `"parameterLocation":"query"`)
	assert.Contains(t, jsonStr, `"performStaging":true`)
	assert.Contains(t, jsonStr, `"sensitiveParameter":true`)
	assert.Contains(t, jsonStr, `"signatureOverrides"`)
	assert.Contains(t, jsonStr, `"signatureId":111`)
	assert.Contains(t, jsonStr, `"signatureId":222`)
}

// TestDataSourceBigipWafEntityParameterReadMaxMinLengthIgnoredForNonAlphaNumeric
// covers the branch in getEPConfig where check_max_value_length/
// check_min_value_length are true but data_type is not "alpha-numeric": in
// that case MaximumLength/MinimumLength must NOT be set on the resulting
// Parameter, even though max_value_length/min_value_length are populated.
func TestDataSourceBigipWafEntityParameterReadMaxMinLengthIgnoredForNonAlphaNumeric(t *testing.T) {
	d := wafEntityParameterResourceData(t, map[string]interface{}{
		"name":                   "test-param",
		"check_max_value_length": true,
		"check_min_value_length": true,
		"max_value_length":       100,
		"min_value_length":       1,
		"data_type":              "integer",
	})

	diags := dataSourceBigipWafEntityParameterRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	jsonStr := d.Get("json").(string)
	assert.NotContains(t, jsonStr, `"maximumLength"`, "maximumLength should be omitted when data_type is not alpha-numeric")
	assert.NotContains(t, jsonStr, `"minimumLength"`, "minimumLength should be omitted when data_type is not alpha-numeric")
}

// TestDataSourceBigipWafEntityParameterReadAllowEmptyValueFieldMismatch
// documents existing behavior in getEPConfig: the schema field is named
// "allow_empty_type" but getEPConfig reads d.Get("allow_empty_value")
// (a key that does not exist in this resource's schema). d.Get on an
// unknown key returns nil, so the `!= nil` guard is always false and
// ep.AllowEmptyValue is never set from user input, regardless of what
// allow_empty_type is set to. This test pins down that current behavior;
// it is not asserting this is the intended behavior, only that it doesn't
// panic and produces no allowEmptyValue key in the output (since
// AllowEmptyValue stays at its zero value and the struct field has
// omitempty).
func TestDataSourceBigipWafEntityParameterReadAllowEmptyValueFieldMismatch(t *testing.T) {
	d := wafEntityParameterResourceData(t, map[string]interface{}{
		"name":             "test-param",
		"allow_empty_type": true,
	})

	diags := dataSourceBigipWafEntityParameterRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	jsonStr := d.Get("json").(string)
	assert.NotContains(t, jsonStr, "allowEmptyValue", "allow_empty_type has no effect on the output due to the field-name mismatch in getEPConfig")
}
