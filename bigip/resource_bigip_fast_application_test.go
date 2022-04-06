/*
Copyright 2021 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"testing"

	"os"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var folder3, _ = os.Getwd()
var template = "examples/simple_http"
var tenant = "sample_tenant"
var app = "sample_app"
var TestFastResource = `
resource "bigip_fast_application"  "foo-app" {
     template = "` + template + `"
     fast_json = "${file("` + folder3 + `/../examples/fast/new_fast_app.json")}"
}
`

func TestAccFastAppCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastAppDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestFastResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckFastAppExists("sample_app", "sample_tenant", true),
					resource.TestCheckResourceAttr("bigip_fast_application.foo-app", "application", app),
					resource.TestCheckResourceAttr("bigip_fast_application.foo-app", "tenant", tenant),
					resource.TestCheckResourceAttr("bigip_fast_application.foo-app", "template", template),
				),
			},
		},
	})
}

func testCheckFastAppExists(app, tenant string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetFastApp(tenant, app)
		if err != nil {
			return err
		}
		if exists && p == "" {
			return fmt.Errorf("fast application %s was not created", app)
		}
		if !exists && p != "" {
			return fmt.Errorf("fast application  %s still exists", app)
		}
		return nil
	}
}

func testCheckFastAppDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_fast_application" {
			continue
		}
		name := rs.Primary.ID
		template, err := client.GetFastApp(tenant, name)
		if err != nil {
			return err
		}
		if template != "" {
			return fmt.Errorf("Fast Application  %s not destroyed.", name)
		}
	}
	return nil
}
