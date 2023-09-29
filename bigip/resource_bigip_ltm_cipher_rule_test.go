package bigip

import (
	"fmt"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testCipherRuleConfigTC1 = `
resource "bigip_ltm_cipher_rule" "test-cipher-rule" {
  name   = "/Common/test-cipher-rule"
  cipher = "aes"
}
`

func TestAccBigipLtmCipherRuleCreateTC1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckCipherRuleDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testCipherRuleConfigTC1,
				Check: resource.ComposeTestCheckFunc(
					testCheckCipherRuleExists("/Common/test-cipher-rule"),
					resource.TestCheckResourceAttr("bigip_ltm_cipher_rule.test-cipher-rule", "name", "/Common/test-cipher-rule"),
				),
			},
		},
	})
}

func testCheckCipherRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.GetLtmCipherRule(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("Pool %s does not exist ", name)
		}

		return nil
	}
}

func testCheckCipherRuleDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_cipher_rule" {
			continue
		}

		name := rs.Primary.ID
		pool, err := client.GetLtmCipherRule(name)
		if err != nil {
			return err
		}
		if pool != nil {
			return fmt.Errorf("Cipher rule %s not destroyed ", name)
		}
	}
	return nil
}
