package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

//var TEST_PROVISION_NAME = fmt.Sprintf("/%s/test-provision", TEST_PARTITION)
var TEST_PROVISION_NAME = "afm"

var TEST_PROVISION_RESOURCE = `
resource "bigip_sys_provision" "test-provision" {
 name = "` + TEST_PROVISION_NAME + `"
 full_path  = "afm"
 cpu_ratio = 0
 disk_ratio = 0
 level = "nominal"
 memory_ratio = 0
}
`

func TestAccBigipSysProvision_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckProvisionsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_PROVISION_NAME, true),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "name", TEST_PROVISION_NAME),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "full_path", "afm"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "cpu_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "disk_ratio", "0"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "level", "nominal"),
					resource.TestCheckResourceAttr("bigip_sys_provision.test-provision", "memory_ratio", "0"),
				),
			},
		},
	})
}

func TestAccBigipSysProvision_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//CheckDestroy: testCheckProvisionsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_PROVISION_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckprovisionExists(TEST_PROVISION_NAME, true),
				),
				ResourceName:      TEST_PROVISION_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckprovisionExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		provision, err := client.Provisions(name)
		if err != nil {
			return err
		}
		if exists && provision == nil {
			return fmt.Errorf("provision %s was not created.", name)

		}
		if !exists && provision != nil {
			return fmt.Errorf("provision %s still exists.", name)

		}
		return nil
	}
}

func testCheckProvisionsDestroyed(s *terraform.State) error {
	/*client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_sys_provision" {
			continue
		}

		name := rs.Primary.ID
		provision, err := client.Provisions(name)
		if err != nil {
			return err
		}
		if provision != nil {
			return fmt.Errorf("provision ", name, " not destroyed.")

		}
	}*/
	return nil
}
