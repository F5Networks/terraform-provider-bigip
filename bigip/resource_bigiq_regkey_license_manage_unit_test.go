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
	"net/http/httptest"
	"testing"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Several PostLicense-driven tests in this file unavoidably block for a
// real 5-second time.Sleep hardcoded inside the vendored
// go-bigip.PostLicense() function (vendor/github.com/f5devcentral/go-bigip/bigiq.go),
// which per repo convention (AGENTS.md) must not be patched. Rather than
// share the package-level mux/server singletons from
// resource_bigip_ltm_node_unit_test.go (which are not safe to use from
// t.Parallel() tests), every test in this file spins up its own local
// httptest server so the PostLicense-driven tests can run in parallel and
// keep the overall suite well under make test's 30-second budget.

// withZeroBigiqRevokeDelays temporarily overrides the atomic-backed revoke
// delay variables (declared in resource_bigiq_regkey_license_manage.go) to 0
// for the duration of the test, restoring the previous values on cleanup.
// The underlying values are atomic.Int64, so individual reads/writes never
// race/torn-read. However, the save-then-restore sequence here (Swap, then
// later Store the saved value back in t.Cleanup) is NOT itself atomic as a
// unit: if two tests both call this helper concurrently (e.g. two
// t.Parallel() tests overlap in time), one test's cleanup can stomp on the
// delay value the other test is still relying on being zeroed, or restore
// the wrong "previous" value. This is safe today only because no two
// callers of this helper are both marked t.Parallel() (as of this writing:
// TestResourceBigiqLicenseManageDeleteNoKeyManaged and
// ...DeleteNoKeyUnreachable are t.Parallel(), while
// TestWaitLicenseRevokeNoEntries and ...RetriesThenGivesUp are not, so the
// two parallel callers never overlap with any other caller). If a future
// change marks more than one caller of this helper t.Parallel()
// simultaneously, this helper will need a real mutex/semaphore around the
// swap-and-restore, not just atomic fields.
func withZeroBigiqRevokeDelays(t *testing.T) {
	t.Helper()
	origUnreachable := bigiqRevokeUnreachableDelayNanos.Swap(0)
	origRetry := bigiqRevokeRetryDelayNanos.Swap(0)
	t.Cleanup(func() {
		bigiqRevokeUnreachableDelayNanos.Store(origUnreachable)
		bigiqRevokeRetryDelayNanos.Store(origRetry)
	})
}

// regkeyLicenseTestClient builds a *bigip.BigIP pointed at the given mock
// server URL, representing the target BIG-IP device (passed as `meta` to
// the resource CRUD functions). This bypasses the provider's
// Client()/NewSession negotiation (and its SelfIP-based ValidateConnection
// call), consistent with the pattern used for bigip_sys_bigiplicense.
func regkeyLicenseTestClient(t *testing.T, url string) *bigip.BigIP {
	t.Helper()
	return bigip.NewSession(&bigip.Config{
		Address:  url,
		Username: "admin",
		Password: "admin",
		ConfigOptions: &bigip.ConfigOptions{
			APICallRetries: 1,
			APICallTimeout: 5 * time.Second,
		},
	})
}

// regkeyLicenseResourceData builds a *schema.ResourceData for
// bigip_bigiq_regkey_license_manage populated with the given raw values.
func regkeyLicenseResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := resourceBigiqLicenseManage()
	return schema.TestResourceDataRaw(t, r.Schema, raw)
}

func baseRegkeyLicenseRaw(url string) map[string]interface{} {
	return map[string]interface{}{
		"bigiq_address":    url,
		"bigiq_user":       "admin",
		"bigiq_password":   "admin",
		"bigiq_token_auth": false,
		"assignment_type":  "MANAGED",
		"license_poolname": "test-pool",
	}
}

