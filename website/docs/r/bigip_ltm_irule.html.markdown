---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_irule"
sidebar_current: "docs-bigip-resource-irule-x"
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
  name  = "/Common/terraform_irule"
  irule = file("myirule.tcl")
}

resource "bigip_ltm_irule" "rule2" {
  name  = "/Common/terraform_irule2"
  irule = <<EOF
when CLIENT_ACCEPTED {
     log local0. "test"
   }
EOF

}

```


##myirule.tcl

```
when HTTP_REQUEST {

  if { [string tolower [HTTP::header value Upgrade]] equals "websocket" } {
    HTTP::disable
#    ASM::disable
    log local0. "[IP::client_addr] - Connection upgraded to websocket protocol. Disabling ASM-checks and HTTP protocol. Traffic is treated as L4 TCP stream."
  } else {
    HTTP::enable
#    ASM::enable
    log local0. "[IP::client_addr] - Regular HTTP request. ASM-checks and HTTP protocol enabled. Traffic is deep-inspected at L7."
  }
}
```      

## Argument Reference


* `name` - (Required) Name of the iRule

* `irule` - (Required) Body of the iRule
