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
  cross_origin_requests_enforcement {
    include_subdomains = true
    origin_name        = "app1.com"
    origin_port        = "80"
    origin_protocol    = "http"
  }
  cross_origin_requests_enforcement {
    include_subdomains = true
    origin_name        = "app2.com"
    origin_port        = "443"
    origin_protocol    = "http"
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
* `cross_origin_requests_enforcement` - (Optional) A list of options that enables your web-application to share data with a website hosted on a
different domain.
  * `include_subdomains` - (Required) Determines whether the subdomains are allowed to receive data from the web application.
  * `origin_name` - (Required) Specifies the name of the origin with which you want to share your data.
  * `origin_port` - (Required) Specifies the port that other web applications are allowed to use to request data from your web application.
  * `origin_protocol` - (Required) Specifies the protocol that other web applications are allowed to use to request data from your web application.


## Attributes Reference

* `json` - Json string representing created WAF entity URL declaration in JSON format

