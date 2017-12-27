package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_DEVICEGROUP_NAME = fmt.Sprintf("/%s/test-devicegroup1", TEST_PARTITION)

var TEST_DEVICEGROUP_RESOURCE = `
resource "bigip_cm_devicegroup" "test-devicegroup"

        {
            name = "/Common/test-devicegroup1"
            auto_sync = "enabled"
            full_load_on_sync = "true"
            type = "sync-only"
        }
`

func TestBigipCmdevicegroup_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckdevicegroupsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DEVICEGROUP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdevicegroupExists(TEST_DEVICEGROUP_NAME, true),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "name", "/Common/test-devicegroup1"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "auto_sync", "enabled"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "full_load_on_sync", "true"),
					resource.TestCheckResourceAttr("bigip_cm_devicegroup.test-devicegroup", "type", "sync-only"),

				),
			},
		},
	})
}

func TestBigipLtmCmdevicegroup_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckdevicegroupsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DEVICEGROUP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdevicegroupExists(TEST_DEVICEGROUP_NAME, true),
				),
				ResourceName:      TEST_DEVICEGROUP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckdevicegroupExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.Devicegroups(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("devicegroup %s was not created.", name)
		}
		if !exists && p != nil {

			return fmt.Errorf("devicegroup %s still exists.", name)
		}
		return nil
	}
}

func testCheckdevicegroupsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_cm_devicegroup" {
			continue
		}

		name := rs.Primary.ID
		devicegroup, err := client.Devicegroups(name)
		if err != nil {
			return err
		}
		if devicegroup == nil {
			return fmt.Errorf("devicegroup %s not destroyed.", name)
		}
	}
	return nil
}
