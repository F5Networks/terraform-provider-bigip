---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_policy"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_policy data source
---

# bigip\_ltm\_policy

Use this data source (`bigip_ltm_policy`) to get the ltm policy details available on BIG-IP


## Example Usage

```hcl

data "bigip_ltm_policy" "test" {
  name = "/Common/test-policy"
}

output "bigip_policy" {
  value = data.bigip_ltm_policy.test.rule
}


```

## Argument Reference

* `name` - (Required) Name of the policy which includes partion ( /partition/policy-name )


## Attributes Reference

Additionally, the following attributes are exported:
* `name` - The name of the policy.

* `strategy` - Specifies the match strategy.

* `requires` - Specifies the protocol.

* `controls` - Specifies the controls.

* `rule` - Rules defined in the policy.


