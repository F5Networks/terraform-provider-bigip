---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_node"
sidebar_current: "docs-bigip-datasource-node-x"
description: |-
    Provides details about bigip_ltm_node data source
---

# bigip\_ltm\_node

Use this data source (`bigip_ltm_node`) to get the ltm node details available on BIG-IP
 
 
## Example Usage
```hcl


data "bigip_ltm_node" "test" {
  name = "terraform_node"
  address = "192.168.1.1"
}


output "bigip_node" {
  value = "${data.bigip_ltm_node.test.address}"
}

```      

## Argument Reference

* `name` - (Required) Name of the node

* `address` - (Required) address of the node


## Attributes Reference

Additionally, the following attributes are exported:

* `description` - User defined description of the node.

* `connection_limit` - Node connection limit.

* `dynamic_ratio` - The dynamic ratio number for the node.

* `monitor` - Specifies the health monitors the system currently uses to monitor this node.

* `rate_limit` - Node rate limit.

* `ratio` - Node ratio weight.

* `state` - The current state of the node.

The `fqdn` block contains:

* `address_family` - The FQDN node's address family.
* `name` - The fully qualified domain name of the node.
* `interval` - The amount of time before sending the next DNS query.
* `downinterval` - The number of attempts to resolve a domain name.
* `autopopulate` - Specifies if the node should scale to the IP address set returned by DNS.

