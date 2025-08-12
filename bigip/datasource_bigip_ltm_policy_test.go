package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBigipLtmPolicy_basic(t *testing.T) {
	policyName := "/Common/test_policy"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipLtmPolicyConfig(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bigip_ltm_policy.test", "name", policyName),
					resource.TestCheckResourceAttrSet("data.bigip_ltm_policy.test", "strategy"),
					resource.TestCheckResourceAttrSet("data.bigip_ltm_policy.test", "controls"),
					resource.TestCheckResourceAttrSet("data.bigip_ltm_policy.test", "requires"),
					resource.TestCheckResourceAttrSet("data.bigip_ltm_policy.test", "rule"),
				),
			},
		},
	})
}

func testAccDataSourceBigipLtmPolicyConfig(name string) string {
	return fmt.Sprintf(`
data "bigip_ltm_policy" "test" {
  name = "%s"
}
`, name)
}

// add test for subfolder policy
func TestAccDataSourceBigipLtmPolicy_subfolder(t *testing.T) {
	policyName := "/Common/folder1/ecopolicy"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipLtmPolicySubfolderConfig(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bigip_ltm_policy.test2", "name", policyName),
				),
			},
		},
	})
}

func testAccDataSourceBigipLtmPolicySubfolderConfig(name string) string {
	return fmt.Sprintf(`
data "bigip_ltm_policy" "test2" {
  name = "%s"
}
`, name)
}

// add test for subfolder policy with different name
func TestAccDataSourceBigipLtmPolicy_subfolder2(t *testing.T) {
	policyName := "/Common/testpolicy"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigipLtmPolicySubfolder2Config(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bigip_ltm_policy.test2", "name", policyName),
				),
			},
		},
	})
}

func testAccDataSourceBigipLtmPolicySubfolder2Config(name string) string {
	return fmt.Sprintf(`
data "bigip_ltm_policy" "test2" {
  name = "%s"
}
`, name)
}
