---
layout: "bigip"
page_title: "BIG-IP: bigip_net_interfaces"
subcategory: "Network"
description: |-
  Provides details about bigip_net_interfaces data source
---

# bigip\_net\_interfaces

Use this data source (`bigip_net_interfaces`) to retrieve a list of every physical interface present on a BIG-IP system, including administrative/operational status and negotiated speed. This is useful for discovering the existing physical interface layout of a device (for example, to map it to the equivalent interfaces on another platform).

~> **Note** TMOS interface names (e.g. `1.1`) do not map 1:1 to F5OS interface names -- rSeries (F5OS-A) uses `<port>.0` (e.g. `1.0`, no blade prefix), while VELOS (F5OS-C, chassis-based) uses `<blade>/<port>.<subport>` (e.g. `1/1.0`). See the [interface/trunk mapping guide](../guides/interface-trunk-mapping.html) for the full naming-convention differences.

## Example Usage

```hcl
data "bigip_net_interfaces" "all" {}

output "interfaces" {
  value = data.bigip_net_interfaces.all.interfaces
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

* `interfaces` - List of every physical interface present on the BIG-IP system. Each entry exports:
  * `name` - TMOS interface name, e.g. `1.1` (blade.port) or `mgmt`.
  * `full_path` - Full path of the interface.
  * `enabled` - Whether the interface is administratively enabled.
  * `status` - Operational link status from the interface's stats (`net/interface/stats`), e.g. `up`, `down`, `uninit`. This is retrieved via a single bulk stats request covering every interface, since it isn't present on the `net/interface` object itself.
  * `media_active` - Currently negotiated media/speed/duplex, e.g. `10000T-FD`, or `none` if the link is down/not negotiated.
  * `media_fixed` - Media type configured for the interface when not using auto-negotiation.
  * `media_max` - Maximum media/speed capability advertised for auto-negotiation.
  * `media_sfp` - Media setting for SFP/SFP+ transceiver slots, where applicable.
  * `mac_address` - MAC address of the interface.
  * `mtu` - Configured MTU of the interface.
  * `bundle` - Whether the interface supports bundling (trunk membership), e.g. `not-supported`, `enabled`, `disabled`.
  * `if_index` - SNMP ifIndex of the interface.
  * `flow_control` - Flow control setting, e.g. `tx-rx`, `none`.
  * `lldp_admin` - LLDP administrative mode for the interface, e.g. `txonly`, `disable`.
  * `stp` - Whether Spanning Tree Protocol is enabled on the interface.
  * `stp_link_type` - Spanning Tree link type for the interface.
  * `prefer_port` - For combination ports, which physical media (sfp/copper) is preferred.
