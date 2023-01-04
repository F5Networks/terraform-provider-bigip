package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func dataWAFEntityURLCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_waf_entity_url" "test_entity_url" {
	  name                        = "test_entity_url"
	  type                        = "explicit"
	  protocol                    = "HTTP"
	  perform_staging             = true
	  signature_overrides_disable = [200002029, 200002092]
	  method_overrides {
		allow  = false
		method = "HTTPS"
	  }
	}	
	`, address)
}

func TestAccBigipWAFEntityURLUnit(t *testing.T) {
	setup()
	defer teardown()

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataWAFEntityURLCfg(server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bigip_waf_entity_url.test_entity_url", "type", "explicit"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_url.test_entity_url", "method_overrides.0.allow", "false"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_url.test_entity_url", "method_overrides.0.method", "HTTPS"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_url.test_entity_url", "perform_staging", "true"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_url.test_entity_url", "signature_overrides_disable.0", "200002029"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_url.test_entity_url", "signature_overrides_disable.1", "200002092"),
					resource.TestCheckNoResourceAttr("data.bigip_waf_entity_url.test_entity_url", "signature_overrides_disable.2"),
				),
			},
		},
	})
}
