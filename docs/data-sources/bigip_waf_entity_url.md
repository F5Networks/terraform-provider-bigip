---
layout: "bigip"
page_title: "BIG-IP: bigip_waf entity url"
subcategory: "Web Application Firewall(WAF)"
description: |-
  Provides details and exports bigip_waf_entity_url data source
---

# bigip\_waf\_pb_suggestions

Use this data source (`bigip_waf_pb_suggestions`) to create JSON for WAF URL to later use with an existing WAF policy.


## Example Usage

```hcl

data "bigip_waf_entity_url" "WAFURL1" {
  name                        = "/foobar"
  description                 = "this is a test"
  type                        = "explicit"
  protocol                    = "HTTP"
  perform_staging             = true
  signature_overrides_disable = [12345678, 87654321]
  method_overrides {
    allow  = false
    method = "BCOPY"
  }
  method_overrides {
    allow  = true
    method = "BDELETE"
  }
}

```

## Argument Reference

* `name` - (Required) WAF entity URL name.
* `description` - (Optional) A description of the URL.
* `type` - (Optional) Specifies whether the parameter is an 'explicit' or a 'wildcard' attribute. Default is: wildcard.
* `protocol` - (Optional) Specifies whether the protocol for the URL is 'http' or 'https'. Default is: http.
* `method` - (Optional) Select a Method for the URL to create an API endpoint. Default is : *.
* `perform_staging` - (Optional) If true then any violation associated to the respective URL will not be enforced, and the request will not be considered illegal.
* `signature_overrides_disable` - (Optional) List of Attack Signature Ids which are disabled for this particular URL. 
* `method_overrides` - (Optional) A list of methods that are allowed or disallowed for a specific URL.
  * `allow` - (Required) Specifies that the system allows or disallows a method for this URL
  * `method` - (Required) Specifies an HTTP method.


## Attributes Reference

* `json` - Json string representing created WAF entity URL declaration in JSON format

