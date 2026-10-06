/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
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

// TestResourceBigipSysSnmpTrapsSchema exercises the schema definition of the
// bigip_sys_snmp_traps resource without needing any BIG-IP connection.
func TestResourceBigipSysSnmpTrapsSchema(t *testing.T) {
	r := resourceBigipSysSnmpTraps()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	optionalStringFields := []string{
		"name",
		"auth_passwordencrypted",
		"auth_protocol",
		"community",
		"description",
		"engine_id",
		"host",
		"privacy_password",
		"privacy_password_encrypted",
		"privacy_protocol",
		"security_level",
		"security_name",
		"version",
	}
	for _, field := range optionalStringFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if s.Required {
			t.Errorf("Expected field '%s' to be optional", field)
		}
		if s.Type != schema.TypeString {
			t.Errorf("Expected field '%s' to be TypeString, got %v", field, s.Type)
		}
		if s.Description == "" {
			t.Errorf("Expected field '%s' to have a description", field)
		}
	}

	portSchema, ok := r.Schema["port"]
	if !ok {
		t.Fatal("Expected field 'port' to exist in schema")
	}
	if portSchema.Type != schema.TypeInt {
		t.Errorf("Expected field 'port' to be TypeInt, got %v", portSchema.Type)
	}
	if portSchema.Required {
		t.Error("Expected field 'port' to be optional")
	}

	// security_level and version are Optional+Computed.
	for _, field := range []string{"security_level", "version"} {
		s := r.Schema[field]
		if !s.Computed {
			t.Errorf("Expected field '%s' to be Computed", field)
		}
		if !s.Optional {
			t.Errorf("Expected field '%s' to be Optional", field)
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

// testSnmpTrapsClient builds a *bigip.BigIP pointed at the given mock server
// URL (typically server.URL from the shared setup()/teardown() mux),
// bypassing the provider's Client()/NewSession negotiation (and its
// SelfIP-based ValidateConnection call) since these unit tests invoke the
// resource CRUD functions directly rather than through resource.Test.
func testSnmpTrapsClient(t *testing.T, url string) *bigip.BigIP {
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

func testSnmpTrapsResourceData(t *testing.T, values map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := resourceBigipSysSnmpTraps()
	return schema.TestResourceDataRaw(t, r.Schema, values)
}

func snmpTrapTestValues(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":                       name,
		"auth_passwordencrypted":     "encryptedpass",
		"auth_protocol":              "sha",
		"community":                  "public",
		"description":                "test trap",
		"engine_id":                  "0x80001f8880",
		"host":                       "192.168.1.100",
		"port":                       162,
		"privacy_password":           "privpass",
		"privacy_password_encrypted": "encryptedprivpass",
		"privacy_protocol":           "aes",
		"security_level":             "authPriv",
		"security_name":              "secname",
		"version":                    "3",
	}
}

// TestResourceBigipSysSnmpTrapsCreateSuccess covers the Create happy path:
// CreateTRAP() succeeds, the resource ID is set to the trap name, and the
// follow-up Read succeeds and populates state from the API response.
func TestResourceBigipSysSnmpTrapsCreateSuccess(t *testing.T) {
	setup()
	defer teardown()

	var sawPOST bool
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		sawPOST = true
		_, _ = fmt.Fprint(w, `{"name":"trap1"}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/trap1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		_, _ = fmt.Fprint(w, `{
  "name": "trap1",
  "authPasswordEncrypted": "encryptedpass",
  "authProtocol": "sha",
  "community": "public",
  "description": "test trap",
  "engineId": "0x80001f8880",
  "host": "192.168.1.100",
  "port": 162,
  "privacyPassword": "privpass",
  "privacyPasswordEncrypted": "encryptedprivpass",
  "privacyProtocol": "aes",
  "securityLevel": "authPriv",
  "SecurityName": "secname",
  "version": "3"
}`)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))

	diags := resourceBigipSysSnmpTrapsCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPOST, "expected CreateTRAP to issue a POST request")
	assert.Equal(t, "trap1", d.Id())
	assert.Equal(t, "192.168.1.100", d.Get("host"))
	assert.Equal(t, 162, d.Get("port"))
	assert.Equal(t, "secname", d.Get("security_name"))
}

// TestResourceBigipSysSnmpTrapsCreateError covers the branch where
// CreateTRAP() returns an error; the resource ID is left unset.
func TestResourceBigipSysSnmpTrapsCreateError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/mgmt/tm/sys/snmp/traps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		http.Error(w, `{"code":400,"message":"invalid host"}`, http.StatusBadRequest)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))

	diags := resourceBigipSysSnmpTrapsCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create when CreateTRAP fails")
	assert.Empty(t, d.Id(), "resource ID should remain unset when Create fails")
}

// TestResourceBigipSysSnmpTrapsCreateReadError covers the branch where
// CreateTRAP() succeeds but the follow-up resourceBigipSysSnmpTrapsRead call
// fails. Create calls Read unconditionally after a successful POST, so a
// failing GET here must propagate as an error from Create even though the
// resource ID has already been set.
func TestResourceBigipSysSnmpTrapsCreateReadError(t *testing.T) {
	setup()
	defer teardown()

	var sawPOST bool
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		sawPOST = true
		_, _ = fmt.Fprint(w, `{"name":"trap1"}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/trap1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "internal error", http.StatusInternalServerError)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))

	diags := resourceBigipSysSnmpTrapsCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected Create to propagate an error when the follow-up Read fails")
	assert.True(t, sawPOST, "expected CreateTRAP to issue a POST request before Read is attempted")
	assert.Equal(t, "trap1", d.Id(), "resource ID is set before Read runs, even though Read subsequently fails")
}

