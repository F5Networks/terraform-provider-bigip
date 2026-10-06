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
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setBigiqAs3RetryDelayForTest atomically overrides bigiqAs3RetryDelayNanos
// (declared in resource_bigiq_as3.go) and returns a restore func that puts
// back the previous value. It goes through the atomic.Int64 rather than a
// bare package variable so overriding it is race-free even if a test in
// this file is later marked t.Parallel().
func setBigiqAs3RetryDelayForTest(d time.Duration) (restore func()) {
	previous := bigiqAs3RetryDelayNanos.Swap(int64(d))
	return func() { bigiqAs3RetryDelayNanos.Store(previous) }
}

// withZeroBigiqAs3RetryDelay temporarily overrides the retry delay to 0 for
// the duration of the test so that Read/Update/Delete (which all wait on it
// before contacting BIG-IQ) can be exercised without blocking the test run.
func withZeroBigiqAs3RetryDelay(t *testing.T) {
	t.Helper()
	restore := setBigiqAs3RetryDelayForTest(0)
	t.Cleanup(restore)
}

// TestWaitBigiqAs3RetryDelayCompletes covers the normal path where the
// delay elapses before ctx is cancelled.
func TestWaitBigiqAs3RetryDelayCompletes(t *testing.T) {
	restore := setBigiqAs3RetryDelayForTest(time.Millisecond)
	defer restore()

	err := waitBigiqAs3RetryDelay(context.Background())

	require.NoError(t, err)
}

// TestWaitBigiqAs3RetryDelayContextCancelled covers the branch where ctx is
// cancelled (or its deadline expires) before the retry delay elapses;
// waitBigiqAs3RetryDelay must return promptly with ctx.Err() rather than
// blocking for the full delay.
func TestWaitBigiqAs3RetryDelayContextCancelled(t *testing.T) {
	restore := setBigiqAs3RetryDelayForTest(time.Hour)
	defer restore()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)
	go func() { done <- waitBigiqAs3RetryDelay(ctx) }()

	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("waitBigiqAs3RetryDelay did not return promptly after ctx was cancelled")
	}
}

// TestResourceBigiqAs3ReadContextCancelled covers Read's propagation of a
// cancelled context as an error diagnostic instead of blocking for the full
// retry delay.
func TestResourceBigiqAs3ReadContextCancelled(t *testing.T) {
	restore := setBigiqAs3RetryDelayForTest(time.Hour)
	defer restore()

	d := bigiqAs3ResourceData(t, baseBigiqAs3Raw("http://192.0.2.1", testBigiqAs3Json))
	d.SetId("192.0.2.50_MyTenant")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan diag.Diagnostics, 1)
	go func() { done <- resourceBigiqAs3Read(ctx, d, nil) }()

	select {
	case diags := <-done:
		require.True(t, diags.HasError(), "expected Read to return an error diagnostic when ctx is cancelled")
	case <-time.After(2 * time.Second):
		t.Fatal("resourceBigiqAs3Read did not return promptly after ctx was cancelled")
	}
}

// bigiqAs3ResourceData builds a *schema.ResourceData for bigip_bigiq_as3
// populated with the given raw values, defaulting bigiq_address/user/password
// to point at the mock server when not overridden.
func bigiqAs3ResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := resourceBigiqAs3()
	return schema.TestResourceDataRaw(t, r.Schema, raw)
}

func baseBigiqAs3Raw(url string, as3Json string) map[string]interface{} {
	return map[string]interface{}{
		"bigiq_address":    url,
		"bigiq_user":       "admin",
		"bigiq_password":   "admin",
		"bigiq_token_auth": false,
		"as3_json":         as3Json,
		"ignore_metadata":  true,
	}
}

const testBigiqAs3Json = `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","MyTenant":{"class":"Tenant","app":{"class":"Application"}}}}`

// testBigiqAs3JsonTwoTenants declares two tenants so that PostAs3Bigiq can
// exercise a genuine "Partial Success" (one tenant succeeds, one fails)
// rather than the all-succeed/all-fail branches.
const testBigiqAs3JsonTwoTenants = `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","MyTenant":{"class":"Tenant","app":{"class":"Application"}},"OtherTenant":{"class":"Tenant","app":{"class":"Application"}}}}`

