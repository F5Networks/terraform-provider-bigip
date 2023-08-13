/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigipDeclarativeOnboardTCs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: loadFixtureString("../examples/bigip_onboard.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput("do_json", regexp.MustCompile("ecosyshyd-bigip02.com")),
				),
			},
			{
				Config: loadFixtureString("../examples/bigip_onboard_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput("do_json", regexp.MustCompile("ecosyshyd-bigip03.com")),
				),
			},
		},
	})
}
