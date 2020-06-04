---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_pool_attachment"
sidebar_current: "docs-bigip-resource-pool-attachment-x"
description: |-
    Provides details about bigip_ltm_pool_attachment resource
---

# bigip\_ltm\_pool\_attachment

`bigip_ltm_pool_attachment` Manages nodes membership in pools

Resources should be named with their "full path". The full path is the combination of the partition + name of the resource. 
For example /Common/my-pool.


## Example Usage

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

## Argument Reference

* `pool` - (Required) Name of the pool, which should be referenced from `bigip_ltm_pool` resource

* `node` - (Required) Name of the Node with service port. (Name of Node should be referenced from `bigip_ltm_node` resource)
