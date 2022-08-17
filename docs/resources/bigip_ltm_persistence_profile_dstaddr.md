---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_persistence_profile_dstaddr"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_persistence_profile_dstaddr resource
---

# bigip_ltm_persistence_profile_dstaddr

Configures a cookie persistence profile

## Example

```
resource "bigip_ltm_persistence_profile_dstaddr" "dstaddr" {
  name                  = "/Common/terraform_ppdstaddr"
  defaults_from         = "/Common/dest_addr"
  match_across_pools    = "enabled"
  match_across_services = "enabled"
  match_across_virtuals = "enabled"
  mirror                = "enabled"
  timeout               = 3600
  override_conn_limit   = "enabled"
  hash_algorithm        = "carp"
  mask                  = "255.255.255.255"
}

```

## Reference

`name` - (Required) Name of the virtual address

`defaults_from` - (Optional) Specifies the existing profile from which the system imports settings for the new profile.

`match_across_pools` (Optional) (enabled or disabled) match across pools with given persistence record

`match_across_services` (Optional) (enabled or disabled) match across services with given persistence record

`match_across_virtuals` (Optional) (enabled or disabled) match across virtual servers with given persistence record

`mirror` (Optional) (enabled or disabled) mirror persistence record

`timeout` (Optional) (enabled or disabled) Timeout for persistence of the session in seconds

`override_conn_limit` (Optional) (enabled or disabled) Enable or dissable pool member connection limits are overridden for persisted clients. Per-virtual connection limits remain hard limits and are not overridden.

## Importing
An dest-addr persistence profile can be imported into this resource by supplying the Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_persistence_profile_dstaddr.dstaddr "/Common/terraform_ppdstaddr"
```
