---
layout: "bigip"
page_title: "BIG-IP: bigip_as3_tenant"
sidebar_current: "docs-bigip-datasource-tenant-x"
description: |-
   Provides details about bigip_as3_tenant datasource
---
 
# bigip\_as3\_tenant
 
`bigip_as3_tenant` Manages a Tenant class, which is the highest level class. It becomes a partition on the BIG-IP.

 Each tenant comprises a set of Applications that belong to one authority (system role). In the following example, Sample_01 is the name of the tenant 
## Example Usage
 
 
```hcl
data "bigip_as3_tenant" "sample"{
  name = "Sample_01"
  app_class_list = ["${data.bigip_as3_app.App1.id}","${data.bigip_as3_app.App2.id}"]
  defaultroutedomain = 0
  enable = true
  label = "this is label for tenant"
  optimisticlockkey = "your key goes here"
  remark = "your remark goes here"
}
```
 
## Argument Reference
 
* `name` - (Required) Name of the tenant
 
* `defaultroutedomain` - (Optional) Selects the default route domain for IP traffic to and from this Tenant’s application resources 
 
* `enable` - (Optional) Tenant handles traffic only when enabled (default)
 
* `label` - (Optional) Optional friendly name for this object
 
* `remark` - (Optional) Arbitrary (brief) text pertaining to this object
 
* `optimisticlockkey` - (Optional) When you deploy a declaration with a non-empty ‘key’ value here, that activates an optimistic lock on changes to this Tenant
 
* `app_class_list`- (Optional) Specifies the list of applications present in the tenant 
