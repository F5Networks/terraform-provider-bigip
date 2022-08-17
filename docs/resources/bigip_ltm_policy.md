---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_policy"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_policy resource
---

# bigip\_ltm\_policy

`bigip_ltm_policy` Configures ltm policies to manage traffic assigned to a virtual server

For resources should be named with their `full path`. The full path is the combination of the `partition + name` of the resource. For example `/Common/test-policy`.

## Example Usage

```hcl

resource "bigip_ltm_pool" "mypool" {
  name                = "/Common/test-pool"
  allow_nat           = "yes"
  allow_snat          = "yes"
  load_balancing_mode = "round-robin"
}
resource "bigip_ltm_policy" "test-policy" {
  name     = "/Common/test-policy"
  strategy = "first-match"
  requires = ["http"]
  controls = ["forwarding"]
  rule {
    name = "rule6"
    action {
      forward    = true
      connection = false
      pool       = bigip_ltm_pool.mypool.name
    }
  }
  depends_on = [bigip_ltm_pool.mypool]
}
```      

## Argument Reference

* `name`- (Required) Name of the Policy ( policy name should be in full path which is combination of partition and policy name )

* `strategy` - (Optional) Specifies the match strategy

* `requires` - (Optional) Specifies the protocol

* `published_copy` - (Optional) If you want to publish the policy else it will be deployed in Drafts mode.

*  `controls` - (Optional) Specifies the controls

* `rule` - (Optional) Rules can be applied using the policy

* `forward` - (Optional) This action will affect forwarding.

* `pool` - (Optional ) This action will direct the stream to this pool.

* `connection` - (Optional) This action is set to `true` by default, it needs to be explicitly set to `false` for actions it conflicts with.


## Importing
An existing monitor can be imported into this resource by supplying monitor Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_monitor.monitor /Common/terraform_monitor
```
