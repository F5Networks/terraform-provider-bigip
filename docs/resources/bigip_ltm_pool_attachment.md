---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_pool_attachment"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_pool_attachment resource
---

# bigip\_ltm\_pool\_attachment

`bigip_ltm_pool_attachment` Manages nodes membership in pools

## Example Usage


There are two ways to use ltm_pool_attachment resource, where we can take node reference from ltm_node or we can specify node directly with ip:port/fqdn:port which will also create node and atach to pool.


### Pool attachment with node directly taking  `ip:port` / `fqdn:port`

```hcl
resource "bigip_ltm_monitor" "monitor" {
  name     = "/Common/terraform_monitor"
  parent   = "/Common/http"
  send     = "GET /some/path\r\n"
  timeout  = "999"
  interval = "998"
}
resource "bigip_ltm_pool" "pool" {
  name                = "/Common/terraform-pool"
  load_balancing_mode = "round-robin"
  monitors            = ["${bigip_ltm_monitor.monitor.name}"]
  allow_snat          = "yes"
  allow_nat           = "yes"
}

resource "bigip_ltm_pool_attachment" "attach_node" {
  pool = bigip_ltm_pool.pool.name
  node = "1.1.1.1:80"
}

```

### Pool attachment with node referenced from `bigip_ltm_node`

```hcl
resource "bigip_ltm_monitor" "monitor" {
  name     = "/Common/terraform_monitor"
  parent   = "/Common/http"
  send     = "GET /some/path\r\n"
  timeout  = "999"
  interval = "998"
}
resource "bigip_ltm_pool" "pool" {
  name                = "/Common/terraform-pool"
  load_balancing_mode = "round-robin"
  monitors            = ["${bigip_ltm_monitor.monitor.name}"]
  allow_snat          = "yes"
  allow_nat           = "yes"
}

resource "bigip_ltm_node" "node" {
  name    = "/Common/terraform_node"
  address = "192.168.30.2"
}

resource "bigip_ltm_pool_attachment" "attach_node" {
  pool = bigip_ltm_pool.pool.name
  node = "${bigip_ltm_node.node.name}:80"
}
```

### Pool attachment resource with attaching multiple nodes in same pool using `for_each`

```hcl
resource "bigip_ltm_node" "node1" {
  name    = "/Common/terraform_node1"
  address = "192.168.30.1"
}
resource "bigip_ltm_node" "node2" {
  name    = "/Common/terraform_node2"
  address = "192.168.30.2"
}
resource "bigip_ltm_pool" "k8s_prod" {
  name = "/Common/k8prod_Pool"
}
resource "bigip_ltm_pool_attachment" "k8sprod" {
  for_each = toset([bigip_ltm_node.node1.name, bigip_ltm_node.node2.name])
  pool     = bigip_ltm_pool.k8s_prod.name
  node     = "${each.key}:80"
}
```


## Argument Reference

* `pool` - (Required) Name of the pool to which members should be attached,it should be "full path".The full path is the combination of the partition + name of the pool.(For example `/Common/my-pool`) or partition + directory + name of the pool (For example `/Common/test/my-pool`).When including directory in fullpath we have to make sure it is created in the given partition before using it.

* `node` - (Required) Pool member address/fqdn with service port, (ex: `1.1.1.1:80/www.google.com:80`). (Note: Member will be in same partition of Pool)

* `connection_limit` - (Optional) Specifies a maximum established connection limit for a pool member or node.The default is 0, meaning that there is no limit to the number of connections.

* `connection_rate_limit` - (Optional) Specifies the maximum number of connections-per-second allowed for a pool member,The default is 0.

* `dynamic_ratio` - (Optional) Specifies the fixed ratio value used for a node during ratio load balancing.

* `ratio`- (Optional) "Specifies the ratio weight to assign to the pool member. Valid values range from 1 through 65535. The default is 1, which means that each pool member has an equal ratio proportion.".

* `priority_group` - (Optional) Specifies a number representing the priority group for the pool member. The default is 0, meaning that the member has no priority

* `monitor` - (Optional) Specifies the health monitors that the system uses to monitor this pool member,value can be `none` (or) `default` (or) list of monitors joined with and ( ex: `/Common/test_monitor_pa_tc1 and /Common/gateway_icmp`).

* `state` - (Optional) Specifies the state the pool member should be in,value can be `enabled` (or) `disabled` (or) `forced_offline`).

* `fqdn_autopopulate` - (Optional) Specifies whether the system automatically creates ephemeral nodes using the IP addresses returned by the resolution of a DNS query for a node defined by an FQDN. The default is enabled

## Importing
An existing pool attachment (i.e. pool membership) can be imported into this resource by supplying both the pool full path, and the node full path with the relevant port. If the pool or node membership is not found, an error will be returned. An example is below:

```sh
$ terraform import bigip_ltm_pool_attachment.node-pool-attach \
	'{"pool": "/Common/terraform-pool", "node": "/Common/node1:80"}'
```
