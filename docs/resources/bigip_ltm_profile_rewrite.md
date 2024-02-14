---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_rewrite"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_rewrite resource
---

# bigip\_ltm\_rewrite\_profile

`bigip_ltm_rewrite_profile` Configures ltm policies to manage traffic assigned to a virtual server

For resources should be named with their `full path`. The full path is the combination of the `partition + name` of the resource. For example `/Common/test-profile`.

## Example Usage

```hcl
resource "bigip_ltm_profile_rewrite" "test-profile" {
  name            = "/Common/tf_profile"
  defaults_from   = "/Common/rewrite"
  bypass_list     = ["http://notouch.com"]
  rewrite_list    = ["http://some.com"]
  rewrite_mode    = "portal"
  cache_type      = "cache-img-css-js"
  ca_file         = "/Common/ca-bundle.crt"
  crl_file        = "none"
  signing_cert    = "/Common/default.crt"
  signing_key     = "/Common/default.key"
  split_tunneling = "true"
}

resource "bigip_ltm_profile_rewrite" "test-profile2" {
  name          = "/Common/tf_profile_translate"
  defaults_from = "/Common/rewrite"
  rewrite_mode  = "uri-translation"
  request {
    insert_xfwd_for      = "enabled"
    insert_xfwd_host     = "disabled"
    insert_xfwd_protocol = "enabled"
    rewrite_headers      = "disabled"
  }
  response {
    rewrite_content = "enabled"
    rewrite_headers = "disabled"
  }
  cookie_rules {
    rule_name     = "cookie1"
    client_domain = "wrong.com"
    client_path   = "/this/"
    server_domain = "wrong.com"
    server_path   = "/this/"
  }
  cookie_rules {
    rule_name     = "cookie2"
    client_domain = "incorrect.com"
    client_path   = "/this/"
    server_domain = "absolute.com"
    server_path   = "/this/"
  }
}
```

## Argument Reference

* `name`- (Required) Name of the rewrite profile. ( profile name should be in full path which is combination of partition and profile name )

* `partition` - (optional,type `string`) Specifies the partition to create resource.

* `defaults_from` - (optional,type `string`) Specifies the profile from which this profile inherits settings. The default is the system-supplied `rewrite` profile.

* `bypass_list` - (Optional,type `list`) Specifies a list of URIs to bypass inside a web page when the page is accessed using Portal Access.

* `cache_type` - (Optional,type `string`) Specifies the type of Client caching. Valid choices are: `cache-css-js, cache-all, no-cache, cache-img-css-js`. Default value: `cache-img-css-js`

* `ca_file` - (Optional, type `string`) Specifies a CA against which to verify signed Java applets signatures. (name should be in full path which is combination of partition and CA file name )

* `crl_file` - (Optional, type `string`) Specifies a CRL against which to verify signed Java applets signature certificates. The default option is `none`.

* `signing_cert` - (Optional, type `string`) Specifies a certificate to use for re-signing of signed Java applets after patching. (name should be in full path which is combination of partition and certificate name )

* `signing_key` - (Optional, type `string`) Specifies a certificate to use for re-signing of signed Java applets after patching. (name should be in full path which is combination of partition and key name )

* `signing_key_password` - (Optional, type `string`) Specifies a pass phrase to use for encrypting the private signing key. Since it's a sensitive entity idempotency will fail in the update call.

* `split_tunneling` - (Optional,type `string`) Specifies the type of Client caching. Valid choices are: `true, false`

* `rewrite_list` - (Optional,type `list`) Specifies a list of URIs to rewrite inside a web page when the page is accessed using Portal Access.

* `rewrite_mode` - (Required,type `string`) Specifies the type of Client caching. Valid choices are: `portal, uri-translation`

* `request` - (Optional,type `set`) Block type. Each request is block type with following arguments.
    * `insert_xfwd_for` -  (Optional,type `string`) Enable to add the X-Forwarded For (XFF) header, to specify the originating IP address of the client. Valid choices are: `enabled, disabled`
    * `insert_xfwd_host` - (Optional,type `string`) Enable to add the X-Forwarded Host header, to specify the originating host of the client. Valid choices are: `enabled, disabled`
    * `insert_xfwd_protocol` - (Optional,type `string`) Enable to add the X-Forwarded Proto header, to specify the originating protocol of the client. Valid choices are: `enabled, disabled`
    * `rewrite_headers` - (Optional,type `string`) Enable to rewrite headers in Request settings. Valid choices are: `enabled, disabled`

* `response` - (Optional,type `set`) Block type. Each request is block type with following arguments.
    * `rewrite_content` -  (Optional,type `string`) Enable to rewrite links in content in the response. Valid choices are: `enabled, disabled`
    * `rewrite_headers` - (Optional,type `string`) Enable to rewrite headers in the response. Valid choices are: `enabled, disabled`

* `cookie_rules` - (Optional,type `set`) Specifies the cookie rewrite rules. Block type. Each request is block type with following arguments.
    * `rule_name` - (Required,type `string`) Name of the cookie rewrite rule.
    * `client_domain` - (Required,type `string`) 
    * `client_path` - (Required,type `string`)
    * `server_domain` - (Required,type `string`)
    * `server_path` - (Required,type `string`)
