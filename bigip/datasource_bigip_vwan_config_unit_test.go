/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE on coverage scope: DownloadVwanConfig (the bulk of this file) makes
// live network calls to Azure Blob Storage (https://<account>.blob.core.
// windows.net) and the Azure Resource Manager API (https://management.
// azure.com) via hardcoded URLs baked into the go-bigip/Azure SDK client
// constructors, with no injectable HTTP transport, base URI override, or
// other seam available from outside the package. It cannot be exercised in
// a unit test without either real Azure credentials/resources or modifying
// datasource_bigip_vwan_config.go to add a testing seam (out of scope
// here). These tests therefore cover everything that is reachable without
// network I/O: the schema definition, Read's environment-variable
// validation (fails fast before any Azure call is attempted), and
// CreateToken's local, non-network config/token construction and its
// input-validation error paths.

// vwanConfigResourceData builds a *schema.ResourceData for the
// bigip_vwan_config data source using schema.TestResourceDataRaw.
func vwanConfigResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := dataSourceBigipVwanconfig()
	return schema.TestResourceDataRaw(t, r.Schema, raw)
}

// clearVwanAzureEnv sets all Azure-related environment variables that
// dataSourceBigipVwanconfigRead checks to empty strings (the same value
// os.Getenv returns for a truly unset variable, and the only condition the
// code under test distinguishes), and restores their original values (if
// any) after the test via t.Cleanup. This isolates the test from whatever
// the ambient environment happens to have set (e.g. a developer's shell or
// a CI secret) so the "missing env vars" branch is deterministic.
func clearVwanAzureEnv(t *testing.T) {
	t.Helper()
	vars := []string{
		"AZURE_SUBSCRIPTION_ID",
		"AZURE_CLIENT_ID",
		"AZURE_CLIENT_SECRET",
		"AZURE_TENANT_ID",
		"STORAGE_ACCOUNT_NAME",
		"STORAGE_ACCOUNT_KEY",
	}
	for _, v := range vars {
		t.Setenv(v, "")
	}
}

// TestDataSourceBigipVwanconfigSchema exercises the schema definition of
// the bigip_vwan_config data source without needing any BIG-IP or Azure
// connection.
func TestDataSourceBigipVwanconfigSchema(t *testing.T) {
	r := dataSourceBigipVwanconfig()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}
	if r.ReadContext == nil {
		t.Fatal("Expected ReadContext to be defined")
	}

	requiredFields := []string{"azure_vwan_resourcegroup", "azure_vwan_name", "azure_vwan_vpnsite"}
	for _, field := range requiredFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
		if s.Type != schema.TypeString {
			t.Errorf("Expected field '%s' to be TypeString, got %v", field, s.Type)
		}
	}

	computedStringFields := []string{"bigip_gw_ip", "hub_address_space"}
	for _, field := range computedStringFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Computed {
			t.Errorf("Expected field '%s' to be Computed", field)
		}
		if s.Type != schema.TypeString {
			t.Errorf("Expected field '%s' to be TypeString, got %v", field, s.Type)
		}
	}

	presharedKeySchema, ok := r.Schema["preshared_key"]
	if !ok {
		t.Fatal("Expected field 'preshared_key' to exist in schema")
	}
	if !presharedKeySchema.Computed {
		t.Error("Expected field 'preshared_key' to be Computed")
	}
	if !presharedKeySchema.Sensitive {
		t.Error("Expected field 'preshared_key' to be Sensitive")
	}

	listFields := []string{"hub_connected_subnets", "vwan_gw_address"}
	for _, field := range listFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Computed {
			t.Errorf("Expected field '%s' to be Computed", field)
		}
		if s.Type != schema.TypeList {
			t.Errorf("Expected field '%s' to be TypeList, got %v", field, s.Type)
		}
	}
}

