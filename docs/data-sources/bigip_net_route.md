---
layout: "bigip"
page_title: "BIG-IP: bigip_net_route"
subcategory: "Network"
description: |-
  Provides details about bigip_net_route data source
---

# bigip\_net\_route

Use this data source (`bigip_net_route`) to look up a single static route already configured on a BIG-IP system, identified by name, and retrieve its destination network, gateway, and description. This is useful for migration discovery -- for example, replacing a `bigip_command` invocation that parsed `tmsh list net route` text output with a first-class data source.

Use [`bigip_net_routes`](./bigip_net_routes.md) instead if you need to enumerate every static route on the system.

## Example Usage

```hcl
data "bigip_net_route" "external" {
  name = "/Common/external-route"
}

output "route_network" {
  value = data.bigip_net_route.external.network
}

output "route_gateway" {
  value = data.bigip_net_route.external.gw
}
```

## Argument Reference

* `name` - (Required) Name of the route to look up. May be a bare name (combined with `partition`) or a full path (e.g. `/Common/external-route`).

* `partition` - (Optional) Partition of the route when `name` is a bare name. Ignored when `name` is already a full path. Defaults to `Common`.

## Attributes Reference

Additionally, the following attributes are exported:

* `full_path` - Full path of the route (partition and name).

* `network` - Destination network of the route.

* `gw` - Gateway address for the route.

* `description` - Description of the route.

* `tm_interface` - Tunnel/interface used to route traffic, if configured instead of a gateway.

* `mtu` - MTU configured for the route.

* `blackhole` - Whether the route is configured to reject/blackhole traffic.
