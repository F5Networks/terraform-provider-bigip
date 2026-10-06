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

// TestDataSourceBigipFastAzureServiceDiscoverySchema exercises the schema
// definition of the bigip_fast_azure_service_discovery data source.
func TestDataSourceBigipFastAzureServiceDiscoverySchema(t *testing.T) {
	ds := dataSourceBigipFastAzureServiceDiscovery()

	if ds.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if ds.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	requiredFields := []string{"resource_group", "subscription_id"}
	for _, field := range requiredFields {
		s, ok := ds.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}

	if ds.Schema["azure_sd_json"].Computed != true {
		t.Error("Expected 'azure_sd_json' to be Computed")
	}
	if ds.Schema["type"].Default != "azure" {
		t.Errorf("Expected 'type' to default to 'azure', got %v", ds.Schema["type"].Default)
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

func azureServiceDiscoveryResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	ds := dataSourceBigipFastAzureServiceDiscovery()
	return schema.TestResourceDataRaw(t, ds.Schema, raw)
}

// TestGetAzureSDConfigDefaults covers getAzureSDConfig with only the
// required fields set, verifying schema defaults and the two
// unconditionally-set fields (sd_rtype, sd_useManagedIdentity).
func TestGetAzureSDConfigDefaults(t *testing.T) {
	d := azureServiceDiscoveryResourceData(t, map[string]interface{}{
		"resource_group":  "my-rg",
		"subscription_id": "sub-123",
	})

	config, err := getAzureSDConfig(d)

	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(config), &parsed))

	assert.Equal(t, "azure", parsed["sd_type"])
	assert.Equal(t, float64(80), parsed["sd_port"])
	assert.Equal(t, "my-rg", parsed["sd_rg"])
	assert.Equal(t, "sub-123", parsed["sd_sid"])
	assert.Equal(t, "tag", parsed["sd_rtype"])
	assert.Equal(t, true, parsed["sd_useManagedIdentity"])
	assert.Equal(t, "private", parsed["sd_addressRealm"])
	assert.Equal(t, "remove", parsed["sd_undetectableAction"])
	_, hasCredentialUpdate := parsed["sd_credentialUpdate"]
	assert.False(t, hasCredentialUpdate, "expected sd_credentialUpdate to be omitted when false")
	for _, key := range []string{"sd_azure_tag_key", "sd_azure_tag_val", "sd_minimumMonitors", "sd_updateInterval"} {
		_, ok := parsed[key]
		assert.False(t, ok, "expected %s to be omitted when unset", key)
	}
}

// TestGetAzureSDConfigAllFields covers getAzureSDConfig with every optional
// field populated.
func TestGetAzureSDConfigAllFields(t *testing.T) {
	d := azureServiceDiscoveryResourceData(t, map[string]interface{}{
		"type":                "azure",
		"port":                8443,
		"resource_group":      "my-rg",
		"subscription_id":     "sub-123",
		"address_realm":       "public",
		"undetectable_action": "revoke",
		"credential_update":   true,
		"tag_key":             "Name",
		"tag_value":           "my-pool",
		"minimum_monitors":    "1",
		"update_interval":     "60",
	})

	config, err := getAzureSDConfig(d)

	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(config), &parsed))

	assert.Equal(t, float64(8443), parsed["sd_port"])
	assert.Equal(t, "public", parsed["sd_addressRealm"])
	assert.Equal(t, "revoke", parsed["sd_undetectableAction"])
	assert.Equal(t, true, parsed["sd_credentialUpdate"])
	assert.Equal(t, "Name", parsed["sd_azure_tag_key"])
	assert.Equal(t, "my-pool", parsed["sd_azure_tag_val"])
	assert.Equal(t, "1", parsed["sd_minimumMonitors"])
	assert.Equal(t, "60", parsed["sd_updateInterval"])
}

// TestDataBigipFastAzureServiceDiscoveryRead covers the Read function's
// happy path: it populates azure_sd_json and derives the resource ID from
// a hash of the resulting JSON string itself.
func TestDataBigipFastAzureServiceDiscoveryRead(t *testing.T) {
	d := azureServiceDiscoveryResourceData(t, map[string]interface{}{
		"resource_group":  "my-rg",
		"subscription_id": "sub-123",
	})

	diags := dataBigipFastAzureServiceDiscoveryRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	jsonVal := d.Get("azure_sd_json").(string)
	assert.NotEmpty(t, jsonVal)
	assert.Equal(t, hashForState(jsonVal), d.Id())
}
