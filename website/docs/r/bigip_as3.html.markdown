---
layout: "bigip"
page_title: "BIG-IP: bigip_as3"
sidebar_current: "docs-bigip-resource-device-x"
description: |-
   Provides details about bigip as3[application service extension 3]
---

# bigip_as3

`bigip_as3` Configures BIG-IP using AS3 Json.

For more info on AS3, please refer [AS3](https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/) 

## Example Usage

```hcl

resource "bigip_as3" "as3_example"
        {
            tenant_name = "as3"
            as3_json = "${file("example.json")}"
        }
```     

## Argument Reference

* `as3_json` - (Required) Name of AS3 Json file/file path,Reference for [AS3_Example](https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/userguide/quick-start.html#quick-start-example-declaration)

* `tenant_name` - (Optional) Arbitary Name used for Unique Identifier of terraform resource,Default value `"as3"`

