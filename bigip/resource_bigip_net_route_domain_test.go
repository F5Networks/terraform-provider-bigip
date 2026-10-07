/*
Copyright 2019 F5 Networks Inc.
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

var TEST_ROUTE_DOMAIN_NAME = fmt.Sprintf("/%s/test-rd", TestPartition)

var TEST_ROUTE_DOMAIN_RESOURCE = `
resource "bigip_net_route_domain" "test-rd" {
	name   = "` + TEST_ROUTE_DOMAIN_NAME + `"
	id_    = 100
	strict = true
}
`

var TEST_ROUTE_DOMAIN_RESOURCE_UPDATE = `
resource "bigip_net_route_domain" "test-rd" {
	name   = "` + TEST_ROUTE_DOMAIN_NAME + `"
	id_    = 100
	strict = false
}
`

var TEST_ROUTE_DOMAIN_RESOURCE_WITH_VLAN = `
resource "bigip_net_vlan" "test-vlan-rd" {
	name = "/Common/test-vlan-rd"
	tag  = 200
}
resource "bigip_net_route_domain" "test-rd-vlan" {
	name   = "/Common/test-rd-vlan"
	id_    = 101
	strict = true
	vlans  = ["/Common/test-vlan-rd"]
	depends_on = ["bigip_net_vlan.test-vlan-rd"]
}
`

func TestAccBigipNetRouteDomain_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckRouteDomainDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ROUTE_DOMAIN_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckRouteDomainExists(TEST_ROUTE_DOMAIN_NAME, true),
					resource.TestCheckResourceAttr("bigip_net_route_domain.test-rd", "name", TEST_ROUTE_DOMAIN_NAME),
					resource.TestCheckResourceAttr("bigip_net_route_domain.test-rd", "id_", "100"),
					resource.TestCheckResourceAttr("bigip_net_route_domain.test-rd", "strict", "true"),
				),
			},
		},
	})
}

func TestAccBigipNetRouteDomain_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckRouteDomainDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ROUTE_DOMAIN_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckRouteDomainExists(TEST_ROUTE_DOMAIN_NAME, true),
					resource.TestCheckResourceAttr("bigip_net_route_domain.test-rd", "strict", "true"),
				),
			},
			{
				Config: TEST_ROUTE_DOMAIN_RESOURCE_UPDATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckRouteDomainExists(TEST_ROUTE_DOMAIN_NAME, true),
					resource.TestCheckResourceAttr("bigip_net_route_domain.test-rd", "strict", "false"),
				),
			},
		},
	})
}

func TestAccBigipNetRouteDomain_withVlan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckRouteDomainDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ROUTE_DOMAIN_RESOURCE_WITH_VLAN,
				Check: resource.ComposeTestCheckFunc(
					testCheckRouteDomainExists("/Common/test-rd-vlan", true),
					resource.TestCheckResourceAttr("bigip_net_route_domain.test-rd-vlan", "id_", "101"),
					resource.TestCheckResourceAttr("bigip_net_route_domain.test-rd-vlan", "strict", "true"),
					resource.TestCheckTypeSetElemAttr("bigip_net_route_domain.test-rd-vlan", "vlans.*", "/Common/test-vlan-rd"),
				),
			},
		},
	})
}

func TestAccBigipNetRouteDomain_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckRouteDomainDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_ROUTE_DOMAIN_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckRouteDomainExists(TEST_ROUTE_DOMAIN_NAME, true),
				),
			},
			{
				ResourceName:      "bigip_net_route_domain.test-rd",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckRouteDomainExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		rd, err := getRouteDomain(client, name)
		if err != nil {
			return err
		}
		if exists && rd == nil {
			return fmt.Errorf("route domain %s was not created", name)
		}
		if !exists && rd != nil {
			return fmt.Errorf("route domain %s still exists", name)
		}
		return nil
	}
}

func testCheckRouteDomainDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_net_route_domain" {
			continue
		}

		name := rs.Primary.ID
		rd, err := getRouteDomain(client, name)
		if err != nil {
			// RouteDomains() failing entirely is a real error
			if !strings.Contains(err.Error(), "not found") {
				return err
			}
			return nil
		}
		if rd != nil {
			return fmt.Errorf("route domain %s was not destroyed", name)
		}
	}
	return nil
}
