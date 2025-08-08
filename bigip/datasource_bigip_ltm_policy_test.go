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
