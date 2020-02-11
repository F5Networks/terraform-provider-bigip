---
layout: "bigip"
page_title: "BIG-IP: bigip_as3_adc"
sidebar_current: "docs-bigip-datasource-adc-x"
description: |-
   Provides details about bigip_as3_adc datasource
---
 
# bigip\_as3\_adc
 
`bigip_as3_adc` Manages an ADC class, that defines general settings for the declaration
 
## Example Usage
 
 
```hcl
data "bigip_as3_adc" "exmpadc"{
  name = "exmpadc"
  label = "your label goes here"
  tenant_class_list = ["${data.bigip_as3_tenant.sample.id}"]
}
```
 
## Argument Reference
 
* `name` - (Required) Name of the adc class
 
* `schema_version` - (Optional) When composing new declarations, you should use the latest schema version. This prevents inadvertently running a declaration on an outdated version of AS3 code.
 
* `tenant_class_list` - (Optional) Pointer to the list of ids of tenant datasources
 
* `label` - (Optional) This value can be anything less than 255 characters and simply labels the declaration.