// registerRegkeySelfIPOK registers a handler for /mgmt/tm/net/self, hit by
// connectBigIq -> Client() -> ValidateConnection() so the BIG-IQ connection
// step succeeds and the resource's own CRUD logic can be exercised.
func registerRegkeySelfIPOK(mux *http.ServeMux) {
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
}

// registerRegPools registers the three GetPoolType/GetRegkeyPoolId lookup
// endpoints (regkey, utility, purchased-pool), returning a single pool
// matching poolName/sortName under the regkey endpoint, and empty item
// lists for the other two so GetPoolType's later checks find nothing there.
func registerRegPools(mux *http.ServeMux, poolName, poolID, sortName string) {
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"items":[{"id":%q,"name":%q,"sortName":%q}]}`, poolID, poolName, sortName)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/utility/licenses", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/purchased-pool/licenses", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
}

// registerRegPoolsNotFound registers all three pool lookup endpoints
// returning empty item lists, simulating a pool name that doesn't exist
// anywhere.
func registerRegPoolsNotFound(mux *http.ServeMux) {
	for _, path := range []string{
		"/mgmt/cm/device/licensing/pool/regkey/licenses",
		"/mgmt/cm/device/licensing/pool/utility/licenses",
		"/mgmt/cm/device/licensing/pool/purchased-pool/licenses",
	} {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"items":[]}`)
		})
	}
}

// registerBigipLicenseOK registers a handler for /mgmt/tm/sys/license
// returning a benign license status (no "entries" key), used by
// GetBigipLiceseStatus in both Read's non-regKey branch and Delete's
// unreachable/waitLicenseRevoke branch. GetBigipLiceseStatus retries up to
// 15 times with a 10-second sleep between attempts if this endpoint ever
// returns a non-2xx response, so every test path that reaches it MUST
// register this handler to avoid a 150-second hang.
func registerBigipLicenseOK(mux *http.ServeMux) {
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"registrationKey":"ABCDE-FGHIJ-KLMNO-PQRST"}`)
	})
}

// TestResourceBigiqLicenseManageSchema exercises the schema definition of
// the bigip_bigiq_regkey_license_manage resource without needing any
// BIG-IQ connection.
func TestResourceBigiqLicenseManageSchema(t *testing.T) {
	r := resourceBigiqLicenseManage()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	requiredFields := []string{"bigiq_address", "bigiq_user", "bigiq_password", "assignment_type", "license_poolname"}
	for _, field := range requiredFields {
		s, ok := r.Schema[field]
		if !ok {
			t.Fatalf("Expected field '%s' to exist in schema", field)
		}
		if !s.Required {
			t.Errorf("Expected field '%s' to be required", field)
		}
	}

	if !r.Schema["bigiq_user"].Sensitive {
		t.Error("Expected bigiq_user to be Sensitive")
	}
	if !r.Schema["bigiq_password"].Sensitive {
		t.Error("Expected bigiq_password to be Sensitive")
	}
	if !r.Schema["device_license_status"].Computed {
		t.Error("Expected device_license_status to be Computed")
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

// ---------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------

// TestResourceBigiqLicenseManageCreateConnectError covers the branch where
// connectBigIq fails inside Create.
func TestResourceBigiqLicenseManageCreateConnectError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create when connectBigIq fails")
}

// TestResourceBigiqLicenseManageCreatePoolNotFound covers the branch where
// GetPoolType finds no matching pool anywhere, returning (nil, nil).
func TestResourceBigiqLicenseManageCreatePoolNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPoolsNotFound(mux)

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create when the pool is not found")
}

// TestResourceBigiqLicenseManageCreateUtilityMissingUnitOfMeasure covers the
// branch where the pool is a Utility pool but unit_of_measure is unset.
func TestResourceBigiqLicenseManageCreateUtilityMissingUnitOfMeasure(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "Utility")

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create for a Utility pool with no unit_of_measure")
}

// TestResourceBigiqLicenseManageCreateUnreachableMissingFields covers the
// branch where assignment_type is UNREACHABLE but mac_address/hypervisor
// are unset.
func TestResourceBigiqLicenseManageCreateUnreachableMissingFields(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "UNREACHABLE"
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create for UNREACHABLE with no mac_address/hypervisor")
}

// TestResourceBigiqLicenseManageCreateRegkeyPoolIdError covers the branch
// where GetRegkeyPoolId's underlying request fails. GetPoolType succeeds
// first (its own call to the regkey endpoint), then GetRegkeyPoolId's own
// call to the same endpoint fails.
func TestResourceBigiqLicenseManageCreateRegkeyPoolIdError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	var callCount int
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			// First call: GetPoolType's regkey lookup succeeds.
			_, _ = fmt.Fprint(w, `{"items":[{"id":"pool-id-1","name":"test-pool","sortName":"RegKey Pool"}]}`)
			return
		}
		// Second call: GetRegkeyPoolId's own lookup fails.
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/utility/licenses", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create when GetRegkeyPoolId fails")
}

// TestResourceBigiqLicenseManageCreateNoKeyPostLicense covers the branch
// where key is unset, so Create issues a PostLicense call (regardless of
// assignment_type) and sets the ID to the returned task ID, then Read
// succeeds via the FINISHED + Purchased Pool short-circuit.
//
// PostLicense (vendored go-bigip) unconditionally sleeps 5 seconds after a
// successful POST, so this test is marked Parallel to keep the overall
// suite fast; it uses its own local mux/server rather than the shared
// package-level singletons.
func TestResourceBigiqLicenseManageCreateNoKeyPostLicense(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "Purchased Pool")

	var sawPOST bool
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		sawPOST = true
		_, _ = fmt.Fprint(w, `{"id":"task-123"}`)
	})
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management/task-123", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"FINISHED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "MANAGED"
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPOST, "expected PostLicense to issue a POST request")
	assert.Equal(t, "task-123", d.Id())
	assert.Equal(t, "LICENSED", d.Get("device_license_status").(string))
}

// TestResourceBigiqLicenseManageCreateNoKeyPostLicenseError covers the
// branch where PostLicense itself returns an error.
func TestResourceBigiqLicenseManageCreateNoKeyPostLicenseError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "Purchased Pool")

	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"code":400,"message":"invalid task"}`, http.StatusBadRequest)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create when PostLicense fails")
}

