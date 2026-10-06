/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDataSourceBigipFastConsulServiceDiscoverySchema exercises the schema
// definition of the bigip_fast_consul_service_discovery data source.
func TestDataSourceBigipFastConsulServiceDiscoverySchema(t *testing.T) {
	ds := dataSourceBigipFastConsulServiceDiscovery()

	if ds.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if ds.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	requiredFields := []string{"uri", "port"}
	for _, field := range requiredFields {
		s, ok := ds.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}

	if ds.Schema["consul_sd_json"].Computed != true {
		t.Error("Expected 'consul_sd_json' to be Computed")
	}
	if ds.Schema["type"].Default != "consul" {
		t.Errorf("Expected 'type' to default to 'consul', got %v", ds.Schema["type"].Default)
	}
	if ds.Schema["address_realm"].Default != "private" {
		t.Errorf("Expected 'address_realm' to default to 'private', got %v", ds.Schema["address_realm"].Default)
	}
	if ds.Schema["undetectable_action"].Default != "remove" {
		t.Errorf("Expected 'undetectable_action' to default to 'remove', got %v", ds.Schema["undetectable_action"].Default)
	}
	if ds.Schema["credential_update"].Default != false {
		t.Errorf("Expected 'credential_update' to default to false, got %v", ds.Schema["credential_update"].Default)
	}
	if ds.Schema["reject_unauthorized"].Default != true {
		t.Errorf("Expected 'reject_unauthorized' to default to true, got %v", ds.Schema["reject_unauthorized"].Default)
	}
}

func consulServiceDiscoveryResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	ds := dataSourceBigipFastConsulServiceDiscovery()
	return schema.TestResourceDataRaw(t, ds.Schema, raw)
}

// TestGetConsulSDConfigDefaults covers getConsulSDConfig with only the
// required fields set, verifying schema defaults. reject_unauthorized
// defaults to true, and even though the underlying struct field DOES have
// omitempty (like credential_update), true is a non-zero bool value, so
// omitempty has no effect here and the field must still be present.
func TestGetConsulSDConfigDefaults(t *testing.T) {
	d := consulServiceDiscoveryResourceData(t, map[string]interface{}{
		"uri":  "https://consul.example.com",
		"port": 8500,
	})

	config, err := getConsulSDConfig(d)

	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(config), &parsed))

	assert.Equal(t, "consul", parsed["sd_type"])
	assert.Equal(t, float64(8500), parsed["sd_port"])
	assert.Equal(t, "https://consul.example.com", parsed["sd_uri"])
	assert.Equal(t, "private", parsed["sd_addressRealm"])
	assert.Equal(t, "remove", parsed["sd_undetectableAction"])
	assert.Equal(t, true, parsed["sd_rejectUnauthorized"], "expected sd_rejectUnauthorized to be present (true) since its default is true and omitempty only drops the zero value (false)")
	_, hasCredentialUpdate := parsed["sd_credentialUpdate"]
	assert.False(t, hasCredentialUpdate, "expected sd_credentialUpdate to be omitted when false")
	for _, key := range []string{"sd_encodedToken", "sd_jmesPathQuery", "sd_minimumMonitors", "sd_trustCA", "sd_updateInterval"} {
		_, ok := parsed[key]
		assert.False(t, ok, "expected %s to be omitted when unset", key)
	}
}

// TestGetConsulSDConfigAllFields covers getConsulSDConfig with every
// optional field populated, including reject_unauthorized set to false
// (its non-default value).
func TestGetConsulSDConfigAllFields(t *testing.T) {
	d := consulServiceDiscoveryResourceData(t, map[string]interface{}{
		"type":                "consul",
		"uri":                 "https://consul.example.com",
		"port":                8500,
		"address_realm":       "public",
		"undetectable_action": "revoke",
		"credential_update":   true,
		"encoded_token":       "dG9rZW4=",
		"jmes_path_query":     "[*]",
		"reject_unauthorized": false,
		"trust_ca":            "-----BEGIN CERTIFICATE-----",
		"minimum_monitors":    "1",
		"update_interval":     "60",
	})

	config, err := getConsulSDConfig(d)

	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(config), &parsed))

	assert.Equal(t, "public", parsed["sd_addressRealm"])
	assert.Equal(t, "revoke", parsed["sd_undetectableAction"])
	assert.Equal(t, true, parsed["sd_credentialUpdate"])
	assert.Equal(t, "dG9rZW4=", parsed["sd_encodedToken"])
	assert.Equal(t, "[*]", parsed["sd_jmesPathQuery"])
	_, hasRejectUnauthorized := parsed["sd_rejectUnauthorized"]
	assert.False(t, hasRejectUnauthorized, "expected sd_rejectUnauthorized to be omitted when false due to omitempty")
	assert.Equal(t, "-----BEGIN CERTIFICATE-----", parsed["sd_trustCA"])
	assert.Equal(t, "1", parsed["sd_minimumMonitors"])
	assert.Equal(t, "60", parsed["sd_updateInterval"])
}

// TestDataBigipFastConsulServiceDiscoveryRead covers the Read function's
// happy path: it populates consul_sd_json and sets the resource ID
// directly to the uri (no hashing, unlike the AWS/Azure/GCE datasources).
func TestDataBigipFastConsulServiceDiscoveryRead(t *testing.T) {
	d := consulServiceDiscoveryResourceData(t, map[string]interface{}{
		"uri":  "https://consul.example.com",
		"port": 8500,
	})

	diags := dataBigipFastConsulServiceDiscoveryRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.NotEmpty(t, d.Get("consul_sd_json").(string))
	assert.Equal(t, "https://consul.example.com", d.Id())
}
