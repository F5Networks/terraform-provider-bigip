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


## Example Usage


```hcl
resource "bigip_ltm_pool_attachment" "node-terraform_pool" {
  pool = "/Common/terraform-pool"
  node = "${bigip_ltm_node.node.name}:80"
}

```      

## Argument Reference

* `pool` - (Required) Name of the pool in /Partition/Name format

* `node` - (Required) Node to add to the pool in /Partition/NodeName:Port format (e.g. /Common/Node01:80)
