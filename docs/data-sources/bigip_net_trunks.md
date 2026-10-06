---
layout: "bigip"
page_title: "BIG-IP: bigip_net_trunks"
subcategory: "Network"
description: |-
  Provides details about bigip_net_trunks data source
---

# bigip\_net\_trunks

Use this data source (`bigip_net_trunks`) to retrieve a list of every trunk (LAG) currently configured on a BIG-IP system, including member interfaces and LACP settings. This is useful for discovering the existing L2 configuration of a device (for example, to recreate an equivalent configuration on another platform).

~> **Note** Trunk (LAG) membership on TMOS does not map 1:1 to F5OS LAG configuration -- see the [interface/trunk mapping guide](../guides/interface-trunk-mapping.html) for the naming-convention differences between TMOS interface names (e.g. `1.1`) and F5OS interface names (rSeries/F5OS-A: `<port>.0`, e.g. `1.0`; VELOS/F5OS-C: `<blade>/<port>.<subport>`, e.g. `1/1.0`) that any migration mapping needs to account for.

## Example Usage

```hcl
data "bigip_net_trunks" "all" {}

output "trunks" {
  value = data.bigip_net_trunks.all.trunks
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `trunks` - List of every trunk (LAG) configured on the BIG-IP system. Each entry exports:
  * `name` - Name of the trunk.
  * `full_path` - Full path of the trunk.
  * `interfaces` - Physical interfaces that are members of the trunk.
  * `lacp` - Whether LACP is enabled on the trunk.
  * `lacp_mode` - LACP mode (active/passive).
  * `lacp_timeout` - LACP timeout (long/short).
  * `distribution_hash` - Load balancing distribution hash used across trunk members.
  * `link_select_policy` - Link selection policy for the trunk.
  * `bandwidth` - Configured bandwidth of the trunk.
  * `id` - Trunk ID.
  * `stp` - Whether Spanning Tree Protocol is enabled on the trunk.
  * `type` - Trunk type.
  * `working_member_count` - Number of currently-working trunk members.
