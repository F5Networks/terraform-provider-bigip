//go:build unit
// +build unit

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

// TestResourceBigipGtmWideipPoolsSchemaExists reproduces the customer-reported
// issue: the bigip_gtm_wideip resource has no "pools" attribute, so there is
// no way to associate pools with a WideIP. Without the fix this test fails
// with: 'Expected "pools" field to exist in schema'.
func TestResourceBigipGtmWideipPoolsSchemaExists(t *testing.T) {
	resource := resourceBigipGtmWideip()

	poolsSchema, ok := resource.Schema["pools"]
	if !ok {
		t.Fatal("Expected 'pools' field to exist in schema. The bigip_gtm_wideip resource does not support associating pools.")
	}

	if poolsSchema.Type != schema.TypeList {
		t.Errorf("Expected 'pools' to be TypeList, got %v", poolsSchema.Type)
	}

	if !poolsSchema.Optional {
		t.Error("Expected 'pools' to be optional")
	}

	// Verify the nested resource schema
	poolElem, ok := poolsSchema.Elem.(*schema.Resource)
	if !ok {
		t.Fatal("Expected 'pools' Elem to be a *schema.Resource")
	}

	// Verify 'name' sub-field
	nameField, ok := poolElem.Schema["name"]
	if !ok {
		t.Fatal("Expected 'pools.name' field to exist")
	}
	if !nameField.Required {
		t.Error("Expected 'pools.name' to be required")
	}
	if nameField.Type != schema.TypeString {
		t.Errorf("Expected 'pools.name' to be TypeString, got %v", nameField.Type)
	}

	// Verify 'order' sub-field
	orderField, ok := poolElem.Schema["order"]
	if !ok {
		t.Fatal("Expected 'pools.order' field to exist")
	}
	if orderField.Type != schema.TypeInt {
		t.Errorf("Expected 'pools.order' to be TypeInt, got %v", orderField.Type)
	}
	if orderField.Default != 0 {
		t.Errorf("Expected 'pools.order' default to be 0, got %v", orderField.Default)
	}

	// Verify 'ratio' sub-field
	ratioField, ok := poolElem.Schema["ratio"]
	if !ok {
		t.Fatal("Expected 'pools.ratio' field to exist")
	}
	if ratioField.Type != schema.TypeInt {
		t.Errorf("Expected 'pools.ratio' to be TypeInt, got %v", ratioField.Type)
	}
	if ratioField.Default != 1 {
		t.Errorf("Expected 'pools.ratio' default to be 1, got %v", ratioField.Default)
	}
}