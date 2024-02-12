/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var rewPorfile = fmt.Sprintf("/%s/%s", TestPartition, "tf_profile")
var rule1 = "tf_rule"
var TestProfileResource = `
resource "bigip_ltm_profile_rewrite" "tftest" {
  name = "` + rewPorfile + `"
  defaults_from = "/Common/rewrite"
  rewrite_mode = "uri-translation"
  
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

  cookie_rules {
    rule_name = "cookie1"
    client_domain = "wrong.com"
    client_path   = "/this/"
    server_domain = "wrong.com"
    server_path   = "/this/"
  }
}

resource "bigip_ltm_profile_rewrite_uri_rules" "tftestrule1" {
  profile_name = "` + rewPorfile + `"
  rule_name = "` + rule1 + `""  
  rule_type = "request"

  client {
    host = "www.foo.com"
    scheme = "https"
  }

  server {
    host = "www.bar.com"
    path = "/this/"
    scheme = "https"
    port = "8888"
  }
}
`

func TestAccLtmRewriteProfileUriRulesCreateOnBigip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckLtmRewriteProfileDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestProfileResource,
				Check: resource.ComposeTestCheckFunc(
					testLtmRewriteProfileExists(rewPorfile, true),
					testLtmRewriteProfileUriRuleExists(rewPorfile, rule1, true),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.tftest", "rewrite_mode", "uri-translation"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.tftest", "request.0.insert_xfwd_for", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.tftest", "request.0.insert_xfwd_host", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.tftest", "request.0.insert_xfwd_protocol", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.tftest", "cookie_rules.0.rule_name", "cookie1"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.tftest", "cookie_rules.0.client_domain", "wrong.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.tftest", "cookie_rules.0.client_path", "/this/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.tftest", "cookie_rules.0.server_domain", "right.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite.tftest", "cookie_rules.0.server_path", "/that/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite_uri_rules.tftestrule1", "rule_name", rule1),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite_uri_rules.tftestrule1", "rule_type", "request"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite_uri_rules.tftestrule1", "client.0.host", "www.foo.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite_uri_rules.tftestrule1", "client.0.scheme", "https"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite_uri_rules.tftestrule1", "server.0.host", "www.bar.com"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite_uri_rules.tftestrule1", "server.0.path", "/this/"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite_uri_rules.tftestrule1", "server.0.scheme", "https"),
					resource.TestCheckResourceAttr("bigip_ltm_profile_rewrite_uri_rules.tftestrule1", "server.0.port", "8888"),
				),
			},
		},
	})
}

func testLtmRewriteProfileUriRuleExists(profile string, uri string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.GetRewriteProfileUriRule(profile, uri)
		if err != nil {
			return err
		}
		if exists && p != nil {
			return fmt.Errorf("rewrite profile uri rule %s was not created", uri)
		}
		if !exists && p != nil {
			return fmt.Errorf("rewrite profile uri rule %s still exists", uri)
		}
		return nil
	}
}
