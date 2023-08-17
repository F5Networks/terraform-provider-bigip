package bigip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testResourceSSLKeyCert = `
resource "bigip_ssl_key_cert" "testkeycert" {
  partition   = "Common"
  key_name    = "ssl-test-key"
  key_content = "${file("` + folder + `/../examples/serverkey.key")}"
  cert_name    = "ssl-test-cert"
  cert_content = "${file("` + folder + `/../examples/servercert.crt")}"
}
`

func TestAccBigipSSLCertKeyCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		// CheckDestroy:
		Steps: []resource.TestStep{
			{
				Config: testResourceSSLKeyCert,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "key_name", "ssl-test-key"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_name", "ssl-test-cert"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "partition", "Common"),
				),
				Destroy: false,
			},
			{
				Config: testResourceSSLKeyCert,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "key_name", "ssl-test-key"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "cert_name", "ssl-test-cert"),
					resource.TestCheckResourceAttr("bigip_ssl_key_cert.testkeycert", "partition", "Common"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
