/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResourceBigipSysBigiplicenseSchema exercises the schema definition of
// the bigip_sys_bigiplicense resource without needing any BIG-IP connection.
func TestResourceBigipSysBigiplicenseSchema(t *testing.T) {
	r := resourceBigipSysBigiplicense()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	requiredFields := []string{"command", "registration_key"}
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
		if s.Description == "" {
			t.Errorf("Expected field '%s' to have a description", field)
		}
	}

	if r.CreateContext == nil {
		t.Error("Expected CreateContext to be defined")
	}
	if r.ReadContext == nil {
		t.Error("Expected ReadContext to be defined")
	}
	if r.UpdateContext == nil {
		t.Error("Expected UpdateContext to be defined")
	}
	if r.DeleteContext == nil {
		t.Error("Expected DeleteContext to be defined")
	}
	if r.Importer == nil {
		t.Error("Expected Importer to be defined")
	}
}

// testBigiplicenseClient builds a *bigip.BigIP pointed at the given mock
// server URL (typically server.URL from the shared setup()/teardown() mux),
// bypassing the provider's Client()/NewSession negotiation (and its
// SelfIP-based ValidateConnection call) since these unit tests invoke the
// resource CRUD functions directly rather than through resource.Test. This
// also avoids Create's unconditional 300-second time.Sleep, which would
// otherwise make any test exercising Create/Update via Terraform's test
// harness (which always calls Create) far too slow for `make test`'s 30s
// timeout.
func testBigiplicenseClient(t *testing.T, url string) *bigip.BigIP {
	t.Helper()
	return bigip.NewSession(&bigip.Config{
		Address:  url,
		Username: "admin",
		Password: "admin",
		ConfigOptions: &bigip.ConfigOptions{
			APICallRetries: 1,
		},
	})
}

func testBigiplicenseResourceData(t *testing.T, command, registrationKey string) *schema.ResourceData {
	t.Helper()
	r := resourceBigipSysBigiplicense()
	return schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"command":          command,
		"registration_key": registrationKey,
	})
}

// withZeroLicenseApplyDelay temporarily overrides bigipLicenseApplyDelay to 0
// for the duration of the test so that resourceBigipSysBigiplicenseCreate
// (which unconditionally sleeps after issuing the license request) can be
// exercised without blocking the test run.
func withZeroLicenseApplyDelay(t *testing.T) {
	t.Helper()
	original := bigipLicenseApplyDelay
	bigipLicenseApplyDelay = 0
	t.Cleanup(func() { bigipLicenseApplyDelay = original })
}

// TestResourceBigipSysBigiplicenseCreateSuccess covers the Create happy path:
// CreateBigiplicense() succeeds, the resource ID is set to the registration
// key, and the follow-up Read succeeds.
func TestResourceBigipSysBigiplicenseCreateSuccess(t *testing.T) {
	withZeroLicenseApplyDelay(t)

	setup()
	defer teardown()
	var sawPOST bool
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			sawPOST = true
			_, _ = fmt.Fprint(w, `{"registrationKey":"ABCDE-FGHIJ-KLMNO-PQRST","command":"install"}`)
		case "GET":
			_, _ = fmt.Fprint(w, `{"registrationKey":"ABCDE-FGHIJ-KLMNO-PQRST","command":"install"}`)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	client := testBigiplicenseClient(t, server.URL)
	d := testBigiplicenseResourceData(t, "install", "ABCDE-FGHIJ-KLMNO-PQRST")

	diags := resourceBigipSysBigiplicenseCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPOST, "expected CreateBigiplicense to issue a POST request")
	assert.Equal(t, "ABCDE-FGHIJ-KLMNO-PQRST", d.Id(), "resource ID should be set to the registration key")
}

// TestResourceBigipSysBigiplicenseCreateError covers the branch where
// CreateBigiplicense() returns an error. Create still calls time.Sleep after
// the API call regardless of error (same code path as production), but the
// delay is zeroed via withZeroLicenseApplyDelay so the sleep is a no-op here.
func TestResourceBigipSysBigiplicenseCreateError(t *testing.T) {
	withZeroLicenseApplyDelay(t)

	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			http.Error(w, `{"code":400,"message":"invalid registration key"}`, http.StatusBadRequest)
			return
		}
		t.Fatalf("unexpected method %s", r.Method)
	})

	client := testBigiplicenseClient(t, server.URL)
	d := testBigiplicenseResourceData(t, "install", "ABCDE-FGHIJ-KLMNO-PQRST")

	diags := resourceBigipSysBigiplicenseCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create when CreateBigiplicense fails")
	assert.Empty(t, d.Id(), "resource ID should remain unset when Create fails")
}

// TestResourceBigipSysBigiplicenseReadSuccess covers the Read happy path:
// Bigiplicenses() succeeds and returns a non-nil license, so the resource ID
// is left untouched.
func TestResourceBigipSysBigiplicenseReadSuccess(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		_, _ = fmt.Fprint(w, `{"registrationKey":"ABCDE-FGHIJ-KLMNO-PQRST","command":"install"}`)
	})

	client := testBigiplicenseClient(t, server.URL)
	d := testBigiplicenseResourceData(t, "install", "ABCDE-FGHIJ-KLMNO-PQRST")
	d.SetId("ABCDE-FGHIJ-KLMNO-PQRST")

	diags := resourceBigipSysBigiplicenseRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "ABCDE-FGHIJ-KLMNO-PQRST", d.Id(), "resource ID should be left untouched on successful read")
}

