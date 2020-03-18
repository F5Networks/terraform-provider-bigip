---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_virtual_address"
sidebar_current: "docs-bigip-resource-virtual_address-x"
description: |-
    Provides details about bigip_ltm_virtual_address resource
---

# bigip\_ltm\_virtual\_address

`bigip_ltm_virtual_address` Configures Virtual Server

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
resource "bigip_ltm_virtual_address" "vs_va" {
  name            = "/Common/vs_va"
  advertize_route = true
}

```      

## Argument Reference

* `name` - (Required) Name of the virtual address

* `description` - (Optional) Description of the virtual address

* `advertize_route` - (Optional) Enabled dynamic routing of the address

* `conn_limit` - (Optional, Default=0) Max number of connections for virtual address

* `enabled` - (Optional, Default=true) Enable or disable the virtual address

* `arp` - (Optional, Default=true) Enable or disable ARP for the virtual address

* `auto_delete` - (Optional, Default=true) Automatically delete the virtual address with the virtual server

* `icmp_echo` - (Optional, Default=true) Enable/Disable ICMP response to the virtual address

* `traffic_group` - (Optional, Default=/Common/traffic-group-1) Specify the partition and traffic group
