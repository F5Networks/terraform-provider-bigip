//go:build unit
// +build unit

package bigip

import (
	"testing"
)

func TestDataSourceBigipGtmDatacenterSchema(t *testing.T) {
	ds := dataSourceBigipGtmDatacenter()

	if ds.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	requiredFields := []string{"name", "partition"}
	for _, field := range requiredFields {
		s, ok := ds.Schema[field]
		if !ok {
			t.Errorf("Expected field '%s' to exist", field)
			continue
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}

	computedFields := []string{"description", "contact", "enabled", "location", "prober_fallback", "prober_preference"}
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
