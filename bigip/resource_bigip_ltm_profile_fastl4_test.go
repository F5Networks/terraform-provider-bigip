package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_FASTL4_NAME = fmt.Sprintf("/%s/test-fastl4", TEST_PARTITION)

var TEST_FASTL4_RESOURCE = `
resource "bigip_fastl4_profile" "sjfastl4profile"

        {
            name = "` + TEST_FASTL4_NAME + `"
            partition = "Common"
            defaults_from = "/Common/fastL4"
            client_timeout = 40
            explicitflow_migration = "enabled"
            hardware_syncookie = "enabled"
            idle_timeout = 200
            iptos_toclient = "pass-through"
            iptos_toserver = "pass-through"
            keepalive_interval = "disabled"  //This cannot take enabled
        }`

func TestBigipLtmFastl4_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPoolsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_FASTL4_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckFastl4Exists(TEST_FASTL4_NAME, true),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "name", TEST_FASTL4_NAME),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "defaults_from", "/Common/fastL4"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "client_timeout", "40"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "explicitflow_migration", "enabled"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "hardware_syncookie", "enabled"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "idle_timeout", "200"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "hardware_syncookie", "enabled"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "iptos_toclient", "pass-through"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "iptos_toserver", "pass-through"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "keepalive_interval", "disabled"),
				),
			},
		},
	})
}

func TestBigipLtmFastl4_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckFastl4Destroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_FASTL4_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckFastl4Exists(TEST_FASTL4_NAME, true),
				),
				ResourceName:      TEST_FASTL4_RESOURCE,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

//TODO: test adding/removing nodes

func testCheckFastl4Exists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.Fastl4(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("FASTL4 ", name, " does not exist.")
		}
		if !exists && p != nil {
			return fmt.Errorf("FASTL4 ", name, " exists.")
		}
		return nil
	}
}

func testCheckFastl4Destroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_fastl4_profile" {
			continue
		}

		name := rs.Primary.ID
		fastl4, err := client.Fastl4(name)
		if err != nil {
			return err
		}
		if fastl4 != nil {
			return fmt.Errorf("FastL4 ", name, " not destroyed.")
		}
	}
	return nil
}
