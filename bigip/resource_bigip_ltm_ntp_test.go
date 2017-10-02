package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_NTP_NAME = fmt.Sprintf("/%s/test-ntp", TEST_PARTITION)

var TEST_NTP_RESOURCE = `
resource "bigip_ltm_ntp" "test-ntp" {
	description = "` + TEST_NTP_NAME + `"
	description = "/Common/ntp1"
	servers = ["time.facebook.com"]
	timezone = "America/Los_Angeles"
}
`

func TestBigipLtmPool_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_NTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckNtpExists(TEST_NTP_NAME, true),
					resource.TestCheckResourceAttr("bigip_ltm_ntp.test-ntp", "description", TEST_NTP_NAME),
					resource.TestCheckResourceAttr("bigip_ltm_ntp.test-ntp", "servers", ["time.facebook.com"]),
					resource.TestCheckResourceAttr("bigip_ltm_ntp.test-ntp", "timezone", "America/Los_Angeles"),
				),
			},
		},
	})
}

func TestBigipLtmNtp_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckNTPDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_NTP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckPoolExists(TEST_NTP_NAME, true),
				),
				ResourceName:      TEST_NTP_RESOURCE,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//TODO: test adding/removing nodes

func testCheckNtpExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.NTPs()
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("Ntp ", description, " does not exist.")
		}
		if !exists && p != nil {
			return fmt.Errorf( "Ntp", name, " exists.")
		}
		return nil
	}
}






func testCheckNTPDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_ntp" {
			continue
		}

		description := rs.Primary.ID
		ntp, err := client.NTPs()
		if err != nil {
			return err
		}
		if ntp != nil {
			return fmt.Errorf("Ntp ", description, " not destroyed.")
		}
	}
	return nil
}
