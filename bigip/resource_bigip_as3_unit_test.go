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

func getDiffSuppressFunc() schema.SchemaDiffSuppressFunc {
	r := resourceBigipAs3()
	return r.Schema["as3_json"].DiffSuppressFunc
}

func TestDiffSuppressFunc_IdenticalJSON(t *testing.T) {
	diffSuppress := getDiffSuppressFunc()

	old := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","MyTenant":{"class":"Tenant"}}}`
	new := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","MyTenant":{"class":"Tenant"}}}`

	r := resourceBigipAs3()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ignore_metadata": false,
	})

	if !diffSuppress("as3_json", old, new, d) {
		t.Error("Expected identical JSON to suppress diff")
	}
}

func TestDiffSuppressFunc_IgnoreMetadataFalse_RealDiff(t *testing.T) {
	diffSuppress := getDiffSuppressFunc()

	old := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","MyTenant":{"class":"Tenant","app":{"class":"Application"}}}}`
	new := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","MyTenant":{"class":"Tenant","app":{"class":"Application","pool":{"class":"Pool"}}}}}`

	r := resourceBigipAs3()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ignore_metadata": false,
	})

	if diffSuppress("as3_json", old, new, d) {
		t.Error("Expected real diff to NOT be suppressed when ignore_metadata=false")
	}
}

func TestDiffSuppressFunc_IgnoreMetadataTrue_MetadataOnlyDiff(t *testing.T) {
	diffSuppress := getDiffSuppressFunc()

	// old has extra metadata fields that BIG-IP added
	old := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.50.0","id":"autogen_id","updateMode":"selective","label":"some label","remark":"auto remark","MyTenant":{"class":"Tenant"}},"persist":true}`
	new := `{"class":"AS3","declaration":{"class":"ADC","MyTenant":{"class":"Tenant"}}}`

	r := resourceBigipAs3()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ignore_metadata": true,
	})

	if !diffSuppress("as3_json", old, new, d) {
		t.Error("Expected metadata-only diff to be suppressed when ignore_metadata=true")
	}
}

func TestDiffSuppressFunc_IgnoreMetadataTrue_RealDiff(t *testing.T) {
	diffSuppress := getDiffSuppressFunc()

	// Real user config change beyond metadata
	old := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","updateMode":"selective","MyTenant":{"class":"Tenant","app":{"class":"Application"}}}}`
	new := `{"class":"AS3","declaration":{"class":"ADC","MyTenant":{"class":"Tenant","app":{"class":"Application","pool":{"class":"Pool"}}}}}`

	r := resourceBigipAs3()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ignore_metadata": true,
	})

	if diffSuppress("as3_json", old, new, d) {
		t.Error("Expected real config diff to NOT be suppressed even when ignore_metadata=true")
	}
}

func TestDiffSuppressFunc_IgnoreMetadataTrue_CommonNotUserDefined(t *testing.T) {
	diffSuppress := getDiffSuppressFunc()

	// old has Common added by BIG-IP, new does NOT define Common
	old := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","Common":{"class":"Tenant","Shared":{"class":"Application"}},"MyTenant":{"class":"Tenant"}}}`
	new := `{"class":"AS3","declaration":{"class":"ADC","MyTenant":{"class":"Tenant"}}}`

	r := resourceBigipAs3()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ignore_metadata": true,
	})

	if !diffSuppress("as3_json", old, new, d) {
		t.Error("Expected auto-created Common to be suppressed when user did not define it and ignore_metadata=true")
	}
}

func TestDiffSuppressFunc_IgnoreMetadataTrue_CommonUserDefined(t *testing.T) {
	diffSuppress := getDiffSuppressFunc()

	// Both old and new define Common, but with different content
	old := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","Common":{"class":"Tenant","Shared":{"class":"Application","myPool":{"class":"Pool","members":[{"servicePort":80}]}}},"MyTenant":{"class":"Tenant"}}}`
	new := `{"class":"AS3","declaration":{"class":"ADC","Common":{"class":"Tenant","Shared":{"class":"Application","myPool":{"class":"Pool","members":[{"servicePort":8080}]}}},"MyTenant":{"class":"Tenant"}}}`

	r := resourceBigipAs3()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ignore_metadata": true,
	})

	if diffSuppress("as3_json", old, new, d) {
		t.Error("Expected user-defined Common changes to NOT be suppressed even with ignore_metadata=true")
	}
}

func TestDiffSuppressFunc_IgnoreMetadataFalse_MetadataOnlyDiff(t *testing.T) {
	diffSuppress := getDiffSuppressFunc()

	// Only metadata differs — but ignore_metadata is false, so strict mode should detect it
	old := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.50.0","id":"autogen_id","updateMode":"selective","MyTenant":{"class":"Tenant"}}}`
	new := `{"class":"AS3","declaration":{"class":"ADC","MyTenant":{"class":"Tenant"}}}`

	r := resourceBigipAs3()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ignore_metadata": false,
	})

	if diffSuppress("as3_json", old, new, d) {
		t.Error("Expected metadata-only diff to NOT be suppressed when ignore_metadata=false (strict mode)")
	}
}