// TestResourceBigiqAs3Schema exercises the schema definition of the
// bigip_bigiq_as3 resource without needing any BIG-IQ connection.
func TestResourceBigiqAs3Schema(t *testing.T) {
	r := resourceBigiqAs3()

	if r.Schema == nil {
		t.Fatal("Expected schema to be defined")
	}

	requiredFields := []string{"bigiq_address", "bigiq_user", "bigiq_password", "as3_json"}
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

	if r.Schema["tenant_list"].Computed != true {
		t.Error("Expected tenant_list to be Computed")
	}
	if r.Schema["ignore_metadata"].Default != true {
		t.Error("Expected ignore_metadata to default to true")
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

// TestResourceBigiqAs3DiffSuppressUnknownVariable covers the special-cased
// "known after apply" placeholder branch of the as3_json DiffSuppressFunc.
func TestResourceBigiqAs3DiffSuppressUnknownVariable(t *testing.T) {
	r := resourceBigiqAs3()
	diffSuppress := r.Schema["as3_json"].DiffSuppressFunc
	d := bigiqAs3ResourceData(t, baseBigiqAs3Raw("http://192.0.2.1", testBigiqAs3Json))

	old := `{"class":"AS3"}`
	assert.False(t, diffSuppress("as3_json", old, unknownVariableValue, d), "expected placeholder diff to NOT be suppressed")
}

// TestResourceBigiqAs3DiffSuppressIdenticalJSON covers the fast-path where
// old and new JSON are already equal.
func TestResourceBigiqAs3DiffSuppressIdenticalJSON(t *testing.T) {
	r := resourceBigiqAs3()
	diffSuppress := r.Schema["as3_json"].DiffSuppressFunc
	d := bigiqAs3ResourceData(t, baseBigiqAs3Raw("http://192.0.2.1", testBigiqAs3Json))

	assert.True(t, diffSuppress("as3_json", testBigiqAs3Json, testBigiqAs3Json, d))
}

// TestResourceBigiqAs3DiffSuppressIgnoreMetadataTrue covers the branch where
// only metadata fields differ and ignore_metadata=true suppresses the diff.
func TestResourceBigiqAs3DiffSuppressIgnoreMetadataTrue(t *testing.T) {
	r := resourceBigiqAs3()
	diffSuppress := r.Schema["as3_json"].DiffSuppressFunc
	d := bigiqAs3ResourceData(t, baseBigiqAs3Raw("http://192.0.2.1", testBigiqAs3Json))

	old := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.50.0","id":"autogen","updateMode":"selective","label":"lbl","remark":"rmk","MyTenant":{"class":"Tenant"}}}`
	new := `{"class":"AS3","declaration":{"class":"ADC","MyTenant":{"class":"Tenant"}}}`

	assert.True(t, diffSuppress("as3_json", old, new, d))
}

// TestResourceBigiqAs3DiffSuppressIgnoreMetadataFalse covers the strict-mode
// branch where ignore_metadata=false surfaces a real diff.
func TestResourceBigiqAs3DiffSuppressIgnoreMetadataFalse(t *testing.T) {
	r := resourceBigiqAs3()
	diffSuppress := r.Schema["as3_json"].DiffSuppressFunc
	raw := baseBigiqAs3Raw("http://192.0.2.1", testBigiqAs3Json)
	raw["ignore_metadata"] = false
	d := bigiqAs3ResourceData(t, raw)

	old := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","MyTenant":{"class":"Tenant","app":{"class":"Application"}}}}`
	new := `{"class":"AS3","declaration":{"class":"ADC","schemaVersion":"3.0.0","MyTenant":{"class":"Tenant","app":{"class":"Application","pool":{"class":"Pool"}}}}}`

	assert.False(t, diffSuppress("as3_json", old, new, d))
}

// -----------------------------------------------------------------------
// connectBigIq (shared with resource_bigiq_regkey_license_manage.go) is
// exercised indirectly through Create/Read/Update/Delete below.
// -----------------------------------------------------------------------

// registerSelfIPOK registers a handler for the /mgmt/tm/net/self endpoint hit
// by connectBigIq -> Client() -> ValidateConnection() so that the connection
// step succeeds and the resource's own CRUD logic can be exercised.
func registerSelfIPOK() {
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})
}

