/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"

	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestTranslateNodeState(t *testing.T) {
	cases := []struct {
		name           string
		state, session string
		wantState      string
		wantSession    string
	}{
		{"canonical enabled", "enabled", "", "user-up", "user-enabled"},
		{"canonical disabled", "disabled", "", "user-up", "user-disabled"},
		{"canonical forced_offline", "forced_offline", "", "user-down", "user-disabled"},
		{"canonical state ignores explicit session", "enabled", "user-disabled", "user-up", "user-enabled"},
		{"legacy user-up passes through", "user-up", "user-enabled", "user-up", "user-enabled"},
		{"legacy user-down passes through", "user-down", "user-disabled", "user-down", "user-disabled"},
		{"empty passes through", "", "", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotState, gotSession := translateNodeState(tc.state, tc.session)
			assert.Equal(t, tc.wantState, gotState, "state")
			assert.Equal(t, tc.wantSession, gotSession, "session")
		})
	}
}

func TestNodeStateForRead(t *testing.T) {
	cases := []struct {
		name                  string
		prior                 string
		apiState, apiSession  string
		wantState, wantSession string
	}{
		// Canonical-prior scenarios (new style)
		{"canonical prior, healthy", "enabled", "user-up", "monitor-enabled", "enabled", "user-enabled"},
		{"canonical prior, user-disabled", "enabled", "user-up", "user-disabled", "disabled", "user-disabled"},
		{"canonical prior, force-offlined", "enabled", "user-down", "user-disabled", "forced_offline", "user-disabled"},
		{"canonical prior, monitor reports down (still user-up)", "enabled", "down", "monitor-enabled", "enabled", "user-enabled"},
		// Legacy-prior scenarios (preserve current behavior)
		{"legacy prior user-up, monitor-enabled", "user-up", "user-up", "monitor-enabled", "user-up", "user-enabled"},
		{"legacy prior user-up, user-disabled", "user-up", "user-up", "user-disabled", "user-up", "user-disabled"},
		{"legacy prior user-down, force offline", "user-down", "user-down", "user-disabled", "user-down", "user-disabled"},
		// Empty prior (fresh import) defaults to canonical
		{"empty prior treated as canonical", "", "user-up", "monitor-enabled", "enabled", "user-enabled"},
		{"empty prior, force-offlined device", "", "user-down", "user-disabled", "forced_offline", "user-disabled"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotState, gotSession := nodeStateForRead(tc.prior, tc.apiState, tc.apiSession)
			assert.Equal(t, tc.wantState, gotState, "state")
			assert.Equal(t, tc.wantSession, gotSession, "session")
		})
	}
}

func testBigipLtmNodeInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_node" "test-node" {
  name       = "%s"
  address    = "10.10.10.10"
  invalidkey = "foo"
}
`, resourceName)
}

func TestAccBigipLtmNodeInvalid(t *testing.T) {
	resourceName := "/Common/test-node"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testBigipLtmNodeInvalid(resourceName),
				ExpectError: regexp.MustCompile("An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func testBigipLtmNodeCreate(resourceName string, url string, address string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_node" "test-node" {
  name    = "%s"
  address = "%s"
}
//provider "bigip" {
////alias = "unitbigip"
//address  = "%s"
//username = "xxxx"
//password = "xxxx"
//}
`, resourceName, address, url)
}

func TestAccBigipLtmNodeCreate(t *testing.T) {
	resourceName := "/Common/test-node"
	address := "10.10.10.10"
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		log.Println(" value of t  ")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{
}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		_, _ = fmt.Fprintf(w, `{"name":"%s","address":"%s"}`, resourceName, address)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~test-node", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","address":"%s","monitor":"/Common/icmp"}`, resourceName, address)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testBigipLtmNodeCreate(resourceName, server.URL, address),
			},
		},
	})
}

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func teardown() {
	server.Close()
}
