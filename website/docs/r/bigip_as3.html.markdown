---
layout: "bigip"
page_title: "BIG-IP: bigip_as3"
sidebar_current: "docs-bigip-resource-x"
description: |-
    Provides details about bigip as3 resource
---

# bigip_as3

`bigip_as3` provides details about bigip as3 resource

This resource is helpful to configure deploy as3 declarative JSON on BIG-IP.
## Example Usage


```hcl

resource "bigip_as3"  "as3-example1" {
     as3_json = "${file("example1.json")}" 
     tenant_name = "as3"
 }

```  

## Argument Reference


* `as3_json` - (Required) Name of the of the Declarative AS3 JSON file

* `tenant_name` - (Required) This is the partition name where the application services will be configured.    

For more information on as3 please refer https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/#
 
