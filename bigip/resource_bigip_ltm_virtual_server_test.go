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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestVsName = fmt.Sprintf("/%s/test-vs", TestPartition)

var TestIruleResource = `
resource "bigip_ltm_irule" "test-rule-vstc1" {
  name  = "/Common/test-rule-vstc1"
  irule = <<EOF
when CLIENT_ACCEPTED {
     log local0. "test"
}
EOF
}
`
var TestVsResource = TestIruleResource + `
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
      connection = false
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
	irules = [bigip_ltm_irule.test-rule-vstc1.name]
	profiles = ["/Common/http"]
	client_profiles = ["/Common/tcp"]
	server_profiles = ["/Common/tcp-lan-optimized"]
	persistence_profiles = ["/Common/source_addr","/Common/hash"]
	default_persistence_profile = "/Common/hash"
	fallback_persistence_profile = "/Common/dest_addr"
    policies = [bigip_ltm_policy.http_to_https_redirect.name]
}
`
var TestVs6Name = fmt.Sprintf("/%s/test-vs6", TestPartition)

var TestVs6Resource = TestIruleResource + `
resource "bigip_ltm_virtual_server" "test-vs" {
	name = "` + TestVs6Name + `"
    destination = "fe80::11"
	description = "VirtualServer-test"
	port = 9999
	source_address_translation = "automap"
	ip_protocol = "tcp"
    irules = [bigip_ltm_irule.test-rule-vstc1.name] 
	//irules = [bigip_ltm_irule.test-rule.name]
	profiles = ["/Common/http"]
	client_profiles = ["/Common/tcp"]
	server_profiles = ["/Common/tcp-lan-optimized"]
	persistence_profiles = ["/Common/source_addr", "/Common/hash"]
	default_persistence_profile = "/Common/hash"
	fallback_persistence_profile = "/Common/dest_addr"
}
`

func TestAccBigipLtmVirtualServerCreateV4V6(t *testing.T) {
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
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "irules.*", "/Common/test-rule-vstc1"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "profiles.*", "/Common/http"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "client_profiles.*", "/Common/tcp"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "server_profiles.*", "/Common/tcp-lan-optimized"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "persistence_profiles.*", "/Common/source_addr"),
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
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "irules.0", "/Common/test-rule-vstc1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "default_persistence_profile", "/Common/hash"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "irules.*", "/Common/test-rule-vstc1"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "profiles.*", "/Common/http"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "client_profiles.*", "/Common/tcp"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "server_profiles.*", "/Common/tcp-lan-optimized"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "persistence_profiles.*", "/Common/source_addr"),
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs", "persistence_profiles.*", "/Common/hash"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "fallback_persistence_profile", "/Common/dest_addr"),
				),
			},
		},
	})
}

func TestAccBigipLtmVirtualServercreate_Defaultstate(t *testing.T) {
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
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					// resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
					//	fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")), "/Common/fastL4"),
				),
			},
		},
	})
}

func TestAccBigipLtmVirtualServerModify_stateDisabledtoEnabled(t *testing.T) {
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
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					// resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
					//	fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
					//	"/Common/fastL4"),
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
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					// resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
					//	fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
					//	"/Common/fastL4"),
				),
			},
		},
	})
}

func TestAccBigipLtmVirtualServerModify_stateEnabledtoDisabled(t *testing.T) {
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
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					// resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
					//	fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
					//	"/Common/fastL4"),
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
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs", "ip_protocol", "tcp"),
					// resource.TestCheckResourceAttr("bigip_ltm_virtual_server.test-vs",
					//	fmt.Sprintf("profiles.%d", schema.HashString("/Common/fastL4")),
					//	"/Common/fastL4"),
				),
			},
		},
	})
}
func TestAccBigipLtmVirtualServerPolicyattach_detach(t *testing.T) {
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
					resource.TestCheckTypeSetElemAttr(rsName, "policies.*", "/Common/test-policy-tc88"),
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
func TestAccBigipLtmVirtualServerPooolattach_detatch(t *testing.T) {
	var poolName = "test-pool-tc7"
	var policyName = "test-policy-tc7"
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

func TestAccBigipLtmVirtualServerVlan_EnabledDisabled(t *testing.T) {
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
					resource.TestCheckTypeSetElemAttr("bigip_ltm_virtual_server.test-vs-vlan", "vlans.*", "/Common/test-vlan-vsenable"),
				),
			},
		},
	})
}

func TestAccBigipLtmVirtualServerTCIssue712(t *testing.T) {
	vsTCIssue712Name1 := fmt.Sprintf("/%s/%s", "Common", "vs_issue_712_a")
	vsTCIssue712Name2 := fmt.Sprintf("/%s/%s", "Common", "vs_issue_712_b")
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
				Config: getVSTCIssue712Config(vsTCIssue712Name1, vsTCIssue712Name2),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(vsTCIssue712Name1),
					testCheckVSExists(vsTCIssue712Name2),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server1", "name", vsTCIssue712Name1),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server1", "destination", "192.168.1.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server1", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server1", "port", "8080"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server2", "name", vsTCIssue712Name2),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server2", "destination", "192.168.1.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server2", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server2", "port", "8080"),
				),
			},
		},
	})
}

