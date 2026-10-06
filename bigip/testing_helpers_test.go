/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

// Shared unit-test infrastructure for bigip/*_unit_test.go files.
//
// This repo has settled on two complementary unit-testing styles (see
// CONTRIBUTING.md for the full writeup and a copy-pasteable template):
//
//  1. SDK-driven acceptance-style unit tests: use resource.Test with
//     IsUnitTest: true, the package-level setup()/teardown()/mux/server
//     singletons (defined in resource_bigip_ltm_node_unit_test.go) and
//     testAcctUnitPreCheck (defined in provider_test.go) to point the
//     provider at a local httptest server instead of a real device.
//     Terraform config is exercised end-to-end through the real provider.
//
//  2. Direct CRUD-function unit tests: call a resource's/data source's
//     Create/Read/Update/Delete function directly against a *schema.ResourceData
//     built with NewTestResourceData (or the data-source-only
//     newDatasourceTestResourceData) and a *bigip.BigIP client built with
//     NewUnitTestClient (or the data-source-only newDatasourceTestClient),
//     talking to a local httptest.Server owned by the test itself (not the
//     shared singleton) so slow/blocking tests can run independently.
//
// The helpers below exist to avoid re-deriving these two patterns (and their
// supporting boilerplate) from scratch in every new *_unit_test.go file.
// Prefer them over hand-rolled equivalents in new tests; existing files that
// predate this helper are not required to be migrated.
//
// # 404 handling in Read tests
//
// A number of resources/data sources (e.g. the GTM ones) have a "not found"
// branch in their Read function that checks for a nil result from a getter
// like GetGTMDatacenter and, on nil, clears the resource ID (or, for data
// sources, returns a "not found" diagnostic) instead of propagating an
// error. When writing a unit test for a 404 response, do not expect that
// branch to execute: go-bigip's getForEntity (vendor/github.com/f5devcentral/
// go-bigip/bigip.go) always returns a non-nil error when APICall fails,
// including on a 404 -- the "not found -> nil result, nil error" path it
// appears to support is unreachable, because the caller's `if err != nil`
// check fires first regardless of the underlying HTTP status. Concretely,
// GetGTMDatacenter (and its siblings) never actually return (nil, nil) via
// the real client for a 404; they return (nil, err). So a unit test that
// mocks a 404 response should assert diags.HasError() (and, for resources,
// that the ID is unchanged), not assert that the ID was cleared. See
// TestUnitGtmDatacenterReadNotFound / TestDataSourceBigipGtmDatacenterReadNotFound
// for worked examples, including a comment that explains this at the
// specific call site.

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NewUnitTestClient returns a *bigip.BigIP configured to talk to the given
// httptest server URL, with retries/timeouts set so the go-bigip APICall
// retry loop executes exactly once per request (rather than retrying for
// real against a mock that isn't going to change its answer). Use this for
// "direct CRUD-function" style unit tests (see package doc above); SDK-driven
// tests should use testAcctUnitPreCheck instead, which configures the
// provider (and therefore its client) via environment variables.
func NewUnitTestClient(serverURL string) *bigip.BigIP {
	return bigip.NewSession(&bigip.Config{
		Address:           serverURL,
		Username:          "admin",
		Password:          "admin",
		CertVerifyDisable: true,
		ConfigOptions: &bigip.ConfigOptions{
			APICallRetries: 1,
			APICallTimeout: 5 * time.Second,
		},
	})
}

// NewTestResourceData builds a *schema.ResourceData for the given resource or
// data source using schema.TestResourceDataRaw. If id is non-empty it is set
// as the Terraform resource ID afterward, mimicking state as it would exist
// after a Create or Import (schema.TestResourceDataRaw itself has no way to
// set an ID). Pass "" for data sources and for resource Create tests, where
// no ID exists yet.
func NewTestResourceData(t *testing.T, r *schema.Resource, raw map[string]interface{}, id string) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, r.Schema, raw)
	if id != "" {
		d.SetId(id)
	}
	return d
}

