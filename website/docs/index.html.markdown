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

- `address` - (Required) Domain name or IP address of the device. May be set via the `BIGIP_HOST` environment variable.
- `username` - (Required) Username for authentication. May be set via the `BIGIP_USER` environment variable.
- `password` - (Required) Password for authentication. May be set via the `BIGIP_PASSWORD` environment variable.
- `token_auth` - (Optional, Default=false) Enable to use an external authentication source (LDAP, TACACS, etc). May be set via the `BIGIP_TOKEN_AUTH` environment variable.
- `login_ref` - (Optional, Default="tmos") Login reference for token authentication (see BIG-IP REST docs for details). May be set via the `BIGIP_LOGIN_REF` environment variable.

### Note
The F5 BIG-IP provider gathers non-identifiable usage data for the purposes of improving the product as outlined in the end user license agreement for BIG-IP. To opt out of data collection, use the following:

`export TEEM_DISABLE=true`