// TestResourceBigiqLicenseManageCreateManaged covers the MANAGED assignment
// path with a non-empty key: GetDeviceId + RegkeylicenseAssign (POST +
// GetMemberStatus), followed by a successful Read via GetMemberStatus.
func TestResourceBigiqLicenseManageCreateManaged(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	mux.HandleFunc("/mgmt/shared/resolver/device-groups/cm-bigip-allBigIpDevices/devices", func(w http.ResponseWriter, r *http.Request) {
		host := getDeviceUri(server.URL)[2]
		_, _ = fmt.Fprintf(w, `{"items":[{"address":%q,"selfLink":"https://localhost/mgmt/shared/resolver/device-groups/cm-bigip-allBigIpDevices/devices/dev-1"}]}`, host)
	})

	var sawAssignPOST bool
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members", func(w http.ResponseWriter, r *http.Request) {
		sawAssignPOST = true
		_, _ = fmt.Fprint(w, `{"id":"mem-1"}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"mem-1","status":"LICENSED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "MANAGED"
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawAssignPOST, "expected RegkeylicenseAssign to issue a POST request")
	assert.Equal(t, "mem-1", d.Id())
}

// TestResourceBigiqLicenseManageCreateUnmanaged covers the UNMANAGED
// assignment path with a non-empty key: RegkeylicenseAssign is called
// directly (no GetDeviceId lookup).
func TestResourceBigiqLicenseManageCreateUnmanaged(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	var sawAssignPOST bool
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members", func(w http.ResponseWriter, r *http.Request) {
		sawAssignPOST = true
		_, _ = fmt.Fprint(w, `{"id":"mem-2"}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-2", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"mem-2","status":"LICENSED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "UNMANAGED"
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawAssignPOST, "expected RegkeylicenseAssign to issue a POST request")
	assert.Equal(t, "mem-2", d.Id())
}

