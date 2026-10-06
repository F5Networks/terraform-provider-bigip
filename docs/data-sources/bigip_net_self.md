---
layout: "bigip"
page_title: "BIG-IP: bigip_net_self"
subcategory: "Network"
description: |-
  Provides details about bigip_net_self data source
---

# bigip\_net\_self

Use this data source (`bigip_net_self`) to look up a single self IP already configured on a BIG-IP system, identified by name, and retrieve its address, VLAN, traffic group, and allowed services. This is useful for migration discovery -- for example, replacing a `bigip_command` invocation that parsed `tmsh list net self` text output with a first-class data source.

Use [`bigip_net_selfips`](./bigip_net_selfips.md) instead if you need to enumerate every self IP on the system.

## Example Usage

```hcl
data "bigip_net_self" "external" {
  name = "/Common/external-self"
}

output "self_ip_address" {
  value = data.bigip_net_self.external.address
}

output "self_ip_allow_service" {
  value = data.bigip_net_self.external.allow_service
}
```

## Argument Reference

* `name` - (Required) Name of the self IP to look up. May be a bare name (combined with `partition`) or a full path (e.g. `/Common/external-self`).

* `partition` - (Optional) Partition of the self IP when `name` is a bare name. Ignored when `name` is already a full path. Defaults to `Common`.

## Attributes Reference

Additionally, the following attributes are exported:

* `full_path` - Full path of the self IP (partition and name).

* `address` - IP address (with subnet/route-domain suffix) of the self IP.

* `vlan` - VLAN the self IP is associated with.

* `traffic_group` - Traffic group the self IP belongs to.

* `floating` - Whether the self IP is a floating address.

* `unit` - Unit ID for non-floating self IPs in a redundant pair.

* `allow_service` - Port lockdown / allowed services for the self IP (e.g. `all`, `none`, or a list of service:port entries).