// TestResourceBigipSysSnmpTrapsReadSuccess covers the Read happy path:
// TRAPs() succeeds and returns a non-nil trap, so all fields are written to
// state.
func TestResourceBigipSysSnmpTrapsReadSuccess(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/trap1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		_, _ = fmt.Fprint(w, `{
  "name": "trap1",
  "authPasswordEncrypted": "encryptedpass",
  "authProtocol": "sha",
  "community": "public",
  "description": "test trap",
  "engineId": "0x80001f8880",
  "host": "192.168.1.100",
  "port": 162,
  "privacyPassword": "privpass",
  "privacyPasswordEncrypted": "encryptedprivpass",
  "privacyProtocol": "aes",
  "securityLevel": "authPriv",
  "SecurityName": "secname",
  "version": "3"
}`)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))
	d.SetId("trap1")

	diags := resourceBigipSysSnmpTrapsRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "trap1", d.Id())
	assert.Equal(t, "trap1", d.Get("name"))
	assert.Equal(t, "encryptedpass", d.Get("auth_passwordencrypted"))
	assert.Equal(t, "sha", d.Get("auth_protocol"))
	assert.Equal(t, "public", d.Get("community"))
	assert.Equal(t, "test trap", d.Get("description"))
	assert.Equal(t, "0x80001f8880", d.Get("engine_id"))
	assert.Equal(t, "192.168.1.100", d.Get("host"))
	assert.Equal(t, 162, d.Get("port"))
	assert.Equal(t, "privpass", d.Get("privacy_password"))
	assert.Equal(t, "encryptedprivpass", d.Get("privacy_password_encrypted"))
	assert.Equal(t, "aes", d.Get("privacy_protocol"))
	assert.Equal(t, "authPriv", d.Get("security_level"))
	assert.Equal(t, "secname", d.Get("security_name"))
	assert.Equal(t, "3", d.Get("version"))
}

// TestResourceBigipSysSnmpTrapsReadError covers the branch where TRAPs()
// itself returns an error (e.g. the BIG-IP is unreachable or returns a
// non-JSON/error payload).
func TestResourceBigipSysSnmpTrapsReadError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/trap1", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))
	d.SetId("trap1")

	diags := resourceBigipSysSnmpTrapsRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when the API call fails")
}

