/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Unit tests for GTM Monitor resources - no F5 BIG-IP connection required

func TestResourceBigipGtmMonitorHttpSchema(t *testing.T) {
	resource := resourceBigipGtmMonitorHttp()

	// Verify schema exists
	if resource.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	// Test required fields
	requiredFields := []string{"name"}
	for _, field := range requiredFields {
		s, ok := resource.Schema[field]
		if !ok {
			t.Errorf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}

	// Test optional fields with defaults
	optionalWithDefaults := map[string]interface{}{
		"defaults_from":        "/Common/http",
		"destination":          "*:*",
		"interval":             30,
		"timeout":              120,
		"probe_timeout":        5,
		"ignore_down_response": "disabled",
		"transparent":          "disabled",
		"reverse":              "disabled",
	}

	for field, expectedDefault := range optionalWithDefaults {
		s, ok := resource.Schema[field]
		if !ok {
			t.Errorf("Expected field '%s' to exist in schema", field)
			continue
		}
		if s.Required {
			t.Errorf("Expected field '%s' to be optional", field)
		}
		if s.Default != expectedDefault {
			t.Errorf("Expected field '%s' default to be %v, got %v", field, expectedDefault, s.Default)
		}
	}

	// Verify CRUD functions are defined
	if resource.CreateContext == nil {
		t.Error("Expected CreateContext to be defined")
	}
	if resource.ReadContext == nil {
		t.Error("Expected ReadContext to be defined")
	}
	if resource.UpdateContext == nil {
		t.Error("Expected UpdateContext to be defined")
	}
	if resource.DeleteContext == nil {
		t.Error("Expected DeleteContext to be defined")
	}

	// Verify Importer is defined
	if resource.Importer == nil {
		t.Error("Expected Importer to be defined")
	}
}

func TestResourceBigipGtmMonitorHttpsSchema(t *testing.T) {
	resource := resourceBigipGtmMonitorHttps()

	if resource.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	// Test HTTPS-specific fields
	httpsFields := []string{"cert", "key", "cipherlist", "compatibility"}
	for _, field := range httpsFields {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("Expected HTTPS-specific field '%s' to exist in schema", field)
		}
	}

	// Verify cipherlist default
	if s, ok := resource.Schema["cipherlist"]; ok {
		expectedDefault := "DEFAULT:+SHA:+3DES:+kEDH"
		if s.Default != expectedDefault {
			t.Errorf("Expected cipherlist default to be '%s', got '%v'", expectedDefault, s.Default)
		}
	}

	// Verify compatibility default
	if s, ok := resource.Schema["compatibility"]; ok {
		expectedDefault := "enabled"
		if s.Default != expectedDefault {
			t.Errorf("Expected compatibility default to be '%s', got '%v'", expectedDefault, s.Default)
		}
	}
}

func TestResourceBigipGtmMonitorTcpSchema(t *testing.T) {
	resource := resourceBigipGtmMonitorTcp()

	if resource.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	// Verify send/receive are optional (not required)
	optionalFields := []string{"send", "receive"}
	for _, field := range optionalFields {
		s, ok := resource.Schema[field]
		if !ok {
			t.Errorf("Expected field '%s' to exist in schema", field)
			continue
		}
		if s.Required {
			t.Errorf("Expected field '%s' to be optional for TCP monitors", field)
		}
	}

	// Verify defaults_from default value
	if s, ok := resource.Schema["defaults_from"]; ok {
		expectedDefault := "/Common/tcp"
		if s.Default != expectedDefault {
			t.Errorf("Expected defaults_from default to be '%s', got '%v'", expectedDefault, s.Default)
		}
	}
}

func TestResourceBigipGtmMonitorPostgresqlSchema(t *testing.T) {
	resource := resourceBigipGtmMonitorPostgresql()

	if resource.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	// Test PostgreSQL-specific fields
	postgresqlFields := []string{"database", "username", "password", "count", "debug"}
	for _, field := range postgresqlFields {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("Expected PostgreSQL-specific field '%s' to exist in schema", field)
		}
	}

	// Verify password is marked as sensitive
	if s, ok := resource.Schema["password"]; ok {
		if !s.Sensitive {
			t.Error("Expected password field to be marked as sensitive")
		}
	}

	// Verify debug default
	if s, ok := resource.Schema["debug"]; ok {
		expectedDefault := "no"
		if s.Default != expectedDefault {
			t.Errorf("Expected debug default to be '%s', got '%v'", expectedDefault, s.Default)
		}
	}

	// Verify defaults_from default value
	if s, ok := resource.Schema["defaults_from"]; ok {
		expectedDefault := "/Common/postgresql"
		if s.Default != expectedDefault {
			t.Errorf("Expected defaults_from default to be '%s', got '%v'", expectedDefault, s.Default)
		}
	}
}

