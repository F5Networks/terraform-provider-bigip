data "bigip_waf_entity_parameter" "Param1" {
  name            = "Param1"
  type            = "explicit"
  data_type       = "alpha-numeric"
  perform_staging = true
}

data "bigip_waf_entity_parameter" "Param2" {
  name            = "Param2"
  type            = "explicit"
  data_type       = "alpha-numeric"
  perform_staging = true
}

data "bigip_waf_entity_parameter" "Param3" {
  name      = "Param3"
  type      = "explicit"
  data_type = "alpha-numeric"
  level     = "url"
  url {
    method   = "*"
    name     = "*"
    protocol = "https"
    type     = "wildcard"
  }
  perform_staging = true
}

data "bigip_waf_entity_url" "URL4" {
  name            = "/www.twitter.com"
  protocol        = "https"
  method          = "GET"
  type            = "explicit"
  perform_staging = true
}

data "bigip_waf_entity_parameter" "Param4" {
  name      = "Param4"
  type      = "explicit"
  data_type = "alpha-numeric"
  level     = "url"
  url {
    method   = "GET"
    name     = "/www.twitter.com"
    protocol = "https"
    type     = "explicit"
  }
  perform_staging = true
}

resource "bigip_waf_policy" "test-awaf" {
  name                 = "mytestpolicy"
  partition            = "Common"
  template_name        = "POLICY_TEMPLATE_API_SECURITY"
  application_language = "utf-8"
  enforcement_mode     = "blocking"
  description          = "Rapid Deployment-2"
  urls                 = [data.bigip_waf_entity_url.URL4.json]
  parameters           = [data.bigip_waf_entity_parameter.Param1.json, data.bigip_waf_entity_parameter.Param2.json, data.bigip_waf_entity_parameter.Param3.json, data.bigip_waf_entity_parameter.Param4.json]
  policy_builder {
    learning_mode = "disabled"
  }
  file_types {
    name    = "testfiletype"
    type    = "explicit"
    allowed = true
  }
  ip_exceptions {
    ip_address = "100.10.10.0"
    ip_mask    = "255.255.255.0"
    //    block_requests = true
  }
  ip_exceptions {
    ip_address = "100.20.10.0"
    ip_mask    = "255.255.255.0"
    //    block_requests = true
  }
  graphql_profiles {
    name                    = "test_graphql"
    attack_signatures_check = true
    defense_attributes {
      maximum_total_length = 100000
      maximum_value_length = 10000
    }
  }
  //  open_api_files = ["file:////ts/var/rest/1676207728307_openapitestfile.json"]
  server_technologies = ["MySQL", "Unix/Linux", "MongoDB"]
}