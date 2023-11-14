package bigip

import (
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TEST_PARTITION_NAME = "test-partition"
var TEST_PARTITION_RESOURCE_1 = `
resource "bigip_partition" "test-partition" {
  name = "` + TEST_PARTITION_NAME + `"
  description = "created by teraform"
  route_domain_id = 0
}`

var TEST_PARTITION_RESOURCE_2 = `
resource "bigip_partition" "test-partition" {
  name = "` + TEST_PARTITION_NAME + `"
  description = "updated by teraform"
  route_domain_id = 2
}`

func TestAccPartitionCreateUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckPartitionsDestroyed,
		Steps: []resource.TestStep{
			{
				Config:  TEST_PARTITION_RESOURCE_1,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckPartitionExists(TEST_PARTITION_NAME),
					resource.TestCheckResourceAttr("bigip_partition.test-partition", "name", TEST_PARTITION_NAME),
					resource.TestCheckResourceAttr("bigip_partition.test-partition", "description", "created by teraform"),
					resource.TestCheckResourceAttr("bigip_partition.test-partition", "route_domain_id", "0"),
				),
			},
			{
				Config: TEST_PARTITION_RESOURCE_2,
				Check: resource.ComposeTestCheckFunc(
					testCheckPartitionExists(TEST_PARTITION_NAME),
					resource.TestCheckResourceAttr("bigip_partition.test-partition", "name", TEST_PARTITION_NAME),
					resource.TestCheckResourceAttr("bigip_partition.test-partition", "description", "updated by teraform"),
					resource.TestCheckResourceAttr("bigip_partition.test-partition", "route_domain_id", "2"),
				),
			},
		},
	})
}

func testCheckPartitionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		partition, err := client.GetPartition(name)
		if err != nil {
			return err
		}
		if partition.Name != name {
			return err
		}
		return nil
	}
}

func testCheckPartitionsDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_partition" {
			continue
		}
		_, err := client.GetPartition(rs.Primary.ID)
		if err == nil {
			return err
		}
	}
	return nil
}
