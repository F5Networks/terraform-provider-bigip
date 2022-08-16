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
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var poolMember = fmt.Sprintf("%s:443", "10.10.10.10")
var poolMemberFqdn = fmt.Sprintf("%s:443", "www.google.com")
var poolMemberFullpath = fmt.Sprintf("/%s/%s", TEST_PARTITION, poolMember)
var poolMemberFqdnFullpath = fmt.Sprintf("/%s/%s", TEST_PARTITION, poolMemberFqdn)

var TestPoolResource1 = `
resource "bigip_ltm_pool" "test-pool" {
	name = "` + TestPoolName + `"
	monitors = ["/Common/http"]
	allow_nat = "yes"
	allow_snat = "yes"
	description = "Test-Pool-Sample"
	load_balancing_mode = "round-robin"
	slow_ramp_time = "5"
	service_down_action = "reset"
	reselect_tries = "2"
}
resource "bigip_ltm_pool_attachment" "test-pool_test-node" {
	pool = bigip_ltm_pool.test-pool.name
	node = "` + poolMember + `"
}
`
var TestPoolResource2 = `
resource "bigip_ltm_pool" "test-pool" {
        name = "` + TestPoolName + `"
        monitors = ["/Common/http"]
        allow_nat = "yes"
        allow_snat = "yes"
        description = "Test-Pool-Sample"
        load_balancing_mode = "round-robin"
        slow_ramp_time = "5"
        service_down_action = "reset"
        reselect_tries = "2"
}
resource "bigip_ltm_pool_attachment" "test-pool_test-node" {
	pool = bigip_ltm_pool.test-pool.name
	node = "` + poolMember + `"
    ratio                 = 2
    connection_limit      = 2
    connection_rate_limit = 2
    priority_group        = 2
    dynamic_ratio         = 3
}
`
var TestPoolResource3 = `
resource "bigip_ltm_pool" "test-pool" {
        name = "` + TestPoolName + `"
        monitors = ["/Common/http"]
        allow_nat = "yes"
        allow_snat = "yes"
        description = "Test-Pool-Sample"
        load_balancing_mode = "round-robin"
        slow_ramp_time = "5"
        service_down_action = "reset"
        reselect_tries = "2"
}
resource "bigip_ltm_pool_attachment" "test-pool_test-node" {
	pool = bigip_ltm_pool.test-pool.name
	node = "` + poolMemberFqdn + `"
}
`
var TestPoolResource4 = `
resource "bigip_ltm_node" "test-node" {
        name = "` + TestNodeName + `"
        address = "10.10.10.11"
        connection_limit = "0"
        dynamic_ratio = "1"
        monitor = "default"
        rate_limit = "disabled"
        fqdn {
    address_family = "ipv4"
    interval       = "3000"
  }
}
resource "bigip_ltm_pool" "test-pool" {
        name = "` + TestPoolName + `"
        monitors = ["/Common/http"]
        allow_nat = "yes"
        allow_snat = "yes"
        description = "Test-Pool-Sample"
        load_balancing_mode = "round-robin"
        slow_ramp_time = "5"
        service_down_action = "reset"
        reselect_tries = "2"
}
`

func TestAccBigipLtmPoolAttachment_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPoolResource1,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", poolMember),
				),
			},
		},
	})
}
func TestAccBigipLtmPoolAttachment_createFqdn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPoolResource3,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFqdnFullpath, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", poolMemberFqdn),
				),
			},
		},
	})
}

func TestAccBigipLtmPoolAttachment_Issue381(t *testing.T) {
	t.Parallel()
	TestPoolName = "/Common/k8s_example_pool"
	poolMemberFullpath2 := "/Common/node2.com:31380"
	poolMemberFullpath1 := "/Common/node1.com:31380"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmPoolattachIssu381(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath1, true),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath2, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_1", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_1", "node", poolMemberFullpath1),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_1", "priority_group", "2"),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_2", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_2", "node", poolMemberFullpath2),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_2", "priority_group", "1"),
				),
			},
		},
	})
}
func TestAccBigipLtmPoolAttachment_Issue380(t *testing.T) {
	t.Parallel()
	TestPoolName = "/Common/test-pool-issue380"
	poolMemberFullpath1 := "/Common/www.perdu.com:31380"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmPoolattachIssu380(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath1, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_1", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_1", "node", poolMemberFullpath1),
				),
			},
		},
	})
}
func TestAccBigipLtmPoolAttachment_Issue92(t *testing.T) {
	t.Parallel()
	TestPoolName = "/Common/test-pool-issue92"
	poolMemberFullpath1 := "/Common/test-node-issue92:0"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmPoolattachIssu92(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath1, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_1", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.example_pool_member_1", "node", poolMemberFullpath1),
				),
			},
		},
	})
}
func TestAccBigipLtmPoolAttachment_Issue661(t *testing.T) {
	t.Parallel()
	TestPoolName = "/TEST2/test-pool-issue661"
	poolMemberFullpath1 := "/TEST2/2.3.2.2%30:8080"
	poolMem2 := "2.3.2.2%30:8080"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testaccbigipltmPoolattachIssu661(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath1, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.attach", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.attach", "node", poolMem2),
				),
			},
		},
	})
}
func TestAccBigipLtmPoolAttachment_Modify(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPoolResource1,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", poolMember),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPoolResource2,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", poolMember),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "connection_limit", "2"),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "connection_rate_limit", "2"),
				),
			},
		},
	})
}

