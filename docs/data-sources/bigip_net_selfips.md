---
layout: "bigip"
page_title: "BIG-IP: bigip_net_selfips"
subcategory: "Network"
description: |-
  Provides details about bigip_net_selfips data source
---

# bigip\_net\_selfips

Use this data source (`bigip_net_selfips`) to retrieve a list of every self IP currently configured on a BIG-IP system. This is useful for discovering the existing L3 configuration of a device (for example, to recreate an equivalent configuration on another platform).

## Example Usage

```hcl
data "bigip_net_selfips" "all" {}

output "self_ips" {
  value = data.bigip_net_selfips.all.self_ips
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `self_ips` - List of every self IP configured on the BIG-IP system. Each entry exports:
  * `name` - Name of the self IP.
  * `partition` - Partition the self IP belongs to.
  * `full_path` - Full path of the self IP (partition and name).
  * `address` - IP address (with subnet/route-domain suffix) of the self IP.
  * `vlan` - VLAN the self IP is associated with.
  * `traffic_group` - Traffic group the self IP belongs to.
  * `floating` - Whether the self IP is a floating address.
  * `unit` - Unit ID for non-floating self IPs in a redundant pair.
  * `allow_service` - Port lockdown / allowed services for the self IP (e.g. `all`, `none`, or a list of service:port entries).
