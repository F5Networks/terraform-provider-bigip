---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_irule"
sidebar_current: "docs-bigip-datasource-irule-x"
description: |-
    Provides details about bigip_ltm_irule resource
---

# bigip\_ltm\_irule

`bigip_ltm_irule` Creates iRule on BIG-IP F5 device

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.


## Example Usage


```hcl
# Loading from a file is the preferred method
resource "bigip_ltm_irule" "rule" {
  name = "/Common/terraform_irule"
  irule = "${file("myirule.tcl")}"
}

resource "bigip_ltm_irule" "rule2" {
  name = "/Common/terraform_irule2"
  irule = <<EOF
when CLIENT_ACCEPTED {
     log local0. "test"
   }
EOF
}

```      

## Argument Reference


* `name` - (Required) Name of the iRule

* `irule` - (Required) Body of the iRule
