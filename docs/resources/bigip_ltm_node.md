---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_node"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_node resource
---

# bigip\_ltm\_node

`bigip_ltm_node` Manages a node configuration

For resources should be named with their `full path`.The full path is the combination of the `partition + name` of the resource( example: `/Common/my-node` ) or `partition + Direcroty + name` of the resource ( example: `/Common/test/my-node` ).
When including directory in `full path` we have to make sure it is created in the given partition before using it.


## Example Usage

```hcl
resource "bigip_ltm_node" "node" {
  name             = "/Common/terraform_node1"
  address          = "192.168.30.1"
  connection_limit = "0"
  dynamic_ratio    = "1"
  monitor          = "/Common/icmp"
  description      = "Test-Node"
  rate_limit       = "disabled"
  fqdn {
    address_family = "ipv4"
    interval       = "3000"
  }
}
```      

## Argument Reference

* `name` - (Required , type `string`) Name of the node

* `address` - (Required, type `string`) IP or hostname of the node

* `description` - (Optional,type `string`) User-defined description give ltm_node

* `connection_limit` - (Optional,type `int`) Specifies the maximum number of connections allowed for the node or node address.

* `dynamic_ratio` - (Optional, type `int`) Specifies the fixed ratio value used for a node during ratio load balancing.

* `monitor` - (Optional) specifies the name of the monitor or monitor rule that you want to associate with the node.

* `rate_limit`- (Optional,type `string`) Specifies the maximum number of connections per second allowed for a node or node address. The default value is 'disabled'.

* `state` - (Optional) Default is "user-up" you can set to "user-down" if you want to disable

 ~> *NOTE* Below attributes needs to be configured under fqdn option.

* `interval` - (Optional, type `string`) Specifies the amount of time before sending the next DNS query. Default is 3600. This needs to be specified inside the fqdn (fully qualified domain name).

* `address_family` - (Optional) Specifies the node's address family. The default is 'unspecified', or IP-agnostic. This needs to be specified inside the fqdn (fully qualified domain name).

## Importing
An existing Node can be imported into this resource by supplying Node Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_node.site2_node "/TEST/testnode"
            (or)
$ terraform import bigip_ltm_node.site2_node "/Common/3.3.3.3"

```