// TestResourceBigipSysBigiplicenseReadError covers the branch where
// Bigiplicenses() itself returns an error (e.g. the BIG-IP is unreachable or
// returns a non-JSON/error payload).
func TestResourceBigipSysBigiplicenseReadError(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})

	client := testBigiplicenseClient(t, server.URL)
	d := testBigiplicenseResourceData(t, "install", "ABCDE-FGHIJ-KLMNO-PQRST")
	d.SetId("ABCDE-FGHIJ-KLMNO-PQRST")

	diags := resourceBigipSysBigiplicenseRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when the API call fails")
}

// TestResourceBigipSysBigiplicenseUpdateSuccess covers the Update happy path:
// ModifyBigiplicense() succeeds, followed by a successful Read.
func TestResourceBigipSysBigiplicenseUpdateSuccess(t *testing.T) {
	setup()
	defer teardown()
	var sawPUT bool
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			sawPUT = true
			_, _ = fmt.Fprint(w, `{"registrationKey":"ABCDE-FGHIJ-KLMNO-PQRST","command":"install"}`)
		case "GET":
			_, _ = fmt.Fprint(w, `{"registrationKey":"ABCDE-FGHIJ-KLMNO-PQRST","command":"install"}`)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	client := testBigiplicenseClient(t, server.URL)
	d := testBigiplicenseResourceData(t, "install", "ABCDE-FGHIJ-KLMNO-PQRST")
	d.SetId("ABCDE-FGHIJ-KLMNO-PQRST")

	diags := resourceBigipSysBigiplicenseUpdate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPUT, "expected ModifyBigiplicense to issue a PUT request")
}

// TestResourceBigipSysBigiplicenseUpdateReadError covers the branch where
// ModifyBigiplicense() succeeds but the follow-up
// resourceBigipSysBigiplicenseRead call fails. Update calls Read
// unconditionally after a successful PUT, so a failing GET here must
// propagate as an error from Update.
func TestResourceBigipSysBigiplicenseUpdateReadError(t *testing.T) {
	setup()
	defer teardown()
	var sawPUT bool
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			sawPUT = true
			_, _ = fmt.Fprint(w, `{"registrationKey":"ABCDE-FGHIJ-KLMNO-PQRST","command":"install"}`)
		case "GET":
			http.Error(w, "internal error", http.StatusInternalServerError)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	})

	client := testBigiplicenseClient(t, server.URL)
	d := testBigiplicenseResourceData(t, "install", "ABCDE-FGHIJ-KLMNO-PQRST")
	d.SetId("ABCDE-FGHIJ-KLMNO-PQRST")

	diags := resourceBigipSysBigiplicenseUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected Update to propagate an error when the follow-up Read fails")
	assert.True(t, sawPUT, "expected ModifyBigiplicense to issue a PUT request before Read is attempted")
}

// TestResourceBigipSysBigiplicenseUpdateError covers the branch where
// ModifyBigiplicense() returns an error.
func TestResourceBigipSysBigiplicenseUpdateError(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			http.Error(w, `{"code":400,"message":"invalid license"}`, http.StatusBadRequest)
			return
		}
		t.Fatalf("unexpected method %s", r.Method)
	})

	client := testBigiplicenseClient(t, server.URL)
	d := testBigiplicenseResourceData(t, "install", "ABCDE-FGHIJ-KLMNO-PQRST")
	d.SetId("ABCDE-FGHIJ-KLMNO-PQRST")

	diags := resourceBigipSysBigiplicenseUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Update when ModifyBigiplicense fails")
}

// TestResourceBigipSysBigiplicenseDelete covers Delete, which is a no-op
// because the underlying API has no delete/unlicense operation.
func TestResourceBigipSysBigiplicenseDelete(t *testing.T) {
	d := testBigiplicenseResourceData(t, "install", "ABCDE-FGHIJ-KLMNO-PQRST")
	d.SetId("ABCDE-FGHIJ-KLMNO-PQRST")

	diags := resourceBigipSysBigiplicenseDelete(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	// Delete does not clear the ID or make any API call; it simply returns nil.
	assert.Equal(t, "ABCDE-FGHIJ-KLMNO-PQRST", d.Id())
}

// TestUnitBigipSysBigiplicenseInvalid confirms Terraform Core rejects an
// unknown schema key on the resource, exercised through the full Terraform
// SDK test harness (schema validation only, no API calls are made).
func testBigipSysBigiplicenseInvalid(registrationKey string) string {
	return fmt.Sprintf(`
resource "bigip_sys_bigiplicense" "test-license" {
  command           = "install"
  registration_key  = "%s"
  invalidkey        = "foo"
}
`, registrationKey)
}

func TestUnitBigipSysBigiplicenseInvalid(t *testing.T) {
	registrationKey := "ABCDE-FGHIJ-KLMNO-PQRST"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysBigiplicenseInvalid(registrationKey),
				ExpectError: regexp.MustCompile(`An argument named "invalidkey" is not expected here`),
			},
		},
	})
}
