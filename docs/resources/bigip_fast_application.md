---
layout: "bigip"
page_title: "BIG-IP: bigip_fast_application"
subcategory: "F5 Automation Tool Chain(ATC)"
description: |-
  Provides details about bigip_fast_application resource
---

# bigip_fast_application

`bigip_fast_application` This resource will create and manage FAST applications on BIG-IP from provided JSON declaration. 


## Example Usage


```hcl

resource "bigip_fast_application" "foo-app" {
  template  = "examples/simple_http"
  fast_json = "${file("new_fast_app.json")}"
}

```      

## Argument Reference


* `fast_json` - (Required) Path/Filename of Declarative FAST JSON which is a json file used with builtin ```file``` function
* `template` - (Optional) Name of installed FAST template used to create FAST application. This parameter is required when creating new resource.
* `tenant` - (Optional) A FAST tenant name on which you want to manage application.
* `application` - (Optional) A FAST application name.



* `FAST documentation` - https://clouddocs.f5.com/products/extensions/f5-appsvcs-templates/latest/
