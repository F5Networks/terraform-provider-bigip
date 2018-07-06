package bigip

import (
	"fmt"
	"testing"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TEST_DATAGROUP_NAME = "/" + TEST_PARTITION + "/test-datagroup"

var TEST_DATAGROUP_RESOURCE = `
	resource "bigip_ltm_datagroup" "test-datagroup" {
		name = "` + TEST_DATAGROUP_NAME + `"
		type = "string"

		record {
			name = "test-record-name"
			data = "test-record-data"
		}
		record {
			name = "test-record-name2"
			data = "test-record-data2"
		}
	}`

func TestAccBigipLtmDataGroup_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DATAGROUP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TEST_DATAGROUP_NAME),
				),
			},
		},
	})
}

func TestAccBigipLtmDataGroup_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckDataGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_DATAGROUP_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataGroupExists(TEST_DATAGROUP_NAME),
				),
				ResourceName:      TEST_DATAGROUP_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckDataGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		datagroup, err := client.GetDatagroup(name)
		if err != nil {
			return fmt.Errorf("Error while fetching Data Group: %v", err)

		}

		datagroup_name := fmt.Sprintf("/%s/%s", datagroup.Partition, datagroup.Name)
		if datagroup_name != name {
			return fmt.Errorf("Data Group name does not match. Expecting %s got %s.", name, datagroup_name)
		}
		return nil
	}
}

func testCheckDataGroupDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_datagroup" {
			continue
		}

		name := rs.Primary.ID
		datagroup, err := client.GetDatagroup(name)

		if err != nil {
			return nil
		}
		if datagroup != nil {
			return fmt.Errorf("Data Group %s not destroyed.", name)
		}
	}
	return nil
}
