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
}

```      

## Argument Reference

* `name` - (Required) Name of the node

* `address` - (Required) IP or hostname of the node