// bigiqAs3TestHandlers wires up the standard mock endpoints used by
// Create/Read/Update: the AS3 declare POST endpoint at
// /mgmt/shared/appsvcs/declare and the tenant-scoped GET endpoint at
// /mgmt/shared/appsvcs/declare/<tenant>. It also registers a passing
// /mgmt/tm/net/self handler so connectBigIq succeeds.
type bigiqAs3TestHandlers struct {
	postFunc func(w http.ResponseWriter, r *http.Request)
	getFunc  func(w http.ResponseWriter, r *http.Request)
}

func registerBigiqAs3Handlers(h bigiqAs3TestHandlers) {
	registerSelfIPOK()
	mux.HandleFunc("/mgmt/shared/appsvcs/declare", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && h.postFunc != nil {
			h.postFunc(w, r)
			return
		}
	})
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/MyTenant", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && h.getFunc != nil {
			h.getFunc(w, r)
			return
		}
	})
}

// TestResourceBigiqAs3CreateSuccess covers the Create happy path: PostAs3Bigiq
// succeeds with no error, and the resource ID is composed from target+tenant.
func TestResourceBigiqAs3CreateSuccess(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()

	var sawPOST bool
	registerBigiqAs3Handlers(bigiqAs3TestHandlers{
		postFunc: func(w http.ResponseWriter, r *http.Request) {
			sawPOST = true
			_, _ = fmt.Fprint(w, `{"code":200,"results":[{"code":200,"tenant":"MyTenant"}]}`)
		},
		getFunc: func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"target":{"address":"192.0.2.50"},"MyTenant":{"class":"Tenant"}}`)
		},
	})

	d := bigiqAs3ResourceData(t, baseBigiqAs3Raw(server.URL, testBigiqAs3Json))

	diags := resourceBigiqAs3Create(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPOST, "expected PostAs3Bigiq to issue a POST request")
	assert.NotEmpty(t, d.Id())
	assert.Contains(t, d.Id(), "MyTenant")
}

// TestResourceBigiqAs3CreateTotalFailure covers the branch where
// PostAs3Bigiq returns an error and no tenants succeeded (successfulTenants
// == ""), so Create returns an error diagnostic and never sets the ID.
func TestResourceBigiqAs3CreateTotalFailure(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()

	registerBigiqAs3Handlers(bigiqAs3TestHandlers{
		postFunc: func(w http.ResponseWriter, r *http.Request) {
			// taskList.Code != 200/0 and no result has Code == 200 -> success_count == 0
			_, _ = fmt.Fprint(w, `{"code":400,"results":[{"code":400,"message":"failure","tenant":"MyTenant"}]}`)
		},
	})

	d := bigiqAs3ResourceData(t, baseBigiqAs3Raw(server.URL, testBigiqAs3Json))

	diags := resourceBigiqAs3Create(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error from Create when all tenants fail")
	assert.Empty(t, d.Id(), "resource ID should remain unset when Create fails entirely")
}

// TestResourceBigiqAs3CreatePartialFailure covers the branch where
// PostAs3Bigiq returns a "Partial Success" error along with a non-empty
// successfulTenants list; Create should still set tenant_list and proceed to
// build an ID (via the follow-up Read).
func TestResourceBigiqAs3CreatePartialFailure(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()

	registerBigiqAs3Handlers(bigiqAs3TestHandlers{
		postFunc: func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"code":400,"results":[{"code":400,"message":"failure","tenant":"OtherTenant"},{"code":200,"tenant":"MyTenant"}]}`)
		},
		getFunc: func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"target":{"address":"192.0.2.50"},"MyTenant":{"class":"Tenant"}}`)
		},
	})

	d := bigiqAs3ResourceData(t, baseBigiqAs3Raw(server.URL, testBigiqAs3JsonTwoTenants))

	diags := resourceBigiqAs3Create(context.Background(), d, nil)

	require.False(t, diags.HasError(), "partial success should not produce an error diagnostic")
	assert.Equal(t, "MyTenant", d.Get("tenant_list").(string))
}

// TestResourceBigiqAs3CreateConnectError covers the branch where
// connectBigIq itself fails (e.g. ValidateConnection against an
// unreachable/erroring BIG-IQ).
func TestResourceBigiqAs3CreateConnectError(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	d := bigiqAs3ResourceData(t, baseBigiqAs3Raw(server.URL, testBigiqAs3Json))

	diags := resourceBigiqAs3Create(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error from Create when connectBigIq fails")
}

// TestResourceBigiqAs3ReadDifferentTenant covers the branch of Read where the
// tenant embedded in the resource ID differs from the current tenant_list
// value, exercising the GetAs3Bigiq(targetRef, tenant_list) call path.
func TestResourceBigiqAs3ReadDifferentTenant(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()
	registerSelfIPOK()

	var sawGET bool
	mux.HandleFunc("/mgmt/shared/appsvcs/declare/OtherTenant", func(w http.ResponseWriter, r *http.Request) {
		sawGET = true
		assert.Equal(t, "GET", r.Method)
		_, _ = fmt.Fprint(w, `{"target":{"address":"192.0.2.50"},"OtherTenant":{"class":"Tenant"}}`)
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "OtherTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Read(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawGET, "expected GetAs3Bigiq to be called with the tenant_list value")
}

// TestResourceBigiqAs3ReadSameTenant covers the branch of Read where the
// tenant embedded in the resource ID matches tenant_list.
func TestResourceBigiqAs3ReadSameTenant(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()

	var sawGET bool
	registerBigiqAs3Handlers(bigiqAs3TestHandlers{
		getFunc: func(w http.ResponseWriter, r *http.Request) {
			sawGET = true
			_, _ = fmt.Fprint(w, `{"target":{"address":"192.0.2.50"},"MyTenant":{"class":"Tenant"}}`)
		},
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "MyTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Read(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawGET)
}

// TestResourceBigiqAs3ReadGetError covers the branch where GetAs3Bigiq itself
// returns an error (e.g. the BIG-IQ endpoint returns a server error).
func TestResourceBigiqAs3ReadGetError(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()

	registerBigiqAs3Handlers(bigiqAs3TestHandlers{
		getFunc: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "internal error", http.StatusInternalServerError)
		},
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "MyTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Read(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error from Read when GetAs3Bigiq fails")
}

// TestResourceBigiqAs3ReadConnectError covers the branch where connectBigIq
// fails inside Read.
func TestResourceBigiqAs3ReadConnectError(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "MyTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Read(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error from Read when connectBigIq fails")
}

// TestResourceBigiqAs3UpdateSuccess covers the Update happy path:
// PostAs3Bigiq succeeds, followed by a successful Read.
func TestResourceBigiqAs3UpdateSuccess(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()

	var sawPOST bool
	registerBigiqAs3Handlers(bigiqAs3TestHandlers{
		postFunc: func(w http.ResponseWriter, r *http.Request) {
			sawPOST = true
			_, _ = fmt.Fprint(w, `{"code":200,"results":[{"code":200,"tenant":"MyTenant"}]}`)
		},
		getFunc: func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"target":{"address":"192.0.2.50"},"MyTenant":{"class":"Tenant"}}`)
		},
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "MyTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Update(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPOST, "expected PostAs3Bigiq to issue a POST request")
}