func TestResourceBigipGtmMonitorBigipSchema(t *testing.T) {
	resource := resourceBigipGtmMonitorBigip()

	if resource.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	// Verify BIG-IP monitor does NOT have probe_timeout
	if _, ok := resource.Schema["probe_timeout"]; ok {
		t.Error("BIG-IP monitor should not have probe_timeout field")
	}

	// Verify BIG-IP monitor has aggregation_type
	if _, ok := resource.Schema["aggregation_type"]; !ok {
		t.Error("Expected BIG-IP monitor to have aggregation_type field")
	}

	// Verify aggregation_type default
	if s, ok := resource.Schema["aggregation_type"]; ok {
		expectedDefault := "none"
		if s.Default != expectedDefault {
			t.Errorf("Expected aggregation_type default to be '%s', got '%v'", expectedDefault, s.Default)
		}
	}

	// Verify timeout default is 90 (different from other monitors)
	if s, ok := resource.Schema["timeout"]; ok {
		expectedDefault := 90
		if s.Default != expectedDefault {
			t.Errorf("Expected timeout default to be %d, got %v", expectedDefault, s.Default)
		}
	}

	// Verify defaults_from default value
	if s, ok := resource.Schema["defaults_from"]; ok {
		expectedDefault := "/Common/bigip"
		if s.Default != expectedDefault {
			t.Errorf("Expected defaults_from default to be '%s', got '%v'", expectedDefault, s.Default)
		}
	}

	// Verify BIG-IP monitor does NOT have send/receive
	if _, ok := resource.Schema["send"]; ok {
		t.Error("BIG-IP monitor should not have send field")
	}
	if _, ok := resource.Schema["receive"]; ok {
		t.Error("BIG-IP monitor should not have receive field")
	}
}

func TestResourceBigipGtmMonitorSchemaTypes(t *testing.T) {
	// Test that all resources have correct schema types
	resources := map[string]*schema.Resource{
		"http":       resourceBigipGtmMonitorHttp(),
		"https":      resourceBigipGtmMonitorHttps(),
		"tcp":        resourceBigipGtmMonitorTcp(),
		"postgresql": resourceBigipGtmMonitorPostgresql(),
		"bigip":      resourceBigipGtmMonitorBigip(),
	}

	for monitorType, resource := range resources {
		t.Run(monitorType, func(t *testing.T) {
			// Verify string fields are TypeString
			stringFields := []string{"name", "defaults_from", "destination"}
			for _, field := range stringFields {
				if s, ok := resource.Schema[field]; ok {
					if s.Type != schema.TypeString {
						t.Errorf("[%s] Expected field '%s' to be TypeString, got %v", monitorType, field, s.Type)
					}
				}
			}

			// Verify integer fields are TypeInt
			intFields := []string{"interval", "timeout"}
			for _, field := range intFields {
				if s, ok := resource.Schema[field]; ok {
					if s.Type != schema.TypeInt {
						t.Errorf("[%s] Expected field '%s' to be TypeInt, got %v", monitorType, field, s.Type)
					}
				}
			}
		})
	}
}

func TestResourceBigipGtmMonitorNameValidation(t *testing.T) {
	resources := map[string]*schema.Resource{
		"http":       resourceBigipGtmMonitorHttp(),
		"https":      resourceBigipGtmMonitorHttps(),
		"tcp":        resourceBigipGtmMonitorTcp(),
		"postgresql": resourceBigipGtmMonitorPostgresql(),
		"bigip":      resourceBigipGtmMonitorBigip(),
	}

	for monitorType, resource := range resources {
		t.Run(monitorType, func(t *testing.T) {
			nameSchema, ok := resource.Schema["name"]
			if !ok {
				t.Fatalf("[%s] Expected 'name' field to exist", monitorType)
			}

			// Verify name field has validation function
			if nameSchema.ValidateFunc == nil {
				t.Errorf("[%s] Expected 'name' field to have ValidateFunc", monitorType)
			}

			// Verify name is ForceNew (cannot be changed after creation)
			if !nameSchema.ForceNew {
				t.Errorf("[%s] Expected 'name' field to be ForceNew", monitorType)
			}
		})
	}
}

