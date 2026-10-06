---
layout: "bigip"
page_title: "BIG-IP: bigip_net_trunk"
subcategory: "Network"
description: |-
  Provides details about bigip_net_trunk data source
---

# bigip\_net\_trunk

Use this data source (`bigip_net_trunk`) to look up a single trunk (LAG) already configured on a BIG-IP system, identified by name, and retrieve its member interfaces and LACP settings. This is useful for LAG mapping during migration -- for example, resolving a specific TMOS trunk's membership before translating it to an F5OS LAG configuration.

Use [`bigip_net_trunks`](./bigip_net_trunks.md) instead if you need to enumerate every trunk on the system.

~> **Note** Trunk (LAG) membership on TMOS does not map 1:1 to F5OS LAG configuration -- see the [interface/trunk mapping guide](../guides/interface-trunk-mapping.html) for the naming-convention differences between TMOS interface names (e.g. `1.1`) and F5OS interface names (rSeries/F5OS-A: `<port>.0`, e.g. `1.0`; VELOS/F5OS-C: `<blade>/<port>.<subport>`, e.g. `1/1.0`) that any migration mapping needs to account for.

## Example Usage

```hcl
data "bigip_net_trunk" "lag1" {
  name = "lag1"
}

output "trunk_interfaces" {
  value = data.bigip_net_trunk.lag1.interfaces
}

output "trunk_lacp_mode" {
  value = data.bigip_net_trunk.lag1.lacp_mode
}
```

## Argument Reference

* `name` - (Required) Name of the trunk to look up.

## Attributes Reference

Additionally, the following attributes are exported:

* `interfaces` - Physical interfaces that are members of the trunk.

* `lacp` - Whether LACP is enabled on the trunk.

* `lacp_mode` - LACP mode (active/passive).

* `lacp_timeout` - LACP timeout (long/short).

* `distribution_hash` - Load balancing distribution hash used across trunk members.

* `link_select_policy` - Link selection policy for the trunk.

* `bandwidth` - Configured bandwidth of the trunk.

* `trunk_id` - Trunk ID.

* `stp` - Whether Spanning Tree Protocol is enabled on the trunk.

* `type` - Trunk type.

* `working_member_count` - Number of currently-working trunk members.
