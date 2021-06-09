---
layout: "bigip"
page_title: "BIG-IP Provider : Index"
sidebar_current: "docs-bigip-index"
description: |-
    Provides details about provider bigip
---

# F5 BIG-IP Provider

A [Terraform](https://terraform.io) provider for F5 BIG-IP. Resources are currently available for LTM.

### Requirements

This provider uses the iControlREST API. All the resources are validated with BigIP v12.1.1 and above.
## Example

```
provider "bigip" {
  address = "${var.url}"
  username = "${var.username}"
  password = "${var.password}"
}
```

## Reference

- `address` - (Required) Address of the device
- `username` - (Required) Username for authentication
- `password` - (Required) Password for authentication
- `token_value` - (Optional) A token generated outside the provider, in place of password
- `token_auth` - (Optional, Default=false) Enable to use an external authentication source (LDAP, TACACS, etc)
- `login_ref` - (Optional, Default="tmos") Login reference for token authentication (see BIG-IP REST docs for details)
- `port` - (Optional) Management Port to connect to Bigip.This is mainly required if we have single nic bigip or any custom managemt port other than 443
