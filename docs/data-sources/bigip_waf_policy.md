---
layout: "bigip"
page_title: "BIG-IP: bigip_waf_policy"
subcategory: "Web Application Firewall(WAF)"
description: |-
  Provides details about deployed waf policy using its ID
---

# bigip\_waf\_policy

Use this data source (`bigip_waf_policy`) to get the details of exist WAF policy BIG-IP.
 
## Example Usage
```hcl

data "bigip_waf_policy" "existpolicy" {
  policy_id = "xxxxx"
}

```

## Argument Reference

* `policy_id` - (Required) ID of the WAF policy deployed in the BIG-IP.


## Attributes Reference

* `policy_json` - Exported WAF policy JSON
