package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_NTP_NAME = fmt.Sprintf("/%s/test-ntp", TEST_PARTITION)

var TEST_NTP_RESOURCE = `
resource "bigip_sys_ntp" "test-ntp" {
	description = "` + TEST_NTP_NAME + `"
	servers = ["10.10.10.10"]
	timezone = "America/Los_Angeles"
}
`

func TestBigipSysNtp_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//CheckDestroy: testCheckntpsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_NTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckntpExists(TEST_NTP_NAME, true),
					resource.TestCheckResourceAttr("bigip_sys_ntp.test-ntp", "description", TEST_NTP_NAME),
					//resource.TestCheckResourceAttr("bigip_sys_ntp.test-ntp", "servers", "[10.10.10.10]"),
					resource.TestCheckResourceAttr("bigip_sys_ntp.test-ntp", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("bigip_sys_ntp.test-ntp",
						fmt.Sprintf("servers.%d", schema.HashString("10.10.10.10")),
						"10.10.10.10"),
				),
			},
		},
	})
}

func TestBigipSysNtp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//	CheckDestroy: testCheckntpsDestroyed, ( No Delet API support)
		Steps: []resource.TestStep{
			{
				Config: TEST_NTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckntpExists(TEST_NTP_NAME, true),
				),
				ResourceName:      TEST_NTP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckntpExists(description string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		ntp, err := client.NTPs()
		if err != nil {
			return err
		}
		if exists && ntp == nil {
			return fmt.Errorf("ntp %s was not created.", description)

		}
		if !exists && ntp != nil {
			return fmt.Errorf("ntp %s still exists.", description)

		}
		return nil
	}
}

func testCheckntpsDestroyed(s *terraform.State) error {
	/* client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_sys_ntp" {
			continue
		}

		description := rs.Primary.ID
		ntp, err := client.NTPs()
		if err != nil {
			return err
		}
		if ntp != nil {
			return fmt.Errorf("ntp ", description, " not destroyed.")

		}
	}*/
	return nil
}