func TestAccBigipLtmVirtualServerTCIssue736(t *testing.T) {
	vsTCIssue736Name1 := fmt.Sprintf("/%s/%s", "Common", "vs_issue_736_a")
	vsTCIssue736Name2 := fmt.Sprintf("/%s/%s", "Common", "vs_issue_736_b")
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
				Config: getVSTCIssue736Config(vsTCIssue736Name1, vsTCIssue736Name2),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(vsTCIssue736Name1),
					testCheckVSExists(vsTCIssue736Name2),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "name", vsTCIssue736Name1),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "destination", "192.168.50.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "name", vsTCIssue736Name2),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "destination", "192.168.60.1"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "mask", "255.255.255.255"),
				),
			},
			{
				Config: getVSTCIssue736ModifyConfig(vsTCIssue736Name1, vsTCIssue736Name2),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists(vsTCIssue736Name1),
					testCheckVSExists(vsTCIssue736Name2),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "name", vsTCIssue736Name1),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "destination", "192.168.48.0"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-a", "mask", "255.255.248.0"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "name", vsTCIssue736Name2),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "destination", "192.168.32.0"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "port", "800"),
					resource.TestCheckResourceAttr("bigip_ltm_virtual_server.server736-b", "mask", "255.255.224.0"),
				),
			},
		},
	})
}

func TestAccBigipLtmVirtualServerTCIssue729(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVSsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccBigipLtmVSImportIssue729(),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("/Common/test-vs-issue729"),
				),
			},
			{
				ResourceName:      "bigip_ltm_virtual_server.test_vs_issue729_import",
				ImportStateId:     "/Common/test-vs-issue729",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigipLtmVirtualServerTestCases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVSsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: loadFixtureString("../examples/bigip_ltm_virtual_server.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckVSExists("/Common/test_vs_tc1"),
					testCheckVSExists("/Common/test_vs_tc3"),
					testCheckVSExists("/Common/test_vs_tc4"),
					testCheckVSExists("/Common/test_vs_tc5"),
					testCheckVSExists("/Common/test_vs_tc6"),
					// resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.pa_tc9", "pool", "/Common/test_pool_pa_tc9"),
					// resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.pa_tc8", "pool", "/Common/test_pool_pa_tc1"),
					// resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.pa_tc8", "node", "1.1.12.2:80"),
					// resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.pa_tc6", "node", "/Common/test3.com:80"),
				),
			},
		},
	})
}

func TestAccBigipLtmVirtualServerimport(t *testing.T) {
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

func testaccBigipLtmVSImportIssue729() string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "test_vs_issue729" {
  name                     = "%s"
  trafficmatching_criteria = "/Common/test-virtualserver_VS_TMC_OBJ"
}
`, "/Common/test-vs-issue729")
}

func testCheckVSExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		vs, err := client.GetVirtualServer(name)
		if err != nil {
			return err
		}
		if vs == nil {
			return fmt.Errorf("Virtual server %s does not exist. ", name)
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
  ip_protocol = "tcp"
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
  ip_protocol = "tcp"
  mask        = "255.255.255.255"
}
`, vsName)
}

func testVSCreatePolicyAttach(vsName string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_pool" "mypool" {
  name                = "/Common/test-pool-policyattach"
  allow_nat           = "yes"
  allow_snat          = "yes"
  load_balancing_mode = "round-robin"
}

resource "bigip_ltm_policy" "test-policy" {
  name     = "/Common/test-policy-tc88"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule6"
    action {
      forward    = true
      connection = false
      pool       = bigip_ltm_pool.mypool.name
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
      forward    = true
      connection = false
      pool       = bigip_ltm_pool.%[1]s.name
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
      forward    = true
      connection = false
      pool       = bigip_ltm_pool.%[1]s.name
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
  name                = "/Common/test-pool-policyattach"
  allow_nat           = "yes"
  allow_snat          = "yes"
  load_balancing_mode = "round-robin"
}

resource "bigip_ltm_policy" "test-policy" {
  name     = "/Common/test-policy-tc88"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule6"
    action {
      forward    = true
      connection = false
      pool       = bigip_ltm_pool.mypool.name
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
resource "bigip_net_vlan" "test-vlan" {
  name = "/Common/test-vlan-vsenable"
  tag  = 1010
  interfaces {
    vlanport = 1.1
    tagged   = true
  }
}
resource "bigip_ltm_virtual_server" "%s" {
  name          = "/Common/%s"
  destination   = "192.168.50.11"
  port          = 80
  vlans_enabled = true
  vlans         = [bigip_net_vlan.test-vlan.name]
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

func getVSTCIssue712Config(profileName1, profileName2 string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "server1" {
  name        = "%v"
  destination = "192.168.1.1"
  port        = 8080
  ip_protocol = "tcp"
}
resource "bigip_ltm_virtual_server" "server2" {
  name        = "%v"
  destination = "192.168.1.1"
  port        = 8080
  ip_protocol = "udp"
}
`, profileName1, profileName2)
}

func getVSTCIssue736Config(vsName1, vsName2 string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "server736-a" {
  name        = "%s"
  destination = "192.168.50.1"
  port        = 800
  ip_protocol = "tcp"
  mask        = "255.255.255.255"
}
resource "bigip_ltm_virtual_server" "server736-b" {
  name        = "%s"
  destination = "192.168.60.1"
  port        = 800
  ip_protocol = "tcp"
  mask        = "32"
}
`, vsName1, vsName2)
}

func getVSTCIssue736ModifyConfig(vsName1, vsName2 string) string {
	return fmt.Sprintf(`
resource "bigip_ltm_virtual_server" "server736-a" {
  name        = "%s"
  destination = "192.168.48.0"
  port        = 800
  ip_protocol = "tcp"
  mask        = "255.255.248.0"
}
resource "bigip_ltm_virtual_server" "server736-b" {
  name        = "%s"
  destination = "192.168.32.0"
  port        = 800
  ip_protocol = "tcp"
  mask        = "19"
}
`, vsName1, vsName2)
}
