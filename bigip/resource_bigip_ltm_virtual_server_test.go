/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TestVsName = fmt.Sprintf("/%s/test-vs", TEST_PARTITION)

var TestVsResource = TEST_IRULE_RESOURCE + `
resource "bigip_ltm_policy" "http_to_https_redirect" {
  name = "/Common/http_to_https_redirect"
  strategy = "first-match"
  requires = ["http"]
  published_copy = "Drafts/http_to_https_redirect"
  controls = ["forwarding"]
  rule  {
    name = "http_to_https_redirect_rule"
    action {
      redirect = true
      location = "tcl:https://[HTTP::host][HTTP::uri]"
      http_reply = true
    }
  }
}
resource "bigip_ltm_virtual_server" "test-vs" {
	name = "` + TestVsName + `"
	destination = "10.255.255.254"
	description = "VirtualServer-test"
	port = 9999
	mask = "255.255.255.255"
	source_address_translation = "automap"
	ip_protocol = "tcp"
	irules = [bigip_ltm_irule.test-rule.name]
	profiles = ["/Common/http"]
	client_profiles = ["/Common/tcp"]
	server_profiles = ["/Common/tcp-lan-optimized"]
	persistence_profiles = ["/Common/source_addr","/Common/hash"]
	default_persistence_profile = "/Common/hash"
	fallback_persistence_profile = "/Common/dest_addr"
    policies = [bigip_ltm_policy.http_to_https_redirect.name]
}
`
var TestVs6Name = fmt.Sprintf("/%s/test-vs6", TEST_PARTITION)

var TestVs6Resource = TEST_IRULE_RESOURCE + `
resource "bigip_ltm_virtual_server" "test-vs" {
	name = "` + TestVs6Name + `"
    destination = "fe80::11"
	description = "VirtualServer-test"
	port = 9999
	source_address_translation = "automap"
	ip_protocol = "tcp"
	irules = [bigip_ltm_irule.test-rule.name]
	profiles = ["/Common/http"]
	client_profiles = ["/Common/tcp"]
	server_profiles = ["/Common/tcp-lan-optimized"]
	persistence_profiles = ["/Common/source_addr", "/Common/hash"]
	default_persistence_profile = "/Common/hash"
	fallback_persistence_profile = "/Common/dest_addr"
}
`

func TestAccBigipLtmVS_CreateV4V6(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
			testCheckIRulesDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: TestVsResource,
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(TestVsName),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", TestVsName),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "10.255.255.254"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "9999"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source_address_translation", "automap"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "irules.0", TEST_IRULE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/http")),
						"/Common/http"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("client_profiles.%d", schema.HashString("/Common/tcp")),
						"/Common/tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("server_profiles.%d", schema.HashString("/Common/tcp-lan-optimized")),
						"/Common/tcp-lan-optimized"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("persistence_profiles.%d", schema.HashString("/Common/source_addr")),
						"/Common/source_addr"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "fallback_persistence_profile", "/Common/dest_addr"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
			testCheckIRulesDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: TestVs6Resource,
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(TestVs6Name),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", TestVs6Name),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "fe80::11"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "9999"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source", "::/0"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", "VirtualServer-test"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "source_address_translation", "automap"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "irules.0", TEST_IRULE_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "default_persistence_profile", "/Common/hash"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/http")),
						"/Common/http"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("client_profiles.%d", schema.HashString("/Common/tcp")),
						"/Common/tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("server_profiles.%d", schema.HashString("/Common/tcp-lan-optimized")),
						"/Common/tcp-lan-optimized"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("persistence_profiles.%d", schema.HashString("/Common/source_addr")),
						"/Common/source_addr"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("persistence_profiles.%d", schema.HashString("/Common/hash")),
						"/Common/hash"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "fallback_persistence_profile", "/Common/dest_addr"),
				),
			},
		},
	})
}

func TestAccBigipLtmVS_create_Defaultstate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testVSCreateDefaultstate("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")), "/Common/fastL4"),
				),
			},
		},
	})
}

func TestAccBigipLtmVS_Modify_stateDisabledtoEnabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testVSCreatestatedisabled("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
						"/Common/fastL4"),
				),
			},
			{
				Config: testVSCreatestateenabled("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
						"/Common/fastL4"),
				),
			},
		},
	})
}

