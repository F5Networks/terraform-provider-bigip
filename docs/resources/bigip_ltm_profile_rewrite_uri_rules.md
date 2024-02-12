---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_rewrite_uri_rules"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_rewrite_uri_rules resource
---

# bigip\_ltm\_rewrite\_profile_\uri\_rules

`bigip_ltm_profile_rewrite_uri_rules` Configures uri rewrite rules attached to the ltm rewrite profile

## Example Usage

```hcl
resource "bigip_ltm_profile_rewrite" "tftest" {
  name          = "/Common/tf_profile"
  defaults_from = "/Common/rewrite"
  rewrite_mode  = "uri-translation"
}

resource "bigip_ltm_profile_rewrite_uri_rules" "tftestrule1" {
  profile_name = bigip_ltm_profile_rewrite.tftest.name
  rule_name    = "tf_rule"
  rule_type    = "request"
  client {
    host   = "www.foo.com"
    scheme = "https"
  }
  server {
    host   = "www.bar.com"
    path   = "/this/"
    scheme = "https"
    port   = "8888"
  }
}

resource "bigip_ltm_profile_rewrite_uri_rules" "tftestrule2" {
  profile_name = bigip_ltm_profile_rewrite.tftest.name
  rule_name    = "tf_rule2"
  client {
    host   = "www.baz.com"
    path   = "/that/"
    scheme = "ftp"
    port   = "8888"
  }
  server {
    host   = "www.buz.com"
    path   = "/those/"
    scheme = "ftps"
  }
}
```

## Argument Reference

* `profile_name`- (Required, type `string`) Name of the rewrite profile. ( policy name should be in full path which is combination of partition and policy name )

* `rule_name` - (Required, type `string`) Specifies the name of the uri rule.

* `rule_type` - (Optional, type `string`) Specifies the type of the uri rule. Valid choices are: `request, response, both`. Default value is: `both`

* `client` - (Optional,type `set`) Block type. Each request is block type with following arguments.
    * `host` -  (Required,type `string`) Host part of the uri, e.g. `www.foo.com`.
    * `path` - (Optional,type `string`) Path part of the uri, must always end with `/`. Default value is: `/`
    * `scheme` - (Required,type `string`) Scheme part of the uri, e.g. `https`, `ftp`.
    * `port` - (Optional,type `string`) Port part of the uri. Default value is: `none`

* `server` - (Optional,type `set`) Block type. Each request is block type with following arguments.
  * `host` -  (Required,type `string`) Host part of the uri, e.g. `www.foo.com`.
  * `path` - (Optional,type `string`) Path part of the uri, must always end with `/`. Default value is: `/`
  * `scheme` - (Required,type `string`) Scheme part of the uri, e.g. `https`, `ftp`.
  * `port` - (Optional,type `string`) Port part of the uri. Default value is: `none`
