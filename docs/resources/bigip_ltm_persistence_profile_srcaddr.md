---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_persistence_profile_srcaddr"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_persistence_profile_srcaddr resource
---

# bigip_ltm_persistence_profile_srcaddr

Configures a source address persistence profile

## Example

```
resource "bigip_ltm_persistence_profile_srcaddr" "srcaddr" {
  name                  = "/Common/terraform_srcaddr"
  defaults_from         = "/Common/source_addr"
  match_across_pools    = "enabled"
  match_across_services = "enabled"
  match_across_virtuals = "enabled"
  mirror                = "enabled"
  timeout               = 3600
  override_conn_limit   = "enabled"
  hash_algorithm        = "carp"
  map_proxies           = "enabled"
  mask                  = "255.255.255.255"
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

`hash_algorithm` (Optional) Specify the hash algorithm

`mask` (Optional) Identify a range of source IP addresses to manage together as a single source address affinity persistent connection when connecting to the pool. Must be a valid IPv4 or IPv6 mask.

`map_proxies` (Optional) (enabled or disabled) Directs all to the same single pool member

## Importing
An source-addr persistence profile can be imported into this resource by supplying the Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_persistence_profile_srcaddr.srcaddr "/Common/terraform_srcaddr"
```
