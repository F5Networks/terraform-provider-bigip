---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_virtual_address"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_virtual_address resource
---

# bigip\_ltm\_virtual\_address

`bigip_ltm_virtual_address` Configures Virtual Server

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/virtual_server.


## Example Usage


```hcl
resource "bigip_ltm_virtual_address" "vs_va" {
  name            = "/Common/xxxxx"
  advertize_route = "enabled"
}

```      

## Argument Reference

* `name` - (Required) Name of the virtual address

* `description` - (Optional) Description of the virtual address

* `advertize_route` - (Optional) Enabled dynamic routing of the address ( In versions prior to BIG-IP 13.0.0 HF1, you can configure the Route Advertisement option for a virtual address to be either Enabled or Disabled only. Beginning with BIG-IP 13.0.0 HF1, F5 added more settings for the Route Advertisement option. In addition, the Enabled setting is deprecated and replaced by the Selective setting. For more information, please look into KB article https://support.f5.com/csp/article/K85543242 )

* `conn_limit` - (Optional, Default=0) Max number of connections for virtual address

* `enabled` - (Optional, Default=true) Enable or disable the virtual address

* `arp` - (Optional, Default=true) Enable or disable ARP for the virtual address

* `auto_delete` - (Optional, Default=true) Automatically delete the virtual address with the virtual server

* `icmp_echo` - (Optional, Default=enabled) Specifies how the system sends responses to ICMP echo requests on a per-virtual address basis.

* `traffic_group` - (Optional, Default=/Common/traffic-group-1) Specify the partition and traffic group
