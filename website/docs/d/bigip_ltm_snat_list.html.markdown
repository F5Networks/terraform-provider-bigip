---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_snat"
sidebar_current: "docs-bigip-datasource-snat-x"
description: |-
    Provides details about bigip_ltm_snat resource
---

# bigip\_ltm\_snat

`bigip_ltm_snat` Manages a snat configuration

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
resource "bigip_ltm_snat" "snat3" {
  // this is using snatpool translation is not required
  name = "snat3"
  origins = ["6.1.6.6"]
  mirror = "false"
  snatpool = "/Common/sanjosesnatpool"
}
resource "bigip_ltm_snat" "snat_list" {
 name = "NewSnatList"
 translation = "136.1.1.1"
 origins = ["2.2.2.2", "3.3.3.3"]
}

```      

## Argument Reference

* `name` - (Required) Name of the snat

* `origins` - (Optional) IP or hostname of the snat

* `snatpool` - (Optional) Specifies the name of a SNAT pool. You can only use this option when automap and translation are not used.

* `mirror` - (Optional) Enables or disables mirroring of SNAT connections.

* `translation` - (Optional) Provide the translation IP address.
