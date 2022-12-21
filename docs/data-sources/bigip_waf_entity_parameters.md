---
layout: "bigip"
page_title: "BIG-IP: bigip_waf entity parameters"
subcategory: "Web Application Firewall(WAF)"
description: |-
  Provides details entity parameters associated with a policy
---

# bigip\_waf\_entity_parameter

Use this data source (`bigip_waf_entity_parameter`) to get the details of entity parameters associated with a policy
 
 
## Example Usage

```hcl

data "bigip_waf_entity_parameter" "EPX" {
  name                        = "testParamX"
  type                        = "explicit"
  data_type                   = "alpha-numeric"
  enable_regular_expression   = true
  perform_staging             = true
  signature_overrides_disable = [200001494, 200001472]
}

```      

## Argument Reference

* `name` - (Required) Name of the Entity Parameter.
* `description` - Description of the Entity Parameter.
* `type` - Specifies whether the parameter is an explicit or a wildcard attribute.
* `value_type` - Specify the valid type for the value of the attribute.
* `allow_empty_type` - Determines whether an empty value is allowed for a parameter.
* `allow_repeated_parameter_name` - Determines whether multiple parameter instances with the same name are allowed in one request.
* `attack_signatures_check` - Determines whether attack signatures and threat campaigns must be detected in a parameter's value.
* `check_max_value_length` - Determines whether a parameter has a restricted maximum length for value.
* `check_min_value_length` - Determines whether a parameter has a restricted minimum length for value.
* `data_type` - Specifies data type of parameter's value.
* `enable_regular_expression` - Determines whether the parameter value includes the pattern defined in regularExpression.
* `is_base64` - Determines whether a parameter’s value contains a Base64 encoded string.
* `is_cookie` - Determines whether a parameter is located in the value of Cookie header.
* `is_header` - Determines whether a parameter is located in headers as one of the headers.
* `level` - Specifies whether the parameter is associated with a URL, a flow, or neither.
* `mandatory` - Determines whether a parameter must exist in the request.
* `metachars_on_parameter_value_check` - Determines whether disallowed metacharacters must be detected in a parameter’s value.
* `parameter_location` - Specifies location of parameter in request.
* `perform_staging` - Determines the staging state of a parameter.
* `sensitive_parameter` - Determines whether a parameter is sensitive and must be not visible in logs nor in the user interface.
* `signature_overrides_disable` - List of Attack Signature Ids which are disabled for this particular parameter.
* `url` - `url` block will provide options to be used for binding urls to parameter entity.See [url](#url) below for more details.

### url
The `url` block supports the following:

* `name` - (Required , `string`) name of url block attribute
* `method` - (Required , `string`) Unique ID of a URL with a protocol type and name. Select a Method for the URL to create an API endpoint: URL + Method
* `protocol` - (Required , `string`) Specifies whether the protocol for the URL is HTTP or HTTPS.The available options are : ["http","https"]
* `type` - (Required , `string`) Determines the type of the name attribute. Only when setting the type to wildcard will the special wildcard characters in the name be interpreted as such. The available types are : ["explicit","wildcard"]



## Attributes Reference

* `json` - JSON string representing the WAF Entity Parameter declaration.

