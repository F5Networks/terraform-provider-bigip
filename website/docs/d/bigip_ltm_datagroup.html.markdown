---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_datagroup"
sidebar_current: "docs-bigip-datasource-datagroup-x"
description: |-
    Provides details about bigip_ltm_datagroup resource
---

# bigip\_ltm\_datagroup

`bigip_ltm_datagroup` Manages a datagroup configuration

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-datagroup.


## Example Usage


```hcl
resource "bigip_ltm_datagroup" "datagroup" {
  name = "dgx2"
  type = "string"
  records  {
   name = "abc.com"
   data = "pool1"
   }
}

```      

## Argument Reference

* `name` - (Required) Name of the datagroup

* `type` -  datagroup is string or address  Format

* `records` - Collections of Data and Value

* `name` - Data Value, this can be URL like www.abc.com

* `data` - This data can be a pool or virtual server etc.
