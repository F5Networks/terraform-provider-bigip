---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_policy"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_policy resource
---

# bigip\_ltm\_rewrite\_profile

`bigip_ltm_rewrite_profile` Configures ltm policies to manage traffic assigned to a virtual server

For resources should be named with their `full path`. The full path is the combination of the `partition + name` of the resource. For example `/Common/test-profile`.

## Example Usage

```hcl

resource "bigip_ltm_profile_rewrite" "test-profile" {
  name = "/Common/tf_profile"
  defaults_from = "/Common/rewrite"
  bypass_list = ["http://notouch.com"]
  rewrite_list = ["http://some.com"]
  rewrite_mode = "portal"
  cache_type = "cache-img-css-js"
  ca_file = "/Common/ca-bundle.crt"
  crl_file = "none"
  signing_cert = "/Common/default.crt"
  signing_key = "/Common/default.key"
  split_tunneling = "false"

  request {
    insert_xfwd_for = "enabled"
    insert_xfwd_host = "disabled"
    insert_xfwd_protocol = "enabled"
    rewrite_headers = "disabled"
  }

  response {
    rewrite_content = "enabled"
    rewrite_headers = "disabled"
  }
}
```

## Argument Reference

> [!NOTE]
> The attribute `published_copy` is not required anymore as the resource automatically publishes the policy, hence it's deprecated and will be removed from future release.

* `name`- (Required) Name of the Policy ( policy name should be in full path which is combination of partition and policy name )

* `defaults_from` - (optional,type `string`) Specifies the profile from which this profile inherits settings. The default is the system-supplied `rewrite` profile.

* `bypass_list` - (Optional,type `list`) Specifies a list of URIs to bypass inside a web page when the page is accessed using Portal Access.

* `cache_type` - (Required,type `string`) Specifies the type of Client caching. Valid choices are: `cache-css-js, cache-all, no-cache, cache-img-css-js`

* `ca_file` - (Optional) Specifies a CA against which to verify signed Java applets signatures. (name should be in full path which is combination of partition and CA file name )

*  `crl_file` - (Optional) Specifies a CRL against which to verify signed Java applets signature certificates. The default option is `none`.

* `signing_cert` - (Optional) Specifies a certificate to use for re-signing of signed Java applets after patching. (name should be in full path which is combination of partition and certificate name )

* `signing_key` - (Optional) Specifies a certificate to use for re-signing of signed Java applets after patching. (name should be in full path which is combination of partition and key name )

* `signing_key_password` - (Optional) Specifies a pass phrase to use for encrypting the private signing key.

* `split_tunneling` - (Optional,type `string`) Specifies the type of Client caching. Valid choices are: `true, false`

* `rewrite_list` - (Optional,type `list`) Specifies a list of URIs to rewrite inside a web page when the page is accessed using Portal Access.

* `rewrite_mode` - (Required,type `string`) Specifies the type of Client caching. Valid choices are: `portal, uri-translation`

* `request` - (Optional,type `set`) Block type. Each request is block type with following arguments.
    * `insert_xfwd_for` -  (Optional,type `string`) Valid choices are: `enabled, disabled`
    * `insert_xfwd_host` - (Optional,type `string`) Valid choices are: `enabled, disabled`
    * `insert_xfwd_protocol` - (Optional,type `string`) Valid choices are: `enabled, disabled`
    * `rewrite_headers` - (Optional,type `string`)  Valid choices are: `enabled, disabled`
* 
* `request` - (Optional,type `set`) Block type. Each request is block type with following arguments.
  * `rewrite_content` -  (Optional,type `string`) Valid choices are: `enabled, disabled`
  * `rewrite_headers` - (Optional,type `string`)  Valid choices are: `enabled, disabled`

## Importing
An existing profile can be imported into this resource by supplying profile Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_profile_rewrite.profile-import-test /Common/tf_profile
```
