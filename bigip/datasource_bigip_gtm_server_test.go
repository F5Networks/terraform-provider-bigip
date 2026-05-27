//go:build unit
// +build unit

package bigip

import (
	"testing"
)

func TestDataSourceBigipGtmServerSchema(t *testing.T) {
	ds := dataSourceBigipGtmServer()

	if ds.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	// name is the only required field (used to look up the server)
	nameField, ok := ds.Schema["name"]
	if !ok {
		t.Fatal("Expected field 'name' to exist")
	}
	if !nameField.Required {
		t.Error("Expected field 'name' to be required")
	}

	computedFields := []string{"datacenter", "description", "product", "enabled", "monitor",
		"virtual_server_discovery", "link_discovery", "prober_preference", "prober_fallback",
		"prober_pool", "addresses", "virtual_servers"}
	for _, field := range computedFields {
		s, ok := ds.Schema[field]
		if !ok {
			t.Errorf("Expected field '%s' to exist", field)
			continue
		}
		if !s.Computed {
			t.Errorf("Expected field '%s' to be computed", field)
		}
	}
}
