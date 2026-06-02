//go:build unit
// +build unit

package bigip

import (
	"testing"
)

// TestDataSourceBigipGtmDatacenterExists proves that a data source for GTM
// datacenters is registered in the provider. Without the fix, this test fails
// because no such data source exists.
func TestDataSourceBigipGtmDatacenterExists(t *testing.T) {
	p := Provider()

	if _, ok := p.DataSourcesMap["bigip_gtm_datacenter"]; !ok {
		t.Fatal("Expected data source 'bigip_gtm_datacenter' to be registered in the provider")
	}
}

// TestDataSourceBigipGtmServerExists proves that a data source for GTM
// servers is registered in the provider. Without the fix, this test fails
// because no such data source exists.
func TestDataSourceBigipGtmServerExists(t *testing.T) {
	p := Provider()

	if _, ok := p.DataSourcesMap["bigip_gtm_server"]; !ok {
		t.Fatal("Expected data source 'bigip_gtm_server' to be registered in the provider")
	}
}