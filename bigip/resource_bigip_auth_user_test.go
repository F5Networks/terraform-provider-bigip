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

func TestAccBigipAuthUser_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_auth_user" "test-user" {
  name        = "tf-acc-test-user"
  description = "Terraform acceptance test user"
  password    = "Testpassw0rd!"
  shell       = "tmsh"
  partition_access {
    partition = "Common"
    role      = "guest"
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_auth_user.test-user", "name", "tf-acc-test-user"),
					resource.TestCheckResourceAttr("bigip_auth_user.test-user", "shell", "tmsh"),
					resource.TestCheckResourceAttr("bigip_auth_user.test-user", "partition_access.0.partition", "Common"),
					resource.TestCheckResourceAttr("bigip_auth_user.test-user", "partition_access.0.role", "guest"),
				),
			},
		},
	})
}

func TestAccBigipAuthUser_import(t *testing.T) {
	resourceName := "bigip_auth_user.test-user"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_auth_user" "test-user" {
  name     = "tf-acc-test-user-import"
  password = "Testpassw0rd!"
  shell    = "tmsh"
  partition_access {
    partition = "Common"
    role      = "guest"
  }
}
`,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           "tf-acc-test-user-import",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
