---
layout: "bigip"
page_title: "BIG-IP: bigip_net_selfip"
subcategory: "Network"
description: |-
  Provides details about bigip_net_selfip resource
---

# bigip\_net\_selfip

`bigip_net_selfip` Manages a selfip configuration

Resource should be named with their `full path`. The full path is the combination of the `partition + name of the resource`, for example `/Common/my-selfip`.


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
resource "bigip_net_selfip" "selfip1" {
  name       = "/Common/internalselfIP"
  ip         = "11.1.1.1/24"
  vlan       = "/Common/internal"
  depends_on = [bigip_net_vlan.vlan1]
}
```
### Example usage with `port_lockdown`

```hcl
resource "bigip_net_selfip" "selfip1" {
  name          = "/Common/internalselfIP"
  ip            = "11.1.1.1/24"
  vlan          = "/Common/internal"
  traffic_group = "traffic-group-1"
  port_lockdown = ["tcp:4040", "udp:5050", "egp:0"]
  depends_on    = [bigip_net_vlan.vlan1]
}
```

### Example usage with `port_lockdown` set to `["none"]`

```hcl
resource "bigip_net_selfip" "selfip1" {
  name          = "/Common/internalselfIP"
  ip            = "11.1.1.1/24"
  vlan          = "/Common/internal"
  traffic_group = "traffic-group-1"
  port_lockdown = ["none"]
  depends_on    = [bigip_net_vlan.vlan1]
}
```

## Argument Reference

* `name` - (Required) Name of the selfip

* `ip` - (Required) The Self IP's address and netmask.

* `vlan` - (Required) Specifies the VLAN for which you are setting a self IP address. This setting must be provided when a self IP is created.

* `traffic_group` - (Optional) Specifies the traffic group, defaults to `traffic-group-local-only` if not specified.

* `port_lockdown` - (Optional) Specifies the port lockdown, defaults to `Allow None` if not specified.
