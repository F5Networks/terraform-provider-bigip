package bigip

//TODO: delete not implemented in virtual address

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
	"testing"
)

var TEST_VA_NAME = fmt.Sprintf("/%s/test-va", TEST_PARTITION)

var TEST_VA_RESOURCE = `
resource "bigip_ltm_virtual_address" "test-va" {
	name = "` + TEST_VA_NAME + `"

}
`

func TestBigipLtmVA_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVAsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_VA_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVAExists(TEST_VA_NAME, true),
				),
			},
		},
	})
}

func TestBigipLtmVA_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckVAsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_VA_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVAExists(TEST_VA_NAME, true),
				),
				ResourceName:      TEST_VA_NAME,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckVAExists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		vas, err := client.VirtualAddresses()
		if err != nil {
			return err
		}

		for _, va := range vas.VirtualAddresses {
			if va.FullPath == name {
				if !exists {
					return fmt.Errorf("Virtual address " + name + " exists.")
				} else {
					return nil
				}
			}
		}
		if exists {
			return fmt.Errorf("Virtual address " + name + " does not exist.")
		}

		return nil
	}
}

func testCheckVAsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_virtual_address" {
			continue
		}

		name := rs.Primary.ID
		vas, err := client.VirtualAddresses()
		if err != nil {
			return err
		}
		for _, va := range vas.VirtualAddresses {
			if va.FullPath == name {
				return fmt.Errorf("Virtual address ", name, " not destroyed.")
			}
		}
	}
	return nil
}
