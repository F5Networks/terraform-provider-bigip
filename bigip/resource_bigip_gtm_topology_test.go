//go:build unit
// +build unit

package bigip

import (
	"encoding/json"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceBigipGtmTopologyRecordExists proves that the bigip_gtm_topology_record
// resource is registered in the provider.
func TestResourceBigipGtmTopologyRecordExists(t *testing.T) {
	p := Provider()

	if _, ok := p.ResourcesMap["bigip_gtm_topology_record"]; !ok {
		t.Fatal("Expected resource 'bigip_gtm_topology_record' to be registered in the provider")
	}
}

// TestResourceBigipGtmTopologyRegionExists proves that the bigip_gtm_topology_region
// resource is registered in the provider.
func TestResourceBigipGtmTopologyRegionExists(t *testing.T) {
	p := Provider()

	if _, ok := p.ResourcesMap["bigip_gtm_topology_region"]; !ok {
		t.Fatal("Expected resource 'bigip_gtm_topology_region' to be registered in the provider")
	}
}

func TestResourceBigipGtmTopologyRecordSchema(t *testing.T) {
	resource := resourceBigipGtmTopologyRecord()

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

	// ldns is required, TypeList, MaxItems 1
	ldnsField, ok := resource.Schema["ldns"]
	if !ok {
		t.Fatal("Expected 'ldns' field to exist")
	}
	if !ldnsField.Required {
		t.Error("Expected 'ldns' to be required")
	}
	if ldnsField.Type != schema.TypeList {
		t.Errorf("Expected 'ldns' to be TypeList, got %v", ldnsField.Type)
	}
	if ldnsField.MaxItems != 1 {
		t.Errorf("Expected 'ldns' MaxItems to be 1, got %d", ldnsField.MaxItems)
	}
	ldnsElem, ok := ldnsField.Elem.(*schema.Resource)
	if !ok {
		t.Fatal("Expected 'ldns' Elem to be a *schema.Resource")
	}
	if _, ok := ldnsElem.Schema["match_type"]; !ok {
		t.Error("Expected 'ldns.match_type' field to exist")
	}
	if _, ok := ldnsElem.Schema["match_value"]; !ok {
		t.Error("Expected 'ldns.match_value' field to exist")
	}
	if _, ok := ldnsElem.Schema["match_negate"]; !ok {
		t.Error("Expected 'ldns.match_negate' field to exist")
	}

	// server is required, TypeList, MaxItems 1
	serverField, ok := resource.Schema["server"]
	if !ok {
		t.Fatal("Expected 'server' field to exist")
	}
	if !serverField.Required {
		t.Error("Expected 'server' to be required")
	}
	if serverField.Type != schema.TypeList {
		t.Errorf("Expected 'server' to be TypeList, got %v", serverField.Type)
	}
	if serverField.MaxItems != 1 {
		t.Errorf("Expected 'server' MaxItems to be 1, got %d", serverField.MaxItems)
	}
	serverElem, ok := serverField.Elem.(*schema.Resource)
	if !ok {
		t.Fatal("Expected 'server' Elem to be a *schema.Resource")
	}
	if _, ok := serverElem.Schema["match_type"]; !ok {
		t.Error("Expected 'server.match_type' field to exist")
	}
	if _, ok := serverElem.Schema["match_value"]; !ok {
		t.Error("Expected 'server.match_value' field to exist")
	}
	if _, ok := serverElem.Schema["match_negate"]; !ok {
		t.Error("Expected 'server.match_negate' field to exist")
	}

	// description is optional
	descField, ok := resource.Schema["description"]
	if !ok {
		t.Fatal("Expected 'description' field to exist")
	}
	if descField.Required {
		t.Error("Expected 'description' to be optional, not required")
	}
	if descField.Type != schema.TypeString {
		t.Errorf("Expected 'description' to be TypeString, got %v", descField.Type)
	}

	// order is optional with default 0
	orderField, ok := resource.Schema["order"]
	if !ok {
		t.Fatal("Expected 'order' field to exist")
	}
	if orderField.Type != schema.TypeInt {
		t.Errorf("Expected 'order' to be TypeInt, got %v", orderField.Type)
	}
	if orderField.Default != 0 {
		t.Errorf("Expected 'order' default to be 0, got %v", orderField.Default)
	}

	// score is optional with default 1
	scoreField, ok := resource.Schema["score"]
	if !ok {
		t.Fatal("Expected 'score' field to exist")
	}
	if scoreField.Type != schema.TypeInt {
		t.Errorf("Expected 'score' to be TypeInt, got %v", scoreField.Type)
	}
	if scoreField.Default != 1 {
		t.Errorf("Expected 'score' default to be 1, got %v", scoreField.Default)
	}
}

func TestResourceBigipGtmTopologyRegionSchema(t *testing.T) {
	resource := resourceBigipGtmTopologyRegion()

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

	// name is required and ForceNew
	nameField, ok := resource.Schema["name"]
	if !ok {
		t.Fatal("Expected 'name' field to exist")
	}
	if !nameField.Required {
		t.Error("Expected 'name' to be required")
	}
	if !nameField.ForceNew {
		t.Error("Expected 'name' to be ForceNew")
	}

	// partition is optional with default Common
	partField, ok := resource.Schema["partition"]
	if !ok {
		t.Fatal("Expected 'partition' field to exist")
	}
	if partField.Default != "Common" {
		t.Errorf("Expected 'partition' default to be 'Common', got %v", partField.Default)
	}

	// members is optional TypeSet with nested name field
	membersField, ok := resource.Schema["members"]
	if !ok {
		t.Fatal("Expected 'members' field to exist")
	}
	if membersField.Type != schema.TypeSet {
		t.Errorf("Expected 'members' to be TypeSet, got %v", membersField.Type)
	}
	memberElem, ok := membersField.Elem.(*schema.Resource)
	if !ok {
		t.Fatal("Expected 'members' Elem to be a *schema.Resource")
	}
	memberName, ok := memberElem.Schema["name"]
	if !ok {
		t.Fatal("Expected 'members.name' field to exist")
	}
	if !memberName.Required {
		t.Error("Expected 'members.name' to be required")
	}
	if memberName.Type != schema.TypeString {
		t.Errorf("Expected 'members.name' to be TypeString, got %v", memberName.Type)
	}
}

func TestGTMTopologyRecordMarshalJSON(t *testing.T) {
	record := bigip.GTMTopologyRecord{
		Name:  "ldns: region /Common/east-coast server: datacenter /Common/dc1",
		Order: 1,
		Score: 100,
	}

	data, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("Failed to marshal GTMTopologyRecord: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["name"] != "ldns: region /Common/east-coast server: datacenter /Common/dc1" {
		t.Errorf("Expected name to match, got '%v'", result["name"])
	}
	if int(result["order"].(float64)) != 1 {
		t.Errorf("Expected order 1, got %v", result["order"])
	}
	if int(result["score"].(float64)) != 100 {
		t.Errorf("Expected score 100, got %v", result["score"])
	}
}

func TestGTMTopologyRecordUnmarshalJSON(t *testing.T) {
	jsonData := `{
		"name": "ldns: region /Common/east-coast server: datacenter /Common/dc1",
		"order": 2,
		"score": 50
	}`

	var record bigip.GTMTopologyRecord
	if err := json.Unmarshal([]byte(jsonData), &record); err != nil {
		t.Fatalf("Failed to unmarshal GTMTopologyRecord: %v", err)
	}

	if record.Name != "ldns: region /Common/east-coast server: datacenter /Common/dc1" {
		t.Errorf("Expected name match, got '%s'", record.Name)
	}
	if record.Order != 2 {
		t.Errorf("Expected order 2, got %d", record.Order)
	}
	if record.Score != 50 {
		t.Errorf("Expected score 50, got %d", record.Score)
	}
}

func TestParseTopologyName(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		ldnsType     string
		ldnsValue    string
		ldnsNegate   bool
		serverType   string
		serverValue  string
		serverNegate bool
		expectErr    bool
	}{
		{
			name:        "region to datacenter",
			input:       "ldns: region /Common/east-coast server: datacenter /Common/dc1",
			ldnsType:    "region",
			ldnsValue:   "/Common/east-coast",
			serverType:  "datacenter",
			serverValue: "/Common/dc1",
		},
		{
			name:        "subnet to pool",
			input:       "ldns: subnet 10.0.0.0/8 server: pool /Common/internal-pool",
			ldnsType:    "subnet",
			ldnsValue:   "10.0.0.0/8",
			serverType:  "pool",
			serverValue: "/Common/internal-pool",
		},
		{
			name:        "country to datacenter",
			input:       "ldns: country US server: datacenter /Common/us-dc",
			ldnsType:    "country",
			ldnsValue:   "US",
			serverType:  "datacenter",
			serverValue: "/Common/us-dc",
		},
		{
			name:        "state to datacenter",
			input:       "ldns: state US/California server: datacenter /Common/west-dc",
			ldnsType:    "state",
			ldnsValue:   "US/California",
			serverType:  "datacenter",
			serverValue: "/Common/west-dc",
		},
		{
			name:        "negated ldns",
			input:       "ldns: not region /Common/east-coast server: datacenter /Common/dc1",
			ldnsType:    "region",
			ldnsValue:   "/Common/east-coast",
			ldnsNegate:  true,
			serverType:  "datacenter",
			serverValue: "/Common/dc1",
		},
		{
			name:         "negated server",
			input:        "ldns: region /Common/east-coast server: not datacenter /Common/dc1",
			ldnsType:     "region",
			ldnsValue:    "/Common/east-coast",
			serverType:   "datacenter",
			serverValue:  "/Common/dc1",
			serverNegate: true,
		},
		{
			name:      "invalid format",
			input:     "some invalid string",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lt, lv, ln, st, sv, sn, err := parseTopologyName(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if lt != tt.ldnsType {
				t.Errorf("ldns type: expected '%s', got '%s'", tt.ldnsType, lt)
			}
			if lv != tt.ldnsValue {
				t.Errorf("ldns value: expected '%s', got '%s'", tt.ldnsValue, lv)
			}
			if ln != tt.ldnsNegate {
				t.Errorf("ldns negate: expected %v, got %v", tt.ldnsNegate, ln)
			}
			if st != tt.serverType {
				t.Errorf("server type: expected '%s', got '%s'", tt.serverType, st)
			}
			if sv != tt.serverValue {
				t.Errorf("server value: expected '%s', got '%s'", tt.serverValue, sv)
			}
			if sn != tt.serverNegate {
				t.Errorf("server negate: expected %v, got %v", tt.serverNegate, sn)
			}
		})
	}
}

