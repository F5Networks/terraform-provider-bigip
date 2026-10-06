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

func TestAccBigipAuthRadius_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_auth_radius_server" "test-radius-server" {
  name    = "test-radius-server"
  server  = "192.0.2.20"
  secret  = "supersecret"
  port    = 1812
  timeout = 3
}

resource "bigip_auth_radius" "test-radius" {
  servers      = [bigip_auth_radius_server.test-radius-server.name]
  service_type = "authenticate-only"
  retries      = 3
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_auth_radius_server.test-radius-server", "server", "192.0.2.20"),
					resource.TestCheckResourceAttr("bigip_auth_radius.test-radius", "servers.0", "test-radius-server"),
					resource.TestCheckResourceAttr("bigip_auth_radius.test-radius", "retries", "3"),
				),
			},
		},
	})
}

func TestAccBigipAuthRadiusServer_import(t *testing.T) {
	resourceName := "bigip_auth_radius_server.test-radius-server"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "bigip_auth_radius_server" "test-radius-server" {
  name   = "test-radius-server"
  server = "192.0.2.20"
  secret = "supersecret"
}
`,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           "test-radius-server",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
		},
	})
}
