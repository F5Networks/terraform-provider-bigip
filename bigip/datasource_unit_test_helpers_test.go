/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"testing"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// newDatasourceTestClient returns a *bigip.BigIP configured to talk to the
// supplied httptest server, with retries/timeouts set so the go-bigip
// APICall retry loop actually executes exactly once per request. Shared by
// the per-data-source unit tests (e.g. datasource_bigip_ltm_pool,
// datasource_bigip_ssl_certificate, datasource_bigip_ltm_irule) so each
// doesn't need its own copy of this client construction.
func newDatasourceTestClient(serverURL string) *bigip.BigIP {
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

// newDatasourceTestResourceData builds a *schema.ResourceData for the given
// data source resource using schema.TestResourceDataRaw. Shared by the
// per-data-source unit tests so each doesn't need its own thin wrapper
// around TestResourceDataRaw.
func newDatasourceTestResourceData(t *testing.T, r *schema.Resource, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, r.Schema, raw)
}
