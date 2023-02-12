/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var TestPartition = "Common"

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.

var testAccProvider = Provider()

var testAccProviders = map[string]func() (*schema.Provider, error){
	"bigip": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func TestAccProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAcctPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if os.Getenv("BIGIP_TOKEN_VALUE") != "" || (os.Getenv("BIGIP_TOKEN_AUTH") != "" && os.Getenv("BIGIP_LOGIN_REF") != "") {
		return
	}
	if v := os.Getenv("BIGIP_TEST_PARTITION"); v != "" {
		TestPartition = v
	}
	for _, s := range [...]string{"BIGIP_HOST", "BIGIP_USER", "BIGIP_PASSWORD"} {
		if os.Getenv(s) == "" {
			t.Fatal("Either BIGIP_TOKEN_AUTH + BIGIP_LOGIN_REF or BIGIP_USER, BIGIP_PASSWORD and BIGIP_HOST are required for tests.")
			return
		}
	}
}
