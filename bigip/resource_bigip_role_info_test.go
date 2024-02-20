package bigip

import (
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TEST_ROLE_INFO_NAME = "test-roleinfo"
var TEST_ROLE_INFO_RESOURCE_1 = `
resource "bigip_role_info" "test-roleinfo" {
  name = "` + TEST_ROLE_INFO_NAME + `"
  description = "created by teraform"
  attribute = "attribute"
  console = "console"
  deny = "deny"
  role = "role"
  user_partition = "user_partition"
  line_order = 1
}`

func TestAccRoleInfoCreateUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckRoleInfoDestroyed,
		Steps: []resource.TestStep{
			{
				Config:  TEST_ROLE_INFO_RESOURCE_1,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleInfoExists(TEST_ROLE_INFO_NAME),
					resource.TestCheckResourceAttr("bigip_role_info.test-roleinfo", "name", TEST_ROLE_INFO_NAME),
					resource.TestCheckResourceAttr("bigip_role_info.test-roleinfo", "description", "created by teraform"),
					resource.TestCheckResourceAttr("bigip_role_info.test-roleinfo", "attribute", "attribute"),
					resource.TestCheckResourceAttr("bigip_role_info.test-roleinfo", "console", "console"),
					resource.TestCheckResourceAttr("bigip_role_info.test-roleinfo", "deny", "deny"),
					resource.TestCheckResourceAttr("bigip_role_info.test-roleinfo", "role", "role"),
					resource.TestCheckResourceAttr("bigip_role_info.test-roleinfo", "user_partition", "user_partition"),
				),
			},
		},
	})
}

func testCheckRoleInfoExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		roleInfo, err := client.GetRoleInfo(name)
		if err != nil {
			return err
		}
		if roleInfo.Name != name {
			return err
		}
		return nil
	}
}

func testCheckRoleInfoDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_role_info" {
			continue
		}
		_, err := client.GetRoleInfo(rs.Primary.ID)
		if err == nil {
			return err
		}
	}
	return nil
}
