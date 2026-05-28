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

	// description is required and ForceNew
	descField, ok := resource.Schema["description"]
	if !ok {
		t.Fatal("Expected 'description' field to exist")
	}
	if !descField.Required {
		t.Error("Expected 'description' to be required")
	}
	if !descField.ForceNew {
		t.Error("Expected 'description' to be ForceNew")
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
		Description: "ldns: region /Common/east-coast server: datacenter /Common/dc1",
		Order:       1,
		Score:       100,
	}

	data, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("Failed to marshal GTMTopologyRecord: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["description"] != "ldns: region /Common/east-coast server: datacenter /Common/dc1" {
		t.Errorf("Expected description to match, got '%v'", result["description"])
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
		"description": "ldns: region /Common/east-coast server: datacenter /Common/dc1",
		"order": 2,
		"score": 50
	}`

	var record bigip.GTMTopologyRecord
	if err := json.Unmarshal([]byte(jsonData), &record); err != nil {
		t.Fatalf("Failed to unmarshal GTMTopologyRecord: %v", err)
	}

	if record.Description != "ldns: region /Common/east-coast server: datacenter /Common/dc1" {
		t.Errorf("Expected description match, got '%s'", record.Description)
	}
	if record.Order != 2 {
		t.Errorf("Expected order 2, got %d", record.Order)
	}
	if record.Score != 50 {
		t.Errorf("Expected score 50, got %d", record.Score)
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
