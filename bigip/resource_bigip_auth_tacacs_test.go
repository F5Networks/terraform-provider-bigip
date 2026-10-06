/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigipAuthTacacs_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_auth_tacacs" "test-tacacs" {
  servers    = ["192.0.2.30"]
  secret     = "tacacssecret"
  service    = "ppp"
  protocol   = "ip"
  encryption = "enabled"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_auth_tacacs.test-tacacs", "servers.0", "192.0.2.30"),
					resource.TestCheckResourceAttr("bigip_auth_tacacs.test-tacacs", "service", "ppp"),
					resource.TestCheckResourceAttr("bigip_auth_tacacs.test-tacacs", "id", "/Common/system-auth"),
				),
			},
		},
	})
}

func TestAccBigipAuthTacacs_import(t *testing.T) {
	resourceName := "bigip_auth_tacacs.test-tacacs"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_auth_tacacs" "test-tacacs" {
  servers = ["192.0.2.30"]
  secret  = "tacacssecret"
}
`,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           "/Common/system-auth",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
		},
	})
}
