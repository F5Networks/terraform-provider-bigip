---
layout: "bigip"
page_title: "BIG-IP: bigip_as3_as3class"
sidebar_current: "docs-bigip-datasource-as3-x"
description: |-
   Provides details about bigip_as3_as3class datasource
---

# bigip\_as3\_as3class

`bigip_as3_as3class` Manages an as3 class, that defines the top level objects

## Example Usage


```hcl
resource "bigip_as3_class" "as3-example" {
  name = "as3-example"
  declaration="${data.bigip_as3_adc.exmpadc.result_map}"
  tenants = ["Sample_01"]
}
```

## Argument Reference

* `name` - (Required) Name of the as3 class

* `declaration` - (Required) Pointer to the result map of adc datasource

* `tenants` - (Required) Specifies the list of tenants
