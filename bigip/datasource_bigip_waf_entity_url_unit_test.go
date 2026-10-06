/*
Copyright 2022 F5 Networks Inc.
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

// wafEntityUrlResourceData builds a *schema.ResourceData for the
// bigip_waf_entity_url data source using schema.TestResourceDataRaw.
func wafEntityUrlResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return newDatasourceTestResourceData(t, dataSourceBigipWafEntityUrl(), raw)
}

// TestDataSourceBigipWafEntityUrlSchema exercises the schema definition of
// the bigip_waf_entity_url data source without needing any BIG-IP
// connection.
func TestDataSourceBigipWafEntityUrlSchema(t *testing.T) {
	r := dataSourceBigipWafEntityUrl()

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
	if !nameSchema.ForceNew {
		t.Error("Expected field 'name' to be ForceNew")
	}

	typeSchema, ok := r.Schema["type"]
	if !ok {
		t.Fatal("Expected field 'type' to exist in schema")
	}
	if typeSchema.Default != "wildcard" {
		t.Errorf("Expected field 'type' to default to 'wildcard', got %v", typeSchema.Default)
	}
	if typeSchema.ValidateFunc == nil {
		t.Error("Expected field 'type' to have a ValidateFunc")
	} else {
		_, errs := typeSchema.ValidateFunc("invalid-type", "type")
		assert.NotEmpty(t, errs, "expected validation error for a type outside {explicit, wildcard}")
		_, errs = typeSchema.ValidateFunc("explicit", "type")
		assert.Empty(t, errs, "expected no validation error for 'explicit'")
	}

	protocolSchema, ok := r.Schema["protocol"]
	if !ok {
		t.Fatal("Expected field 'protocol' to exist in schema")
	}
	if protocolSchema.Default != "http" {
		t.Errorf("Expected field 'protocol' to default to 'http', got %v", protocolSchema.Default)
	}

	methodSchema, ok := r.Schema["method"]
	if !ok {
		t.Fatal("Expected field 'method' to exist in schema")
	}
	if methodSchema.Default != "*" {
		t.Errorf("Expected field 'method' to default to '*', got %v", methodSchema.Default)
	}

	jsonSchema, ok := r.Schema["json"]
	if !ok {
		t.Fatal("Expected field 'json' to exist in schema")
	}
	if !jsonSchema.Computed {
		t.Error("Expected field 'json' to be Computed")
	}

	methodOverridesSchema, ok := r.Schema["method_overrides"]
	if !ok {
		t.Fatal("Expected field 'method_overrides' to exist in schema")
	}
	if methodOverridesSchema.Type != schema.TypeList {
		t.Errorf("Expected field 'method_overrides' to be TypeList, got %v", methodOverridesSchema.Type)
	}

	corsSchema, ok := r.Schema["cross_origin_requests_enforcement"]
	if !ok {
		t.Fatal("Expected field 'cross_origin_requests_enforcement' to exist in schema")
	}
	if corsSchema.Type != schema.TypeList {
		t.Errorf("Expected field 'cross_origin_requests_enforcement' to be TypeList, got %v", corsSchema.Type)
	}

	sigOverridesSchema, ok := r.Schema["signature_overrides_disable"]
	if !ok {
		t.Fatal("Expected field 'signature_overrides_disable' to exist in schema")
	}
	if sigOverridesSchema.Type != schema.TypeList {
		t.Errorf("Expected field 'signature_overrides_disable' to be TypeList, got %v", sigOverridesSchema.Type)
	}
}

// TestDataSourceBigipWafEntityUrlReadMinimal covers the Read happy path
// with only the required "name" field set: Read never calls the BIG-IP
// (meta is explicitly ignored), so it must succeed purely from schema
// defaults, set the resource ID to the URL name, and build a JSON payload
// reflecting those defaults.
func TestDataSourceBigipWafEntityUrlReadMinimal(t *testing.T) {
	d := wafEntityUrlResourceData(t, map[string]interface{}{
		"name": "/test-url",
	})

	diags := dataSourceBigipWafEntityUrlRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/test-url", d.Id())
	jsonStr := d.Get("json").(string)
	assert.Contains(t, jsonStr, `"name":"/test-url"`)
	assert.Contains(t, jsonStr, `"type":"wildcard"`)
	assert.Contains(t, jsonStr, `"protocol":"http"`)
	assert.Contains(t, jsonStr, `"method":"*"`)
}

// TestDataSourceBigipWafEntityUrlReadFull covers Read with every optional
// field populated, including method_overrides,
// cross_origin_requests_enforcement, and signature_overrides_disable, to
// exercise the loops that build those nested structures.
func TestDataSourceBigipWafEntityUrlReadFull(t *testing.T) {
	d := wafEntityUrlResourceData(t, map[string]interface{}{
		"name":            "/full-url",
		"description":     "a full url entity",
		"type":            "explicit",
		"protocol":        "https",
		"method":          "GET",
		"perform_staging": true,
		"method_overrides": []interface{}{
			map[string]interface{}{
				"allow":  true,
				"method": "POST",
			},
			map[string]interface{}{
				"allow":  false,
				"method": "DELETE",
			},
		},
		"cross_origin_requests_enforcement": []interface{}{
			map[string]interface{}{
				"include_subdomains": true,
				"origin_name":        "example.com",
				"origin_port":        "443",
				"origin_protocol":    "https",
			},
		},
		"signature_overrides_disable": []interface{}{1111, 2222},
	})

	diags := dataSourceBigipWafEntityUrlRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "/full-url", d.Id())
	jsonStr := d.Get("json").(string)
	assert.Contains(t, jsonStr, `"description":"a full url entity"`)
	assert.Contains(t, jsonStr, `"type":"explicit"`)
	assert.Contains(t, jsonStr, `"protocol":"https"`)
	assert.Contains(t, jsonStr, `"performStaging":true`)
	assert.Contains(t, jsonStr, `"methodOverrides"`)
	assert.Contains(t, jsonStr, `"method":"POST"`)
	assert.Contains(t, jsonStr, `"method":"DELETE"`)
	assert.Contains(t, jsonStr, `"signatureOverrides"`)
	assert.Contains(t, jsonStr, `"signatureId":1111`)
	assert.Contains(t, jsonStr, `"signatureId":2222`)
	assert.Contains(t, jsonStr, `"enforcementMode":"enforce"`)
	assert.Contains(t, jsonStr, `"originName":"example.com"`)
}

// TestDataSourceBigipWafEntityUrlReadNoCrossOrigin covers the branch where
// cross_origin_requests_enforcement is empty (the default): the
// HTML5CrossOriginRequestsEnforcement.EnforcementMode must remain unset
// rather than being forced to "enforce".
func TestDataSourceBigipWafEntityUrlReadNoCrossOrigin(t *testing.T) {
	d := wafEntityUrlResourceData(t, map[string]interface{}{
		"name": "/no-cors-url",
	})

	diags := dataSourceBigipWafEntityUrlRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	jsonStr := d.Get("json").(string)
	assert.NotContains(t, jsonStr, `"enforcementMode":"enforce"`)
}