// TestResourceBigiqLicenseManageCreateManagedAssignError covers the branch
// where RegkeylicenseAssign fails for a MANAGED assignment.
func TestResourceBigiqLicenseManageCreateManagedAssignError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	mux.HandleFunc("/mgmt/shared/resolver/device-groups/cm-bigip-allBigIpDevices/devices", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"code":400,"message":"assign failed"}`, http.StatusBadRequest)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "MANAGED"
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create when RegkeylicenseAssign fails")
}

// TestResourceBigiqLicenseManageCreateUnreachableFailed covers the branch
// where assignment_type is UNREACHABLE, GetLicenseStatus reports FAILED,
// and Create clears the resource ID and returns an error.
//
// PostLicense unconditionally sleeps 5 seconds; marked Parallel.
func TestResourceBigiqLicenseManageCreateUnreachableFailed(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"task-999"}`)
	})
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management/task-999", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"FAILED","errorMessage":"license invalid"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "UNREACHABLE"
	raw["mac_address"] = "00:11:22:33:44:55"
	raw["hypervisor"] = "aws"
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Create when license status is FAILED")
	assert.Empty(t, d.Id(), "resource ID should be cleared when license status is FAILED")
}

// TestResourceBigiqLicenseManageCreateUnreachableSuccess covers the branch
// where assignment_type is UNREACHABLE and the license status resolves
// successfully, triggering InstallLicense on the target BIG-IP.
//
// PostLicense unconditionally sleeps 5 seconds; marked Parallel.
func TestResourceBigiqLicenseManageCreateUnreachableSuccess(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")
	registerBigipLicenseOK(mux)

	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"task-777"}`)
	})
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management/task-777", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"FINISHED","licenseText":"license-blob","licenseAssignmentReference":{"link":"https://localhost/mgmt/cm/device/tasks/licensing/pool/member-management/task-777/status"}}`)
	})
	var sawInstallPUT bool
	mux.HandleFunc("/mgmt/tm/shared/licensing/registration", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			sawInstallPUT = true
		}
	})
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management/task-777/status", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"LICENSED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "UNREACHABLE"
	raw["mac_address"] = "00:11:22:33:44:55"
	raw["hypervisor"] = "aws"
	raw["license_poolname"] = "test-pool"
	d := regkeyLicenseResourceData(t, raw)

	diags := resourceBigiqLicenseManageCreate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawInstallPUT, "expected InstallLicense to issue a PUT request")
	assert.Equal(t, "task-777", d.Id())
}

// ---------------------------------------------------------------------
// Read
// ---------------------------------------------------------------------

// TestResourceBigiqLicenseManageReadConnectError covers the branch where
// connectBigIq fails inside Read.
func TestResourceBigiqLicenseManageReadConnectError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-123")

	diags := resourceBigiqLicenseManageRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when connectBigIq fails")
}

// TestResourceBigiqLicenseManageReadPoolNotFound covers the branch where
// GetPoolType finds no matching pool.
func TestResourceBigiqLicenseManageReadPoolNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPoolsNotFound(mux)

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-123")

	diags := resourceBigiqLicenseManageRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when the pool is not found")
}

// TestResourceBigiqLicenseManageReadFinishedPurchasedPool covers the
// short-circuit branch where the license task status is FINISHED and the
// pool is a Purchased Pool, setting device_license_status to LICENSED
// without further BIG-IP status calls.
func TestResourceBigiqLicenseManageReadFinishedPurchasedPool(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "Purchased Pool")
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management/task-123", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"FINISHED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-123")

	diags := resourceBigiqLicenseManageRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "LICENSED", d.Get("device_license_status").(string))
}

