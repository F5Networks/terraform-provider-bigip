---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_node"
sidebar_current: "docs-bigip-resource-node-x"
description: |-
    Provides details about bigip_ltm_node resource
---

# bigip\_ltm\_node

`bigip_ltm_node` Manages a node configuration

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
resource "bigip_ltm_node" "node" {
  name = "/Common/terraform_node1"
  address = "10.10.10.10"
  connection_limit = "0"
	dynamic_ratio = "1"
	monitor = "default"
	rate_limit = "disabled"
	fqdn = { interval = "3000"}
}

```      

## Argument Reference

* `name` - (Required) Name of the node

* `address` - (Required) IP or hostname of the node

* `state` - (Optional) Default is "user-up" you can set to "user-down" if you want to disable

* `connection_limit` - (Optional) Specifies the maximum number of connections allowed for the node or node address, default is 0

* `monitor` - (Optional) Specifies the name of the monitor or monitor rule that you want to associate with the node.

* `dynamic_ratio` - (Optional) Sets the dynamic ratio number for the node. Used for dynamic ratio load balancing. The ratio weights are based on continuous monitoring of the servers and are therefore continually changing. Dynamic Ratio load balancing may currently be implemented on RealNetworks RealServer platforms, on Windows platforms equipped with Windows Management Instrumentation (WMI), or on a server equipped with either the UC Davis SNMP agent or Windows 2000 Server SNMP agent.

* `rate_limit` - (Optional) Specifies the maximum number of connections per second allowed for a node or node address. The default value is 'disabled'.

* `interval` - (Optional) Specifies the amount of time before sending the next DNS query. It can also take value as "ttl" when "ttl" is specified the interval option is disabled, vice versa.