func TestAccBigipLtmPoolAttachment_Delete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TestPoolResource1,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath, true),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "pool", TestPoolName),
					resource.TestCheckResourceAttr("bigip_ltm_pool_attachment.test-pool_test-node", "node", poolMember),
				),
			},
			{
				Config: TestPoolResource4,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TestPoolName),
					testCheckPoolAttachment(TestPoolName, poolMemberFullpath, false),
				),
			},
		},
	})
}
func testCheckPoolAttachment(poolName string, expected string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		pool, err := client.GetPool(poolName)
		if err != nil {
			return err
		}
		if pool == nil {
			return fmt.Errorf("Pool %s does not exist ", poolName)
		}

		nodes, err := client.PoolMembers(poolName)
		if err != nil {
			return fmt.Errorf("Error retrieving pool (%s) members: %s ", poolName, err)
		}
		if nodes == nil {
			return fmt.Errorf("Pool member %s does not exist. ", expected)
		}
		found := false
		for _, node := range nodes.PoolMembers {
			if expected == node.FullPath {
				found = true
				break
			}
		}

		if !found && exists {
			return fmt.Errorf("Node %s is not a member of pool %s ", expected, poolName)
		}
		if found && !exists {
			return fmt.Errorf("Node %s is still  a member of pool %s", expected, poolName)
		}

		return nil
	}
}

func testaccbigipltmPoolattachIssu381() string {
	tfConfig := `
		resource "bigip_ltm_pool" "example_pool" {
  			name = "/Common/k8s_example_pool"
  			load_balancing_mode = "round-robin"
  			description = "Example pool"
  			monitors = [
    			"/Common/tcp"
  			]
  			allow_snat = "yes"
  			allow_nat = "yes"
  			minimum_active_members = 1
		}
		resource "bigip_ltm_node" "node1" {
  			name = "/Common/node1.com"
  			description = "Terraform managed"
  			address = "2.22.22.22"
		}
		resource "bigip_ltm_node" "node2" {
		  name = "/Common/node2.com"
		  description = "Terraform managed"
		  address = "3.23.23.23"
		}
		resource "bigip_ltm_pool_attachment" "example_pool_member_1" {
		  pool = bigip_ltm_pool.example_pool.name
		  node = format("%s:%s",bigip_ltm_node.node1.name,31380)
		  priority_group = 2
		}
		resource "bigip_ltm_pool_attachment" "example_pool_member_2" {
		  pool = bigip_ltm_pool.example_pool.name
		  node = format("%s:%s",bigip_ltm_node.node2.name,31380)
		  priority_group = 1
		}`
	return tfConfig
}

func testaccbigipltmPoolattachIssu380() string {
	tfConfig := `
		resource "bigip_ltm_pool" "example_pool" {
  			name = "/Common/test-pool-issue380"
  			load_balancing_mode = "round-robin"
  			description = "Example pool"
  			monitors = [
    			"/Common/tcp"
  			]
  			allow_snat = "yes"
  			allow_nat = "yes"
		}
		resource "bigip_ltm_node" "node1" {
  			name = "/Common/www.perdu.com"
  			description = "Terraform managed"
  			address = "100.22.22.22"
		}
		resource "bigip_ltm_pool_attachment" "example_pool_member_1" {
		  pool = bigip_ltm_pool.example_pool.name
		  node = format("%s:%s",bigip_ltm_node.node1.name,31380)
		}`
	return tfConfig
}

func testaccbigipltmPoolattachIssu92() string {
	tfConfig := `
		resource "bigip_ltm_pool" "example_pool" {
  			name = "/Common/test-pool-issue92"
  			load_balancing_mode = "round-robin"
  			description = "Example pool"
  			allow_snat = "yes"
  			allow_nat = "yes"
		}
		resource "bigip_ltm_node" "node1" {
  			name = "/Common/test-node-issue92"
  			description = "Terraform managed"
  			address = "172.17.240.182%10"
		}
		resource "bigip_ltm_pool_attachment" "example_pool_member_1" {
		  pool = bigip_ltm_pool.example_pool.name
		  node = format("%s:%s",bigip_ltm_node.node1.name,0)
		}`
	return tfConfig
}

func testaccbigipltmPoolattachIssu661() string {
	tfConfig := `
		resource "bigip_command" "create_partition" {
			commands = ["create auth partition TEST2","create net route-domain /TEST2/testdomain id 30"]
			when     = "apply"
		}
		resource "bigip_ltm_pool" "example" {
  			name = "/TEST2/test-pool-issue661"
			depends_on = [bigip_command.create_partition]
		}
		resource "bigip_ltm_pool_attachment" "attach" {
  			pool = bigip_ltm_pool.example.name
  			node = "2.3.2.2%30:8080"
  			connection_limit = 11
		}`
	return tfConfig
}
