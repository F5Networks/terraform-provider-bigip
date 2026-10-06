---
layout: "bigip"
page_title: "BIG-IP: bigip_net_routes"
subcategory: "Network"
description: |-
  Provides details about bigip_net_routes data source
---

# bigip\_net\_routes

Use this data source (`bigip_net_routes`) to retrieve a list of every static route currently configured on a BIG-IP system. This is useful for discovering the existing L3 configuration of a device (for example, to recreate an equivalent configuration on another platform).

## Example Usage

```hcl
data "bigip_net_routes" "all" {}

output "routes" {
  value = data.bigip_net_routes.all.routes
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `routes` - List of every static route configured on the BIG-IP system. Each entry exports:
  * `name` - Name of the route.
  * `partition` - Partition the route belongs to.
  * `full_path` - Full path of the route (partition and name).
  * `network` - Destination network of the route.
  * `gw` - Gateway address for the route.
  * `tm_interface` - Tunnel/interface used to route traffic, if configured instead of a gateway.
  * `mtu` - MTU configured for the route.
  * `blackhole` - Whether the route is configured to reject/blackhole traffic.
