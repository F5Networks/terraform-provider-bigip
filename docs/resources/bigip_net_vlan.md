---
layout: "bigip"
page_title: "BIG-IP: bigip_net_vlan"
subcategory: "Network"
description: |-
  Provides details about bigip_net_vlan resource
---

# bigip\_net\_vlan

`bigip_net_vlan` Manages a vlan configuration

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
resource "bigip_net_vlan" "vlan1" {
  name = "/Common/Internal"
  tag  = 101
  interfaces {
    vlanport = 1.2
    tagged   = false
  }
}


```      

## Argument Reference

* `name` - (Required) Name of the vlan

* `tag` - (Optional) Specifies a number that the system adds into the header of any frame passing through the VLAN.

* `interfaces` - (Optional) Specifies which interfaces you want this VLAN to use for traffic management.

* `vlanport` - Physical or virtual port used for traffic

* `cmp_hash` - (Optional,type `string`) Specifies how the traffic on the VLAN will be disaggregated. The value selected determines the traffic disaggregation method. possible options: [`default`, `src-ip`, `dst-ip`]

* `tagged` - Specifies a list of tagged interfaces or trunks associated with this VLAN. Note that you can associate tagged interfaces or trunks with any number of VLANs.
