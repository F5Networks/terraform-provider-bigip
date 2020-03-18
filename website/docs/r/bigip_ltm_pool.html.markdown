---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_pool"
sidebar_current: "docs-bigip-resource-pool-x"
description: |-
   Provides details about bigip_ltm_pool resource
---

# bigip\_ltm\_pool

`bigip_ltm_pool` Manages a pool configuration.

Resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
resource "bigip_ltm_pool" "pool" {
  name                = "/Common/terraform-pool"
  load_balancing_mode = "round-robin"
  description         = "Test-Pool"
  monitors            = [bigip_ltm_monitor.monitor.name, bigip_ltm_monitor.monitor2.name]
  allow_snat          = "yes"
  allow_nat           = "yes"
}
```      

## Argument Reference

* `name` - (Required) Name of the pool

* `monitors` - (Optional) List of monitor names to associate with the pool

* `description` - (Optional) Userdefined value to describe the pool 

* `allow_nat` - (Optional)

* `allow_snat` - (Optional)

* `load_balancing_mode` - (Optional, Default = round-robin)
