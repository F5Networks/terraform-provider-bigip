package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func dataWAFEntityParametersCfg(address string) string {
	return fmt.Sprintf(`
	provider "bigip" {
	  address   = "%s"
	  username  = ""
	  password  = ""
	  login_ref = ""
	}
	data "bigip_waf_entity_parameter" "test_ep" {
	  name = "test_entity_param"
	  type = "explicit"
	  value_type = "user-input"
	  allow_empty_type = true
	  allow_repeated_parameter_name = true
	  attack_signatures_check = true
	  check_max_value_length = false
	  check_min_value_length = false
	  data_type = "alpha_numeric"
	  enable_regular_expression = false
	  is_base64 = false
	  is_cookie = false
	  is_header = false
	  level = "url"
	  mandatory = false
	  metachars_on_parameter_value_check = false
	  parameter_location = "any"
	  perform_staging = true
	  sensitive_parameter = false
	  signature_overrides_disable = [200002290, 200002292]
	  url {
		name = "test.com"
		method = "GET"
		protocol = "HTTP"
		type = "explicit"
	  }
	}	
	`, address)
}

func TestAccBigipWAFEntityParametersUnit(t *testing.T) {
	setup()
	defer teardown()

	js := `{"name":"test_entity_param","type":"explicit","valueType":"user-input","allowRepeatedParameterName":true,"attackSignaturesCheck":true,"dataType":"alpha_numeric","level":"url","parameterLocation":"any","performStaging":true,"signatureOverrides":[{"enabled":false,"signatureId":200002290},{"enabled":false,"signatureId":200002292}],"url":{"method":"GET","name":"test.com","protocol":"HTTP","type":"explicit"}}`

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: dataWAFEntityParametersCfg(server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.bigip_waf_entity_parameter.test_ep", "value_type", "user-input"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_parameter.test_ep", "url.0.method", "GET"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_parameter.test_ep", "url.0.name", "test.com"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_parameter.test_ep", "url.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_parameter.test_ep", "url.0.type", "explicit"),
					resource.TestCheckResourceAttr("data.bigip_waf_entity_parameter.test_ep", "json", js),
				),
			},
		},
	})
}
