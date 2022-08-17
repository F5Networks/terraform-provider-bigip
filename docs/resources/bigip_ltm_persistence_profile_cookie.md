---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_persistence_profile_cookie"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_persistence_profile_cookie resource
---

# bigip_ltm_persistence_profile_cookie

Configures a cookie persistence profile

## Example

```hcl
resource "bigip_ltm_persistence_profile_cookie" "test_ppcookie" {
  name                         = "/Common/terraform_cookie"
  defaults_from                = "/Common/cookie"
  match_across_pools           = "enabled"
  match_across_services        = "enabled"
  match_across_virtuals        = "enabled"
  timeout                      = 3600
  override_conn_limit          = "enabled"
  always_send                  = "enabled"
  cookie_encryption            = "required"
  cookie_encryption_passphrase = "iam"
  cookie_name                  = "ham"
  expiration                   = "1:0:0"
  hash_length                  = 0

  lifecycle {
    ignore_changes = [cookie_encryption_passphrase]
  }
}

```

## Reference

`name` - (Required) Name of the virtual address

`defaults_from` - (Required) Parent cookie persistence profile

`match_across_pools` (Optional) (enabled or disabled) match across pools with given persistence record

`match_across_services` (Optional) (enabled or disabled) match across services with given persistence record

`match_across_virtuals` (Optional) (enabled or disabled) match across virtual servers with given persistence record

`method` (Optional) Specifies the type of cookie processing that the system uses. The default value is insert

`mirror` (Optional) (enabled or disabled) mirror persistence record

`timeout` (Optional) (enabled or disabled) Timeout for persistence of the session in seconds

`override_conn_limit` (Optional) (enabled or disabled) Enable or dissable pool member connection limits are overridden for persisted clients. Per-virtual connection limits remain hard limits and are not overridden.

`always_send` (Optional) (enabled or disabled) always send cookies

`cookie_encryption` (Optional) (required, preferred, or disabled) To required, preferred, or disabled policy for cookie encryption

`cookie_encryption_passphrase` (Optional) (required, preferred, or disabled) Passphrase for encrypted cookies. The field is encrypted on the server and will always return differently then set.
If this is configured specify `ignore_changes` under the `lifecycle` block to ignore returned encrypted value.

`cookie_name` (Optional) Name of the cookie to track persistence

`expiration` (Optional) Expiration TTL for cookie specified in DAY:HOUR:MIN:SECONDS (Examples: 1:0:0:0 one day, 1:0:0 one hour, 30:0 thirty minutes)

`hash_length` (Optional) (Integer) Length of hash to apply to cookie

`hash_offset` (Optional) (Integer) Number of characters to skip in the cookie for the hash

`httponly` (Optional) (enabled or disabled) Sending only over http

## Importing
An cookie persistence profile can be imported into this resource by supplying the Name in `full path` as `id`.
An example is below:
```sh
$ terraform import bigip_ltm_persistence_profile_cookie.test_ppcookie "/Common/terraform_cookie"
```