func TestResourceBigipGtmMonitorDescriptions(t *testing.T) {
	resources := map[string]*schema.Resource{
		"http":       resourceBigipGtmMonitorHttp(),
		"https":      resourceBigipGtmMonitorHttps(),
		"tcp":        resourceBigipGtmMonitorTcp(),
		"postgresql": resourceBigipGtmMonitorPostgresql(),
		"bigip":      resourceBigipGtmMonitorBigip(),
	}

	for monitorType, resource := range resources {
		t.Run(monitorType, func(t *testing.T) {
			// Verify all fields have descriptions
			for fieldName, fieldSchema := range resource.Schema {
				if fieldSchema.Description == "" {
					t.Errorf("[%s] Field '%s' missing description", monitorType, fieldName)
				}
			}
		})
	}
}

func TestResourceBigipGtmMonitorHttpDefaultSendString(t *testing.T) {
	resource := resourceBigipGtmMonitorHttp()

	sendSchema, ok := resource.Schema["send"]
	if !ok {
		t.Fatal("Expected 'send' field to exist")
	}

	// HTTP monitor should have default send string
	expectedDefault := "GET /\\r\\n"
	if sendSchema.Default != expectedDefault {
		t.Errorf("Expected send default to be '%s', got '%v'", expectedDefault, sendSchema.Default)
	}
}

func TestNormalizeGtmBigipDestination(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string returns default",
			input:    "",
			expected: "*:*",
		},
		{
			name:     "wildcard without port gets port appended",
			input:    "*",
			expected: "*:*",
		},
		{
			name:     "IP without port gets wildcard port appended",
			input:    "10.1.1.100",
			expected: "10.1.1.100:*",
		},
		{
			name:     "already has port - no change",
			input:    "*:*",
			expected: "*:*",
		},
		{
			name:     "IP with specific port - no change",
			input:    "10.1.1.100:8080",
			expected: "10.1.1.100:8080",
		},
		{
			name:     "wildcard IP with specific port - no change",
			input:    "*:5432",
			expected: "*:5432",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeGtmBigipDestination(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeGtmBigipDestination(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateGtmBigipDestination(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		expectError bool
	}{
		{
			name:        "empty string is valid",
			destination: "",
			expectError: false,
		},
		{
			name:        "all wildcards is valid",
			destination: "*:*",
			expectError: false,
		},
		{
			name:        "wildcard IP with specific port is valid",
			destination: "*:80",
			expectError: false,
		},
		{
			name:        "specific IP with specific port is valid",
			destination: "10.1.1.50:80",
			expectError: false,
		},
		{
			name:        "specific IP with wildcard port is INVALID",
			destination: "10.1.1.50:*",
			expectError: true,
		},
		{
			name:        "another specific IP with wildcard port is INVALID",
			destination: "192.168.1.100:*",
			expectError: true,
		},
		{
			name:        "specific IP with port 443 is valid",
			destination: "10.1.1.50:443",
			expectError: false,
		},
		{
			name:        "wildcard IP with port 5432 is valid",
			destination: "*:5432",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGtmBigipDestination(tt.destination)
			if tt.expectError && err == nil {
				t.Errorf("validateGtmBigipDestination(%q) expected error, got nil", tt.destination)
			}
			if !tt.expectError && err != nil {
				t.Errorf("validateGtmBigipDestination(%q) unexpected error: %v", tt.destination, err)
			}
		})
	}
}

func TestResourceBigipGtmMonitorIntervalTimeoutValidation(t *testing.T) {
	// Verify timeout is always greater than interval in defaults
	resources := map[string]*schema.Resource{
		"http":       resourceBigipGtmMonitorHttp(),
		"https":      resourceBigipGtmMonitorHttps(),
		"tcp":        resourceBigipGtmMonitorTcp(),
		"postgresql": resourceBigipGtmMonitorPostgresql(),
		"bigip":      resourceBigipGtmMonitorBigip(),
	}

	for monitorType, resource := range resources {
		t.Run(monitorType, func(t *testing.T) {
			intervalSchema, okInterval := resource.Schema["interval"]
			timeoutSchema, okTimeout := resource.Schema["timeout"]

			if !okInterval || !okTimeout {
				t.Fatalf("[%s] Expected both interval and timeout fields to exist", monitorType)
			}

			interval, okIntervalInt := intervalSchema.Default.(int)
			timeout, okTimeoutInt := timeoutSchema.Default.(int)

			if !okIntervalInt || !okTimeoutInt {
				t.Fatalf("[%s] Expected interval and timeout defaults to be integers", monitorType)
			}

			if timeout <= interval {
				t.Errorf("[%s] Expected timeout (%d) to be greater than interval (%d)", monitorType, timeout, interval)
			}
		})
	}
}
