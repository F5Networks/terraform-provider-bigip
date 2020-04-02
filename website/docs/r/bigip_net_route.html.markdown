---
layout: "bigip"
page_title: "BIG-IP: bigip_net_route"
sidebar_current: "docs-bigip-resource-route-x"
description: |-
    Provides details about bigip_net_route resource
---

# bigip\_net\_route

`bigip_net_route` Manages a route configuration

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
resource "bigip_net_route" "route2" {
  name    = "external-route"
  network = "10.10.10.0/24"
  gw      = "1.1.1.2"
}

```      

## Argument Reference

* `name` - (Required) Name of the route

* `network` - (Optional) The destination subnet and netmask for the route.

* `network` - (Optional) Specifies a gateway address for the route.
