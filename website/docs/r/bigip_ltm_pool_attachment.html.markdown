---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_pool_attachment"
sidebar_current: "docs-bigip-resource-pool-attachment-x"
description: |-
    Provides details about bigip_ltm_pool_attachment resource
---

# bigip\_ltm\_pool\_attachment

`bigip_ltm_pool_attachment` Manages nodes membership in pools

Resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

Note: node must be the full path to the node followed by the port. For example /Common/my-node:80

## Example Usage


```hcl
resource "bigip_ltm_pool_attachment" "node-pool-attach" {
  pool = "/Common/terraform-pool"
  node = "${bigip_ltm_node.node.name}:80"
}

```      

## Argument Reference

* `id` - (Computed) the `id` of the resource is a combination of the pool and node member full path, joined by a hyphen (e.g. "/Common/terraform-pool-/Common/node1:80")
* `pool` - (Required) Name of the pool in /Partition/Name format

* `node` - (Required) Node to add to the pool in /Partition/NodeName:Port format (e.g. /Common/Node01:80)

## Importing

An existing pool attachment (i.e. pool membership) can be imported into this resource by supplying both the pool full path, and the node full path with the relevant port. If the pool or node membership is not found, an error will be returned. An example is below:

```sh
$ terraform import bigip_ltm_pool_attachment.node-pool-attach \
	'{"pool": "/Common/terraform-pool", "node": "/Common/node1:80"}'
