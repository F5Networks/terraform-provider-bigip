---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_snat"
sidebar_current: "docs-bigip-resource-snat-x"
description: |-
    Provides details about bigip_ltm_snat resource
---

# bigip\_ltm\_snat

`bigip_ltm_snat` Manages a snat configuration

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
resource "bigip_ltm_snat" "test-snat" {
  name        = "TEST_SNAT_NAME"
  translation = "/Common/136.1.1.1"
  origins {
    name = "2.2.2.2"
  }
  origins {
    name = "3.3.3.3"
  }
  vlansdisabled = true
  autolasthop   = "default"
  mirror        = "disabled"
  partition     = "Common"
  full_path     = "/Common/test-snat"
}

```      

## Argument Reference

* `name` - (Required) Name of the snat

* `partition` - (Optional) Displays the administrative partition within which this profile resides

* `origins` - (Optional) IP or hostname of the snat

* `snatpool` - (Optional) Specifies the name of a SNAT pool. You can only use this option when automap and translation are not used.

* `mirror` - (Optional) Enables or disables mirroring of SNAT connections.

* `autolasthop` -(Optional) Specifies whether to automatically map last hop for pools or not. The default is to use next level's default.

* `sourceport` - (Optional) Specifies whether the system preserves the source port of the connection. The default is preserve. Use of the preserve-strict setting should be restricted to UDP only under very special circumstances such as nPath or transparent (that is, no translation of any other L3/L4 field), where there is a 1:1 relationship between virtual IP addresses and node addresses, or when clustered multi-processing (CMP) is disabled. The change setting is useful for obfuscating internal network addresses.

* `translation` - (Optional) Specifies the name of a translated IP address. Note that translated addresses are outside the traffic management system. You can only use this option when automap and snatpool are not used.

* `vlansdisabled` - (Optional) Disables the SNAT on all VLANs.

* `vlans` - (Optional) Specifies the name of the VLAN to which you want to assign the SNAT. The default is vlans-enabled.