// TestResourceBigipSysSnmpTrapsUpdateSuccess covers the Update happy path:
// ModifyTRAP() succeeds, followed by a successful Read.
func TestResourceBigipSysSnmpTrapsUpdateSuccess(t *testing.T) {
	setup()
	defer teardown()

	var sawPATCH bool
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Expected method 'PATCH', got %s", r.Method)
		sawPATCH = true
		_, _ = fmt.Fprint(w, `{"name":"trap1"}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/trap1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		_, _ = fmt.Fprint(w, `{"name":"trap1","host":"192.168.1.200","port":163}`)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))
	d.SetId("trap1")

	diags := resourceBigipSysSnmpTrapsUpdate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPATCH, "expected ModifyTRAP to issue a PATCH request")
	assert.Equal(t, "192.168.1.200", d.Get("host"))
	assert.Equal(t, 163, d.Get("port"))
}

// TestResourceBigipSysSnmpTrapsUpdateReadError covers the branch where
// ModifyTRAP() succeeds but the follow-up resourceBigipSysSnmpTrapsRead call
// fails. Update calls Read unconditionally after a successful PATCH, so a
// failing GET here must propagate as an error from Update.
func TestResourceBigipSysSnmpTrapsUpdateReadError(t *testing.T) {
	setup()
	defer teardown()

	var sawPATCH bool
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Expected method 'PATCH', got %s", r.Method)
		sawPATCH = true
		_, _ = fmt.Fprint(w, `{"name":"trap1"}`)
	})
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/trap1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		http.Error(w, "internal error", http.StatusInternalServerError)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))
	d.SetId("trap1")

	diags := resourceBigipSysSnmpTrapsUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected Update to propagate an error when the follow-up Read fails")
	assert.True(t, sawPATCH, "expected ModifyTRAP to issue a PATCH request before Read is attempted")
}

// TestResourceBigipSysSnmpTrapsUpdateError covers the branch where
// ModifyTRAP() returns an error.
func TestResourceBigipSysSnmpTrapsUpdateError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/mgmt/tm/sys/snmp/traps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Expected method 'PATCH', got %s", r.Method)
		http.Error(w, `{"code":400,"message":"invalid trap configuration"}`, http.StatusBadRequest)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))
	d.SetId("trap1")

	diags := resourceBigipSysSnmpTrapsUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Update when ModifyTRAP fails")
}

// TestResourceBigipSysSnmpTrapsDeleteSuccess covers the Delete happy path:
// DeleteTRAP() succeeds and the resource ID is cleared.
func TestResourceBigipSysSnmpTrapsDeleteSuccess(t *testing.T) {
	setup()
	defer teardown()

	var sawDELETE bool
	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/trap1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "Expected method 'DELETE', got %s", r.Method)
		sawDELETE = true
		w.WriteHeader(http.StatusOK)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))
	d.SetId("trap1")

	diags := resourceBigipSysSnmpTrapsDelete(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawDELETE, "expected DeleteTRAP to issue a DELETE request")
	assert.Empty(t, d.Id(), "resource ID should be cleared after successful delete")
}

// TestResourceBigipSysSnmpTrapsDeleteError covers the branch where
// DeleteTRAP() returns an error; the resource ID is left untouched.
func TestResourceBigipSysSnmpTrapsDeleteError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/mgmt/tm/sys/snmp/traps/trap1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "Expected method 'DELETE', got %s", r.Method)
		http.Error(w, `{"code":404,"message":"trap not found"}`, http.StatusNotFound)
	})

	client := testSnmpTrapsClient(t, server.URL)
	d := testSnmpTrapsResourceData(t, snmpTrapTestValues("trap1"))
	d.SetId("trap1")

	diags := resourceBigipSysSnmpTrapsDelete(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Delete when DeleteTRAP fails")
	assert.Equal(t, "trap1", d.Id(), "resource ID should remain set when Delete fails")
}

func testBigipSysSnmpTrapsInvalid(name string) string {
	return fmt.Sprintf(`
resource "bigip_sys_snmp_traps" "test-trap" {
  name       = "%s"
  host       = "192.168.1.100"
  invalidkey = "foo"
}
`, name)
}

// TestUnitBigipSysSnmpTrapsInvalid confirms Terraform Core rejects an
// unknown schema key on the resource (schema validation only, no API calls
// are made).
func TestUnitBigipSysSnmpTrapsInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipSysSnmpTrapsInvalid("trap1"),
				ExpectError: regexp.MustCompile(`An argument named "invalidkey" is not expected here`),
			},
		},
	})
}