// NewTestResourceDataWithRawConfig builds a *schema.ResourceData for r whose
// GetRawConfig() actually returns rawConfig, unlike schema.TestResourceDataRaw
// (used by NewTestResourceData), which always leaves GetRawConfig() null.
// This is needed by any Read/Update/getConfig-style code path that reads
// exclusively via d.GetRawConfig() -- rather than the ordinary d.Get/d.GetOk,
// which can't distinguish "the user set this in their .tf" from "Read
// backfilled this into state from an inherited/computed device value" -- to
// decide whether to include an Optional+Computed field in an outgoing
// Create/Update payload. See resource_bigip_ltm_profile_ssl_client.go's
// tm_options handling (SFDC #01262589) and resource_bigip_ltm_policy.go's
// dataToPolicy rule parsing for the two motivating cases; both
// clientSslResourceDataWithRawConfig and policyResourceDataWithRawConfig
// used to duplicate this exact construction before being consolidated here.
func NewTestResourceDataWithRawConfig(t *testing.T, r *schema.Resource, raw map[string]interface{}, rawConfig cty.Value) *schema.ResourceData {
	t.Helper()
	im := schema.InternalMap(r.Schema)
	rc := terraform.NewResourceConfigRaw(raw)
	diff, err := im.Diff(context.Background(), nil, rc, nil, nil, true)
	require.NoError(t, err)
	if diff == nil {
		diff = new(terraform.InstanceDiff)
	}
	diff.RawConfig = rawConfig
	d, err := im.Data(nil, diff)
	require.NoError(t, err)
	return d
}

// MangleFullPath converts a BIG-IP full path such as "/Common/my-object" into
// the mangled form ("~Common~my-object") iControl REST uses in URLs, e.g.
// "/mgmt/tm/ltm/pool/" + MangleFullPath(name). This mirrors (but does not
// call) go-bigip's unexported iControlPath, since the mangling rule itself
// (replace "/" with "~") is a stable part of the iControl REST API, not an
// implementation detail of the client. Use this instead of hand-writing
// mangled path literals in mock handler registrations.
func MangleFullPath(fullPath string) string {
	mangled := make([]byte, 0, len(fullPath))
	for i := 0; i < len(fullPath); i++ {
		if fullPath[i] == '/' {
			mangled = append(mangled, '~')
		} else {
			mangled = append(mangled, fullPath[i])
		}
	}
	return string(mangled)
}

// ExpectUnexpectedArgumentError is the regexp Terraform's config parser
// produces when a resource block sets an argument the schema doesn't define.
// Shared by the "invalidkey" style schema-rejection tests (e.g.
// TestUnitBigipLtmNodeInvalid): a Config containing an undefined
// "invalidkey = ..." argument should fail a resource.Test step with this
// ExpectError.
var ExpectUnexpectedArgumentError = regexp.MustCompile(`An argument named "invalidkey" is not expected here`)

// AssertRequestMethod asserts that r was made with the given HTTP method
// (e.g. AssertRequestMethod(t, r, http.MethodPost)). Intended for use inside
// mux.HandleFunc mock handlers in direct CRUD-function unit tests.
func AssertRequestMethod(t *testing.T, r *http.Request, method string) {
	t.Helper()
	assert.Equal(t, method, r.Method, "expected method %q, got %q", method, r.Method)
}

// AssertJSONContentType asserts that r declares a JSON request body
// (Content-Type: application/json), as go-bigip always does for requests
// with a body. Intended for use inside mux.HandleFunc mock handlers.
func AssertJSONContentType(t *testing.T, r *http.Request) {
	t.Helper()
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
}

// NewUnitTestServer collapses the "direct CRUD-function" unit test
// boilerplate of creating a *http.ServeMux, wrapping it in an
// httptest.Server, registering the server's Close via t.Cleanup, and
// building a *bigip.BigIP client pointed at it. It returns the mux (so the
// caller can still register handlers on it, since that part is inherently
// test-specific) and the client. For example:
//
//	mux, client := NewUnitTestServer(t)
//	mux.HandleFunc("/mgmt/tm/gtm/datacenter", func(w http.ResponseWriter, r *http.Request) {
//		_, _ = fmt.Fprint(w, `{"name":"test-dc"}`)
//	})
//	diags := resourceBigipGtmDatacenterCreate(ctx, d, client)
//
// Using t.Cleanup instead of `defer server.Close()` means this only saves a
// line when a test has no other defer to combine it with, but it also means
// callers can't forget the Close(). Existing tests that already spell out
// mux/server/client construction inline are not required to migrate.
func NewUnitTestServer(t *testing.T) (*http.ServeMux, *bigip.BigIP) {
	t.Helper()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return mux, NewUnitTestClient(server.URL)
}

// IsNotFoundError reports whether err looks like the "object does not
// exist" error go-bigip's getForEntity returns for a 404 response (see the
// package doc above: it returns a non-nil error rather than (nil, nil) for
// missing objects, and the BIG-IP-supplied message text varies by endpoint,
// e.g. "01020036:3: The requested VLAN (...) was not found." or the net
// route endpoint's own "route not found: <path>"). CheckDestroy functions
// that call a client Get*/GetXyz method should treat this as "successfully
// destroyed" instead of a fatal error -- the object==nil branch such
// functions may also check is unreachable in practice for the same reason.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}
