/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"strings"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLtmRewriteProfileCreateOnBigipTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckLtmRewriteProfileDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getLtmRewritePortalProfileConfig(),
				Check: resource.ComposeTestCheckFunc(
					testLtmRewriteProfileExists("/Common/tf_profile-tc1", true),
					testLtmRewriteProfileExists("/Common/tf_profile-tc11", false),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "rewrite_mode", "portal"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cache_type", "cache-img-css-js"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "ca_file", "/Common/ca-bundle.crt"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "split_tunneling", "false"),
				),
			},
		},
	})
}
func TestAccLtmRewriteProfileCreateOnBigipTC2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckLtmRewriteProfileDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getLtmRewriteUriRewriteProfileConfig(),
				Check: resource.ComposeTestCheckFunc(
					testLtmRewriteProfileExists("/Common/tf_profile_translate", true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "rewrite_mode", "uri-translation"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "request.0.insert_xfwd_for", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "request.0.insert_xfwd_host", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "request.0.insert_xfwd_protocol", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.rule_name", "cookie1"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.client_domain", "wrong.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.client_path", "/this/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.server_domain", "right.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.server_path", "/that/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.rule_name", "cookie2"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.client_domain", "incorrect.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.client_path", "/this/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.server_domain", "absolute.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.server_path", "/that/"),
				),
			},
		},
	})
}

func TestAccLtmRewriteProfileCreateOnBigipTC3(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckLtmRewriteProfileDestroyed,
		Steps: []resource.TestStep{
			{
				Config: getLtmRewriteUriRewriteProfileConfig(),
				Check: resource.ComposeTestCheckFunc(
					testLtmRewriteProfileExists("/Common/tf_profile_translate", true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "rewrite_mode", "uri-translation"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "request.0.insert_xfwd_for", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "request.0.insert_xfwd_host", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "request.0.insert_xfwd_protocol", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.rule_name", "cookie1"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.client_domain", "wrong.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.client_path", "/this/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.server_domain", "right.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.server_path", "/that/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.rule_name", "cookie2"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.client_domain", "incorrect.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.client_path", "/this/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.server_domain", "absolute.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.1.server_path", "/that/"),
				),
			},
			{
				Config: getLtmRewriteUriRewriteProfileConfigChanged(),
				Check: resource.ComposeTestCheckFunc(
					testLtmRewriteProfileExists("/Common/tf_profile_translate", true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "request.0.insert_xfwd_for", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "request.0.insert_xfwd_host", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "request.0.insert_xfwd_protocol", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.rule_name", "cookie1"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.client_domain", "totallywrong.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.client_path", "/these/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.server_domain", "totallyright.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.test-profile", "cookie_rules.0.server_path", "/those/"),
				),
			},
		},
	})
}

func getLtmRewritePortalProfileConfig() string {
	return fmt.Sprintf(`resource "bigip_ltm_profile_rewrite" "test-profile" {
		name            = "/Common/%v"
		defaults_from   = "/Common/rewrite"
		rewrite_mode    = "portal"
		cache_type      = "cache-img-css-js"
		ca_file         = "/Common/ca-bundle.crt"
		signing_cert    = "/Common/default.crt"
		signing_key     = "/Common/default.key"
		split_tunneling = "false"}`, "tf_profile-tc1")
}
func getLtmRewriteUriRewriteProfileConfig() string {
	return fmt.Sprintf(`resource "bigip_ltm_profile_rewrite" "test-profile2" {
		name          = "%v"
		defaults_from = "/Common/rewrite"
		rewrite_mode  = "uri-translation"
	  
		request {
		  insert_xfwd_for      = "enabled"
		  insert_xfwd_host     = "disabled"
		  insert_xfwd_protocol = "enabled"
		}
		cookie_rules {
		  rule_name     = "cookie1"
		  client_domain = "wrong.com"
		  client_path   = "/this/"
		  server_domain = "right.com"
		  server_path   = "/that/"
		}
		cookie_rules {
		  rule_name     = "cookie2"
		  client_domain = "incorrect.com"
		  client_path   = "/this/"
		  server_domain = "absolute.com"
		  server_path   = "/that/"
		}
	  }`, "/Common/tf_profile_translate")
}
func getLtmRewriteUriRewriteProfileConfigChanged() string {
	return fmt.Sprintf(`resource "bigip_ltm_profile_rewrite" "test-profile2" {
	  name = "%v"

	  request {
		insert_xfwd_for = "disabled"
		insert_xfwd_host = "enabled"
		insert_xfwd_protocol = "disabled"
	  }
	 
	  cookie_rules {
		rule_name = "cookie1"
		client_domain = "totallywrong.com"
		client_path   = "/these/"
		server_domain = "totallyright.com"
		server_path   = "/those/"
	  }
	}`, "/Common/tf_profile_translate")
}
func testLtmRewriteProfileExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetRewriteProfile(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
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
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return err
		}
		if p != nil {
			return fmt.Errorf("Ltm Rewrite Profile  %s not destroyed.", name)
		}
	}
	return nil
}
