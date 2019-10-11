---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_persistence_profile_ssl"
sidebar_current: "docs-bigip-resource-persistence_profile_ssl-x"
description: |-
    Provides details about bigip_ltm_persistence_profile_ssl resource
---

# bigip_ltm_persistence_profile_ssl

Configures an SSL persistence profile

## Example

```hcl
resource "bigip_ltm_persistence_profile_ssl" "ppssl" {
  name                  = "/Common/terraform_ssl"
  defaults_from         = "/Common/ssl"
  match_across_pools    = "enabled"
  match_across_services = "enabled"
  match_across_virtuals = "enabled"
  mirror                = "enabled"
  timeout               = 3600
  override_conn_limit   = "enabled"
}
```

## Reference

`name` - (Required) Name of the virtual address

`defaults_from` - (Required) Parent cookie persistence profile

`match_across_pools` (Optional) (enabled or disabled) match across pools with given persistence record

`match_across_services` (Optional) (enabled or disabled) match across services with given persistence record

`match_across_virtuals` (Optional) (enabled or disabled) match across virtual servers with given persistence record

`mirror` (Optional) (enabled or disabled) mirror persistence record

`timeout` (Optional) (enabled or disabled) Timeout for persistence of the session in seconds

`override_conn_limit` (Optional) (enabled or disabled) Enable or dissable pool member connection limits are overridden for persisted clients. Per-virtual connection limits remain hard limits and are not overridden.
