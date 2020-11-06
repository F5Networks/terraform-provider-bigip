---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_datagroup"
sidebar_current: "docs-bigip-datasource-datagroup-x"
description: |-
    Provides details about bigip_ltm_datagroup data source
---

# bigip\_ltm\_datagroup

Use this data source (`bigip_ltm_datagroup`) to get the data group details available on BIG-IP
 
 
## Example Usage
```hcl

data "bigip_ltm_datagroup" "DG-TC3" {
  name = "test-dg"
  partition = "Common"
}

```      

## Argument Reference

* `name` - (Required) Name of the datagroup

* `partition` - (Required) partition of the datagroup
