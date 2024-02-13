data "bigip_waf_entity_parameter" "Param1" {
  name            = "Param1"
  type            = "explicit"
  data_type       = "alpha-numeric"
  check_max_value_length = true
  check_min_value_length = true
  max_value_length = 30
  min_value_length = 15
  perform_staging = true
}

data "bigip_waf_entity_parameter" "Param2" {
  name            = "Param2"
  type            = "explicit"
  data_type       = "alpha-numeric"
  check_max_value_length = true
  perform_staging = true
}

data "bigip_waf_entity_parameter" "Param3" {
  name            = "Param3"
  type            = "explicit"
  data_type       = "alpha-numeric"
  max_value_length = 30
  min_value_length = 15
  perform_staging = true
}

resource "bigip_waf_policy" "github925-awaf" {
  name                 = "github925-awaf"
  partition            = "Common"
  template_name        = "POLICY_TEMPLATE_API_SECURITY"
  application_language = "utf-8"
  enforcement_mode     = "blocking"
  description          = "Rapid Deployment-2"
  server_technologies  = ["MySQL", "Unix/Linux", "MongoDB"]
  parameters           = [data.bigip_waf_entity_parameter.Param1.json, data.bigip_waf_entity_parameter.Param2.json,data.bigip_waf_entity_parameter.Param3.json]
}
