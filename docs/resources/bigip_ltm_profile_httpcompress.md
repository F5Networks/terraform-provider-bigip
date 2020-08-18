---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_httpcompress"
sidebar_current: "docs-bigip-resource-profile_httpcompress-x"
description: |-
    Provides details about bigip_ltm_profile_httpcompress resource
---

# bigip\_ltm\_profile_httpcompress

`bigip_ltm_profile_httpcompress`  Virtual server HTTP compression profile configuration


For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

## Example Usage


```hcl
 resource "bigip_ltm_profile_httpcompress" "sjhttpcompression" {
   name                 = "/Common/sjhttpcompression2"
   defaults_from        = "/Common/httpcompression"
   uri_exclude          = ["www.abc.f5.com", "www.abc2.f5.com"]
   uri_include          = ["www.xyzbc.cisco.com"]
   content_type_include = ["nicecontent.com"]
   content_type_exclude = ["nicecontentexclude.com"]
 }

```      

## Argument Reference

* `name` (Required) Name of the profile_httpcompress

* `defaults_from` - (Optional) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `uri_exclude`  - (Optional) Disables compression on a specified list of HTTP Request-URI responses. Use a regular expression to specify a list of URIs you do not want to compress.

* `uri_include`  - (Optional) Enables compression on a specified list of HTTP Request-URI responses. Use a regular expression to specify a list of URIs you want to compress.

* `content_type_include` - (Optional) Specifies a list of content types for compression of HTTP Content-Type responses. Use a string list to specify a list of content types you want to compress.

* `content_type_exclude` - (Optional) Excludes a specified list of content types from compression of HTTP Content-Type responses. Use a string list to specify a list of content types you want to compress.
