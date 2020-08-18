---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_http"
sidebar_current: "docs-bigip-resource-profile_http-x"
description: |-
    Provides details about bigip_ltm_profile_http resource
---

# bigip\_ltm\_profile_http

`bigip_ltm_profile_http` Configures a custom profile_http for use by health checks.

For resources should be named with their "full path". The full path is the combination of the partition + name of the resource. For example /Common/my-pool.

## Example Usage


```hcl
resource "bigip_ltm_profile_http" "sanjose-http" {
  name                  = "/Common/sanjose-http"
  defaults_from         = "/Common/http"
  description           = "some http"
  fallback_host         = "titanic"
  fallback_status_codes = ["400", "500", "300"]
}

```      

## Argument Reference

* `name` (Required) Name of the profile_http

* `defaults_from` - (Required) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `fallback_host` - (Optional) Specifies an HTTP fallback host. HTTP redirection allows you to redirect HTTP traffic to another protocol identifier, host name, port number

* `fallback_status_codes` - (Optional) Specifies one or more three-digit status codes that can be returned by an HTTP server.

* `oneconnect_transformations` - (Optional) Enables the system to perform HTTP header transformations for the purpose of  keeping server-side connections open. This feature requires configuration of a OneConnect profile

* `basic_auth_realm` - (Optional) Specifies a quoted string for the basic authentication realm. The system sends this string to a client whenever authorization fails. The default value is none

* `head_insert` - (Optional) Specifies a quoted header string that you want to insert into an HTTP request

* `insert_xforwarded_for` - (Optional) When using connection pooling, which allows clients to make use of other client requests' server-side connections, you can insert the X-Forwarded-For header and specify a client IP address
