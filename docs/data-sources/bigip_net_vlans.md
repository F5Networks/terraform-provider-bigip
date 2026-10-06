---
layout: "bigip"
page_title: "BIG-IP: bigip_net_vlans"
subcategory: "Network"
description: |-
  Provides details about bigip_net_vlans data source
---

# bigip\_net\_vlans

Use this data source (`bigip_net_vlans`) to retrieve a list of every VLAN currently configured on a BIG-IP system, along with its tagged/untagged interface or trunk membership. This is useful for discovering the existing L2 configuration of a device (for example, to recreate an equivalent configuration on another platform).

## Example Usage

```hcl
data "bigip_net_vlans" "all" {}

output "vlans" {
  value = data.bigip_net_vlans.all.vlans
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `vlans` - List of every VLAN configured on the BIG-IP system. Each entry exports:
  * `name` - Name of the VLAN.
  * `partition` - Partition the VLAN belongs to.
  * `full_path` - Full path of the VLAN (partition and name).
  * `tag` - VLAN ID (tag).
  * `mtu` - Maximum Transmission Unit (MTU) for the VLAN.
  * `cmp_hash` - Specifies how traffic on the VLAN is disaggregated.
  * `auto_lasthop` - Specifies whether auto lasthop is enabled or disabled on the VLAN.
  * `failsafe` - Whether failsafe is enabled on the VLAN.
  * `failsafe_action` - Action to take when failsafe is triggered.
  * `failsafe_timeout` - Failsafe timeout, in seconds.
  * `learning` - Specifies the VLAN's learning mode.
  * `source_checking` - Specifies whether source checking is enabled on the VLAN.
  * `interfaces` - Interfaces (or trunks) attached to the VLAN. Each entry exports:
    * `name` - Name of the interface or trunk attached to the VLAN.
    * `tagged` - Whether the interface is tagged on this VLAN.
