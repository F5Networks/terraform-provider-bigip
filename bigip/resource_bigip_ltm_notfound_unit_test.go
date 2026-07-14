/*
Copyright 2026 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Regression test for objects deleted out-of-band (outside Terraform) on the
// BIG-IP: the Read functions must treat the iControl 404 ("... was not found")
// as "resource gone -> remove from state" so the next plan can recreate it,
// instead of returning a fatal error that aborts refresh/plan entirely.
// https://github.com/F5Networks/terraform-provider-bigip/issues/1169
func TestAccBigipLtmNodeDeletedOutOfBand(t *testing.T) {
	resourceName := "/Common/test-node"
	address := "10.10.10.10"
	deleted := false
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"%s","address":"%s"}`, resourceName, address)
	})
	mux.HandleFunc("/mgmt/tm/ltm/node/~Common~test-node", func(w http.ResponseWriter, r *http.Request) {
		if deleted {
			// iControl-style 404 for a manually deleted object
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, `{"code":404,"message":"01020036:3: The requested Node (/Common/test-node) was not found.","errorStack":[],"apiError":3}`)
			return
		}
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
				Check: resource.TestCheckResourceAttr(
					"bigip_ltm_node.test-node", "name", resourceName),
			},
			{
				// Simulate the object being deleted on the box outside
				// Terraform: refresh must SUCCEED (not error out) and the
				// follow-up plan must be non-empty (recreate).
				PreConfig:          func() { deleted = true },
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Same scenario for an iRule, whose client call wraps the iControl error, to
// cover the wrapped-error path through the not-found guard.
func TestAccBigipLtmIruleDeletedOutOfBand(t *testing.T) {
	deleted := false
	setup()
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/rule", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"name":"test-rule","fullPath":"/Common/test-rule","apiAnonymous":"when HTTP_REQUEST {}"}`)
	})
	mux.HandleFunc("/mgmt/tm/ltm/rule/~Common~test-rule", func(w http.ResponseWriter, r *http.Request) {
		if deleted {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, `{"code":404,"message":"01020036:3: The requested iRule (/Common/test-rule) was not found.","errorStack":[],"apiError":3}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"name":"test-rule","fullPath":"/Common/test-rule","apiAnonymous":"when HTTP_REQUEST {}"}`)
	})
	defer teardown()
	config := `
resource "bigip_ltm_irule" "test-rule" {
  name  = "/Common/test-rule"
  irule = "when HTTP_REQUEST {}"
}
`
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck:   func() { testAcctUnitPreCheck(t, server.URL) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.TestCheckResourceAttr(
					"bigip_ltm_irule.test-rule", "name", "/Common/test-rule"),
			},
			{
				PreConfig:          func() { deleted = true },
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
