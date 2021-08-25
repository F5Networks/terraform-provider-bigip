---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_node"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_node data source
---

# bigip\_ltm\_node

Use this data source (`bigip_ltm_node`) to get the ltm node details available on BIG-IP


## Example Usage

```hcl


data "bigip_ltm_node" "test" {
  name      = "terraform_node"
  partition = "Common"
}


output "bigip_node" {
  value = "${data.bigip_ltm_node.test.address}"
}


# if it is fqdn address we can get fqdn elements as below

output "bigip_node" {
  value = "${data.bigip_ltm_node.test.fqdn[0].address_family}"
}

```

## Argument Reference

* `name` - (Required) Name of the node.
* `partition` - (Required) partition of the node.


## Attributes Reference

Additionally, the following attributes are exported:
* `address` - The address of the node.

* `description` - User defined description of the node.

* `connection_limit` - Node connection limit.

* `dynamic_ratio` - The dynamic ratio number for the node.

* `full_path` - Full path of the node (partition and name)

* `monitor` - Specifies the health monitors the system currently uses to monitor this node.

* `rate_limit` - Node rate limit.

* `ratio` - Node ratio weight.

* `state` - The current state of the node.

The `fqdn` block contains:

* `address_family` - The FQDN node's address family.
* `interval` - The amount of time before sending the next DNS query.
* `downinterval` - The number of attempts to resolve a domain name.
* `autopopulate` - Specifies if the node should scale to the IP address set returned by DNS.

