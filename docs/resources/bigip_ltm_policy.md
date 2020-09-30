---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_policy"
sidebar_current: "docs-bigip-resource-policy-x"
description: |-
    Provides details about bigip_ltm_policy resource
---

# bigip\_ltm\_policy

`bigip_ltm_policy` Configures Virtual Server

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl

resource "bigip_ltm_policy" "test-policy" {
  name           = "my_policy"
  strategy       = "first-match"
  requires       = ["http"]
  published_copy = "Drafts/my_policy"
  controls       = ["forwarding"]
  rule {
    name = "rule6"

    action {
      tm_name = "20"
      forward = true
      pool    = "/Common/mypool"
    }
  }
  depends_on = [bigip_ltm_pool.mypool]
}
```      

## Argument Reference


* `name`- (Required) Name of the Policy

* `strategy` - (Optional) Specifies the match strategy

* `requires` - (Optional) Specifies the protocol

* `published_copy` - (Optional) If you want to publish the policy else it will be deployed in Drafts mode.

*  `controls` - (Optional) Specifies the controls

* `rule` - (Optional) Rules can be applied using the policy

* `tm_name` - (Required) If Rule is used then you need to provide the tm_name it can be any value

* `forward` - (Optional) This action will affect forwarding.

* `pool` - (Optional ) This action will direct the stream to this pool.
