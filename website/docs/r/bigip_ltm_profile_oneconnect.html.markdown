---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_oneconnect"
sidebar_current: "docs-bigip-resource-profile_oneconnect-x"
description: |-
    Provides details about bigip_ltm_profile_oneconnect resource
---

# bigip\_ltm\_profile_oneconnect

`bigip_ltm_profile_oneconnect` Configures a custom profile_oneconnect for use by health checks.

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

## Example Usage


```hcl
resource "bigip_ltm_profile_oneconnect" "oneconnect-sanjose" {
  name                  = "sanjose"
  partition             = "Common"
  defaults_from         = "/Common/oneconnect"
  idle_timeout_override = "disabled"
  max_age               = 3600
  max_reuse             = 1000
  max_size              = 1000
  share_pools           = "disabled"
  source_mask           = "255.255.255.255"
}


```      

## Argument Reference

* `name` (Required) Name of the profile_oneconnect

* `partition` - (Optional) Displays the administrative partition within which this profile resides

* `defaults_from` - (Optional) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `idle_timeout_override` - (Optional) Specifies the number of seconds that a connection is idle before the connection flow is eligible for deletion. Possible values are disabled, indefinite, or a numeric value that you specify. The default value is disabled.

* `share_pools` - (Optional) Specify if you want to share the pool, default value is "disabled"

* `max_age` - (Optional) Specifies the maximum age in number of seconds allowed for a connection in the connection reuse pool. For any connection with an age higher than this value, the system removes that connection from the reuse pool. The default value is 86400.

* `max_reuse` - (Optional) Specifies the maximum number of times that a server-side connection can be reused. The default value is 1000.

* `max_size` - (Optional) Specifies the maximum number of connections that the system holds in the connection reuse pool. If the pool is already full, then the server-side connection closes after the response is completed. The default value is 10000.

* `source_mask` - (Optional) Specifies a source IP mask. The default value is 0.0.0.0. The system applies the value of this option to the source address to determine its eligibility for reuse. A mask of 0.0.0.0 causes the system to share reused connections across all clients. A host mask (all 1's in binary), causes the system to share only those reused connections originating from the same client IP address.
