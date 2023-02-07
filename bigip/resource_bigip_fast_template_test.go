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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var folder2, _ = os.Getwd()
var TEST_TEMPLATE = `foo_template`
var TEST_FAST_TEMPLATE = `
resource "bigip_fast_template" "foo-template" {
  name		= "` + TEST_TEMPLATE + `"
  source   = "${"` + folder2 + `/../examples/fast/foo_template.zip"}"
  md5_hash = "89011331d11ac8bac2a1ad3235f38c80"
}
`

func TestAccFastTemplateCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastTemplateDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FAST_TEMPLATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckFastTemplateExists(TEST_TEMPLATE, true),
					resource.TestCheckResourceAttr("bigip_fast_template.foo-template", "name", TEST_TEMPLATE),
					resource.TestCheckResourceAttr("bigip_fast_template.foo-template", "md5_hash", "89011331d11ac8bac2a1ad3235f38c80"),
				),
			},
		},
	})
}

func testCheckFastTemplateExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetTemplateSet(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("Fast Template set %s was not created.", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("Fast Template set  %s still exists.", name)
		}
		return nil
	}
}

func testCheckFastTemplateDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_fast_template" {
			continue
		}
		name := rs.Primary.ID
		template, err := client.GetTemplateSet(name)
		if err != nil {
			return err
		}
		if template != nil {
			return fmt.Errorf("Fast Template set  %s not destroyed.", name)
		}
	}
	return nil
}
