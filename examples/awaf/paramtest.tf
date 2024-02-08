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