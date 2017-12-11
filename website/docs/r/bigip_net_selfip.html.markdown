---
layout: "bigip"
page_title: "BIG-IP: bigip_net_selfip"
sidebar_current: "docs-bigip-resource-selfip-x"
description: |-
    Provides details about bigip_net_selfip resource
---

# bigip\_net\_selfip

`bigip_net_selfip` Manages a selfip configuration

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
resource "bigip_net_selfip" "selfip1" {
	        name = "/Common/internalselfIP"
	        ip = "11.1.1.1/24"
	        vlan = "/Common/internal"
	        depends_on = ["bigip_net_vlan.vlan1"]
	}

```      

## Argument Reference

* `name` - (Required) Name of the selfip

* `ip` - (Optional) The Self IP's address and netmask.

* `vlan` - Specifies the VLAN for which you are setting a self IP address. This setting must be provided when a self IP is created.

* `depends_on` - (Optional) You need to provide this if you have not created the vlan.
