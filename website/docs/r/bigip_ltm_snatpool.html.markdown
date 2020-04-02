---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_snatpool"
sidebar_current: "docs-bigip-resource-snatpool-x"
description: |-
    Provides details about bigip_ltm_snatpool resource
---

# bigip\_ltm\_snatpool

`bigip_ltm_snatpool` Collections of SNAT translation addresses

Resource should be named with their "full path". The full path is the combination of the partition + name of the resource, for example /Common/my-snatpool. 


## Example Usage


```hcl

resource "bigip_ltm_snatpool" "snatpool_sanjose" {
  name    = "/Common/snatpool_sanjose"
  members = ["191.1.1.1", "194.2.2.2"]
}

```      

## Argument Reference

* `name` - (Required) Name of the snatpool

* `members` - (Required) Specifies a translation address to add to or delete from a SNAT pool (at least one address is required)
