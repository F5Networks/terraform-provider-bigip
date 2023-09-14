package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testResourceCipherRule = `
resource "bigip_ltm_cipher_rule" "testcipher" {
  name = "testcipher"
  partition = "Common"
  cipher_suites = "fips"
  dh_groups = "P256:P384:FFDHE2048:FFDHE3072:FFDHE4096"
  signature_algorithms = "DEFAULT"
}`

func TestAccCipherRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceCipherRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ltm_cipher_rule.testcipher", "name", "testcipher"),
					resource.TestCheckResourceAttr("bigip_ltm_cipher_rule.testcipher", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_ltm_cipher_rule.testcipher", "cipher_suites", "fips"),
					resource.TestCheckResourceAttr("bigip_ltm_cipher_rule.testcipher", "dh_groups", "P256:P384:FFDHE2048:FFDHE3072:FFDHE4096"),
					resource.TestCheckResourceAttr("bigip_ltm_cipher_rule.testcipher", "signature_algorithms", "DEFAULT"),
				),
			},
		},
	})
}
