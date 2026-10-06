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

// TestDataSourceBigipFastAwsServiceDiscoverySchema exercises the schema
// definition of the bigip_fast_aws_service_discovery data source.
func TestDataSourceBigipFastAwsServiceDiscoverySchema(t *testing.T) {
	ds := dataSourceBigipFastAwsServiceDiscovery()

	if ds.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if ds.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	requiredFields := []string{"tag_key", "tag_value"}
	for _, field := range requiredFields {
		s, ok := ds.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}

	if ds.Schema["aws_sd_json"].Computed != true {
		t.Error("Expected 'aws_sd_json' to be Computed")
	}
	if ds.Schema["type"].Default != "aws" {
		t.Errorf("Expected 'type' to default to 'aws', got %v", ds.Schema["type"].Default)
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

	if !ds.Schema["aws_access_key"].Sensitive {
		t.Error("Expected 'aws_access_key' to be Sensitive")
	}
	if !ds.Schema["aws_secret_access_key"].Sensitive {
		t.Error("Expected 'aws_secret_access_key' to be Sensitive")
	}
}

func awsServiceDiscoveryResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	ds := dataSourceBigipFastAwsServiceDiscovery()
	return schema.TestResourceDataRaw(t, ds.Schema, raw)
}

// TestGetAwsSDConfigDefaults covers getAwsSDConfig with only the required
// fields set, verifying the marshalled JSON reflects schema defaults and
// omits empty optional fields (omitempty).
func TestGetAwsSDConfigDefaults(t *testing.T) {
	d := awsServiceDiscoveryResourceData(t, map[string]interface{}{
		"tag_key":   "Name",
		"tag_value": "my-pool",
	})

	config, err := getAwsSDConfig(d)

	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(config), &parsed))

	assert.Equal(t, "aws", parsed["sd_type"])
	assert.Equal(t, float64(80), parsed["sd_port"])
	assert.Equal(t, "Name", parsed["sd_tag_key"])
	assert.Equal(t, "my-pool", parsed["sd_tag_val"])
	assert.Equal(t, "private", parsed["sd_addressRealm"])
	assert.Equal(t, "remove", parsed["sd_undetectableAction"])
	// credential_update defaults to false, and the struct field has no
	// omitempty tag override for bool, but false + omitempty means it's
	// dropped from the JSON entirely.
	_, hasCredentialUpdate := parsed["sd_credentialUpdate"]
	assert.False(t, hasCredentialUpdate, "expected sd_credentialUpdate to be omitted when false")
	// aws_access_key/aws_secret_access_key/external_id/role_arn/etc are
	// unset and all have omitempty, so should be absent.
	for _, key := range []string{"sd_accessKeyId", "sd_secretAccessKey", "sd_externalId", "sd_roleARN", "sd_aws_region", "sd_minimumMonitors", "sd_updateInterval"} {
		_, ok := parsed[key]
		assert.False(t, ok, "expected %s to be omitted when unset", key)
	}
}

// TestGetAwsSDConfigAllFields covers getAwsSDConfig with every optional
// field populated, verifying each maps to the expected JSON key.
func TestGetAwsSDConfigAllFields(t *testing.T) {
	d := awsServiceDiscoveryResourceData(t, map[string]interface{}{
		"type":                  "aws",
		"port":                  8443,
		"tag_key":               "Name",
		"tag_value":             "my-pool",
		"address_realm":         "public",
		"undetectable_action":   "revoke",
		"credential_update":     true,
		"aws_region":            "us-west-2",
		"aws_access_key":        "AKIAEXAMPLE",
		"aws_secret_access_key": "supersecret",
		"external_id":           "ext-1",
		"role_arn":              "arn:aws:iam::123:role/foo",
		"minimum_monitors":      "1",
		"update_interval":       "60",
	})

	config, err := getAwsSDConfig(d)

	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(config), &parsed))

	assert.Equal(t, "aws", parsed["sd_type"])
	assert.Equal(t, float64(8443), parsed["sd_port"])
	assert.Equal(t, "Name", parsed["sd_tag_key"])
	assert.Equal(t, "my-pool", parsed["sd_tag_val"])
	assert.Equal(t, "public", parsed["sd_addressRealm"])
	assert.Equal(t, "revoke", parsed["sd_undetectableAction"])
	assert.Equal(t, true, parsed["sd_credentialUpdate"])
	assert.Equal(t, "us-west-2", parsed["sd_aws_region"])
	assert.Equal(t, "AKIAEXAMPLE", parsed["sd_accessKeyId"])
	assert.Equal(t, "supersecret", parsed["sd_secretAccessKey"])
	assert.Equal(t, "ext-1", parsed["sd_externalId"])
	assert.Equal(t, "arn:aws:iam::123:role/foo", parsed["sd_roleARN"])
	assert.Equal(t, "1", parsed["sd_minimumMonitors"])
	assert.Equal(t, "60", parsed["sd_updateInterval"])
}

// TestDataBigipFastCAwsServiceDiscoveryRead covers the Read function's
// happy path: it populates aws_sd_json and derives the resource ID from a
// hash of tag_key.
func TestDataBigipFastCAwsServiceDiscoveryRead(t *testing.T) {
	d := awsServiceDiscoveryResourceData(t, map[string]interface{}{
		"tag_key":   "Name",
		"tag_value": "my-pool",
	})

	diags := dataBigipFastCAwsServiceDiscoveryRead(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.NotEmpty(t, d.Get("aws_sd_json").(string))
	assert.Equal(t, hashForState("Name"), d.Id())
}
