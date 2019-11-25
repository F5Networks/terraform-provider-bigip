/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

var TEST_PARTITION = "Common"

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

var testProviders = map[string]terraform.ResourceProvider{
	"bigip": Provider(),
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"bigip": testAccProvider,
	}
	if v := os.Getenv("BIGIP_TEST_PARTITION"); v != "" {
		TEST_PARTITION = v
	}
}

func TestAccProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAcctPreCheck(t *testing.T) {
	if os.Getenv("BIGIP_TOKEN_AUTH") != "" && os.Getenv("BIGIP_LOGIN_REF") != "" {
		return
	}
	for _, s := range [...]string{"BIGIP_HOST", "BIGIP_USER", "BIGIP_PASSWORD"} {
		if os.Getenv(s) == "" {
			t.Fatal("Either BIGIP_TOKEN_AUTH + BIGIP_LOGIN_REF or BIGIP_USER, BIGIP_PASSWORD and BIGIP_HOST are required for tests.")
			return
		}
	}
}
