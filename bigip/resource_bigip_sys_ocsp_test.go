package bigip

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testSysOcspDNS = `
resource "bigip_sys_ocsp" "test-ocsp" {
  name              = "/Common/test-ocsp"
  dns_resolver      = "/Common/f5-aws-dns"
  signer_key        = "/Common/le-ssl"
  signer_cert       = "/Common/le-ssl"
  passphrase        = "testabcdef"
}
`

const testSysOcspProxy = `
resource "bigip_sys_ocsp" "test-ocsp" {
  name              = "/Common/test-ocsp"
  proxy_server_pool = "/Common/test-poolxyz"
  signer_key        = "/Common/le-ssl"
  signer_cert       = "/Common/le-ssl"
  passphrase        = "testabcdef"
}
`

func TestAccBigipSysOCSP_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckOCSPDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testSysOcspDNS,
				Check: resource.ComposeTestCheckFunc(
					testCheckOCSPExists("~Common~test-ocsp"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "name", "/Common/test-ocsp"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "dns_resolver", "/Common/f5-aws-dns"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "signer_key", "/Common/le-ssl"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "signer_cert", "/Common/le-ssl"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "passphrase", "testabcdef"),
				),
			},
			{
				Config: testSysOcspProxy,
				Check: resource.ComposeTestCheckFunc(
					testCheckOCSPExists("~Common~test-ocsp"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "name", "/Common/test-ocsp"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "proxy_server_pool", "/Common/test-poolxyz"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "signer_key", "/Common/le-ssl"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "signer_cert", "/Common/le-ssl"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "passphrase", "testabcdef"),
				),
			},
			{
				Config: testSysOcspProxy,
				Check: resource.ComposeTestCheckFunc(
					testCheckOCSPExists("~Common~test-ocsp"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "name", "/Common/test-ocsp"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "proxy_server_pool", "/Common/test-poolxyz"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "signer_key", "/Common/le-ssl"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "signer_cert", "/Common/le-ssl"),
					resource.TestCheckResourceAttr("bigip_sys_ocsp.test-ocsp", "passphrase", "testabcdef"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testCheckOCSPExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		p, err := client.GetOCSP(name)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("OCSP %s does not exist ", name)
		}

		return nil
	}
}

func testCheckOCSPDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_sys_ocsp" {
			continue
		}

		id := rs.Primary.ID
		id = strings.Trim(id, "/")
		splitArr := strings.Split(id, "/")
		if len(splitArr) != 2 {
			return fmt.Errorf("Invalid ID %s", id)
		}

		name := splitArr[1]
		partition := splitArr[0]
		ocspFqdn := fmt.Sprintf("~%s~%s", partition, name)
		// client.DeleteOCSP(ocspFqdn)
		ocsp, err := client.GetOCSP(ocspFqdn)
		js, _ := json.Marshal(ocsp)
		if err != nil {
			return err
		}
		if ocsp != nil {
			return fmt.Errorf("OCSP %s not destroyed, struct %+v, js %s", name, ocsp, string(js))
		}
	}
	return nil
}
