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

func TestAccBigipAuthLdap_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_auth_ldap" "test-ldap" {
  servers        = ["ldap1.example.com"]
  port           = 389
  bind_dn        = "cn=admin,dc=example,dc=com"
  search_base_dn = "dc=example,dc=com"
  ssl            = "disabled"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_auth_ldap.test-ldap", "servers.0", "ldap1.example.com"),
					resource.TestCheckResourceAttr("bigip_auth_ldap.test-ldap", "port", "389"),
					resource.TestCheckResourceAttr("bigip_auth_ldap.test-ldap", "id", "/Common/system-auth"),
				),
			},
		},
	})
}

func TestAccBigipAuthLdap_import(t *testing.T) {
	resourceName := "bigip_auth_ldap.test-ldap"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_auth_ldap" "test-ldap" {
  servers = ["ldap1.example.com"]
}
`,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           "/Common/system-auth",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bind_pw"},
			},
		},
	})
}