// TestResourceBigiqAs3UpdateTotalFailure covers the branch where
// PostAs3Bigiq fails entirely during Update.
func TestResourceBigiqAs3UpdateTotalFailure(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()

	registerBigiqAs3Handlers(bigiqAs3TestHandlers{
		postFunc: func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"code":400,"results":[{"code":400,"message":"failure","tenant":"MyTenant"}]}`)
		},
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "MyTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Update(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error from Update when all tenants fail")
}

// TestResourceBigiqAs3UpdatePartialFailure covers the branch where
// PostAs3Bigiq returns a "Partial Success" error during Update (one tenant
// succeeds, one fails), which should set tenant_list to the successful
// subset rather than returning an error.
func TestResourceBigiqAs3UpdatePartialFailure(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()

	registerBigiqAs3Handlers(bigiqAs3TestHandlers{
		postFunc: func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"code":400,"results":[{"code":400,"message":"failure","tenant":"OtherTenant"},{"code":200,"tenant":"MyTenant"}]}`)
		},
		getFunc: func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"target":{"address":"192.0.2.50"},"MyTenant":{"class":"Tenant"}}`)
		},
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3JsonTwoTenants)
	raw["tenant_list"] = "MyTenant,OtherTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Update(context.Background(), d, nil)

	require.False(t, diags.HasError(), "partial success should not produce an error diagnostic")
	assert.Equal(t, "MyTenant", d.Get("tenant_list").(string))
}

// TestResourceBigiqAs3UpdateRecomputesTenantList covers the branch where the
// tenant list computed from as3_json differs from the stored tenant_list
// value, so tenant_list is refreshed before the PostAs3Bigiq call.
func TestResourceBigiqAs3UpdateRecomputesTenantList(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()

	registerBigiqAs3Handlers(bigiqAs3TestHandlers{
		postFunc: func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"code":200,"results":[{"code":200,"tenant":"MyTenant"}]}`)
		},
		getFunc: func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"target":{"address":"192.0.2.50"},"MyTenant":{"class":"Tenant"}}`)
		},
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	// tenant_list in state deliberately differs from what GetTenantList will
	// compute from as3_json ("MyTenant"), exercising the refresh branch.
	raw["tenant_list"] = "StaleTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Update(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
}

// TestResourceBigiqAs3UpdateConnectError covers the branch where connectBigIq
// fails inside Update.
func TestResourceBigiqAs3UpdateConnectError(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "MyTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Update(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error from Update when connectBigIq fails")
}

// TestResourceBigiqAs3DeleteSuccess covers the Delete happy path:
// DeleteAs3Bigiq succeeds with no failed tenants, so the resource ID is
// cleared without a follow-up Read.
func TestResourceBigiqAs3DeleteSuccess(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()
	registerSelfIPOK()

	// DeleteAs3Bigiq (via tenantTrimToDelete + b.post) issues a POST to the
	// base declare endpoint, not a tenant sub-path or an actual HTTP DELETE.
	var sawPOST bool
	mux.HandleFunc("/mgmt/shared/appsvcs/declare", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			sawPOST = true
			w.WriteHeader(http.StatusOK)
			return
		}
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "MyTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Delete(context.Background(), d, nil)

	require.False(t, diags.HasError(), "unexpected error: %v", diags)
	assert.True(t, sawPOST, "expected DeleteAs3Bigiq to issue a POST request")
	assert.Empty(t, d.Id(), "resource ID should be cleared on successful delete")
}

// TestResourceBigiqAs3DeleteError covers the branch where DeleteAs3Bigiq
// itself returns an error (e.g. tenantTrimToDelete or the underlying POST
// fails).
func TestResourceBigiqAs3DeleteError(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()
	registerSelfIPOK()

	mux.HandleFunc("/mgmt/shared/appsvcs/declare", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			http.Error(w, `{"code":400,"message":"invalid declaration"}`, http.StatusBadRequest)
			return
		}
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "MyTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Delete(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error from Delete when DeleteAs3Bigiq fails")
	assert.Equal(t, "192.0.2.50_MyTenant", d.Id(), "resource ID should remain set when Delete fails")
}

// TestResourceBigiqAs3DeleteConnectError covers the branch where
// connectBigIq fails inside Delete.
func TestResourceBigiqAs3DeleteConnectError(t *testing.T) {
	withZeroBigiqAs3RetryDelay(t)
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["tenant_list"] = "MyTenant"
	d := bigiqAs3ResourceData(t, raw)
	d.SetId("192.0.2.50_MyTenant")

	diags := resourceBigiqAs3Delete(context.Background(), d, nil)

	require.True(t, diags.HasError(), "expected an error from Delete when connectBigIq fails")
}

// TestConnectBigIqTokenAuth exercises the bigiq_token_auth=true branch of
// connectBigIq, which routes through bigip.NewTokenSession instead of
// bigip.NewSession.
func TestConnectBigIqTokenAuth(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"token":{"token":"unit-test-token"},"timeout":{"timeout":1200}}`)
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})

	raw := baseBigiqAs3Raw(server.URL, testBigiqAs3Json)
	raw["bigiq_token_auth"] = true
	raw["bigiq_login_ref"] = "tmos"
	d := bigiqAs3ResourceData(t, raw)

	client, err := connectBigIq(d)

	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "unit-test-token", client.Token)
}

// TestConnectBigIqBasicAuth exercises the default (bigiq_token_auth=false)
// branch of connectBigIq, which routes through bigip.NewSession.
func TestConnectBigIqBasicAuth(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"items":[]}`)
	})

	d := bigiqAs3ResourceData(t, baseBigiqAs3Raw(server.URL, testBigiqAs3Json))

	client, err := connectBigIq(d)

	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Empty(t, client.Token)
}
