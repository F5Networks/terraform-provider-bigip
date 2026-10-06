---
layout: "bigip"
page_title: "BIG-IP: bigip_net_vlan"
subcategory: "Network"
description: |-
  Provides details about bigip_net_vlan data source
---

# bigip\_net\_vlan

Use this data source (`bigip_net_vlan`) to look up a single VLAN already configured on a BIG-IP system, identified by either its name or its VLAN ID (tag), and retrieve its tagged/untagged interface membership. This is useful for migration discovery -- for example, replacing a `bigip_command` invocation that parsed `tmsh list net vlan` text output with a first-class data source.

Use [`bigip_net_vlans`](./bigip_net_vlans.md) instead if you need to enumerate every VLAN on the system.

## Example Usage

```hcl
# Look up by name
data "bigip_net_vlan" "by_name" {
  name = "/Common/external"
}

# Look up by VLAN ID (tag)
data "bigip_net_vlan" "by_id" {
  vlan_id = 100
}

output "tagged_interfaces" {
  value = data.bigip_net_vlan.by_name.tagged_interfaces
}

output "untagged_interfaces" {
  value = data.bigip_net_vlan.by_name.untagged_interfaces
}
```

## Argument Reference

Exactly one of `name` or `vlan_id` must be specified.

* `name` - (Optional) Name of the VLAN to look up. May be a bare name (combined with `partition`) or a full path (e.g. `/Common/external`).

* `partition` - (Optional) Partition of the VLAN when looking up by `name` as a bare name. Ignored when `name` is already a full path or when looking up by `vlan_id`. Defaults to `Common`.

* `vlan_id` - (Optional) VLAN ID (tag) to look up.

## Attributes Reference

Additionally, the following attributes are exported:

* `full_path` - Full path of the VLAN (partition and name).

* `tagged_interfaces` - Names of the interfaces or trunks attached to the VLAN as tagged members.

* `untagged_interfaces` - Names of the interfaces or trunks attached to the VLAN as untagged members.
