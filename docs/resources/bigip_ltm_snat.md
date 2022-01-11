---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_snat"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_snat resource
---

# bigip\_ltm\_snat

`bigip_ltm_snat` Manages a SNAT configuration

For resources should be named with their `full path`. The full path is the combination of the `partition + name` of the resource.For example `/Common/test-snat`.

## Example Usage

```hcl
resource "bigip_ltm_snat" "test-snat" {
  name        = "/Common/test-snat"
  translation = "/Common/136.1.1.2"
  sourceport  = "preserve"
  origins {
    name = "0.0.0.0/0"
  }
  vlans = [
    "/Common/internal",
  ]
  vlansdisabled = false
}
```      

## Argument Reference

* `name` - (Required) Name of the SNAT, name of SNAT should be full path. Full path is the combination of the `partition + SNAT name`,For example `/Common/test-snat`.

* `origins` - (Required) Specifies, for each SNAT that you create, the origin addresses that are to be members of that SNAT. Specify origin addresses by their IP addresses and service ports

* `translation` - (Optional) Specifies the IP address configured for translation. Note that translated addresses are outside the traffic management system. You can only use this option when `automap` and `snatpool` are not used.

* `snatpool` - (Optional) Specifies the name of a SNAT pool. You can only use this option when `automap` and `translation` are not used.

* `mirror` - (Optional) Enables or disables mirroring of SNAT connections.

* `autolasthop` -(Optional) Specifies whether to automatically map last hop for pools or not. The default is to use next level's default.

* `sourceport` - (Optional) Specifies how the SNAT object handles the client's source port. The default is `preserve`.

* `vlansdisabled` - (Optional,bool) Specifies the VLANs or tunnels for which the SNAT is enabled or disabled. The default is `true`, vlandisabled on VLANS specified by `vlans`,if set to `false` vlanEnabled set on VLANS specified by `vlans` .

* `vlans` - (Optional) Specifies the available VLANs or tunnels and those for which the SNAT is enabled or disabled.
