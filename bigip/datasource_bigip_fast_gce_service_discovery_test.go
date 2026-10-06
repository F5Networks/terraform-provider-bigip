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

// TestDataSourceBigipFastGceServiceDiscoverySchema exercises the schema
// definition of the bigip_fast_gce_service_discovery data source.
func TestDataSourceBigipFastGceServiceDiscoverySchema(t *testing.T) {
	ds := dataSourceBigipFastGceServiceDiscovery()

	if ds.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if ds.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	requiredFields := []string{"tag_key", "tag_value", "region"}
	for _, field := range requiredFields {
		s, ok := ds.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}

	if ds.Schema["gce_sd_json"].Computed != true {
		t.Error("Expected 'gce_sd_json' to be Computed")
	}
	if ds.Schema["type"].Default != "gce" {
		t.Errorf("Expected 'type' to default to 'gce', got %v", ds.Schema["type"].Default)
	}
	if ds.Schema["port"].Default != 80 {
		t.Errorf("Expected 'port' to default to 80, got %v", ds.Schema["port"].Default)
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
}

func gceServiceDiscoveryResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	ds := dataSourceBigipFastGceServiceDiscovery()
	return schema.TestResourceDataRaw(t, ds.Schema, raw)
}

// TestGetGceConfigDefaults covers getGceConfig with only the required
// fields set, verifying schema defaults.
func TestGetGceConfigDefaults(t *testing.T) {
	d := gceServiceDiscoveryResourceData(t, map[string]interface{}{
		"tag_key":   "Name",
		"tag_value": "my-pool",
		"region":    "us-central1",
	})

	config, err := getGceConfig(d)

	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(config), &parsed))

	assert.Equal(t, "gce", parsed["sd_type"])
	assert.Equal(t, float64(80), parsed["sd_port"])
	assert.Equal(t, "Name", parsed["sd_tag_key"])
	assert.Equal(t, "my-pool", parsed["sd_tag_val"])
	assert.Equal(t, "us-central1", parsed["sd_region"])
	assert.Equal(t, "private", parsed["sd_addressRealm"])
	assert.Equal(t, "remove", parsed["sd_undetectableAction"])
	_, hasCredentialUpdate := parsed["sd_credentialUpdate"]
	assert.False(t, hasCredentialUpdate, "expected sd_credentialUpdate to be omitted when false")
	for _, key := range []string{"sd_encodedCredentials", "sd_projectId", "sd_minimumMonitors", "sd_updateInterval"} {
		_, ok := parsed[key]
		assert.False(t, ok, "expected %s to be omitted when unset", key)
	}
}

// TestGetGceConfigAllFields covers getGceConfig with every optional field
// populated.
func TestGetGceConfigAllFields(t *testing.T) {
	d := gceServiceDiscoveryResourceData(t, map[string]interface{}{
		"type":                "gce",
		"port":                8443,
		"tag_key":             "Name",
		"tag_value":           "my-pool",
		"region":              "us-central1",
		"address_realm":       "public",
		"undetectable_action": "revoke",
		"credential_update":   true,
		"encoded_credentials": "eyJhbGciOiJSUzI1NiJ9",
		"project_id":          "my-project",
		"minimum_monitors":    "1",
		"update_interval":     "60",
	})

	config, err := getGceConfig(d)

	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(config), &parsed))

	assert.Equal(t, float64(8443), parsed["sd_port"])
	assert.Equal(t, "public", parsed["sd_addressRealm"])
	assert.Equal(t, "revoke", parsed["sd_undetectableAction"])
	assert.Equal(t, true, parsed["sd_credentialUpdate"])
	assert.Equal(t, "eyJhbGciOiJSUzI1NiJ9", parsed["sd_encodedCredentials"])
	assert.Equal(t, "my-project", parsed["sd_projectId"])
	assert.Equal(t, "1", parsed["sd_minimumMonitors"])
	assert.Equal(t, "60", parsed["sd_updateInterval"])
}

// TestDataBigipFastGceServiceDiscoveryRead covers the Read function's
// happy path: it populates gce_sd_json and derives the resource ID from a
// hash of the resulting JSON string itself.
func TestDataBigipFastGceServiceDiscoveryRead(t *testing.T) {
	d := gceServiceDiscoveryResourceData(t, map[string]interface{}{
		"tag_key":   "Name",
		"tag_value": "my-pool",
		"region":    "us-central1",
	})

	diags := dataBigipFastGceServiceDiscoveryRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	jsonVal := d.Get("gce_sd_json").(string)
	assert.NotEmpty(t, jsonVal)
	assert.Equal(t, hashForState(jsonVal), d.Id())
}