// TestResourceBigiqLicenseManageReadFailed covers the branch where the
// license task status is FAILED, clearing the resource ID and returning an
// error.
func TestResourceBigiqLicenseManageReadFailed(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management/task-123", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"FAILED","errorMessage":"boom"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-123")

	diags := resourceBigiqLicenseManageRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when license status is FAILED")
	assert.Empty(t, d.Id())
}

// TestResourceBigiqLicenseManageReadDeviceStatus covers the general (non
// Purchased-Pool-FINISHED) status branch: GetDeviceLicenseStatus and
// GetBigipLiceseStatus are both called, and device_license_status is set to
// the device status.
func TestResourceBigiqLicenseManageReadDeviceStatus(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")
	registerBigipLicenseOK(mux)
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management/task-123", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"FINISHED","licenseAssignmentReference":{"link":"https://localhost/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-1"}}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"LICENSED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-123")

	diags := resourceBigiqLicenseManageRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.Equal(t, "LICENSED", d.Get("device_license_status").(string))
}

// TestResourceBigiqLicenseManageReadWithKey covers the branch where key is
// set, so Read calls GetMemberStatus directly using the poolId/regKey/memID.
func TestResourceBigiqLicenseManageReadWithKey(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"mem-1","status":"LICENSED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("mem-1")

	diags := resourceBigiqLicenseManageRead(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

// TestResourceBigiqLicenseManageReadWithKeyError covers the branch where
// key is set and GetMemberStatus fails (INSTALLATION_FAILED).
func TestResourceBigiqLicenseManageReadWithKeyError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"mem-1","status":"INSTALLATION_FAILED","message":"boom"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("mem-1")

	diags := resourceBigiqLicenseManageRead(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Read when GetMemberStatus reports INSTALLATION_FAILED")
}

// ---------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------

// TestResourceBigiqLicenseManageUpdateConnectError covers the branch where
// connectBigIq fails inside Update.
func TestResourceBigiqLicenseManageUpdateConnectError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-123")

	diags := resourceBigiqLicenseManageUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Update when connectBigIq fails")
}

// TestResourceBigiqLicenseManageUpdateNoKeyPostLicense covers the branch
// where key is unset, so Update issues a PostLicense call and sets the ID
// to the returned task ID, then Read succeeds via the FINISHED + Purchased
// Pool short-circuit.
//
// PostLicense unconditionally sleeps 5 seconds; marked Parallel.
func TestResourceBigiqLicenseManageUpdateNoKeyPostLicense(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "Purchased Pool")

	var sawPOST bool
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		sawPOST = true
		_, _ = fmt.Fprint(w, `{"id":"task-456"}`)
	})
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management/task-456", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"status":"FINISHED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-old")

	diags := resourceBigiqLicenseManageUpdate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPOST, "expected PostLicense to issue a POST request")
	assert.Equal(t, "task-456", d.Id())
}

// TestResourceBigiqLicenseManageUpdateNoKeyPostLicenseError covers the
// branch where PostLicense itself returns an error during Update.
func TestResourceBigiqLicenseManageUpdateNoKeyPostLicenseError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "Purchased Pool")

	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"code":400,"message":"invalid task"}`, http.StatusBadRequest)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-old")

	diags := resourceBigiqLicenseManageUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Update when PostLicense fails")
}

// TestResourceBigiqLicenseManageUpdateManaged covers the MANAGED assignment
// path with a non-empty key during Update.
func TestResourceBigiqLicenseManageUpdateManaged(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	mux.HandleFunc("/mgmt/shared/resolver/device-groups/cm-bigip-allBigIpDevices/devices", func(w http.ResponseWriter, r *http.Request) {
		host := getDeviceUri(server.URL)[2]
		_, _ = fmt.Fprintf(w, `{"items":[{"address":%q,"selfLink":"https://localhost/mgmt/shared/resolver/device-groups/cm-bigip-allBigIpDevices/devices/dev-1"}]}`, host)
	})
	var sawAssignPOST bool
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members", func(w http.ResponseWriter, r *http.Request) {
		sawAssignPOST = true
		_, _ = fmt.Fprint(w, `{"id":"mem-3"}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-3", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"mem-3","status":"LICENSED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "MANAGED"
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("mem-old")

	diags := resourceBigiqLicenseManageUpdate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawAssignPOST, "expected RegkeylicenseAssign to issue a POST request")
	assert.Equal(t, "mem-3", d.Id())
}