func TestBuildEndpointString(t *testing.T) {
	tests := []struct {
		name     string
		endpoint map[string]interface{}
		expected string
	}{
		{
			name:     "simple region",
			endpoint: map[string]interface{}{"match_type": "region", "match_value": "/Common/east-coast", "match_negate": false},
			expected: "region /Common/east-coast",
		},
		{
			name:     "negated region",
			endpoint: map[string]interface{}{"match_type": "region", "match_value": "/Common/east-coast", "match_negate": true},
			expected: "not region /Common/east-coast",
		},
		{
			name:     "subnet",
			endpoint: map[string]interface{}{"match_type": "subnet", "match_value": "10.0.0.0/8", "match_negate": false},
			expected: "subnet 10.0.0.0/8",
		},
		{
			name:     "country",
			endpoint: map[string]interface{}{"match_type": "country", "match_value": "US", "match_negate": false},
			expected: "country US",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildEndpointString(tt.endpoint)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGTMRegionMarshalJSON(t *testing.T) {
	region := bigip.GTMRegion{
		Name:      "east-coast",
		Partition: "Common",
		Members: []bigip.GTMRegionMember{
			{Name: "subnet 10.0.0.0/8"},
			{Name: "state US/New-York"},
		},
	}

	data, err := json.Marshal(region)
	if err != nil {
		t.Fatalf("Failed to marshal GTMRegion: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["name"] != "east-coast" {
		t.Errorf("Expected name 'east-coast', got '%v'", result["name"])
	}

	members, ok := result["regionMembers"]
	if !ok {
		t.Fatal("Expected 'regionMembers' key in marshalled JSON")
	}

	membersArr := members.([]interface{})
	if len(membersArr) != 2 {
		t.Fatalf("Expected 2 members, got %d", len(membersArr))
	}

	first := membersArr[0].(map[string]interface{})
	if first["name"] != "subnet 10.0.0.0/8" {
		t.Errorf("Expected first member 'subnet 10.0.0.0/8', got '%v'", first["name"])
	}
}

func TestGTMRegionUnmarshalJSON(t *testing.T) {
	jsonData := `{
		"name": "east-coast",
		"partition": "Common",
		"fullPath": "/Common/east-coast",
		"regionMembers": [
			{"name": "subnet 10.0.0.0/8"},
			{"name": "country US"}
		]
	}`

	var region bigip.GTMRegion
	if err := json.Unmarshal([]byte(jsonData), &region); err != nil {
		t.Fatalf("Failed to unmarshal GTMRegion: %v", err)
	}

	if region.Name != "east-coast" {
		t.Errorf("Expected name 'east-coast', got '%s'", region.Name)
	}
	if region.Partition != "Common" {
		t.Errorf("Expected partition 'Common', got '%s'", region.Partition)
	}
	if len(region.Members) != 2 {
		t.Fatalf("Expected 2 members, got %d", len(region.Members))
	}
	if region.Members[0].Name != "subnet 10.0.0.0/8" {
		t.Errorf("Expected first member 'subnet 10.0.0.0/8', got '%s'", region.Members[0].Name)
	}
	if region.Members[1].Name != "country US" {
		t.Errorf("Expected second member 'country US', got '%s'", region.Members[1].Name)
	}
}

func TestGTMRegionEmptyMembersOmittedFromJSON(t *testing.T) {
	region := bigip.GTMRegion{
		Name:      "empty-region",
		Partition: "Common",
	}

	data, err := json.Marshal(region)
	if err != nil {
		t.Fatalf("Failed to marshal GTMRegion without members: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if _, ok := result["regionMembers"]; ok {
		t.Error("Expected 'regionMembers' to be omitted from JSON when empty")
	}
}