// TestDataSourceBigipVwanconfigReadMissingEnv covers the branch where none
// of the required Azure/storage environment variables are set: Read must
// fail fast with the azureEnvErr message before attempting any network
// call, and must leave the resource ID cleared.
func TestDataSourceBigipVwanconfigReadMissingEnv(t *testing.T) {
	clearVwanAzureEnv(t)

	d := vwanConfigResourceData(t, map[string]interface{}{
		"azure_vwan_resourcegroup": "test-rg",
		"azure_vwan_name":          "test-vwan",
		"azure_vwan_vpnsite":       "test-site",
	})
	d.SetId("preexisting-id")

	diags := dataSourceBigipVwanconfigRead(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error when required Azure env vars are unset")
	assert.Contains(t, diags[0].Summary, "Azure Environment is not set")
	assert.Empty(t, d.Id(), "resource ID should be cleared when required env vars are missing")
}

// TestDataSourceBigipVwanconfigReadMissingEnvPartial covers the same branch
// but with some (not all) of the required environment variables set,
// confirming the check requires every variable rather than just one.
func TestDataSourceBigipVwanconfigReadMissingEnvPartial(t *testing.T) {
	clearVwanAzureEnv(t)
	t.Setenv("AZURE_SUBSCRIPTION_ID", "sub-id")
	t.Setenv("AZURE_CLIENT_ID", "client-id")
	// AZURE_CLIENT_SECRET, AZURE_TENANT_ID, STORAGE_ACCOUNT_NAME, and
	// STORAGE_ACCOUNT_KEY are intentionally left unset.

	d := vwanConfigResourceData(t, map[string]interface{}{
		"azure_vwan_resourcegroup": "test-rg",
		"azure_vwan_name":          "test-vwan",
		"azure_vwan_vpnsite":       "test-site",
	})

	diags := dataSourceBigipVwanconfigRead(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error when any required Azure env var is unset")
	assert.Contains(t, diags[0].Summary, "Azure Environment is not set")
}

// TestDataSourceBigipVwanconfigReadDownloadError covers the branch past the
// env-var check: with all six variables set (so Read proceeds to call
// DownloadVwanConfig), an invalid (non-base64) STORAGE_ACCOUNT_KEY makes
// azblob.NewSharedKeyCredential fail immediately on a local base64-decode
// error -- before DownloadVwanConfig issues any network request -- so this
// exercises Read's "err != nil" handling around the DownloadVwanConfig call
// without needing real Azure credentials or network access.
func TestDataSourceBigipVwanconfigReadDownloadError(t *testing.T) {
	clearVwanAzureEnv(t)
	t.Setenv("AZURE_SUBSCRIPTION_ID", "sub-id")
	t.Setenv("AZURE_CLIENT_ID", "client-id")
	t.Setenv("AZURE_CLIENT_SECRET", "client-secret")
	t.Setenv("AZURE_TENANT_ID", "tenant-id")
	t.Setenv("STORAGE_ACCOUNT_NAME", "teststorageaccount")
	t.Setenv("STORAGE_ACCOUNT_KEY", "not-valid-base64!!")

	d := vwanConfigResourceData(t, map[string]interface{}{
		"azure_vwan_resourcegroup": "test-rg",
		"azure_vwan_name":          "test-vwan",
		"azure_vwan_vpnsite":       "test-site",
	})
	d.SetId("preexisting-id")

	diags := dataSourceBigipVwanconfigRead(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error when STORAGE_ACCOUNT_KEY is not valid base64")
	assert.Empty(t, d.Id(), "resource ID should be cleared on a download error")
}

// TestCreateTokenSuccess covers CreateToken's happy path: valid
// tenant/client/secret values let it build an OAuthConfig and a
// ServicePrincipalToken purely locally (no network call is made by these
// adal constructors), returning a nil error and a non-nil token.
func TestCreateTokenSuccess(t *testing.T) {
	token, err := CreateToken("test-tenant-id", "test-client-id", "test-client-secret")

	require.NoError(t, err)
	assert.NotNil(t, token)
}

// TestCreateTokenEmptyClientID covers the validation-error branch:
// NewServicePrincipalToken rejects an empty clientID before any network
// call would occur.
func TestCreateTokenEmptyClientID(t *testing.T) {
	_, err := CreateToken("test-tenant-id", "", "test-client-secret")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "clientID")
}

// TestCreateTokenEmptySecret covers the validation-error branch where the
// client secret is empty.
func TestCreateTokenEmptySecret(t *testing.T) {
	_, err := CreateToken("test-tenant-id", "test-client-id", "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "secret")
}

// TestCreateTokenInvalidTenantID covers the branch where NewOAuthConfig
// itself fails to construct (an invalid tenantID that url.Parse rejects),
// which causes CreateToken to return early with that error.
func TestCreateTokenInvalidTenantID(t *testing.T) {
	_, err := CreateToken("\x7f", "test-client-id", "test-client-secret")

	require.Error(t, err)
}

// TestDownloadVwanConfigInvalidAccountKey covers DownloadVwanConfig's own
// earliest error branch directly: an invalid (non-base64) accountKey makes
// azblob.NewSharedKeyCredential fail on a local base64-decode error, before
// any network call (container creation, blob upload/download, or the Azure
// Resource Manager vpnconfigClient.Download call) is attempted.
func TestDownloadVwanConfigInvalidAccountKey(t *testing.T) {
	_, err := DownloadVwanConfig(azureConfig{
		subscriptionID:    "sub-id",
		clientID:          "client-id",
		clientPassword:    "client-secret",
		tenantID:          "tenant-id",
		resourceGroupName: "test-rg",
		virtualWANName:    "test-vwan",
		siteName:          "test-site",
		accountName:       "teststorageaccount",
		accountKey:        "not-valid-base64!!",
	})

	require.Error(t, err)
}