// TestResourceBigiqLicenseManageUpdateUnmanaged covers the UNMANAGED
// assignment path with a non-empty key during Update.
func TestResourceBigiqLicenseManageUpdateUnmanaged(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	var sawAssignPOST bool
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members", func(w http.ResponseWriter, r *http.Request) {
		sawAssignPOST = true
		_, _ = fmt.Fprint(w, `{"id":"mem-4"}`)
	})
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-4", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"mem-4","status":"LICENSED"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "UNMANAGED"
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("mem-old")

	diags := resourceBigiqLicenseManageUpdate(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawAssignPOST, "expected RegkeylicenseAssign to issue a POST request")
	assert.Equal(t, "mem-4", d.Id())
}

// TestResourceBigiqLicenseManageUpdatePoolNotFound covers the branch where
// GetPoolType finds no matching pool during Update.
func TestResourceBigiqLicenseManageUpdatePoolNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPoolsNotFound(mux)

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-old")

	diags := resourceBigiqLicenseManageUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Update when the pool is not found")
}

// TestResourceBigiqLicenseManageUpdateUtilityMissingUnitOfMeasure covers the
// branch where the pool is a Utility pool but unit_of_measure is unset,
// during Update.
func TestResourceBigiqLicenseManageUpdateUtilityMissingUnitOfMeasure(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "Utility")

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-old")

	diags := resourceBigiqLicenseManageUpdate(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Update for a Utility pool with no unit_of_measure")
}

// ---------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------

// TestResourceBigiqLicenseManageDeleteConnectError covers the branch where
// connectBigIq fails inside Delete.
func TestResourceBigiqLicenseManageDeleteConnectError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-123")

	diags := resourceBigiqLicenseManageDelete(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Delete when connectBigIq fails")
}

// TestResourceBigiqLicenseManageDeletePoolIdError covers the branch where
// looking up the pool ID via GetRegkeyPoolId fails during Delete.
func TestResourceBigiqLicenseManageDeletePoolIdError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-123")

	diags := resourceBigiqLicenseManageDelete(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Delete when GetRegkeyPoolId fails")
}

// TestResourceBigiqLicenseManageDeleteNoKeyManaged covers the branch where
// key is unset (revoke path): PostLicense with command=revoke, followed by
// waitLicenseRevoke succeeding (GetBigipLiceseStatus with no "entries" key).
//
// PostLicense unconditionally sleeps 5 seconds; marked Parallel.
func TestResourceBigiqLicenseManageDeleteNoKeyManaged(t *testing.T) {
	t.Parallel()
	withZeroBigiqRevokeDelays(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")
	registerBigipLicenseOK(mux)

	var sawPOST bool
	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		sawPOST = true
		_, _ = fmt.Fprint(w, `{"id":"task-revoke-1"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "MANAGED"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-revoke-1")

	diags := resourceBigiqLicenseManageDelete(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPOST, "expected PostLicense to issue a revoke POST request")
	assert.Empty(t, d.Id(), "resource ID should be cleared on successful delete")
}

// TestResourceBigiqLicenseManageDeleteNoKeyUnreachable covers the branch
// where assignment_type is UNREACHABLE and key is unset: PostLicense,
// RevokeLicense (DELETE on the target BIG-IP), then waitLicenseRevoke.
//
// PostLicense unconditionally sleeps 5 seconds; marked Parallel.
func TestResourceBigiqLicenseManageDeleteNoKeyUnreachable(t *testing.T) {
	t.Parallel()
	withZeroBigiqRevokeDelays(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")
	registerBigipLicenseOK(mux)

	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"task-revoke-2"}`)
	})
	var sawRevokeDELETE bool
	mux.HandleFunc("/mgmt/tm/shared/licensing/registration", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			sawRevokeDELETE = true
		}
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "UNREACHABLE"
	raw["mac_address"] = "00:11:22:33:44:55"
	raw["hypervisor"] = "aws"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-revoke-2")

	diags := resourceBigiqLicenseManageDelete(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawRevokeDELETE, "expected RevokeLicense to issue a DELETE request")
	assert.Empty(t, d.Id())
}

