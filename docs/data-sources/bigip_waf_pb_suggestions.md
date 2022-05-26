---
layout: "bigip"
page_title: "BIG-IP: bigip_waf pb suggestions"
subcategory: "Web Application Firewall(WAF)"
description: |-
  Provides details and exports bigip_waf_pb_suggestions data source
---

# bigip\_waf\_pb_suggestions

Use this data source (`bigip_waf_pb_suggestions`) to export PB suggestions from an existing WAF policy.
 
 
## Example Usage

```hcl

data "bigip_waf_pb_suggestions" "PBWAF1" {
  policy_name            = "protect_me_policy"
  partition              = "Common"
  minimum_learning_score = 20
}

```

## Argument Reference

* `policy_name` - (Required) WAF policy name from which PB suggestions should be exported.
* `partition` - (Required) Partition on which WAF policy is located.
* `minimum_learning_score` - (Required) The minimum learning score for suggestions.


## Attributes Reference

* `policy_id` - System generated id of the WAF policy
* `json` - Json string representing exported PB suggestions ready to be used in WAF policy declaration

