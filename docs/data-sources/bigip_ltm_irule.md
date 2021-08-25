---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_irule"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_irule data source
---

# bigip\_ltm\_irule

Use this data source (`bigip_ltm_irule`) to get the ltm irule details available on BIG-IP
 
 
## Example Usage
```hcl

data "bigip_ltm_irule" "test" {
  name      = "terraform_irule"
  partition = "Common"
}


output "bigip_irule" {
  value = "${data.bigip_ltm_irule.test.irule}"
}

```      

## Argument Reference

* `name` - (Required) Name of the irule

* `partition` - (Required) partition of the ltm irule


## Attributes Reference

* `irule` - Irule configured on bigip

* `name` - Name of irule configured on bigip with full path

* `partition` - Bigip partition in which rule is configured