// TestResourceBigiqLicenseManageDeleteNoKeyPostLicenseError covers the
// branch where PostLicense (revoke) itself returns an error.
func TestResourceBigiqLicenseManageDeleteNoKeyPostLicenseError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"code":400,"message":"revoke failed"}`, http.StatusBadRequest)
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-revoke-3")

	diags := resourceBigiqLicenseManageDelete(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Delete when PostLicense (revoke) fails")
}

// TestResourceBigiqLicenseManageDeleteNoKeyUnreachableRevokeError covers
// the branch where RevokeLicense on the target BIG-IP fails.
//
// PostLicense unconditionally sleeps 5 seconds; marked Parallel.
func TestResourceBigiqLicenseManageDeleteNoKeyUnreachableRevokeError(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	mux.HandleFunc("/mgmt/cm/device/tasks/licensing/pool/member-management", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"id":"task-revoke-4"}`)
	})
	mux.HandleFunc("/mgmt/tm/shared/licensing/registration", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			http.Error(w, `{"code":400,"message":"revoke failed"}`, http.StatusBadRequest)
		}
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "UNREACHABLE"
	raw["mac_address"] = "00:11:22:33:44:55"
	raw["hypervisor"] = "aws"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("task-revoke-4")

	diags := resourceBigiqLicenseManageDelete(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Delete when RevokeLicense fails")
}

// TestResourceBigiqLicenseManageDeleteWithKeyManaged covers the branch
// where key is set and assignment_type is MANAGED: RegkeylicenseRevoke
// (DELETE + GET) is called directly.
func TestResourceBigiqLicenseManageDeleteWithKeyManaged(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	var sawDELETE bool
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-5", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			sawDELETE = true
		case "GET":
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "MANAGED"
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("mem-5")

	diags := resourceBigiqLicenseManageDelete(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawDELETE, "expected RegkeylicenseRevoke to issue a DELETE request")
	assert.Empty(t, d.Id())
}

// TestResourceBigiqLicenseManageDeleteWithKeyManagedError covers the branch
// where RegkeylicenseRevoke itself returns an error.
func TestResourceBigiqLicenseManageDeleteWithKeyManagedError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-6", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "MANAGED"
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("mem-6")

	diags := resourceBigiqLicenseManageDelete(context.Background(), d, client)

	require.True(t, diags.HasError(), "expected an error from Delete when RegkeylicenseRevoke fails")
}

// TestResourceBigiqLicenseManageDeleteWithKeyUnmanaged covers the branch
// where key is set and assignment_type is UNMANAGED: LicenseRevoke (DELETE
// with body) is called; its error is intentionally ignored by the
// resource, so Delete always succeeds along this path.
func TestResourceBigiqLicenseManageDeleteWithKeyUnmanaged(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)
	registerRegPools(mux, "test-pool", "pool-id-1", "RegKey Pool")

	var sawDELETE bool
	mux.HandleFunc("/mgmt/cm/device/licensing/pool/regkey/licenses/pool-id-1/offerings/REGKEY123/members/mem-7", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			sawDELETE = true
		case "GET":
			_, _ = fmt.Fprint(w, `{}`)
		}
	})

	client := regkeyLicenseTestClient(t, server.URL)
	raw := baseRegkeyLicenseRaw(server.URL)
	raw["assignment_type"] = "UNMANAGED"
	raw["key"] = "REGKEY123"
	d := regkeyLicenseResourceData(t, raw)
	d.SetId("mem-7")

	diags := resourceBigiqLicenseManageDelete(context.Background(), d, client)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawDELETE, "expected LicenseRevoke to issue a DELETE request")
	assert.Empty(t, d.Id())
}