func TestAccBigipLtmVS_Modify_stateEnabledtoDisabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testVSCreatestateenabled("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "enabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
						"/Common/fastL4"),
				),
			},
			{
				Config: testVSCreatestatedisabled("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "description", ""),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "state", "disabled"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "any"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
						fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
						"/Common/fastL4"),
				),
			},
		},
	})
}
func TestAccBigipLtmVS_Policyattach_detach(t *testing.T) {
	var rsName = "bigip_ltm_virtual_server.test_virtual_server_policyattch_detach"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testVSCreatePolicyAttach("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample"),
					resource.TestCheckResourceAttr(rsName, "name", "/Common/test-vs-sample"),
					resource.TestCheckResourceAttr(rsName, fmt.Sprintf("policies.%d", schema.HashString("/Common/test-policy")), "/Common/test-policy"),
				),
			},
			{
				Config: testVSCreatePolicyDettach("test-vs-sample"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-sample"),
					resource.TestCheckResourceAttr(rsName, "name", "/Common/test-vs-sample"),
				),
			},
		},
	})
}
func TestAccBigipLtmVS_Pooolattach_detatch(t *testing.T) {
	var poolName = "test-pool"
	var policyName = "test-policy"
	var vsName = "test-vs-sample"
	var partition = "Common"
	var vsFullname = fmt.Sprintf("/%s/%s", partition, vsName)
	var rsName = "bigip_ltm_virtual_server." + vsName
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testVSCreateAttach(poolName, policyName, vsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(vsName),
					resource.TestCheckResourceAttr(rsName, "name", vsFullname),
				),
			},
			{
				Config: testVSCreateDetatch(poolName, policyName, vsName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(vsName),
					resource.TestCheckResourceAttr(rsName, "name", vsFullname),
				),
			},
		},
	})
}

func TestAccBigipLtmVS_Vlan_EnabledDisabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckVSsDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: testVSCreatevlanEnabled("test-vs-vlan"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("test-vs-vlan"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs-vlan", "name", "/Common/test-vs-vlan"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs-vlan", "destination", "192.168.50.11"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs-vlan", "port", "80"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs-vlan", "vlans_enabled", "true"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs-vlan", fmt.Sprintf("vlans.%d", schema.HashString("/Common/external")), "/Common/external"),
				),
			},
		},
	})
}

func TestAccBigipLtmVS_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVSsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccBigipLtmVSImportConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("/Common/test-vs"),
				),
			},
			{
				ResourceName:      "bigip_ltm_virtual_server.test_vs_import",
				ImportStateId:     "/Common/test-vs",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testaccBigipLtmVSImportConfig() string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test_vs_import" {
  name        = "%s"
  destination = "192.168.11.11"
  port        = 80
}
`, "/Common/test-vs")
}

func testCheckVSExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		vs, err := client.GetVirtualServer(name)
		if err != nil {
			return err
		}
		if vs == nil {
			return fmt.Errorf("Virtual server %s does not exist.", name)
		}

		return nil
	}
}
func testVSCreatestatedisabled(vsName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test-vs" {
  name        = "/Common/%s"
  destination = "192.168.50.1"
  port        = 800
  mask        = "255.255.255.255"
  state       = "disabled"
}
`, vsName)
}

func testVSCreatestateenabled(vsName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test-vs" {
  name        = "/Common/%s"
  destination = "192.168.50.1"
  port        = 800
  mask        = "255.255.255.255"
  state       = "enabled"
}
`, vsName)
}

func testVSCreateDefaultstate(vsName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test-vs" {
  name        = "/Common/%s"
  destination = "192.168.50.1"
  port        = 800
  mask        = "255.255.255.255"
}
`, vsName)
}

func testVSCreatePolicyAttach(vsName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool" "mypool" {
  name                = "/Common/test-pool"
  allow_nat           = "yes"
  allow_snat          = "yes"
  load_balancing_mode = "round-robin"
}

resource "bigip_ltm_policy" "test-policy" {
  name     = "/Common/test-policy"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule6"
    action {
      forward = true
      pool    = bigip_ltm_pool.mypool.name
    }
  }
  depends_on = [bigip_ltm_pool.mypool]
}

