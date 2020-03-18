---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_http2"
sidebar_current: "docs-bigip-resource-profile_http2-x"
description: |-
    Provides details about bigip_ltm_profile_http2 resource
---

# bigip\_ltm\_profile_http2

`bigip_ltm_profile_http2` Configures a custom profile_http2 for use by health checks.

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

## Example Usage


```hcl
resource "bigip_ltm_profile_http2" "nyhttp2" {
  name                              = "/Common/NewYork_http2"
  defaults_from                     = "/Common/http2"
  concurrent_streams_per_connection = 10
  connection_idle_timeout           = 30
  activation_modes                  = ["alpn", "npn"]
}

```      

## Argument Reference

* `name` (Required) Name of the profile_http2

* `defaults_from` - (Required) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `concurrent_streams_per_connection` - (Optional) Specifies how many concurrent requests are allowed to be outstanding on a single HTTP/2 connection.

* `connection_idle_timeout` - (Optional) Specifies the number of seconds that a connection is idle before the connection is eligible for deletion..

* `connpool_maxsize` - (Optional) Specifies the maximum number of connections to a load balancing pool. A setting of 0 specifies that a pool can accept an unlimited number of connections. The default value is 2048.

* `activation_modes` - (Optional) Specifies what will cause an incoming connection to be handled as a HTTP/2 connection. The default values npn and alpn specify that the TLS next-protocol-negotiation and application-layer-protocol-negotiation extensions will be used.