// ---------------------------------------------------------------------
// waitLicenseRevoke
// ---------------------------------------------------------------------

// TestWaitLicenseRevokeNoEntries covers the branch where
// GetBigipLiceseStatus returns a status with no "entries" key on the first
// call, so waitLicenseRevoke returns immediately without retrying.
func TestWaitLicenseRevokeNoEntries(t *testing.T) {
	withZeroBigiqRevokeDelays(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	var callCount int
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		_, _ = fmt.Fprint(w, `{"registrationKey":"ABCDE"}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)

	result, err := waitLicenseRevoke(client)

	require.NoError(t, err)
	assert.NotContains(t, result, "entries")
	assert.Equal(t, 1, callCount, "expected exactly one GetBigipLiceseStatus call when entries is absent")
}

// TestWaitLicenseRevokeRetriesThenGivesUp covers the branch where
// GetBigipLiceseStatus keeps returning an "entries" key, so
// waitLicenseRevoke retries up to 3 times (bounded by the retry delay,
// zeroed here) before giving up and returning the last (still-"entries")
// status.
func TestWaitLicenseRevokeRetriesThenGivesUp(t *testing.T) {
	withZeroBigiqRevokeDelays(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	var callCount int
	mux.HandleFunc("/mgmt/tm/sys/license", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		_, _ = fmt.Fprint(w, `{"entries":{"foo":"bar"}}`)
	})

	client := regkeyLicenseTestClient(t, server.URL)

	result, err := waitLicenseRevoke(client)

	require.NoError(t, err)
	assert.Contains(t, result, "entries")
	assert.Equal(t, 4, callCount, "expected the initial call plus 3 retries")
}

// ---------------------------------------------------------------------
// connectBigIq
// ---------------------------------------------------------------------

// TestConnectBigIqRegkeyTokenAuth exercises the bigiq_token_auth=true
// branch of connectBigIq, which routes through bigip.NewTokenSession
// instead of bigip.NewSession.
func TestConnectBigIqRegkeyTokenAuth(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"token":{"token":"unit-test-token"},"timeout":{"timeout":1200}}`)
	})
	registerRegkeySelfIPOK(mux)

	raw := baseRegkeyLicenseRaw(server.URL)
	raw["bigiq_token_auth"] = true
	raw["bigiq_login_ref"] = "tmos"
	d := regkeyLicenseResourceData(t, raw)

	client, err := connectBigIq(d)

	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "unit-test-token", client.Token)
}

// TestConnectBigIqRegkeyBasicAuth exercises the default
// (bigiq_token_auth=false) branch of connectBigIq, which routes through
// bigip.NewSession.
func TestConnectBigIqRegkeyBasicAuth(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	registerRegkeySelfIPOK(mux)

	d := regkeyLicenseResourceData(t, baseRegkeyLicenseRaw(server.URL))

	client, err := connectBigIq(d)

	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Empty(t, client.Token)
}

// TestConnectBigIqRegkeyConnectionError covers the branch where
// ValidateConnection fails (e.g. the mock BIG-IQ returns an error for
// /mgmt/tm/net/self).
func TestConnectBigIqRegkeyConnectionError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	d := regkeyLicenseResourceData(t, baseRegkeyLicenseRaw(server.URL))

	_, err := connectBigIq(d)

	require.Error(t, err)
}
