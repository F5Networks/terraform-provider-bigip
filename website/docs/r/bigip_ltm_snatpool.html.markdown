---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_snatpool"
sidebar_current: "docs-bigip-resource-snatpool-x"
description: |-
    Provides details about bigip_ltm_snatpool resource
---

# bigip\_ltm\_snatpool

`bigip_ltm_snatpool` Collections of SNAT translation addresses

 


## Example Usage


```hcl
 resource "bigip_ltm_snatpoolpool" "snatpool_sanjose" {
  name = "/Common/snatpool_sanjose"
  members = ["191.1.1.1","194.2.2.2"]
}


```      

## Argument Reference

* `name` - (Required) Name of the snatpool

* ` members` - (Optional) Specifies a translation address to add to or delete from a SNAT pool.
