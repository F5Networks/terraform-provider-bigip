package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_TCP_NAME = fmt.Sprintf("/%s/test-tcp", TEST_PARTITION)

var TEST_TCP_RESOURCE = `
resource "bigip_tcp_profile" "sanjose-tcp-wan-profile"

        {
            name = "sanjose-tcp-wan-profile"
            defaults_from = "/Common/tcp-wan-optimized"
            idle_timeout = 300
            close_wait_timeout = 5
            finwait_2timeout = 5
            finwait_timeout = 300
            keepalive_interval = 1700
            deferred_accept = "enabled"
            fast_open = "enabled"
        }
`

func TestBigipLtmtcp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_TCP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(TEST_TCP_NAME, true),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "name", "sanjose-tcp-wan-profile"),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "defaults_from", "/Common/tcp-wan-optimized"),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "idle_timeout", "300"),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "close_wait_timeout", "5"),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "finwait_2timeout", "5"),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "finwait_timeout", "300"),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "keepalive_interval", "1700"),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "deferred_accept", "enabled"),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "manual_resume", "false"),
					resource.TestCheckResourceAttr("bigip_tcp_profile.test-tcp", "fast_open", "enabled"),
				),
			},
		},
	})
}

func TestBigipLtmTcp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckTcpsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_TCP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckTcpExists(TEST_TCP_NAME, true),
				),
				ResourceName:      TEST_TCP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckTcpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		tcps, err := client.Tcp()
		if err != nil {
			return err
		}
		if exists && tcps == nil {
			return fmt.Errorf("tcp profile ", name, " was not created.")
		}
		if !exists && tcps != nil {
			return fmt.Errorf("tcp profile ", name, " still exists.")
		}
		return nil

	}
}

func testCheckTcpsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_tcp_profile" {
			continue
		}

		name := rs.Primary.ID
		tcp, err := client.Tcp()
		if err != nil {
			return err
		}
		if tcp == nil {
			return fmt.Errorf("fasthttp ", name, " not destroyed.")
		}
	}
	return nil
}