resource "bigip_ltm_virtual_server" "test_virtual_server_policyattch_detach" {
  name        = "/Common/%s"
  destination = "192.168.10.11"
  port        = 80
  description = "Test virtual server"
  pool        = bigip_ltm_pool.mypool.name
  ip_protocol = "tcp"
  profiles = [
    "/Common/tcp",
    "/Common/http"
  ]
  persistence_profiles = [
    "/Common/cookie"
  ]
  source_address_translation = "automap"
  translate_address          = "enabled"
  policies = [
    bigip_ltm_policy.test-policy.name
  ]
}
`, vsName)
}

func testVSCreateAttach(poolName, policyName, vsName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool" "%[1]s" {
  name                = "/Common/%[1]s"
  allow_nat           = "yes"
  allow_snat          = "yes"
  load_balancing_mode = "round-robin"
}

resource "bigip_ltm_policy" "%[2]s" {
  name     = "/Common/%[2]s"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule6"
    action {
      forward = true
      pool    = bigip_ltm_pool.%[1]s.name
    }
  }
  depends_on = [bigip_ltm_pool.%[1]s]
}

resource "bigip_ltm_virtual_server" "%[3]s" {
  name        = "/Common/%[3]s"
  destination = "192.168.10.11"
  port        = 80
  description = "Test virtual server"
  pool        = bigip_ltm_pool.%[1]s.name
  ip_protocol = "tcp"
  profiles = [
    "/Common/tcp",
    "/Common/http"
  ]
  persistence_profiles = [
    "/Common/cookie"
  ]
  source_address_translation = "automap"
  translate_address          = "enabled"
  policies = [
    bigip_ltm_policy.%[2]s.name
  ]
}
`, poolName, policyName, vsName)
}

func testVSCreateDetatch(poolName, policyName, vsName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool" "%[1]s" {
  name                = "/Common/%[1]s"
  allow_nat           = "yes"
  allow_snat          = "yes"
  load_balancing_mode = "round-robin"
}

resource "bigip_ltm_policy" "%[2]s" {
  name     = "/Common/%[2]s"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule6"
    action {
      forward = true
      pool    = bigip_ltm_pool.%[1]s.name
    }
  }
  depends_on = [bigip_ltm_pool.%[1]s]
}

resource "bigip_ltm_virtual_server" "%[3]s" {
  name        = "/Common/%[3]s"
  destination = "192.168.10.11"
  port        = 80
  description = "Test virtual server"
  //pool = bigip_ltm_pool.%[1]s.name
  ip_protocol = "tcp"
  profiles = [
    "/Common/tcp",
    "/Common/http"
  ]
  persistence_profiles = [
    "/Common/cookie"
  ]
  source_address_translation = "automap"
  translate_address          = "enabled"
  policies = [
    bigip_ltm_policy.%[2]s.name
  ]
}
`, poolName, policyName, vsName)
}

func testVSCreatePolicyDettach(vsName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool" "mypool" {
  name                = "/Common/test-pool"
  allow_nat           = "yes"
  allow_snat          = "yes"
  load_balancing_mode = "round-robin"
}

resource "bigip_ltm_policy" "test-policy" {
  name     = "/Common/test-policy"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule6"
    action {
      forward = true
      pool    = bigip_ltm_pool.mypool.name
    }
  }
  depends_on = [bigip_ltm_pool.mypool]
}

resource "bigip_ltm_virtual_server" "test_virtual_server_policyattch_detach" {
  name        = "/Common/%s"
  destination = "192.168.10.11"
  port        = 80
  description = "Test virtual server"
  pool        = bigip_ltm_pool.mypool.name
  ip_protocol = "tcp"
  profiles = [
    "/Common/tcp",
    "/Common/http"
  ]
  persistence_profiles = [
    "/Common/cookie"
  ]
  source_address_translation = "automap"
  translate_address          = "enabled"
  //policies = [
  //  bigip_ltm_policy.test-policy.name
  //]
}
`, vsName)
}

func testVSCreatevlanEnabled(vsName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "%s" {
  name          = "/Common/%s"
  destination   = "192.168.50.11"
  port          = 80
  vlans_enabled = true
  vlans         = ["/Common/external"]
}
`, vsName, vsName)
}

func testCheckVSsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_virtual_address" {
			continue
		}

		name := rs.Primary.ID
		vs, err := client.GetVirtualServer(name)
		if err != nil {
			return err
		}
		if vs != nil {
			return fmt.Errorf("Virtual server %s not destroyed. ", name)
		}
	}
	return nil
}
