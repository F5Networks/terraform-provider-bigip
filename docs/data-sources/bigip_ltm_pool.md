---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_pool"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_pool data source
---

# bigip\_ltm\_monitor

Use this data source (`bigip_ltm_pool`) to get the ltm monitor details available on BIG-IP
 
 
## Example Usage
```hcl

data "bigip_ltm_pool" "Pool-Example" {
  name      = "example-pool"
  partition = "Common"
}

```      

## Argument Reference

* `name` - (Required) Name of the ltm monitor

* `partition` - (Required) partition of the ltm monitor


## Attributes Reference

Additionally, the following attributes are exported:

* `full_path` - Full path to the pool.
