/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLtmRewriteProfileCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckLtmRewriteProfileDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getLtmRewriteProfileConfig(),
				Check: resource.ComposeTestCheckFunc(
					testLtmRewriteProfileExists("/Common/tf_profile", true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "rewrite_mode", "portal"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cache_type", "cache-img-css-js"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "ca_file", "/Common/ca-bundle.crt"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "split_tunneling", "false"),
				),
			},
		},
	})
}

func getLtmRewriteProfileConfig() string {
	return fmt.Sprintf(`
	resource "bigip_ltm_profile_rewrite" "test-profile" {
	  name = "%v"
	  defaults_from = "/Common/rewrite"
	  bypass_list = ["http://notouch.com"]
	  rewrite_list = ["http://some.com"]
	  rewrite_mode = "portal"
	  cache_type = "cache-img-css-js"
	  ca_file = "/Common/ca-bundle.crt"
	  crl_file = "none"
	  signing_cert = "/Common/default.crt"
	  signing_key = "/Common/default.key"
	  split_tunneling = "false"
	  
	  request {
		insert_xfwd_for = "enabled"
		insert_xfwd_host = "disabled"
		insert_xfwd_protocol = "enabled"
		rewrite_headers = "disabled"
	  }
	
	  response {
		rewrite_content = "enabled"
		rewrite_headers = "disabled"
	  }
	}
`, "/Common/tf_profile")
}

func testLtmRewriteProfileExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetRewriteProfile(name)
		if err != nil {
			return err
		}
		if exists && p != nil {
			return fmt.Errorf("rewrite profile %s was not created", name)
		}
		if !exists && p != nil {
			return fmt.Errorf("rewrite profile %s still exists", name)
		}
		return nil
	}
}

func testCheckLtmRewriteProfileDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_profile_rewrite" {
			continue
		}
		name := rs.Primary.ID
		p, err := client.GetRewriteProfile(name)
		if err != nil {
			return err
		}
		if p != nil {
			return fmt.Errorf("Ltm Rewrite Profile  %s not destroyed.", name)
		}
	}
	return nil
}
